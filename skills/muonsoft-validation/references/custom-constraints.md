# Custom constraints and adapters

First check [built-in constraints](builtin-constraints.md) and existing project
rules. Create a custom rule when those do not express the required semantics.
Choose the smallest abstraction that owns the rule:

- An existing typed constraint: pass it directly to the matching argument.
- One safe boolean expression with no matching built-in: `Check` or `CheckProperty`.
- A pure string predicate with no matching built-in: `OfStringBy`, with explicit error and message.
- A reusable rule with options or dependencies: the appropriate typed constraint
  interface, or `Constraint[T]` used with `This`/`Each`.
- An existing method with extra inputs that is not already a constraint: `ValidatableFunc`, reusing the validator
  supplied to the callback.

An application-owned policy implementing `StringConstraint` already fits
`StringProperty`. Use `Sequentially` to gate it on requiredness/format; an extra
`ValidatableFunc` or a direct call to its internal callback adds no semantics.

```go
package main

import (
    "context"
    "fmt"

    "github.com/muonsoft/validation"
    "github.com/muonsoft/validation/it"
)

type EmailPolicy struct { Check func(context.Context, *validation.Validator, string) error }
func (p EmailPolicy) ValidateString(ctx context.Context, v *validation.Validator, value *string) error {
    if value == nil || *value == "" { return nil }
    return p.Check(ctx, v, *value)
}

func checkInput(ctx context.Context, v *validation.Validator, name, email string, policy validation.StringConstraint) error {
    return v.Validate(ctx,
        validation.StringProperty("name", name, it.IsNotBlank()),
        validation.Sequentially(
            validation.StringProperty("email", email, it.IsNotBlank(), it.IsEmail()),
            validation.StringProperty("email", email, policy),
        ),
    )
}

func main() {
    v, err := validation.NewValidator()
    if err != nil { panic(err) }
    calls := 0
    policy := EmailPolicy{Check: func(ctx context.Context, child *validation.Validator, email string) error {
        calls++
        return nil
    }}
    // Invalid name does not suppress the policy for a valid email.
    _ = checkInput(context.Background(), v, "", "ada@example.com", policy)
    fmt.Println(calls)
    _ = checkInput(context.Background(), v, "Ada", "invalid", policy)
    fmt.Println(calls)
}
// Output:
// 1
// 1
```

Use an adapter when the existing method actually has a different signature:

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
