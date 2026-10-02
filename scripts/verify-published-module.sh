#!/usr/bin/env bash
# Verify remote publication without workspace or local replacement leakage.
set -euo pipefail
version=${1:?release version required}
repo_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
[[ "$version" == v* ]] || { echo 'version must have v prefix' >&2; exit 1; }
bash "$repo_root/scripts/prepare-release.sh" "${version#v}" --semver-only
module=$(awk '/^module / { print $2; exit }' "$repo_root/go.mod")
[[ -n "$module" ]] || exit 1
scratch=$(mktemp -d)
trap 'rm -rf -- "$scratch"' EXIT
export GOWORK=off GOFLAGS=-modcacherw GOTOOLCHAIN=local
export GOMODCACHE="$scratch/modcache"
mkdir "$scratch/consumer"
cd "$scratch/consumer"
go mod init example.com/release-smoke
# Retry only fetching: public proxies may not observe the new tag immediately.
for attempt in 1 2 3 4 5; do
  if go mod download "$module@$version"; then break; fi
  if [[ "$attempt" == 5 ]]; then
    echo 'published version did not resolve; inspect publication, never move the tag' >&2
    exit 1
  fi
  sleep "$((attempt * 5))"
done
go mod edit "-require=$module@$version"
printf 'package main\n' > main.go
while IFS= read -r package; do
  [[ "$package" == "$module" || "$package" == "$module/"* ]] || {
    echo "consumer imports must belong to $module" >&2; exit 1;
  }
  printf 'import _ "%s"\n' "$package" >> main.go
done < "$repo_root/scripts/release-consumer-imports.txt"
printf 'func main() {}\n' >> main.go
go mod tidy
go list -m -json "$module" > resolved.json
python3 - "$module" "$version" <<'PY'
import json, sys
with open('resolved.json') as f:
    resolved = json.load(f)
if resolved.get('Path') != sys.argv[1] or resolved.get('Version') != sys.argv[2] or 'Replace' in resolved:
    raise SystemExit('consumer resolved an unexpected module/version or replacement')
PY
go build -buildvcs=false ./...
go version
printf 'Published consumer verified: %s@%s\n' "$module" "$version"
