# Agent Guide for Validation Library

This document provides guidance for AI coding agents working with this Go validation library.

## Project Overview

This is a comprehensive Go validation library that provides:
- Declarative validation using struct tags
- Chainable validation constraints via `it` package
- Conditional checks via `is` package
- Custom error messages and translations
- Validation groups for different contexts
- Support for nested structures and collections

## Key Architecture

### Core Components

- **`validation.go`** - Main validation entry points and `Validatable` interface
- **`validator.go`** - Core validation logic and execution
- **`constraints.go`** - Constraint interface and execution
- **`it/`** - Constraint builders for assertions (e.g., `it.IsEmail()`, `it.MinLength()`)
- **`is/`** - Boolean check functions (e.g., `is.Email()`, `is.URL()`)
- **`validate/`** - Standalone validation functions
- **`violations.go`** - Validation error handling
- **`message/`** - Message templating and translation system

### Validation Flow

1. User calls `validation.Validate(value)` or implements `Validatable` interface
2. Validator discovers constraints from struct tags or `Validate()` methods
3. Constraints execute and collect violations
4. Violations are formatted using message templates

## Common Tasks

### Adding New Constraints

When adding new validation constraints:

1. **Add boolean check to `is/` package** (if needed)
   - Pure functions returning `bool`
   - No error handling, just true/false
   
2. **Add constraint builder to `it/` package**
   - Returns a `validation.Constraint`
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
        value     any
        constraint validation.Constraint
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

## File Organization

When adding new functionality:

- **Constraints** → `it/` package
- **Check functions** → `is/` package  
- **Standalone validators** → `validate/` package
- **Tests** → `test/` directory
- **Examples** → `example_*_test.go` in root
- **Messages** → `message/translations/`

## Common Patterns

### Implementing Validatable

```go
type User struct {
    Email string
    Age   int
}

func (u User) Validate() error {
    return validation.ValidateValue(
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

1. Run tests: `go test ./...`
2. Run linter: `golangci-lint run`
3. Check test coverage
4. Update documentation if adding public APIs
5. Add examples for new features
6. Ensure all tests pass in CI

## Resources

- **README.md** - User-facing documentation
- **CONTRIBUTING.md** - Contribution guidelines
- **CODE_OF_CONDUCT.md** - Community standards
- **pkg.go.dev** - Auto-generated API documentation

## Cursor Cloud specific instructions

This is a pure Go library with no external services or infrastructure dependencies. The entire dev workflow is:

- **Install deps:** `go mod download`
- **Lint:** `golangci-lint run` (requires `golangci-lint` v2 on `PATH`; installed to `$(go env GOPATH)/bin`)
- **Test:** `go test -race ./...`
- **Build:** `go build ./...`

### Caveats

- `golangci-lint` is installed to `$(go env GOPATH)/bin`. Ensure this is on `PATH` (the VM's `~/.bashrc` exports it).
- The CI workflow (`.github/workflows/tests.yml`) pins `golangci-lint` at **v2.6.1** and Go at **^1.24**. Match these versions locally.
- The `.golangci.yml` uses config **version: "2"** (golangci-lint v2 format). Do not use golangci-lint v1.
- No Makefile, Docker, or docker-compose is used. No services need to be started.
