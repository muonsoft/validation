package test

import (
	"context"
	"testing"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/code"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/message"
)

type property struct {
	Name  string
	Type  string
	Value []*property
}

func (p property) Validate(ctx context.Context, validator *validation.Validator) error {
	arguments := []validation.Argument{
		validation.StringProperty(
			"name",
			p.Name,
			it.IsNotBlank(),
		),
		validation.StringProperty(
			"type",
			p.Type,
			it.IsNotBlank(),
		),
	}

	if p.Value != nil {
		arguments = append(arguments,
			validation.IterableProperty(
				"value",
				p.Value,
				it.HasMinCount(1),
			),
		)
	}

	return validator.Validate(ctx, arguments...)
}

func TestValidate_AtProperty_WhenGivenRecursiveProperties_ExpectViolationWithProperty(t *testing.T) {
	validator := newValidator(t)
	properties := []*property{
		{
			Name: "first level",
			Type: "first type",
			Value: []*property{
				{
					Name: "second level",
					Type: "second type",
					Value: []*property{
						{
							Name:  "",
							Type:  "blank",
							Value: nil,
						},
					},
				},
			},
		},
	}

	err := validator.Validate(context.Background(), validation.Iterable(properties))

	assertHasOneViolationAtPath(code.NotBlank, message.NotBlank, "[0].value[0].value[0].name")(t, err)
}

func TestValidate_WhenPathIsSetViaOptions_ExpectViolationAtPath(t *testing.T) {
	validator := newValidator(t)
	v := ""

	err := validator.Validate(
		context.Background(),
		validation.String(
			v,
			validation.PropertyName("properties"),
			validation.ArrayIndex(0),
			validation.PropertyName("value"),
			it.IsNotBlank(),
		),
	)

	assertHasOneViolationAtPath(code.NotBlank, message.NotBlank, customPath)(t, err)
}

func TestValidate_AtProperty_WhenGivenProperty_ExpectViolationWithProperty(t *testing.T) {
	validator := newValidator(t)

	err := validator.
		AtProperty("property").
		Validate(context.Background(), validation.String("", it.IsNotBlank()))

	assertHasOneViolationAtPath(code.NotBlank, message.NotBlank, "property")(t, err)
}

func TestValidate_AtIndex_WhenGivenIndex_ExpectViolationWithIndex(t *testing.T) {
	validator := newValidator(t)

	err := validator.
		AtIndex(1).
		Validate(context.Background(), validation.String("", it.IsNotBlank()))

	assertHasOneViolationAtPath(code.NotBlank, message.NotBlank, "[1]")(t, err)
}
