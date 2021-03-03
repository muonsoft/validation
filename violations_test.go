package validation

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestViolation_Error_MessageOnly_ErrorWithMessage(t *testing.T) {
	violation := internalViolation{Message: "Message"}

	err := violation.Error()

	assert.Equal(t, "violation: Message", err)
}

func TestViolation_Error_MessageAndPropertyPath_ErrorWithPropertyPathAndMessage(t *testing.T) {
	violation := internalViolation{
		Message:      "Message",
		PropertyPath: PropertyPath{PropertyNameElement("propertyPath")},
	}

	err := violation.Error()

	assert.Equal(t, "violation at 'propertyPath': Message", err)
}

func TestViolationList_Error_CoupleOfViolations_JoinedMessage(t *testing.T) {
	violations := ViolationList{
		internalViolation{
			Message:      "first Message",
			PropertyPath: PropertyPath{PropertyNameElement("path"), ArrayIndexElement(0)},
		},
		internalViolation{
			Message:      "second Message",
			PropertyPath: PropertyPath{PropertyNameElement("path"), ArrayIndexElement(1)},
		},
	}

	err := violations.Error()

	assert.Equal(t, "violation at 'path[0]': first Message; violation at 'path[1]': second Message", err)
}

func TestIsViolation_CustomError_False(t *testing.T) {
	err := errors.New("error")

	is := IsViolation(err)

	assert.False(t, is)
}

func TestIsViolation_Violation_True(t *testing.T) {
	err := &internalViolation{Message: "Message"}

	is := IsViolation(err)

	assert.True(t, is)
}

func TestIsViolationList_CustomError_False(t *testing.T) {
	err := errors.New("error")

	is := IsViolationList(err)

	assert.False(t, is)
}

func TestIsViolationList_Violation_True(t *testing.T) {
	err := ViolationList{internalViolation{Message: "Message"}}

	is := IsViolationList(err)

	assert.True(t, is)
}

func TestUnwrapViolation_WrappedViolation_UnwrappedViolation(t *testing.T) {
	wrapped := &internalViolation{Message: "Message"}
	err := fmt.Errorf("error: %w", wrapped)

	unwrapped, ok := UnwrapViolation(err)

	assert.True(t, ok)
	assert.Equal(t, wrapped, unwrapped)
}

func TestUnwrapViolationList_WrappedViolationList_UnwrappedViolationList(t *testing.T) {
	wrapped := ViolationList{internalViolation{Message: "Message"}}
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
				Code:            "code",
				Message:         "message",
				MessageTemplate: "messageTemplate",
				Parameters:      map[string]string{"name": "value"},
				PropertyPath:    PropertyPath{PropertyNameElement("properties"), ArrayIndexElement(1), PropertyNameElement("name")},
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
