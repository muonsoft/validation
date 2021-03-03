package examples

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

func (p Product) Validate(scope validation.Scope) error {
	return validator.InScope(scope).Validate(
		validation.String(
			&p.Name,
			validation.PropertyName("name"),
			it.IsNotBlank(),
		),
		validation.Iterable(
			p.Tags,
			validation.PropertyName("tags"),
			it.HasMinCount(1),
		),
		validation.Iterable(
			p.Components,
			validation.PropertyName("components"),
			it.HasMinCount(1),
		),
	)
}

type Component struct {
	ID   int
	Name string
	Tags []string
}

func (c Component) Validate(scope validation.Scope) error {
	return validator.InScope(scope).Validate(
		validation.String(
			&c.Name,
			validation.PropertyName("name"),
			it.IsNotBlank(),
		),
		validation.Iterable(
			c.Tags,
			validation.PropertyName("tags"),
			it.HasMinCount(1),
		),
	)
}

func ExampleValidateValidatable_withSingletonValidator() {
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

	violations := err.(validation.ViolationList)
	for _, violation := range violations {
		fmt.Println(violation.Error())
	}
	// Output:
	// violation at 'name': This value should not be blank.
	// violation at 'tags': This collection should contain 1 element or more.
	// violation at 'components[0].name': This value should not be blank.
	// violation at 'components[0].tags': This collection should contain 1 element or more.
}
