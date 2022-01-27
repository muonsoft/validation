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

	validationtest.Assert(t, err).IsViolationList().WithCodes(code.Blank, code.LengthTooFew)
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

	validationtest.Assert(t, err).IsViolationList().WithCodes(code.Blank, code.LengthTooFew)
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

	assert.EqualError(t, err, `failed to set up constraint "ConditionalConstraint": then branch of conditional constraint not set`)
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

	validationtest.Assert(t, err).IsViolationList().WithCodes(code.Blank)
}

func TestValidate_WhenSequentiallyConstraintIsDisabled_ExpectNoErrors(t *testing.T) {
	value := "foo"

	err := validator.Validate(
		context.Background(),
		validation.String(
			value,
			validation.Sequentially(it.IsBlank()).When(false),
		),
	)

	assert.NoError(t, err)
}

func TestValidate_WhenSequentiallyGroupsNotMatch_ExpectNoErrors(t *testing.T) {
	value := "foo"

	err := validator.Validate(
		context.Background(),
		validation.String(
			value,
			validation.Sequentially(it.IsBlank()).WhenGroups(testGroup),
		),
	)

	assert.NoError(t, err)
}

func TestValidate_WhenSequentiallyConstraintsNotSet_ExpectError(t *testing.T) {
	value := "bar"

	err := validator.Validate(
		context.Background(),
		validation.String(value, validation.Sequentially()),
	)

	assert.EqualError(t, err, `failed to set up constraint "SequentialConstraint": constraints for sequentially validation not set`)
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

	validationtest.Assert(t, err).IsViolationList().WithCodes(code.Blank, code.LengthTooFew)
}

func TestValidate_WhenAtLeastOneOfConstraintIsDisabled_ExpectNoError(t *testing.T) {
	value := "foo"

	err := validator.Validate(
		context.Background(),
		validation.String(
			value,
			validation.AtLeastOneOf(it.IsBlank()).When(false),
		),
	)

	assert.NoError(t, err)
}

func TestValidate_WhenAtLeastOneOfConstraintGroupsNotMatch_ExpectNoError(t *testing.T) {
	value := "foo"

	err := validator.Validate(
		context.Background(),
		validation.String(
			value,
			validation.AtLeastOneOf(it.IsBlank()).WhenGroups(testGroup),
		),
	)

	assert.NoError(t, err)
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

	assert.EqualError(t, err, `failed to set up constraint "AtLeastOneOfConstraint": constraints for at least one of validation not set`)
}

func TestValidate_WhenCompoundWithFailingConstraint_ExpectViolation(t *testing.T) {
	value := "bar"
	isEmployeeEmail := validation.Compound(it.HasMinLength(5), it.IsEmail())

	err := validator.Validate(
		context.Background(),
		validation.String(value, isEmployeeEmail),
	)

	validationtest.Assert(t, err).IsViolationList().WithCodes(code.LengthTooFew, code.InvalidEmail)
}

func TestValidate_WhenCompoundIsDisabled_ExpectNoError(t *testing.T) {
	value := "bar"
	isEmployeeEmail := validation.Compound(it.HasMinLength(5), it.IsEmail())

	err := validator.Validate(
		context.Background(),
		validation.String(value, isEmployeeEmail.When(false)),
	)

	assert.NoError(t, err)
}

func TestValidate_WhenCompoundGroupsNotMatch_ExpectNoError(t *testing.T) {
	value := "bar"
	isEmployeeEmail := validation.Compound(it.HasMinLength(5), it.IsEmail())

	err := validator.Validate(
		context.Background(),
		validation.String(value, isEmployeeEmail.WhenGroups(testGroup)),
	)

	assert.NoError(t, err)
}

func TestValidate_WhenCompoundConstraintsNotSet_ExpectError(t *testing.T) {
	value := "bar"
	isEmployeeEmail := validation.Compound()

	err := validator.Validate(
		context.Background(),
		validation.String(value, isEmployeeEmail),
	)

	assert.EqualError(t, err, `failed to set up constraint "CompoundConstraint": constraints for compound validation not set`)
}
