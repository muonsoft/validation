# Skills for Validation Library

This document outlines common skills and workflows for working with this Go validation library.

## Table of Contents

- [Adding New Validation Constraints](#adding-new-validation-constraints)
- [Writing Idiomatic Tests](#writing-idiomatic-tests)
- [Internationalization Workflow](#internationalization-workflow)
- [Performance Optimization](#performance-optimization)
- [Creating Examples](#creating-examples)

---

## Adding New Validation Constraints

### Workflow

1. **Identify the constraint category:**
   - String validation → `it/string.go` and `is/data.go`
   - Numeric validation → `it/comparison.go` and `is/numeric.go`
   - Web/URL validation → `it/web.go` and `is/web.go`
   - Date/time validation → `it/date_time.go`
   - Barcode/identifier → `it/identifiers.go` and `is/identifiers.go`
   - Iterable/collection → `it/iterable.go`

2. **Add check function to `is/` package:**

```go
// Pure boolean check
func YourCheck(value any) bool {
    str, ok := value.(string)
    if !ok {
        return false
    }
    // implement validation logic
    return true
}
```

3. **Add constraint builder to `it/` package:**

```go
// YourConstraint validates that value satisfies your check.
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

4. **Add message translations:**
   - English: `message/translations/english/messages.go`
   - Russian: `message/translations/russian/messages.go`

5. **Write tests in `test/`:**
   - Add test cases to appropriate `constraints_*_cases_test.go`
   - Follow table-driven test pattern

6. **Create usage example:**
   - Add to relevant `example_*_test.go` file
   - Include `// Output:` for testable examples

---

## Writing Idiomatic Tests

### Table-Driven Test Pattern

```go
func TestYourConstraint(t *testing.T) {
    t.Parallel()
    
    cases := []struct {
        name       string
        value      any
        constraint validation.Constraint
        wantErr    bool
        wantMsg    string
    }{
        {
            name:       "valid case",
            value:      "valid",
            constraint: it.YourConstraint(),
            wantErr:    false,
        },
        {
            name:       "invalid case",
            value:      "invalid",
            constraint: it.YourConstraint(),
            wantErr:    true,
            wantMsg:    "must satisfy your constraint",
        },
        {
            name:       "blank allowed",
            value:      "",
            constraint: it.YourConstraint(),
            wantErr:    false,
        },
    }
    
    for _, tc := range cases {
        tc := tc // capture range variable
        t.Run(tc.name, func(t *testing.T) {
            t.Parallel()
            
            err := newValidator(t).Validate(context.Background(), tc.value, tc.constraint)
            
            if tc.wantErr {
                assert.Error(t, err)
                if tc.wantMsg != "" {
                    assert.Contains(t, err.Error(), tc.wantMsg)
                }
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

### Benchmark Tests

Add benchmarks for performance-critical constraints:

```go
func BenchmarkYourConstraint(b *testing.B) {
    constraint := it.YourConstraint()
    value := "test value"
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = newValidator(b).Validate(value, constraint)
    }
}
```

---

## Internationalization Workflow

### Adding New Message Keys

1. **Define message template in constraint:**

`message/messages.go`:
```go
const (
    NotValid = "{{label}} is not a valid."
)
```

2. **Add English translation:**

`message/translations/english/messages.go`:
```go
var Messages = map[language.Tag]map[string]catalog.Message{
     language.English: {
        message.NotValid: catalog.String("{{label}} is not a valid."),
    }
}
```

3. **Add Russian translation:**

`message/translations/russian/messages.go`:
```go
var Messages = map[language.Tag]map[string]catalog.Message{
    language.Russian: {
        message.NotValid: catalog.String("Значение {{label}} невалидно."),
    }
}
```

4. **Test both languages:**

```go
func TestYourConstraintTranslation(t *testing.T) {
    // Test English
    validator := newValidator(
        t,
        validation.DefaultLanguage(language.English),
        validation.Translations(english.Messages),
    )
    validator.ValidateString(
        context.Background(),
        validation.String("", it.YourConstraint()),
    )
    
    // Test Russian
    validator := newValidator(
        t,
        validation.DefaultLanguage(language.Russian),
        validation.Translations(russian.Messages),
    )
    validator.ValidateString(
        context.Background(),
        validation.String("", it.YourConstraint()),
    )
}
```

---

## Performance Optimization

### Guidelines

1. **Avoid allocations in hot paths:**
   - Reuse buffers where possible
   - Use value receivers for small types
   - Minimize string concatenation

2. **Lazy evaluation:**
   - Don't compute violation messages until needed
   - Use message templates instead of preformatted strings

3. **Early returns:**
   - Check for blank/nil values first
   - Return early on first violation when appropriate

4. **Benchmark before and after:**

```go
func BenchmarkConstraintBefore(b *testing.B) {
    // baseline implementation
}

func BenchmarkConstraintAfter(b *testing.B) {
    // optimized implementation
}
```

5. **Profile if needed:**
```bash
go test -cpuprofile=cpu.prof -memprofile=mem.prof -bench=.
go tool pprof cpu.prof
```

---

## Creating Examples

### Example Test Pattern

Place examples in `example_*_test.go` files:

```go
func ExampleYourConstraint() {
    type User struct {
        Field string
    }
    
    user := User{Field: "value"}

    validator := validation.NewValidator()
    err := validator.Validate(
        context.Background(),
        validation.String(user.Field, it.YourConstraint()),
    )
    
    if err != nil {
        fmt.Println(err)
    }
    // Output:
    // value must satisfy your constraint
}
```

### Example Guidelines

1. **Keep examples simple and focused**
   - Demonstrate one concept per example
   - Use realistic but minimal data

2. **Include expected output**
   - Use `// Output:` comment for testable examples
   - Show both success and error cases

3. **Name examples descriptively**
   - `ExampleConstraintName`
   - `ExampleConstraintName_withOption`
   - `ExampleConstraintName_customMessage`

4. **Test examples automatically**
   - Examples with `// Output:` are tested by `go test`
   - Keep output deterministic

---

## Quick Reference

### Running Tests

```bash
# All tests
go test ./...

# Specific package
go test ./test

# With coverage
go test -cover ./...

# Verbose output
go test -v ./...

# Run specific test
go test -run TestYourConstraint ./test
```

### Running Linter

```bash
golangci-lint run

# Auto-fix issues
golangci-lint run --fix
```

### Generating Documentation

```bash
# View docs locally
godoc -http=:6060

# Then visit http://localhost:6060/pkg/github.com/yourusername/validation/
```

### Common Git Workflows

```bash
# Create feature branch
git checkout -b feature/your-constraint

# Run tests before commit
go test ./...
golangci-lint run

# Commit changes
git add .
git commit -m "feat: add YourConstraint validation"

# Push and create PR
git push origin feature/your-constraint
```
