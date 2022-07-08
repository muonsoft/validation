package test

import (
	"context"
	"testing"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/validator"
	"github.com/stretchr/testify/assert"
)

type errConstraint struct{}

func (c errConstraint) ValidateString(ctx context.Context, validator *validation.Validator, value *string) error {
	return validator.CreateConstraintError("errConstraint", "description")
}

func TestValidator_Validate_WhenInvalidConstraintAtPropertyPath_ExpectErrorWithPropertyPath(t *testing.T) {
	err := validator.Validate(
		context.Background(),
		validation.String("", errConstraint{}).With(
			validation.PropertyName("properties"),
			validation.ArrayIndex(1),
			validation.PropertyName("error"),
		),
	)

	assert.EqualError(t, err, `failed to validate by errConstraint at path "properties[1].error": description`)
}
