# Golang validation framework

[![Go Reference](https://pkg.go.dev/badge/github.com/muonsoft/validation.svg)](https://pkg.go.dev/github.com/muonsoft/validation)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/muonsoft/validation)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/muonsoft/validation)
![GitHub](https://img.shields.io/github/license/muonsoft/validation)
[![tests](https://github.com/muonsoft/validation/actions/workflows/tests.yml/badge.svg)](https://github.com/muonsoft/validation/actions/workflows/tests.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/muonsoft/validation)](https://goreportcard.com/report/github.com/muonsoft/validation)
[![Code Coverage](https://scrutinizer-ci.com/g/muonsoft/validation/badges/coverage.png?b=main)](https://scrutinizer-ci.com/g/muonsoft/validation/?branch=main)
[![Scrutinizer Code Quality](https://scrutinizer-ci.com/g/muonsoft/validation/badges/quality-score.png?b=main)](https://scrutinizer-ci.com/g/muonsoft/validation/?branch=main)
[![Maintainability](https://api.codeclimate.com/v1/badges/1385bcb467b6e43bff8d/maintainability)](https://codeclimate.com/github/muonsoft/validation/maintainability)
[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-2.0-4baaaa.svg)](CODE_OF_CONDUCT.md)

The package provides tools for data validation. It is designed to create complex validation rules with abilities to hook into the validation process.

This project is inspired by [Symfony Validator component](https://symfony.com/index.php/doc/current/validation.html).

## Key features

* Flexible and customizable API built in mind to use benefits of static typing
* Nice and readable way to describe validation process in code
* Validation of different types: booleans, numbers, strings, slices, maps, and time
* Validation of custom data types that implements `Validatable` interface
* Customizable validation errors with translations and pluralization supported out of the box
* Easy way to create own validation rules with context propagation and message translations

## Work-in-progress notice

This package is under active development and API may be changed until the first major version will be released. Minor versions `n` 0.n.m may contain breaking changes. Patch versions `m` 0.n.m may contain only bug fixes.

First stable version aimed to be released after [generics implementation](https://go.googlesource.com/proposal/+/refs/heads/master/design/go2draft-contracts.md) in Golang.

## Installation

Run the following command to install the package

```bash
go get -u github.com/muonsoft/validation
```

## How to use

### Basic concepts

The validation process is built around functional options and passing values by specific typed arguments. A common way to use validation is to call the `validator.Validate` method and pass the argument option with the list of validation constraints.

```golang
err := validator.Validate(context.Background(), validation.String("", it.IsNotBlank()))

violations := err.(validation.ViolationList)
for _, violation := range violations {
    fmt.Println(violation.Error())
}
// Output:
// violation: This value should not be blank.
```

List of common [validation arguments](https://pkg.go.dev/github.com/muonsoft/validation#Argument)

* `validation.Value()` - passes any value. It uses reflection to detect the type of the argument and pass it to a specific validation method.
* `validation.Bool()` - passes boolean value.
* `validation.NilBool()` - passes nillable boolean value.
* `validation.Number()` - passes any numeric value. At the moment it uses reflection for executing validation process.
* `validation.String()` - passes string value.
* `validation.NilString()` - passes nillable string value.
* `validation.Strings()` - passes slice of strings value.
* `validation.Iterable()` - passes array, slice or a map. At the moment it uses reflection for executing validation process.
* `validation.Countable()` - you can pass result of `len()` to use easy way of iterable validation based only on count of the elements.
* `validation.Time()` - passes `time.Time` value.
* `validation.NilTime()` - passes nillable `time.Time` value.
* `validation.Each()` - passes array, slice or a map. Used to validate each value of iterable. It uses reflection.
* `validation.EachString()` - passes slice of strings. This is more performant version than `Each`.
* `validation.Valid()` - passes `Validatable` value to run embedded validation.

For single value validation, you can use shorthand versions of the validation method.

* `validator.ValidateValue()`
* `validator.ValidateBool()`
* `validator.ValidateNumber()`
* `validator.ValidateString()`
* `validator.ValidateStrings()`
* `validator.ValidateIterable()`
* `validator.ValidateCountable()`
* `validator.ValidateTime()`
* `validator.ValidateEach()`
* `validator.ValidateEachString()`
* `validator.ValidateValidatable()`

See usage examples in the [documentation](https://pkg.go.dev/github.com/muonsoft/validation#Validator.Validate).

### How to use the validator

There are two ways to use the validator service. You can build your instance of validator service by using `validation.NewValidator()` or use singleton service from package `github.com/muonsoft/validation/validator`.

Example of creating a new instance of the validator service.

```golang
// import "github.com/muonsoft/validation"

validator, err := validation.NewValidator(
	validation.DefaultLanguage(language.Russian), // passing default language of translations
	validation.Translations(russian.Messages), // setting up custom or built-in translations
	validation.SetViolationFactory(userViolationFactory), // if you want to override creation of violations
)

// don't forget to check for errors
if err != nil {
	fmt.Println(err)
}
```

If you want to use a singleton service make sure to set up your configuration once during the initialization of your application.

```golang
// import "github.com/muonsoft/validation/validator"

err := validator.SetOptions(
    validation.DefaultLanguage(language.Russian), // passing default language of translations
    validation.Translations(russian.Messages), // setting up custom or built-in translations
    validation.SetViolationFactory(userViolationFactory), // if you want to override creation of violations
)

// don't forget to check for errors
if err != nil {
    fmt.Println(err)
}
```

### Processing property paths

One of the main concepts of the package is to provide helpful violation descriptions for complex data structures. For example, if you have lots of structures used in other structures you want somehow to describe property paths to violated attributes. 

The [property path](https://pkg.go.dev/github.com/muonsoft/validation#PropertyPath) generated by the validator indicates how it reached the invalid value from the root element. Property path is denoted by dots, while array access is denoted by square brackets. For example, `book.keywords[0]` means that the violation occurred on the first element of array `keywords` in the `book` object. 

You can pass a property name or an array index via `validation.PropertyName()` and `validation.ArrayIndex()` options.

```golang
err := validator.Validate(
    context.Background(),
    validation.String(
        "",
        validation.PropertyName("properties"),
        validation.ArrayIndex(1),
        validation.PropertyName("tag"),
        it.IsNotBlank(),
    ),
)

violation := err.(validation.ViolationList)[0]
fmt.Println("property path:", violation.GetPropertyPath().Format())
// Output:
// property path: properties[1].tag
```

Also, you can create scoped validator by using `valdiator.AtProperty()` or `validator.AtIndex()` methods. It can be used to validate a couple of attributes of one object.

```golang
err := validator.
    AtProperty("properties").
    AtIndex(1).
    AtProperty("tag").
    Validate(context.Background(), validation.String("", it.IsNotBlank()))

violation := err.(validation.ViolationList)[0]
fmt.Println("property path:", violation.GetPropertyPath().Format())
// Output:
// property path: properties[1].tag
```

For a better experience with struct validation, you can use shorthand versions of validation arguments with passing property names.

* `validation.PropertyValue()`
* `validation.BoolProperty()`
* `validation.NilBoolProperty()`
* `validation.NumberProperty()`
* `validation.StringProperty()`
* `validation.NilStringProperty()`
* `validation.StringsProperty()`
* `validation.IterableProperty()`
* `validation.CountableProperty()`
* `validation.TimeProperty()`
* `validation.NilTimeProperty()`
* `validation.EachProperty()`
* `validation.EachStringProperty()`
* `validation.ValidProperty()`

```golang
err := validator.Validate(
    context.Background(),
    validation.StringProperty("property", "", it.IsNotBlank()),
)

violation := err.(validation.ViolationList)[0]
fmt.Println("property path:", violation.GetPropertyPath().Format())
// Output:
// property path: property
```

### Validation of structs

There are few ways to validate structs. The simplest one is to call the `validator.Validate` method with property arguments.

```golang
document := Document{
    Title:    "",
    Keywords: []string{"", "book", "fantasy", "book"},
}

err := validator.Validate(
    context.Background(),
    validation.StringProperty("title", document.Title, it.IsNotBlank()),
    validation.CountableProperty("keywords", len(document.Keywords), it.HasCountBetween(5, 10)),
    validation.StringsProperty("keywords", document.Keywords, it.HasUniqueValues()),
    validation.EachStringProperty("keywords", document.Keywords, it.IsNotBlank()),
)

if violations, ok := validation.UnwrapViolationList(err); ok {
    for violation := violations.First(); violation != nil; violation = violation.Next() {
        fmt.Println(violation)
    }
}
// Output:
// violation at 'title': This value should not be blank.
// violation at 'keywords': This collection should contain 5 elements or more.
// violation at 'keywords': This collection should contain only unique elements.
// violation at 'keywords[0]': This value should not be blank.
```

The recommended way is to implement the `validation.Validatable` interface for your structures. By using it you can build complex validation rules on a set of objects used in other objects.

```golang
type Product struct {
    Name       string
    Tags       []string
    Components []Component
}

func (p Product) Validate(ctx context.Context, validator *validation.Validator) error {
    return validator.Validate(
        ctx,
        validation.StringProperty("name", p.Name, it.IsNotBlank()),
        validation.CountableProperty("tags", len(p.Tags), it.HasMinCount(5)),
        validation.StringsProperty("tags", p.Tags, it.HasUniqueValues()),
        validation.EachStringProperty("tags", p.Tags, it.IsNotBlank()),
        // this also runs validation on each of the components
        validation.IterableProperty("components", p.Components, it.HasMinCount(1)),
    )
}

type Component struct {
    ID   int
    Name string
    Tags []string
}

func (c Component) Validate(ctx context.Context, validator *validation.Validator) error {
    return validator.Validate(
        ctx,
        validation.StringProperty("name", c.Name, it.IsNotBlank()),
        validation.CountableProperty("tags", len(c.Tags), it.HasMinCount(1)),
    )
}

func main() {
    p := Product{
        Name: "",
        Tags: []string{"device", "", "phone", "device"},
        Components: []Component{
            {
                ID:   1,
                Name: "",
            },
        },
    }
    
    err := validator.ValidateValidatable(context.Background(), p)
    
    violations := err.(validation.ViolationList)
    for _, violation := range violations {
        fmt.Println(violation.Error())
    }
    // Output:
    // violation at 'name': This value should not be blank.
    // violation at 'tags': This collection should contain 5 elements or more.
    // violation at 'tags': This collection should contain only unique elements.
    // violation at 'tags[1]': This value should not be blank.
    // violation at 'components[0].name': This value should not be blank.
    // violation at 'components[0].tags': This collection should contain 1 element or more.
}
```

### Conditional validation

You can use the `When()` method on any of the built-in constraints to execute conditional validation on it. 

```golang
err := validator.Validate(
    context.Background(),
    validation.StringProperty("text", note.Text, it.IsNotBlank().When(note.IsPublic)),
)

violations := err.(validation.ViolationList)
for _, violation := range violations {
    fmt.Println(violation.Error())
}
// Output:
// violation at 'text': This value should not be blank.
```

### Working with violations and errors

There are two types of errors returned from the validator. One is validation violations and another is internal errors (for example, when attempting to apply a constraint on not applicable argument type). The best way to handle validation errors is to check for implementing the `validation.ViolationList` struct. You can use the default way to unwrap errors.

```golang
err := validator.Validate(/* validation arguments */)

var violations ViolationList
if err != nil {
    if errors.As(err, &violations) {
        // handle violations
    } else {
        // handle internal error
    }
}
```

Also, you can use helper function `validation.UnwrapViolationList()`.

```golang
err := validator.Validate(/* validation arguments */)
if violations, ok := validation.UnwrapViolationList(err); ok {
    // handle violations
} else if err != nil {
    // handle internal error
}
```

The validation error called violation consists of a few parameters.

* `code` - unique, short, and semantic violation code that can be used to programmatically test for specific violation. All `code` values are defined in the `github.com/muonsoft/validation/code` package and are protected by backward compatibility rules.
* `message` - translated message with injected values from constraint. It can be used to show a description of a violation to the end-user. Possible values for build-in constraints are defined in the `github.com/muonsoft/validation/message` package and can be changed at any time, even in patch versions.
* `messageTemplate` - template for rendering message. Alongside `parameters` it can be used to render the message on the client-side of the library.
* `parameters` is the map of the template variables and their values provided by the specific constraint.
* `propertyPath` points to the violated property as it described in the [previous section](#processing-property-paths).

You can hook into process of violation generation by implementing `validation.ViolationFactory` interface and passing it via `validation.SetViolationFactory()` option. Custom violation must implement `validation.Violation` interface.

### How to use translations

By default, all violation messages are generated in the English language with pluralization capabilities. To use a custom language you have to load translations on validator initialization. Built-in translations are available in the sub-packages of the package `github.com/muonsoft/message/translations`. The translation mechanism is provided by the `golang.org/x/text` package (be aware, it has no stable version yet).

```golang
// import "github.com/muonsoft/validation/message/translations/russian"

validator, err := validation.NewValidator(
	validation.Translations(russian.Messages),
)
```

There are different ways to initialize translation to a specific language.

The first one is to use the default language. In that case, all messages will be translated to this language.

```golang
validator, _ := validation.NewValidator(
    validation.Translations(russian.Messages),
    validation.DefaultLanguage(language.Russian),
)

err := validator.ValidateString(context.Background(), "", it.IsNotBlank())

violations := err.(validation.ViolationList)
for _, violation := range violations {
    fmt.Println(violation.Error())
}
// Output:
// violation: Значение не должно быть пустым.
```

The second way is to use the `validation.Language()` argument. Be aware that this method works only on a specific scope. Also, you can use the `validator.WithLanguage()` method to create scoped validator and use it in different places.

```golang
validator, _ := validation.NewValidator(
    validation.Translations(russian.Messages),
)

err := validator.Validate(
    context.Background(),
    validation.Language(language.Russian),
    validation.String("", it.IsNotBlank()),
)

violations := err.(validation.ViolationList)
for _, violation := range violations {
    fmt.Println(violation.Error())
}
// Output:
// violation: Значение не должно быть пустым.
```

The last way is to pass language via context. It is provided by the `github.com/muonsoft/language` package and can be useful in combination with [language middleware](https://github.com/muonsoft/language/blob/main/middleware.go).

```golang
// import "github.com/muonsoft/language"

validator, _ := validation.NewValidator(
    validation.Translations(russian.Messages),
)
ctx := language.WithContext(context.Background(), language.Russian)

err := validator.ValidateString(ctx, "", it.IsNotBlank())

violations := err.(validation.ViolationList)
for _, violation := range violations {
    fmt.Println(violation.Error())
}
// Output:
// violation: Значение не должно быть пустым.
```

You can see the complex example with handling HTTP request [here](https://pkg.go.dev/github.com/muonsoft/validation#example-Validator.Validate-HttpHandler).

### Customizing violation messages

You may customize the violation message on any of the built-in constraints by calling the `Message()` method or similar if the constraint has more than one template. Also, you can include template parameters in it. See details of a specific constraint to know what parameters are available.

```golang
err := validator.ValidateString(context.Background(), "", it.IsNotBlank().Message("this value is required"))

violations := err.(validation.ViolationList)
for _, violation := range violations {
    fmt.Println(violation.Error())
}
// Output:
// violation: this value is required
```

To use pluralization and message translation you have to load up your translations via `validation.Translations()` option to the validator. See `golang.org/x/text` package [documentation](https://pkg.go.dev/golang.org/x/text) for details of translations.

```golang
const customMessage = "tags should contain more than {{ limit }} element(s)"
validator, _ := validation.NewValidator(
    validation.Translations(map[language.Tag]map[string]catalog.Message{
        language.Russian: {
            customMessage: plural.Selectf(1, "",
                plural.One, "теги должны содержать {{ limit }} элемент и более",
                plural.Few, "теги должны содержать более {{ limit }} элемента",
                plural.Other, "теги должны содержать более {{ limit }} элементов"),
        },
    }),
)

var tags []string
err := validator.ValidateIterable(
    context.Background(),
    tags,
    validation.Language(language.Russian),
    it.HasMinCount(1).MinMessage(customMessage),
)

violations := err.(validation.ViolationList)
for _, violation := range violations {
    fmt.Println(violation.Error())
}
// Output:
// violation: теги должны содержать 1 элемент и более
```

### Creating custom constraints

Everything you need to create a custom constraint is to implement one of the interfaces:

* `BoolConstraint` - for validating boolean values;
* `NumberConstraint` - for validating numeric values;
* `StringConstraint` - for validating string values;
* `StringsConstraint` - for validating slice of strings;
* `IterableConstraint` - for validating iterable values: arrays, slices, or maps;
* `CountableConstraint` - for validating iterable values based only on the count of elements;
* `TimeConstraint` - for validating date\time values.

Also, you can combine several types of constraints. See examples for more details:

* [custom static constraint](https://pkg.go.dev/github.com/muonsoft/validation#example-Validator.Validate-CustomConstraint);
* [custom constraint as a service](https://pkg.go.dev/github.com/muonsoft/validation#example-Validator.ValidateBy-CustomServiceConstraint).
* [custom constraint with custom argument for domain type](https://pkg.go.dev/github.com/muonsoft/validation#example-NewArgument-CustomArgumentConstraintValidator).

## Contributing

You may help this project by

* reporting an [issue](https://github.com/muonsoft/validation/issues);
* making translations for error messages;
* suggest an improvement or [discuss](https://github.com/muonsoft/validation/discussions) the usability of the package.

If you'd like to contribute, see [the contribution guide](CONTRIBUTING.md). Pull requests are welcome.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
