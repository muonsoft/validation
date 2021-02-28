package internal

import (
	"testing"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/validationtest"
	"github.com/stretchr/testify/assert"
)

func TestWhenGlobalValidatorWithOverriddenNewViolation_ExpectCustomViolation(t *testing.T) {
	validation.OverrideDefaults(validation.OverrideViolationFactory(mockNewViolationFunc()))
	defer validation.ResetDefaults()

	err := validation.ValidateString(nil, it.IsNotBlank())

	validationtest.AssertIsViolationList(t, err, func(t *testing.T, violations validation.ViolationList) bool {
		t.Helper()
		return assert.Len(t, violations, 1) && assert.IsType(t, &mockViolation{}, violations[0])
	})
}

func TestWhenValidatorWithOverriddenNewViolation_ExpectCustomViolation(t *testing.T) {
	validator, err := validation.NewValidator(
		validation.OverrideViolationFactory(mockNewViolationFunc()),
	)
	if err != nil {
		t.Fatal(err)
	}

	err = validator.ValidateString(nil, it.IsNotBlank())

	validationtest.AssertIsViolationList(t, err, func(t *testing.T, violations validation.ViolationList) bool {
		t.Helper()
		return assert.Len(t, violations, 1) && assert.IsType(t, &mockViolation{}, violations[0])
	})
}
