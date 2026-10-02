package test

import (
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/message"
	"github.com/muonsoft/validation/validate"
)

var cardSchemeConstraintTestCases = []ConstraintValidationTestCase{
	{
		name:            "CardScheme nil",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsCardScheme(validate.CardSchemeVisa),
		assert:          assertNoError,
	},
	{
		name:            "CardScheme empty",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue(""),
		constraint:      it.IsCardScheme(validate.CardSchemeVisa),
		assert:          assertNoError,
	},
	{
		name:            "CardScheme valid",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("4111111111111111"),
		constraint:      it.IsCardScheme(validate.CardSchemeVisa),
		assert:          assertNoError,
	},
	{
		name:            "CardScheme wrong scheme",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("5100000000000000"),
		constraint:      it.IsCardScheme(validate.CardSchemeVisa),
		assert:          assertHasOneViolation(validation.ErrInvalidCardScheme, message.InvalidCardScheme),
	},
	{
		name:            "CardScheme nonnumeric",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("invalid"),
		constraint:      it.IsCardScheme(validate.CardSchemeVisa),
		assert:          assertHasOneViolation(validation.ErrInvalidCardScheme, message.InvalidCardScheme),
	},
	{
		name:            "CardScheme multiple",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("5100000000000000"),
		constraint:      it.IsCardScheme(validate.CardSchemeVisa, validate.CardSchemeMastercard),
		assert:          assertNoError,
	},
	{
		name:            "CardScheme disabled",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("bad"),
		constraint:      it.IsCardScheme(validate.CardSchemeVisa).When(false),
		assert:          assertNoError,
	},
	{
		name:            "CardScheme enabled",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("bad"),
		constraint:      it.IsCardScheme(validate.CardSchemeVisa).When(true),
		assert:          assertHasOneViolation(validation.ErrInvalidCardScheme, message.InvalidCardScheme),
	},
	{
		name:            "CardScheme groups",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("bad"),
		constraint:      it.IsCardScheme(validate.CardSchemeVisa).WhenGroups(testGroup),
		assert:          assertNoError,
	},
	{
		name:            "CardScheme missing schemes",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsCardScheme(),
		assert:          assertError("validate by CardSchemeConstraint: at least one known card scheme is required"),
	},
	{
		name:            "CardScheme unknown schemes",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue(""),
		constraint:      it.IsCardScheme("unknown"),
		assert:          assertError("validate by CardSchemeConstraint: at least one known card scheme is required"),
	},
	{
		name:            "CardScheme custom",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("bad"),
		constraint:      it.IsCardScheme(validate.CardSchemeVisa).WithError(ErrCustom).WithMessage("{{ value }} {{ custom }}", validation.TemplateParameter{Key: "{{ custom }}", Value: "parameter"}),
		assert:          assertHasOneViolation(ErrCustom, "bad parameter"),
	},
}
