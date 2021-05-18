package test

import (
	"github.com/muonsoft/validation/code"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/message"
)

var choiceConstraintTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsOneOfStrings error on empty list",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsOneOfStrings(),
		assert:          assertError(`failed to set up constraint "ChoiceConstraint": empty list of choices`),
	},
	{
		name:            "IsOneOfStrings passes on nil",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsOneOfStrings("expected"),
		assert:          assertNoError,
	},
	{
		name:            "IsOneOfStrings passes on empty string",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue(""),
		constraint:      it.IsOneOfStrings("expected"),
		assert:          assertNoError,
	},
	{
		name:            "IsOneOfStrings passes on expected string",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("expected"),
		constraint:      it.IsOneOfStrings("expected"),
		assert:          assertNoError,
	},
	{
		name:            "IsOneOfStrings violation on missing value",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("not-expected"),
		constraint:      it.IsOneOfStrings("expected"),
		assert:          assertHasOneViolation(code.NoSuchChoice, message.NoSuchChoice),
	},
	{
		name:            "IsOneOfStrings violation on missing value with custom message",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("unexpected"),
		constraint: it.IsOneOfStrings("alpha", "beta", "gamma").
			Message(`Unexpected value "{{ value }}", expected values are: {{ choices }}.`),
		assert: assertHasOneViolation(
			code.NoSuchChoice,
			`Unexpected value "unexpected", expected values are: alpha, beta, gamma.`,
		),
	},
	{
		name:            "IsOneOfStrings passes when condition is false",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("not-expected"),
		constraint:      it.IsOneOfStrings("expected").When(false),
		assert:          assertNoError,
	},
	{
		name:            "IsOneOfStrings violation on missing value when condition is true",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("not-expected"),
		constraint:      it.IsOneOfStrings("expected").When(true),
		assert:          assertHasOneViolation(code.NoSuchChoice, message.NoSuchChoice),
	},
}
