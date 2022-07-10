package test

import (
	"context"
	"testing"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/validationtest"
	"github.com/stretchr/testify/assert"
)

const testGroup = "testGroup"

func TestValidationGroups_WhenBothDefaultNonSetGroups_ExpectViolation(t *testing.T) {
	err := newValidator(t).Validate(
		context.Background(),
		validation.String("", it.IsNotBlank()),
	)

	validationtest.Assert(t, err).IsViolationList().WithOneViolation().WithError(validation.ErrIsBlank)
}

func TestValidationGroups_WhenBothDefaultGroups_ExpectViolation(t *testing.T) {
	err := newValidator(t).WithGroups(validation.DefaultGroup).Validate(
		context.Background(),
		validation.String("", it.IsNotBlank().WhenGroups(validation.DefaultGroup)),
	)

	validationtest.Assert(t, err).IsViolationList().WithOneViolation().WithError(validation.ErrIsBlank)
}

func TestValidationGroups_WhenValidatorNotSetAndConstraintDefault_ExpectViolation(t *testing.T) {
	err := newValidator(t).Validate(
		context.Background(),
		validation.String("", it.IsNotBlank().WhenGroups(validation.DefaultGroup)),
	)

	validationtest.Assert(t, err).IsViolationList().WithOneViolation().WithError(validation.ErrIsBlank)
}

func TestValidationGroups_WhenValidatorDefaultAndConstraintNotSet_ExpectViolation(t *testing.T) {
	err := newValidator(t).WithGroups(validation.DefaultGroup).Validate(
		context.Background(),
		validation.String("", it.IsNotBlank().WhenGroups()),
	)

	validationtest.Assert(t, err).IsViolationList().WithOneViolation().WithError(validation.ErrIsBlank)
}

func TestValidationGroups_WhenBothGroupsAreSet_ExpectViolation(t *testing.T) {
	err := newValidator(t).WithGroups(testGroup).Validate(
		context.Background(),
		validation.String("", it.IsNotBlank().WhenGroups(testGroup)),
	)

	validationtest.Assert(t, err).IsViolationList().WithOneViolation().WithError(validation.ErrIsBlank)
}

func TestValidationGroups_WhenConstraintWithNonDefaultGroupAndValidationGroupsIsNotSet_ExpectNoViolations(t *testing.T) {
	err := newValidator(t).Validate(
		context.Background(),
		validation.String("", it.IsNotBlank().WhenGroups(testGroup)),
	)

	assert.NoError(t, err)
}

func TestValidationGroups_WhenConstraintWithNonDefaultGroupAndValidatorWithDefaultGroup_ExpectNoViolations(t *testing.T) {
	err := newValidator(t).WithGroups(validation.DefaultGroup).Validate(
		context.Background(),
		validation.String("", it.IsNotBlank().WhenGroups(testGroup)),
	)

	assert.NoError(t, err)
}

func TestValidationGroups_WhenValidatorWithNonDefaultGroup_ExpectNoViolations(t *testing.T) {
	err := newValidator(t).WithGroups(testGroup).Validate(
		context.Background(),
		validation.String("", it.IsNotBlank()),
	)

	assert.NoError(t, err)
}

func TestValidationGroups_WhenValidatorWithNonDefaultGroupAndConstraintWithDefaultGroup_ExpectNoViolations(t *testing.T) {
	err := newValidator(t).WithGroups(testGroup).Validate(
		context.Background(),
		validation.String("", it.IsNotBlank().WhenGroups(validation.DefaultGroup)),
	)

	assert.NoError(t, err)
}
