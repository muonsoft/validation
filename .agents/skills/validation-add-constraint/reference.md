# Reference: Translations and Test Cases

## Translations

### Keys and files

- **Key**: Always the constant from `message` package (e.g. `message.InvalidUUID`). The same key is used in English and Russian maps.
- **English**: `message/translations/english/messages.go` — `var Messages = map[language.Tag]map[string]catalog.Message{ language.English: { ... } }`
- **Russian**: `message/translations/russian/messages.go` — same structure with `language.Russian`.

### Simple message (no plural)

```go
message.InvalidUUID: catalog.String(message.InvalidUUID),           // English: reuse default
message.InvalidUUID: catalog.String("Значение не соответствует формату UUID."), // Russian
```

### Plural (English)

Use `plural.Selectf(1, "", plural.One, "...", plural.Other, "...")`. The first argument is the placeholder index for the number (1 = first `{{ limit }}` etc.):

```go
message.TooShort: plural.Selectf(1, "",
	plural.One, "This value is too short. It should have {{ limit }} character or more.",
	plural.Other, "This value is too short. It should have {{ limit }} characters or more."),
```

### Plural (Russian)

Russian has One, Few, Other (and optionally Many). Use the same key and `plural.Selectf(1, "", ...)`:

```go
message.TooShort: plural.Selectf(1, "",
	plural.One, "Значение слишком короткое. Должно быть равно {{ limit }} символу или больше.",
	plural.Few, "Значение слишком короткое. Должно быть равно {{ limit }} символам или больше.",
	plural.Other, "Значение слишком короткое. Должно быть равно {{ limit }} символам или больше."),
```

### Placeholders

- Common: `{{ value }}`, `{{ limit }}`, `{{ comparedValue }}`, `{{ min }}`, `{{ max }}`, `{{ divisibleBy }}`.
- Keep placeholder names identical to those used in `message/messages.go` and in constraint `WithMessage`/`WithParameters`.

---

## Test cases (test package)

### ConstraintValidationTestCase

Defined in `test/constraints_test.go`. Relevant fields:

- `name`: test name (e.g. `"IsUUID passes on valid value"`).
- `isApplicableFor`: `specificValueTypes(stringType)` or `specificValueTypes(intType)`, etc.
- `stringValue`: `stringValue("")` or `stringValue("valid")` (helper in `test/mocks_test.go`).
- `constraint`: e.g. `it.IsUUID()`, `it.IsUUID().NotNil()`.
- `assert`: `assertNoError` or `assertHasOneViolation(validation.ErrInvalidUUID, message.InvalidUUID)`.

For custom message assertion use the exact expected message string:

```go
assert: assertHasOneViolation(ErrCustom, `Invalid value "invalid" for parameter.`),
```

### Helpers

- `stringValue(s string) *string` — in `test/mocks_test.go`.
- `specificValueTypes(types ...string) func(string) bool` — in `test/constraints_test.go`; use `stringType`, `intType`, `floatType`, `boolType`, `timeType`, `stringsType`, `iterableType`, `countableType`, `comparableType`, `nilType`.
- `assertNoError`, `assertHasOneViolation(err, message)` — in `test/assertions_test.go`.

### Merging into validateTestCases

In `test/constraints_test.go`, `validateTestCases` is built with `mergeTestCases(...)`. Add your slice:

```go
var validateTestCases = mergeTestCases(
	// ...
	myConstraintTestCases,
)
```

The runners `TestValidateString`, `TestValidateNilString`, etc. iterate over `validateTestCases` and filter by `isApplicableFor`; they call `newValidator(t).Validate(...)` and then `test.assert(t, err)`.

### New constraint file

Create or use a file like `test/constraints_identifiers_cases_test.go`:

```go
package test

import (
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/message"
)

var myConstraintTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsMy passes on valid value",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("valid"),
		constraint:      it.IsMy(),
		assert:          assertNoError,
	},
	{
		name:            "IsMy violation on invalid value",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("invalid"),
		constraint:      it.IsMy(),
		assert:          assertHasOneViolation(validation.ErrInvalidMy, message.InvalidMy),
	},
	// empty, nil, options, When(false), WithError/WithMessage...
}
```

Then add `myConstraintTestCases` to `mergeTestCases` in `test/constraints_test.go`.
