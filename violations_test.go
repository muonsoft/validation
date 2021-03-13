package validation

import (
	"encoding/json"
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
		propertyPath: PropertyPath{PropertyNameElement("propertyPath")},
	}

	err := violation.Error()

	assert.Equal(t, "violation at 'propertyPath': message", err)
}

func TestViolationList_Error_CoupleOfViolations_JoinedMessage(t *testing.T) {
	violations := ViolationList{
		internalViolation{
			message:      "first message",
			propertyPath: PropertyPath{PropertyNameElement("path"), ArrayIndexElement(0)},
		},
		internalViolation{
			message:      "second message",
			propertyPath: PropertyPath{PropertyNameElement("path"), ArrayIndexElement(1)},
		},
	}

	err := violations.Error()

	assert.Equal(t, "violation at 'path[0]': first message; violation at 'path[1]': second message", err)
}

func TestViolationList_Error_EmptyList_ErrorWithHelpMessage(t *testing.T) {
	violations := ViolationList{}

	err := violations.Error()

	assert.Equal(t, "the list of violations is empty, it looks like you forgot to use the AsError method somewhere", err)
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

func TestMarshalInternalViolationToJSON(t *testing.T) {
	tests := []struct {
		name         string
		violation    internalViolation
		expectedJSON string
	}{
		{
			name: "full data",
			violation: internalViolation{
				code:            "code",
				message:         "message",
				messageTemplate: "messageTemplate",
				parameters:      map[string]string{"name": "value"},
				propertyPath:    PropertyPath{PropertyNameElement("properties"), ArrayIndexElement(1), PropertyNameElement("name")},
			},
			expectedJSON: `{
				"code": "code",
				"message": "message",
				"propertyPath": "properties[1].name"
			}`,
		},
		{
			name:         "empty data",
			violation:    internalViolation{},
			expectedJSON: `{"code": "", "message": ""}`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			data, err := json.Marshal(test.violation)

			if assert.NoError(t, err) {
				assert.JSONEq(t, test.expectedJSON, string(data))
			}
		})
	}
}
