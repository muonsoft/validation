package test

import (
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/message"
)

var weekConstraintTestCases = []ConstraintValidationTestCase{
	{
		name:            "Week nil",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsWeek(),
		assert:          assertNoError,
	},
	{
		name:            "Week empty",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue(""),
		constraint:      it.IsWeek(),
		assert:          assertNoError,
	},
	{
		name:            "Week valid",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("2020-W53"),
		constraint:      it.IsWeek(),
		assert:          assertNoError,
	},
	{
		name:            "Week format",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("2020-W00"),
		constraint:      it.IsWeek(),
		assert:          assertHasOneViolation(validation.ErrInvalidWeekFormat, message.InvalidWeekFormat),
	},
	{
		name:            "Week number",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("2021-W53"),
		constraint:      it.IsWeek(),
		assert:          assertHasOneViolation(validation.ErrInvalidWeekNumber, "Week 2021-W53 does not exist in its ISO year."),
	},
	{
		name:            "Week min",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("2020-W53"),
		constraint:      it.IsWeek().WithMin("2021-W01"),
		assert:          assertHasOneViolation(validation.ErrWeekTooEarly, "This value should be on or after week 2021-W01."),
	},
	{
		name:            "Week max",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("2021-W01"),
		constraint:      it.IsWeek().WithMax("2020-W53"),
		assert:          assertHasOneViolation(validation.ErrWeekTooLate, "This value should be on or before week 2020-W53."),
	},
	{
		name:            "Week bounds",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("2020-W53"),
		constraint:      it.IsWeek().WithMin("2020-W53").WithMax("2020-W53"),
		assert:          assertNoError,
	},
	{
		name:            "Week disabled",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("invalid"),
		constraint:      it.IsWeek().When(false),
		assert:          assertNoError,
	},
	{
		name:            "Week enabled",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("invalid"),
		constraint:      it.IsWeek().When(true),
		assert:          assertHasOneViolation(validation.ErrInvalidWeekFormat, message.InvalidWeekFormat),
	},
	{
		name:            "Week groups",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("invalid"),
		constraint:      it.IsWeek().WhenGroups(testGroup),
		assert:          assertNoError,
	},
	{
		name:            "Week invalid bounds",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsWeek().WithMin("invalid"),
		assert:          assertError("validate by WeekConstraint: bounds must be valid ISO weeks and minimum must not exceed maximum"),
	},
	{
		name:            "Week custom Format",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("invalid"),
		constraint:      it.IsWeek().WithFormatError(ErrCustom).WithFormatMessage(customMessage, validation.TemplateParameter{Key: "{{ custom }}", Value: "parameter"}),
		assert:          assertHasOneViolation(ErrCustom, renderedCustomMessage),
	},
	{
		name:            "Week custom WeekNumber",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("2021-W53"),
		constraint:      it.IsWeek().WithWeekNumberError(ErrCustom).WithWeekNumberMessage(customMessage, validation.TemplateParameter{Key: "{{ custom }}", Value: "parameter"}),
		assert:          assertHasOneViolation(ErrCustom, renderedCustomMessage),
	},
	{
		name:            "Week custom Min",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("2020-W53"),
		constraint:      it.IsWeek().WithMin("2021-W01").WithMinError(ErrCustom).WithMinMessage(customMessage, validation.TemplateParameter{Key: "{{ custom }}", Value: "parameter"}),
		assert:          assertHasOneViolation(ErrCustom, renderedCustomMessage),
	},
	{
		name:            "Week custom Max",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("2021-W01"),
		constraint:      it.IsWeek().WithMax("2020-W53").WithMaxError(ErrCustom).WithMaxMessage(customMessage, validation.TemplateParameter{Key: "{{ custom }}", Value: "parameter"}),
		assert:          assertHasOneViolation(ErrCustom, renderedCustomMessage),
	},
}
