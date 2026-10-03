# Errors and translations

`validation.NewError(code, messageTemplate)` creates a package-level sentinel.
`Error()` returns its machine code; `Message()` returns its message template.
`CreateViolation` uses the configured factory and translator and implements
`error`; it does not choose an HTTP status. Translate violations at the transport
boundary using the application's established response contract.

```go
package main

import (
    "context"
    "errors"
    "fmt"

    "github.com/muonsoft/validation"
    "golang.org/x/text/language"
    "golang.org/x/text/message/catalog"
)

var ErrLimit = validation.NewError("limit_exceeded", "Use at most {{ limit }} entries.")
func main() {
    v, err := validation.NewValidator(validation.Translations(map[language.Tag]map[string]catalog.Message{
        language.Russian: {ErrLimit.Message(): catalog.String("Используйте не более {{ limit }} элементов.")},
    }))
    if err != nil { panic(err) }
    violation := v.WithLanguage(language.Russian).BuildViolation(context.Background(), ErrLimit, ErrLimit.Message()).
        AtProperty("entries").WithParameter("{{ limit }}", "3").Create()
    fmt.Println(errors.Is(violation, ErrLimit))
    fmt.Println(violation.Message())
}
// Output:
// true
// Используйте не более 3 элементов.
```

Keep machine codes stable and user messages actionable. Do not include internal
requirement identifiers. Language choices and code naming conventions belong to
the application; English fallback with registered translations is one useful
pattern, not a requirement to add a particular locale to every project.

Use `UnwrapViolations` or `errors.As` to distinguish violations from technical
errors; use `errors.Is` for sentinel identity. Do not branch on rendered text.
Prefer code and path assertions in tests; assert text when testing translation
itself. Custom message overrides need corresponding translation keys.

When accumulating manually, propagate an error from `AppendFromError` and use
`AsError()` on the completed list. Do not return a non-nil empty list as successful
validation. Never reconstruct an ordinary dependency error as a user violation
just to fit an API response shape.
