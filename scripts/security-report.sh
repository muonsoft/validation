#!/usr/bin/env bash
# govulncheck v1.8.0 text mode: 0 = clean, 3 = findings, other = failure.
set -euo pipefail
scanner=${1:?expected scanner path}
report=${2:?expected report directory}
mkdir -p "$report"
status=error
code=1
{
  date -u &&
  go version &&
  "$scanner" -version
} > "$report/scan.log" 2>&1 && {
  set +e
  "$scanner" ./... >> "$report/scan.log" 2>&1
  code=$?
  set -e
  case "$code" in
    0) status=clean ;;
    3) status=findings ;;
  esac
}
cat "$report/scan.log"
{
  echo "## Go security"
  echo
  echo "Result: **$status** (exit $code)."
  echo
  echo 'Full scanner output and effective toolchain are in the security report artifact.'
} > "$report/summary.md"
if [[ -n "${GITHUB_STEP_SUMMARY:-}" ]]; then
  cat "$report/summary.md" >> "$GITHUB_STEP_SUMMARY"
fi
if [[ "$status" == error ]]; then
  echo '::error::Security scan failed or could not start; this is not a clean result.'
fi
exit "$code"
