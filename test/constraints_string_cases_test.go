package test

import (
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/code"
	"github.com/muonsoft/validation/it"
)

var lengthConstraintTestCases = []ConstraintValidationTestCase{
	{
		name:            "HasMinLength passes on nil",
		isApplicableFor: specificValueTypes(stringType),
		options:         []validation.Option{it.HasMinLength(1)},
		assert:          assertNoError,
	},
	{
		name:            "HasMinLength passes on empty value",
		isApplicableFor: specificValueTypes(stringType),
		options:         []validation.Option{it.HasMinLength(1)},
		stringValue:     stringValue(""),
		assert:          assertNoError,
	},
	{
		name:            "HasMinLength violation ignored when condition false",
		isApplicableFor: specificValueTypes(stringType),
		options:         []validation.Option{it.HasMinLength(2).When(false)},
		stringValue:     stringValue("a"),
		assert:          assertNoError,
	},
	{
		name:            "HasMinLength violation when condition true",
		isApplicableFor: specificValueTypes(stringType),
		options:         []validation.Option{it.HasMinLength(2).When(true)},
		stringValue:     stringValue("a"),
		assert: assertHasOneViolation(
			code.LengthTooFew,
			"This value is too short. It should have 2 characters or more.",
			"",
		),
	},
	{
		name:            "HasMinLength violation with custom property path",
		isApplicableFor: specificValueTypes(stringType),
		options: []validation.Option{
			validation.PropertyName("properties"),
			validation.ArrayIndex(0),
			validation.PropertyName("value"),
			it.HasMinLength(2),
		},
		stringValue: stringValue("a"),
		assert: assertHasOneViolation(
			code.LengthTooFew,
			"This value is too short. It should have 2 characters or more.",
			customPath,
		),
	},
	{
		name:            "HasMinLength violation with custom message",
		isApplicableFor: specificValueTypes(stringType),
		options: []validation.Option{
			it.HasMinLength(2).MinMessage(customMessage),
		},
		stringValue: stringValue("a"),
		assert:      assertHasOneViolation(code.LengthTooFew, customMessage, ""),
	},
	{
		name:            "HasMinLength passes on equal length",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("aa"),
		options:         []validation.Option{it.HasMinLength(2)},
		assert:          assertNoError,
	},
	{
		name:            "HasMaxLength violation on max",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("aaa"),
		options:         []validation.Option{it.HasMaxLength(2)},
		assert: assertHasOneViolation(
			code.LengthTooMany,
			"This value is too long. It should have 2 characters or less.",
			"",
		),
	},
}
