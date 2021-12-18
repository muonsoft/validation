package test

import (
	"context"
	"testing"
	"time"

	"github.com/muonsoft/validation/code"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/validationtest"
)

func TestValidate_ArgumentAliases_WhenAliasMethodForGivenType_ExpectValidationExecuted(t *testing.T) {
	validator := newValidator(t)

	tests := []struct {
		name string
		err  error
	}{
		{"ValidateValue", validator.ValidateValue(context.Background(), "", it.IsNotBlank())},
		{"ValidateBool", validator.ValidateBool(context.Background(), false, it.IsNotBlank())},
		{"ValidateNumber", validator.ValidateNumber(context.Background(), 0, it.IsNotBlank())},
		{"ValidateString", validator.ValidateString(context.Background(), "", it.IsNotBlank())},
		{"ValidateIterable", validator.ValidateIterable(context.Background(), []string{}, it.IsNotBlank())},
		{"ValidateCountable", validator.ValidateCountable(context.Background(), 0, it.IsNotBlank())},
		{"ValidateTime", validator.ValidateTime(context.Background(), time.Time{}, it.IsNotBlank())},
		{"ValidateEach", validator.ValidateEach(context.Background(), []string{""}, it.IsNotBlank())},
		{"ValidateEachString", validator.ValidateEachString(context.Background(), []string{""}, it.IsNotBlank())},
		{"ValidateValidatable", validator.ValidateValidatable(context.Background(), mockValidatableString{""}, it.IsNotBlank())},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			validationtest.Assert(t, test.err).IsViolationList().WithCodes(code.NotBlank)
		})
	}
}
