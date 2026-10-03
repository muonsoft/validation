#!/usr/bin/env python3
"""Manual skill evals. Python standard library only; never called by CI."""

import argparse
import base64
from collections import Counter
from contextlib import contextmanager
import difflib
import hashlib
import io
import json
import os
from pathlib import Path
import re
import secrets
import shutil
import signal
import socket
import sqlite3
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


def cases():
    data = json.loads((HERE / "cases.json").read_text())
    if data["schema_version"] != 1:
        raise ValueError("unsupported case schema")
    ids = set()
    for case in data["cases"]:
        identifier = case["id"]
        if not re.fullmatch(r"[0-9]{2}-[a-z]+", identifier) or identifier in ids:
            raise ValueError(f"invalid/duplicate case ID: {identifier}")
        ids.add(identifier)
        path = HERE / "testdata" / identifier
        for name in ("prompt.md", "workspace/contract.go", "workspace/solution.go",
                     "reference/solution.go", "checks/requirements_test.go"):
            if not (path / name).is_file():
                raise ValueError(f"missing {identifier}/{name}")
        tests = re.findall(r"^func (Test\w+)\(", (path / "checks/requirements_test.go").read_text(), re.M)
        if tests != case["requirements"] or not tests:
            raise ValueError(f"test manifest mismatch: {identifier}")
    return data["cases"]


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
        shutil.copyfile(HERE / "testdata" / case["id"] / "reference/solution.go", target / "solution.go")
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
            "requirements_total": len(requirements),
            "passed": process["status"] == "completed" and len(passed) == len(requirements)}


def grade(case, source, output, snapshots, env):
    results = {}
    for version, library in snapshots.items():
        target = output / version
        target.mkdir(parents=True)
        # Only immutable contract and allowed submission are compiled, not worker tests.
        shutil.copyfile(HERE / "testdata" / case["id"] / "workspace/contract.go", target / "contract.go")
        shutil.copyfile(source / "solution.go", target / "solution.go")
        prepare_module(target, Path(library["path"]), negative=not case["should_trigger"])
        shutil.copyfile(HERE / "testdata" / case["id"] / "checks/requirements_test.go", target / "requirements_test.go")
        if case["should_trigger"]:
            shutil.copyfile(HERE / "testdata/helpers_test.go", target / "helpers_test.go")
        log = target / "tests.jsonl"
        process = command(["go", "test", "-mod=mod", "-race", "-count=1", "-json", "./..."], target, log, env)
        results[version] = test_result(log, case["requirements"], process)
        if process["status"] == "cancelled":
            break
    return results


def events_result(path, completion=None):
    parsed, malformed, stopped, errors = 0, 0, False, []
    evidence, usage = [], []
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
            if "tokens" in part or "cost" in part:
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
        if stopped and not completion.get("finish_logged") and any(key in completion for key in ("tokens", "cost")):
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
    for case in cases():
        directory = output / case["id"]
        for mode in ("reference", "starter"):
            src = directory / mode / "workspace"
            workspace(case, src, Path(snapshots["v0.19.0"]["path"]), reference=mode == "reference")
            grades = grade(case, src, directory / mode / "checks", snapshots, go_env(output))
            report["cases"].setdefault(case["id"], {})[mode] = grades
        print(f"checked {case['id']}", flush=True)
    write_json(output / "check.json", report)
    failures = list(report["materials"]["errors"])
    for identifier, modes in report["cases"].items():
        for version in snapshots:
            if not modes["reference"][version]["passed"]:
                failures.append(f"reference failed: {identifier}/{version}")
            if modes["starter"][version]["passed"]:
                failures.append(f"checks did not reject starter: {identifier}/{version}")
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


def run(args):
    selected = [case for case in cases() if (not args.case or case["id"] in args.case)
                and (args.profile == "full" or case["smoke"] or args.case)]
    if args.case and set(args.case) - {case["id"] for case in selected}:
        raise ValueError("unknown case ID")
    output = new_output(args.output)
    env, major, version = prepare_runtime(args, output)
    model_preflight(args, output, env, major)
    snapshots = library_snapshots(output)
    manifest = provenance(snapshots)
    manifest.update(model=args.model, profile=args.profile, timeout=args.timeout,
                    repeats=1 if args.profile == "smoke" else 3, attempts=[], status="running",
                    opencode_version=version)
    write_json(output / "run.json", manifest)
    try:
        for case in selected:
            for repeat in range(manifest["repeats"]):
                # Alternate order to reduce systematic ordering effects.
                variants = ("without", "with") if repeat % 2 == 0 else ("with", "without")
                for variant in variants:
                    identifier = f"{case['id']}-{variant}-{repeat + 1}"
                    directory = output / "attempts" / identifier
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
                        unexpected = [name for name in changed if name not in case["editable"]]
                        attempt["unexpected_changes"] = unexpected
                        original = (HERE / "testdata" / case["id"] / "workspace/solution.go").read_text()
                        solution = (src / "solution.go").read_text()
                        (directory / "solution.diff").write_text("".join(difflib.unified_diff(original.splitlines(True), solution.splitlines(True), fromfile="before/solution.go", tofile="after/solution.go")))
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
        report(output, None)
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
        lines.append(f"| [{attempt['id']}](attempts/{attempt['id']}/solution.diff) | {attempt['status']} | {' | '.join(cells)} | {attempt.get('events', {}).get('skill_loading', 'unknown')} |")
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
    requirements = {case["id"]: case["requirements"] for case in cases()}
    for a in attempts:
        for version, grade_result in a.get("grades", {}).items():
            failed = [name for name in requirements[a["case"]] if name not in grade_result["requirements_passed"]]
            if failed or not grade_result["passed"]:
                missing.append({"attempt": a["id"], "version": version, "requirements": failed, "status": grade_result["status"]})
                lines.append(f"- [{a['id']}/{version}](attempts/{a['id']}/checks/{version}/tests.jsonl): {', '.join(failed) or 'сбой процесса тестирования'}.")
    summary["failed_requirements"] = missing
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
    write_json(output / "results.json", summary)
    (output / "report.md").write_text("\n".join(lines))
    return summary


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    commands = parser.add_subparsers(dest="command", required=True)
    checker = commands.add_parser("check", help="check docs, references and failing starters, without a model")
    checker.add_argument("--output", required=True)
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
    runner.add_argument("--case", action="append")
    runner.add_argument("--timeout", type=int, default=900)
    reporter = commands.add_parser("report", help="rebuild report without rerunning agents")
    reporter.add_argument("--run", required=True)
    reporter.add_argument("--review", type=Path)
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
        report(Path(args.run).resolve(), args.review)
        return 0
    except (ValueError, OSError, subprocess.CalledProcessError) as exc:
        print(str(exc), file=sys.stderr)
        return 1


if __name__ == "__main__":
    sys.exit(main())
