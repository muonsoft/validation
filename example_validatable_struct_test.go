package validation_test

import (
	"context"
	"fmt"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/validator"
)

type Product struct {
	Name       string
	Tags       []string
	Components []Component
}

func (p Product) Validate(ctx context.Context, validator *validation.Validator) error {
	return validator.Validate(
		ctx,
		validation.StringProperty("name", p.Name, it.IsNotBlank()),
		validation.AtProperty(
			"tags",
			validation.Countable(len(p.Tags), it.HasMinCount(5)),
			validation.Comparables[string](p.Tags, it.HasUniqueValues[string]()),
			validation.EachString(p.Tags, it.IsNotBlank()),
		),
		validation.AtProperty(
			"components",
			validation.Countable(len(p.Components), it.HasMinCount(1)),
			// this runs validation on each of the components
			validation.ValidSlice(p.Components),
		),
	)
}

type Component struct {
	ID   int
	Name string
	Tags []string
}

func (c Component) Validate(ctx context.Context, validator *validation.Validator) error {
	return validator.Validate(
		ctx,
		validation.StringProperty("name", c.Name, it.IsNotBlank()),
		validation.CountableProperty("tags", len(c.Tags), it.HasMinCount(1)),
	)
}

func ExampleValid_validatableStruct() {
	p := Product{
		Name: "",
		Tags: []string{"device", "", "phone", "device"},
		Components: []Component{
			{
				ID:   1,
				Name: "",
			},
		},
	}

	err := validator.Validate(context.Background(), validation.Valid(p))

	if violations, ok := validation.UnwrapViolationList(err); ok {
		for violation := violations.First(); violation != nil; violation = violation.Next() {
			fmt.Println(violation)
		}
	}
	// Output:
	// violation at "name": "This value should not be blank."
	// violation at "tags": "This collection should contain 5 elements or more."
	// violation at "tags": "This collection should contain only unique elements."
	// violation at "tags[1]": "This value should not be blank."
	// violation at "components[0].name": "This value should not be blank."
	// violation at "components[0].tags": "This collection should contain 1 element or more."
}
