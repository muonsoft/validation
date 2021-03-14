package test

import (
	"errors"
	"testing"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/validator"
	"github.com/stretchr/testify/assert"
)

func TestValidator_Validate_WhenNotSupportedType_ExpectError(t *testing.T) {
	err := validator.Validate(validation.Value(func() {}))

	var notValidatable *validation.NotValidatableError
	assert.True(t, errors.As(err, &notValidatable))
	assert.EqualError(t, err, "value of type func is not validatable")
}
