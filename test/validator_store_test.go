package test

import (
	"context"
	"testing"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/message"
	"github.com/stretchr/testify/assert"
)

func TestValidator_GetConstraint_WhenConstraintExists_ExpectValidationByStoredConstraint(t *testing.T) {
	validator := newValidator(t, validation.StoredConstraint("notBlank", it.IsNotBlank()))

	err := validator.Validate(
		context.Background(),
		validation.String("", validator.GetConstraint("notBlank").(validation.StringConstraint)),
	)

	assertHasOneViolation(validation.ErrIsBlank, message.IsBlank)(t, err)
}

func TestValidator_StoreConstraint_WhenConstraintExists_ExpectError(t *testing.T) {
	validator, err := validation.NewValidator(
		validation.StoredConstraint("key", it.IsNotBlank()),
		validation.StoredConstraint("key", it.IsBlank()),
	)

	assert.Nil(t, validator)
	assert.EqualError(t, err, `constraint with key "key" already stored`)
}
