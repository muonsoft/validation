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

func TestValidateString_WhenConditionIsTrue_ExpectAllConstraintsOfThenBranchApplied(t *testing.T) {
	value := foo

	err := validator.ValidateString(
		&value,
		validation.When(true).
			Then(
				it.IsBlank(),
				it.HasMinLength(5),
			),
	)

	validationtest.AssertIsViolationList(t, err, func(t *testing.T, violations []validation.Violation) bool {
		t.Helper()
		return assert.Len(t, violations, 2) &&
			assert.Equal(t, code.Blank, violations[0].Code()) &&
			assert.Equal(t, code.LengthTooFew, violations[1].Code())
	})
}

func TestValidateString_WhenConditionIsFalse_ExpectAllConstraintsOfElseBranchApplied(t *testing.T) {
	value := bar

	err := validator.ValidateString(
		&value,
		validation.When(false).
			Then(
				it.IsNotBlank(),
			).
			Else(
				it.IsBlank(),
				it.HasMinLength(5),
			),
	)

	validationtest.AssertIsViolationList(t, err, func(t *testing.T, violations []validation.Violation) bool {
		t.Helper()
		return assert.Len(t, violations, 2) &&
			assert.Equal(t, code.Blank, violations[0].Code()) &&
			assert.Equal(t, code.LengthTooFew, violations[1].Code())
	})
}

func TestValidateString_WhenConditionIsFalseAndNoElseBranch_ExpectNoViolations(t *testing.T) {
	value := foo

	err := validator.ValidateString(
		&value,
		validation.When(false).
			Then(
				it.IsNotBlank(),
			),
	)

	assertNoError(t, err)
}

func TestValidateString_WhenThenBranchIsNotSet_ExpectError(t *testing.T) {
	value := bar

	err := validator.ValidateString(
		&value,
		validation.When(true),
	)

	assert.Error(t, err, "then branch of conditional constraint not set")
}
