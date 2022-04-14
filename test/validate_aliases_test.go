package test

import (
	"context"
	"testing"
	"time"

	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/validationtest"
)

func TestValidate_ArgumentAliases_WhenAliasMethodForGivenType_ExpectValidationExecuted(t *testing.T) {
	validator := newValidator(t)

	tests := []struct {
		name string
		err  error
	}{
		{"ValidateBool", validator.ValidateBool(context.Background(), false, it.IsNotBlank())},
		{"ValidateInt", validator.ValidateInt(context.Background(), 0, it.IsNotBlankNumber[int]())},
		{"ValidateFloat", validator.ValidateFloat(context.Background(), 0, it.IsNotBlankNumber[float64]())},
		{"ValidateString", validator.ValidateString(context.Background(), "", it.IsNotBlank())},
		{"ValidateStrings", validator.ValidateStrings(context.Background(), []string{"foo", "foo"}, it.HasUniqueValues[string]())},
		{"ValidateCountable", validator.ValidateCountable(context.Background(), 0, it.IsNotBlank())},
		{"ValidateTime", validator.ValidateTime(context.Background(), time.Time{}, it.IsNotBlank())},
		{"ValidateEachString", validator.ValidateEachString(context.Background(), []string{""}, it.IsNotBlank())},
		{"ValidateIt", validator.ValidateIt(context.Background(), mockValidatableString{""})},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			validationtest.Assert(t, test.err).IsViolationList().WithOneViolation()
		})
	}
}
