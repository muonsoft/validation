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
    "fmt"

    "github.com/muonsoft/validation"
    "github.com/muonsoft/validation/it"
)

type Status string
func (s Status) Validate(ctx context.Context, v *validation.Validator) error {
    return v.Validate(ctx, validation.Comparable(s, it.IsOneOf(Status("draft"), Status("ready")).WithoutBlank()))
}
type Input struct { Note *string; Status Status; Title string }
func (x Input) Validate(ctx context.Context, v *validation.Validator) error {
    return v.Validate(ctx,
        validation.NilStringProperty("note", x.Note, it.IsNotBlank().WithAllowedNil()),
        validation.ValidProperty("status", x.Status),
        validation.StringProperty("title", x.Title, it.IsNotBlank().WhenGroups("publish")),
    )
}
func main() {
    v, err := validation.NewValidator()
    if err != nil { panic(err) }
    x := Input{Status: "draft"}
    fmt.Println(v.ValidateIt(context.Background(), x) == nil)
    fmt.Println(v.WithGroups(validation.DefaultGroup, "publish").ValidateIt(context.Background(), x) != nil)
}
// Output:
// true
// true
```

`IsOneOf` normally accepts a zero value even when absent from the choices. Use
`WithoutBlank` when enum membership must include rejecting that zero value.

Use `CheckProperty` for a simple predicate attributed to one field and `Check`
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
