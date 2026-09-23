#!/usr/bin/env bash
set -euo pipefail

repo_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
cd "$repo_root"

echo "==> format check"
unformatted=$(gofmt -l .)
if [[ -n "$unformatted" ]]; then
  echo "gofmt needed for:" >&2
  echo "$unformatted" >&2
  exit 1
fi

echo "==> vet"
go vet ./...

if command -v golangci-lint >/dev/null 2>&1; then
  echo "==> lint"
  golangci-lint run
fi

echo "==> unit tests"
go test ./...

echo "==> race tests"
go test -race ./...

echo "==> module tidy"
go mod tidy -diff

echo "==> module verify"
go mod verify

echo "==> all checks passed"
