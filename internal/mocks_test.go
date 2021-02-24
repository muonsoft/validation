package internal

import (
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"golang.org/x/text/language"
)

var (
	nilBool    *bool
	nilInt     *int64
	nilUint    *uint64
	nilFloat   *float64
	nilString  *string
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
	pluralCount int,
	parameters map[string]string,
	propertyPath validation.PropertyPath,
	lang language.Tag,
) validation.Violation {
	return func(
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
	}
}

type mockValidatableString struct {
	value string
}

func (mock mockValidatableString) Validate(options ...validation.Option) error {
	return validation.Filter(
		validation.ValidateString(
			&mock.value,
			validation.PassOptions(options),
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

func (mock mockValidatableStruct) Validate(options ...validation.Option) error {
	validator, err := validation.WithOptions(options...)
	if err != nil {
		return err
	}

	return validation.Filter(
		validator.ValidateNumber(
			mock.intValue,
			validation.PropertyName("intValue"),
			it.IsNotBlank(),
		),
		validator.ValidateNumber(
			mock.floatValue,
			validation.PropertyName("floatValue"),
			it.IsNotBlank(),
		),
		validator.ValidateString(
			&mock.stringValue,
			validation.PropertyName("stringValue"),
			it.IsNotBlank(),
		),
		validator.Validate(
			&mock.structValue,
			validation.PropertyName("structValue"),
		),
	)
}
