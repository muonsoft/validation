package test

import (
	"fmt"
	"testing"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/validationtest"
	"github.com/stretchr/testify/assert"
)

func TestFilter_WhenNoViolations_ExpectNil(t *testing.T) {
	err := validation.Filter(nil, nil)

	assert.NoError(t, err)
}

func TestFilter_WhenSingleViolation_ExpectViolationInList(t *testing.T) {
	violation := validation.GetScope().BuildViolation("code", "message").GetViolation()
	wrapped := fmt.Errorf("error: %w", violation)

	err := validation.Filter(nil, wrapped)

	validationtest.AssertIsViolationList(t, err, func(t *testing.T, violations validation.ViolationList) bool {
		t.Helper()
		return assert.Len(t, violations, 1) && assert.Equal(t, violation, violations[0])
	})
}

func TestFilter_WhenViolationList_ExpectViolationsInList(t *testing.T) {
	violation := validation.GetScope().BuildViolation("code", "message").GetViolation()
	violations := validation.ViolationList{violation}
	wrapped := fmt.Errorf("error: %w", violations)

	err := validation.Filter(nil, wrapped)

	validationtest.AssertIsViolationList(t, err, func(t *testing.T, violations validation.ViolationList) bool {
		t.Helper()
		return assert.Len(t, violations, 1) && assert.Equal(t, violation, violations[0])
	})
}

func TestFilter_WhenUnexpectedError_ExpectError(t *testing.T) {
	unexpectedError := fmt.Errorf("error")

	err := validation.Filter(unexpectedError)

	assert.Equal(t, unexpectedError, err)
}
