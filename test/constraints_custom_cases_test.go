package test

import (
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/message"
)

func validString(value string) bool {
	return true
}

func invalidString(value string) bool {
	return false
}

var customStringConstraintTestCases = []ConstraintValidationTestCase{
	{
		name:            "StringFuncConstraint passes on nil",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      validation.OfStringBy(invalidString),
		assert:          assertNoError,
	},
	{
		name:            "StringFuncConstraint passes on empty value",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      validation.OfStringBy(invalidString),
		stringValue:     stringValue(""),
		assert:          assertNoError,
	},
	{
		name:            "StringFuncConstraint passes on valid value",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      validation.OfStringBy(validString),
		stringValue:     stringValue("foo"),
		assert:          assertNoError,
	},
	{
		name:            "StringFuncConstraint violation on invalid value",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      validation.OfStringBy(invalidString),
		stringValue:     stringValue("foo"),
		assert:          assertHasOneViolation(validation.ErrNotValid, message.NotValid),
	},
	{
		name:            "StringFuncConstraint violation with given error and message",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      validation.OfStringBy(invalidString).WithError(ErrCustom).WithMessage("message"),
		stringValue:     stringValue("foo"),
		assert:          assertHasOneViolation(ErrCustom, "message"),
	},
	{
		name:            "StringFuncConstraint violation with custom message",
		isApplicableFor: specificValueTypes(stringType),
		constraint: validation.
			OfStringBy(invalidString).
			WithError(ErrCustom).
			WithMessage(
				`Unexpected value "{{ value }}" for {{ custom }}.`,
				validation.TemplateParameter{Key: "{{ custom }}", Value: "parameter"},
			),
		stringValue: stringValue("foo"),
		assert:      assertHasOneViolation(ErrCustom, `Unexpected value "foo" for parameter.`),
	},
	{
		name:            "StringFuncConstraint passes when condition is false",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      validation.OfStringBy(invalidString).When(false),
		stringValue:     stringValue("foo"),
		assert:          assertNoError,
	},
	{
		name:            "StringFuncConstraint passes when groups not match",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      validation.OfStringBy(invalidString).WhenGroups(testGroup),
		stringValue:     stringValue("foo"),
		assert:          assertNoError,
	},
	{
		name:            "StringFuncConstraint violation when condition is true",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      validation.OfStringBy(invalidString).When(true),
		stringValue:     stringValue("foo"),
		assert:          assertHasOneViolation(validation.ErrNotValid, message.NotValid),
	},
}
