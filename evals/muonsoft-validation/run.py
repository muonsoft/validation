#!/usr/bin/env python3
"""Manual skill evals. Python standard library only; never called by CI."""

import argparse
import base64
from collections import Counter
from contextlib import contextmanager
import difflib
import fcntl
import gzip
import hashlib
import io
import json
import math
import os
from pathlib import Path
import re
import secrets
import shutil
import signal
import socket
import sqlite3
import statistics
import subprocess
import sys
import tarfile
import time
import urllib.parse
import urllib.request

HERE = Path(__file__).resolve().parent
ROOT = HERE.parent.parent
SKILL = ROOT / "skills/muonsoft-validation"
MODULE = "example.com/validation-eval/scenario"
RUBRIC = ("api", "rule_ownership", "execution_flow", "scope")


def write_json(path, value):
    path.parent.mkdir(parents=True, exist_ok=True)
    temporary = path.with_suffix(path.suffix + ".tmp")
    temporary.write_text(json.dumps(value, ensure_ascii=False, indent=2) + "\n")
    temporary.replace(path)


def digest(data):
    return hashlib.sha256(data).hexdigest()


def files_digest(directory):
    files = {}
    for path in sorted(directory.rglob("*")):
        if ".git" in path.relative_to(directory).parts:
            continue
        if path.is_symlink():
            raise ValueError(f"symlink is not an eval artifact: {path}")
        if path.is_file():
            files[path.relative_to(directory).as_posix()] = digest(path.read_bytes())
    return files


def case_manifest(suite):
    if suite not in ("main", "legacy", "transfer"):
        raise ValueError("unknown suite")
    return HERE / ("cases.json" if suite == "main" else f"{suite}-cases.json")


def cases(suite="main"):
    data = json.loads(case_manifest(suite).read_text())
    if data["schema_version"] != 1:
        raise ValueError("unsupported case schema")
    ids = set()
    for case in data["cases"]:
        identifier = case["id"]
        if not re.fullmatch(r"[0-9]{2}-[a-z]+", identifier) or identifier in ids:
            raise ValueError(f"invalid/duplicate case ID: {identifier}")
        ids.add(identifier)
        path = HERE / "testdata" / identifier
        for name in ("prompt.md", "workspace/contract.go", "checks/requirements_test.go",
                     *(f"{folder}/{name}" for folder in ("workspace", "reference") for name in case["editable"])):
            if not (path / name).is_file():
                raise ValueError(f"missing {identifier}/{name}")
        tests = re.findall(r"^func (Test\w+)\(", (path / "checks/requirements_test.go").read_text(), re.M)
        if tests != case["requirements"] or not tests:
            raise ValueError(f"test manifest mismatch: {identifier}")
    return data["cases"]


def editable(case, name):
    original = HERE / "testdata" / case["id"] / "workspace" / name
    return name in case["editable"] or (case.get("allow_helpers", False)
        and Path(name).name == name and name.endswith(".go")
        and not name.endswith("_test.go") and not original.exists())


def command(argv, cwd, output, env=None, timeout=180):
    """Keep partial logs and stop the entire process group on timeout/cancellation."""
    output.parent.mkdir(parents=True, exist_ok=True)
    started = time.monotonic()
    result = {"status": "running", "exit_code": None}
    with output.open("w") as stdout, output.with_suffix(".stderr").open("w") as stderr:
        try:
            process = subprocess.Popen(argv, cwd=cwd, env=env, stdout=stdout, stderr=stderr,
                                       start_new_session=True)
        except OSError as exc:
            stderr.write(str(exc))
            result["status"] = "launch_error"
        else:
            try:
                result["exit_code"] = process.wait(timeout=timeout)
                result["status"] = "completed" if process.returncode == 0 else "process_error"
            except (subprocess.TimeoutExpired, KeyboardInterrupt) as exc:
                try:
                    os.killpg(process.pid, signal.SIGTERM)
                except ProcessLookupError:
                    pass
                try:
                    process.wait(timeout=3)
                except subprocess.TimeoutExpired:
                    os.killpg(process.pid, signal.SIGKILL)
                    process.wait()
                result["status"] = "cancelled" if isinstance(exc, KeyboardInterrupt) else "timeout"
    result["seconds"] = round(time.monotonic() - started, 3)
    return result


def go_env(output):
    return {**os.environ, "GOWORK": "off", "GOTOOLCHAIN": "local",
            "GOCACHE": str(output / "go-cache")}


def library_snapshots(output):
    """Copy Go implementation only: the worker must not discover eval answers."""
    snapshots = {}
    archive = subprocess.run(["git", "archive", "v0.19.0"], cwd=ROOT,
                             capture_output=True, check=True).stdout
    with tarfile.open(fileobj=io.BytesIO(archive)) as tar:
        release = {entry.name: tar.extractfile(entry).read() for entry in tar.getmembers()
                   if entry.isfile() and library_file(entry.name)}
    tracked = subprocess.run(["git", "ls-files", "-z"], cwd=ROOT,
                             capture_output=True, check=True).stdout.decode().split("\0")
    current = {name: (ROOT / name).read_bytes() for name in tracked
               if name and library_file(name) and (ROOT / name).is_file()}
    # Include newly added library implementation files, but never skills or eval fixtures.
    for path in ROOT.rglob("*.go"):
        name = path.relative_to(ROOT).as_posix()
        if library_file(name):
            current[name] = path.read_bytes()
    for version, contents in (("v0.19.0", release), ("current", current)):
        target = output / "libraries" / version
        target.mkdir(parents=True)
        for name, data in contents.items():
            destination = target / name
            destination.parent.mkdir(parents=True, exist_ok=True)
            destination.write_bytes(data)
        snapshots[version] = {"path": str(target), "sha256": digest(json.dumps(files_digest(target), sort_keys=True).encode())}
    return snapshots


def library_file(name):
    parts = Path(name).parts
    return (not any(p.startswith(".") or p in {"evals", "skills", "test", "testdata", "vendor"} for p in parts)
            and (name in {"go.mod", "go.sum"} or name.endswith(".go") and not name.endswith("_test.go")))


def prepare_module(workspace, library, negative=False):
    workspace.mkdir(parents=True, exist_ok=True)
    if negative:
        (workspace / "go.mod").write_text(f"module {MODULE}\n\ngo 1.24.0\n")
        return
    text = (library / "go.mod").read_text().replace("module github.com/muonsoft/validation", f"module {MODULE}", 1)
    text += '\nrequire github.com/muonsoft/validation v0.19.0\n'
    # Go accepts quoted replacement directories, including paths with spaces.
    text += 'replace github.com/muonsoft/validation => ' + json.dumps(str(library)) + '\n'
    (workspace / "go.mod").write_text(text)
    shutil.copyfile(library / "go.sum", workspace / "go.sum")


def workspace(case, target, library, reference=False):
    shutil.copytree(HERE / "testdata" / case["id"] / "workspace", target)
    if reference:
        shutil.copytree(HERE / "testdata" / case["id"] / "reference", target, dirs_exist_ok=True)
    prepare_module(target, library, negative=not case["should_trigger"])


def test_result(path, requirements, process):
    completed = {}
    for line in path.read_text().splitlines():
        try:
            event = json.loads(line)
        except json.JSONDecodeError:
            continue
        if event.get("Test") in requirements and event.get("Action") in ("pass", "fail", "skip"):
            completed[event["Test"]] = event["Action"]
    passed = [name for name in requirements if completed.get(name) == "pass"]
    return {**process, "requirements_passed": passed,
            "requirements_failed": [name for name in requirements if completed.get(name) == "fail"],
            "requirements_total": len(requirements),
            "passed": process["status"] == "completed" and len(passed) == len(requirements)}


