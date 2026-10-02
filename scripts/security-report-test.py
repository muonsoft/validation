"""Exercise scanner status handling without network access."""
import os
from pathlib import Path
import subprocess
import tempfile
import unittest

SCRIPT = Path(__file__).with_name('security-report.sh')

class SecurityReportTest(unittest.TestCase):
    def run_case(self, mode, scan_code, expected, status, preflight=False, missing=False, go_failure=False):
        with tempfile.TemporaryDirectory() as tmp:
            root = Path(tmp)
            scanner = root / 'scanner'
            if not missing:
                scanner.write_text(f'#!/bin/sh\nif [ "$1" = -version ]; then exit {1 if preflight else 0}; fi\necho scan-output\nexit {scan_code}\n')
                scanner.chmod(0o755)
            go = root / 'go'
            go.write_text(f'#!/bin/sh\necho go-version-fixture\nexit {1 if go_failure else 0}\n')
            go.chmod(0o755)
            env = dict(os.environ, PATH=f'{tmp}:'+os.environ['PATH'], GITHUB_STEP_SUMMARY=str(root/'summary'))
            result = subprocess.run(['bash', str(SCRIPT), mode, str(scanner), str(root/'report')], env=env, capture_output=True, text=True)
            self.assertEqual(result.returncode, expected, result.stdout+result.stderr)
            self.assertIn(f'**{status}**', (root/'summary').read_text())
            self.assertIn('go-version-fixture', (root/'report/scan.log').read_text())
            if status == 'error':
                self.assertIn('::error::', result.stdout)
            if status == 'findings' and mode == 'legacy':
                self.assertIn('::warning::', result.stdout)

    def test_result_matrix(self):
        for mode in ['strict', 'legacy']:
            for code, status in [(0, 'clean'), (3, 'findings'), (1, 'error'), (2, 'error'), (137, 'error')]:
                with self.subTest(mode=mode, code=code):
                    self.run_case(mode, code, 0 if mode == 'legacy' and code == 3 else code, status)

    def test_failed_preflight(self):
        for mode in ['strict', 'legacy']:
            self.run_case(mode, 0, 1, 'error', preflight=True)

    def test_failed_go(self):
        self.run_case('legacy', 0, 1, 'error', go_failure=True)

    def test_missing_scanner(self):
        self.run_case('legacy', 0, 1, 'error', missing=True)

if __name__ == '__main__':
    unittest.main()
