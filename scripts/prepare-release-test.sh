#!/usr/bin/env bash
# Isolated tests for scripts/prepare-release.sh. CHANGELOG.md is never modified.

set -euo pipefail

repo_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
prepare="$repo_root/scripts/prepare-release.sh"
test_root=$(mktemp -d "${TMPDIR:-/tmp}/validation-prepare-release-test.XXXXXX")
passed=0
failed=0

cleanup() {
  case "$test_root" in
    "${TMPDIR:-/tmp}"/validation-prepare-release-test.*) rm -rf -- "$test_root" ;;
    *) echo "refusing to remove unexpected test directory: $test_root" >&2 ;;
  esac
}
trap cleanup EXIT INT HUP TERM

pass() {
  echo "PASS: $1"
  passed=$((passed + 1))
}

fail() {
  echo "FAIL: $1" >&2
  failed=$((failed + 1))
}

run_prepare() {
  local changelog_file=$1
  local release_version=$2
  shift 2
  CHANGELOG_FILE="$changelog_file" "$prepare" "$release_version" "$@"
}

fixture() {
  local name=$1
  echo "$test_root/${name}.md"
}

test_invalid_versions() {
  local file
  file=$(fixture invalid)
  printf '## [Unreleased]\n' >"$file"
  if run_prepare "$file" 01.2.3 --semver-only >/dev/null 2>&1 ||
    run_prepare "$file" 1.2.3-01 --semver-only >/dev/null 2>&1 ||
    run_prepare "$file" 1.2.3+build --semver-only >/dev/null 2>&1; then
    fail "invalid SemVer rejected"
  else
    pass "invalid SemVer rejected"
  fi
}

test_planned_check_and_finalize() {
  local file
  file=$(fixture planned)
  printf '%s\n' \
    '## [Unreleased]' \
    '' \
    '## [0.1.0] — planned' \
    '' \
    '### Added' \
    '' \
    '- first release' >"$file"

  if run_prepare "$file" 0.1.0 --check-only >/dev/null &&
    run_prepare "$file" 0.1.0 2026-09-04 >/dev/null &&
    grep -qFx '## [0.1.0] - 2026-09-04' "$file"; then
    pass "planned section finalized"
  else
    fail "planned section finalized"
  fi
}

test_finalize_is_idempotent() {
  local file
  file=$(fixture idempotent)
  printf '%s\n' \
    '## [Unreleased]' \
    '' \
    '## [0.2.0] - 2026-09-04' \
    '' \
    '### Changed' \
    '' \
    '- existing' >"$file"

  if run_prepare "$file" 0.2.0 2026-09-05 >/dev/null &&
    grep -qFx '## [0.2.0] - 2026-09-04' "$file"; then
    pass "finalization is idempotent"
  else
    fail "finalization is idempotent"
  fi
}

test_unreleased_promotion() {
  local file
  file=$(fixture unreleased)
  printf '%s\n' \
    '## [Unreleased]' \
    '' \
    '### Fixed' \
    '' \
    '- future fix' \
    '' \
    '## [0.1.0] - 2026-09-01' >"$file"

  if run_prepare "$file" 0.1.1 2026-09-04 >/dev/null &&
    grep -qFx '## [0.1.1] - 2026-09-04' "$file" &&
    grep -qFx -- '- future fix' "$file"; then
    pass "unreleased section promoted"
  else
    fail "unreleased section promoted"
  fi
}

test_require_final() {
  local file
  file=$(fixture require-final)
  printf '%s\n' \
    '## [Unreleased]' \
    '' \
    '## [0.3.0] — planned' \
    '' \
    '### Added' \
    '' \
    '- pending' >"$file"

  if run_prepare "$file" 0.3.0 --require-final >/dev/null 2>&1; then
    fail "require-final rejects planned section"
  else
    pass "require-final rejects planned section"
  fi
}

test_require_final_accepts_finalized() {
  local file
  file=$(fixture require-finalized)
  printf '%s\n' \
    '## [Unreleased]' \
    '' \
    '## [0.3.1] - 2026-09-04' \
    '' \
    '### Fixed' \
    '' \
    '- ready' >"$file"

  if run_prepare "$file" 0.3.1 --require-final >/dev/null; then
    pass "require-final accepts finalized section"
  else
    fail "require-final accepts finalized section"
  fi
}

test_empty_sections_fail() {
  local file
  file=$(fixture empty)
  printf '%s\n' '## [Unreleased]' '' '## [0.4.0] — planned' >"$file"
  if run_prepare "$file" 0.4.0 --check-only >/dev/null 2>&1; then
    fail "empty sections rejected"
  else
    pass "empty sections rejected"
  fi
}

test_invalid_date_fails() {
  local file
  file=$(fixture date)
  printf '%s\n' '## [Unreleased]' '' '### Added' '' '- item' >"$file"
  if run_prepare "$file" 0.5.0 04-09-2026 >/dev/null 2>&1; then
    fail "invalid date rejected"
  else
    pass "invalid date rejected"
  fi
}

test_invalid_versions
test_planned_check_and_finalize
test_finalize_is_idempotent
test_unreleased_promotion
test_require_final
test_require_final_accepts_finalized
test_empty_sections_fail
test_invalid_date_fails

echo "prepare-release-test: ${passed} passed, ${failed} failed"
if [[ "$failed" -ne 0 ]]; then
  exit 1
fi
