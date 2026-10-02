package test

import (
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/message"
	"github.com/muonsoft/validation/validate"
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
	{
		name:            "IsTimezone passes on empty value",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue(""),
		constraint:      it.IsTimezone(),
		assert:          assertNoError,
	},
	{
		name:            "IsTimezone passes on UTC",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("UTC"),
		constraint:      it.IsTimezone(),
		assert:          assertNoError,
	},
	{
		name:            "IsTimezone passes on valid IANA identifier",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("Europe/Berlin"),
		constraint:      it.IsTimezone(),
		assert:          assertNoError,
	},
	{
		name:            "IsTimezone violation on unknown identifier",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("Invalid/Zone"),
		constraint:      it.IsTimezone(),
		assert:          assertHasOneViolation(validation.ErrInvalidTimezone, message.InvalidTimezone),
	},
	{
		name:            "IsTimezone violation on Local",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("Local"),
		constraint:      it.IsTimezone(),
		assert:          assertHasOneViolation(validation.ErrInvalidTimezone, message.InvalidTimezone),
	},
	{
		name:            "IsTimezone violation on abbreviation",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("EST"),
		constraint:      it.IsTimezone(),
		assert:          assertHasOneViolation(validation.ErrInvalidTimezone, message.InvalidTimezone),
	},
	{
		name:            "IsTimezone passes with matching zone",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("Europe/Paris"),
		constraint:      it.IsTimezone().WithZone(validate.TimezoneZoneEurope),
		assert:          assertNoError,
	},
	{
		name:            "IsTimezone violation with non-matching zone",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("America/New_York"),
		constraint:      it.IsTimezone().WithZone(validate.TimezoneZoneEurope),
		assert:          assertHasOneViolation(validation.ErrInvalidTimezone, message.InvalidTimezone),
	},
	{
		name:            "IsTimezone violation with custom error and message",
		isApplicableFor: specificValueTypes(stringType),
		constraint: it.IsTimezone().
			WithError(ErrCustom).
			WithMessage(
				`Invalid timezone "{{ value }}" for {{ custom }}.`,
				validation.TemplateParameter{Key: "{{ custom }}", Value: "parameter"},
			),
		stringValue: stringValue("EST"),
		assert:      assertHasOneViolation(ErrCustom, `Invalid timezone "EST" for parameter.`),
	},
	{
		name:            "IsTimezone passes when condition is false",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsTimezone().When(false),
		stringValue:     stringValue("EST"),
		assert:          assertNoError,
	},
	{
		name:            "IsTimezone violation when condition is true",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsTimezone().When(true),
		stringValue:     stringValue("EST"),
		assert:          assertHasOneViolation(validation.ErrInvalidTimezone, message.InvalidTimezone),
	},
	{
		name:            "IsTimezone passes when groups not match",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsTimezone().WhenGroups(testGroup),
		stringValue:     stringValue("EST"),
		assert:          assertNoError,
	},
}
