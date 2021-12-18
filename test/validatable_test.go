package test

import (
	"context"
	"testing"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/code"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/message"
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
		validation.String(
			p.Name,
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

func (c Component) Validate(ctx context.Context, validator *validation.Validator) error {
	return validator.Validate(
		ctx,
		validation.String(
			c.Name,
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

func TestValidateValue_WhenValidatableString_ExpectValidationExecutedWithPassedOptionsWithoutConstraints(t *testing.T) {
	validatable := mockValidatableString{value: ""}

	err := validator.Validate(
		context.Background(),
		validation.Value(
			validatable,
			validation.PropertyName("top"),
			it.IsNotBlank().Message("ignored"),
		),
	)

	assertHasOneViolationAtPath(code.NotBlank, message.Templates[code.NotBlank], "top.value")(t, err)
}

func TestValidateValidatable_WhenValidatableString_ExpectValidationExecutedWithPassedOptionsWithoutConstraints(t *testing.T) {
	validatable := mockValidatableString{value: ""}

	err := validator.Validate(
		context.Background(),
		validation.Valid(
			validatable,
			validation.PropertyName("top"),
			it.IsNotBlank().Message("ignored"),
		),
	)

	assertHasOneViolationAtPath(code.NotBlank, message.Templates[code.NotBlank], "top.value")(t, err)
}

func TestValidateValue_WhenValidatableStruct_ExpectValidationExecutedWithPassedOptionsWithoutConstraints(t *testing.T) {
	validatable := mockValidatableStruct{}

	err := validator.Validate(
		context.Background(),
		validation.Value(
			validatable,
			validation.PropertyName("top"),
			it.IsNotBlank().Message("ignored"),
		),
	)

	validationtest.Assert(t, err).IsViolationList().WithAttributes(
		validationtest.ViolationAttributes{PropertyPath: "top.intValue"},
		validationtest.ViolationAttributes{PropertyPath: "top.floatValue"},
		validationtest.ViolationAttributes{PropertyPath: "top.stringValue"},
		validationtest.ViolationAttributes{PropertyPath: "top.structValue.value"},
	)
}
