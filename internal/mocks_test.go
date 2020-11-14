package internal

import "github.com/muonsoft/validation"

var (
	nilInt    *int64
	nilUint   *uint64
	nilFloat  *float64
	nilString *string
)

func intValue(i int64) *int64 {
	return &i
}

func uintValue(u uint64) *uint64 {
	return &u
}

func floatValue(f float64) *float64 {
	return &f
}

func stringValue(s string) *string {
	return &s
}

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
