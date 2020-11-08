package validation

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestViolation_Error_MessageOnly_ErrorWithMessage(t *testing.T) {
	violation := internalViolation{message: "message"}

	err := violation.Error()

	assert.Equal(t, "violation: message", err)
}

func TestViolation_Error_MessageAndPropertyPath_ErrorWithPropertyPathAndMessage(t *testing.T) {
	violation := internalViolation{
		message:      "message",
		propertyPath: PropertyPath{PropertyNameElement{"propertyPath"}},
	}

	err := violation.Error()

	assert.Equal(t, "violation at 'propertyPath': message", err)
}

func TestViolationList_Error_CoupleOfViolations_JoinedMessage(t *testing.T) {
	violations := ViolationList{
		internalViolation{
			message:      "first message",
			propertyPath: PropertyPath{PropertyNameElement{"path"}, ArrayIndexElement{0}},
		},
		internalViolation{
			message:      "second message",
			propertyPath: PropertyPath{PropertyNameElement{"path"}, ArrayIndexElement{1}},
		},
	}

	err := violations.Error()

	assert.Equal(t, "violation at 'path[0]': first message; violation at 'path[1]': second message", err)
}

func TestIsViolation_CustomError_False(t *testing.T) {
	err := errors.New("error")

	is := IsViolation(err)

	assert.False(t, is)
}

func TestIsViolation_Violation_True(t *testing.T) {
	err := &internalViolation{message: "message"}

	is := IsViolation(err)

	assert.True(t, is)
}

func TestIsViolationList_CustomError_False(t *testing.T) {
	err := errors.New("error")

	is := IsViolationList(err)

	assert.False(t, is)
}

func TestIsViolationList_Violation_True(t *testing.T) {
	err := ViolationList{internalViolation{message: "message"}}

	is := IsViolationList(err)

	assert.True(t, is)
}

func TestUnwrapViolation_WrappedViolation_UnwrappedViolation(t *testing.T) {
	wrapped := &internalViolation{message: "message"}
	err := fmt.Errorf("error: %w", wrapped)

	unwrapped, ok := UnwrapViolation(err)

	assert.True(t, ok)
	assert.Equal(t, wrapped, unwrapped)
}

func TestUnwrapViolationList_WrappedViolationList_UnwrappedViolationList(t *testing.T) {
	wrapped := ViolationList{internalViolation{message: "message"}}
	err := fmt.Errorf("error: %w", wrapped)

	unwrapped, ok := UnwrapViolationList(err)

	assert.True(t, ok)
	assert.Equal(t, wrapped, unwrapped)
}
