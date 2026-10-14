# Optional and conditional values

For primitive pointers use `NilStringProperty`, `NilNumberProperty`, and related
arguments. Most value constraints skip nil. Requiredness and allowed empty values
are separate decisions: `IsNotNil` permits an empty string, while
`IsNotBlank().WithAllowedNil()` permits omission but rejects a supplied empty string.
Whitespace passes `IsNotBlank`; normalize explicitly if whitespace is blank in
the application's contract. Check individual constraints rather than assuming all
format or length checks share empty-value behavior.

```go
package main

import (
    "context"
    "errors"
    "fmt"

    "github.com/muonsoft/validation"
    "github.com/muonsoft/validation/it"
)

type Status string
func (s Status) Validate(ctx context.Context, v *validation.Validator) error {
    return v.Validate(ctx, validation.Comparable(s, it.IsOneOf(Status("draft"), Status("ready")).WithoutBlank()))
}
type Input struct { Note *string; Status Status; Title string; Mode string }
func (x Input) Validate(ctx context.Context, v *validation.Validator) error {
    mode := it.IsOneOf("quiet", "active")
    if !v.IsIgnoredForGroups("publish") { mode = mode.WithoutBlank() }
    return v.Validate(ctx,
        validation.NilStringProperty("note", x.Note, it.IsNotBlank().WithAllowedNil()),
        validation.ValidProperty("status", x.Status),
        validation.StringProperty("mode", x.Mode, mode),
        validation.StringProperty("title", x.Title, it.IsNotBlank().WhenGroups("publish")),
    )
}
func main() {
    v, err := validation.NewValidator()
    if err != nil { panic(err) }
    x := Input{Status: "draft"}
    fmt.Println(v.ValidateIt(context.Background(), x) == nil)
    x.Title = "Example"
    err = v.WithGroups(validation.DefaultGroup, "publish").ValidateIt(context.Background(), x)
    violations, ok := validation.UnwrapViolations(err)
    if !ok || violations.Len() != 1 { panic("expected one mode violation") }
    violation := violations.First().Violation()
    fmt.Println(violation.PropertyPath().String(), errors.Is(violation, validation.ErrNoSuchChoice))
}
// Output:
// true
// mode true
```

`IsOneOf` normally accepts a zero value even when absent from the choices. Use
`WithoutBlank` when enum membership must include rejecting that zero value.
The violation remains `ErrNoSuchChoice`; adding `IsNotBlank` instead produces
`ErrIsBlank`, which is not interchangeable when callers depend on error identity.
In the example, empty mode is allowed by default but rejected for publishing.
Nonempty unknown modes must still be rejected in either group. Select the blank
policy on one choice constraint rather than stacking overlapping membership
checks that could produce duplicate violations.

Use `CheckProperty` for a simple predicate with no matching built-in, attributed to one field, and `Check`
for an object-level invariant, with a meaningful code/message. Use `When` when
the input expression is safe to evaluate. For a range with optional endpoints,
construct comparisons only after checking both pointers; do not dereference in
an argument and expect `.When` to guard it.

Select groups at the operation boundary. Include `DefaultGroup` alongside a
scenario group when the operation should retain default rules. Test both selected
and unselected groups and ensure derived validators reach all children.

A nullable value object's present state can delegate to the base value's
validation. Keep one owner for the underlying rules. Distinguish an input patch
from the resulting complete entity when deciding which rules apply.

## Patch fields with format and business choices

A business allowlist is still a choice constraint: configure `IsOneOf` with the
policy's error and message. Use a custom predicate only when membership alone
does not express the policy. `Sequentially` makes a valid format a prerequisite
without evaluating the same format again through `is` or `validate` in `When`.
Keep that sequence local to the field so independent violations still accumulate.

```go
package main

import (
    "context"
    "errors"
    "fmt"

    "github.com/muonsoft/validation"
    "github.com/muonsoft/validation/it"
)

var ErrCountryPolicy = validation.NewError("country_policy", "Country is not supported.")

type TextPatch struct { Present bool; Value *string }
type Input struct { Country TextPatch; Name string }

func (x Input) Validate(ctx context.Context, v *validation.Validator) error {
    args := []validation.Argument{
        validation.StringProperty("name", x.Name, it.IsNotBlank()),
    }
    // Omission ignores even a stale pointer; present nil means clearing.
    if x.Country.Present && x.Country.Value != nil {
        args = append(args, validation.Sequentially(
            validation.NilStringProperty("country", x.Country.Value, it.IsNotBlank(), it.IsCountry()),
            validation.NilStringProperty("country", x.Country.Value,
                it.IsOneOf("CA", "NZ").WithError(ErrCountryPolicy).WithMessage(ErrCountryPolicy.Message())),
        ))
    }
    return v.Validate(ctx, args...)
}

func main() {
    v, err := validation.NewValidator()
    if err != nil { panic(err) }
    ctx := context.Background()
    country := "US"
    x := Input{Country: TextPatch{Present: true, Value: &country}}
    err = x.Validate(ctx, v.AtProperty("payload"))
    violations, ok := validation.UnwrapViolations(err)
    if !ok { panic(err) }
    fmt.Println(violations.Len(), errors.Is(err, ErrCountryPolicy))
    country = "??"
    err = x.Validate(ctx, v)
    violations, ok = validation.UnwrapViolations(err)
    if !ok { panic(err) }
    fmt.Println(violations.Len(), errors.Is(err, validation.ErrInvalidCountry), errors.Is(err, ErrCountryPolicy))
    x.Name = "Ada"
    x.Country.Present = false
    fmt.Println(x.Validate(ctx, v) == nil)
    x.Country = TextPatch{Present: true}
    fmt.Println(x.Validate(ctx, v) == nil)
}
// Output:
// 2 true
// 2 true false
// true
// true
```
