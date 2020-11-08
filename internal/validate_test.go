package internal

import (
	"fmt"
	"testing"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/validationtest"
	"github.com/stretchr/testify/assert"
)

func TestValidate_NoViolations_Nil(t *testing.T) {
	err := validation.Validate(nil, nil)

	assert.NoError(t, err)
}

func TestValidate_SingleViolation_ViolationInList(t *testing.T) {
	violation := validation.NewViolation("code", "message", nil, nil)
	wrapped := fmt.Errorf("error: %w", violation)

	err := validation.Validate(nil, wrapped)

	validationtest.AssertIsViolationList(t, err, func(t *testing.T, violations validation.ViolationList) bool {
		return assert.Len(t, violations, 1) && assert.Equal(t, violation, violations[0])
	})
}

func TestValidate_ViolationList_ViolationsInList(t *testing.T) {
	violation := validation.NewViolation("code", "message", nil, nil)
	violations := validation.ViolationList{violation}
	wrapped := fmt.Errorf("error: %w", violations)

	err := validation.Validate(nil, wrapped)

	validationtest.AssertIsViolationList(t, err, func(t *testing.T, violations validation.ViolationList) bool {
		return assert.Len(t, violations, 1) && assert.Equal(t, violation, violations[0])
	})
}

func TestValidate_UnexpectedError_Error(t *testing.T) {
	unexpectedError := fmt.Errorf("error")

	err := validation.Validate(unexpectedError)

	assert.Equal(t, unexpectedError, err)
}
