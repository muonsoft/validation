package internal

import (
	"errors"
	"testing"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/code"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/message"
	"github.com/muonsoft/validation/validationtest"
	"github.com/stretchr/testify/assert"
)

const (
	customMessage = "Custom message."
	customPath    = "properties[0].value"
)

type ValidateTestCase struct {
	name            string
	isApplicableFor func(valueType string) bool
	intValue        *int
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
		intValue:        intValue(0),
		floatValue:      floatValue(0),
		stringValue:     stringValue(""),
		options:         []validation.Option{it.IsNotBlank()},
		assert:          assertHasOneViolation(code.NotBlank, message.NotBlank, ""),
	},
	{
		name:            "IsNotBlank violation on empty value when condition is true",
		isApplicableFor: anyValueType,
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
		intValue:        intValue(1),
		floatValue:      floatValue(0.1),
		stringValue:     stringValue("a"),
		options:         []validation.Option{it.IsBlank()},
		assert:          assertHasOneViolation(code.Blank, message.Blank, ""),
	},
	{
		name:            "IsBlank violation on value when condition is true",
		isApplicableFor: anyValueType,
		intValue:        intValue(1),
		floatValue:      floatValue(0.1),
		stringValue:     stringValue("a"),
		options:         []validation.Option{it.IsBlank().When(true)},
		assert:          assertHasOneViolation(code.Blank, message.Blank, ""),
	},
	{
		name:            "IsBlank violation on value with custom path",
		isApplicableFor: anyValueType,
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
		intValue:        intValue(0),
		floatValue:      floatValue(0.0),
		stringValue:     stringValue(""),
		options:         []validation.Option{it.IsBlank()},
		assert:          assertNoError,
	},
	{
		name:            "IsBlank passes on value when condition is false",
		isApplicableFor: anyValueType,
		intValue:        intValue(1),
		floatValue:      floatValue(0.1),
		stringValue:     stringValue("a"),
		options:         []validation.Option{it.IsBlank().When(false)},
		assert:          assertNoError,
	},
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

func TestValidateInt(t *testing.T) {
	for _, test := range validateTestCases {
		t.Run(test.name, func(t *testing.T) {
			err := validation.ValidateInt(test.intValue, test.options...)

			if test.isApplicableFor("int") {
				test.assert(t, err)
			} else {
				assertIsInapplicableConstraintError(t, err, "int")
			}
		})
	}
}

func TestValidateFloat(t *testing.T) {
	for _, test := range validateTestCases {
		t.Run(test.name, func(t *testing.T) {
			err := validation.ValidateFloat(test.floatValue, test.options...)

			if test.isApplicableFor("float") {
				test.assert(t, err)
			} else {
				assertIsInapplicableConstraintError(t, err, "float")
			}
		})
	}
}

func anyValueType(valueType string) bool {
	return true
}

func assertHasOneViolation(code, message, path string) func(t *testing.T, err error) {
	return func(t *testing.T, err error) {
		validationtest.AssertIsViolationList(t, err, func(t *testing.T, violations validation.ViolationList) bool {
			if assert.Len(t, violations, 1) {
				return assert.Equal(t, code, violations[0].GetCode()) &&
					assert.Equal(t, message, violations[0].GetMessage()) &&
					assert.Equal(t, path, violations[0].GetPropertyPath().Format())
			}

			return false
		})
	}
}

func assertNoError(t *testing.T, err error) {
	assert.NoError(t, err)
}

func assertIsInapplicableConstraintError(t *testing.T, err error, valueType string) {
	var inapplicableConstraint *validation.ErrInapplicableConstraint

	if !errors.As(err, &inapplicableConstraint) {
		t.Errorf("failed asserting that error is ErrInapplicableConstraint")
		return
	}

	assert.Equal(t, valueType, inapplicableConstraint.Type)
}

func intValue(i int) *int {
	return &i
}

func floatValue(f float64) *float64 {
	return &f
}

func stringValue(s string) *string {
	return &s
}
