package test

import (
	"time"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/code"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/message"
)

var isNotBlankNumberConstraintTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsNotBlankNumber violation on nil",
		isApplicableFor: specificValueTypes(intType),
		constraint:      it.IsNotBlankNumber[int](),
		assert:          assertHasOneViolation(code.NotBlank, message.Templates[code.NotBlank]),
	},
	{
		name:            "IsNotBlankNumber violation on empty int value",
		isApplicableFor: specificValueTypes(intType),
		intValue:        intValue(0),
		constraint:      it.IsNotBlankNumber[int](),
		assert:          assertHasOneViolation(code.NotBlank, message.Templates[code.NotBlank]),
	},
	{
		name:            "IsNotBlankNumber violation on empty float value",
		isApplicableFor: specificValueTypes(floatType),
		floatValue:      floatValue(0),
		constraint:      it.IsNotBlankNumber[float64](),
		assert:          assertHasOneViolation(code.NotBlank, message.Templates[code.NotBlank]),
	},
	{
		name:            "IsNotBlankNumber violation on empty value when condition is true",
		isApplicableFor: specificValueTypes(intType),
		boolValue:       boolValue(false),
		intValue:        intValue(0),
		floatValue:      floatValue(0),
		stringValue:     stringValue(""),
		stringsValue:    []string{},
		sliceValue:      []string{},
		mapValue:        map[string]string{},
		constraint:      it.IsNotBlankNumber[int]().When(true),
		assert:          assertHasOneViolation(code.NotBlank, message.Templates[code.NotBlank]),
	},
	{
		name:            "IsNotBlankNumber violation on nil with custom message",
		isApplicableFor: specificValueTypes(intType),
		constraint: it.IsNotBlankNumber[int]().
			Code(customCode).
			Message(
				customMessage,
				validation.TemplateParameter{Key: "{{ custom }}", Value: "parameter"},
			),
		assert: assertHasOneViolation(customCode, renderedCustomMessage),
	},
	{
		name:            "IsNotBlankNumber passes on value",
		isApplicableFor: specificValueTypes(intType),
		boolValue:       boolValue(true),
		intValue:        intValue(1),
		floatValue:      floatValue(0.1),
		stringValue:     stringValue("a"),
		stringsValue:    []string{""},
		timeValue:       timeValue(time.Now()),
		sliceValue:      []string{"a"},
		mapValue:        map[string]string{"a": "a"},
		constraint:      it.IsNotBlankNumber[int](),
		assert:          assertNoError,
	},
	{
		name:            "IsNotBlankNumber passes on nil when allowed",
		isApplicableFor: specificValueTypes(intType),
		constraint:      it.IsNotBlankNumber[int]().AllowNil(),
		assert:          assertNoError,
	},
	{
		name:            "IsNotBlankNumber passes on nil when condition is false",
		isApplicableFor: specificValueTypes(intType),
		constraint:      it.IsNotBlankNumber[int]().When(false),
		assert:          assertNoError,
	},
	{
		name:            "IsNotBlankNumber passes on nil when groups not match",
		isApplicableFor: specificValueTypes(intType),
		constraint:      it.IsNotBlankNumber[int]().WhenGroups(testGroup),
		assert:          assertNoError,
	},
}

var isBlankNumberConstraintTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsBlankNumber violation on value",
		isApplicableFor: specificValueTypes(intType),
		boolValue:       boolValue(true),
		intValue:        intValue(1),
		floatValue:      floatValue(0.1),
		stringValue:     stringValue("a"),
		stringsValue:    []string{""},
		timeValue:       timeValue(time.Now()),
		sliceValue:      []string{"a"},
		mapValue:        map[string]string{"a": "a"},
		constraint:      it.IsBlankNumber[int](),
		assert:          assertHasOneViolation(code.Blank, message.Templates[code.Blank]),
	},
	{
		name:            "IsBlankNumber violation on value when condition is true",
		isApplicableFor: specificValueTypes(intType),
		boolValue:       boolValue(true),
		intValue:        intValue(1),
		floatValue:      floatValue(0.1),
		stringValue:     stringValue("a"),
		stringsValue:    []string{""},
		timeValue:       timeValue(time.Now()),
		sliceValue:      []string{"a"},
		mapValue:        map[string]string{"a": "a"},
		constraint:      it.IsBlankNumber[int]().When(true),
		assert:          assertHasOneViolation(code.Blank, message.Templates[code.Blank]),
	},
	{
		name:            "IsBlankNumber violation on value with custom message",
		isApplicableFor: specificValueTypes(intType),
		boolValue:       boolValue(true),
		intValue:        intValue(1),
		floatValue:      floatValue(0.1),
		stringValue:     stringValue("a"),
		stringsValue:    []string{""},
		timeValue:       timeValue(time.Now()),
		sliceValue:      []string{"a"},
		mapValue:        map[string]string{"a": "a"},
		constraint: it.IsBlankNumber[int]().
			Code(customCode).
			Message(
				customMessage,
				validation.TemplateParameter{Key: "{{ custom }}", Value: "parameter"},
			),
		assert: assertHasOneViolation(customCode, renderedCustomMessage),
	},
	{
		name:            "IsBlankNumber passes on nil",
		isApplicableFor: specificValueTypes(intType),
		constraint:      it.IsBlankNumber[int](),
		assert:          assertNoError,
	},
	{
		name:            "IsBlankNumber passes on empty value",
		isApplicableFor: specificValueTypes(intType),
		boolValue:       boolValue(false),
		intValue:        intValue(0),
		floatValue:      floatValue(0.0),
		stringValue:     stringValue(""),
		timeValue:       timeValue(time.Time{}),
		stringsValue:    []string{},
		sliceValue:      []string{},
		mapValue:        map[string]string{},
		constraint:      it.IsBlankNumber[int](),
		assert:          assertNoError,
	},
	{
		name:            "IsBlankNumber passes on value when condition is false",
		isApplicableFor: specificValueTypes(intType),
		boolValue:       boolValue(true),
		intValue:        intValue(1),
		floatValue:      floatValue(0.1),
		stringValue:     stringValue("a"),
		timeValue:       timeValue(time.Now()),
		stringsValue:    []string{""},
		sliceValue:      []string{"a"},
		mapValue:        map[string]string{"a": "a"},
		constraint:      it.IsBlankNumber[int]().When(false),
		assert:          assertNoError,
	},
	{
		name:            "IsBlankNumber passes on value when groups not match",
		isApplicableFor: specificValueTypes(intType),
		boolValue:       boolValue(true),
		intValue:        intValue(1),
		floatValue:      floatValue(0.1),
		stringValue:     stringValue("a"),
		timeValue:       timeValue(time.Now()),
		stringsValue:    []string{""},
		sliceValue:      []string{"a"},
		mapValue:        map[string]string{"a": "a"},
		constraint:      it.IsBlankNumber[int]().WhenGroups(testGroup),
		assert:          assertNoError,
	},
}
