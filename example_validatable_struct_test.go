package validation_test

import (
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

func (p Product) Validate(validator *validation.Validator) error {
	return validator.Validate(
		validation.StringProperty("name", &p.Name, it.IsNotBlank()),
		validation.IterableProperty("tags", p.Tags, it.HasMinCount(1)),
		// this also runs validation on each of the components
		validation.IterableProperty("components", p.Components, it.HasMinCount(1)),
	)
}

type Component struct {
	ID   int
	Name string
	Tags []string
}

func (c Component) Validate(validator *validation.Validator) error {
	return validator.Validate(
		validation.StringProperty("name", &c.Name, it.IsNotBlank()),
		validation.CountableProperty("tags", len(c.Tags), it.HasMinCount(1)),
	)
}

func ExampleValidator_ValidateValidatable_validatableStruct() {
	p := Product{
		Name: "",
		Components: []Component{
			{
				ID:   1,
				Name: "",
			},
		},
	}

	err := validator.ValidateValidatable(p)

	if violations, ok := validation.UnwrapViolationList(err); ok {
		for violation := violations.First(); violation != nil; violation = violation.Next() {
			fmt.Println(violation)
		}
	}
	// Output:
	// violation at 'name': This value should not be blank.
	// violation at 'tags': This collection should contain 1 element or more.
	// violation at 'components[0].name': This value should not be blank.
	// violation at 'components[0].tags': This collection should contain 1 element or more.
}
