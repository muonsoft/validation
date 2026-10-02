# How to Use

## Basic concepts

Pass typed arguments such as `validation.String` to `validator.Validate`, with
constraints from `it`. The examples below use the `validator` package unless a
local validator is explicitly constructed. See [Installation](installation.md)
for imports and the [README](../README.md#basic-example) for a complete program.

| Package | Use |
|---------|-----|
| `validation` | Typed arguments, interfaces, errors, and configurable validator instances |
| `validator` | A default validator and convenience functions |
| `it` | Composable constraints with messages, conditions, and groups |
| `validate` | Standalone checks returning `error` |
| `is` | Standalone checks returning `bool` |

```go
err := validator.Validate(context.Background(), validation.String("", it.IsNotBlank()))

fmt.Println(err)
// Output:
// violation: "This value should not be blank."
```

List of common [validation arguments](https://pkg.go.dev/github.com/muonsoft/validation#Argument):

* `validation.Nil()` - passes result of comparison to nil to test against nil constraints;
* `validation.Bool()` - passes boolean value;
* `validation.NilBool()` - passes boolean pointer value;
* `validation.Number[T]()` - passes generic numeric value;
* `validation.NilNumber[T]()` - passes generic numeric pointer value;
* `validation.String()` - passes string value;
* `validation.NilString()` - passes string pointer value;
* `validation.Countable()` - passes result of `len()` to test against constraints based on count of the elements;
* `validation.Time()` - passes `time.Time` value;
* `validation.NilTime()` - passes `time.Time` pointer value;
* `validation.EachNumber[T]()` - passes slice of generic numbers to test each of the element against numeric constraints;
* `validation.EachString()` - passes slice of strings to test each of the element against string constraints;
* `validation.Valid()` - passes `Validatable` value to run embedded validation;
* `validation.ValidSlice[T]()` - passes `[]T` where `T` implements `Validatable` to run embedded validation on each of the elements;
* `validation.ValidMap[T]()` - passes `map[string]T` where `T` implements `Validatable` to run embedded validation on each of the elements;
* `validation.Comparable[T]()` - passes generic comparable value to test against comparable constraints;
* `validation.NilComparable[T]()` - passes generic comparable pointer value to test against comparable constraints;
* `validation.Comparables[T]()` - passes generic slice of comparable values (can be used to check for uniqueness of the elements);
* `validation.Slice[T]()` - passes generic slice to test against `SliceConstraint[T]` list (e.g. uniqueness by key via `it.HasUniqueValuesBy()`);
* `validation.SliceProperty[T]()` - same as `Slice` with property name in violation path;
* `validation.Check()` - passes result of any boolean expression;
* `validation.This[T]()` - passes any value to `Constraint[T]` implementations;
* `validation.Each[T]()` - validates each slice element with `Constraint[T]` implementations;
* `validation.CheckNoViolations()` - accepts errors from embedded validation, collects violations, and propagates other errors.

These methods combine `Validate(ctx, ...)` with the corresponding typed argument:

* `validator.ValidateBool()` - shorthand for `validation.Bool()`;
* `validator.ValidateInt()` - shorthand for `validation.Number[int]()`;
* `validator.ValidateFloat()` - shorthand for `validation.Number[float64]()`;
* `validator.ValidateString()` - shorthand for `validation.String()`;
* `validator.ValidateStrings()` - shorthand for `validation.Comparables[string]()`;
* `validator.ValidateCountable()` - shorthand for `validation.Countable()`;
* `validator.ValidateTime()` - shorthand for `validation.Time()`;
* `validator.ValidateEachString()` - shorthand for `validation.EachString()`;
* `validator.ValidateIt()` - shorthand for `validation.Valid()`.

See usage examples in the [documentation](https://pkg.go.dev/github.com/muonsoft/validation#Validator.Validate).

## Empty and required values

Most string format constraints accept `""` and skip nil pointers. Combine them
with `it.IsNotBlank()` when a field is required:

```go
err := validator.Validate(ctx,
    validation.StringProperty("email", email, it.IsNotBlank(), it.IsEmail()),
)
```

Use `validation.NilString` for `*string` values; `String` accepts a `string`.
`it.IsNotNil()` requires a non-nil pointer while allowing an empty string.
`it.IsNotBlank()` rejects nil and empty strings by default; whitespace-only
strings pass. Normalize with `strings.TrimSpace` before validation if your
application treats whitespace as empty. `WithAllowedNil()` permits nil without
permitting an empty string.

Empty-value behavior belongs to each constraint. For example,
`it.HasPasswordStrength()` evaluates empty strings, and `it.HasMinCount(1)` rejects
a zero count. Constraints that validate their options may return configuration
errors even for empty input. Consult the [constraint catalog](constraints.md)
and each rule's documentation.

Validation normally collects all violations. Adding `IsNotBlank` does not stop
later rules from running. Use `validation.Sequentially` to stop after an argument
produces violations. Non-violation errors stop validation; see
[Violations and errors](violations-and-errors.md).

## Configure a validator

Create an instance with `validation.NewValidator`. For example, to use Russian
translations by default:

```go
// Also import "golang.org/x/text/language" and
// "github.com/muonsoft/validation/message/translations/russian".
v, err := validation.NewValidator(
    validation.DefaultLanguage(language.Russian),
    validation.Translations(russian.Messages),
)
if err != nil {
    return err
}
return v.Validate(ctx, validation.String("", it.IsNotBlank()))
```

The `return` statements above assume an enclosing function returning `error`.
Options also support a custom translator or violation factory; see
[Translations](translations.md) and [Violations and errors](violations-and-errors.md).

To configure the default validator, create the instance during application
initialization and install it after checking the constructor error:

```go
v, err := validation.NewValidator(validation.Translations(russian.Messages))
if err != nil {
    return err
}
validator.SetDefault(v)
```

`validator.Default()` returns the current instance. `SetDefault` replaces it;
call it once during initialization. `SetUp` and `Instance` are deprecated.
