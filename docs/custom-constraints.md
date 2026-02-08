# Creating Custom Constraints

This guide describes how to build custom validation constraints, from simple checks based on a boolean function to complex constraints with options, dependencies, and domain types.

## Interfaces overview

To create a custom constraint, implement one of these interfaces:

| Interface | Use case |
|-----------|----------|
| `BoolConstraint` | Boolean values |
| `NumberConstraint[T]` | Numeric values (int, float, etc.) |
| `StringConstraint` | String values |
| `ComparableConstraint[T]` | Generic comparable values |
| `ComparablesConstraint[T]` | Slices of comparable values |
| `CountableConstraint` | Count of elements (e.g. `len(slice)`) |
| `TimeConstraint` | `time.Time` values |
| `NilConstraint` | Nil checks for nillable types |
| `SliceConstraint[T]` | Generic slice validation |
| `Constraint[T]` | Any type (used with [This], [Each], [EachProperty]) |

You can implement several interfaces on the same type (e.g. a constraint that validates both strings and count).

---

## Simple checks

### 1. Function-based string constraint (`OfStringBy`)

When the rule is a pure `func(string) bool`, use [OfStringBy] and a check from the `is` package (or your own function):

```go
// IsMyFormat validates that the value matches the desired format.
func IsMyFormat() validation.StringFuncConstraint {
	return validation.OfStringBy(is.MyFormat).
		WithError(validation.ErrInvalidMyFormat).
		WithMessage(validation.ErrInvalidMyFormat.Message())
}
```

- No struct needed.
- Empty/nil values are skipped by default; use `it.IsNotBlank()` (or similar) to require non-empty.
- Use [WithError] and [WithMessage] for violations and translations.

See `it.IsJSON`, `it.IsNumeric`, `it.IsInteger` in the codebase.

### 2. Struct implementing a typed interface

For a self-contained rule (e.g. regex or fixed logic), define a struct and implement the appropriate interface (e.g. [StringConstraint]):

```go
var ErrNotNumeric = errors.New("not numeric")

type NumericConstraint struct {
	matcher *regexp.Regexp
}

func IsNumeric() NumericConstraint {
	return NumericConstraint{matcher: regexp.MustCompile("^[0-9]+$")}
}

func (c NumericConstraint) ValidateString(ctx context.Context, validator *validation.Validator, value *string) error {
	// Skip empty; use it.NotBlank() to require non-empty.
	if value == nil || *value == "" {
		return nil
	}
	if c.matcher.MatchString(*value) {
		return nil
	}
	return validator.CreateViolation(ctx, ErrNotNumeric, "This value should be numeric.")
}
```

Usage with other constraints:

```go
err := validator.Validate(ctx, validation.String(s, it.IsNotBlank(), IsNumeric()))
```

**Reporting violations:**

- Use [BuildViolation] when the message comes from the message/translation system (e.g. `validation.ErrNotNumeric` and its template). That way the validator can translate it.
- Use [CreateViolation] for a one-off, non-translatable message (as in the snippet above).

