#!/usr/bin/env python3
"""Exercise release ancestry checks in disposable Git repositories."""
import os
from pathlib import Path
import subprocess
import tempfile
import unittest

SCRIPT = Path(__file__).resolve().with_name('verify-release-candidate.sh')

class CandidateTest(unittest.TestCase):
    def setUp(self):
        self.temp = tempfile.TemporaryDirectory()
        self.addCleanup(self.temp.cleanup)
        self.root = Path(self.temp.name)
        self.git('init', '-q')
        self.git('config', 'user.email', 'test@example.invalid')
        self.git('config', 'user.name', 'Release test')
        self.commit('CHANGELOG.md', 'base')
        self.base = self.git('rev-parse', 'HEAD')

    def git(self, *args):
        return subprocess.check_output(['git', *args], cwd=self.root, text=True,
                                       stderr=subprocess.STDOUT).strip()

    def commit(self, path, text):
        (self.root / path).write_text(text)
        self.git('add', path)
        self.git('commit', '-qm', text)

    def check(self, expected):
        result = subprocess.run(['bash', str(SCRIPT), self.base], cwd=self.root,
                                capture_output=True, text=True)
        self.assertEqual(result.returncode == 0, expected, result.stderr)

    def test_unchanged_candidate(self):
        self.check(True)

    def test_one_changelog_child(self):
        self.commit('CHANGELOG.md', 'release')
        self.check(True)

    def test_code_child_rejected(self):
        self.commit('source.go', 'code')
        self.check(False)

    def test_two_children_rejected(self):
        self.commit('CHANGELOG.md', 'release')
        self.commit('CHANGELOG.md', 'another')
        self.check(False)

    def test_dirty_tracked_file_rejected(self):
        (self.root / 'CHANGELOG.md').write_text('uncommitted')
        self.check(False)

    def test_merge_rejected(self):
        self.git('checkout', '-qb', 'side')
        self.commit('side.txt', 'side')
        self.git('checkout', '--detach', self.base)
        self.commit('CHANGELOG.md', 'release')
        self.git('merge', '--no-ff', '-m', 'merge', 'side')
        self.check(False)

    def test_annotated_candidate_resolves(self):
        self.git('tag', '-a', 'candidate', '-m', 'candidate')
        self.base = 'candidate'
        self.check(True)


class PublicationTest(CandidateTest):
    def setUp(self):
        super().setUp()
        workflow = SCRIPT.parent.parent / '.github/workflows/release.yml'
        text = workflow.read_text()
        block = text.split('      - name: Push verified changelog commit\n', 1)[1]
        block = block.split('      - name:', 1)[0].split('        run: |\n', 1)[1]
        self.command = '\n'.join(line[10:] for line in block.splitlines() if line.strip())
        self.branch = 'master' if 'origin/master' in self.command else 'main'
        self.remote = self.root / 'remote.git'
        self.git('init', '--bare', '-q', str(self.remote))
        self.git('remote', 'add', 'origin', str(self.remote))
        self.git('push', '-q', 'origin', 'HEAD:refs/heads/' + self.branch)
        (self.root / 'scripts').mkdir()
        (self.root / 'scripts/verify-release-candidate.sh').write_text(SCRIPT.read_text())

    def publish(self):
        env = os.environ.copy()
        env['DISPATCH_SHA'] = self.base
        return subprocess.run(['bash', '-euo', 'pipefail', '-c', self.command],
                              cwd=self.root, env=env, capture_output=True, text=True)

    def remote_head(self):
        return self.git('--git-dir=' + str(self.remote), 'rev-parse', self.branch)

    def test_valid_child_pushes(self):
        self.commit('CHANGELOG.md', 'release')
        result = self.publish()
        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertEqual(self.remote_head(), self.git('rev-parse', 'HEAD'))

    def test_stale_branch_does_not_push(self):
        self.commit('other.txt', 'remote work')
        remote_head = self.git('rev-parse', 'HEAD')
        self.git('push', '-q', 'origin', 'HEAD:refs/heads/' + self.branch)
        self.git('checkout', '--detach', self.base)
        self.commit('CHANGELOG.md', 'release')
        self.assertNotEqual(self.publish().returncode, 0)
        self.assertEqual(self.remote_head(), remote_head)

    def test_code_change_does_not_push(self):
        self.commit('other.txt', 'unexpected code')
        self.assertNotEqual(self.publish().returncode, 0)
        self.assertEqual(self.remote_head(), self.base)

if __name__ == '__main__':
    unittest.main()
