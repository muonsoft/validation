#!/usr/bin/env python3
"""Check repository Markdown links and run complete Go examples (no network crawl).

Supports inline links, reference definitions, and GitHub-style heading anchors used
by this repository. Fenced code and inline code are excluded from link checks.
Go fences containing package main are run from the module root; an optional final
// Output: comment is compared with stdout. Partial snippets remain the author's
responsibility and should link to executable Example functions.
"""

import os
from pathlib import Path
import re
import subprocess
import sys
import tempfile
from urllib.parse import unquote, urlsplit

ROOT = Path(__file__).resolve().parent.parent
FENCE = re.compile(r"^```([^\n]*)\n(.*?)^```[ \t]*$", re.MULTILINE | re.DOTALL)
LINK = re.compile(r"\[[^\]\n]*\]\(([^\s)]+)(?:\s+\"[^\"]*\")?\)")
REFERENCE = re.compile(r"^\[[^\]]+\]:\s*(\S+)", re.MULTILINE)


def prose(text):
    text = FENCE.sub("", text)
    return re.sub(r"`[^`\n]+`", "", text)


def anchors(text):
    result = set()
    counts = {}
    for heading in re.findall(r"^#{1,6}\s+(.+?)\s*#*\s*$", FENCE.sub("", text), re.MULTILINE):
        heading = re.sub(r"\[([^\]]+)\]\([^)]+\)", r"\1", heading)
        slug = re.sub(r"[^\w\- ]", "", heading.lower()).replace(" ", "-")
        count = counts.get(slug, 0)
        counts[slug] = count + 1
        result.add(f"{slug}-{count}" if count else slug)
    return result


def check_links(path, text):
    errors = []
    for match in list(LINK.finditer(prose(text))) + list(REFERENCE.finditer(prose(text))):
        url = urlsplit(match.group(1).strip("<>"))
        if url.scheme or url.netloc:
            continue
        target = (path.parent / unquote(url.path)).resolve() if url.path else path
        if not target.exists():
            errors.append(f"{path.relative_to(ROOT)}: missing target {match.group(1)}")
        elif url.fragment and target.suffix == ".md":
            if unquote(url.fragment) not in anchors(target.read_text()):
                errors.append(f"{path.relative_to(ROOT)}: missing heading {match.group(1)}")
    return errors


def check_programs(path, text):
    errors = []
    count = 0
    for block in FENCE.finditer(text):
        source = block.group(2)
        if block.group(1).strip() != "go" or not re.search(r"^package main\s*$", source, re.MULTILINE):
            continue
        count += 1
        line = text.count("\n", 0, block.start()) + 1
        label = f"{path.relative_to(ROOT)}:{line}"
        with tempfile.TemporaryDirectory(prefix="validation-docs-") as temporary:
            program = Path(temporary) / "main.go"
            program.write_text(source)
            try:
                run = subprocess.run(
                    ["go", "run", str(program)], cwd=ROOT,
                    env={**os.environ, "GOWORK": "off"},
                    capture_output=True, text=True, timeout=120, check=False,
                )
            except subprocess.TimeoutExpired:
                errors.append(f"{label}: Go example timed out")
                continue
        if run.returncode:
            errors.append(f"{label}: Go example failed\n{run.stderr}")
            continue
        output = re.search(r"^// Output:\n((?://[^\n]*\n?)*)\s*$", source, re.MULTILINE)
        if output:
            expected = "\n".join(re.sub(r"^// ?", "", line) for line in output.group(1).splitlines()) + "\n"
            if run.stdout != expected:
                errors.append(f"{label}: output mismatch\nexpected: {expected!r}\nactual: {run.stdout!r}")
    return errors, count


def main():
    paths = sorted(ROOT.glob("*.md")) + sorted((ROOT / "docs").rglob("*.md"))
    errors = []
    count = 0
    for path in paths:
        text = path.read_text()
        errors.extend(check_links(path, text))
        program_errors, programs = check_programs(path, text)
        errors.extend(program_errors)
        count += programs
    if errors:
        print("\n".join(errors), file=sys.stderr)
        return 1
    print(f"Documentation checks passed: {len(paths)} Markdown files, {count} runnable Go examples")
    return 0


if __name__ == "__main__":
    sys.exit(main())
