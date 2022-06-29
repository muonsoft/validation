package validation_test

import (
	"context"
	"errors"
	"fmt"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/validator"
)

var ErrNestingLimitReached = errors.New("nesting limit reached")

// It is recommended to make a custom constraint to check for nesting limit.
type NestingLimitConstraint struct {
	limit int
}

func (c NestingLimitConstraint) ValidateProperty(property *Property, scope validation.Scope) error {
	// You can read any passed context value from scope.
	level, ok := scope.Context().Value(nestingLevelKey).(int)
	if !ok {
		// Don't forget to handle missing value.
		return fmt.Errorf("nesting level not found in context")
	}

	if level >= c.limit {
		return scope.CreateViolation(ErrNestingLimitReached, "Maximum nesting level reached.")
	}

	return nil
}

func ItIsNotDeeperThan(limit int) NestingLimitConstraint {
	return NestingLimitConstraint{limit: limit}
}

// Properties can be nested.
type Property struct {
	Name       string
	Properties []Property
}

// You can declare you own constraint interface to create custom constraints.
type PropertyConstraint interface {
	ValidateProperty(property *Property, scope validation.Scope) error
}

// To create your own functional argument for validation simply create a function with
// a typed value and use the validation.NewArgument constructor.
func ValidProperty(property *Property, constraints ...PropertyConstraint) validation.ValidatorArgument {
	return validation.NewArgument(func(scope validation.Scope) (*validation.ViolationList, error) {
		violations := validation.NewViolationList()

		for i := range constraints {
			err := violations.AppendFromError(constraints[i].ValidateProperty(property, scope))
			if err != nil {
				return nil, err
			}
		}

		return violations, nil
	})
}

type recursionKey string

const nestingLevelKey recursionKey = "nestingLevel"

func (p Property) Validate(ctx context.Context, validator *validation.Validator) error {
	return validator.Validate(
		// Incrementing nesting level in context with special function.
		contextWithNextNestingLevel(ctx),
		// Executing validation for maximum nesting level of properties.
		ValidProperty(&p, ItIsNotDeeperThan(3)),
		validation.StringProperty("name", p.Name, it.IsNotBlank()),
		// This should run recursive validation for properties.
		validation.ValidSliceProperty("properties", p.Properties),
	)
}

// This function increments current nesting level.
func contextWithNextNestingLevel(ctx context.Context) context.Context {
	level, ok := ctx.Value(nestingLevelKey).(int)
	if !ok {
		level = -1
	}

	return context.WithValue(ctx, nestingLevelKey, level+1)
}

func ExampleValidator_Validate_usingContextWithRecursion() {
	properties := []Property{
		{
			Name: "top",
			Properties: []Property{
				{
					Name: "middle",
					Properties: []Property{
						{
							Name: "low",
							Properties: []Property{
								// This property should cause a violation.
								{Name: "limited"},
							},
						},
					},
				},
			},
		},
	}

	err := validator.Validate(context.Background(), validation.ValidSlice(properties))

	fmt.Println(err)
	fmt.Println("errors.Is(err, ErrNestingLimitReached) =", errors.Is(err, ErrNestingLimitReached))
	// Output:
	// violation at '[0].properties[0].properties[0].properties[0]': Maximum nesting level reached.
	// errors.Is(err, ErrNestingLimitReached) = true
}
