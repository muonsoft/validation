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
	v := newValidator(t,
		validation.SetViolationFactory(mockNewViolationFunc()),
	)
	validator.SetDefault(v)
	defer func() {
		d := newValidator(t)
		validator.SetDefault(d)
	}()

	err := validator.Validate(context.Background(), validation.String("", it.IsNotBlank()))

	validationtest.Assert(t, err).IsViolationList().Assert(func(tb testing.TB, violations []validation.Violation) {
		tb.Helper()
		assert.Len(tb, violations, 1)
		assert.IsType(tb, &mockViolation{}, violations[0])
	})
}

func TestWhenValidatorWithOverriddenNewViolation_ExpectCustomViolation(t *testing.T) {
	v := newValidator(t, validation.SetViolationFactory(mockNewViolationFunc()))

	err := v.Validate(context.Background(), validation.String("", it.IsNotBlank()))

	validationtest.Assert(t, err).IsViolationList().Assert(func(tb testing.TB, violations []validation.Violation) {
		tb.Helper()
		assert.Len(tb, violations, 1)
		assert.IsType(tb, &mockViolation{}, violations[0])
	})
}