def grade(case, source, output, snapshots, env):
    results = {}
    for name in case["editable"]:
        if not (source / name).is_file() or (source / name).is_symlink():
            raise ValueError(f"missing or symlink implementation: {name}")
    for version, library in snapshots.items():
        target = output / version
        shutil.copytree(HERE / "testdata" / case["id"] / "workspace", target)
        # Restore protected files; include all allowed implementation files, never worker tests.
        for path in source.glob("*.go"):
            if editable(case, path.name):
                if path.is_symlink():
                    raise ValueError("symlink submission")
                shutil.copyfile(path, target / path.name)
        prepare_module(target, Path(library["path"]), negative=not case["should_trigger"])
        shutil.copyfile(HERE / "testdata" / case["id"] / "checks/requirements_test.go", target / "requirements_test.go")
        if case["should_trigger"]:
            shutil.copyfile(HERE / "testdata/helpers_test.go", target / "helpers_test.go")
            if case.get("allow_helpers"):
                shutil.copyfile(HERE / "testdata/behavior_test.go", target / "behavior_test.go")
        log = target / "tests.jsonl"
        process = command(["go", "test", "-mod=mod", "-race", "-count=1", "-json", "./..."], target, log, env)
        results[version] = test_result(log, case["requirements"], process)
        if process["status"] == "cancelled":
            break
    return results


def events_result(path, completion=None):
    parsed, malformed, stopped, errors = 0, 0, False, []
    evidence, usage = [], []
    step_ids = set()
    for number, line in enumerate(path.read_text().splitlines(), 1):
        if not line.strip():
            continue
        try:
            event = json.loads(line)
        except json.JSONDecodeError:
            malformed += 1
            continue
        if not isinstance(event, dict):
            malformed += 1
            continue
        parsed += 1
        if event.get("type") == "error":
            errors.append(number)
        part = event.get("part") or {}
        if not isinstance(part, dict):
            continue
        if event.get("type") == "step_finish":
            stopped = part.get("reason") == "stop"
            step_id = part.get("id")
            if step_id is not None:
                if step_id in step_ids:
                    errors.append(number)
                step_ids.add(step_id)
            # Retain a placeholder when a finished step omitted usage.
            usage.append({key: part[key] for key in ("tokens", "cost") if key in part})
        state = part.get("state") or {}
        if not isinstance(state, dict):
            continue
        args = state.get("input") or {}
        if not isinstance(args, dict) or state.get("status") != "completed":
            continue
        if part.get("tool") == "skill" and (args.get("name") or args.get("id")) == "muonsoft-validation":
            evidence.append(number)
        if part.get("tool") == "read" and str(args.get("filePath") or args.get("path", "")).endswith("/muonsoft-validation/SKILL.md"):
            evidence.append(number)
    if completion is not None:
        stopped = completion.get("verified") is True
        if stopped and not completion.get("finish_logged"):
            usage.append({key: completion[key] for key in ("tokens", "cost") if key in completion})
    return {"complete": bool(parsed and stopped and not malformed and not errors),
            "malformed_lines": malformed, "error_lines": errors,
            "skill_loading": "observed" if evidence else "not_observed",
            "skill_evidence_lines": evidence, "reported_usage": usage}


def v2_completion(log, workspace, database):
    """V2 may reconcile final text without emitting step_finish; verify persisted completion."""
    events = []
    for line in log.read_text().splitlines():
        if not line.strip():
            continue
        try:
            event = json.loads(line)
        except ValueError:
            return {"verified": False, "reason": "malformed event log"}
        if not isinstance(event, dict):
            return {"verified": False, "reason": "invalid event"}
        events.append(event)
    sessions = {e["sessionID"] for e in events if e.get("sessionID")}
    if len(sessions) != 1:
        return {"verified": False, "reason": "expected one logged session"}
    session_id = sessions.pop()
    try:
        connection = sqlite3.connect(database.as_uri() + "?mode=ro", uri=True)
        try:
            session = connection.execute("SELECT directory, idle_outcome, time_idle FROM session_v2 WHERE id = ?",
                                         (session_id,)).fetchone()
            final = connection.execute("SELECT id, data FROM session_message WHERE session_id = ? AND type = 'assistant' ORDER BY seq DESC LIMIT 1",
                                       (session_id,)).fetchone()
        finally:
            connection.close()
    except sqlite3.Error:
        return {"verified": False, "reason": "v2 completion database unavailable or unsupported"}
    if not session or not final:
        return {"verified": False, "reason": "missing persisted session or assistant message"}
    message_id, encoded = final
    message = json.loads(encoded)
    expected = "".join(item["text"] for item in message.get("content", []) if item.get("type") == "text")
    logged = "".join(e["part"].get("text", "") for e in events if e.get("type") == "text"
                     and e.get("part", {}).get("messageID") == message_id)
    verified = (Path(session[0]).resolve() == workspace.resolve() and session[1] == "succeeded"
                and bool(session[2]) and message.get("finish") == "stop"
                and bool(message.get("time", {}).get("completed")) and bool(expected) and logged == expected)
    return {"verified": verified, "source": "opencode-v2-session-db", "session_id": session_id,
            "directory": session[0], "outcome": session[1], "time_idle": session[2],
            "message_id": message_id, "finish": message.get("finish"),
            "text_sha256": digest(expected.encode()), "logged_text_sha256": digest(logged.encode()),
            "finish_logged": any(e.get("type") == "step_finish" and e.get("part", {}).get("messageID") == message_id for e in events),
            **{key: message[key] for key in ("tokens", "cost") if key in message}}


def check_materials(output, snapshots):
    """Validate standalone distribution and execute every complete Go fence twice."""
    import importlib.util
    spec = importlib.util.spec_from_file_location("docs", ROOT / "scripts/check-docs.py")
    docs = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(docs)
    errors, count = [], 0
    copied = output / "installed-skill/muonsoft-validation"
    shutil.copytree(SKILL, copied)
    docs.ROOT = copied
    entry = (copied / "SKILL.md").read_text()
    if not entry.startswith("---\nname: muonsoft-validation\n") or "\ndescription: " not in entry:
        errors.append("invalid public skill frontmatter")
    for path in sorted(copied.rglob("*.md")):
        text = path.read_text()
        errors.extend(docs.check_links(path, text))
        for link in list(docs.LINK.finditer(docs.prose(text))) + list(docs.REFERENCE.finditer(docs.prose(text))):
            url = docs.urlsplit(link.group(1).strip("<>"))
            if not url.scheme and not url.netloc:
                target = (path.parent / docs.unquote(url.path)).resolve()
                if not target.is_relative_to(copied.resolve()):
                    errors.append(f"link escapes installed skill: {path.name}: {link.group(1)}")
        for block in docs.FENCE.finditer(text):
            source = block.group(2)
            if block.group(1).strip() != "go" or not re.search(r"^package main\s*$", source, re.M):
                continue
            count += 1
            for version, snapshot in snapshots.items():
                directory = output / "examples" / str(count) / version
                prepare_module(directory, Path(snapshot["path"]))
                (directory / "main.go").write_text(source)
                log = directory / "output.txt"
                result = command(["go", "run", "-mod=mod", "."], directory, log, go_env(output))
                expected = re.search(r"^// Output:\n((?://[^\n]*\n?)*)\s*$", source, re.M)
                if result["status"] != "completed":
                    errors.append(f"{path.name}: example {count} failed on {version}; see {log}")
                elif expected:
                    want = "\n".join(re.sub(r"^// ?", "", line) for line in expected.group(1).splitlines()) + "\n"
                    if log.read_text() != want:
                        errors.append(f"{path.name}: output mismatch on {version}: {log.read_text()!r}")
    return {"examples": count, "errors": errors}