See [example_custom_constraint_test.go](https://pkg.go.dev/github.com/muonsoft/validation#example-Validator.Validate-CustomConstraint) in the repo.

---

## Complex checks

### 3. Configurable constraints (options, multiple messages)

For constraints with options (min/max, regex, multiple error types), use a struct with constructor(s) and fluent options:

- Store configuration (regex, limits, flags).
- Implement the right interface (e.g. [StringConstraint]).
- Skip when `value == nil` or empty unless the constraint is explicitly “required” (e.g. use `it.IsNotBlank()` separately).
- Use [BuildViolation] with message templates and [WithParameters] / [WithParameter] so messages are translatable.
- Optionally support [When] and [WhenGroups] for conditional execution.

Example shape (see `it.RegexpConstraint`, `it.LengthConstraint` in the codebase):

```go
type MyConstraint struct {
	isIgnored       bool
	groups          []string
	err             error
	messageTemplate string
	messageParameters validation.TemplateParameterList
	// ... your options
}

func IsMy(option MyOption) MyConstraint { ... }

func (c MyConstraint) WithError(err error) MyConstraint { ... }
func (c MyConstraint) WithMessage(template string, parameters ...validation.TemplateParameter) MyConstraint { ... }
func (c MyConstraint) When(condition bool) MyConstraint { ... }
func (c MyConstraint) WhenGroups(groups ...string) MyConstraint { ... }

func (c MyConstraint) ValidateString(ctx context.Context, validator *validation.Validator, value *string) error {
	if c.isIgnored || validator.IsIgnoredForGroups(c.groups...) || value == nil || *value == "" {
		return nil
	}
	// ... your logic
	return validator.BuildViolation(ctx, c.err, c.messageTemplate).
		WithParameters(c.messageParameters.Prepend(validation.TemplateParameter{Key: "{{ value }}", Value: *value})...).
		Create()
}
```

If the same constraint is used with [Each] or [This], also implement [Constraint][T] by delegating to the typed method:

```go
func (c RegexpConstraint) Validate(ctx context.Context, validator *validation.Validator, v string) error {
	return c.ValidateString(ctx, validator, &v)
}
```

### 4. Custom types with `Constraint[T]` and `This`

For domain types (structs, not primitives), use the generic [Constraint][T] and [This]:

1. Define `Constraint[YourType]` (e.g. a struct with a `Validate(ctx, validator, v YourType) error` method).
2. Pass the value and constraints via [This][T](value, constraint1, constraint2, ...).
3. In the constraint, use [BuildViolation] (or [CreateViolation]) and optional [WithParameter] for message placeholders.

Example: uniqueness check that needs a repository (see [example_custom_argument_constraint_test.go](https://pkg.go.dev/github.com/muonsoft/validation#example-NewArgument-CustomArgumentConstraintValidator)):

```go
type Brand struct { Name string }

type UniqueBrandConstraint struct {
	brands *BrandRepository
}

func (c *UniqueBrandConstraint) Validate(ctx context.Context, validator *validation.Validator, brand *Brand) error {
	if brand == nil {
		return nil
	}
	brands, err := c.brands.FindByName(ctx, brand.Name)
	if err != nil {
		return err  // Stops validation
	}
	if len(brands) == 0 {
		return nil
	}
	return validator.
		BuildViolation(ctx, ErrNotUniqueBrand, `Brand with name "{{ name }}" already exists.`).
		WithParameter("{{ name }}", brand.Name).
		Create()
}
```

Custom argument constructor (optional but convenient):

```go
func ValidBrand(brand *Brand, constraints ...validation.Constraint[*Brand]) validation.ValidatorArgument {
	return validation.This[*Brand](brand, constraints...)
}

// Usage
err := validator.Validate(ctx, ValidBrand(&brand, isUnique))
```

### 5. Functions as constraints (`Func[T]`)

[Func] wraps a function so it implements [Constraint][T]. Use it for one-off or inline rules and with [Each] / [EachProperty]:

```go
notBlank := validation.Func[string](func(ctx context.Context, v *validation.Validator, s string) error {
	return v.Validate(ctx, validation.String(s, it.IsNotBlank()))
})
err := validator.Validate(ctx, validation.This("", notBlank))
```

With [Each] and a custom element type:

```go
type Item struct { Code string }

validCode := validation.Func[Item](func(ctx context.Context, v *validation.Validator, item Item) error {
	return v.Validate(ctx, validation.StringProperty("code", item.Code, it.IsNotBlank()))
})
items := []Item{{Code: "A"}, {Code: ""}}
err := validator.Validate(ctx, validation.Each(items, validCode))
// Violation path: "[1].code"
```

See [ExampleFunc](https://pkg.go.dev/github.com/muonsoft/validation#example-Func) and [ExampleEach_withCustomType](https://pkg.go.dev/github.com/muonsoft/validation#example-Func) in the repo.

### 6. Slices: `Each` / `EachProperty` and typed helpers

- **Primitives:** Use [EachString], [EachNumber], [EachComparable] (and `*Property` variants) with the corresponding constraint types. Paths include the index (e.g. `[0]`, `[1]`).
- **Any type:** Use [Each] or [EachProperty] with [Constraint][E]. Constraints from `it` that implement [Constraint][T] (e.g. string constraints with `Validate`) can be passed directly.
- **Custom elements:** Use [Each](items, validation.Func[E](...)) as in the `Item` example above. Paths look like `[index]` or `property[index].field`.

---

## Summary

| Need | Approach |
|------|----------|
| Simple string rule `func(string) bool` | [OfStringBy] + [WithError] / [WithMessage] |
| One-off string/struct rule, no options | Struct implementing [StringConstraint] or [Constraint][T], [CreateViolation] or [BuildViolation] |
| Configurable constraint (min/max, regex, etc.) | Struct with options, [BuildViolation], optional [When] / [WhenGroups] |
| Domain type (e.g. Brand) | [Constraint][*YourType], [This], optional custom argument (e.g. `ValidBrand`) |
| Inline / ad-hoc rule | [Func][T] |
| Validate each element of a slice | [Each] / [EachProperty] with [Constraint][E] or [Func][E]; for primitives use [EachString] etc. |

---

## References and examples

- [Custom static constraint (string)](https://pkg.go.dev/github.com/muonsoft/validation#example-Validator.Validate-CustomConstraint) — struct implementing [StringConstraint].
- [Custom constraint as a service](https://pkg.go.dev/github.com/muonsoft/validation#example-Validator.GetConstraint-CustomServiceConstraint) — constraint resolved from a validator/service.
- [Custom argument for domain type](https://pkg.go.dev/github.com/muonsoft/validation#example-NewArgument-CustomArgumentConstraintValidator) — [This] + [Constraint][*Brand], repository dependency.
- [Func, Each, EachProperty](https://pkg.go.dev/github.com/muonsoft/validation#example-Func) — function as constraint and validating slices of any type.

For adding new constraints to the library (messages, translations, tests), see [.cursor/skills/validation-add-constraint/SKILL.md](../.cursor/skills/validation-add-constraint/SKILL.md).
