package internal

import (
	"testing"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/validationtest"
	"github.com/stretchr/testify/assert"
)

type mockViolation struct {
	err             string
	code            string
	message         string
	messageTemplate string
	parameters      map[string]string
	propertyPath    validation.PropertyPath
}

func (mock *mockViolation) Error() string {
	return mock.err
}

func (mock *mockViolation) GetCode() string {
	return mock.code
}

func (mock *mockViolation) GetMessage() string {
	return mock.message
}

func (mock *mockViolation) GetMessageTemplate() string {
	return mock.messageTemplate
}

func (mock *mockViolation) GetParameters() map[string]string {
	return mock.parameters
}

func (mock *mockViolation) GetPropertyPath() validation.PropertyPath {
	return mock.propertyPath
}

func TestWhenValidatorWithOverriddenNewViolation_ExpectCustomViolation(t *testing.T) {
	validator := validation.NewValidator(
		validation.OverrideNewViolation(mockNewViolationFunc()),
	)

	err := validator.ValidateString(nil, it.IsNotBlank())

	validationtest.AssertIsViolationList(t, err, func(t *testing.T, violations validation.ViolationList) bool {
		return assert.Len(t, violations, 1) && assert.IsType(t, &mockViolation{}, violations[0])
	})
}

func mockNewViolationFunc() func(
	code string,
	messageTemplate string,
	parameters map[string]string,
	propertyPath validation.PropertyPath,
) validation.Violation {
	return func(code, messageTemplate string, parameters map[string]string, propertyPath validation.PropertyPath) validation.Violation {
		return &mockViolation{code: code, messageTemplate: messageTemplate, parameters: parameters, propertyPath: propertyPath}
	}
}
