# Validation flow

Choose the dependency boundary deliberately:

- `v.Validate(ctx, args...)` and `validation.All(args...)` collect independent
  violations. A non-violation error stops execution and is returned instead.
- `validation.Sequentially(args...)` stops at the first argument with violations.
  Wrap a whole independent stage in `All` or `Valid` to collect that stage fully.
- A `Sequentially` nested among sibling arguments stops only its own branch on
  violations. Other independent branches still run.
- Two explicit calls separated by `if err != nil { return err }` are appropriate
  for a boundary between stages, but lose completeness when used per independent field.

```go
package main

import (
    "context"
    "fmt"

    "github.com/muonsoft/validation"
    "github.com/muonsoft/validation/it"
)

func main() {
    v, err := validation.NewValidator()
    if err != nil { panic(err) }
    calls := 0
    dependent := validation.ValidatableFunc(func(ctx context.Context, v *validation.Validator) error {
        calls++
        return nil
    })
    err = v.Validate(context.Background(), validation.Sequentially(
        validation.All(
            validation.StringProperty("name", "", it.IsNotBlank()),
            validation.NumberProperty("quantity", -1, it.IsPositive[int]()),
        ),
        validation.Valid(dependent),
    ))
    violations, ok := validation.UnwrapViolations(err)
    if !ok { panic(err) }
    fmt.Println(violations.Len(), calls)
}
// Output:
// 2 0
```

## Go evaluates arguments before calling a function

`CheckNoViolations(check(ctx))` receives an already computed error. Putting it
inside `Sequentially` does not defer `check`. Likewise `Filter(checkA(), checkB())`
evaluates both calls first. These APIs are useful for merging existing results,
not for scheduling operations. Defer work with a constraint or `ValidatableFunc`.

`StringProperty("note", *note, ...).When(note != nil)` still dereferences `note`
before `When` runs. Prefer `NilStringProperty`; when constructing a constraint
requires dereferencing another optional value, use an `if` to append arguments
only when that value exists. A helper returning `[]validation.Argument` is fine.

Do not put side effects in argument constructors or boolean expressions supplied
to `Check`. Keep checks read-only; perform persistence after validation succeeds.
Asynchronous validation is not required for eager accumulation.
