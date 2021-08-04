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
	"github.com/stretchr/testify/assert"
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

	validationtest.AssertIsViolationList(t, err, func(t *testing.T, violations []validation.Violation) bool {
		t.Helper()
		if assert.Len(t, violations, 4) {
			assert.Equal(t, code.NotBlank, violations[0].Code())
			assert.Equal(t, "name", violations[0].PropertyPath().String())
			assert.Equal(t, code.CountTooFew, violations[1].Code())
			assert.Equal(t, "tags", violations[1].PropertyPath().String())
			assert.Equal(t, code.NotBlank, violations[2].Code())
			assert.Equal(t, "components[0].name", violations[2].PropertyPath().String())
			assert.Equal(t, code.CountTooFew, violations[3].Code())
			assert.Equal(t, "components[0].tags", violations[3].PropertyPath().String())
		}
		return true
	})
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

	assertHasOneViolationAtPath(code.NotBlank, message.NotBlank, "top.value")(t, err)
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

	assertHasOneViolationAtPath(code.NotBlank, message.NotBlank, "top.value")(t, err)
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

	validationtest.AssertIsViolationList(t, err, func(t *testing.T, violations []validation.Violation) bool {
		t.Helper()
		if assert.Len(t, violations, 4) {
			assert.Equal(t, "top.intValue", violations[0].PropertyPath().String())
			assert.Equal(t, "top.floatValue", violations[1].PropertyPath().String())
			assert.Equal(t, "top.stringValue", violations[2].PropertyPath().String())
			assert.Equal(t, "top.structValue.value", violations[3].PropertyPath().String())
		}
		return true
	})
}
