package test

import (
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/code"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/validationtest"
	"github.com/muonsoft/validation/validator"
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestValidate_WhenInvalidValueAtFirstConstraintOfSequentiallyConstraint_ExpectOneViolation(t *testing.T) {
	value := foo

	err := validator.ValidateString(
		&value,
		validation.Sequentially(
			it.IsBlank(),
			it.HasMinLength(5),
		),
	)

	validationtest.AssertIsViolationList(t, err, func(t *testing.T, violations validation.ViolationList) bool {
		t.Helper()
		return assert.Len(t, violations, 1) &&
			assert.Equal(t, code.Blank, violations[0].Code())
	})
}

func TestValidate_WhenSequentiallyConstraintsNotSet_ExpectError(t *testing.T) {
	value := bar

	err := validator.ValidateString(
		&value,
		validation.Sequentially(),
	)

	assert.Error(t, err, "constraints for sequentially validation not set")
}
