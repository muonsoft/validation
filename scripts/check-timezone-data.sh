#!/usr/bin/env bash
# Compare identifiers only: changes to offsets, rules, version, or archive hash
# do not require updating the validation table. CI reports failures as warnings.
set -euo pipefail

repo_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
cd "$repo_root"

timezone_check_dir=$(mktemp -d "${TMPDIR:-/tmp}/validation-timezone-check.XXXXXX")
trap 'rm -rf -- "$timezone_check_dir"' EXIT

curl --fail --silent --show-error --location \
  --connect-timeout 10 --max-time 30 \
  https://data.iana.org/time-zones/tzdata-latest.tar.gz \
  --output "$timezone_check_dir/tzdata.tar.gz"

go run validate/generate_timezones.go "$timezone_check_dir/tzdata.tar.gz" \
  > "$timezone_check_dir/latest.go"

# Ignore generated metadata; compare only the sorted identifier table.
sed -n '/^var timezoneIdentifiers = /,$p' validate/timezone_identifiers.go \
  > "$timezone_check_dir/bundled.txt"
sed -n '/^var timezoneIdentifiers = /,$p' "$timezone_check_dir/latest.go" \
  > "$timezone_check_dir/latest.txt"
test -s "$timezone_check_dir/bundled.txt"
test -s "$timezone_check_dir/latest.txt"

if cmp -s "$timezone_check_dir/bundled.txt" "$timezone_check_dir/latest.txt"; then
  echo "Bundled timezone identifiers match the latest IANA list."
else
  echo "::warning file=validate/timezone_identifiers.go,title=Timezone identifiers need updating::The bundled identifier list differs from the latest IANA tzdata. Regenerate it with validate/generate_timezones.go, update the documented version in validate/timezone.go, and review the diff."
  diff -u "$timezone_check_dir/bundled.txt" "$timezone_check_dir/latest.txt" || true
fi