def check(args):
    output = new_output(args.output)
    snapshots = library_snapshots(output)
    report = {"materials": check_materials(output, snapshots), "cases": {}}
    for case in cases(getattr(args, "suite", "main")):
        directory = output / case["id"]
        for mode in ("reference", "starter"):
            src = directory / mode / "workspace"
            workspace(case, src, Path(snapshots["v0.19.0"]["path"]), reference=mode == "reference")
            grades = grade(case, src, directory / mode / "checks", snapshots, go_env(output))
            report["cases"].setdefault(case["id"], {})[mode] = grades
        for mutation in case.get("mutations", []):
            src = directory / mutation["id"] / "workspace"
            workspace(case, src, Path(snapshots["v0.19.0"]["path"]), reference=True)
            path = src / mutation["file"]
            text = path.read_text()
            if text.count(mutation["old"]) != 1:
                raise ValueError(f"ambiguous mutation: {case['id']}/{mutation['id']}")
            path.write_text(text.replace(mutation["old"], mutation["new"]))
            report["cases"][case["id"]][mutation["id"]] = grade(
                case, src, directory / mutation["id"] / "checks", snapshots, go_env(output))
        print(f"checked {case['id']}", flush=True)
    write_json(output / "check.json", report)
    failures = list(report["materials"]["errors"])
    for identifier, modes in report["cases"].items():
        for version in snapshots:
            if not modes["reference"][version]["passed"]:
                failures.append(f"reference failed: {identifier}/{version}")
            if not modes["starter"][version]["requirements_failed"]:
                failures.append(f"starter did not execute failing behavioral tests: {identifier}/{version}")
            if modes["starter"][version]["passed"]:
                failures.append(f"checks did not reject starter: {identifier}/{version}")
    for case in cases(getattr(args, "suite", "main")):
        for mutation in case.get("mutations", []):
            for version in snapshots:
                result = report["cases"][case["id"]][mutation["id"]][version]
                if mutation["test"] not in result["requirements_failed"]:
                    failures.append(f"mutation did not fail its behavioral test: {case['id']}/{mutation['id']}/{version}")
    if getattr(args, "retention", "compact") == "compact":
        shutil.rmtree(output / "go-cache", ignore_errors=True)
    print("\n".join(failures) if failures else "Material, compatibility, and negative checks passed.")
    print(output / "check.json")
    return bool(failures)


def new_output(path):
    output = Path(path).expanduser().resolve()
    if output.is_relative_to(ROOT) or ROOT.is_relative_to(output):
        raise ValueError("output must be outside the library checkout and not its ancestor")
    output.mkdir(mode=0o700, parents=True, exist_ok=False)
    return output


def isolated_env(output, config):
    # Shell tool output is model-visible. Never inherit unrelated service secrets
    # or shell startup hooks from the orchestrator's environment.
    allowed = {"PATH", "HOME", "USER", "LOGNAME", "SHELL", "LANG", "LC_ALL", "LC_CTYPE",
               "TERM", "TMPDIR", "SSL_CERT_FILE", "SSL_CERT_DIR", "NODE_EXTRA_CA_CERTS",
               "HTTP_PROXY", "HTTPS_PROXY", "ALL_PROXY", "NO_PROXY",
               "http_proxy", "https_proxy", "all_proxy", "no_proxy"}
    env = {key: value for key, value in os.environ.items() if key in allowed}
    env.update(GOWORK="off", GOTOOLCHAIN="local", GOCACHE=str(output / "go-cache"))
    for key, directory in (("XDG_CONFIG_HOME", "config"), ("XDG_DATA_HOME", "data"),
                           ("XDG_CACHE_HOME", "cache"), ("XDG_STATE_HOME", "state")):
        location = output / "profile" / directory
        location.mkdir(parents=True, exist_ok=True)
        env[key] = str(location)
    env.update(OPENCODE_CONFIG=str(config), OPENCODE_CONFIG_DIR=str(output / "profile/config/opencode"),
               OPENCODE_DISABLE_CLAUDE_CODE="true", OPENCODE_DISABLE_EXTERNAL_SKILLS="true",
               OPENCODE_DISABLE_AUTOUPDATE="true", OPENCODE_TEST_MANAGED_CONFIG_DIR=str(output / "profile/managed"))
    return env


def opencode_config(path, major=1):
    return configure_provider(json.loads(path.read_text()), major)


def configure_provider(config, major):
    # Only provider setup is inherited, never prompts, plugins, MCP, skills, or agents.
    if set(config) - {"$schema", "provider", "providers"}:
        raise ValueError("eval provider config may contain only $schema and provider/providers")
    if not config.get("provider") and not config.get("providers"):
        raise ValueError("provider config is empty; use explicit provider setup with env credentials")
    # Authentication must remain in environment, not generated artifacts.
    def check_credentials(value):
        if not isinstance(value, dict):
            return
        for key, item in value.items():
            if key in {"apiKey", "authToken", "accessToken"} and item:
                if not isinstance(item, str) or not re.fullmatch(r"\{env:[A-Z_][A-Z_0-9]*\}", item):
                    raise ValueError("use an {env:VARIABLE} credential reference, not a literal credential")
            elif key == "headers":
                for header, content in item.items():
                    if not isinstance(content, str) or not re.fullmatch(r"\{env:[A-Z_][A-Z_0-9]*\}", content):
                        raise ValueError(f"custom header {header} must use an environment reference")
            elif isinstance(item, dict):
                check_credentials(item)
            elif isinstance(item, list):
                for child in item:
                    check_credentials(child)
    check_credentials(config)
    if major == 2:
        config.update(update="disable", share="disabled", instructions=[], plugins=[],
                      mcp={"servers": {}}, snapshots=False,
                      permissions=[{"action": "*", "resource": "*", "effect": "deny"}]
                      + [{"action": action, "resource": "*", "effect": "allow"}
                         for action in ("read", "glob", "grep", "edit", "shell")]
                      + [{"action": "skill", "resource": "muonsoft-validation", "effect": "allow"}])
        return config
    config.update(autoupdate=False, share="disabled", instructions=[], plugin=[], mcp={},
                  permission={"*": "deny", "read": "allow", "glob": "allow", "grep": "allow",
                              "edit": "allow", "bash": "allow", "skill": {"*": "deny", "muonsoft-validation": "allow"},
                              "external_directory": "deny"})
    return config


def local_opencode_config(model, env):
    """Reuse an active Console connection without copying login storage or user settings."""
    state = Path(os.environ.get("XDG_STATE_HOME", Path.home() / ".local/state")) / "opencode"
    data = Path(os.environ.get("XDG_DATA_HOME", Path.home() / ".local/share")) / "opencode"
    service = json.loads((state / "service.json").read_text())
    url = urllib.parse.urlsplit(service["url"])
    if url.scheme != "http" or url.hostname not in {"127.0.0.1", "localhost", "::1"}:
        raise ValueError("automatic configuration requires a local OpenCode service")
    authorization = base64.b64encode(("opencode:" + service["password"]).encode()).decode()
    request = urllib.request.Request(service["url"] + "/api/model",
                                     headers={"Authorization": "Basic " + authorization})
    with urllib.request.urlopen(request, timeout=30) as response:
        models = json.load(response)["data"]
    base, separator, variant = model.partition("#")
    candidates = [m for m in models if m.get("enabled") and
                  (f"{m['providerID']}/{m['id']}" == base or m["id"] == base)]
    if len(candidates) > 1 and "/" not in base:
        preferences = json.loads((state / "model.json").read_text()) if (state / "model.json").exists() else {}
        for recent in preferences.get("recent", []):
            matches = [m for m in candidates if m["providerID"] == recent.get("providerID")
                       and m["id"] == recent.get("modelID")]
            if matches:
                candidates = matches
                break
    if len(candidates) != 1:
        raise ValueError(f"model must resolve uniquely in the local catalog: {model}; specify provider/model; no fallback")
    selected = candidates[0]
    provider = selected["providerID"]
    if provider not in {"opencode", "opencode-go"}:
        raise ValueError("automatic login reuse currently supports OpenCode Console; use --config for other providers")
    if separator and variant not in {v["id"] for v in selected.get("variants", [])}:
        raise ValueError(f"unknown model variant: {variant}")
    connection = sqlite3.connect((data / "opencode.db").as_uri() + "?mode=ro", uri=True)
    try:
        rows = connection.execute("SELECT value FROM credential WHERE integration_id = ? AND active = 1",
                                  ("opencode",)).fetchall()
    finally:
        connection.close()
    if len(rows) != 1:
        raise ValueError("expected one active OpenCode Console credential; use --config otherwise")
    credential = json.loads(rows[0][0])
    if credential["type"] == "oauth":
        if credential.get("expires", 0) <= time.time() * 1000:
            raise ValueError("OpenCode login token expired; refresh the connection in OpenCode and retry")
        token = credential["access"]
    elif credential["type"] == "key":
        token = credential["key"]
    else:
        raise ValueError("unsupported OpenCode credential type; use --config")
    env["EVAL_OPENCODE_TOKEN"] = token
    # Copy only the selected model's runtime definition, never configuration documents.
    definition = {key: selected[key] for key in
                  ("modelID", "name", "package", "settings", "capabilities", "limit", "compatibility", "variants", "cost")
                  if key in selected}
    definition.setdefault("settings", {})["apiKey"] = "{env:EVAL_OPENCODE_TOKEN}"
    headers = {}
    for index, (key, value) in enumerate(selected.get("headers", {}).items()):
        variable = f"EVAL_OPENCODE_HEADER_{index}"
        env[variable] = value
        headers[key] = "{env:" + variable + "}"
    definition["headers"] = headers
    resolved = f"{provider}/{selected['id']}" + ("#" + variant if separator else "")
    return resolved, {"providers": {provider: {"models": {selected["id"]: definition},
                                               "package": selected["package"],
                                               "settings": definition["settings"]}}}


