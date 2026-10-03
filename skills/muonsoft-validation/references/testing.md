# Testing application validation

Test through the exported entrypoint used by the application. Prefer black-box
tests and the project's existing testing framework. `validationtest` provides
assertions for codes, paths, and messages; standard Go assertions also work.

For a DTO or entity, cover valid input, multiple simultaneous independent errors,
nil and empty inputs, and applicable boundaries. For collections, combine a parent
error with errors in more than one element. Verify indices and relative prefixes;
do not rely on map traversal order.

This complete program demonstrates the same observable assertions a Go test needs:

```go
package main

import (
    "context"
    "errors"
    "fmt"

    "github.com/muonsoft/validation"
    "github.com/muonsoft/validation/it"
)

func main() {
    v, err := validation.NewValidator()
    if err != nil { panic(err) }
    err = v.Validate(context.Background(),
        validation.StringProperty("name", "", it.IsNotBlank()),
        validation.NumberProperty("quantity", 0, it.IsPositive[int]()),
    )
    violations, ok := validation.UnwrapViolations(err)
    if !ok || violations.Len() != 2 { panic("expected both violations") }
    first := violations.First().Violation()
    if !errors.Is(first, validation.ErrIsBlank) || first.PropertyPath().String() != "name" {
        panic("unexpected code or path")
    }
    fmt.Println("both violations preserved")
}
// Output:
// both violations preserved
```

In a test using `validationtest`, the equivalent assertion is
`validationtest.Assert(t, err).IsViolationList().WithAttributes(...)`, supplying
`ViolationAttributes{Error: ..., PropertyPath: ...}` for each expected violation.
Use an explicit path assertion to test an empty root path; empty attribute values
are not compared by `ViolationAttributes`.

For infrastructure checks, use fakes that record inputs, contexts, and call counts:

- Local invalid input: all expected violations; zero dependent lookup/write calls.
- Valid local input with multiple missing references: every required violation,
  original indices, the full collection path with and without a caller prefix,
  and the intended number and contents of batch requests. Allow deduplicated IDs
  when the contract permits them, while checking every original occurrence.
- Timeout/cancellation: the original error remains discoverable with `errors.Is`;
  later checks and writes do not run.
- Update: necessary load occurs, current identity is excluded from uniqueness,
  and the normalized state being validated is the state passed to persistence.
- Adaptation: nested paths, selected groups, and translations survive the adapter.

Test one failing prerequisite alongside a separate failing sibling to distinguish
local branch sequencing from accidental global short-circuiting. An early-return
test involving only one invalid field cannot prove eager accumulation.
