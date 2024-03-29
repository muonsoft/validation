package test

import (
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/message"
)

var choiceConstraintTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsOneOf error on empty list",
		isApplicableFor: specificValueTypes(stringType, comparableType),
		constraint:      it.IsOneOf[string](),
		assert:          assertError(`validate by ChoiceConstraint: empty list of choices`),
	},
	{
		name:            "IsOneOf passes on nil",
		isApplicableFor: specificValueTypes(stringType, comparableType),
		constraint:      it.IsOneOf("expected"),
		assert:          assertNoError,
	},
	{
		name:            "IsOneOf violation on empty string",
		isApplicableFor: specificValueTypes(stringType, comparableType),
		stringValue:     stringValue(""),
		constraint:      it.IsOneOf("expected"),
		assert:          assertNoError,
	},
	{
		name:            "IsOneOf passes on expected string",
		isApplicableFor: specificValueTypes(stringType, comparableType),
		stringValue:     stringValue("expected"),
		constraint:      it.IsOneOf("expected"),
		assert:          assertNoError,
	},
	{
		name:            "IsOneOf violation on missing value",
		isApplicableFor: specificValueTypes(stringType, comparableType),
		stringValue:     stringValue("not-expected"),
		constraint:      it.IsOneOf("expected"),
		assert:          assertHasOneViolation(validation.ErrNoSuchChoice, message.NoSuchChoice),
	},
	{
		name:            "IsOneOf violation on missing value with custom message",
		isApplicableFor: specificValueTypes(stringType, comparableType),
		stringValue:     stringValue("unexpected"),
		constraint: it.IsOneOf("alpha", "beta", "gamma").
			WithError(ErrCustom).
			WithMessage(
				`Unexpected value "{{ value }}" at {{ custom }}, expected values are: {{ choices }}.`,
				validation.TemplateParameter{Key: "{{ custom }}", Value: "parameter"},
			),
		assert: assertHasOneViolation(
			ErrCustom,
			`Unexpected value "unexpected" at parameter, expected values are: alpha, beta, gamma.`,
		),
	},
	{
		name:            "IsOneOf passes when condition is false",
		isApplicableFor: specificValueTypes(stringType, comparableType),
		stringValue:     stringValue("not-expected"),
		constraint:      it.IsOneOf("expected").When(false),
		assert:          assertNoError,
	},
	{
		name:            "IsOneOf passes when groups not match",
		isApplicableFor: specificValueTypes(stringType, comparableType),
		stringValue:     stringValue("not-expected"),
		constraint:      it.IsOneOf("expected").WhenGroups(testGroup),
		assert:          assertNoError,
	},
	{
		name:            "IsOneOf violation on missing value when condition is true",
		isApplicableFor: specificValueTypes(stringType, comparableType),
		stringValue:     stringValue("not-expected"),
		constraint:      it.IsOneOf("expected").When(true),
		assert:          assertHasOneViolation(validation.ErrNoSuchChoice, message.NoSuchChoice),
	},
	{
		name:            "IsOneOf violation on missing integer value",
		isApplicableFor: specificValueTypes(intType),
		intValue:        intValue(4),
		constraint:      it.IsOneOf(1, 2, 3, 5),
		assert:          assertHasOneViolation(validation.ErrNoSuchChoice, message.NoSuchChoice),
	},
	{
		name:            "IsOneOf violation on blank integer",
		isApplicableFor: specificValueTypes(intType),
		intValue:        intValue(0),
		constraint:      it.IsOneOf(1, 2),
		assert:          assertNoError,
	},
	{
		name:            "IsOneOf passes on blank integer when blank is allowed",
		isApplicableFor: specificValueTypes(intType),
		intValue:        intValue(0),
		constraint:      it.IsOneOf(1, 2).WithoutBlank(),
		assert:          assertHasOneViolation(validation.ErrNoSuchChoice, message.NoSuchChoice),
	},
}
