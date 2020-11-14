package internal

import (
	"testing"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/validationtest"
	"github.com/stretchr/testify/assert"
)

func TestWhenGlobalValidatorWithOverriddenNewViolation_ExpectCustomViolation(t *testing.T) {
	validation.OverrideDefaults(validation.OverrideNewViolation(mockNewViolationFunc()))

	err := validation.ValidateString(nil, it.IsNotBlank())

	validationtest.AssertIsViolationList(t, err, func(t *testing.T, violations validation.ViolationList) bool {
		return assert.Len(t, violations, 1) && assert.IsType(t, &mockViolation{}, violations[0])
	})
}

func TestWhenValidatorWithOverriddenNewViolation_ExpectCustomViolation(t *testing.T) {
	validator := validation.NewValidator(
		validation.OverrideNewViolation(mockNewViolationFunc()),
	)

	err := validator.ValidateString(nil, it.IsNotBlank())

	validationtest.AssertIsViolationList(t, err, func(t *testing.T, violations validation.ViolationList) bool {
		return assert.Len(t, violations, 1) && assert.IsType(t, &mockViolation{}, violations[0])
	})
}
