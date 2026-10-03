# Nested objects and collections

Use `ValidProperty` for one child, `ValidSliceProperty` for validatable slice
elements, and `ValidMapProperty` for string-keyed maps. Children use relative
names. Use `AtProperty` to combine count, uniqueness, and element checks at one
prefix. Introduce a named collection implementing `Validatable` when those rules
are reusable, rather than requiring a wrapper for every slice.

```go
package main

import (
    "context"
    "fmt"

    "github.com/muonsoft/validation"
    "github.com/muonsoft/validation/it"
)

type Entry struct { Code string; Quantity int }
func (e Entry) Validate(ctx context.Context, v *validation.Validator) error {
    return v.Validate(ctx,
        validation.StringProperty("code", e.Code, it.IsNotBlank()),
        validation.NumberProperty("quantity", e.Quantity, it.IsPositive[int]()),
    )
}
type Input struct { Entries []Entry }
func (x Input) Validate(ctx context.Context, v *validation.Validator) error {
    return v.Validate(ctx, validation.AtProperty("entries",
        validation.Countable(len(x.Entries), it.HasMinCount(1)),
        validation.Slice(x.Entries, it.HasUniqueValuesBy(func(e Entry) string { return e.Code })),
        validation.ValidSlice(x.Entries),
    ))
}
func main() {
    v, err := validation.NewValidator()
    if err != nil { panic(err) }
    err = v.ValidateIt(context.Background(), Input{Entries: []Entry{{Code: "A", Quantity: 0}}})
    violations, ok := validation.UnwrapViolations(err)
    if !ok { panic(err) }
    fmt.Println(violations.First().Violation().PropertyPath().String())
}
// Output:
// entries[0].quantity
```

For comparable slices use `Comparables` with `HasUniqueValues`; for object keys
use `Slice` with `HasUniqueValuesBy`. Decide whether uniqueness is an error on the
collection or on individual fields before choosing path options or a custom rule.
Do not substitute artificial zero keys for nil elements: they can create false
duplicates. Specify and test missing-element semantics.

`Valid` does not automatically skip nil. A pointer receiver can explicitly accept
nil; alternatively skip the argument with `When`. Requiredness still needs its
own check. Slice elements need the same deliberate treatment.

For heterogeneous collections or methods with extra arguments, a manual loop can
call each child with `v.AtIndex(i)` and collect via `ViolationList.AppendFromError`.
Return a non-violation error immediately and finish with `list.AsError()` so an
empty list becomes nil. See [custom constraints](custom-constraints.md).

Never encode a path as `PropertyName("entries[0].code")`. It is one literal key,
not an index and property. Use `PropertyName("entries"), ArrayIndex(0),
PropertyName("code")`, or the composition helpers. Test JSON Pointer as well as
display paths when the external API consumes it. Map iteration order is not stable.
