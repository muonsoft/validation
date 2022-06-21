package validation_test

import (
	"context"
	"fmt"

	"github.com/muonsoft/validation"
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
	ValidateBrand(brand *Brand, scope validation.Scope) error
}

// To create your own functional argument for validation simply create a function with
// a typed value and use the validation.NewArgument constructor.
func ValidBrand(brand *Brand, constraints ...BrandConstraint) validation.ValidatorArgument {
	return validation.NewArgument(func(scope validation.Scope) (*validation.ViolationList, error) {
		violations := validation.NewViolationList()

		for i := range constraints {
			err := violations.AppendFromError(constraints[i].ValidateBrand(brand, scope))
			if err != nil {
				return nil, err
			}
		}

		return violations, nil
	})
}

// UniqueBrandConstraint implements BrandConstraint.
type UniqueBrandConstraint struct {
	brands *BrandRepository
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
		WithParameter("{{ name }}", brand.Name).
		Create()
}

func ExampleNewArgument_customArgumentConstraintValidator() {
	repository := &BrandRepository{brands: []Brand{{"Apple"}, {"Orange"}}}
	isEntityUnique := &UniqueBrandConstraint{brands: repository}

	brand := Brand{Name: "Apple"}

	err := validator.Validate(
		// you can pass here the context value to the validation scope
		context.WithValue(context.Background(), exampleKey, "value"),
		ValidBrand(&brand, isEntityUnique),
	)

	fmt.Println(err)
	// Output:
	// violation: Brand with name "Apple" already exists.
}
