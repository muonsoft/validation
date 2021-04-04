package test

import (
	"testing"

	"github.com/muonsoft/validation/code"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/message"
	"github.com/muonsoft/validation/validator"
	"github.com/stretchr/testify/assert"
)

func TestValidator_ValidateBy_WhenConstraintExists_ExpectValidationByStoredConstraint(t *testing.T) {
	defer validator.Reset()
	constraint := it.IsNotBlank()
	err := validator.StoreConstraint("notBlank", constraint)
	if err != nil {
		t.Fatal("failed to store constraint:", err)
	}

	s := ""
	err = validator.ValidateString(&s, validator.ValidateBy("notBlank"))

	assertHasOneViolation(code.NotBlank, message.NotBlank, "")(t, err)
}

func TestValidator_ValidateBy_WhenConstraintDoesNotExist_ExpectError(t *testing.T) {
	defer validator.Reset()

	s := ""
	err := validator.ValidateString(&s, validator.ValidateBy("notBlank"))

	assert.EqualError(t, err, `failed to set up constraint "notFoundConstraint": constraint with key "notBlank" is not stored in the validator`)
}

func TestValidator_StoreConstraint_WhenConstraintExists_ExpectError(t *testing.T) {
	defer validator.Reset()
	constraint := it.IsNotBlank()
	err := validator.StoreConstraint("notBlank", constraint)
	if err != nil {
		t.Fatal("failed to store constraint:", err)
	}

	err = validator.StoreConstraint("notBlank", constraint)

	assert.EqualError(t, err, `constraint with key "notBlank" already stored`)
}
