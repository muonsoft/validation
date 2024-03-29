package test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/muonsoft/validation"
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
		{boolType, validation.Bool(false, nilConstraint{})},
		{"number", validation.Number(intValue(0), nilConstraint{})},
		{stringType, validation.String("", nilConstraint{})},
		{iterableType, validation.Iterable([]string{}, nilConstraint{})},
		{countableType, validation.Countable(0, nilConstraint{})},
		{timeType, validation.Time(time.Time{}, nilConstraint{})},
	}
	for _, test := range tests {
		t.Run(test.valueType, func(t *testing.T) {
			err := validator.Validate(context.Background(), test.argument)

			assertIsInapplicableConstraintError(t, err, test.valueType)
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
			err := validator.Validate(context.Background(), test.argument)

			assert.EqualError(t, err, test.expectedError)
		})
	}
}

func TestValidator_Validate_WhenInvalidConstraintAtPropertyPath_ExpectErrorWithPropertyPath(t *testing.T) {
	err := validator.Validate(
		context.Background(),
		validation.String(
			"",
			validation.PropertyName("properties"),
			validation.ArrayIndex(1),
			validation.PropertyName("error"),
			errConstraint{},
		),
	)

	assert.EqualError(t, err, `failed to set up constraint "errConstraint" at path "properties[1].error": error`)
}
