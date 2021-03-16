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
		name:            "HasMinLength violation with custom min message",
		isApplicableFor: specificValueTypes(stringType),
		options: []validation.Option{
			it.HasMinLength(2).MinMessage(customMessage),
		},
		stringValue: stringValue("a"),
		assert:      assertHasOneViolation(code.LengthTooFew, customMessage, ""),
	},
	{
		name:            "HasMinLength violation with custom max message",
		isApplicableFor: specificValueTypes(stringType),
		options: []validation.Option{
			it.HasMaxLength(2).MaxMessage(customMessage),
		},
		stringValue: stringValue("aaa"),
		assert:      assertHasOneViolation(code.LengthTooMany, customMessage, ""),
	},
	{
		name:            "HasMinLength violation with custom exact message",
		isApplicableFor: specificValueTypes(stringType),
		options: []validation.Option{
			it.HasExactLength(2).ExactMessage(customMessage),
		},
		stringValue: stringValue("aaa"),
		assert:      assertHasOneViolation(code.LengthExact, customMessage, ""),
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
	{
		name:            "HasLengthBetween passes on expected string",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("aaa"),
		options:         []validation.Option{it.HasLengthBetween(1, 5)},
		assert:          assertNoError,
	},
	{
		name:            "HasExactLength passes on expected string",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("aaa"),
		options:         []validation.Option{it.HasExactLength(3)},
		assert:          assertNoError,
	},
}

var regexConstraintTestCases = []ConstraintValidationTestCase{
	{
		name:            "Matches passes on nil",
		isApplicableFor: specificValueTypes(stringType),
		options:         []validation.Option{it.Matches("^[a-z]+$")},
		assert:          assertNoError,
	},
	{
		name:            "Matches passes on empty value",
		isApplicableFor: specificValueTypes(stringType),
		options:         []validation.Option{it.Matches("^[a-z]+$")},
		stringValue:     stringValue(""),
		assert:          assertNoError,
	},
	{
		name:            "Matches violation ignored when condition false",
		isApplicableFor: specificValueTypes(stringType),
		options:         []validation.Option{it.Matches("^[a-z]+$").When(false)},
		stringValue:     stringValue("1"),
		assert:          assertNoError,
	},
	{
		name:            "Matches violation when condition true",
		isApplicableFor: specificValueTypes(stringType),
		options:         []validation.Option{it.Matches("^[a-z]+$").When(true)},
		stringValue:     stringValue("1"),
		assert: assertHasOneViolation(
			code.NotMatches,
			"This value is not valid.",
			"",
		),
	},
	{
		name:            "Matches violation with custom message",
		isApplicableFor: specificValueTypes(stringType),
		options: []validation.Option{
			it.Matches("^[a-z]+$").Message(customMessage),
		},
		stringValue: stringValue("1"),
		assert:      assertHasOneViolation(code.NotMatches, customMessage, ""),
	},
	{
		name:            "Matches passes on expected string",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("a"),
		options:         []validation.Option{it.Matches("^[a-z]+$")},
		assert:          assertNoError,
	},
}
