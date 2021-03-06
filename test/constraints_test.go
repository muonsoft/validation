package test

import (
	"testing"
	"time"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/code"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/message"
	"github.com/muonsoft/validation/validator"
)

const (
	customMessage = "Custom message."
	customPath    = "properties[0].value"

	// Value types.
	boolType      = "bool"
	intType       = "int"
	floatType     = "float"
	stringType    = "string"
	iterableType  = "iterable"
	countableType = "countable"
	timeType      = "time"
)

type ConstraintValidationTestCase struct {
	name            string
	isApplicableFor func(valueType string) bool
	boolValue       *bool
	intValue        *int64
	floatValue      *float64
	stringValue     *string
	timeValue       *time.Time
	sliceValue      []string
	mapValue        map[string]string
	options         []validation.Option
	assert          func(t *testing.T, err error)
}

var validateTestCases = mergeTestCases(
	isNotBlankConstraintTestCases,
	isBlankConstraintTestCases,
	isNotNilConstraintTestCases,
	countConstraintTestCases,
)

func TestValidateBool(t *testing.T) {
	for _, test := range validateTestCases {
		if !test.isApplicableFor(boolType) {
			continue
		}

		t.Run(test.name, func(t *testing.T) {
			err := validator.ValidateBool(test.boolValue, test.options...)

			test.assert(t, err)
		})
	}
}

func TestValidateNumber_AsInt(t *testing.T) {
	for _, test := range validateTestCases {
		t.Run(test.name, func(t *testing.T) {
			err := validator.ValidateNumber(test.intValue, test.options...)

			if test.isApplicableFor(intType) {
				test.assert(t, err)
			} else {
				assertIsInapplicableConstraintError(t, err, "number")
			}
		})
	}
}

func TestValidateNumber_AsFloat(t *testing.T) {
	for _, test := range validateTestCases {
		t.Run(test.name, func(t *testing.T) {
			err := validator.ValidateNumber(test.floatValue, test.options...)

			if test.isApplicableFor(floatType) {
				test.assert(t, err)
			} else {
				assertIsInapplicableConstraintError(t, err, "number")
			}
		})
	}
}

func TestValidateString(t *testing.T) {
	for _, test := range validateTestCases {
		t.Run(test.name, func(t *testing.T) {
			err := validator.ValidateString(test.stringValue, test.options...)

			if test.isApplicableFor(stringType) {
				test.assert(t, err)
			} else {
				assertIsInapplicableConstraintError(t, err, stringType)
			}
		})
	}
}

func TestValidateIterable_AsSlice(t *testing.T) {
	for _, test := range validateTestCases {
		t.Run(test.name, func(t *testing.T) {
			err := validator.ValidateIterable(test.sliceValue, test.options...)

			if test.isApplicableFor(iterableType) {
				test.assert(t, err)
			} else {
				assertIsInapplicableConstraintError(t, err, iterableType)
			}
		})
	}
}

func TestValidateIterable_AsMap(t *testing.T) {
	for _, test := range validateTestCases {
		t.Run(test.name, func(t *testing.T) {
			err := validator.ValidateIterable(test.mapValue, test.options...)

			if test.isApplicableFor(iterableType) {
				test.assert(t, err)
			} else {
				assertIsInapplicableConstraintError(t, err, iterableType)
			}
		})
	}
}

func TestValidateCountable(t *testing.T) {
	for _, test := range validateTestCases {
		if !test.isApplicableFor(countableType) {
			continue
		}

		t.Run(test.name, func(t *testing.T) {
			err := validator.ValidateCountable(len(test.sliceValue), test.options...)

			if test.isApplicableFor(countableType) {
				test.assert(t, err)
			} else {
				assertIsInapplicableConstraintError(t, err, countableType)
			}
		})
	}
}

func TestValidateTime(t *testing.T) {
	for _, test := range validateTestCases {
		if !test.isApplicableFor(timeType) {
			continue
		}

		t.Run(test.name, func(t *testing.T) {
			err := validator.ValidateTime(test.timeValue, test.options...)

			if test.isApplicableFor(timeType) {
				test.assert(t, err)
			} else {
				assertIsInapplicableConstraintError(t, err, timeType)
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
		{"not blank", it.IsNotBlank(), assertHasOneViolation(code.NotBlank, message.NotBlank, "")},
		{"not blank when true", it.IsNotBlank().When(true), assertHasOneViolation(code.NotBlank, message.NotBlank, "")},
		{"not blank when false", it.IsNotBlank().When(false), assertNoError},
		{"not blank when nil allowed", it.IsNotBlank().AllowNil(), assertNoError},
		{"blank", it.IsBlank(), assertNoError},
		{"blank when true", it.IsBlank().When(true), assertNoError},
		{"blank when false", it.IsBlank().When(false), assertNoError},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var v *bool

			err := validator.ValidateValue(v, test.nilConstraint)

			test.assert(t, err)
		})
	}
}

func anyValueType(valueType string) bool {
	return true
}

func specificValueTypes(types ...string) func(valueType string) bool {
	return func(valueType string) bool {
		for _, t := range types {
			if valueType == t {
				return true
			}
		}

		return false
	}
}

func exceptValueTypes(types ...string) func(valueType string) bool {
	return func(valueType string) bool {
		for _, t := range types {
			if valueType == t {
				return false
			}
		}

		return true
	}
}

func mergeTestCases(testCases ...[]ConstraintValidationTestCase) []ConstraintValidationTestCase {
	merged := make([]ConstraintValidationTestCase, 0)

	for _, testCase := range testCases {
		merged = append(merged, testCase...)
	}

	return merged
}
