#!/usr/bin/env bash
# govulncheck v1.8.0 text mode: 0 = clean, 3 = findings, other = failure.
set -euo pipefail
mode=${1:?expected strict or legacy}
scanner=${2:?expected scanner path}
report=${3:?expected report directory}
[[ "$mode" == strict || "$mode" == legacy ]] || exit 2
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
  echo "## Go security: $mode"
  echo
  echo "Result: **$status** (exit $code)."
  echo
  if [[ "$mode" == legacy ]]; then
    echo 'Diagnostic only. Minimum-Go build/tests remain required.'
    echo 'A clean modern-Go scan does not establish safety on this older toolchain.'
  fi
  echo 'Full scanner output and effective toolchain are in the security report artifact.'
} > "$report/summary.md"
if [[ -n "${GITHUB_STEP_SUMMARY:-}" ]]; then
  cat "$report/summary.md" >> "$GITHUB_STEP_SUMMARY"
fi
if [[ "$status" == findings && "$mode" == legacy ]]; then
  echo '::warning::Legacy Go vulnerability findings; see security report.'
  exit 0
fi
if [[ "$status" == error ]]; then
  echo '::error::Security scan failed or could not start; this is not a clean result.'
fi
exit "$code"
