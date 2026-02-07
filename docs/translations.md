# Translations and Custom Messages

## How to use translations

By default, all violation messages are generated in the English language with pluralization capabilities. To use a
custom language you have to load translations on validator initialization. Built-in translations are available in the
sub-packages of the package `github.com/muonsoft/validation/message/translations`. The translation mechanism is provided by
the `golang.org/x/text` package (be aware, it has no stable version yet).

```go
// import "github.com/muonsoft/validation/message/translations/russian"

validator, err := validation.NewValidator(
    validation.Translations(russian.Messages),
)
```

There are different ways to initialize translation to a specific language.

The first one is to use the default language. In that case, all messages will be translated to this language.

```go
validator, _ := validation.NewValidator(
    validation.Translations(russian.Messages),
    validation.DefaultLanguage(language.Russian),
)

err := validator.ValidateString(context.Background(), "", it.IsNotBlank())

if violations, ok := validation.UnwrapViolations(err); ok {
    for _, violation := range violations.All() {
        fmt.Println(violation.Error())
    }
}
// Output:
// violation: Значение не должно быть пустым.
```

The second way is to use the `validator.WithLanguage()` method to create context validator and use it in different places.

```go
validator, _ := validation.NewValidator(
    validation.Translations(russian.Messages),
)

err := validator.WithLanguage(language.Russian).Validate(
    context.Background(),
    validation.String("", it.IsNotBlank()),
)

if violations, ok := validation.UnwrapViolations(err); ok {
    for _, violation := range violations.All() {
        fmt.Println(violation.Error())
    }
}
// Output:
// violation: Значение не должно быть пустым.
```

The last way is to pass language via context. It is provided by the `github.com/muonsoft/language` package and can be
useful in combination with [language middleware](https://github.com/muonsoft/language/blob/main/middleware.go).

```go
// import "github.com/muonsoft/language"

validator, _ := validation.NewValidator(
    validation.Translations(russian.Messages),
)

ctx := language.WithContext(context.Background(), language.Russian)
err := validator.ValidateString(ctx, "", it.IsNotBlank())

if violations, ok := validation.UnwrapViolations(err); ok {
    for _, violation := range violations.All() {
        fmt.Println(violation.Error())
    }
}
// Output:
// violation: Значение не должно быть пустым.
```

You can see the complex example with handling HTTP
request [here](https://pkg.go.dev/github.com/muonsoft/validation#example-Validator.Validate-HttpHandler).

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
    // your implementation of translation mechanism
}

translator := &CustomTranslator{}

validator, err := validation.NewValidator(validation.SetTranslator(translator))
if err != nil {
    log.Fatal(err)
}
```

## Customizing violation messages

You may customize the violation message on any of the built-in constraints by calling the `Message()` method or similar
if the constraint has more than one template. Also, you can include template parameters in it. See details of a specific
constraint to know what parameters are available.

```go
err := validator.ValidateString(context.Background(), "", it.IsNotBlank().Message("this value is required"))

if violations, ok := validation.UnwrapViolations(err); ok {
    for _, violation := range violations.All() {
        fmt.Println(violation.Error())
    }
}
// Output:
// violation: this value is required
```

To use pluralization and message translation you have to load up your translations via `validation.Translations()`
option to the validator. See `golang.org/x/text` package [documentation](https://pkg.go.dev/golang.org/x/text) for
details of translations.

```go
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

if violations, ok := validation.UnwrapViolations(err); ok {
    for _, violation := range violations.All() {
        fmt.Println(violation.Error())
    }
}
// Output:
// violation: теги должны содержать 1 элемент и более
```