def prepare_runtime(args, output):
    config_path = output / "provider.json"
    env = isolated_env(output, config_path)
    version = subprocess.run([args.opencode, "--version"], env=env, capture_output=True,
                             text=True, check=True).stdout.strip()
    match = re.search(r"(?:^|\s)v?([12])\.\d+\.\d+", version)
    if not match:
        raise ValueError(f"unsupported OpenCode version: {version}")
    major = int(match[1])
    if args.config:
        config = opencode_config(Path(args.config), major)
        # Inherit only variables explicitly referenced by the provider config.
        for key in re.findall(r"\{env:([A-Z_][A-Z_0-9]*)\}", json.dumps(config)):
            if key in os.environ:
                env[key] = os.environ[key]
    else:
        if major != 2:
            raise ValueError("OpenCode v1 requires --config")
        args.model, provider = local_opencode_config(args.model, env)
        config = configure_provider(provider, major)
    write_json(config_path, config)
    if major == 2:
        write_json(Path(env["OPENCODE_CONFIG_DIR"]) / "opencode.json", config)
        # V2 discovers ~/.claude/skills independently of the V1 disable flags.
        home = output / "profile/home"
        home.mkdir()
        env["OPENCODE_TEST_HOME"] = str(home)
        paths = json.loads(subprocess.run(["go", "env", "-json", "GOMODCACHE", "GOROOT"],
                                          capture_output=True, text=True, check=True).stdout)
        env.update({key: value for key, value in paths.items()})
        env["EVAL_LIBRARY_ROOT"] = str(output / "libraries")
    if "/" not in args.model:
        raise ValueError("--model must be provider/model when using --config")
    env["EVAL_OPENCODE_MODEL"] = args.model.split("#", 1)[0]
    return env, major, version


def worker_profile(env, directory):
    """Each V2 server has its own session database and tool cache."""
    directory.mkdir(parents=True)
    config_path = directory / "provider.json"
    config = json.loads(Path(env["OPENCODE_CONFIG"]).read_text())
    write_json(config_path, config)
    result = dict(env)
    for key, name in (("XDG_CONFIG_HOME", "config"), ("XDG_DATA_HOME", "data"),
                      ("XDG_CACHE_HOME", "cache"), ("XDG_STATE_HOME", "state")):
        target = directory / name
        target.mkdir()
        result[key] = str(target)
    home = directory / "home"
    home.mkdir()
    result.update(OPENCODE_CONFIG=str(config_path), OPENCODE_CONFIG_DIR=str(directory / "config/opencode"),
                  OPENCODE_TEST_HOME=str(home), GOCACHE=str(directory / "go-cache"),
                  GOPATH=str(directory / "go"))
    write_json(Path(result["OPENCODE_CONFIG_DIR"]) / "opencode.json", config)
    return result


def sandbox_command(binary, cwd, profile, env):
    """Expose no checkout, evaluator artifacts, other attempts, or real home directory."""
    bubblewrap = shutil.which("bwrap")
    if not bubblewrap:
        raise ValueError("OpenCode v2 eval requires bubblewrap (bwrap) for filesystem isolation")
    arguments = [bubblewrap, "--unshare-user", "--unshare-pid", "--die-with-parent", "--new-session"]
    roots = [Path(p) for p in ("/usr", "/bin", "/lib", "/lib64") if Path(p).exists()]
    for root in roots:
        arguments += ["--ro-bind", str(root), str(root)]
    for name in ("/etc/ssl", "/etc/resolv.conf", "/etc/hosts", "/etc/nsswitch.conf", "/etc/passwd", "/etc/group"):
        if Path(name).exists():
            arguments += ["--ro-bind", name, name]
    arguments += ["--proc", "/proc", "--dev", "/dev", "--tmpfs", "/tmp"]
    executable = Path(shutil.which(binary) or binary).resolve()
    readonly = [executable, Path(env["GOMODCACHE"]), Path(env["GOROOT"]), Path(env["EVAL_LIBRARY_ROOT"])]
    for key in ("SSL_CERT_FILE", "NODE_EXTRA_CA_CERTS"):
        if env.get(key):
            readonly.append(Path(env[key]))
    for path in readonly:
        if path.exists() and not any(path.is_relative_to(root) for root in roots):
            arguments += ["--ro-bind", str(path), str(path)]
    for path in (cwd, profile):
        arguments += ["--bind", str(path), str(path)]
    return arguments + ["--chdir", str(cwd), "--", str(executable)]


@contextmanager
def v2_server(binary, cwd, log, env):
    """Wait for asynchronous catalogs and own the server for this operation only."""
    log.parent.mkdir(parents=True, exist_ok=True)
    password = secrets.token_urlsafe(32)
    profile = log.with_suffix(".profile")
    server_env = {**worker_profile(env, profile), "OPENCODE_SERVER_PASSWORD": password, "PWD": str(cwd)}
    with socket.socket() as listener:
        listener.bind(("127.0.0.1", 0))
        port = listener.getsockname()[1]
    url = f"http://127.0.0.1:{port}"
    authorization = base64.b64encode(("opencode:" + password).encode()).decode()
    query = urllib.parse.urlencode({"location[directory]": str(cwd)})
    request = urllib.request.Request(url + "/api/model?" + query,
                                     headers={"Authorization": "Basic " + authorization})
    with log.open("w") as stream:
        process = subprocess.Popen(sandbox_command(binary, cwd, profile, server_env) +
                                   ["serve", "--hostname", "127.0.0.1", "--port", str(port)],
                                   cwd=cwd, env=server_env, stdout=stream, stderr=stream, start_new_session=True)
        try:
            deadline = time.monotonic() + 45
            while time.monotonic() < deadline and process.poll() is None:
                try:
                    with urllib.request.urlopen(request, timeout=3) as response:
                        models = json.load(response)["data"]
                    if any(f"{m['providerID']}/{m['id']}" == env["EVAL_OPENCODE_MODEL"] for m in models):
                        break
                except (OSError, ValueError, KeyError):
                    pass
                time.sleep(0.25)
            else:
                raise ValueError(f"v2 model catalog did not become ready; see {log}; no fallback")
            yield url, server_env
        finally:
            if process.poll() is None:
                os.killpg(process.pid, signal.SIGTERM)
                try:
                    process.wait(timeout=3)
                except subprocess.TimeoutExpired:
                    os.killpg(process.pid, signal.SIGKILL)
                    process.wait()


def opencode_command(binary, arguments, cwd, log, env, major, timeout=180):
    # V2's run command prefers PWD over process.cwd(), even with --server.
    env = {**env, "PWD": str(cwd)}
    if major == 1:
        return command([binary, "--pure"] + arguments, cwd, log, env, timeout)
    with v2_server(binary, cwd, log.with_suffix(".server.log"), env) as (url, server_env):
        result = command([binary, arguments[0], "--server", url] + arguments[1:], cwd, log, server_env, timeout)
        if arguments[0] == "run" and result["status"] == "completed":
            result["completion"] = v2_completion(log, cwd, Path(server_env["XDG_DATA_HOME"]) / "opencode/opencode.db")
        return result


