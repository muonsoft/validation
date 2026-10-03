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
import sys
import tempfile
import unittest
from unittest.mock import patch

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

    def fake_opencode(self):
        executable = self.root / "fake-opencode"
        executable.write_text('''#!/usr/bin/env python3
import json, pathlib, sys
args=sys.argv[1:]
with_skill=pathlib.Path('.opencode/skills/muonsoft-validation/SKILL.md').exists()
if '--version' in args: print('fake-no-model')
elif 'models' in args: print('fixture/model')
elif 'debug' in args: print(json.dumps([{'name':'muonsoft-validation'}] if with_skill else []))
elif 'run' in args:
    if with_skill:
        pathlib.Path('solution.go').write_text('package scenario\\n// simulated worker change\\n')
        print(json.dumps({'type':'tool_use','part':{'tool':'skill','state':{'status':'completed','input':{'name':'muonsoft-validation'}}}}))
    print(json.dumps({'type':'step_finish','part':{'reason':'stop'}}))
else: sys.exit(3)
''')
        executable.chmod(0o755)
        return str(executable)

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
