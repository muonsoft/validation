package validation_test

import (
	"context"
	"fmt"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/validator"
)

// It is recommended to make a custom constraint to check for nesting limit.
type NestingLimitConstraint struct {
	limit int
}

func (c NestingLimitConstraint) SetUp() error {
	return nil
}

func (c NestingLimitConstraint) Name() string {
	return "NestingLimitConstraint"
}

func (c NestingLimitConstraint) ValidateProperty(property *Property, scope validation.Scope) error {
	// You can read any passed context value from scope.
	level, ok := scope.Context().Value(nestingLevelKey).(int)
	if !ok {
		// Don't forget to handle missing value.
		return fmt.Errorf("nesting level not found in context")
	}

	if level >= c.limit {
		return scope.
			BuildViolation("nestingLimitReached", "Maximum nesting level reached.").
			CreateViolation()
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
	validation.Constraint
	ValidateProperty(property *Property, scope validation.Scope) error
}

// To create your own functional argument for validation simply create a function with
// a typed value and use the validation.NewArgument constructor.
func PropertyArgument(property *Property, options ...validation.Option) validation.Argument {
	return validation.NewArgument(options, func(constraint validation.Constraint, scope validation.Scope) error {
		if c, ok := constraint.(PropertyConstraint); ok {
			return c.ValidateProperty(property, scope)
		}
		// If you want to use built-in constraints for checking for nil or empty values
		// such as it.IsNil() or it.IsBlank().
		if c, ok := constraint.(validation.NilConstraint); ok {
			if property == nil {
				return c.ValidateNil(scope)
			}
			return nil
		}

		return validation.NewInapplicableConstraintError(constraint, "Property")
	})
}

type recursionKey string

const nestingLevelKey recursionKey = "nestingLevel"

func (p Property) Validate(ctx context.Context, validator *validation.Validator) error {
	return validator.Validate(
		// Incrementing nesting level in context with special function.
		contextWithNextNestingLevel(ctx),
		// Executing validation for maximum nesting level of properties.
		PropertyArgument(&p, ItIsNotDeeperThan(3)),
		validation.StringProperty("name", p.Name, it.IsNotBlank()),
		// This should run recursive validation for properties.
		validation.IterableProperty("properties", p.Properties),
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

	err := validator.Validate(context.Background(), validation.Iterable(properties))

	fmt.Println(err)
	// Output:
	// violation at '[0].properties[0].properties[0].properties[0]': Maximum nesting level reached.
}
