---
name: muonsoft-validation
description: Implement, extend, debug, and test Go entity and DTO validation with github.com/muonsoft/validation. Use for Validatable, composed constraints, nested violation paths, or checks backed by repositories and clients in projects using this library.
---

# Apply muonsoft/validation

Read the project's `go.mod`, existing validation entrypoints, error handling, and
validator construction before changing rules. Use the installed API; do not
upgrade dependencies or introduce this library into an unrelated project unless
requested. These examples target v0.19.0 and are also checked against the library
checkout accompanying this skill.

## Default pattern

Put reusable object rules on a type implementing
`Validate(context.Context, *validation.Validator) error`. Pass independent field
checks to one `v.Validate` so the caller receives all their violations.

```go
package main

import (
    "context"
    "fmt"

    "github.com/muonsoft/validation"
    "github.com/muonsoft/validation/it"
)

type Input struct { Name string; Quantity int }

func (x Input) Validate(ctx context.Context, v *validation.Validator) error {
    return v.Validate(ctx,
        validation.StringProperty("name", x.Name, it.IsNotBlank()),
        validation.NumberProperty("quantity", x.Quantity, it.IsPositive[int]()),
    )
}

var _ validation.Validatable = Input{}

func main() {
    v, err := validation.NewValidator()
    if err != nil { panic(err) }
    err = v.ValidateIt(context.Background(), Input{})
    if violations, ok := validation.UnwrapViolations(err); ok {
        fmt.Println(violations.Len())
    } else if err != nil { panic(err) }
}
// Output:
// 2
```

## Preserve these semantics

- Collect independent violations within each stage. Use a stage boundary only
  when later checks require valid earlier input. An ordinary, non-violation error
  stops validation; preserve it instead of reporting invalid user input.
- Reuse the supplied validator in children and adapters. Creating a new validator
  inside them loses the caller's path, groups, translations, and violation factory.
- Paths describe the external data contract. Supply relative property names and
  typed array indices; neither Go fields nor JSON tags automatically define rules.
- Most format constraints permit empty values. Express requiredness separately
  and verify each constraint's nil/empty semantics. `IsNotBlank` does not trim whitespace.
- Validate the state that will actually be used or saved. Follow the application's
  normalization and update semantics; do not silently change input in `Validate`.
- Helpers, conditional argument construction, and manual loops are useful when
  they preserve accumulation, paths, and technical errors. Do not rewrite sound
  code merely to make every check look alike.

## Read by task

| Task | Reference |
| --- | --- |
| Add rules to an entity, command, or query | [Entities and DTOs](references/entity-and-dto.md) |
| Collect errors while sequencing dependent checks | [Validation flow](references/validation-flow.md) |
| Validate children, slices, maps, and unique keys | [Nested objects and collections](references/nested-collections.md) |
| Handle optional fields, enums, conditions, and groups | [Optional and conditional values](references/optional-and-conditional.md) |
| Check references or uniqueness using DB/HTTP | [Infrastructure validation](references/infrastructure-validation.md) |
| Add a reusable rule or adapt an existing method | [Custom constraints](references/custom-constraints.md) |
| Return codes, messages, translations, and paths | [Errors and translations](references/errors-and-translations.md) |
| Verify observable behavior | [Testing](references/testing.md) |

Read only the relevant references. Prefer existing project constraints when their
semantics match; distinguish them from the public `it` package. Use `is` for a
standalone boolean predicate and `validate` for a standalone error check; use
`validation` arguments and `it` for composed object validation.

Before finishing, test multiple simultaneous violations, exact paths, valid input,
and relevant nil/boundary cases. For dependent checks, also test call counts,
technical failures, and that invalid input is not saved.
