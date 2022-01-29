package test

import (
	"context"
	"testing"
	"time"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/code"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/message"
	"github.com/stretchr/testify/assert"
)

const (
	customCode            = "customCode"
	customMessage         = "Custom message at {{ custom }}."
	renderedCustomMessage = "Custom message at parameter."
	customPath            = "properties[0].value"

	// Value types.
	boolType      = "bool"
	intType       = "int"
	floatType     = "float"
	stringType    = "string"
	stringsType   = "strings"
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
	stringsValue    []string
	timeValue       *time.Time
	sliceValue      []string
	mapValue        map[string]string
	constraint      validation.Constraint
	assert          func(t *testing.T, err error)
}

var validateTestCases = mergeTestCases(
	barcodeConstraintsTestCases,
	isNotBlankConstraintTestCases,
	isBlankConstraintTestCases,
	isNotNilConstraintTestCases,
	isNilConstraintTestCases,
	isTrueConstraintTestCases,
	isFalseConstraintTestCases,
	lengthConstraintTestCases,
	regexConstraintTestCases,
	countConstraintTestCases,
	choiceConstraintTestCases,
	numberComparisonTestCases,
	stringComparisonTestCases,
	hasUniqueValuesTestCases,
	customStringConstraintTestCases,
	timeComparisonTestCases,
	rangeComparisonTestCases,
	isBetweenTimeTestCases,
	urlConstraintTestCases,
	emailConstraintTestCases,
	ipConstraintTestCases,
	hostnameConstraintTestCases,
	jsonConstraintTestCases,
	numericConstraintTestCases,
)

func TestValidateBool(t *testing.T) {
	for _, test := range validateTestCases {
		if !test.isApplicableFor(boolType) {
			continue
		}

		t.Run(test.name, func(t *testing.T) {
			err := newValidator(t).Validate(context.Background(), validation.NilBool(test.boolValue, test.constraint))

			test.assert(t, err)
		})
	}
}

func TestValidateNumber_AsInt(t *testing.T) {
	for _, test := range validateTestCases {
		if !test.isApplicableFor(intType) {
			continue
		}

		t.Run(test.name, func(t *testing.T) {
			err := newValidator(t).Validate(context.Background(), validation.Number(test.intValue, test.constraint))

			test.assert(t, err)
		})
	}
}

func TestValidateNumber_AsFloat(t *testing.T) {
	for _, test := range validateTestCases {
		if !test.isApplicableFor(floatType) {
			continue
		}

		t.Run(test.name, func(t *testing.T) {
			err := newValidator(t).Validate(context.Background(), validation.Number(test.floatValue, test.constraint))

			test.assert(t, err)
		})
	}
}

func TestValidateString(t *testing.T) {
	for _, test := range validateTestCases {
		if !test.isApplicableFor(stringType) {
			continue
		}

		t.Run(test.name, func(t *testing.T) {
			err := newValidator(t).Validate(context.Background(), validation.NilString(test.stringValue, test.constraint))

			test.assert(t, err)
		})
	}
}

func TestValidateStrings(t *testing.T) {
	for _, test := range validateTestCases {
		if !test.isApplicableFor(stringsType) {
			continue
		}

		t.Run(test.name, func(t *testing.T) {
			err := newValidator(t).Validate(context.Background(), validation.Strings(test.stringsValue, test.constraint))

			test.assert(t, err)
		})
	}
}

func TestValidateIterable_AsSlice(t *testing.T) {
	for _, test := range validateTestCases {
		if !test.isApplicableFor(iterableType) {
			continue
		}

		t.Run(test.name, func(t *testing.T) {
			err := newValidator(t).Validate(context.Background(), validation.Iterable(test.sliceValue, test.constraint))

			test.assert(t, err)
		})
	}
}

func TestValidateIterable_AsMap(t *testing.T) {
	for _, test := range validateTestCases {
		if !test.isApplicableFor(iterableType) {
			continue
		}

		t.Run(test.name, func(t *testing.T) {
			err := newValidator(t).Validate(context.Background(), validation.Iterable(test.mapValue, test.constraint))

			test.assert(t, err)
		})
	}
}

func TestValidateCountable(t *testing.T) {
	for _, test := range validateTestCases {
		if !test.isApplicableFor(countableType) {
			continue
		}

		t.Run(test.name, func(t *testing.T) {
			err := newValidator(t).Validate(context.Background(), validation.Countable(len(test.sliceValue), test.constraint))

			test.assert(t, err)
		})
	}
}

func TestValidateTime(t *testing.T) {
	for _, test := range validateTestCases {
		if !test.isApplicableFor(timeType) {
			continue
		}

		t.Run(test.name, func(t *testing.T) {
			err := newValidator(t).Validate(context.Background(), validation.NilTime(test.timeValue, test.constraint))

			test.assert(t, err)
		})
	}
}

func TestValidateNil(t *testing.T) {
	tests := []struct {
		name          string
		nilConstraint validation.NilConstraint
		assert        func(t *testing.T, err error)
	}{
		{"not blank", it.IsNotBlank(), assertHasOneViolation(code.NotBlank, message.Templates[code.NotBlank])},
		{"not blank when true", it.IsNotBlank().When(true), assertHasOneViolation(code.NotBlank, message.Templates[code.NotBlank])},
		{"not blank when false", it.IsNotBlank().When(false), assertNoError},
		{"not blank when nil allowed", it.IsNotBlank().AllowNil(), assertNoError},
		{"blank", it.IsBlank(), assertNoError},
		{"blank when true", it.IsBlank().When(true), assertNoError},
		{"blank when false", it.IsBlank().When(false), assertNoError},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var v *bool

			err := newValidator(t).Validate(context.Background(), validation.Value(v, test.nilConstraint))

			test.assert(t, err)
		})
	}
}

func TestCustomStringConstraint_Name_WhenNoNameIsSet_ExpectDefaultName(t *testing.T) {
	constraint := validation.NewCustomStringConstraint(validString)

	name := constraint.Name()

	assert.Equal(t, "CustomStringConstraint", name)
}

func TestCustomStringConstraint_Name_WhenNameIsSet_ExpectGivenName(t *testing.T) {
	constraint := validation.NewCustomStringConstraint(validString, "CustomName")

	name := constraint.Name()

	assert.Equal(t, "CustomName", name)
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
