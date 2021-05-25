package validation_test

import (
	"context"
	"fmt"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/validator"
)

type Brand struct {
	Name string
}

type BrandRepository struct {
	brands []Brand
}

func (repository *BrandRepository) FindByName(ctx context.Context, name string) ([]Brand, error) {
	found := make([]Brand, 0)

	for _, brand := range repository.brands {
		if brand.Name == name {
			found = append(found, brand)
		}
	}

	return found, nil
}

// You can declare you own constraint interface to create custom constraints.
type BrandConstraint interface {
	validation.Constraint
	ValidateBrand(brand *Brand, scope validation.Scope) error
}

// To create your own functional argument for validation simply create a function with
// a typed value and use the validation.NewArgument constructor.
func BrandArgument(brand *Brand, options ...validation.Option) validation.Argument {
	return validation.NewArgument(options, func(constraint validation.Constraint, scope validation.Scope) error {
		if c, ok := constraint.(BrandConstraint); ok {
			return c.ValidateBrand(brand, scope)
		}
		// If you want to use built-in constraints for checking for nil or empty values
		// such as it.IsNil() or it.IsBlank().
		if c, ok := constraint.(validation.NilConstraint); ok {
			if brand == nil {
				return c.ValidateNil(scope)
			}
			return nil
		}

		return validation.NewInapplicableConstraintError(constraint, "Brand")
	})
}

// UniqueBrandConstraint implements BrandConstraint.
type UniqueBrandConstraint struct {
	brands *BrandRepository
}

func (c *UniqueBrandConstraint) SetUp() error {
	return nil
}

func (c *UniqueBrandConstraint) Name() string {
	return "UniqueBrandConstraint"
}

func (c *UniqueBrandConstraint) ValidateBrand(brand *Brand, scope validation.Scope) error {
	// usually, you should ignore empty values
	// to check for an empty value you should use it.NotBlankConstraint
	if brand == nil {
		return nil
	}

	// you can pass the context value from the scope
	brands, err := c.brands.FindByName(scope.Context(), brand.Name)
	// here you can return a service error so that the validation process
	// is stopped immediately
	if err != nil {
		return err
	}
	if len(brands) == 0 {
		return nil
	}

	// use the scope to build violation with translations
	return scope.
		BuildViolation("notUniqueBrand", `Brand with name "{{ name }}" already exists.`).
		// you can inject parameter value to the message here
		AddParameter("{{ name }}", brand.Name).
		CreateViolation()
}

func ExampleNewArgument_customArgumentConstraintValidator() {
	repository := &BrandRepository{brands: []Brand{{"Apple"}, {"Orange"}}}
	isEntityUnique := &UniqueBrandConstraint{brands: repository}

	brand := Brand{Name: "Apple"}
	ctx := context.WithValue(context.Background(), exampleKey, "value")

	err := validator.Validate(
		// you can pass here the context value to the validation scope
		validation.Context(ctx),
		BrandArgument(&brand, it.IsNotBlank(), isEntityUnique),
	)

	fmt.Println(err)
	// Output:
	// violation: Brand with name "Apple" already exists.
}
