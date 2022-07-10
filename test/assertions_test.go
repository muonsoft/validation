package test

import (
	"testing"

	"github.com/muonsoft/validation/validationtest"
	"github.com/stretchr/testify/assert"
)

func assertHasOneViolation(expectedError error, message string) func(t *testing.T, err error) {
	return func(t *testing.T, err error) {
		t.Helper()
		validationtest.Assert(t, err).IsViolationList().WithOneViolation().WithError(expectedError).WithMessage(message)
	}
}

func assertHasOneViolationAtPath(expectedError error, message, path string) func(t *testing.T, err error) {
	return func(t *testing.T, err error) {
		t.Helper()
		validationtest.Assert(t, err).IsViolationList().
			WithOneViolation().
			WithError(expectedError).
			WithMessage(message).
			WithPropertyPath(path)
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