def discover(binary, cwd, output, env, expected, major=1):
    log = output / "discovery.json"
    arguments = (["api", "GET", "/api/skill?" +
                  urllib.parse.urlencode({"location[directory]": str(cwd)})] if major == 2
                 else ["debug", "skill"])
    result = opencode_command(binary, arguments, cwd, log, env, major)
    if result["status"] != "completed":
        raise ValueError(f"cannot inspect available skills; see {log}")
    try:
        entries = json.loads(log.read_text())
        if major == 2:
            entries = entries["data"]
            builtins = {"opencode": "/builtin/opencode.md", "report": "/builtin/report.md"}
            names = sorted(entry["id"] for entry in entries
                           if builtins.get(entry["id"]) != entry["path"])
        else:
            names = sorted(entry["name"] for entry in entries)
    except (ValueError, KeyError, TypeError) as exc:
        raise ValueError(f"unsupported OpenCode skill discovery output: {log}") from exc
    if names != expected:
        raise ValueError(f"unexpected skills {names}; expected {expected}. Use a clean environment.")


def model_preflight(args, output, env, major=1):
    catalog = output / "models.txt"
    arguments = ["models"] if major == 2 else ["models", args.model.split("/", 1)[0]]
    status = opencode_command(args.opencode, arguments, output, catalog, env, major)
    if status["status"] != "completed" or args.model.split("#", 1)[0] not in catalog.read_text().splitlines():
        raise ValueError(f"requested model unavailable: {args.model}; see {catalog}; no fallback")


def doctor(args):
    output = new_output(args.output)
    env, major, version = prepare_runtime(args, output)
    model_preflight(args, output, env, major)
    for variant in ("without", "with"):
        directory = output / variant
        directory.mkdir()
        subprocess.run(["git", "init", "-q", str(directory)], check=True)
        if variant == "with":
            shutil.copytree(SKILL, directory / ".opencode/skills/muonsoft-validation")
        discover(args.opencode, directory, output / (variant + "-logs"), env,
                 ["muonsoft-validation"] if variant == "with" else [], major)
    write_json(output / "doctor.json", {"model": args.model, "opencode_version": version,
                                        "discovery": "passed", "inference_requests": 0})
    print("Model ID and isolated skill discovery verified; no inference requests made.")
    return 0


def provenance(snapshots):
    return {"schema_version": 1, "created_at": time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime()),
            "revision": subprocess.run(["git", "rev-parse", "HEAD"], cwd=ROOT, capture_output=True, text=True, check=True).stdout.strip(),
            "skill_sha256": digest(json.dumps(files_digest(SKILL), sort_keys=True).encode()),
            "suite_sha256": digest(json.dumps(files_digest(HERE / "testdata"), sort_keys=True).encode()),
            "runner_sha256": digest(Path(__file__).read_bytes()),
            "cases_sha256": digest((HERE / "cases.json").read_bytes()),
            "libraries": snapshots,
            "go_version": subprocess.run(["go", "version"], capture_output=True, text=True, check=True).stdout.strip()}


TOKEN_KEYS = ("input", "output", "reasoning", "cache_read", "cache_write")


def token_usage(attempt):
    """Missing counters are unknown, including on otherwise successful attempts."""
    records = attempt.get("events", {}).get("reported_usage", [])
    totals = dict.fromkeys(TOKEN_KEYS, 0)
    known = dict.fromkeys(TOKEN_KEYS, bool(records))
    for record in records:
        tokens = record.get("tokens") or {}
        cache = tokens.get("cache") or {}
        values = {**{key: tokens.get(key) for key in TOKEN_KEYS[:3]},
                  "cache_read": cache.get("read"), "cache_write": cache.get("write")}
        for key, value in values.items():
            if type(value) not in (int, float) or not math.isfinite(value) or value < 0:
                known[key] = False
            else:
                totals[key] += value
    totals = {key: value if known[key] else None for key, value in totals.items()}
    for key, components in (("uncached_plus_generated", TOKEN_KEYS[:3]), ("all_categories", TOKEN_KEYS)):
        totals[key] = sum(totals[k] for k in components) if all(known[k] for k in components) else None
    return {"complete": attempt.get("events", {}).get("complete") is True and all(known.values()),
            "records": len(records), "tokens": totals}


def token_comparison(attempts):
    """Compare only matched, graded pairs with complete normalized usage."""
    pairs, excluded, matched = {}, [], []
    usage = {a["id"]: token_usage(a) for a in attempts}
    for attempt in attempts:
        pair = pairs.setdefault((attempt["case"], attempt["repeat"]), {})
        if attempt["variant"] in pair:
            raise ValueError("duplicate case/repeat/variant in token comparison")
        pair[attempt["variant"]] = attempt
    for (case, repeat), pair in pairs.items():
        if set(pair) != {"with", "without"} or any(
                set(a.get("grades", {})) != {"v0.19.0", "current"} or a["status"] != "completed"
                or not usage[a["id"]]["complete"] for a in pair.values()):
            excluded.append({"case": case, "repeat": repeat, "reason": "missing graded variant or complete usage"})
        else:
            matched.append(pair)

    def summarize(group):
        keys = (*TOKEN_KEYS, "uncached_plus_generated", "all_categories")
        totals = {variant: {key: sum(usage[pair[variant]["id"]]["tokens"][key] for pair in group)
                            for key in keys} for variant in ("without", "with")}
        changes = {key: 100 * (totals["with"][key] / totals["without"][key] - 1)
                   if totals["without"][key] else None for key in keys}
        percent = [100 * (usage[pair["with"]["id"]]["tokens"]["all_categories"] /
                         usage[pair["without"]["id"]]["tokens"]["all_categories"] - 1)
                   for pair in group if usage[pair["without"]["id"]]["tokens"]["all_categories"]]
        return {"pairs": len(group), "totals": totals, "change_percent": changes,
                "median_per_attempt": {variant: statistics.median(
                    usage[pair[variant]["id"]]["tokens"]["all_categories"] for pair in group) if group else None
                    for variant in ("without", "with")},
                "median_paired_percent_change": statistics.median(percent) if percent else None,
                "pairs_with_fewer_tokens": sum(usage[pair["with"]["id"]]["tokens"]["all_categories"] <
                                               usage[pair["without"]["id"]]["tokens"]["all_categories"] for pair in group)}

    return {"method": "Normalized per-step counters; verified final completion is included only if not logged. "
                       "Matched graded pairs only; failed tests included, execution failures excluded symmetrically. "
                       "Missing usage is unknown. Token volume is not currency cost.",
            "matched": summarize(matched), "excluded_pairs": excluded,
            "by_case": {case: summarize([p for p in matched if p["with"]["case"] == case])
                        for case in sorted({a["case"] for a in attempts})},
            "matched_applicable": summarize([p for p in matched if p["with"]["should_trigger"]]),
            "matched_both_passed": summarize([p for p in matched if all(
                g["passed"] for a in p.values() for g in a["grades"].values())]),
            "attempts": usage}


@contextmanager
def run_lock(path, create=False):
    """A separate inode survives atomic manifest replacement and process death."""
    with path.open("a" if create else "r") as stream:
        try:
            fcntl.flock(stream, fcntl.LOCK_EX | fcntl.LOCK_NB)
        except BlockingIOError as exc:
            raise ValueError("run is active or already being continued") from exc
        try:
            yield
        finally:
            fcntl.flock(stream, fcntl.LOCK_UN)


def compress_log(path):
    """Lossless compression; remove the source only after validating the archive."""
    if path.is_symlink() or not path.is_file():
        return
    target = path.with_name(path.name + ".gz")
    temporary = target.with_name(target.name + ".tmp")
    checksum = hashlib.sha256()
    with path.open("rb") as source, gzip.open(temporary, "wb") as archive:
        while block := source.read(1024 * 1024):
            checksum.update(block)
            archive.write(block)
    verified = hashlib.sha256()
    with gzip.open(temporary, "rb") as archive:
        while block := archive.read(1024 * 1024):
            verified.update(block)
    if verified.digest() != checksum.digest():
        raise ValueError(f"compressed log verification failed: {path}")
    temporary.replace(target)
    path.unlink()


