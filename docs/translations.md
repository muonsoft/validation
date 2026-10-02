# Translations and Custom Messages

## How to use translations

Snippets use `context`, `fmt`, `log`, `validation`, `it`, and
`golang.org/x/text/language`, plus `russian` below. The context-language example
uses `github.com/muonsoft/language` instead. The pluralization example is a
complete program with all imports.

By default, all violation messages are generated in the English language with pluralization capabilities. To use a
custom language you have to load translations on validator initialization. Built-in translations are available in the
sub-packages of the package `github.com/muonsoft/validation/message/translations`. The translation mechanism is provided by
the `golang.org/x/text` package.

```go
// Also import "github.com/muonsoft/validation/message/translations/russian"

validator, err := validation.NewValidator(
    validation.Translations(russian.Messages),
)
if err != nil {
    log.Fatal(err)
}
```

There are different ways to initialize translation to a specific language.

The first one is to use the default language. This language is used unless overridden by a contextual validator or request context.

```go
validator, err := validation.NewValidator(
    validation.Translations(russian.Messages),
    validation.DefaultLanguage(language.Russian),
)
if err != nil {
    log.Fatal(err)
}

err = validator.ValidateString(context.Background(), "", it.IsNotBlank())

if violations, ok := validation.UnwrapViolations(err); ok {
    for _, violation := range violations.All() {
        fmt.Println(violation.Error())
    }
}
// Output:
// violation: "Значение не должно быть пустым."
```

The second way is to use the `validator.WithLanguage()` method to create context validator and use it in different places.

```go
validator, err := validation.NewValidator(
    validation.Translations(russian.Messages),
)
if err != nil {
    log.Fatal(err)
}

err = validator.WithLanguage(language.Russian).Validate(
    context.Background(),
    validation.String("", it.IsNotBlank()),
)

if violations, ok := validation.UnwrapViolations(err); ok {
    for _, violation := range violations.All() {
        fmt.Println(violation.Error())
    }
}
// Output:
// violation: "Значение не должно быть пустым."
```

The last way is to pass language via context. It is provided by the `github.com/muonsoft/language` package and can be
useful in combination with [language middleware](https://github.com/muonsoft/language/blob/main/middleware.go).

```go
// import "github.com/muonsoft/language"

validator, err := validation.NewValidator(
    validation.Translations(russian.Messages),
)
if err != nil {
    log.Fatal(err)
}

ctx := language.WithContext(context.Background(), language.Russian)
err = validator.ValidateString(ctx, "", it.IsNotBlank())

if violations, ok := validation.UnwrapViolations(err); ok {
    for _, violation := range violations.All() {
        fmt.Println(violation.Error())
    }
}
// Output:
// violation: "Значение не должно быть пустым."
```

You can see the complex example with handling HTTP
request in [example_http_handler_test.go](../example_http_handler_test.go).

The priority of language selection methods:

* `validator.WithLanguage()` has the highest priority and will override any other options;
* if the validator language is not specified, the validator will try to get the language from the context;
* in all other cases, the default language specified in the translator will be used.

Also, there is an ability to totally override translations behaviour. You can use your own translator by
implementing `validation.Translator` interface and passing it to validator constructor via `SetTranslator` option.

```go
type CustomTranslator struct {
    // your attributes
}

func (t *CustomTranslator) Translate(tag language.Tag, message string, pluralCount int) string {
    return message // Replace with a lookup in your application's translation catalog.
}

translator := &CustomTranslator{}

validator, err := validation.NewValidator(validation.SetTranslator(translator))
if err != nil {
    log.Fatal(err)
}
```

## Customizing violation messages

You may customize the violation message on any of the built-in constraints by calling the `WithMessage()` method or similar
if the constraint has more than one template. Also, you can include template parameters in it. See details of a specific
constraint to know what parameters are available.

```go
err := validator.ValidateString(context.Background(), "", it.IsNotBlank().WithMessage("this value is required"))

if violations, ok := validation.UnwrapViolations(err); ok {
    for _, violation := range violations.All() {
        fmt.Println(violation.Error())
    }
}
// Output:
// violation: "this value is required"
```

To use pluralization and message translation you have to load up your translations via `validation.Translations()`
option to the validator. See `golang.org/x/text` package [documentation](https://pkg.go.dev/golang.org/x/text) for
details of translations.

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/muonsoft/validation"
    "github.com/muonsoft/validation/it"
    "golang.org/x/text/feature/plural"
    "golang.org/x/text/language"
    "golang.org/x/text/message/catalog"
)

func main() {
    const customMessage = "tags should contain at least {{ limit }} element(s)"
    validator, err := validation.NewValidator(
        validation.Translations(map[language.Tag]map[string]catalog.Message{
            language.Russian: {
                customMessage: plural.Selectf(1, "",
                    plural.One, "теги должны содержать {{ limit }} элемент и более",
                    plural.Few, "теги должны содержать {{ limit }} элемента и более",
                    plural.Other, "теги должны содержать {{ limit }} элементов и более"),
            },
        }),
    )
    if err != nil {
        log.Fatal(err)
    }

    var tags []string
    err = validator.WithLanguage(language.Russian).Validate(
        context.Background(),
        validation.Countable(len(tags), it.HasMinCount(1).WithMinMessage(customMessage)),
    )

    if violations, ok := validation.UnwrapViolations(err); ok {
        for _, violation := range violations.All() {
            fmt.Println(violation.Error())
        }
    }
}
// Output:
// violation: "теги должны содержать 1 элемент и более"
```

Executable versions of these examples are in [example_test.go](../example_test.go)
(`ExampleValidator_Validate_translationForCustomMessage` and the translation examples).
