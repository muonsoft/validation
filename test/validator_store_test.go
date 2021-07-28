package test

import (
	"context"
	"testing"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/code"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/message"
	"github.com/stretchr/testify/assert"
)

func TestValidator_ValidateBy_WhenConstraintExists_ExpectValidationByStoredConstraint(t *testing.T) {
	validator := newValidator(t, validation.StoredConstraint("notBlank", it.IsNotBlank()))

	s := ""
	err := validator.ValidateString(context.Background(), &s, validator.ValidateBy("notBlank"))

	assertHasOneViolation(code.NotBlank, message.NotBlank)(t, err)
}

func TestValidator_ValidateBy_WhenConstraintDoesNotExist_ExpectError(t *testing.T) {
	validator := newValidator(t)

	s := ""
	err := validator.ValidateString(context.Background(), &s, validator.ValidateBy("notBlank"))

	assert.EqualError(t, err, `failed to set up constraint "notFoundConstraint": constraint with key "notBlank" is not stored in the validator`)
}

func TestValidator_StoreConstraint_WhenConstraintExists_ExpectError(t *testing.T) {
	validator, err := validation.NewValidator(
		validation.StoredConstraint("key", it.IsNotBlank()),
		validation.StoredConstraint("key", it.IsBlank()),
	)

	assert.Nil(t, validator)
	assert.EqualError(t, err, `constraint with key "key" already stored`)
}
