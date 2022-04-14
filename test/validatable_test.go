package test

import (
	"context"
	"testing"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/code"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/validationtest"
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
		validation.StringProperty(
			"name",
			p.Name,
			it.IsNotBlank(),
		),
		validation.CountableProperty(
			"tags",
			len(p.Tags),
			it.HasMinCount(1),
		),
		validation.CountableProperty(
			"components",
			len(p.Components),
			it.HasMinCount(1),
		),
		validation.ValidSliceProperty(
			"components",
			p.Components,
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
		validation.StringProperty(
			"name",
			c.Name,
			it.IsNotBlank(),
		),
		validation.CountableProperty(
			"tags",
			len(c.Tags),
			it.HasMinCount(1),
		),
	)
}

func TestValidateValue_WhenStructWithComplexRules_ExpectViolations(t *testing.T) {
	p := Product{
		Name: "",
		Components: []Component{
			{
				ID:   1,
				Name: "",
			},
		},
	}

	err := validator.Validate(context.Background(), validation.Valid(p))

	validationtest.Assert(t, err).IsViolationList().WithAttributes(
		validationtest.ViolationAttributes{Code: code.NotBlank, PropertyPath: "name"},
		validationtest.ViolationAttributes{Code: code.CountTooFew, PropertyPath: "tags"},
		validationtest.ViolationAttributes{Code: code.NotBlank, PropertyPath: "components[0].name"},
		validationtest.ViolationAttributes{Code: code.CountTooFew, PropertyPath: "components[0].tags"},
	)
}
