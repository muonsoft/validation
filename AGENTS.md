# Agent Guide for Validation Library

This document provides guidance for AI coding agents working with this Go validation library.

This repository is `github.com/muonsoft/validation`.

## Sources of truth

- **README.md** — public API and usage.
- **CHANGELOG.md** — Keep a Changelog record of user-visible changes.
- **docs/release-checklist.md** — release procedure and verification.
- **`.agents/skills/`** — agent skills (for example `validation-add-constraint`, `golang-code-review-comments`).

## Release policy

- Update `CHANGELOG.md` `[Unreleased]` when behavior or public docs change.
- Local agents and scripts never create or push release tags. The maintainer-dispatched **Release** workflow is the only path authorized to push its changelog-only commit and create a release tag.

## Work discipline

- Preserve unrelated user changes.
- Keep changes focused and commits atomic.
- Run checks proportional to the change (`gofmt`, `go test`, `go test -race`) and `bash scripts/test-all.sh` before publication.

## Project Overview

This is a comprehensive Go validation library that provides:

- Declarative validation using typed arguments and constraints in Go code
- Chainable validation constraints via `it` package
- Conditional checks via `is` package
- Custom error messages and translations
- Validation groups for different contexts
- Support for nested structures and collections

## Key Architecture

### Core Components

- `validation.go` - Main validation entry points and `Validatable` interface
- `validator.go` - Core validation logic and execution
- `contraint.go` - Typed constraint interfaces and function adapters
- `it/` - Constraint builders for assertions (e.g., `it.IsEmail()`, `it.HasMinLength()`)
- `is/` - Boolean check functions (e.g., `is.Email()`, `is.URL()`)
- `validate/` - Standalone validation functions
- `violations.go` - Validation error handling
- `message/` - Message templating and translation system

### Validation Flow

1. User calls `validator.Validate(ctx, arguments...)` or an instance's `Validate` method
2. Typed arguments execute explicit constraints; `validation.Valid` calls the
   value's `Validate(ctx context.Context, validator *validation.Validator) error` method.
   The library does not discover validation rules from struct tags
3. Constraints execute and collect violations
4. Violations are formatted using message templates

## Common Tasks

### Adding New Constraints

When adding new validation constraints, follow `.agents/skills/validation-add-constraint/SKILL.md` (and
`golang-code-review-comments` for exported naming).

1. **Add boolean check to `is/` package** (if needed)
  - Pure functions returning `bool`
  - No error handling, just true/false
2. **Add constraint builder to `it/` package**
  - Implements the appropriate typed constraint interface (for example, `validation.StringConstraint`)
  - Uses corresponding `is/` function
  - Defines violation message template
3. **Add tests in `test/constraints_*_cases_test.go`**
  - Use table-driven tests
  - Test both valid and invalid cases
  - Include edge cases
4. **Add examples** in relevant `example_*_test.go` files

### Writing Tests

This project uses table-driven tests extensively:

```go
func TestConstraintName(t *testing.T) {
    cases := []struct {
        name      string
        value     string
        constraint validation.StringConstraint
        wantErr   bool
    }{
        // test cases here
    }
    
    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            // test implementation
        })
    }
}
```

**Black-box tests next to package code** (`is/`, `it/`, `validate/`, root `validation` package, etc.):

- Prefer **black-box testing**: put tests in `*_test.go` files that declare `package foo_test` (the package name must end with the `_test` suffix), and import `github.com/muonsoft/validation/foo` to exercise only the **exported** API.
- Use `package foo` in the same directory only when the test must call **unexported** identifiers (unusual); keep those cases minimal and justified.

Integration-style constraint suites stay in `test/` as today (`test/constraints_*_cases_test.go`).

### Message Templates

When defining new constraints, use message templates:

```go
validation.NewError(
    "your constraint unique code",
    "{{ label }} must satisfy your constraint",
)
```

Add translations in:

- `message/translations/english/messages.go`
- `message/translations/russian/messages.go`

## Code Style Guidelines

1. **Naming Conventions**
  - Use clear, descriptive names
  - Constraint builders: `it.IsXxx()`, `it.HasXxx()`, `it.Xxx()`
  - Check functions: `is.Xxx()`
2. **Documentation**
  - All exported functions must have godoc comments
  - Include examples in `example_*_test.go` files
  - Use `// Output:` comments for testable examples
3. **Error Handling**
  - Use `violations.go` types for validation errors
  - Preserve constraint paths for nested validations
  - Provide clear, actionable error messages
