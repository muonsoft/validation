package test

import (
	"time"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/code"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/message"
)

var isNotBlankConstraintTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsNotBlank violation on nil",
		isApplicableFor: anyValueType,
		constraint:      it.IsNotBlank(),
		assert:          assertHasOneViolation(code.NotBlank, message.NotBlank),
	},
	{
		name:            "IsNotBlank violation on empty value",
		isApplicableFor: anyValueType,
		boolValue:       boolValue(false),
		intValue:        intValue(0),
		floatValue:      floatValue(0),
		stringValue:     stringValue(""),
		stringsValue:    []string{},
		sliceValue:      []string{},
		mapValue:        map[string]string{},
		constraint:      it.IsNotBlank(),
		assert:          assertHasOneViolation(code.NotBlank, message.NotBlank),
	},
	{
		name:            "IsNotBlank violation on empty value when condition is true",
		isApplicableFor: anyValueType,
		boolValue:       boolValue(false),
		intValue:        intValue(0),
		floatValue:      floatValue(0),
		stringValue:     stringValue(""),
		stringsValue:    []string{},
		sliceValue:      []string{},
		mapValue:        map[string]string{},
		constraint:      it.IsNotBlank().When(true),
		assert:          assertHasOneViolation(code.NotBlank, message.NotBlank),
	},
	{
		name:            "IsNotBlank violation on nil with custom message",
		isApplicableFor: anyValueType,
		constraint: it.IsNotBlank().
			Code(customCode).
			Message(
				customMessage,
				validation.TemplateParameter{Key: "{{ custom }}", Value: "parameter"},
			),
		assert: assertHasOneViolation(customCode, renderedCustomMessage),
	},
	{
		name:            "IsNotBlank passes on value",
		isApplicableFor: anyValueType,
		boolValue:       boolValue(true),
		intValue:        intValue(1),
		floatValue:      floatValue(0.1),
		stringValue:     stringValue("a"),
		stringsValue:    []string{""},
		timeValue:       timeValue(time.Now()),
		sliceValue:      []string{"a"},
		mapValue:        map[string]string{"a": "a"},
		constraint:      it.IsNotBlank(),
		assert:          assertNoError,
	},
	{
		name:            "IsNotBlank passes on nil when allowed",
		isApplicableFor: exceptValueTypes("countable"),
		constraint:      it.IsNotBlank().AllowNil(),
		assert:          assertNoError,
	},
	{
		name:            "IsNotBlank passes on nil when condition is false",
		isApplicableFor: exceptValueTypes("countable"),
		constraint:      it.IsNotBlank().When(false),
		assert:          assertNoError,
	},
}

var isBlankConstraintTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsBlank violation on value",
		isApplicableFor: anyValueType,
		boolValue:       boolValue(true),
		intValue:        intValue(1),
		floatValue:      floatValue(0.1),
		stringValue:     stringValue("a"),
		stringsValue:    []string{""},
		timeValue:       timeValue(time.Now()),
		sliceValue:      []string{"a"},
		mapValue:        map[string]string{"a": "a"},
		constraint:      it.IsBlank(),
		assert:          assertHasOneViolation(code.Blank, message.Blank),
	},
	{
		name:            "IsBlank violation on value when condition is true",
		isApplicableFor: anyValueType,
		boolValue:       boolValue(true),
		intValue:        intValue(1),
		floatValue:      floatValue(0.1),
		stringValue:     stringValue("a"),
		stringsValue:    []string{""},
		timeValue:       timeValue(time.Now()),
		sliceValue:      []string{"a"},
		mapValue:        map[string]string{"a": "a"},
		constraint:      it.IsBlank().When(true),
		assert:          assertHasOneViolation(code.Blank, message.Blank),
	},
	{
		name:            "IsBlank violation on value with custom message",
		isApplicableFor: anyValueType,
		boolValue:       boolValue(true),
		intValue:        intValue(1),
		floatValue:      floatValue(0.1),
		stringValue:     stringValue("a"),
		stringsValue:    []string{""},
		timeValue:       timeValue(time.Now()),
		sliceValue:      []string{"a"},
		mapValue:        map[string]string{"a": "a"},
		constraint: it.IsBlank().
			Code(customCode).
			Message(
				customMessage,
				validation.TemplateParameter{Key: "{{ custom }}", Value: "parameter"},
			),
		assert: assertHasOneViolation(customCode, renderedCustomMessage),
	},
	{
		name:            "IsBlank passes on nil",
		isApplicableFor: anyValueType,
		constraint:      it.IsBlank(),
		assert:          assertNoError,
	},
	{
		name:            "IsBlank passes on empty value",
		isApplicableFor: anyValueType,
		boolValue:       boolValue(false),
		intValue:        intValue(0),
		floatValue:      floatValue(0.0),
		stringValue:     stringValue(""),
		timeValue:       timeValue(time.Time{}),
		stringsValue:    []string{},
		sliceValue:      []string{},
		mapValue:        map[string]string{},
		constraint:      it.IsBlank(),
		assert:          assertNoError,
	},
	{
		name:            "IsBlank passes on value when condition is false",
		isApplicableFor: anyValueType,
		boolValue:       boolValue(true),
		intValue:        intValue(1),
		floatValue:      floatValue(0.1),
		stringValue:     stringValue("a"),
		timeValue:       timeValue(time.Now()),
		stringsValue:    []string{""},
		sliceValue:      []string{"a"},
		mapValue:        map[string]string{"a": "a"},
		constraint:      it.IsBlank().When(false),
		assert:          assertNoError,
	},
}

var isNotNilConstraintTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsNotNil violation on nil",
		isApplicableFor: exceptValueTypes(countableType),
		constraint:      it.IsNotNil(),
		assert:          assertHasOneViolation(code.NotNil, message.NotNil),
	},
	{
		name:            "IsNotNil passes on empty value",
		isApplicableFor: exceptValueTypes(countableType),
		boolValue:       boolValue(false),
		intValue:        intValue(0),
		floatValue:      floatValue(0),
		stringValue:     stringValue(""),
		stringsValue:    []string{},
		timeValue:       &time.Time{},
		sliceValue:      []string{},
		mapValue:        map[string]string{},
		constraint:      it.IsNotNil(),
		assert:          assertNoError,
	},
	{
		name:            "IsNotNil passes on empty value when condition is true",
		isApplicableFor: exceptValueTypes(countableType),
		boolValue:       boolValue(false),
		intValue:        intValue(0),
		floatValue:      floatValue(0),
		stringValue:     stringValue(""),
		stringsValue:    []string{},
		timeValue:       &time.Time{},
		sliceValue:      []string{},
		mapValue:        map[string]string{},
		constraint:      it.IsNotNil().When(true),
		assert:          assertNoError,
	},
	{
		name:            "IsNotNil violation on nil with custom message",
		isApplicableFor: exceptValueTypes(countableType),
		constraint: it.IsNotNil().
			Code(customCode).
			Message(
				customMessage,
				validation.TemplateParameter{Key: "{{ custom }}", Value: "parameter"},
			),
		assert: assertHasOneViolation(customCode, renderedCustomMessage),
	},
	{
		name:            "IsNotNil passes on value",
		isApplicableFor: exceptValueTypes(countableType),
		boolValue:       boolValue(true),
		intValue:        intValue(1),
		floatValue:      floatValue(0.1),
		stringValue:     stringValue("a"),
		stringsValue:    []string{},
		timeValue:       timeValue(time.Now()),
		sliceValue:      []string{},
		mapValue:        map[string]string{},
		constraint:      it.IsNotNil(),
		assert:          assertNoError,
	},
	{
		name:            "IsNotNil passes on nil when condition is false",
		isApplicableFor: exceptValueTypes(countableType),
		constraint:      it.IsNotNil().When(false),
		assert:          assertNoError,
	},
}

var isNilConstraintTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsNil passes on nil",
		isApplicableFor: exceptValueTypes(countableType),
		constraint:      it.IsNil(),
		assert:          assertNoError,
	},
	{
		name:            "IsNil violation on empty value",
		isApplicableFor: exceptValueTypes(countableType),
		boolValue:       boolValue(false),
		intValue:        intValue(0),
		floatValue:      floatValue(0),
		stringValue:     stringValue(""),
		stringsValue:    []string{},
		timeValue:       &time.Time{},
		sliceValue:      []string{},
		mapValue:        map[string]string{},
		constraint:      it.IsNil(),
		assert:          assertHasOneViolation(code.Nil, message.Nil),
	},
	{
		name:            "IsNil passes on nil when condition is true",
		isApplicableFor: exceptValueTypes(countableType),
		constraint:      it.IsNil().When(true),
		assert:          assertNoError,
	},
	{
		name:            "IsNil violation on empty value with custom message",
		isApplicableFor: exceptValueTypes(countableType),
		boolValue:       boolValue(false),
		intValue:        intValue(0),
		floatValue:      floatValue(0),
		stringValue:     stringValue(""),
		stringsValue:    []string{},
		timeValue:       &time.Time{},
		sliceValue:      []string{},
		mapValue:        map[string]string{},
		constraint: it.IsNil().
			Code(customCode).
			Message(
				customMessage,
				validation.TemplateParameter{Key: "{{ custom }}", Value: "parameter"},
			),
		assert: assertHasOneViolation(customCode, renderedCustomMessage),
	},
	{
		name:            "IsNil violation on value",
		isApplicableFor: exceptValueTypes(countableType),
		boolValue:       boolValue(true),
		intValue:        intValue(1),
		floatValue:      floatValue(0.1),
		stringValue:     stringValue("a"),
		stringsValue:    []string{},
		timeValue:       timeValue(time.Now()),
		sliceValue:      []string{},
		mapValue:        map[string]string{},
		constraint:      it.IsNil(),
		assert:          assertHasOneViolation(code.Nil, message.Nil),
	},
	{
		name:            "IsNil passes on empty value when condition is false",
		isApplicableFor: exceptValueTypes(countableType),
		boolValue:       boolValue(false),
		intValue:        intValue(0),
		floatValue:      floatValue(0),
		stringValue:     stringValue(""),
		stringsValue:    []string{},
		timeValue:       &time.Time{},
		sliceValue:      []string{},
		mapValue:        map[string]string{},
		constraint:      it.IsNil().When(false),
		assert:          assertNoError,
	},
}

var isTrueConstraintTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsTrue passes on nil",
		isApplicableFor: specificValueTypes(boolType),
		constraint:      it.IsTrue(),
		assert:          assertNoError,
	},
	{
		name:            "IsTrue violation on empty value",
		isApplicableFor: specificValueTypes(boolType),
		boolValue:       boolValue(false),
		constraint:      it.IsTrue(),
		assert:          assertHasOneViolation(code.True, message.True),
	},
	{
		name:            "IsTrue violation on empty value when condition is true",
		isApplicableFor: specificValueTypes(boolType),
		boolValue:       boolValue(false),
		constraint:      it.IsTrue().When(true),
		assert:          assertHasOneViolation(code.True, message.True),
	},
	{
		name:            "IsTrue violation on empty value with custom message",
		isApplicableFor: specificValueTypes(boolType),
		boolValue:       boolValue(false),
		constraint: it.IsTrue().
			Code(customCode).
			Message(
				customMessage,
				validation.TemplateParameter{Key: "{{ custom }}", Value: "parameter"},
			),
		assert: assertHasOneViolation(customCode, renderedCustomMessage),
	},
	{
		name:            "IsTrue passes on value",
		isApplicableFor: specificValueTypes(boolType),
		boolValue:       boolValue(true),
		constraint:      it.IsTrue(),
		assert:          assertNoError,
	},
	{
		name:            "IsTrue passes on empty value when condition is false",
		isApplicableFor: specificValueTypes(boolType),
		boolValue:       boolValue(false),
		constraint:      it.IsTrue().When(false),
		assert:          assertNoError,
	},
}

var isFalseConstraintTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsFalse passes on nil",
		isApplicableFor: specificValueTypes(boolType),
		constraint:      it.IsFalse(),
		assert:          assertNoError,
	},
	{
		name:            "IsFalse passes on empty value",
		isApplicableFor: specificValueTypes(boolType),
		boolValue:       boolValue(false),
		constraint:      it.IsFalse(),
		assert:          assertNoError,
	},
	{
		name:            "IsFalse violation on error value when condition is true",
		isApplicableFor: specificValueTypes(boolType),
		boolValue:       boolValue(true),
		constraint:      it.IsFalse().When(true),
		assert:          assertHasOneViolation(code.False, message.False),
	},
	{
		name:            "IsFalse violation on error value with custom message",
		isApplicableFor: specificValueTypes(boolType),
		boolValue:       boolValue(true),
		constraint:      it.IsFalse().Message(customMessage),
		assert:          assertHasOneViolation(code.False, customMessage),
	},
	{
		name:            "IsFalse passes on value",
		isApplicableFor: specificValueTypes(boolType),
		boolValue:       boolValue(false),
		constraint:      it.IsFalse(),
		assert:          assertNoError,
	},
	{
		name:            "IsFalse passes on error value when condition is false",
		isApplicableFor: specificValueTypes(boolType),
		boolValue:       boolValue(true),
		constraint:      it.IsFalse().When(false),
		assert:          assertNoError,
	},
}
