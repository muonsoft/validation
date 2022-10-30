package test

import (
	"regexp"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/message"
)

var lengthConstraintTestCases = []ConstraintValidationTestCase{
	{
		name:            "HasMinLength passes on nil",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.HasMinLength(1),
		assert:          assertNoError,
	},
	{
		name:            "HasMinLength passes on empty value",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.HasMinLength(1),
		stringValue:     stringValue(""),
		assert:          assertNoError,
	},
	{
		name:            "HasMinLength violation ignored when condition false",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.HasMinLength(2).When(false),
		stringValue:     stringValue("a"),
		assert:          assertNoError,
	},
	{
		name:            "HasMinLength violation ignored when groups not match",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.HasMinLength(2).WhenGroups(testGroup),
		stringValue:     stringValue("a"),
		assert:          assertNoError,
	},
	{
		name:            "HasMinLength violation when condition true",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.HasMinLength(2).When(true),
		stringValue:     stringValue("a"),
		assert: assertHasOneViolation(
			validation.ErrTooShort,
			"This value is too short. It should have 2 characters or more.",
		),
	},
	{
		name:            "HasMinLength violation with custom min message",
		isApplicableFor: specificValueTypes(stringType),
		constraint: it.HasMinLength(2).
			WithMinError(ErrMin).
			WithMinMessage(
				"Unexpected length {{ length }} at {{ custom }} value {{ value }}, should not be less than {{ limit }}.",
				validation.TemplateParameter{Key: "{{ custom }}", Value: "parameter"},
			),
		stringValue: stringValue("a"),
		assert: assertHasOneViolation(
			ErrMin,
			`Unexpected length 1 at parameter value "a", should not be less than 2.`,
		),
	},
	{
		name:            "HasMinLength violation with custom max message",
		isApplicableFor: specificValueTypes(stringType),
		constraint: it.HasMaxLength(2).
			WithMaxError(ErrMax).
			WithMaxMessage(
				"Unexpected length {{ length }} at {{ custom }} value {{ value }}, should not be greater than {{ limit }}.",
				validation.TemplateParameter{Key: "{{ custom }}", Value: "parameter"},
			),
		stringValue: stringValue("aaa"),
		assert: assertHasOneViolation(
			ErrMax,
			`Unexpected length 3 at parameter value "aaa", should not be greater than 2.`,
		),
	},
	{
		name:            "HasMinLength violation with custom exact message",
		isApplicableFor: specificValueTypes(stringType),
		constraint: it.HasExactLength(2).
			WithExactError(ErrExact).
			WithExactMessage(
				"Unexpected length {{ length }} at {{ custom }} value {{ value }}, should be exactly {{ limit }}.",
				validation.TemplateParameter{Key: "{{ custom }}", Value: "parameter"},
			),
		stringValue: stringValue("aaa"),
		assert: assertHasOneViolation(
			ErrExact,
			`Unexpected length 3 at parameter value "aaa", should be exactly 2.`,
		),
	},
	{
		name:            "HasMinLength passes on equal length",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("aa"),
		constraint:      it.HasMinLength(2),
		assert:          assertNoError,
	},
	{
		name:            "HasMaxLength violation on max",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("aaa"),
		constraint:      it.HasMaxLength(2),
		assert: assertHasOneViolation(
			validation.ErrTooLong,
			"This value is too long. It should have 2 characters or less.",
		),
	},
	{
		name:            "HasLengthBetween passes on expected string",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("aaa"),
		constraint:      it.HasLengthBetween(1, 5),
		assert:          assertNoError,
	},
	{
		name:            "HasExactLength passes on expected string",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("aaa"),
		constraint:      it.HasExactLength(3),
		assert:          assertNoError,
	},
}

var regexConstraintTestCases = []ConstraintValidationTestCase{
	{
		name:            "Matches error on nil regex",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.Matches(nil),
		assert:          assertError(`validate by RegexpConstraint: nil regex`),
	},
	{
		name:            "Matches passes on nil",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.Matches(regexp.MustCompile("^[a-z]+$")),
		assert:          assertNoError,
	},
	{
		name:            "Matches passes on empty value",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.Matches(regexp.MustCompile("^[a-z]+$")),
		stringValue:     stringValue(""),
		assert:          assertNoError,
	},
	{
		name:            "Matches violation ignored when condition false",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.Matches(regexp.MustCompile("^[a-z]+$")).When(false),
		stringValue:     stringValue("1"),
		assert:          assertNoError,
	},
	{
		name:            "Matches violation ignored when groups not match",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.Matches(regexp.MustCompile("^[a-z]+$")).WhenGroups(testGroup),
		stringValue:     stringValue("1"),
		assert:          assertNoError,
	},
	{
		name:            "Matches violation when condition true",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.Matches(regexp.MustCompile("^[a-z]+$")).When(true),
		stringValue:     stringValue("1"),
		assert:          assertHasOneViolation(validation.ErrNotValid, message.NotValid),
	},
	{
		name:            "Matches violation with custom message",
		isApplicableFor: specificValueTypes(stringType),
		constraint: it.Matches(regexp.MustCompile("^[a-z]+$")).
			WithError(ErrCustom).
			WithMessage(
				`Unexpected value "{{ value }}" at {{ custom }}.`,
				validation.TemplateParameter{Key: "{{ custom }}", Value: "parameter"},
			),
		stringValue: stringValue("1"),
		assert:      assertHasOneViolation(ErrCustom, `Unexpected value "1" at parameter.`),
	},
	{
		name:            "Matches passes on expected string",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("a"),
		constraint:      it.Matches(regexp.MustCompile("^[a-z]+$")),
		assert:          assertNoError,
	},
	{
		name:            "DoesNotMatch violation on expected string",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("a"),
		constraint:      it.DoesNotMatch(regexp.MustCompile("^[a-z]+$")),
		assert:          assertHasOneViolation(validation.ErrNotValid, message.NotValid),
	},
	{
		name:            "DoesNotMatch passes on expected string",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("1"),
		constraint:      it.DoesNotMatch(regexp.MustCompile("^[a-z]+$")),
		assert:          assertNoError,
	},
}

var jsonConstraintTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsJSON passes on valid JSON",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsJSON(),
		stringValue:     stringValue(`{"valid": true}`),
		assert:          assertNoError,
	},
	{
		name:            "IsJSON violation on invalid JSON",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsJSON(),
		stringValue:     stringValue(`"invalid": true`),
		assert:          assertHasOneViolation(validation.ErrInvalidJSON, message.InvalidJSON),
	},
}

var numericConstraintTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsInteger passes on valid integer",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsInteger(),
		stringValue:     stringValue("123"),
		assert:          assertNoError,
	},
	{
		name:            "IsInteger violation on invalid integer",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsInteger(),
		stringValue:     stringValue("foo"),
		assert:          assertHasOneViolation(validation.ErrNotInteger, message.NotInteger),
	},
	{
		name:            "IsNumeric passes on valid number",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsNumeric(),
		stringValue:     stringValue("123.123"),
		assert:          assertNoError,
	},
	{
		name:            "IsNumeric violation on invalid number",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsNumeric(),
		stringValue:     stringValue("foo.bar"),
		assert:          assertHasOneViolation(validation.ErrNotNumeric, message.NotNumeric),
	},
}
