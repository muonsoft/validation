#!/usr/bin/env python3
"""Manual deterministic runner checks; never invoke a real coding model."""

import argparse
import contextlib
import importlib.util
import io
import json
import os
from pathlib import Path
import signal
import sqlite3
import sys
import tempfile
import time
import unittest
from unittest.mock import Mock, patch

SPEC = importlib.util.spec_from_file_location("runner", Path(__file__).with_name("run.py"))
runner = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(runner)


class RunnerTests(unittest.TestCase):
    def setUp(self):
        self.temporary = tempfile.TemporaryDirectory(prefix="validation-runner-test-")
        self.addCleanup(self.temporary.cleanup)
        self.root = Path(self.temporary.name)

    def test_command_success_and_failure(self):
        log = self.root / "success.jsonl"
        result = runner.command([sys.executable, "-c", "print('partial')"], self.root, log)
        self.assertEqual(result["status"], "completed")
        self.assertEqual(log.read_text(), "partial\n")
        result = runner.command([sys.executable, "-c", "import sys; print('bad', file=sys.stderr); sys.exit(7)"], self.root, log)
        self.assertEqual(result["exit_code"], 7)
        self.assertEqual(result["status"], "process_error")
        self.assertIn("bad", log.with_suffix(".stderr").read_text())

    def test_missing_executable(self):
        result = runner.command([str(self.root / "missing")], self.root, self.root / "events")
        self.assertEqual(result["status"], "launch_error")

    def test_timeout_preserves_partial_output(self):
        log = self.root / "partial.jsonl"
        result = runner.command([sys.executable, "-c", "import time; print('partial', flush=True); time.sleep(30)"], self.root, log, timeout=0.2)
        self.assertEqual(result["status"], "timeout")
        self.assertEqual(log.read_text(), "partial\n")

    @unittest.skipUnless(hasattr(signal, "SIGINT"), "POSIX cancellation")
    def test_cancellation_preserves_partial_output(self):
        log = self.root / "cancel.jsonl"
        source = "import os,signal,time; print('partial',flush=True); os.kill(os.getppid(),signal.SIGINT); time.sleep(30)"
        result = runner.command([sys.executable, "-c", source], self.root, log)
        self.assertEqual(result["status"], "cancelled")
        self.assertEqual(log.read_text(), "partial\n")

    def test_event_completion_and_loading_evidence(self):
        log = self.root / "events.jsonl"
        events = [
            {"type": "tool_use", "part": {"tool": "skill", "state": {"status": "completed", "input": {"name": "muonsoft-validation"}}}},
            {"type": "step_finish", "part": {"reason": "stop", "cost": 0.02, "tokens": {"input": 12, "output": 34}}},
        ]
        log.write_text("\n".join(json.dumps(e) for e in events))
        result = runner.events_result(log)
        self.assertTrue(result["complete"])
        self.assertEqual(result["skill_evidence_lines"], [1])
        self.assertEqual(result["reported_usage"][0]["tokens"]["output"], 34)
        log.write_text(json.dumps({"type": "text", "part": {"text": "I used muonsoft-validation"}}))
        result = runner.events_result(log)
        self.assertFalse(result["complete"])
        self.assertEqual(result["skill_loading"], "not_observed")
        for tail in ('\nnot json', '\n{"type":"error"}', '\nnull'):
            log.write_text(json.dumps(events[-1]) + tail)
            self.assertFalse(runner.events_result(log)["complete"])

    def test_named_tests_must_actually_finish(self):
        log = self.root / "tests.jsonl"
        log.write_text('\n'.join(json.dumps(e) for e in [
            {"Action": "pass", "Test": "TestA"},
            {"Action": "skip", "Test": "TestB"},
        ]))
        result = runner.test_result(log, ["TestA", "TestB"], {"status": "completed"})
        self.assertFalse(result["passed"])
        self.assertEqual(result["requirements_passed"], ["TestA"])

    def test_isolation_does_not_inherit_opencode_settings(self):
        with patch.dict(os.environ, {"OPENCODE_CONFIG_CONTENT": '{"instructions":["contamination"]}', "OPENCODE_CONFIG": "/unrelated"}):
            env = runner.isolated_env(self.root, self.root / "provider.json")
        self.assertNotIn("OPENCODE_CONFIG_CONTENT", env)
        self.assertEqual(env["OPENCODE_CONFIG"], str(self.root / "provider.json"))
        self.assertEqual(env["HOME"], os.environ["HOME"])
        self.assertEqual(env["OPENCODE_DISABLE_EXTERNAL_SKILLS"], "true")

    def test_credentials_and_instructions_are_not_inherited(self):
        config = self.root / "config.json"
        for data in ({"provider": {"p": {}}, "instructions": ["unexpected"]},
                     {"provider": {"p": {"options": {"apiKey": "literal-secret"}}}},
                     {"provider": {"p": {"options": {"headers": {"Authorization": "literal-secret"}}}}}):
            config.write_text(json.dumps(data))
            with self.assertRaises(ValueError):
                runner.opencode_config(config)
        config.write_text(json.dumps({"provider": {"p": {"options": {"apiKey": "{env:EVAL_KEY}"}}}}))
        self.assertEqual(runner.opencode_config(config)["share"], "disabled")

    def test_worker_environment_excludes_unrelated_secrets_and_startup_hooks(self):
        unwanted = {"UNRELATED_API_KEY": "private-key", "AGENTMEM_MCP_KEY": "private-key",
                    "CURSOR_API_KEY": "private-key", "BASH_ENV": "/private/startup.sh",
                    "PYTHONPATH": "/private/modules", "LD_PRELOAD": "/private/hook.so"}
        with patch.dict(os.environ, unwanted):
            env = runner.isolated_env(self.root, self.root / "provider.json")
        self.assertFalse(set(unwanted) & env.keys())
        self.assertEqual(env["PATH"], os.environ["PATH"])

    def test_provider_config_inherits_only_referenced_environment_credentials(self):
        config = self.root / "input.json"
        config.write_text('{"provider":{"fixture":{"options":{"apiKey":"{env:EVAL_KEY}"}}}}')
        args = argparse.Namespace(opencode=self.fake_opencode(), config=str(config), model="fixture/model")
        with patch.dict(os.environ, {"EVAL_KEY": "requested-secret", "OTHER_API_KEY": "unrelated-secret"}):
            env, major, _ = runner.prepare_runtime(args, self.root / "output")
        self.assertEqual(major, 1)
        self.assertEqual(env["EVAL_KEY"], "requested-secret")
        self.assertNotIn("OTHER_API_KEY", env)
        self.assertNotIn("requested-secret", (self.root / "output/provider.json").read_text())

    def test_symlink_submission_rejected(self):
        (self.root / "solution.go").symlink_to(self.root.parent / "outside.go")
        with self.assertRaises(ValueError):
            runner.files_digest(self.root)

    def test_output_cannot_overwrite_checkout_or_existing_run(self):
        for path in (runner.ROOT, runner.ROOT / "eval-results", runner.ROOT.parent):
            with self.assertRaises(ValueError):
                runner.new_output(path)
        with self.assertRaises(FileExistsError):
            runner.new_output(self.root)

    def fake_opencode(self, major=1):
        executable = self.root / "fake-opencode"
        executable.write_text('''#!/usr/bin/env python3
import json, os, pathlib, sys
args=sys.argv[1:]
with_skill=pathlib.Path('.opencode/skills/muonsoft-validation/SKILL.md').exists()
if '--version' in args: print('opencode vVERSION.0.0')
elif 'models' in args: print('fixture/model')
elif 'debug' in args: print(json.dumps([{'name':'muonsoft-validation'}] if with_skill else []))
elif 'api' in args:
    assert '--server' in args and '--pure' not in args
    print(json.dumps({'data':[{'id':'opencode','path':'/builtin/opencode.md'}, {'id':'report','path':'/builtin/report.md'}] + ([{'id':'muonsoft-validation','path':str(pathlib.Path('.opencode/skills/muonsoft-validation/SKILL.md').resolve())}] if with_skill else [])}))
elif 'run' in args:
    assert os.environ['PWD'] == str(pathlib.Path.cwd())
    if 'VERSION' == '2': assert '--server' in args and '--pure' not in args
    if with_skill:
        pathlib.Path('solution.go').write_text('package scenario\\n// simulated worker change\\n')
        print(json.dumps({'type':'tool_use','part':{'tool':'skill','state':{'status':'completed','input':{'name':'muonsoft-validation'}}}}))
    print(json.dumps({'type':'step_finish','part':{'reason':'stop'}}))
else: sys.exit(3)
'''.replace('VERSION', str(major)))
        executable.chmod(0o755)
        return str(executable)

    def test_v2_doctor_and_native_credentials(self):
        config = self.root / "provider.json"
        config.write_text('{"providers":{"fixture":{"settings":{"apiKey":"{env:EVAL_KEY}"}}}}')
        args = argparse.Namespace(model="fixture/model", config=str(config), output=str(self.root / "doctor"),
                                  opencode=self.fake_opencode(2))
        with patch.object(runner, 'v2_server', return_value=contextlib.nullcontext(('http://127.0.0.1:1234', {}))), contextlib.redirect_stdout(io.StringIO()):
            self.assertEqual(runner.doctor(args), 0)
        data = json.loads((Path(args.output) / "doctor.json").read_text())
        self.assertEqual(data["opencode_version"], "opencode v2.0.0")
        self.assertEqual(data["inference_requests"], 0)
        native = runner.opencode_config(config, 2)
        self.assertEqual(native["permissions"][0]["effect"], "deny")
        self.assertNotIn("permission", native)
        config.write_text('{"providers":{"fixture":{"models":{"m":{"settings":{"apiKey":"secret"}}}}}}')
        with self.assertRaises(ValueError):
            runner.opencode_config(config, 2)

    def test_v2_skill_evidence_uses_id_and_completed_state(self):
        log = self.root / "events.jsonl"
        tool = {"type": "tool_use", "part": {"tool": "skill", "state": {"status": "completed", "input": {"id": "muonsoft-validation"}}}}
        finish = {"type": "step_finish", "part": {"reason": "stop"}}
        log.write_text(json.dumps(tool) + '\n' + json.dumps(finish))
        self.assertEqual(runner.events_result(log)["skill_loading"], "observed")
        tool["part"]["state"]["status"] = "error"
        log.write_text(json.dumps(tool) + '\n' + json.dumps(finish))
        self.assertEqual(runner.events_result(log)["skill_loading"], "not_observed")

    def test_v2_paired_run_and_report_without_a_model(self):
        fake = self.fake_opencode
        with patch.object(self, 'fake_opencode', side_effect=lambda: fake(2)), patch.object(
                runner, 'v2_server', side_effect=lambda *args: contextlib.nullcontext(('http://127.0.0.1:1234', args[-1]))), patch.object(
                runner, 'v2_completion', return_value={'verified': True}):
            self.test_paired_run_and_report_without_a_model()

    def test_v2_completion_requires_matching_workspace_and_final_text(self):
        database = self.root / 'opencode.db'
        db = sqlite3.connect(database)
        db.execute('CREATE TABLE session_v2 (id TEXT, directory TEXT, idle_outcome TEXT, time_idle INTEGER)')
        db.execute('CREATE TABLE session_message (id TEXT, session_id TEXT, type TEXT, seq INTEGER, data TEXT)')
        db.execute('INSERT INTO session_v2 VALUES (?, ?, ?, ?)', ('ses_fixture', str(self.root), 'succeeded', 2))
        final = {'content': [{'type': 'text', 'text': 'Done.'}], 'time': {'completed': 2}, 'finish': 'stop', 'tokens': {'input': 10, 'output': 5}}
        db.execute('INSERT INTO session_message VALUES (?, ?, ?, ?, ?)', ('msg_fixture', 'ses_fixture', 'assistant', 1, json.dumps(final)))
        db.commit()
        db.close()
        log = self.root / 'events.jsonl'
        event = {'type': 'text', 'sessionID': 'ses_fixture', 'part': {'messageID': 'msg_fixture', 'text': 'Done.'}}
        log.write_text(json.dumps(event))
        self.assertFalse(runner.events_result(log)['complete'])
        completion = runner.v2_completion(log, self.root, database)
        self.assertTrue(completion['verified'])
        self.assertTrue(runner.events_result(log, completion)['complete'])
        self.assertEqual(runner.events_result(log, completion)['reported_usage'][0]['tokens']['output'], 5)
        self.assertFalse(runner.v2_completion(log, self.root / 'wrong', database)['verified'])
        event['part']['text'] = 'Do'
        log.write_text(json.dumps(event))
        self.assertFalse(runner.v2_completion(log, self.root, database)['verified'])
        log.write_text(json.dumps({'type': 'step_finish', 'part': {'reason': 'stop'}}))
        self.assertFalse(runner.events_result(log, {'verified': False})['complete'])

    def test_invalid_auto_config_is_not_written(self):
        args = argparse.Namespace(opencode=self.fake_opencode(2), config=None, model='fixture/model')
        provider = {'providers': {'fixture': {'settings': {'authToken': 'secret'}}}}
        with patch.object(runner, 'local_opencode_config', return_value=('fixture/model', provider)):
            with self.assertRaises(ValueError):
                runner.prepare_runtime(args, self.root)
        self.assertFalse((self.root / 'provider.json').exists())

    def test_v2_server_waits_for_model_and_stops_on_worker_failure(self):
        process = Mock(pid=123456)
        process.poll.return_value = None
        listener = Mock()
        listener.getsockname.return_value = ('127.0.0.1', 1234)
        replies = [io.BytesIO(b'{"data":[]}'), io.BytesIO(b'{"data":[{"providerID":"fixture","id":"model"}]}')]
        with patch.object(runner.socket, 'socket') as socket_factory, patch.object(
                runner.subprocess, 'Popen', return_value=process), patch.object(
                runner.urllib.request, 'urlopen', side_effect=replies) as request, patch.object(
                runner.time, 'sleep'), patch.object(runner.os, 'killpg') as kill, patch.object(
                runner, 'worker_profile', side_effect=lambda env, directory: dict(env)), patch.object(
                runner, 'sandbox_command', return_value=['opencode']):
            socket_factory.return_value.__enter__.return_value = listener
            with self.assertRaisesRegex(RuntimeError, 'worker failure'):
                with runner.v2_server('opencode', self.root, self.root / 'server.log',
                                      {'EVAL_OPENCODE_MODEL': 'fixture/model'}) as (url, env):
                    self.assertEqual(url, 'http://127.0.0.1:1234')
                    self.assertTrue(env['OPENCODE_SERVER_PASSWORD'])
                    raise RuntimeError('worker failure')
            self.assertEqual(request.call_count, 2)
            kill.assert_called_once_with(process.pid, signal.SIGTERM)
            process.wait.assert_called_once_with(timeout=3)

    def test_v2_profiles_do_not_copy_sessions_or_share_tool_caches(self):
        config = self.root / 'provider.json'
        config.write_text('{"providers":{"fixture":{}}}')
        env = runner.isolated_env(self.root, config)
        database = Path(env['XDG_DATA_HOME']) / 'opencode/opencode.db'
        database.parent.mkdir()
        database.write_text('previous session contents')
        first = runner.worker_profile(env, self.root / 'first')
        second = runner.worker_profile(env, self.root / 'second')
        self.assertNotEqual(first['XDG_DATA_HOME'], second['XDG_DATA_HOME'])
        self.assertNotEqual(first['GOCACHE'], second['GOCACHE'])
        self.assertFalse((Path(first['XDG_DATA_HOME']) / 'opencode/opencode.db').exists())
        self.assertEqual(first['HOME'], os.environ['HOME'])

    def test_v2_sandbox_does_not_bind_checkout_or_sibling_attempts(self):
        workspace, profile, library, modules = [self.root / name for name in ('workspace', 'profile', 'libraries', 'modules')]
        for path in (workspace, profile, library, modules):
            path.mkdir()
        env = {'GOMODCACHE': str(modules), 'GOROOT': '/usr/local/go', 'EVAL_LIBRARY_ROOT': str(library)}
        with patch.object(runner.shutil, 'which', side_effect=lambda name: '/usr/bin/' + name):
            argv = runner.sandbox_command('opencode', workspace, profile, env)
        writable = [argv[i + 1] for i, value in enumerate(argv) if value == '--bind']
        self.assertEqual(writable, [str(workspace), str(profile)])
        self.assertNotIn(str(runner.ROOT), argv)
        self.assertNotIn(str(self.root), argv)
        self.assertIn('--unshare-pid', argv)

    def test_existing_console_login_keeps_secrets_out_of_config(self):
        state, data = self.root / "state/opencode", self.root / "data/opencode"
        state.mkdir(parents=True)
        data.mkdir(parents=True)
        (state / "service.json").write_text('{"url":"http://127.0.0.1:1234","password":"local-secret"}')
        (state / "model.json").write_text('{"recent":[{"providerID":"opencode-go","modelID":"deepseek-v4.1-flash"}]}')
        db = sqlite3.connect(data / "opencode.db")
        db.execute('CREATE TABLE credential (integration_id TEXT, active INTEGER, value TEXT)')
        credential = {"type": "oauth", "access": "test-secret", "expires": (time.time() + 3600) * 1000}
        db.execute('INSERT INTO credential VALUES (?, ?, ?)', ('opencode', 1, json.dumps(credential)))
        db.commit()
        db.close()
        models = [{"id": "deepseek-v4.1-flash", "modelID": "deepseek-v4.1-flash", "providerID": p,
                   "enabled": True, "package": "@opencode/ai/providers/openai-compatible",
                   "headers": {"x-opencode-org-id": "private-org"}, "settings": {"baseURL": "https://example.invalid"}}
                  for p in ('opencode', 'opencode-go')]
        response = json.dumps({"data": models}).encode()
        environment = {"XDG_STATE_HOME": str(state.parent), "XDG_DATA_HOME": str(data.parent)}
        with patch.dict(os.environ, environment), patch.object(runner.urllib.request, 'urlopen', return_value=io.BytesIO(response)):
            env = {}
            model, config = runner.local_opencode_config('deepseek-v4.1-flash', env)
        self.assertEqual(model, 'opencode-go/deepseek-v4.1-flash')
        self.assertEqual(env['EVAL_OPENCODE_TOKEN'], 'test-secret')
        self.assertNotIn('test-secret', json.dumps(config))
        self.assertNotIn('private-org', json.dumps(config))
        self.assertNotIn('local-secret', json.dumps(config))
        with patch.dict(os.environ, environment), patch.object(runner.urllib.request, 'urlopen', return_value=io.BytesIO(response)):
            with self.assertRaisesRegex(ValueError, 'no fallback'):
                runner.local_opencode_config('missing-model', {})

    def test_paired_run_and_report_without_a_model(self):
        binary = self.fake_opencode()
        config = self.root / "provider.json"
        config.write_text('{"provider":{"fixture":{"options":{"apiKey":"{env:EVAL_KEY}"}}}}')
        library = self.root / "library"
        library.mkdir()
        (library / "go.mod").write_text("module github.com/muonsoft/validation\n\ngo 1.24.0\n")
        (library / "go.sum").write_text("")
        snapshots = {v: {"path": str(library), "sha256": "fixture"} for v in ("v0.19.0", "current")}
        def fake_grade(case, source, output, snapshots, env):
            passes = "simulated worker" in (source / "solution.go").read_text()
            return {v: {"status": "completed" if passes else "process_error", "passed": passes,
                        "requirements_passed": case["requirements"] if passes else [],
                        "requirements_total": len(case["requirements"])} for v in snapshots}
        args = argparse.Namespace(model="fixture/model", config=str(config), output=str(self.root / "run"),
                                  opencode=binary, case=["01-eager"], profile="smoke", timeout=5)
        with patch.object(runner, "library_snapshots", return_value=snapshots), patch.object(runner, "grade", side_effect=fake_grade), contextlib.redirect_stdout(io.StringIO()):
            self.assertFalse(runner.run(args))
        output = Path(args.output)
        data = json.loads((output / "run.json").read_text())
        self.assertEqual([a["variant"] for a in data["attempts"]], ["without", "with"])
        self.assertEqual(data["attempts"][1]["events"]["skill_loading"], "observed")
        summary = json.loads((output / "results.json").read_text())
        self.assertEqual(summary["pairs"], {"improved": 1})
        self.assertEqual(summary["paired_delta_percentage_points"], 100)
        review = self.root / "review.json"
        review.write_text(json.dumps({data["attempts"][1]["id"]: {key: {"score": 2, "evidence": "[solution](attempts/01-eager-with-1/workspace/solution.go), line 2"} for key in runner.RUBRIC}}))
        runner.report(output, review)
        self.assertIn("**2/2**", (output / "report.md").read_text())
        self.assertEqual(json.loads((output / "results.json").read_text())["pairs"], summary["pairs"])
        # A cancellation remains ungraded and cannot turn into a false success.
        data["attempts"][1].pop("grades")
        data["attempts"][1]["status"] = "cancelled"
        data["status"] = "cancelled"
        runner.write_json(output / "run.json", data)
        summary = runner.report(output, None)
        self.assertEqual(summary["pairs"], {"incomplete": 1})
        self.assertEqual(summary["variants"]["with"]["ungraded"], {"cancelled": 1})
        with self.assertRaises(ValueError):
            runner.report(output, review)

        # Cancellation during execution stops the pair and persists its report.
        real_command = runner.command
        def cancel_worker(argv, cwd, log, env=None, timeout=180):
            if "run" not in argv:
                return real_command(argv, cwd, log, env, timeout)
            log.write_text("")
            return {"status": "cancelled", "exit_code": None, "seconds": 0.1}
        args.output = str(self.root / "cancelled-run")
        with patch.object(runner, "library_snapshots", return_value=snapshots), patch.object(runner, "command", side_effect=cancel_worker), patch.object(runner, "grade") as grade, contextlib.redirect_stdout(io.StringIO()):
            self.assertTrue(runner.run(args))
            grade.assert_not_called()
        cancelled = json.loads((Path(args.output) / "run.json").read_text())
        self.assertEqual(cancelled["status"], "cancelled")
        self.assertEqual(len(cancelled["attempts"]), 1)
        self.assertEqual(cancelled["attempts"][0]["status"], "cancelled")

        # Cancellation in Go checks must not start testing the next version.
        def cancel_tests(argv, cwd, log, env=None, timeout=180):
            log.write_text("")
            return {"status": "cancelled", "exit_code": None}
        with patch.object(runner, "command", side_effect=cancel_tests) as execute:
            grades = runner.grade(runner.cases()[0], output / "attempts/01-eager-with-1/workspace", self.root / "cancelled-checks", snapshots, {})
            self.assertEqual(execute.call_count, 1)
        self.assertFalse(grades["v0.19.0"]["passed"])

    def test_token_pairs_exclude_unknown_and_execution_failures(self):
        def attempt(case, variant, records, status="completed", passed=True):
            return {"id": f"{case}-{variant}-1", "case": case, "repeat": 1, "variant": variant,
                    "status": status, "should_trigger": True,
                    "grades": {version: {"passed": passed} for version in ("v0.19.0", "current")},
                    "events": {"complete": True, "reported_usage": records}}
        def record(value):
            return {"tokens": {"input": value, "output": 2, "reasoning": 3, "cache": {"read": 10, "write": 0}}}
        attempts = [attempt("01-eager", "without", [record(10), record(20)]),
                    attempt("01-eager", "with", [record(10)], passed=False),
                    attempt("02-collection", "without", [record(10)]),
                    attempt("02-collection", "with", [{"tokens": {"input": 1}}]),
                    attempt("03-optional", "without", [record(10)]),
                    attempt("03-optional", "with", [record(10)], status="timeout")]
        result = runner.token_comparison(attempts)
        self.assertEqual(result["matched"]["pairs"], 1)
        self.assertEqual(result["matched"]["totals"]["without"]["all_categories"], 60)
        self.assertEqual(result["matched"]["totals"]["with"]["all_categories"], 25)
        self.assertEqual(result["matched"]["totals"]["without"]["uncached_plus_generated"], 40)
        self.assertEqual(result["matched_both_passed"]["pairs"], 0)
        self.assertEqual(len(result["excluded_pairs"]), 2)
        self.assertIsNone(result["attempts"]["02-collection-with-1"]["tokens"]["all_categories"])
        self.assertFalse(result["attempts"]["02-collection-with-1"]["complete"])
        with self.assertRaisesRegex(ValueError, "duplicate"):
            runner.token_comparison(attempts + [attempts[0]])

    def test_final_usage_is_not_counted_twice(self):
        log = self.root / "events.jsonl"
        tokens = {"input": 10, "output": 5, "reasoning": 2, "cache": {"read": 8, "write": 0}}
        log.write_text(json.dumps({"type": "step_finish", "part": {"reason": "stop", "tokens": tokens}}))
        completion = {"verified": True, "finish_logged": True, "tokens": tokens}
        self.assertEqual(len(runner.events_result(log, completion)["reported_usage"]), 1)
        log.write_text(json.dumps({"type": "text", "part": {"text": "done"}}))
        completion["finish_logged"] = False
        self.assertEqual(runner.events_result(log, completion)["reported_usage"], [{"tokens": tokens}])

    def test_missing_step_usage_and_duplicate_events_are_not_complete_usage(self):
        log = self.root / "events.jsonl"
        tokens = {"input": 10, "output": 5, "reasoning": 2, "cache": {"read": 8, "write": 0}}
        first = {"type": "step_finish", "part": {"id": "step1", "reason": "tool-calls", "tokens": tokens}}
        last = {"type": "step_finish", "part": {"id": "step2", "reason": "stop"}}
        log.write_text(json.dumps(first) + "\n" + json.dumps(last))
        events = runner.events_result(log)
        self.assertTrue(events["complete"])
        self.assertFalse(runner.token_usage({"events": events})["complete"])
        last["part"]["tokens"] = tokens
        log.write_text("\n".join(json.dumps(e) for e in (first, first, last)))
        events = runner.events_result(log)
        self.assertFalse(events["complete"])
        self.assertEqual(events["error_lines"], [2])

    def test_resume_preserves_attempts_and_checks_provenance(self):
        source = self.root / "original"
        original = source / "attempts/01-eager-without-1"
        original.mkdir(parents=True)
        (original / "events.jsonl").write_text("partial output")
        manifest = {key: "same" for key in ("skill_sha256", "suite_sha256", "runner_sha256", "cases_sha256",
                    "go_version", "model", "profile", "timeout", "opencode_version", "config_sha256")}
        manifest.update(repeats=1, selected_cases=["01-eager"], libraries={"current": {"sha256": "lib"}}, attempts=[])
        previous = {**manifest, "attempts": [{"id": "01-eager-without-1", "case": "01-eager", "variant": "without",
                                             "repeat": 1, "status": "running"}]}
        runner.write_json(source / "run.json", previous)
        original_bytes = (source / "run.json").read_bytes()
        output = self.root / "continued"
        output.mkdir()
        seen = runner.inherit_run(source, output, manifest)
        self.assertEqual(seen, {"01-eager-without-1"})
        self.assertEqual(manifest["attempts"][0]["status"], "interrupted")
        self.assertEqual((source / "run.json").read_bytes(), original_bytes)
        self.assertEqual((output / "attempts/01-eager-without-1/events.jsonl").read_text(), "partial output")
        for key in ("skill_sha256", "suite_sha256", "runner_sha256", "cases_sha256", "go_version", "model",
                    "profile", "timeout", "opencode_version", "config_sha256", "selected_cases", "repeats"):
            with self.subTest(key=key), self.assertRaisesRegex(ValueError, key):
                runner.inherit_run(source, output, {**manifest, key: "different"})
        with self.assertRaisesRegex(ValueError, "library snapshots"):
            runner.inherit_run(source, output, {**manifest, "libraries": {"current": {"sha256": "changed"}}})

    def test_active_run_cannot_be_resumed(self):
        path = self.root / "runner.lock"
        with runner.run_lock(path, create=True):
            with self.assertRaisesRegex(ValueError, "active"):
                with runner.run_lock(path):
                    self.fail("acquired active run")
        with runner.run_lock(path):
            pass

    def test_resume_executes_only_unstarted_attempts(self):
        config = self.root / "provider.json"
        config.write_text('{"provider":{"fixture":{}}}')
        library = self.root / "library"
        library.mkdir()
        (library / "go.mod").write_text("module github.com/muonsoft/validation\n\ngo 1.24.0\n")
        (library / "go.sum").write_text("")
        snapshots = {v: {"path": str(library), "sha256": "fixture"} for v in ("v0.19.0", "current")}
        args = argparse.Namespace(model="fixture/model", config=str(config), output=str(self.root / "first"),
                                  opencode=self.fake_opencode(), case=["01-eager"], profile="smoke", timeout=5)
        execute = runner.command
        calls = []
        def cancel(argv, cwd, log, env=None, timeout=180):
            if "run" not in argv:
                return execute(argv, cwd, log, env, timeout)
            calls.append(cwd)
            log.write_text("")
            return {"status": "cancelled"}
        with patch.object(runner, "library_snapshots", return_value=snapshots), patch.object(
                runner, "command", side_effect=cancel), contextlib.redirect_stdout(io.StringIO()):
            self.assertTrue(runner.run(args))
            source = Path(args.output)
            before = (source / "run.json").read_bytes()
            args.resume = source
            args.output = str(self.root / "second")
            self.assertTrue(runner.run(args))
        self.assertEqual(len(calls), 2)
        self.assertIn("01-eager-without-1", str(calls[0]))
        self.assertIn("01-eager-with-1", str(calls[1]))
        self.assertEqual((source / "run.json").read_bytes(), before)
        manifest = json.loads((Path(args.output) / "run.json").read_text())
        self.assertEqual(len(manifest["attempts"]), 2)
        self.assertEqual(manifest["continuation"]["inherited_attempts"], ["01-eager-without-1"])

    def test_doctor_never_invokes_inference(self):
        config = self.root / "provider.json"
        config.write_text('{"provider":{"fixture":{}}}')
        args = argparse.Namespace(model="fixture/model", config=str(config), output=str(self.root / "doctor"), opencode=self.fake_opencode())
        with contextlib.redirect_stdout(io.StringIO()):
            self.assertEqual(runner.doctor(args), 0)
        data = json.loads((Path(args.output) / "doctor.json").read_text())
        self.assertEqual(data["inference_requests"], 0)
        args.model = "fixture/unavailable"
        args.output = str(self.root / "missing-model")
        with self.assertRaisesRegex(ValueError, "no fallback"):
            runner.doctor(args)


if __name__ == "__main__":
    unittest.main()