4. **Testing**
  - Aim for high test coverage
  - Use table-driven tests
  - Test edge cases and error conditions
  - Include benchmarks for performance-critical code
  - For unit tests in `is/`, `it/`, `validate/`, etc., use `package foo_test` black-box tests (see [Writing Tests](#writing-tests))

## File Organization

When adding new functionality:

- **Constraints** → `it/` package
- **Check functions** → `is/` package  
- **Standalone validators** → `validate/` package
- **Tests** → shared constraint suites in `test/`; package unit tests in `*_test.go` with `package foo_test` where applicable
- **Examples** → `example_*_test.go` in root
- **Messages** → `message/translations/`

## Common Patterns

### Implementing Validatable

```go
type User struct {
    Email string
    Age   int
}

func (u User) Validate(ctx context.Context, v *validation.Validator) error {
    return v.Validate(ctx,
        validation.String(u.Email, it.IsEmail()),
        validation.Number(u.Age, it.IsGreaterThanOrEqual(18)),
    )
}
```

### Custom Constraints

```go
var ErrNotNumeric = errors.New("not numeric")

type NumericConstraint struct {
    matcher *regexp.Regexp
}

// it is recommended to use semantic constructors for constraints.
func IsNumeric() NumericConstraint {
    return NumericConstraint{matcher: regexp.MustCompile("^[0-9]+$")}
}

func (c NumericConstraint) ValidateString(ctx context.Context, validator *validation.Validator, value *string) error {
    // usually, you should ignore empty values
    // to check for an empty value you should use it.NotBlankConstraint
    if value == nil || *value == "" {
        return nil
    }
    
    if c.matcher.MatchString(*value) {
        return nil
    }
    
    // use the validator to build violation with translations
    return validator.CreateViolation(ctx, ErrNotNumeric, "This value should be numeric.")
}
```

## Before Submitting Changes

1. Run `bash scripts/test-all.sh` (or at minimum `gofmt`, `go vet ./...`, `golangci-lint run`, `go test -race ./...`, `go mod tidy -diff`, `go mod verify`)
2. Check test coverage when adding substantial logic
3. Update documentation if adding public APIs
4. Add examples for new features
5. Ensure CI would pass (see `.github/workflows/tests.yml`)
6. Update **CHANGELOG.md** when the change is user-visible (see [Changelog](#changelog) below)

## Changelog

The project uses **[Keep a Changelog](https://keepachangelog.com/)** in `CHANGELOG.md`.

### Rules for agents and contributors

1. **Always edit `CHANGELOG.md`** in the same branch/PR when your work would matter to library users: new or changed public API, new constraints, message or translation changes, behavior changes, deprecations, removals, security fixes, or notable bug fixes.
2. **Use the `[Unreleased]` section** at the top. Maintainers move entries under a version heading and date when they cut a release.
3. **Pick the right subsection**: `Added`, `Changed`, `Deprecated`, `Removed`, `Fixed`, `Security`. Use `Breaking` only for incompatible changes (or describe breaking impact under `Changed` if you prefer a single list).
4. **Write for consumers**: short, imperative bullets; mention package or symbol names (`it.IsXxx`, `validate.Xxx`, `validation.ErrXxx`) when helpful. Linking to PRs/issues is optional.
5. **Do not rewrite published versions**: never change the bullet list under a released version tag except to fix obvious typos or incorrect facts.
6. **Skip the changelog** only for internal-only changes (refactors, tests, CI, comments) with no user-visible effect.
7. **Release links**: when adding a new version section, update the comparison links at the bottom of the file (`[Unreleased]: ...compare/vX.Y.Z...HEAD` and `[X.Y.Z]: ...releases/tag/vX.Y.Z`).

## Resources

- **CHANGELOG.md** - Release history for users and upgraders
- **README.md** - User-facing documentation
- **docs/release-checklist.md** - Release workflow and verification
- **CONTRIBUTING.md** - Contribution guidelines
- **CODE_OF_CONDUCT.md** - Community standards
- **pkg.go.dev** - Auto-generated API documentation

## Cursor Cloud specific instructions

This is a pure Go library with no external services or infrastructure dependencies. The entire dev workflow is:

- **Full gate:** `bash scripts/test-all.sh` (requires Python 3 for documentation checks;
  install golangci-lint before publication because the script skips lint when absent)
- **Lint:** `golangci-lint run` (requires `golangci-lint` v2 on `PATH`; installed to `$(go env GOPATH)/bin`)
- **Build:** `go build ./...`

### Caveats

- `golangci-lint` is installed to `$(go env GOPATH)/bin` (typically `$HOME/go/bin`). Ensure this is on `PATH`; one-time install: `go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.13.2`.
- The shared CI workflow (`.github/workflows/verify.yml`), called by Tests and Release,
  uses Go **1.26.6** for quality/security and **1.24.0** from `go.mod` for minimum compatibility
  with `GOTOOLCHAIN: local`. It pins `golangci-lint` at **v2.13.2** (`gomodguard_v2`)
  and runs documentation checks, gofmt, vet, lint, race tests, `go mod tidy -diff`, and `go mod verify`.
- The `.golangci.yml` uses config **version: "2"** (golangci-lint v2 format). Do not use golangci-lint v1.
- No Makefile, Docker, or docker-compose is used. No services need to be started.
- This is a library, not a runnable server. Verify behavior with `go test -run '^Example' .` (see `example_*_test.go`).

