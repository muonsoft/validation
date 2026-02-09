---
name: validation-add-constraint
description: Adds new validation constraints to the Go validation library. Use when adding a new constraint to package it, translations (message, english, russian), optional validate/is helpers, tests in test/, and examples in it/example_test.go or validate/example_test.go.
---

# Adding a New Constraint to the Validation Library

Follow this workflow when adding a new constraint. Mandatory steps: constraint in `it`, message constant, translations (english + russian), tests, examples, godoc. Optional: `validate` and `is` when useful for standalone validation (e.g. string codes, identifiers).

---

## 1. Message and Error (mandatory)

### 1.1 Add message constant

In **`message/messages.go`** add a new constant (English default text):

```go
const (
	// ... existing
	InvalidMyFormat = "This value is not a valid my format."
)
```

- Use existing style: `InvalidXxx`, `NotXxx`, `TooXxx`, etc.
- Text is the default (English) template; placeholders like `{{ value }}`, `{{ limit }}` are allowed.

### 1.2 Add validation error

In **`errors.go`** (root package) add:

```go
var (
	// ... existing
	ErrInvalidMyFormat = NewError("invalid my format", message.InvalidMyFormat)
)
```

- First argument: stable **code** (backward compatibility).
- Second argument: **message constant** from `message` package (used for translation key and default text).

---

## 2. Translations (mandatory)

Add the **same key** (the message constant) in both translation files.

### 2.1 English

In **`message/translations/english/messages.go`**:

```go
var Messages = map[language.Tag]map[string]catalog.Message{
	language.English: {
		// ... existing
		message.InvalidMyFormat: catalog.String(message.InvalidMyFormat),
	},
}
```

- For simple messages use `catalog.String(message.Const)` or `catalog.String("Your English text.")`.
- For plurals (e.g. "N elements") use `plural.Selectf(1, "", plural.One, "...", plural.Other, "...")`. See [reference.md](reference.md) for plural and Russian forms.

### 2.2 Russian

In **`message/translations/russian/messages.go`**:

```go
var Messages = map[language.Tag]map[string]catalog.Message{
	language.Russian: {
		// ... existing
		message.InvalidMyFormat: catalog.String("Значение не является допустимым форматом."),
	},
}
```

- Key is always the **message constant** from `message` package.
- Value: Russian text; same placeholders as in the message constant (e.g. `{{ value }}`).
- For plurals use `plural.Selectf(1, "", plural.One, "...", plural.Few, "...", plural.Other, "...")` — Russian has Few. See [reference.md](reference.md).

---

## 3. Constraint in package `it` (mandatory)

Choose the right file: `it/string.go`, `it/identifiers.go`, `it/web.go`, `it/comparison.go`, `it/basic.go`, `it/iterable.go`, `it/date_time.go`, `it/choice.go`, `it/barcodes.go`.

### 3.1 Simple string constraint (func(string) bool)

If the check is a pure `func(string) bool`, use `OfStringBy` and the `is` helper:

```go
// IsMyFormat validates whether the value is in my format.
// See [link] for specification.
func IsMyFormat() validation.StringFuncConstraint {
	return validation.OfStringBy(is.MyFormat).
		WithError(validation.ErrInvalidMyFormat).
		WithMessage(validation.ErrInvalidMyFormat.Message())
}
```

### 3.2 Custom struct constraint

When you need options (e.g. versions, formats), define a struct and implement `ValidateString`:

```go
// MyConstraint validates whether the string value satisfies my format.
// Use [MyConstraint.Option] to configure.
type MyConstraint struct {
	isIgnored       bool
	groups          []string
	options         []func(o *validate.MyOptions)
	err             error
	messageTemplate string
	messageParameters validation.TemplateParameterList
}

// IsMy creates the constraint.
func IsMy() MyConstraint {
	return MyConstraint{
		err:             validation.ErrInvalidMyFormat,
		messageTemplate: validation.ErrInvalidMyFormat.Message(),
	}
}

// WithError overrides default error for produced violation.
func (c MyConstraint) WithError(err error) MyConstraint { ... }

// WithMessage sets the violation message template.
func (c MyConstraint) WithMessage(template string, parameters ...validation.TemplateParameter) MyConstraint { ... }

// When / WhenGroups for conditional validation.
func (c MyConstraint) When(condition bool) MyConstraint { ... }
func (c MyConstraint) WhenGroups(groups ...string) MyConstraint { ... }

func (c MyConstraint) ValidateString(ctx context.Context, validator *validation.Validator, value *string) error {
	if c.isIgnored || validator.IsIgnoredForGroups(c.groups...) || value == nil || *value == "" {
		return nil
	}
	if is.My(*value, c.options...) {
		return nil
	}
	return validator.BuildViolation(ctx, c.err, c.messageTemplate).
		WithParameters(
			c.messageParameters.Prepend(
				validation.TemplateParameter{Key: "{{ value }}", Value: *value},
			)...,
		).
		Create()
}
```

