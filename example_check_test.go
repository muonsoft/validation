package validation_test

import (
	"context"
	"errors"
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

var ErrUnsupportedCommodity = errors.New("unsupported commodity")

func ExampleCheckProperty() {
	outlet := Outlet{
		Type:          "offline",
		MainCommodity: DigitalMovie{name: "Digital movie"},
	}

	err := validator.Validate(
		context.Background(),
		validation.
			CheckProperty("mainCommodity", outlet.MainCommodity.Supports(outlet.Type)).
			WithError(ErrUnsupportedCommodity).
			WithMessage(
				`Commodity "{{ value }}" cannot be sold at outlet.`,
				validation.TemplateParameter{Key: "{{ value }}", Value: outlet.MainCommodity.Name()},
			),
	)

	if violations, ok := validation.UnwrapViolationList(err); ok {
		for violation := violations.First(); violation != nil; violation = violation.Next() {
			fmt.Println("violation underlying error:", violation.Unwrap())
			fmt.Println(violation)
		}
	}
	fmt.Println("errors.Is(err, ErrUnsupportedCommodity) =", errors.Is(err, ErrUnsupportedCommodity))
	// Output:
	// violation underlying error: unsupported commodity
	// violation at "mainCommodity": "Commodity "Digital movie" cannot be sold at outlet."
	// errors.Is(err, ErrUnsupportedCommodity) = true
}
