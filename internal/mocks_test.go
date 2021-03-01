package internal

import (
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"golang.org/x/text/language"

	"time"
)

var (
	nilBool    *bool
	nilInt     *int64
	nilUint    *uint64
	nilFloat   *float64
	nilString  *string
	nilTime    *time.Time
	emptyArray [0]string
	emptySlice []string
	emptyMap   map[string]string
)

func boolValue(b bool) *bool {
	return &b
}

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

func timeValue(t time.Time) *time.Time {
	return &t
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

func mockNewViolationFunc() validation.ViolationFactory {
	return validation.NewViolationFunc(func(
		code, messageTemplate string,
		pluralCount int,
		parameters map[string]string,
		propertyPath validation.PropertyPath,
		lang language.Tag,
	) validation.Violation {
		return &mockViolation{
			code:            code,
			messageTemplate: messageTemplate,
			parameters:      parameters,
			propertyPath:    propertyPath,
		}
	})
}

type mockValidatableString struct {
	value string
}

func (mock mockValidatableString) Validate(scope validation.Scope) error {
	return validation.InScope(scope).Validate(
		validation.String(
			&mock.value,
			validation.PropertyName("value"),
			it.IsNotBlank(),
		),
	)
}

type mockValidatableStruct struct {
	intValue    int64
	floatValue  float64
	stringValue string
	structValue mockValidatableString
}

func (mock mockValidatableStruct) Validate(scope validation.Scope) error {
	return validation.InScope(scope).Validate(
		validation.Number(
			mock.intValue,
			validation.PropertyName("intValue"),
			it.IsNotBlank(),
		),
		validation.Number(
			mock.floatValue,
			validation.PropertyName("floatValue"),
			it.IsNotBlank(),
		),
		validation.String(
			&mock.stringValue,
			validation.PropertyName("stringValue"),
			it.IsNotBlank(),
		),
		validation.Value(
			&mock.structValue,
			validation.PropertyName("structValue"),
		),
	)
}
