#!/usr/bin/env python3
"""Manual skill evals. Python standard library only; never called by CI."""

import argparse
from collections import Counter
import difflib
import hashlib
import io
import json
import os
from pathlib import Path
import re
import shutil
import signal
import subprocess
import sys
import tarfile
import time

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


def events_result(path):
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
        if part.get("tool") == "skill" and args.get("name") == "muonsoft-validation":
            evidence.append(number)
        if part.get("tool") == "read" and str(args.get("filePath", "")).endswith("/muonsoft-validation/SKILL.md"):
            evidence.append(number)
    return {"complete": bool(parsed and stopped and not malformed and not errors),
            "malformed_lines": malformed, "error_lines": errors,
            "skill_loading": "observed" if evidence else "not_observed",
            "skill_evidence_lines": evidence, "reported_usage": usage}


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
    output.mkdir(parents=True, exist_ok=False)
    return output


def isolated_env(output, config):
    env = {key: value for key, value in os.environ.items() if not key.startswith("OPENCODE_")}
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


def opencode_config(path):
    config = json.loads(path.read_text())
    # Only provider setup is inherited, never prompts, plugins, MCP, skills, or agents.
    if set(config) - {"$schema", "provider"}:
        raise ValueError("eval provider config may contain only $schema and provider")
    if not config.get("provider"):
        raise ValueError("provider config is empty; use explicit provider setup with env credentials")
    # Authentication must remain in environment, not generated artifacts.
    for provider in config["provider"].values():
        options = provider.get("options", {})
        key = options.get("apiKey")
        if key and not re.fullmatch(r"\{env:[A-Z_][A-Z_0-9]*\}", key):
            raise ValueError("use an {env:VARIABLE} apiKey reference, not a literal credential")
        for key, value in options.get("headers", {}).items():
            if not re.fullmatch(r"\{env:[A-Z_][A-Z_0-9]*\}", value):
                raise ValueError(f"custom header {key} must use an environment reference")
    config.update(autoupdate=False, share="disabled", instructions=[], plugin=[], mcp={},
                  permission={"*": "deny", "read": "allow", "glob": "allow", "grep": "allow",
                              "edit": "allow", "bash": "allow", "skill": {"*": "deny", "muonsoft-validation": "allow"},
                              "external_directory": "deny"})
    return config


def discover(binary, cwd, output, env, expected):
    log = output / "discovery.json"
    result = command([binary, "--pure", "debug", "skill"], cwd, log, env)
    if result["status"] != "completed":
        raise ValueError(f"cannot inspect available skills; see {log}")
    try:
        entries = json.loads(log.read_text())
        names = sorted(entry["name"] for entry in entries)
    except (ValueError, KeyError, TypeError) as exc:
        raise ValueError(f"unsupported OpenCode skill discovery output: {log}") from exc
    if names != expected:
        raise ValueError(f"unexpected skills {names}; expected {expected}. Use a clean environment.")


def model_preflight(args, output, env):
    catalog = output / "models.txt"
    status = command([args.opencode, "--pure", "models", args.model.split("/", 1)[0]], output, catalog, env)
    if status["status"] != "completed" or args.model not in catalog.read_text().splitlines():
        raise ValueError(f"requested model unavailable: {args.model}; see {catalog}; no fallback")


def doctor(args):
    if "/" not in args.model:
        raise ValueError("--model must be the exact provider/model ID")
    config = opencode_config(Path(args.config))
    output = new_output(args.output)
    config_path = output / "provider.json"
    write_json(config_path, config)
    env = isolated_env(output, config_path)
    model_preflight(args, output, env)
    for variant in ("without", "with"):
        directory = output / variant
        directory.mkdir()
        subprocess.run(["git", "init", "-q", str(directory)], check=True)
        if variant == "with":
            shutil.copytree(SKILL, directory / ".opencode/skills/muonsoft-validation")
        discover(args.opencode, directory, output / (variant + "-logs"), env,
                 ["muonsoft-validation"] if variant == "with" else [])
    write_json(output / "doctor.json", {"model": args.model, "discovery": "passed", "inference_requests": 0})
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
    if not args.model or "/" not in args.model:
        raise ValueError("--model must be the exact provider/model ID for the requested model")
    config = opencode_config(Path(args.config))
    selected = [case for case in cases() if (not args.case or case["id"] in args.case)
                and (args.profile == "full" or case["smoke"] or args.case)]
    if args.case and set(args.case) - {case["id"] for case in selected}:
        raise ValueError("unknown case ID")
    output = new_output(args.output)
    config_path = output / "provider.json"
    write_json(config_path, config)
    env = isolated_env(output, config_path)
    model_preflight(args, output, env)
    snapshots = library_snapshots(output)
    manifest = provenance(snapshots)
    manifest.update(model=args.model, profile=args.profile, timeout=args.timeout,
                    repeats=1 if args.profile == "smoke" else 3, attempts=[], status="running",
                    opencode_version=subprocess.run([args.opencode, "--version"], capture_output=True, text=True, check=True).stdout.strip())
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
                    discover(args.opencode, src, directory, env, ["muonsoft-validation"] if variant == "with" else [])
                    before = files_digest(src)
                    prompt = (HERE / "testdata" / case["id"] / "prompt.md").read_text()
                    (directory / "prompt.md").write_text(prompt)
                    attempt["status"] = "running"
                    write_json(output / "run.json", manifest)
                    result = command([args.opencode, "--pure", "run", "--agent", "build", "--model", args.model,
                                      "--format", "json", prompt], src, directory / "events.jsonl", env, args.timeout)
                    events = events_result(directory / "events.jsonl")
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
    preflight.add_argument("--config", required=True)
    preflight.add_argument("--opencode", default="opencode")
    runner = commands.add_parser("run", help="manually execute paired OpenCode evals")
    runner.add_argument("--output", required=True)
    runner.add_argument("--model", required=True)
    runner.add_argument("--config", required=True, help="provider-only JSON config using environment credentials")
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
