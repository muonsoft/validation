# Entities and DTOs

Use a DTO/command/query's `Validate` for its input contract: pagination, selected
IDs, operation flags, and field presence. Put reusable state invariants on the
entity or value object that owns them. A DTO that maps directly to an entity can
delegate or validate the assembled entity instead of duplicating every field rule.

When extending existing validation, retain its independent rules as well as the
new ones. A replacement by a built-in must preserve nil/empty behavior, error
identity, groups, and paths; a maximum length check does not replace requiredness.
Delegate child rules to their existing owner and keep injected project constraints
in the call chain rather than copying their current implementation.

For creation, assemble and explicitly normalize the object before validating the
state to be saved. For updates, load the current state when required, apply the
operation's actual replacement or patch semantics, normalize, then validate.
Input checks that protect loading can run earlier. Loading and access checks need
not be deferred until every entity invariant is checked.

A partial DTO must distinguish omission from a supplied zero value. A pointer can
express two states, but does not by itself distinguish JSON omission from explicit
null. Use the application's existing presence representation when all three matter.

```go
package main

import (
    "context"
    "fmt"
    "strings"

    "github.com/muonsoft/validation"
    "github.com/muonsoft/validation/it"
)

type Record struct { Name string }

func (r Record) Validate(ctx context.Context, v *validation.Validator) error {
    return v.Validate(ctx, validation.StringProperty("name", r.Name, it.IsNotBlank()))
}

func prepare(ctx context.Context, v *validation.Validator, name string) (Record, error) {
    r := Record{Name: strings.TrimSpace(name)}
    if err := v.ValidateIt(ctx, r); err != nil { return Record{}, err }
    return r, nil
}

func main() {
    v, err := validation.NewValidator()
    if err != nil { panic(err) }
    r, err := prepare(context.Background(), v, "  Example  ")
    if err != nil { panic(err) }
    fmt.Println(r.Name)
}
// Output:
// Example
```

Use the application's configured validator at its existing boundary. Dependency
injection is convenient but no particular container, service layout, or HTTP
framework is required. `validator.Default()` is also available; configure a
default once at initialization, not during each request.

Do not require every domain operation to become a validation list. A failure to
load the primary object, an authorization decision, and an input violation have
different meanings. Preserve the application's error contract.
