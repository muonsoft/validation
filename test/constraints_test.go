package test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/muonsoft/validation"
)

var (
	ErrCustom = errors.New("custom")
	ErrMin    = errors.New("min")
	ErrMax    = errors.New("max")
	ErrExact  = errors.New("exact")
)

const (
	customMessage         = "Custom message at {{ custom }}."
	renderedCustomMessage = "Custom message at parameter."
	customPath            = "properties[0].value"

	// Value types.
	nilType        = "nil"
	boolType       = "bool"
	intType        = "int"
	floatType      = "float"
	stringType     = "string"
	stringsType    = "strings"
	iterableType   = "iterable"
	countableType  = "countable"
	comparableType = "comparable"
	timeType       = "time"
)

type ConstraintValidationTestCase struct {
	name            string
	isApplicableFor func(valueType string) bool
	boolValue       *bool
	intValue        *int
	floatValue      *float64
	stringValue     *string
	stringsValue    []string
	timeValue       *time.Time
	sliceValue      []string
	mapValue        map[string]string
	constraint      interface{}
	assert          func(t *testing.T, err error)
}

var validateTestCases = mergeTestCases(
	barcodeConstraintsTestCases,
	choiceConstraintTestCases,
	comparableComparisonTestCases,
	countConstraintTestCases,
	customStringConstraintTestCases,
	dateTimeConstraintTestCases,
	emailConstraintTestCases,
	hasUniqueValuesTestCases,
	hostnameConstraintTestCases,
	identifierConstraintsTestCases,
	ipConstraintTestCases,
	isBetweenTimeTestCases,
	isBlankComparableConstraintTestCases,
	isBlankConstraintTestCases,
	isBlankNumberConstraintTestCases,
	isFalseConstraintTestCases,
	isNilComparableConstraintTestCases,
	isNilConstraintTestCases,
	isNilNumberConstraintTestCases,
	isNotBlankComparableConstraintTestCases,
	isNotBlankConstraintTestCases,
	isNotBlankNumberConstraintTestCases,
	isNotNilComparableConstraintTestCases,
	isNotNilConstraintTestCases,
	isNotNilNumberConstraintTestCases,
	isTrueConstraintTestCases,
	jsonConstraintTestCases,
	lengthConstraintTestCases,
	numberComparisonTestCases,
	numericConstraintTestCases,
	rangeComparisonTestCases,
	regexConstraintTestCases,
	timeComparisonTestCases,
	urlConstraintTestCases,
)

func TestValidateNil(t *testing.T) {
	for _, test := range validateTestCases {
		if !test.isApplicableFor(nilType) {
			continue
		}

		t.Run(test.name, func(t *testing.T) {
			err := newValidator(t).Validate(
				context.Background(),
				validation.Nil(test.stringValue == nil, test.constraint.(validation.NilConstraint)),
			)

			test.assert(t, err)
		})
	}
}

func TestValidateNilBool(t *testing.T) {
	for _, test := range validateTestCases {
		if !test.isApplicableFor(boolType) {
			continue
		}

		t.Run(test.name, func(t *testing.T) {
			err := newValidator(t).Validate(
				context.Background(),
				validation.NilBool(test.boolValue, test.constraint.(validation.BoolConstraint)),
			)

			test.assert(t, err)
		})
	}
}

func TestValidateNilNumber_AsInt(t *testing.T) {
	for _, test := range validateTestCases {
		if !test.isApplicableFor(intType) {
			continue
		}

		t.Run(test.name, func(t *testing.T) {
			err := newValidator(t).Validate(
				context.Background(),
				validation.NilNumber(test.intValue, test.constraint.(validation.NumberConstraint[int])),
			)

			test.assert(t, err)
		})
	}
}

func TestValidateNilNumber_AsFloat(t *testing.T) {
	for _, test := range validateTestCases {
		if !test.isApplicableFor(floatType) {
			continue
		}

		t.Run(test.name, func(t *testing.T) {
			err := newValidator(t).Validate(
				context.Background(),
				validation.NilNumber(test.floatValue, test.constraint.(validation.NumberConstraint[float64])),
			)

			test.assert(t, err)
		})
	}
}

func TestValidateNilString(t *testing.T) {
	for _, test := range validateTestCases {
		if !test.isApplicableFor(stringType) {
			continue
		}

		t.Run(test.name, func(t *testing.T) {
			err := newValidator(t).Validate(
				context.Background(),
				validation.NilString(test.stringValue, test.constraint.(validation.StringConstraint)),
			)

			test.assert(t, err)
		})
	}
}

func TestValidateNilComparable(t *testing.T) {
	for _, test := range validateTestCases {
		if !test.isApplicableFor(comparableType) {
			continue
		}

		t.Run(test.name, func(t *testing.T) {
			err := newValidator(t).Validate(
				context.Background(),
				validation.NilComparable[string](test.stringValue, test.constraint.(validation.ComparableConstraint[string])),
			)

			test.assert(t, err)
		})
	}
}

func TestValidateComparables(t *testing.T) {
	for _, test := range validateTestCases {
		if !test.isApplicableFor(stringsType) {
			continue
		}

		t.Run(test.name, func(t *testing.T) {
			err := newValidator(t).Validate(
				context.Background(),
				validation.Comparables(test.stringsValue, test.constraint.(validation.ComparablesConstraint[string])),
			)

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
			err := newValidator(t).Validate(
				context.Background(),
				validation.Countable(len(test.sliceValue), test.constraint.(validation.CountableConstraint)),
			)

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
			err := newValidator(t).Validate(
				context.Background(),
				validation.NilTime(test.timeValue, test.constraint.(validation.TimeConstraint)),
			)

			test.assert(t, err)
		})
	}
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

func mergeTestCases(testCases ...[]ConstraintValidationTestCase) []ConstraintValidationTestCase {
	merged := make([]ConstraintValidationTestCase, 0)

	for _, testCase := range testCases {
		merged = append(merged, testCase...)
	}

	return merged
}