def disposable_paths(output):
    """Only runner-owned caches/profiles; never follow inherited attempt links."""
    paths = [output / name for name in ("go-cache", "profile", "models.server.profile")]
    attempts = output / "attempts"
    if attempts.is_dir() and not attempts.is_symlink():
        for attempt in attempts.iterdir():
            if not attempt.is_dir() or attempt.is_symlink():
                continue
            paths.extend(attempt / name for name in ("discovery.server.profile", "events.server.profile"))
            workspace = attempt / "workspace"
            if workspace.is_dir() and not workspace.is_symlink():
                paths.extend(workspace / name for name in (".git", ".opencode"))
    return [path for path in paths if path.exists() or path.is_symlink()]


def disk_size(path):
    if path.is_symlink():
        return path.lstat().st_size
    if path.is_file():
        return path.stat().st_size
    return sum(file.stat().st_size for file in path.rglob("*") if file.is_file() and not file.is_symlink())


def compact(output, attempt=None):
    paths = disposable_paths(output)
    if attempt is not None:
        paths = [path for path in paths if path.is_relative_to(attempt)]
    removed = sum(disk_size(path) for path in paths)
    for path in paths:
        # Old runs did not snapshot the skill separately. Preserve their only
        # copy before deleting the workspace installation/runtime directory.
        if path.name == ".opencode" and path.is_dir() and not path.is_symlink() and not (output / "materials/skill").is_dir():
            skill = path / "skills/muonsoft-validation"
            saved = path.parent.parent / "installed-skill"
            if skill.is_dir() and not skill.is_symlink() and not saved.exists():
                shutil.copytree(skill, saved, symlinks=True)
        if path.is_symlink() or path.is_file():
            path.unlink()
        else:
            shutil.rmtree(path)
    attempts = output / "attempts"
    roots = []
    if attempts.is_dir() and not attempts.is_symlink():
        roots = [attempt] if attempt is not None else [p for p in attempts.iterdir() if p.is_dir() and not p.is_symlink()]
    for root in roots:
        for path in root.rglob("*"):
            # A submitted symlink is invalid, and must never expose outside logs.
            if any(parent.is_symlink() for parent in (path, *path.parents) if parent.is_relative_to(root)):
                continue
            if path.suffix in (".jsonl", ".stderr", ".log"):
                compress_log(path)
    return removed


def cleanup(args):
    output = Path(args.run).expanduser().resolve()
    if output.is_relative_to(ROOT) or ROOT.is_relative_to(output):
        raise ValueError("cleanup requires a run outside the checkout")
    with run_lock(output / "runner.lock"):
        manifest = json.loads((output / "run.json").read_text())
        if "attempts" not in manifest or "model" not in manifest:
            raise ValueError("not an eval run")
        paths = disposable_paths(output)
        print(json.dumps({"apply": args.apply, "disposable_bytes": sum(disk_size(p) for p in paths),
                          "paths": [str(p.relative_to(output)) for p in paths]}, indent=2))
        if args.apply:
            compact(output)
            if manifest.get("status") == "running":
                manifest["status"] = "interrupted"
                for attempt in manifest["attempts"]:
                    if attempt["status"] in ("running", "preparing"):
                        attempt["status"] = "interrupted"
                write_json(output / "run.json", manifest)
            report(output, output / "review.json" if (output / "review.json").is_file() else None)
    return 0


def snapshot_materials(output, selected, suite):
    target = output / "materials"
    target.mkdir()
    shutil.copyfile(case_manifest(suite), target / "cases.json")
    shutil.copyfile(Path(__file__), target / "run.py")
    for name in ("README.md", "SUITE.md"):
        shutil.copyfile(HERE / name, target / name)
    shutil.copytree(SKILL, target / "skill")
    for case in selected:
        shutil.copytree(HERE / "testdata" / case["id"], target / "testdata" / case["id"])
    for name in ("helpers_test.go", "behavior_test.go"):
        shutil.copyfile(HERE / "testdata" / name, target / "testdata" / name)
    write_json(target / "hashes.json", files_digest(target))


def inherit_run(source, output, manifest):
    """Carry attempts forward without rewriting or retrying their artifacts."""
    raw = (source / "run.json").read_bytes()
    previous = json.loads(raw)
    keys = ("skill_sha256", "suite_sha256", "runner_sha256", "cases_sha256", "go_version",
            "model", "profile", "timeout", "repeats", "opencode_version", "config_sha256", "selected_cases")
    for key in keys:
        if key not in previous or previous[key] != manifest[key]:
            raise ValueError(f"cannot resume: {key} differs or is absent")
    if {k: v["sha256"] for k, v in previous["libraries"].items()} != {
            k: v["sha256"] for k, v in manifest["libraries"].items()}:
        raise ValueError("cannot resume: library snapshots differ")
    expected = {f"{case}-{variant}-{repeat}" for case in manifest["selected_cases"]
                for variant in ("without", "with") for repeat in range(1, manifest["repeats"] + 1)}
    seen = set()
    for attempt in previous["attempts"]:
        identifier = attempt["id"]
        if (identifier not in expected or identifier in seen or identifier !=
                f"{attempt['case']}-{attempt['variant']}-{attempt['repeat']}"):
            raise ValueError("cannot resume: unexpected or duplicate attempt")
        if not (source / "attempts" / identifier).is_dir():
            raise ValueError(f"cannot resume: missing artifacts for {identifier}")
        seen.add(identifier)
    (output / "attempts").mkdir()
    for attempt in previous["attempts"]:
        identifier = attempt["id"]
        # Continuations remain portable when the original run is moved or removed.
        shutil.copytree((source / "attempts" / identifier).resolve(), output / "attempts" / identifier, symlinks=True,
                        ignore=shutil.ignore_patterns("*.profile", "go-cache", ".git", ".opencode"))
        if attempt["status"] in ("running", "preparing"):
            attempt["status"] = "interrupted"
        manifest["attempts"].append(attempt)
    manifest["continuation"] = {"source": str(source), "source_manifest_sha256": digest(raw),
                                "inherited_attempts": sorted(seen), "policy": "recorded attempts are never retried"}
    if (source / "review.json").is_file():
        shutil.copyfile(source / "review.json", output / "review.json")
    return seen


def run(args):
    output = new_output(args.output)
    source = Path(args.resume).expanduser().resolve() if getattr(args, "resume", None) else None
    with run_lock(output / "runner.lock", create=True):
        if not getattr(args, "skip_check", False):
            if check(argparse.Namespace(output=output / "preparation", suite=getattr(args, "suite", "main"), retention="compact")):
                raise ValueError("material verification failed; no model attempts started")
        if source is not None:
            with run_lock(source / "runner.lock"):
                return run_attempts(args, output, source)
        return run_attempts(args, output, None)


