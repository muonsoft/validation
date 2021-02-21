package internal

import (
	"fmt"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
)

type Product struct {
	Name       string
	Tags       []string
	Components []Component
}

func (p Product) Validate(options ...validation.Option) error {
	validator, err := validation.WithOptions(options...)
	if err != nil {
		return err
	}

	return validation.Filter(
		validator.ValidateString(
			&p.Name,
			validation.PropertyName("name"),
			it.IsNotBlank(),
		),
		validator.ValidateIterable(
			p.Components,
			validation.PropertyName("components"),
		),
	)
}

type Component struct {
	ID   int
	Name string
}

func (c Component) Validate(options ...validation.Option) error {
	validator, err := validation.WithOptions(options...)
	if err != nil {
		return err
	}

	return validation.Filter(
		validator.ValidateString(
			&c.Name,
			validation.PropertyName("name"),
			it.IsNotBlank(),
		),
	)
}

func ExampleValidate() {
	p := Product{
		Name: "",
		Components: []Component{
			{
				ID:   1,
				Name: "",
			},
		},
	}

	err := validation.Validate(p)

	fmt.Println(err)
	// Output: violation at 'name': This value should not be blank.; violation at 'components[0].name': This value should not be blank.
}
