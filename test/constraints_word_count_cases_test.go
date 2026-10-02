package test

import (
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
)

var wordCountConstraintTestCases = []ConstraintValidationTestCase{
	{
		name:            "WordCount nil",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.HasMinWordCount(1),
		assert:          assertNoError,
	},
	{
		name:            "WordCount empty",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue(""),
		constraint:      it.HasMinWordCount(1),
		assert:          assertNoError,
	},
	{
		name:            "WordCount whitespace",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue(" \t\n"),
		constraint:      it.HasMinWordCount(1),
		assert:          assertHasOneViolation(validation.ErrTooFewWords, "This value should contain at least 1 word."),
	},
	{
		name:            "WordCount min boundary",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("one two"),
		constraint:      it.HasWordCountBetween(2, 3),
		assert:          assertNoError,
	},
	{
		name:            "WordCount max boundary",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("one two three"),
		constraint:      it.HasWordCountBetween(2, 3),
		assert:          assertNoError,
	},
	{
		name:            "WordCount below min",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("one"),
		constraint:      it.HasWordCountBetween(2, 3),
		assert:          assertHasOneViolation(validation.ErrTooFewWords, "This value should contain at least 2 words."),
	},
	{
		name:            "WordCount above max",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("one two three four"),
		constraint:      it.HasWordCountBetween(2, 3),
		assert:          assertHasOneViolation(validation.ErrTooManyWords, "This value should contain at most 3 words."),
	},
	{
		name:            "WordCount zero max",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("one"),
		constraint:      it.HasMaxWordCount(0),
		assert:          assertHasOneViolation(validation.ErrTooManyWords, "This value should contain at most 0 words."),
	},
	{
		name:            "WordCount disabled",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("one"),
		constraint:      it.HasMinWordCount(2).When(false),
		assert:          assertNoError,
	},
	{
		name:            "WordCount enabled",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("one two"),
		constraint:      it.HasMaxWordCount(1).When(true),
		assert:          assertHasOneViolation(validation.ErrTooManyWords, "This value should contain at most 1 word."),
	},
	{
		name:            "WordCount ignored groups",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("one"),
		constraint:      it.HasMinWordCount(2).WhenGroups(testGroup),
		assert:          assertNoError,
	},
	{
		name:            "WordCount custom min",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("one"),
		constraint:      it.HasMinWordCount(2).WithMinError(ErrMin).WithMinMessage("{{ count }}/{{ limit }} {{ value }} {{ custom }}", validation.TemplateParameter{Key: "{{ custom }}", Value: "custom"}),
		assert:          assertHasOneViolation(ErrMin, `1/2 "one" custom`),
	},
	{
		name:            "WordCount custom max",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("one two"),
		constraint:      it.HasMaxWordCount(1).WithMaxError(ErrMax).WithMaxMessage("{{ count }}/{{ limit }} {{ value }} {{ custom }}", validation.TemplateParameter{Key: "{{ custom }}", Value: "custom"}),
		assert:          assertHasOneViolation(ErrMax, `2/1 "one two" custom`),
	},
}