def run_attempts(args, output, source):
    suite = getattr(args, "suite", "main")
    retention = getattr(args, "retention", "compact")
    selected = [case for case in cases(suite) if (not args.case or case["id"] in args.case)
                and (args.profile == "full" or case["smoke"] or args.case)]
    if args.case and set(args.case) - {case["id"] for case in selected}:
        raise ValueError("unknown case ID")
    env, major, version = prepare_runtime(args, output)
    snapshots = library_snapshots(output)
    manifest = provenance(snapshots)
    manifest["cases_sha256"] = digest(case_manifest(suite).read_bytes())
    manifest.update(model=args.model, profile=args.profile, timeout=args.timeout,
                    suite=suite, retention=retention,
                    repeats=1 if args.profile == "smoke" else 3, attempts=[], status="running",
                    opencode_version=version,
                    config_sha256=digest((output / "provider.json").read_bytes()),
                    selected_cases=[case["id"] for case in selected],
                    case_requirements={case["id"]: case["requirements"] for case in selected},
                    case_categories={case["id"]: case.get("categories", {}) for case in selected})
    snapshot_materials(output, selected, suite)
    inherited = inherit_run(source, output, manifest) if source else set()
    write_json(output / "run.json", manifest)
    try:
        model_preflight(args, output, env, major)
        for case in selected:
            for repeat in range(manifest["repeats"]):
                # Alternate order to reduce systematic ordering effects.
                variants = ("without", "with") if repeat % 2 == 0 else ("with", "without")
                for variant in variants:
                    identifier = f"{case['id']}-{variant}-{repeat + 1}"
                    if identifier in inherited:
                        continue
                    directory = output / "attempts" / identifier
                    directory.mkdir(parents=True)
                    attempt = {"id": identifier, "case": case["id"], "variant": variant,
                               "repeat": repeat + 1, "should_trigger": case["should_trigger"], "status": "preparing"}
                    manifest["attempts"].append(attempt)
                    write_json(output / "run.json", manifest)
                    src = directory / "workspace"
                    workspace(case, src, Path(snapshots["v0.19.0"]["path"]))
                    subprocess.run(["git", "init", "-q", str(src)], check=True)
                    if variant == "with":
                        shutil.copytree(SKILL, src / ".opencode/skills/muonsoft-validation")
                    discover(args.opencode, src, directory, env, ["muonsoft-validation"] if variant == "with" else [], major)
                    before = files_digest(src)
                    prompt = (HERE / "testdata" / case["id"] / "prompt.md").read_text()
                    (directory / "prompt.md").write_text(prompt)
                    attempt["status"] = "running"
                    write_json(output / "run.json", manifest)
                    result = opencode_command(args.opencode, ["run", "--agent", "build", "--model", args.model,
                                              "--format", "json", prompt], src, directory / "events.jsonl",
                                              env, major, args.timeout)
                    completion = result.get("completion")
                    if completion is not None:
                        write_json(directory / "completion.json", completion)
                    events = events_result(directory / "events.jsonl", completion)
                    attempt.update(result)
                    attempt["events"] = events
                    if attempt["status"] == "completed" and not events["complete"]:
                        attempt["status"] = "incomplete_log"
                    try:
                        after = files_digest(src)
                        changed = sorted(name for name in before.keys() | after.keys() if before.get(name) != after.get(name))
                        unexpected = [name for name in changed if not editable(case, name)]
                        attempt["unexpected_changes"] = unexpected
                        attempt["submission_hashes"] = after
                        chunks = []
                        for name in changed:
                            if not name.endswith(".go"):
                                continue
                            original_path = HERE / "testdata" / case["id"] / "workspace" / name
                            original = original_path.read_text() if original_path.is_file() else ""
                            solution = (src / name).read_text() if (src / name).is_file() else ""
                            chunks.extend(difflib.unified_diff(original.splitlines(True), solution.splitlines(True), fromfile=f"before/{name}", tofile=f"after/{name}"))
                        (directory / "solution.diff").write_text("".join(chunks))
                        if unexpected:
                            attempt["status"] = "invalid_submission"
                        elif attempt["status"] == "completed":
                            attempt["grades"] = grade(case, src, directory / "checks", snapshots, go_env(output))
                    except (OSError, ValueError) as exc:
                        if attempt["status"] == "completed":
                            attempt["status"] = "invalid_submission"
                        attempt["submission_error"] = str(exc)
                    if result["status"] == "cancelled" or any(g["status"] == "cancelled" for g in attempt.get("grades", {}).values()):
                        attempt["status"] = "cancelled"
                        raise KeyboardInterrupt
                    write_json(output / "run.json", manifest)
                    if retention == "compact":
                        compact(output, directory)
                    report(output, None)
                    print(identifier, attempt["status"], flush=True)
        manifest["status"] = "completed"
    except KeyboardInterrupt:
        manifest["status"] = "cancelled"
        if manifest["attempts"] and manifest["attempts"][-1]["status"] in ("preparing", "running"):
            manifest["attempts"][-1]["status"] = "cancelled"
    except Exception:
        manifest["status"] = "interrupted"
        if manifest["attempts"] and manifest["attempts"][-1]["status"] == "preparing":
            manifest["attempts"][-1]["status"] = "preflight_error"
        raise
    finally:
        write_json(output / "run.json", manifest)
        if retention == "compact":
            compact(output)
        review = output / "review.json"
        report(output, review if review.is_file() else None)
    return manifest["status"] != "completed"


