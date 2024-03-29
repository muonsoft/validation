package test

import (
	"context"
	"testing"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/code"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/validationtest"
	"github.com/muonsoft/validation/validator"
	"github.com/stretchr/testify/assert"
)

func TestValidate_WhenConditionIsTrue_ExpectAllConstraintsOfThenBranchApplied(t *testing.T) {
	value := "foo"

	err := validator.Validate(
		context.Background(),
		validation.String(
			value,
			validation.When(true).
				Then(
					it.IsBlank(),
					it.HasMinLength(5),
				),
		),
	)

	validationtest.AssertIsViolationList(t, err, func(t *testing.T, violations []validation.Violation) bool {
		t.Helper()
		return assert.Len(t, violations, 2) &&
			assert.Equal(t, code.Blank, violations[0].Code()) &&
			assert.Equal(t, code.LengthTooFew, violations[1].Code())
	})
}

func TestValidate_WhenConditionIsFalse_ExpectAllConstraintsOfElseBranchApplied(t *testing.T) {
	value := "bar"

	err := validator.Validate(
		context.Background(),
		validation.String(
			value,
			validation.When(false).
				Then(
					it.IsNotBlank(),
				).
				Else(
					it.IsBlank(),
					it.HasMinLength(5),
				),
		),
	)

	validationtest.AssertIsViolationList(t, err, func(t *testing.T, violations []validation.Violation) bool {
		t.Helper()
		return assert.Len(t, violations, 2) &&
			assert.Equal(t, code.Blank, violations[0].Code()) &&
			assert.Equal(t, code.LengthTooFew, violations[1].Code())
	})
}

func TestValidate_WhenConditionIsFalseAndNoElseBranch_ExpectNoViolations(t *testing.T) {
	value := "foo"

	err := validator.Validate(
		context.Background(),
		validation.String(
			value,
			validation.When(false).Then(it.IsNotBlank()),
		),
	)

	assertNoError(t, err)
}

func TestValidate_WhenThenBranchIsNotSet_ExpectError(t *testing.T) {
	value := "bar"

	err := validator.Validate(
		context.Background(),
		validation.String(value, validation.When(true)),
	)

	assert.Error(t, err, "then branch of conditional constraint not set")
}

func TestValidate_WhenInvalidValueAtFirstConstraintOfSequentiallyConstraint_ExpectOneViolation(t *testing.T) {
	value := "foo"

	err := validator.Validate(
		context.Background(),
		validation.String(
			value,
			validation.Sequentially(
				it.IsBlank(),
				it.HasMinLength(5),
			),
		),
	)

	validationtest.AssertIsViolationList(t, err, func(t *testing.T, violations []validation.Violation) bool {
		t.Helper()
		return assert.Len(t, violations, 1) &&
			assert.Equal(t, code.Blank, violations[0].Code())
	})
}

func TestValidate_WhenSequentiallyConstraintsNotSet_ExpectError(t *testing.T) {
	value := "bar"

	err := validator.Validate(
		context.Background(),
		validation.String(value, validation.Sequentially()),
	)

	assert.Error(t, err, "constraints for sequentially validation not set")
}

func TestValidate_WhenInvalidValueAtFirstConstraintOfAtLeastOneOfConstraint_ExpectAllViolation(t *testing.T) {
	value := "foo"

	err := validator.Validate(
		context.Background(),
		validation.String(
			value,
			validation.AtLeastOneOf(
				it.IsBlank(),
				it.HasMinLength(5),
			),
		),
	)

	validationtest.AssertIsViolationList(t, err, func(t *testing.T, violations []validation.Violation) bool {
		t.Helper()
		return assert.Len(t, violations, 2) &&
			assert.Equal(t, code.Blank, violations[0].Code()) &&
			assert.Equal(t, code.LengthTooFew, violations[1].Code())
	})
}

func TestValidate_WhenInvalidValueAtSecondConstraintOfAtLeastOneOfConstraint_ExpectNoViolation(t *testing.T) {
	value := "foo"

	err := validator.Validate(
		context.Background(),
		validation.String(
			value,
			validation.AtLeastOneOf(
				it.IsEqualToString("bar"),
				it.IsEqualToString("foo"),
			),
		),
	)

	assertNoError(t, err)
}

func TestValidate_WhenAtLeastOneOfConstraintsNotSet_ExpectError(t *testing.T) {
	value := "bar"

	err := validator.Validate(
		context.Background(),
		validation.String(value, validation.AtLeastOneOf()),
	)

	assert.Error(t, err, "constraints for at least one of validation not set")
}

func TestValidate_Compound_ExpectNoViolation(t *testing.T) {
	value := "bar"
	isEmployeeEmail := validation.Compound(it.HasMinLength(5), it.IsEmail())

	err := validator.Validate(
		context.Background(),
		validation.String(value, isEmployeeEmail),
	)

	validationtest.AssertIsViolationList(t, err, func(t *testing.T, violations []validation.Violation) bool {
		t.Helper()
		return assert.Len(t, violations, 2) &&
			assert.Equal(t, code.LengthTooFew, violations[0].Code()) &&
			assert.Equal(t, code.InvalidEmail, violations[1].Code())
	})
}

func TestValidate_WhenCompoundConstraintsNotSet_ExpectError(t *testing.T) {
	value := "bar"
	isEmployeeEmail := validation.Compound()

	err := validator.Validate(
		context.Background(),
		validation.String(value, isEmployeeEmail),
	)

	assert.Error(t, err, "constraints for compound validation not set")
}
