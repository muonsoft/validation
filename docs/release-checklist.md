# Release checklist

Release publication is a maintainer action through the GitHub Actions **Release**
workflow. Local scripts and autonomous agents validate and prepare releases but never
create or push release tags. The source-only GitHub Release creates the tag at the
workflow-verified release commit.

## Preflight

- [ ] Work is merged to `main` and the branch is not moving during publication.
- [ ] Normal CI is green.
- [ ] `CHANGELOG.md` has either a non-empty exact planned version section or non-empty
      `[Unreleased]` section.
- [ ] README and public docs are current.
- [ ] No local or remote tag with the requested version points to another commit.
- [ ] Repository Actions settings allow `GITHUB_TOKEN` write access, and branch rules
      permit `github-actions[bot]` to push the changelog-only release commit.

## Local validation

Run from the repository root:

```bash
bash -n scripts/*.sh
bash scripts/prepare-release-test.sh
bash scripts/prepare-release.sh 0.20.0 --check-only
bash scripts/test-all.sh
git diff --check
git status --short
git tag --list v0.20.0
```

## Dispatch the release

1. Open **Actions → Release → Run workflow** in GitHub.
2. Select `main` and enter a `v`-prefixed SemVer version, for example `v0.20.0`.
3. Confirm the validation job succeeds.
4. Confirm the publish job creates at most one changelog-only commit on `main`.
5. Confirm the source-only GitHub Release and its tag point to that release commit.

The workflow will stop if it was dispatched from a branch other than current `main`,
if `main` moves after validation, if the changelog cannot be finalized, or if the tag
already points to a different commit.

## What the workflow does

1. Validates strict SemVer and records the exact dispatch SHA.
2. Runs changelog script tests and checks that release notes are ready.
3. Runs lint and the full unit/race/format/module gate.
4. Rechecks `origin/main`, finalizes the changelog with the UTC release date, and
   pushes a changelog-only release commit when necessary.
5. Verifies that the release commit is either the validated SHA or its single
   changelog-only child.
6. Extracts release notes and publishes a source-only GitHub Release. GitHub creates
   the missing tag at the verified release commit.

## Post-release verification

```bash
gh release view v0.20.0 --repo muonsoft/validation
git ls-remote --tags origin refs/tags/v0.20.0

tmp_dir=$(mktemp -d)
cd "$tmp_dir"
go mod init release-smoke
GOTOOLCHAIN=local go get github.com/muonsoft/validation@v0.20.0
go list -m github.com/muonsoft/validation
```

After verification, pull the changelog release commit into local clones. Do not create
a second local tag or move the published tag.

## Failure and rollback

- A validation failure creates neither a changelog commit nor a tag.
- A failure before GitHub Release publication may leave a valid changelog-only commit;
  rerun the workflow for the same version after diagnosing the failure.
- If a release or tag was published incorrectly, follow repository governance,
  document a retraction, and publish a corrective version. Do not rewrite public
  history or silently move a consumed tag.

## Shared CI contract

Normal CI and Release both call `.github/workflows/verify.yml` at the caller's
revision. Its required jobs cover the existing project gate, minimum Go from
`go.mod`, and pinned govulncheck v1.8.0 on both minimum and tooling Go. The scanner
is installed with tooling Go, then executed with the selected analysis toolchain.
No vulnerability failures are suppressed. Resolving findings or raising the support
floor is a separate compatibility decision; a newer clean scan cannot replace a
failed minimum lane. A two-component Go directive selects that release family's
available patch; a three-component directive pins the specified minimum patch.
Each job logs the effective toolchain. Tooling remains on the pinned Go 1.26.6 lane.

Project-specific checks remain required, including Squirrel's database, nested
modules, differential/API and dependency gates where present. Modern tooling and
`go mod tidy -diff` do not run in the old minimum-Go compatibility lane.

The publish job prepares from the immutable validated SHA and checks ancestry,
changed files, nonempty finalized notes and tag identity **before** pushing anything.
It rechecks the branch and uses a normal non-force push of that exact commit.
Branch movement rejects publication; no pull/rebase incorporates unverified work.
A subsequent unrelated branch commit does not change the already selected release
SHA. Repository release concurrency does not lock human pushes.

After publication, the `consumer` job resolves the exact version in a temporary
external module, with workspace use disabled and an isolated module cache, and
compiles the public imports listed in `scripts/release-consumer-imports.txt`.
The cache is outside the consumer source directory and writable for cleanup. Download retries are bounded to accommodate proxy propagation;
a failed smoke check does not delete or move the published tag. Because bot pushes
need not trigger push workflows, this check is an explicit dependent job.

## Recovery and rollout

If the changelog commit was pushed but publication failed, inspect the tag/release
first and start a **fresh dispatch** on the current release branch with the same
version. Re-running the old event still has its old SHA and can fail the freshness
check. Existing same-commit tags may be reconciled; conflicting tags stop the flow.
The version input has no prefilled default to avoid reusing a previous release.

Before merging this CI change, review branch protection required-check names:
reusable workflow jobs have new check names. Configure rules to require the new
quality, minimum and both vulnerability lanes; do not disable existing protections
without their replacements. Confirm the release bot can push its changelog commit.
No local validation can establish those hosted repository settings.

Local release regression checks:

```bash
bash scripts/prepare-release-test.sh
python3 scripts/release-candidate-test.py
```

The tests reject heading-only/planned/finalized empty notes, duplicate versions,
non-changelog or multiple child commits, merges and dirty tracked files; validation
modes do not rewrite the changelog. Hosted publication itself still needs a real
maintainer-authorized run after the workflow is merged.
