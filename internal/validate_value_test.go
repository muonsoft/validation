package internal

import (
	"testing"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/code"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/message"
)

const (
	customMessage = "Custom message."
	customPath    = "properties[0].value"
)

type ValidateTestCase struct {
	name            string
	isApplicableFor func(valueType string) bool
	boolValue       *bool
	intValue        *int64
	floatValue      *float64
	stringValue     *string
	options         []validation.Option
	assert          func(t *testing.T, err error)
}

var validateTestCases = []ValidateTestCase{
	// IsNotBlank
	{
		name:            "IsNotBlank violation on nil",
		isApplicableFor: anyValueType,
		options:         []validation.Option{it.IsNotBlank()},
		assert:          assertHasOneViolation(code.NotBlank, message.NotBlank, ""),
	},
	{
		name:            "IsNotBlank violation on empty value",
		isApplicableFor: anyValueType,
		boolValue:       boolValue(false),
		intValue:        intValue(0),
		floatValue:      floatValue(0),
		stringValue:     stringValue(""),
		options:         []validation.Option{it.IsNotBlank()},
		assert:          assertHasOneViolation(code.NotBlank, message.NotBlank, ""),
	},
	{
		name:            "IsNotBlank violation on empty value when condition is true",
		isApplicableFor: anyValueType,
		boolValue:       boolValue(false),
		intValue:        intValue(0),
		floatValue:      floatValue(0),
		stringValue:     stringValue(""),
		options:         []validation.Option{it.IsNotBlank().When(true)},
		assert:          assertHasOneViolation(code.NotBlank, message.NotBlank, ""),
	},
	{
		name:            "IsNotBlank violation on nil with custom path",
		isApplicableFor: anyValueType,
		options: []validation.Option{
			validation.PropertyName("properties"),
			validation.ArrayIndex(0),
			validation.PropertyName("value"),
			it.IsNotBlank(),
		},
		assert: assertHasOneViolation(code.NotBlank, message.NotBlank, customPath),
	},
	{
		name:            "IsNotBlank violation on nil with custom message",
		isApplicableFor: anyValueType,
		options:         []validation.Option{it.IsNotBlank().Message(customMessage)},
		assert:          assertHasOneViolation(code.NotBlank, customMessage, ""),
	},
	{
		name:            "IsNotBlank passes on value",
		isApplicableFor: anyValueType,
		boolValue:       boolValue(true),
		intValue:        intValue(1),
		floatValue:      floatValue(0.1),
		stringValue:     stringValue("a"),
		options:         []validation.Option{it.IsNotBlank()},
		assert:          assertNoError,
	},
	{
		name:            "IsNotBlank passes on nil when allowed",
		isApplicableFor: anyValueType,
		options:         []validation.Option{it.IsNotBlank().AllowNil()},
		assert:          assertNoError,
	},
	{
		name:            "IsNotBlank passes on nil when condition is false",
		isApplicableFor: anyValueType,
		options:         []validation.Option{it.IsNotBlank().When(false)},
		assert:          assertNoError,
	},

	// IsBlank
	{
		name:            "IsBlank violation on value",
		isApplicableFor: anyValueType,
		boolValue:       boolValue(true),
		intValue:        intValue(1),
		floatValue:      floatValue(0.1),
		stringValue:     stringValue("a"),
		options:         []validation.Option{it.IsBlank()},
		assert:          assertHasOneViolation(code.Blank, message.Blank, ""),
	},
	{
		name:            "IsBlank violation on value when condition is true",
		isApplicableFor: anyValueType,
		boolValue:       boolValue(true),
		intValue:        intValue(1),
		floatValue:      floatValue(0.1),
		stringValue:     stringValue("a"),
		options:         []validation.Option{it.IsBlank().When(true)},
		assert:          assertHasOneViolation(code.Blank, message.Blank, ""),
	},
	{
		name:            "IsBlank violation on value with custom path",
		isApplicableFor: anyValueType,
		boolValue:       boolValue(true),
		intValue:        intValue(1),
		floatValue:      floatValue(0.1),
		stringValue:     stringValue("a"),
		options: []validation.Option{
			validation.PropertyName("properties"),
			validation.ArrayIndex(0),
			validation.PropertyName("value"),
			it.IsBlank(),
		},
		assert: assertHasOneViolation(code.Blank, message.Blank, customPath),
	},
	{
		name:            "IsBlank violation on value with custom message",
		isApplicableFor: anyValueType,
		boolValue:       boolValue(true),
		intValue:        intValue(1),
		floatValue:      floatValue(0.1),
		stringValue:     stringValue("a"),
		options:         []validation.Option{it.IsBlank().Message(customMessage)},
		assert:          assertHasOneViolation(code.Blank, customMessage, ""),
	},
	{
		name:            "IsBlank passes on nil",
		isApplicableFor: anyValueType,
		options:         []validation.Option{it.IsBlank()},
		assert:          assertNoError,
	},
	{
		name:            "IsBlank passes on empty value",
		isApplicableFor: anyValueType,
		boolValue:       boolValue(false),
		intValue:        intValue(0),
		floatValue:      floatValue(0.0),
		stringValue:     stringValue(""),
		options:         []validation.Option{it.IsBlank()},
		assert:          assertNoError,
	},
	{
		name:            "IsBlank passes on value when condition is false",
		isApplicableFor: anyValueType,
		boolValue:       boolValue(true),
		intValue:        intValue(1),
		floatValue:      floatValue(0.1),
		stringValue:     stringValue("a"),
		options:         []validation.Option{it.IsBlank().When(false)},
		assert:          assertNoError,
	},
}