Template rules:

- **Empty/nil**: Usually skip (return `nil`); use `it.IsNotBlank()` (or similar) to reject empty.
- **Violation**: Use `validator.BuildViolation(ctx, c.err, c.messageTemplate).WithParameters(...).Create()`. Do not use `CreateViolation` for translatable constraints — use `BuildViolation` so the message is translated.
- **Godoc**: Document the constraint type and constructor; document options; add `See ...` for specs if applicable.

---

## 4. Tests (mandatory)

In **`test/constraints_*_cases_test.go`** (create or extend the right file, e.g. `constraints_identifiers_cases_test.go`):

- Define a slice of `ConstraintValidationTestCase`.
- Use `name`, `isApplicableFor: specificValueTypes(stringType)` (or other type), `stringValue: stringValue("...")`, `constraint: it.IsMyFormat()`, `assert: assertNoError` or `assertHasOneViolation(validation.ErrInvalidMyFormat, message.InvalidMyFormat)`.
- Cover: valid, invalid, empty/nil (if applicable), options (e.g. WithError/WithMessage), When(false)/When(true).

Add the slice to **`validateTestCases`** in **`test/constraints_test.go`** via `mergeTestCases(...)` so the shared test runners pick it up.

---

## 5. Examples (mandatory)

In **`it/example_test.go`** add testable examples:

```go
func ExampleIsMyFormat_valid() {
	err := validator.Validate(context.Background(), validation.String("valid-value", it.IsMyFormat()))
	fmt.Println(err)
	// Output:
	// <nil>
}

func ExampleIsMyFormat_invalid() {
	err := validator.Validate(context.Background(), validation.String("invalid", it.IsMyFormat()))
	fmt.Println(err)
	// Output:
	// violation: "This value is not a valid my format."
}
```

Use `// Output:` so `go test` runs them. Prefer `ExampleXxx_valid` / `ExampleXxx_invalid` naming.

---

## 6. Optional: package `validate`

Add when the constraint is useful for **standalone validation** (e.g. string codes, identifiers), without the full validator.

- **File**: `validate/identifiers.go` or new file (e.g. `validate/myformat.go`).
- **Signature**: `func MyFormat(value string) error` (or with options).
- **Return**: `nil` if valid; otherwise a sentinel error from `validate` package (e.g. `ErrTooShort`, or custom).
- **Godoc**: Describe when it returns which error.
- **Tests**: `validate/*_test.go` table or cases.
- **Examples**: `validate/example_test.go` with `ExampleMyFormat` and `// Output:`.

If the `it` constraint needs options, define option types and funcs in `validate` (e.g. `validate.MyOptions`, `validate.AllowXxx()`), and use them from `it` and `is`.

---

## 7. Optional: package `is`

Add when useful for **standalone boolean checks** (e.g. in conditions, or for `OfStringBy`).

- **File**: `is/identifiers.go` or same area as related `validate` logic.
- **Signature**: `func MyFormat(value string) bool` or `func MyFormat(value string, options ...func(o *validate.MyOptions)) bool`.
- **Implementation**: Usually `return validate.MyFormat(value, options...) == nil`.
- **Godoc**: Short description; point to `validate` for options and semantics.
- **Tests**: `is/*_test.go`.
- **Examples**: `is/example_test.go` with `ExampleMyFormat` and `// Output:`.

---

## Checklist

- [ ] `message/messages.go`: new constant
- [ ] `errors.go`: `ErrXxx = NewError("code", message.Xxx)`
- [ ] `message/translations/english/messages.go`: key = message const, value = English
- [ ] `message/translations/russian/messages.go`: key = message const, value = Russian
- [ ] `it/*.go`: constraint (StringFuncConstraint or custom struct), godoc, empty/nil handling, BuildViolation
- [ ] `test/constraints_*_cases_test.go`: test cases + merge into `validateTestCases`
- [ ] `it/example_test.go`: ExampleXxx_valid, ExampleXxx_invalid with `// Output:`
- [ ] Optional: `validate`: function, tests, examples, godoc
- [ ] Optional: `is`: function, tests, examples, godoc
- [ ] `go test ./...` and `golangci-lint run`

---

## Additional resources

- Plural forms and Russian translation details: [reference.md](reference.md)
- Existing patterns: `it/identifiers.go` (UUID/ULID), `it/string.go` (OfStringBy), `test/constraints_identifiers_cases_test.go`, `message/translations/`.
