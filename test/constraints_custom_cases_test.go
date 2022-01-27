package test

import (
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/code"
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
		name:            "CustomStringConstraint passes on nil",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      validation.NewCustomStringConstraint(invalidString),
		assert:          assertNoError,
	},
	{
		name:            "CustomStringConstraint passes on empty value",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      validation.NewCustomStringConstraint(invalidString),
		stringValue:     stringValue(""),
		assert:          assertNoError,
	},
	{
		name:            "CustomStringConstraint passes on valid value",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      validation.NewCustomStringConstraint(validString),
		stringValue:     stringValue("foo"),
		assert:          assertNoError,
	},
	{
		name:            "CustomStringConstraint violation on invalid value",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      validation.NewCustomStringConstraint(invalidString),
		stringValue:     stringValue("foo"),
		assert:          assertHasOneViolation(code.NotValid, message.Templates[code.NotValid]),
	},
	{
		name:            "CustomStringConstraint violation with given code and message",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      validation.NewCustomStringConstraint(invalidString, "name", "code", "message"),
		stringValue:     stringValue("foo"),
		assert:          assertHasOneViolation("code", "message"),
	},
	{
		name:            "CustomStringConstraint violation with custom message",
		isApplicableFor: specificValueTypes(stringType),
		constraint: validation.
			NewCustomStringConstraint(invalidString).
			Code(customCode).
			Message(
				`Unexpected value "{{ value }}" for {{ custom }}.`,
				validation.TemplateParameter{Key: "{{ custom }}", Value: "parameter"},
			),
		stringValue: stringValue("foo"),
		assert:      assertHasOneViolation(customCode, `Unexpected value "foo" for parameter.`),
	},
	{
		name:            "CustomStringConstraint passes when condition is false",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      validation.NewCustomStringConstraint(invalidString).When(false),
		stringValue:     stringValue("foo"),
		assert:          assertNoError,
	},
	{
		name:            "CustomStringConstraint passes when groups not match",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      validation.NewCustomStringConstraint(invalidString).WhenGroups(testGroup),
		stringValue:     stringValue("foo"),
		assert:          assertNoError,
	},
	{
		name:            "CustomStringConstraint violation when condition is true",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      validation.NewCustomStringConstraint(invalidString).When(true),
		stringValue:     stringValue("foo"),
		assert:          assertHasOneViolation(code.NotValid, message.Templates[code.NotValid]),
	},
}
