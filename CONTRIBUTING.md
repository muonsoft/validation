# Contributing

For substantial changes, discuss the proposal in an
[issue](https://github.com/muonsoft/validation/issues) before implementation.
Follow the [code of conduct](CODE_OF_CONDUCT.md) in project discussions.

## Development setup

The library has no external services. Use Go 1.24.0 or later; the quality CI job
uses Go 1.26.6. Install Python 3 for documentation and release-script checks,
and golangci-lint v2.13.2 for linting. Use a toolchain compatible with the linter.

```bash
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.13.2
export PATH="$(go env GOPATH)/bin:$PATH"
```

## Before opening a pull request

Keep changes focused and preserve unrelated work. Format Go changes with `gofmt`,
then run the local gate:

```bash
bash scripts/test-all.sh
git diff --check
```

The gate checks documentation, formatting, vet, unit and race tests, module
consistency, and lint **when golangci-lint is on PATH**. A local run that skips
lint is incomplete for publication; CI requires it. Shared CI also runs security
and minimum-Go checks. See [.github/workflows/verify.yml](.github/workflows/verify.yml).

For release tooling changes, also run:

```bash
bash scripts/prepare-release-test.sh
python3 scripts/release-candidate-test.py
python3 scripts/security-report-test.py
```

## Tests and documentation

Use black-box package tests (`package foo_test`) for public APIs. Shared constraint
suites belong in `test/`. Include runnable `Example` functions with `// Output:`
for new public APIs and link relevant guides to those examples.

Run examples alone with `go test -run '^Example' ./...`. The documentation check
(`python3 scripts/check-docs.py`) checks local Markdown links and executes complete
Go programs in Markdown, comparing any `// Output:` block with their output.
Keep partial snippets clearly contextualized; link them to executable examples.

Update README or the appropriate guide when public behavior changes. Add a concise
consumer-facing entry under `[Unreleased]` in [CHANGELOG.md](CHANGELOG.md), including
public documentation corrections. Use Added, Changed, Deprecated, Removed, Fixed,
or Security as appropriate. Do not rewrite published release notes.

## Review and releases

Describe the problem, resulting behavior, and checks performed in your pull request.
Obtain maintainer review before merging.

Only the maintainer-dispatched **Release** workflow publishes releases and creates
tags. Contributors and local scripts must not create or push release tags. Follow
the [release checklist](docs/release-checklist.md) for validation and publication.
