# Custom constraints and adapters

Choose the smallest abstraction that owns the rule:

- One safe boolean expression: `Check` or `CheckProperty`.
- A pure string predicate: `OfStringBy`, with explicit error and message.
- A reusable rule with options or dependencies: the appropriate typed constraint
  interface, or `Constraint[T]` used with `This`/`Each`.
- An existing method with extra inputs: `ValidatableFunc`, reusing the validator
  supplied to the callback.

```go
package main

import (
    "context"
    "fmt"

    "github.com/muonsoft/validation"
    "github.com/muonsoft/validation/it"
)

type Limits struct { Max int }
type Detail struct { Quantity int }
func (d Detail) Validate(ctx context.Context, v *validation.Validator, limits Limits) error {
    return v.Validate(ctx, validation.NumberProperty("quantity", d.Quantity, it.IsLessThanOrEqual(limits.Max)))
}
func main() {
    v, err := validation.NewValidator()
    if err != nil { panic(err) }
    detail := Detail{Quantity: 5}
    err = v.Validate(context.Background(), validation.ValidProperty("detail",
        validation.ValidatableFunc(func(ctx context.Context, child *validation.Validator) error {
            return detail.Validate(ctx, child, Limits{Max: 3})
        }),
    ))
    violations, ok := validation.UnwrapViolations(err)
    if !ok { panic(err) }
    fmt.Println(violations.First().Violation().PropertyPath().String())
}
// Output:
// detail.quantity
```

The richer method above is deliberately **not** `validation.Validatable`: its
signature has an extra parameter. The adapter implements that interface.

For a string constraint, implement
`ValidateString(context.Context, *validation.Validator, *string) error`.
Decide its nil/empty contract explicitly, commonly skipping them and composing
requiredness separately. Report violations through `CreateViolation` or
`BuildViolation`; return operational/configuration failures as ordinary errors.

Keep error sentinels at package scope. `WithError` changes the code; it does not
also change the message. Set both when a custom rule needs a custom message.
Use template parameters instead of interpolating a new translation key per value.

Avoid copying the library's entire constraint implementation pattern into every
application rule. Add fluent options, group support, or multiple interfaces only
when callers actually need them. Preserve existing project helpers and wrappers
whose semantics already match.
