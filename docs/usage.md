# How to Use

## Basic concepts

The validation process is built around functional options and passing values by specific typed arguments. A common way
to use validation is to call the `validator.Validate` method and pass the argument option with the list of validation
constraints.

```go
err := validator.Validate(context.Background(), validation.String("", it.IsNotBlank()))

fmt.Println(err)
// Output:
// violation: This value should not be blank.
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
* `validation.ValidSlice[T]()` - passes slice of `[]Validatable` value to run embedded validation on each of the elements;
* `validation.ValidMap[T]()` - passes `map[string]Validatable` value to run embedded validation on each of the elements;
* `validation.Comparable[T]()` - passes generic comparable value to test against comparable constraints;
* `validation.NilComparable[T]()` - passes generic comparable pointer value to test against comparable constraints;
* `validation.Comparables[T]()` - passes generic slice of comparable values (can be used to check for uniqueness of the elements);
* `validation.Check()` - passes result of any boolean expression;
* `validation.CheckNoViolations()` - passes `error` to check err for violations, can be used for embedded validation.

For single value validation, you can use shorthand versions of the validation method:

* `validator.ValidateBool()` - shorthand for `validator.Bool()`;
* `validator.ValidateInt()` - shorthand for `validation.Number[int]()`;
* `validator.ValidateFloat()` - shorthand for `validation.Number[float64]()`;
* `validator.ValidateString()` - shorthand for `validation.String()`;
* `validator.ValidateStrings()` - shorthand for `validation.Comparables[[]string]()`;
* `validator.ValidateCountable()` - shorthand for `validation.Countable()`;
* `validator.ValidateTime()` - shorthand for `validation.Time()`;
* `validator.ValidateEachString()` - shorthand for `validation.EachString()`;
* `validator.ValidateIt()` - shorthand for `validation.Valid()`.

See usage examples in the [documentation](https://pkg.go.dev/github.com/muonsoft/validation#Validator.Validate).

## How to use the validator

There are two ways to use the validator service. You can build your instance of validator service by
using `validation.NewValidator()` or use singleton service from package `github.com/muonsoft/validation/validator`.

Example of creating a new instance of the validator service:

```go
// import "github.com/muonsoft/validation"

validator, err := validation.NewValidator(
    validation.DefaultLanguage(language.Russian), // passing default language of translations
    validation.Translations(russian.Messages),    // setting up custom or built-in translations
    validation.SetViolationFactory(userViolationFactory), // if you want to override creation of violations
)

// don't forget to check for errors
if err != nil {
    fmt.Println(err)
}
```

If you want to use a singleton service make sure to set up your configuration once during the initialization of your
application.

```go
// import "github.com/muonsoft/validation/validator"

err := validator.SetUp(
    validation.DefaultLanguage(language.Russian), // passing default language of translations
    validation.Translations(russian.Messages),    // setting up custom or built-in translations
    validation.SetViolationFactory(userViolationFactory), // if you want to override creation of violations
)

// don't forget to check for errors
if err != nil {
    fmt.Println(err)
}
```
