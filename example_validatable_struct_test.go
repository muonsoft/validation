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
		validation.CountableProperty("tags", len(p.Tags), it.HasMinCount(5)),
		validation.StringsProperty("tags", p.Tags, it.HasUniqueValues()),
		validation.EachStringProperty("tags", p.Tags, it.IsNotBlank()),
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
		Tags: []string{"device", "", "phone", "device"},
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
	// violation at 'tags': This collection should contain 5 elements or more.
	// violation at 'tags': This collection should contain only unique elements.
	// violation at 'tags[1]': This value should not be blank.
	// violation at 'components[0].name': This value should not be blank.
	// violation at 'components[0].tags': This collection should contain 1 element or more.
}
