package test

import (
	"errors"
	"testing"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/validationtest"
	"github.com/stretchr/testify/assert"
)

func assertHasOneViolation(code, message, path string) func(t *testing.T, err error) {
	return func(t *testing.T, err error) {
		t.Helper()
		validationtest.AssertIsViolationList(t, err, func(t *testing.T, violations validation.ViolationList) bool {
			t.Helper()

			if assert.Len(t, violations, 1) {
				return assert.Equal(t, code, violations[0].Code()) &&
					assert.Equal(t, message, violations[0].Message()) &&
					assert.Equal(t, path, violations[0].PropertyPath().String())
			}

			return false
		})
	}
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	assert.NoError(t, err)
}

func assertError(expectedError string) func(t *testing.T, err error) {
	return func(t *testing.T, err error) {
		t.Helper()
		assert.EqualError(t, err, expectedError)
	}
}

func assertIsInapplicableConstraintError(t *testing.T, err error, valueType string) {
	t.Helper()
	var inapplicableConstraint validation.InapplicableConstraintError

	if !errors.As(err, &inapplicableConstraint) {
		t.Errorf("failed asserting that error is InapplicableConstraintError")
		return
	}

	assert.Equal(t, valueType, inapplicableConstraint.ValueType)
}
