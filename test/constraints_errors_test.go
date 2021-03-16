package test

import (
	"errors"
	"testing"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/validator"
	"github.com/stretchr/testify/assert"
)

type nilConstraint struct {
}

func (c nilConstraint) SetUp() error {
	return nil
}

func (c nilConstraint) Name() string {
	return "nilConstraint"
}

func (c nilConstraint) ValidateNil(scope validation.Scope) error {
	return nil
}

type errConstraint struct {
}

func (c errConstraint) SetUp() error {
	return errors.New("error")
}

func (c errConstraint) Name() string {
	return "errConstraint"
}

func (c errConstraint) ValidateNil(scope validation.Scope) error {
	return nil
}

func TestValidator_Validate_WhenInapplicableConstraint_ExpectError(t *testing.T) {
	tests := []struct {
		valueType string
		argument  validation.Argument
	}{
		{boolType, validation.Bool(nil, nilConstraint{})},
		{"number", validation.Number(intValue(0), nilConstraint{})},
		{stringType, validation.String(nilString, nilConstraint{})},
		{iterableType, validation.Iterable([]string{}, nilConstraint{})},
		{countableType, validation.Countable(0, nilConstraint{})},
		{timeType, validation.Time(nil, nilConstraint{})},
	}
	for _, test := range tests {
		t.Run(test.valueType, func(t *testing.T) {
			err := validator.Validate(test.argument)

			assertIsInapplicableConstraintError(t, err, test.valueType)
		})
	}
}

func TestValidator_Validate_WhenInvalidConstraint_ExpectError(t *testing.T) {
	tests := []struct {
		constraint    string
		argument      validation.Argument
		expectedError string
	}{
		{
			constraint:    it.ChoiceConstraint{}.Name(),
			argument:      validation.String(nil, it.IsOneOfStrings()),
			expectedError: `failed to set up constraint "ChoiceConstraint": empty list of choices`,
		},
		{
			constraint:    it.RegexConstraint{}.Name(),
			argument:      validation.String(nil, it.Matches("")),
			expectedError: `failed to set up constraint "RegexConstraint": empty pattern`,
		},
		{
			constraint:    it.RegexConstraint{}.Name(),
			argument:      validation.String(nil, it.Matches("invalid_path")),
			expectedError: `failed to set up constraint "RegexConstraint": invalid pattern`,
		},
	}
	for _, test := range tests {
		t.Run(test.constraint, func(t *testing.T) {
			err := validator.Validate(test.argument)

			assert.EqualError(t, err, test.expectedError)
		})
	}
}

func TestValidator_Validate_WhenInvalidValue_ExpectError(t *testing.T) {
	tests := []struct {
		name          string
		argument      validation.Argument
		expectedError string
	}{
		{
			name:          "invalid number",
			argument:      validation.Number("string"),
			expectedError: `cannot convert value "string" to number: value of type string is not numeric`,
		},
		{
			name:          "invalid iterable",
			argument:      validation.Iterable("string"),
			expectedError: `cannot convert value "string" to iterable: value of type string is not iterable`,
		},
		{
			name:          "invalid each",
			argument:      validation.Each("string"),
			expectedError: `cannot convert value "string" to iterable: value of type string is not iterable`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validator.Validate(test.argument)

			assert.EqualError(t, err, test.expectedError)
		})
	}
}

func TestValidator_Validate_WhenInvalidConstraintAtPropertyPath_ExpectErrorWithPropertyPath(t *testing.T) {
	err := validator.Validate(
		validation.String(
			nil,
			validation.PropertyName("properties"),
			validation.ArrayIndex(1),
			validation.PropertyName("error"),
			errConstraint{},
		),
	)

	assert.EqualError(t, err, `failed to set up constraint "errConstraint" at path "properties[1].error": error`)
}