def report(output, review_path):
    manifest = json.loads((output / "run.json").read_text())
    attempts = manifest["attempts"]
    reviews = json.loads(review_path.read_text()) if review_path else {}
    for identifier, review in reviews.items():
        if identifier not in {attempt["id"] for attempt in attempts if "grades" in attempt} or set(review) != set(RUBRIC):
            raise ValueError(f"unknown attempt or invalid rubric: {identifier}")
        for item in review.values():
            if type(item.get("score")) is not int or item["score"] not in (0, 1, 2) or not item.get("evidence"):
                raise ValueError(f"score and evidence required: {identifier}")
    summary = {"status": manifest["status"], "model": manifest["model"], "variants": {}, "reviews": reviews}
    summary["suite"] = manifest.get("suite", "legacy")
    summary["execution_seconds"] = {variant: sum(a.get("seconds", 0) for a in attempts if a["variant"] == variant)
                                    for variant in ("without", "with")}
    lines = ["# Оценка скилла muonsoft-validation", "", f"Модель: `{manifest['model']}`. Состояние запуска: **{manifest['status']}**.",
             "", "## Объективные результаты", "",
             "| Вариант | Попыток | Оценено | Прошли обе версии | Остальные статусы |",
             "| --- | ---: | ---: | ---: | --- |"]
    for variant in ("without", "with"):
        group = [a for a in attempts if a["variant"] == variant]
        graded = [a for a in group if "grades" in a]
        passed = sum(all(g["passed"] for g in a["grades"].values()) for a in graded)
        statuses = dict(Counter(a["status"] for a in group if "grades" not in a))
        summary["variants"][variant] = {"attempts": len(group), "graded": len(graded), "passed_both_versions": passed, "ungraded": statuses}
        lines.append(f"| {variant} | {len(group)} | {len(graded)} | {passed} | {statuses} |")
    lines += ["", "Неоценённые попытки показаны отдельно; они не считаются успешными. Баллы ревью не смешиваются с тестами.", "",
              "| Попытка | Статус | v0.19.0 | Текущая | Загрузка скилла |", "| --- | --- | --- | --- | --- |"]
    for attempt in attempts:
        cells = []
        for version in ("v0.19.0", "current"):
            g = attempt.get("grades", {}).get(version)
            cells.append(f"{len(g['requirements_passed'])}/{g['requirements_total']} ({g['status']})" if g else "не оценено")
        identifier = attempt["id"]
        artifact = next((name for name in ("solution.diff", "events.jsonl", "events.jsonl.gz")
                         if (output / "attempts" / identifier / name).is_file()), None)
        label = f"[{identifier}](attempts/{identifier}/{artifact})" if artifact else identifier
        lines.append(f"| {label} | {attempt['status']} | {' | '.join(cells)} | {attempt.get('events', {}).get('skill_loading', 'unknown')} |")
    pairs = {}
    for a in attempts:
        pairs.setdefault((a["case"], a["repeat"]), {})[a["variant"]] = a
    changes = Counter()
    for pair in pairs.values():
        if set(pair) != {"with", "without"} or any("grades" not in a for a in pair.values()):
            changes["incomplete"] += 1
            continue
        success = {k: all(g["passed"] for g in a["grades"].values()) for k, a in pair.items()}
        changes["improved" if success["with"] and not success["without"] else "regressed" if success["without"] and not success["with"] else "unchanged"] += 1
    summary["pairs"] = dict(changes)
    matched = sum(changes[key] for key in ("improved", "regressed", "unchanged"))
    summary["paired_delta_percentage_points"] = round(100 * (changes["improved"] - changes["regressed"]) / matched, 2) if matched else None
    lines += ["", f"Парные результаты: улучшений {changes['improved']}, регрессий {changes['regressed']}, без изменения {changes['unchanged']}, неполных пар {changes['incomplete']}.",
              f"Разница доли успешных сценариев в оценённых парах: {summary['paired_delta_percentage_points']} п.п. (пар: {matched}).",
              "", "## Разброс и невыполненные требования", ""]
    summary["by_case"] = {}
    for case_id in sorted({a["case"] for a in attempts}):
        summary["by_case"][case_id] = {}
        for variant in ("without", "with"):
            group = [a for a in attempts if a["case"] == case_id and a["variant"] == variant]
            rates = []
            for a in group:
                grades = a.get("grades", {})
                if grades:
                    total = sum(g["requirements_total"] for g in grades.values())
                    rates.append(round(100 * sum(len(g["requirements_passed"]) for g in grades.values()) / total, 2))
            summary["by_case"][case_id][variant] = {"requirement_percent_by_repeat": rates, "ungraded": sum("grades" not in a for a in group)}
            lines.append(f"- {case_id}/{variant}: доля пройденных требований по оценённым повторам {rates}; не оценено {sum('grades' not in a for a in group)}.")
    missing = []
    requirements = manifest.get("case_requirements")
    if requirements is None:
        # Never reinterpret an old run with a changed suite's expectations.
        requirements = ({case["id"]: case["requirements"] for case in cases()}
                        if manifest.get("cases_sha256") == digest((HERE / "cases.json").read_bytes()) else {})
    for a in attempts:
        for version, grade_result in a.get("grades", {}).items():
            failed = [name for name in requirements.get(a["case"], []) if name not in grade_result["requirements_passed"]]
            if failed or not grade_result["passed"]:
                missing.append({"attempt": a["id"], "version": version, "requirements": failed, "status": grade_result["status"]})
                log = f"attempts/{a['id']}/checks/{version}/tests.jsonl"
                if not (output / log).is_file() and (output / (log + ".gz")).is_file():
                    log += ".gz"
                lines.append(f"- [{a['id']}/{version}]({log}): {', '.join(failed) or 'сбой процесса тестирования'}.")
    summary["failed_requirements"] = missing
    categories = {}
    for attempt in attempts:
        if "grades" not in attempt:
            continue
        for requirement, category in manifest.get("case_categories", {}).get(attempt["case"], {}).items():
            result = categories.setdefault(category, {}).setdefault(attempt["variant"], {"passed": 0, "total": 0})
            result["total"] += 1
            result["passed"] += int(all(requirement in g["requirements_passed"] for g in attempt["grades"].values()))
    summary["categories"] = categories
    if categories:
        lines += ["", "## Категории требований", "", "Один результат требования учитывает обе версии библиотеки; категории не заменяют успешность целого задания.", ""]
        for category, values in categories.items():
            lines.append(f"- {category}: {values}")
    lines += ["", f"Время выполнения попыток, секунды (включая незавершённые с известной длительностью): {summary['execution_seconds']}."]
    false_positives = [a["id"] for a in attempts if a["variant"] == "with" and not a["should_trigger"] and a.get("events", {}).get("skill_loading") == "observed"]
    summary["potential_discovery_false_positives"] = false_positives
    lines += ["", f"Наблюдаемая загрузка скилла в отрицательном контроле: {false_positives}.",
              "", "## Ревью Codex", "", "Шкала каждого критерия: 0 — не выполнен; 1 — частично; 2 — выполнен. Без ссылки на код или журнал оценка недействительна.", ""]
    for a in attempts:
        review = reviews.get(a["id"])
        lines.append(f"### {a['id']}")
        lines.append("")
        if not review:
            lines.append("Ревью не выполнено.")
        else:
            for key in RUBRIC:
                lines.append(f"- {key}: **{review[key]['score']}/2** — {review[key]['evidence']}")
        lines.append("")
    lines += ["## Воспроизводимость", "", f"Ревизия: `{manifest.get('revision', 'unknown')}`.",
              f"Скилл SHA-256: `{manifest.get('skill_sha256', 'unknown')}`.",
              "Параметры, версии, длительность, токены и стоимость (если переданы OpenCode) сохранены в [run.json](run.json).",
              "", "Результаты относятся к этой модели, конфигурации и выборке. Отсутствие наблюдаемой загрузки скилла не доказывает её наличие или отсутствие.", ""]
    tokens = token_comparison(attempts)
    summary["tokens"] = tokens
    write_json(output / "tokens.json", tokens)
    matched_tokens = tokens["matched"]
    lines += ["## Токены", "", f"Полных оценённых пар с известным usage: {matched_tokens['pairs']}.",
              "Кешированный вход считается отдельно; сумма категорий — объём токенов, не денежная стоимость.",
              "Ошибки тестов включены; сбои выполнения и неизвестный usage исключены попарно.", "",
              "| Категория | Без скилла | Со скиллом | Изменение |", "| --- | ---: | ---: | ---: |"]
    for key in (*TOKEN_KEYS, "uncached_plus_generated", "all_categories"):
        delta = matched_tokens["change_percent"][key]
        change = f"{delta:+.1f}%" if delta is not None else "не определено"
        lines.append(f"| {key} | {matched_tokens['totals']['without'][key]} | {matched_tokens['totals']['with'][key]} | {change} |")
    lines += ["", f"Медианы объёма на попытку: {matched_tokens['median_per_attempt']}.",
              f"Пар с меньшим объёмом со скиллом: {matched_tokens['pairs_with_fewer_tokens']}.",
              f"Исключённые пары: {tokens['excluded_pairs']}.", "",
              "| Сценарий | Пар | Без скилла | Со скиллом | Изменение |", "| --- | ---: | ---: | ---: | ---: |"]
    for case, values in tokens["by_case"].items():
        delta = values["change_percent"]["all_categories"]
        change = f"{delta:+.1f}%" if delta is not None else "не определено"
        lines.append(f"| {case} | {values['pairs']} | {values['totals']['without']['all_categories']} | {values['totals']['with']['all_categories']} | {change} |")
    lines += ["", "Категории, неполный usage и дополнительные срезы: [tokens.json](tokens.json).", ""]
    write_json(output / "results.json", summary)
    (output / "report.md").write_text("\n".join(lines))
    return summary


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    commands = parser.add_subparsers(dest="command", required=True)
    checker = commands.add_parser("check", help="check docs, references and failing starters, without a model")
    checker.add_argument("--output", required=True)
    checker.add_argument("--suite", choices=("main", "legacy", "transfer"), default="main")
    checker.add_argument("--retention", choices=("compact", "debug"), default="compact")
    preflight = commands.add_parser("doctor", help="inspect exact model and isolated skill discovery; no inference")
    preflight.add_argument("--output", required=True)
    preflight.add_argument("--model", required=True)
    preflight.add_argument("--config", help="provider-only JSON; omit on v2 to reuse the local Console connection")
    preflight.add_argument("--opencode", default="opencode")
    runner = commands.add_parser("run", help="manually execute paired OpenCode evals")
    runner.add_argument("--output", required=True)
    runner.add_argument("--model", required=True)
    runner.add_argument("--config", help="provider-only JSON; omit on v2 to reuse the local Console connection")
    runner.add_argument("--opencode", default="opencode")
    runner.add_argument("--profile", choices=("smoke", "full"), default="smoke")
    runner.add_argument("--suite", choices=("main", "legacy", "transfer"), default="main")
    runner.add_argument("--retention", choices=("compact", "debug"), default="compact")
    runner.add_argument("--skip-check", action="store_true", help="skip deterministic material checks after separately verifying this exact revision")
    runner.add_argument("--case", action="append")
    runner.add_argument("--resume", type=Path, help="continue unstarted attempts from an unchanged run into a new output")
    runner.add_argument("--timeout", type=int, default=900)
    reporter = commands.add_parser("report", help="rebuild report without rerunning agents")
    reporter.add_argument("--run", required=True)
    reporter.add_argument("--review", type=Path)
    cleaner = commands.add_parser("cleanup", help="preview disposable run data; --apply compacts it")
    cleaner.add_argument("--run", required=True)
    cleaner.add_argument("--apply", action="store_true")
    args = parser.parse_args()
    try:
        if args.command == "check":
            return check(args)
        if args.command == "doctor":
            return doctor(args)
        if args.command == "run":
            if args.timeout <= 0:
                raise ValueError("timeout must be positive")
            return run(args)
        if args.command == "cleanup":
            return cleanup(args)
        report(Path(args.run).resolve(), args.review)
        return 0
    except (ValueError, OSError, subprocess.CalledProcessError) as exc:
        print(str(exc), file=sys.stderr)
        return 1


if __name__ == "__main__":
    sys.exit(main())
