#!/usr/bin/env bash
# Prepare CHANGELOG.md for a SemVer release (Keep a Changelog).
#
# Supports an exact planned section (`## [X.Y.Z] — planned`), promotion of a
# non-empty [Unreleased] section, validation-only modes, and idempotent reruns.
#
# Usage:
#   scripts/prepare-release.sh 0.1.0 [YYYY-MM-DD]
#   scripts/prepare-release.sh 0.1.0 --check-only
#   scripts/prepare-release.sh 0.1.0 --require-final
#   scripts/prepare-release.sh 0.1.0 --semver-only

set -euo pipefail

version=${1:?version required (for example 0.1.0)}
changelog=${CHANGELOG_FILE:-CHANGELOG.md}
release_date=
check_only=false
require_final=false
semver_only=false

shift
while [[ $# -gt 0 ]]; do
  case "$1" in
    --check-only)
      check_only=true
      ;;
    --require-final)
      require_final=true
      ;;
    --semver-only)
      semver_only=true
      ;;
    --*)
      echo "unknown option: $1" >&2
      exit 1
      ;;
    *)
      if [[ -n "$release_date" ]]; then
        echo "unexpected extra argument: $1" >&2
        exit 1
      fi
      release_date=$1
      ;;
  esac
  shift
done

validate_numeric_identifier() {
  local identifier=$1
  local label=$2
  if [[ ! "$identifier" =~ ^(0|[1-9][0-9]*)$ ]]; then
    echo "invalid ${label} identifier: ${identifier}" >&2
    return 1
  fi
}

validate_semver() {
  local candidate=$1
  local core prerelease major minor patch identifier

  if [[ "$candidate" == *+* ]]; then
    echo "build metadata is not supported in release versions: ${candidate}" >&2
    return 1
  fi

  core=${candidate%%-*}
  prerelease=
  if [[ "$candidate" == *-* ]]; then
    prerelease=${candidate#*-}
  fi

  IFS=. read -r major minor patch <<<"$core"
  if [[ -z "${major:-}" || -z "${minor:-}" || -z "${patch:-}" || "$core" == *.*.*.* ]]; then
    echo "invalid SemVer core: ${candidate}" >&2
    return 1
  fi
  validate_numeric_identifier "$major" major || return 1
  validate_numeric_identifier "$minor" minor || return 1
  validate_numeric_identifier "$patch" patch || return 1

  if [[ "$candidate" == *-* ]]; then
    if [[ -z "$prerelease" || "$prerelease" == .* || "$prerelease" == *. || "$prerelease" == *..* ]]; then
      echo "invalid prerelease identifiers: ${candidate}" >&2
      return 1
    fi
    IFS=. read -ra prerelease_parts <<<"$prerelease"
    for identifier in "${prerelease_parts[@]}"; do
      if [[ "$identifier" =~ ^[0-9]+$ ]]; then
        validate_numeric_identifier "$identifier" prerelease || return 1
      elif [[ ! "$identifier" =~ ^[0-9A-Za-z-]+$ ]]; then
        echo "invalid prerelease identifier: ${identifier}" >&2
        return 1
      fi
    done
  fi
}

validate_semver "$version"
if [[ "$semver_only" == true ]]; then
  exit 0
fi

if [[ ! -f "$changelog" ]]; then
  echo "changelog file not found: $changelog" >&2
  exit 1
fi

if [[ -z "$release_date" ]]; then
  release_date=$(date -u +%Y-%m-%d)
fi
if [[ ! "$release_date" =~ ^[0-9]{4}-[0-9]{2}-[0-9]{2}$ ]]; then
  echo "invalid release date: $release_date (expected YYYY-MM-DD)" >&2
  exit 1
fi

final_prefix="## [${version}] - "
planned_header="## [${version}] — planned"

is_finalized() {
  awk -v prefix="$final_prefix" '
    index($0, prefix) == 1 {
      value = substr($0, length(prefix) + 1)
      if (value ~ /^[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]$/) {
        found = 1
        exit
      }
    }
    END { exit(found ? 0 : 1) }
  ' "$changelog"
}

is_planned() {
  grep -qFx "$planned_header" "$changelog"
}

section_has_content() {
  local header=$1
  awk -v header="$header" '
    $0 == header { in_section = 1; next }
    in_section && /^## \[/ { exit }
    in_section {
      line = $0
      gsub(/<!--.*-->/, "", line)
      if (line ~ /^[[:space:]]*$/ || line ~ /^#/ || line ~ /^\[[^]]+\]:/ || line ~ /^[[:space:]]*[-*+][[:space:]]*$/) next
      found = 1; exit
    }
    END { exit(found ? 0 : 1) }
  ' "$changelog"
}

# Duplicate version sections make extraction and idempotent reruns ambiguous.
section_count=$(awk -v prefix="## [${version}]" 'index($0, prefix) == 1 { n++ } END { print n+0 }' "$changelog")
if (( section_count > 1 )); then
  echo "duplicate changelog sections for ${version}" >&2
  exit 1
fi

if is_finalized; then
  final_header=$(awk -v prefix="$final_prefix" 'index($0, prefix) == 1 { print; exit }' "$changelog")
  if ! section_has_content "$final_header"; then
    echo "finalized changelog section is empty" >&2
    exit 1
  fi
  echo "changelog section [${version}] is finalized"
  exit 0
fi

if [[ "$require_final" == true ]]; then
  echo "changelog section [${version}] is not finalized" >&2
  exit 1
fi

source_section=
if is_planned && section_has_content "$planned_header"; then
  source_section=planned
elif grep -qFx '## [Unreleased]' "$changelog" && section_has_content '## [Unreleased]'; then
  source_section=unreleased
fi

if [[ -z "$source_section" ]]; then
  echo "no non-empty planned [${version}] or [Unreleased] section in $changelog" >&2
  exit 1
fi

if [[ "$check_only" == true ]]; then
  echo "${source_section} changelog section is ready for ${version}"
  exit 0
fi

temporary=$(mktemp "${TMPDIR:-/tmp}/validation-changelog.XXXXXX")
cleanup() {
  rm -f "$temporary"
}
trap cleanup EXIT INT HUP TERM

if [[ "$source_section" == planned ]]; then
  awk -v planned="$planned_header" -v final="${final_prefix}${release_date}" '
    $0 == planned { print final; next }
    { print }
  ' "$changelog" >"$temporary"
else
  awk -v final="${final_prefix}${release_date}" '
    $0 == "## [Unreleased]" {
      print "## [Unreleased]"
      print ""
      print final
      next
    }
    { print }
  ' "$changelog" >"$temporary"
fi

mv "$temporary" "$changelog"
trap - EXIT INT HUP TERM
echo "finalized changelog section [${version}] - ${release_date}"
