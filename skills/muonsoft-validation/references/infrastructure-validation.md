# Infrastructure validation

Prefer two stages when DB/HTTP checks require structurally valid input:

1. Collect all independent local violations on the object.
2. If that stage succeeds, run independent reference, availability, and uniqueness
   constraints, collecting their violations until a technical error occurs.

Use constraints with injected repository/client interfaces. Keep dependencies out
of data-only DTOs. An application service or a validatable wrapper can compose the
stages. Adapt an existing richer method instead of restructuring the application
solely to enforce this preference.

```go
package main

import (
    "context"
    "fmt"

    "github.com/muonsoft/validation"
    "github.com/muonsoft/validation/it"
)

var ErrUnavailable = validation.NewError("reference_unavailable", "Reference is unavailable.")
type Lookup interface { Exists(context.Context, string) (bool, error) }
type ExistsConstraint struct { Lookup Lookup }
func (c ExistsConstraint) Validate(ctx context.Context, v *validation.Validator, id string) error {
    exists, err := c.Lookup.Exists(ctx, id)
    if err != nil { return fmt.Errorf("lookup reference: %w", err) }
    if !exists { return v.CreateViolation(ctx, ErrUnavailable, ErrUnavailable.Message()) }
    return nil
}
type fakeLookup struct { calls int }
func (f *fakeLookup) Exists(context.Context, string) (bool, error) { f.calls++; return false, nil }
func main() {
    v, err := validation.NewValidator()
    if err != nil { panic(err) }
    lookup := &fakeLookup{}
    err = v.Validate(context.Background(), validation.Sequentially(
        validation.All(validation.StringProperty("reference", "", it.IsNotBlank())),
        validation.This("", ExistsConstraint{Lookup: lookup}).At(validation.PropertyName("reference")),
    ))
    fmt.Println(err != nil, lookup.calls)
}
// Output:
// true 0
```

For a dependency limited to one branch, put only that branch in `Sequentially`.
Other fields can still contribute violations. Limit collection size before an
expensive batch lookup when the request contract specifies that limit.

## Batch references

Collect IDs, make a batch request if the dependency supports it, and iterate the
original collection to identify unavailable references. Preserve original indices
even if lookup IDs were deduplicated. Compose the error path from the caller's
prefix, the collection property from the external contract, and the original
index. Scope the supplied validator to the collection before adding indices. If
the caller already scoped it to that collection, add only the indices; do not
repeat the property or rebuild the validator from the root.

This example's caller supplies `input`; the function owns `references`:

```go
package main

import (
    "context"
    "fmt"

    "github.com/muonsoft/validation"
)

var ErrMissing = validation.NewError("missing_reference", "Reference is missing.")

func missingReferences(ctx context.Context, v *validation.Validator, ids []string, found map[string]bool) error {
    list := v.AtProperty("references").BuildViolationList(ctx)
    for i, id := range ids {
        if !found[id] {
            list.AddViolation(ErrMissing, ErrMissing.Message(), validation.ArrayIndex(i))
        }
    }
    return list.Create().AsError()
}

func main() {
    v, err := validation.NewValidator()
    if err != nil { panic(err) }
    err = missingReferences(context.Background(), v.AtProperty("input"),
        []string{"missing", "ok", "missing"}, map[string]bool{"ok": true})
    violations, ok := validation.UnwrapViolations(err)
    if !ok { panic(err) }
    for _, violation := range violations.AsSlice() {
        fmt.Println(violation.PropertyPath().String())
    }
}
// Output:
// input.references[0]
// input.references[2]
```

`Create().AsError()` returns nil for an empty list. Existence, ownership, and
availability checks use the application's disclosure rules.

Distinguish a missing reference from timeout, cancellation, and connection errors.
Return technical errors with `%w`, preserving `errors.Is`. Pass the original
context to the dependency. Do not convert every lookup failure to “not found”.

For uniqueness on update, exclude the current object's identity and preserve the
scope (for example, an account or tenant). A preflight lookup improves feedback
but cannot prevent races; keep the storage guarantee and map an actual conflict
according to the application's error contract. Validation itself should not save.

Tests should count validation lookups and writes separately from prerequisite
loads. A necessary load before validating an update is not a violation of the
two-stage pattern.
