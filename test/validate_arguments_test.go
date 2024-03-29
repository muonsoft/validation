package test

import (
	"context"
	"testing"
	"time"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/code"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/validationtest"
	"github.com/muonsoft/validation/validator"
	"github.com/stretchr/testify/assert"
)

func TestValidate_WhenArgumentForGivenType_ExpectValidationExecuted(t *testing.T) {
	tests := []struct {
		name     string
		argument validation.Argument
	}{
		{"Value", validation.Value(stringValue(""), it.IsNotBlank())},
		{"Bool", validation.Bool(false, it.IsNotBlank())},
		{"NilBool", validation.NilBool(boolValue(false), it.IsNotBlank())},
		{"Number", validation.Number(0, it.IsNotBlank())},
		{"String", validation.String("", it.IsNotBlank())},
		{"NilString", validation.NilString(stringValue(""), it.IsNotBlank())},
		{"Iterable", validation.Iterable([]string{}, it.IsNotBlank())},
		{"Countable", validation.Countable(0, it.IsNotBlank())},
		{"Time", validation.Time(time.Time{}, it.IsNotBlank())},
		{"NilTime", validation.NilTime(nilTime, it.IsNotBlank())},
		{"Each", validation.Each([]string{""}, it.IsNotBlank())},
		{"EachString", validation.EachString([]string{""}, it.IsNotBlank())},
		{"Valid", validation.Valid(mockValidatableString{""}, it.IsNotBlank())},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validator.Validate(context.Background(), test.argument)

			validationtest.AssertIsViolationList(t, err, func(t *testing.T, violations []validation.Violation) bool {
				t.Helper()
				return assert.Len(t, violations, 1) && assert.Equal(t, code.NotBlank, violations[0].Code())
			})
		})
	}
}

func TestValidate_WhenPropertyArgument_ExpectValidPathInViolation(t *testing.T) {
	opts := []validation.Option{validation.PropertyName("internal"), it.IsNotBlank()}

	tests := []struct {
		name         string
		argument     validation.Argument
		expectedPath string
	}{
		{"PropertyValue", validation.PropertyValue("property", stringValue(""), opts...), "property.internal"},
		{"BoolProperty", validation.BoolProperty("property", false, opts...), "property.internal"},
		{"NilBoolProperty", validation.NilBoolProperty("property", boolValue(false), opts...), "property.internal"},
		{"NumberProperty", validation.NumberProperty("property", 0, opts...), "property.internal"},
		{"StringProperty", validation.StringProperty("property", "", opts...), "property.internal"},
		{"NilStringProperty", validation.NilStringProperty("property", stringValue(""), opts...), "property.internal"},
		{"IterableProperty", validation.IterableProperty("property", []string{}, opts...), "property.internal"},
		{"CountableProperty", validation.CountableProperty("property", 0, opts...), "property.internal"},
		{"TimeProperty", validation.TimeProperty("property", time.Time{}, opts...), "property.internal"},
		{"NilTimeProperty", validation.NilTimeProperty("property", nilTime, opts...), "property.internal"},
		{"EachProperty", validation.EachProperty("property", []string{""}, opts...), "property.internal[0]"},
		{"EachStringProperty", validation.EachStringProperty("property", []string{""}, opts...), "property.internal[0]"},
		{"ValidProperty", validation.ValidProperty("property", mockValidatableString{""}, opts...), "property.internal.value"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validator.Validate(context.Background(), test.argument)

			validationtest.AssertIsViolationList(t, err, func(t *testing.T, violations []validation.Violation) bool {
				t.Helper()
				return assert.Len(t, violations, 1) &&
					assert.Equal(t, test.expectedPath, violations[0].PropertyPath().String())
			})
		})
	}
}
