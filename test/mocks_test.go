package test

import (
	"context"
	"time"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"golang.org/x/text/language"
)

var nilTime *time.Time

func boolValue(b bool) *bool {
	return &b
}

func intValue(i int) *int {
	return &i
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

func givenLocation(name string) *time.Location {
	loc, _ := time.LoadLocation(name)
	return loc
}

type mockViolation struct {
	err             string
	code            string
	message         string
	messageTemplate string
	parameters      []validation.TemplateParameter
	propertyPath    *validation.PropertyPath
}

func (mock *mockViolation) Is(codes ...string) bool {
	return false
}

func (mock *mockViolation) Error() string {
	return mock.err
}

func (mock *mockViolation) Code() string {
	return mock.code
}

func (mock *mockViolation) Message() string {
	return mock.message
}

func (mock *mockViolation) MessageTemplate() string {
	return mock.messageTemplate
}

func (mock *mockViolation) Parameters() []validation.TemplateParameter {
	return mock.parameters
}

func (mock *mockViolation) PropertyPath() *validation.PropertyPath {
	return mock.propertyPath
}

func mockNewViolationFunc() validation.ViolationFactory {
	return validation.NewViolationFunc(func(
		code, messageTemplate string,
		pluralCount int,
		parameters []validation.TemplateParameter,
		propertyPath *validation.PropertyPath,
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

func (mock mockValidatableString) Validate(ctx context.Context, validator *validation.Validator) error {
	return validator.Validate(
		ctx,
		validation.StringProperty(
			"value",
			mock.value,
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

func (mock mockValidatableStruct) Validate(ctx context.Context, validator *validation.Validator) error {
	return validator.Validate(
		ctx,
		validation.NumberProperty[int64](
			"intValue",
			mock.intValue,
			it.IsNotBlankNumber[int64](),
		),
		validation.NumberProperty[float64](
			"floatValue",
			mock.floatValue,
			it.IsNotBlankNumber[float64](),
		),
		validation.StringProperty(
			"stringValue",
			mock.stringValue,
			it.IsNotBlank(),
		),
		validation.ValidProperty(
			"structValue",
			&mock.structValue,
		),
	)
}

type mockTranslator struct {
	translate func(tag language.Tag, message string, pluralCount int) string
}

func (m mockTranslator) Translate(tag language.Tag, message string, pluralCount int) string {
	return m.translate(tag, message, pluralCount)
}
