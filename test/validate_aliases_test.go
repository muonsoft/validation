package test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/validationtest"
	"github.com/stretchr/testify/assert"
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

// singleViolationValidatable returns a single Violation from Validate (not ViolationList).
// Used to reproduce the bug where validateIt uses UnwrapViolationList and misses single Violation.
type singleViolationValidatable struct{}

var errSingleViolation = errors.New("single violation")

func (singleViolationValidatable) Validate(ctx context.Context, v *validation.Validator) error {
	return v.CreateViolation(ctx, errSingleViolation, "single violation message")
}

func TestValidateIt_WhenValidatableReturnsSingleViolation_ExpectViolationList(t *testing.T) {
	validator := newValidator(t)

	err := validator.ValidateIt(context.Background(), singleViolationValidatable{})

	assert.True(t, validation.IsViolationList(err), "ValidateIt must return ViolationList when Validatable returns single Violation")
	validationtest.Assert(t, err).IsViolationList().WithOneViolation()
}
