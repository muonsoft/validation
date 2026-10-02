#!/usr/bin/env python3
"""Regression checks for Markdown links and executable documentation examples."""

import importlib.util
from pathlib import Path
import subprocess
import tempfile
import unittest
from unittest.mock import patch

spec = importlib.util.spec_from_file_location("check_docs", Path(__file__).with_name("check-docs.py"))
docs = importlib.util.module_from_spec(spec)
spec.loader.exec_module(docs)


class DocumentationChecks(unittest.TestCase):
    def setUp(self):
        self.temporary = tempfile.TemporaryDirectory()
        self.addCleanup(self.temporary.cleanup)
        self.root = Path(self.temporary.name)
        self.path = self.root / "README.md"
        self.patch = patch.object(docs, "ROOT", self.root)
        self.patch.start()
        self.addCleanup(self.patch.stop)

    def test_links_and_heading_anchors(self):
        (self.root / "guide.md").write_text("# Guide\n## `Some` API\n## `Some` API\n")
        self.assertEqual(docs.check_links(self.path, "[guide](guide.md#some-api-1)"), [])
        self.assertEqual(len(docs.check_links(self.path, "[guide](guide.md#absent)")), 1)
        self.assertEqual(len(docs.check_links(self.path, "[guide](missing.md)")), 1)
        self.assertEqual(len(docs.check_links(self.path, "[guide]: missing.md")), 1)

    def test_code_and_external_links_are_not_local_targets(self):
        text = "`[skip](missing)`\n```go\n[skip](missing)\n```\n[web](https://example.com)"
        self.assertEqual(docs.check_links(self.path, text), [])

    def test_partial_snippets_are_not_executed(self):
        with patch.object(docs.subprocess, "run") as run:
            self.assertEqual(docs.check_programs(self.path, '```go\nfmt.Println("x")\n```'), ([], 0))
            run.assert_not_called()

    def test_program_output_is_checked(self):
        text = '```go\npackage main\nfunc main() {}\n// Output:\n// hello\n```'
        with patch.object(docs.subprocess, "run", return_value=subprocess.CompletedProcess([], 0, "hello\n", "")):
            self.assertEqual(docs.check_programs(self.path, text), ([], 1))
        with patch.object(docs.subprocess, "run", return_value=subprocess.CompletedProcess([], 0, "wrong\n", "")):
            errors, count = docs.check_programs(self.path, text)
            self.assertEqual(count, 1)
            self.assertIn("output mismatch", errors[0])

    def test_compile_errors_and_timeouts_fail(self):
        text = '```go\npackage main\nfunc main() {}\n```'
        with patch.object(docs.subprocess, "run", return_value=subprocess.CompletedProcess([], 1, "", "compile error")):
            errors, _ = docs.check_programs(self.path, text)
            self.assertIn("compile error", errors[0])
        with patch.object(docs.subprocess, "run", side_effect=subprocess.TimeoutExpired("go", 120)):
            errors, _ = docs.check_programs(self.path, text)
            self.assertIn("timed out", errors[0])


if __name__ == "__main__":
    unittest.main()
