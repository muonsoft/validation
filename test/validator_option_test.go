package test

import (
	"context"
	"testing"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/validationtest"
	"github.com/muonsoft/validation/validator"
	"github.com/stretchr/testify/assert"
)

func TestWhenGlobalValidatorWithOverriddenNewViolation_ExpectCustomViolation(t *testing.T) {
	err := validator.SetUp(validation.SetViolationFactory(mockNewViolationFunc()))
	if err != nil {
		t.Fatal(err)
	}
	defer validator.SetUp()

	err = validator.Validate(context.Background(), validation.String("", it.IsNotBlank()))

	validationtest.AssertIsViolationList(t, err, func(t *testing.T, violations []validation.Violation) bool {
		t.Helper()
		return assert.Len(t, violations, 1) && assert.IsType(t, &mockViolation{}, violations[0])
	})
}

func TestWhenValidatorWithOverriddenNewViolation_ExpectCustomViolation(t *testing.T) {
	v := newValidator(t, validation.SetViolationFactory(mockNewViolationFunc()))

	err := v.Validate(context.Background(), validation.String("", it.IsNotBlank()))

	validationtest.AssertIsViolationList(t, err, func(t *testing.T, violations []validation.Violation) bool {
		t.Helper()
		return assert.Len(t, violations, 1) && assert.IsType(t, &mockViolation{}, violations[0])
	})
}
