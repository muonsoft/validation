package validation_test

import (
	"context"
	"errors"
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

// To create your own functional argument for validation simply create a function with
// a typed value and use the validation.This constructor.
func ValidBrand(brand *Brand, constraints ...validation.Constraint[*Brand]) validation.ValidatorArgument {
	return validation.This[*Brand](brand, constraints...)
}

var ErrNotUniqueBrand = errors.New("not unique brand")

// UniqueBrandConstraint implements BrandConstraint.
type UniqueBrandConstraint struct {
	brands *BrandRepository
}

func (c *UniqueBrandConstraint) Validate(ctx context.Context, validator *validation.Validator, brand *Brand) error {
	// usually, you should ignore empty values
	// to check for an empty value you should use it.NotBlankConstraint
	if brand == nil {
		return nil
	}

	brands, err := c.brands.FindByName(ctx, brand.Name)
	// here you can return a service error so that the validation process
	// is stopped immediately
	if err != nil {
		return err
	}
	if len(brands) == 0 {
		return nil
	}

	// use the validator to build violation with translations
	return validator.
		BuildViolation(ctx, ErrNotUniqueBrand, `Brand with name "{{ name }}" already exists.`).
		// you can inject parameter value to the message here
		WithParameter("{{ name }}", brand.Name).
		Create()
}

func ExampleThis_customArgumentConstraintValidator() {
	repository := &BrandRepository{brands: []Brand{{"Apple"}, {"Orange"}}}
	isUnique := &UniqueBrandConstraint{brands: repository}

	brand := Brand{Name: "Apple"}

	err := validator.Validate(
		// you can pass here the context value to the validation context
		context.WithValue(context.Background(), exampleKey, "value"),
		ValidBrand(&brand, isUnique),
		// it is full equivalent of
		// validation.This[*Brand](&brand, isUnique),
	)

	fmt.Println(err)
	fmt.Println("errors.Is(err, ErrNotUniqueBrand) =", errors.Is(err, ErrNotUniqueBrand))
	// Output:
	// violation: Brand with name "Apple" already exists.
	// errors.Is(err, ErrNotUniqueBrand) = true
}
