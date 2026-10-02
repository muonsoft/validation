#!/usr/bin/env bash
# Read-only check: HEAD is the validated commit or its single changelog-only child.
set -euo pipefail
candidate=${1:?validated commit required}
candidate=$(git rev-parse --verify "${candidate}^{commit}")
git diff --exit-code --quiet
git diff --cached --exit-code --quiet
head=$(git rev-parse HEAD)
if [[ "$head" == "$candidate" ]]; then
  exit 0
fi
parents=$(git show -s --format=%P HEAD)
changed=$(git diff-tree --no-commit-id --name-only -r HEAD)
if [[ "$parents" != "$candidate" || "$changed" != CHANGELOG.md ]]; then
  echo "release HEAD must be the validated SHA or its single changelog-only child" >&2
  exit 1
fi
