package test

import (
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
)

var dateTimeConstraintTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsDateTime passes on nil",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     nil,
		constraint:      it.IsDateTime(),
		assert:          assertNoError,
	},
	{
		name:            "IsDateTime passes on empty value",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue(""),
		constraint:      it.IsDateTime(),
		assert:          assertNoError,
	},
	{
		name:            "IsDateTime violation on invalid value",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("invalid"),
		constraint:      it.IsDateTime(),
		assert:          assertHasOneViolation(validation.ErrInvalidDateTime, "This value is not a valid datetime."),
	},
	{
		name:            "IsDateTime passes on valid value",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("2022-07-12T12:34:56+00:00"),
		constraint:      it.IsDateTime(),
		assert:          assertNoError,
	},
	{
		name:            "IsDateTime passes when condition is false",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("invalid"),
		constraint:      it.IsDateTime().When(false),
		assert:          assertNoError,
	},
	{
		name:            "IsDateTime violation when condition is true",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("invalid"),
		constraint:      it.IsDateTime().When(true),
		assert:          assertHasOneViolation(validation.ErrInvalidDateTime, "This value is not a valid datetime."),
	},
	{
		name:            "IsDateTime passes when groups not match",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("invalid"),
		constraint:      it.IsDateTime().WhenGroups(testGroup),
		assert:          assertNoError,
	},
	{
		name:            "IsDateTime violation when groups match",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("invalid"),
		constraint:      it.IsDateTime().WhenGroups(validation.DefaultGroup),
		assert:          assertHasOneViolation(validation.ErrInvalidDateTime, "This value is not a valid datetime."),
	},
	{
		name:            "IsDateTime violation with custom message",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("invalid"),
		constraint: it.IsDateTime().
			WithError(ErrCustom).
			WithMessage(
				`Invalid date time at {{ custom }} value {{ value }} with layout {{ layout }}.`,
				validation.TemplateParameter{Key: "{{ custom }}", Value: "parameter"},
			),
		assert: assertHasOneViolation(
			ErrCustom,
			`Invalid date time at parameter value invalid with layout 2006-01-02T15:04:05Z07:00.`,
		),
	},
	{
		name:            "IsDateTime passes with custom layout",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("2022-07-12 12:34:56"),
		constraint:      it.IsDateTime().WithLayout("2006-01-02 15:04:05"),
		assert:          assertNoError,
	},
	{
		name:            "IsDate passes on valid value",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("2022-07-12"),
		constraint:      it.IsDate(),
		assert:          assertNoError,
	},
	{
		name:            "IsDate violation on invalid value",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("invalid"),
		constraint:      it.IsDate(),
		assert:          assertHasOneViolation(validation.ErrInvalidDate, "This value is not a valid date."),
	},
	{
		name:            "IsTime passes on valid value",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("12:34:56"),
		constraint:      it.IsTime(),
		assert:          assertNoError,
	},
	{
		name:            "IsTime violation on invalid value",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("invalid"),
		constraint:      it.IsTime(),
		assert:          assertHasOneViolation(validation.ErrInvalidTime, "This value is not a valid time."),
	},
}
