package validation_test

import (
	"context"
	"fmt"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/validator"
)

type Outlet struct {
	Type          string
	MainCommodity OutletCommodity
}

type OutletCommodity interface {
	Name() string
	Supports(outletType string) bool
}

type DigitalMovie struct {
	name string
}

func (m DigitalMovie) Name() string {
	return m.name
}

func (m DigitalMovie) Supports(outletType string) bool {
	return outletType == "digital"
}

func ExampleCheckProperty() {
	outlet := Outlet{
		Type:          "offline",
		MainCommodity: DigitalMovie{name: "Digital movie"},
	}

	err := validator.Validate(
		context.Background(),
		validation.
			CheckProperty("mainCommodity", outlet.MainCommodity.Supports(outlet.Type)).
			Code("unsupportedCommodity").
			Message(
				`Commodity "{{ value }}" cannot be sold at outlet.`,
				validation.TemplateParameter{Key: "{{ value }}", Value: outlet.MainCommodity.Name()},
			),
	)

	if violations, ok := validation.UnwrapViolationList(err); ok {
		for violation := violations.First(); violation != nil; violation = violation.Next() {
			fmt.Println("violation code:", violation.Code())
			fmt.Println(violation)
		}
	}
	// Output:
	// violation code: unsupportedCommodity
	// violation at 'mainCommodity': Commodity "Digital movie" cannot be sold at outlet.
}
