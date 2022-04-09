package test

import (
	"context"
	"errors"
	"testing"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/validator"
	"github.com/stretchr/testify/assert"
)

type nilConstraint struct{}

func (c nilConstraint) SetUp() error {
	return nil
}

func (c nilConstraint) Name() string {
	return "nilConstraint"
}

func (c nilConstraint) ValidateNil(scope validation.Scope) error {
	return nil
}

type errConstraint struct{}

func (c errConstraint) SetUp() error {
	return errors.New("error")
}

func (c errConstraint) Name() string {
	return "errConstraint"
}

func (c errConstraint) ValidateNil(scope validation.Scope) error {
	return nil
}

func (c errConstraint) ValidateString(value *string, scope validation.Scope) error {
	return scope.NewConstraintError("errConstraint", "description")
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