func TestValidateBool(t *testing.T) {
	tests := make([]ValidateTestCase, 0)
	for _, test := range validateTestCases {
		if test.isApplicableFor("bool") {
			tests = append(tests, test)
		}
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validation.ValidateBool(test.boolValue, test.options...)

			test.assert(t, err)
		})
	}
}

func TestValidateNumber_AsInt(t *testing.T) {
	for _, test := range validateTestCases {
		t.Run(test.name, func(t *testing.T) {
			err := validation.ValidateNumber(test.intValue, test.options...)

			if test.isApplicableFor("int") {
				test.assert(t, err)
			} else {
				assertIsInapplicableConstraintError(t, err, "int")
			}
		})
	}
}

func TestValidateNumber_AsFloat(t *testing.T) {
	for _, test := range validateTestCases {
		t.Run(test.name, func(t *testing.T) {
			err := validation.ValidateNumber(test.floatValue, test.options...)

			if test.isApplicableFor("float") {
				test.assert(t, err)
			} else {
				assertIsInapplicableConstraintError(t, err, "float")
			}
		})
	}
}

func TestValidateString(t *testing.T) {
	for _, test := range validateTestCases {
		t.Run(test.name, func(t *testing.T) {
			err := validation.ValidateString(test.stringValue, test.options...)

			if test.isApplicableFor("string") {
				test.assert(t, err)
			} else {
				assertIsInapplicableConstraintError(t, err, "string")
			}
		})
	}
}

func TestValidateNil(t *testing.T) {
	tests := []struct {
		name          string
		nilConstraint validation.NilConstraint
		assert        func(t *testing.T, err error)
	}{
		{"not blank", it.IsNotBlank(), assertIsViolation(code.NotBlank, message.NotBlank, "")},
		{"not blank when true", it.IsNotBlank().When(true), assertIsViolation(code.NotBlank, message.NotBlank, "")},
		{"not blank when false", it.IsNotBlank().When(false), assertNoError},
		{"not blank when nil allowed", it.IsNotBlank().AllowNil(), assertNoError},
		{"blank", it.IsBlank(), assertNoError},
		{"blank when true", it.IsBlank().When(true), assertNoError},
		{"blank when false", it.IsBlank().When(false), assertNoError},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.nilConstraint.ValidateNil(validation.Options{
				NewViolation: validation.NewViolation,
			})

			test.assert(t, err)
		})
	}
}

func anyValueType(valueType string) bool {
	return true
}
