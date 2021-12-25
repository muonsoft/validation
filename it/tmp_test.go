package it_test

import (
	"testing"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/stretchr/testify/assert"
)

func Number[T validation.Numeric](value T, constraints ...validation.NumberConstraint[T]) error {
	for _, constraint := range constraints {
		err := constraint.ValidateNumber(&value, validation.Scope{})
		if err != nil {
			return err
		}
	}
	return nil
}

func Test(t *testing.T) {
	err := Number[int](1, it.IsGreaterThan(0))
	assert.NoError(t, err)
}
