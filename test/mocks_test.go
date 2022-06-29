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
	err             error
	message         string
	messageTemplate string
	parameters      []validation.TemplateParameter
	propertyPath    *validation.PropertyPath
}

func (mock *mockViolation) Is(target error) bool                       { return mock.err == target }
func (mock *mockViolation) Unwrap() error                              { return mock.err }
func (mock *mockViolation) Error() string                              { return mock.err.Error() }
func (mock *mockViolation) Message() string                            { return mock.message }
func (mock *mockViolation) MessageTemplate() string                    { return mock.messageTemplate }
func (mock *mockViolation) Parameters() []validation.TemplateParameter { return mock.parameters }
func (mock *mockViolation) PropertyPath() *validation.PropertyPath     { return mock.propertyPath }

func mockNewViolationFunc() validation.ViolationFactory {
	return validation.NewViolationFunc(func(
		err error,
		messageTemplate string,
		pluralCount int,
		parameters []validation.TemplateParameter,
		propertyPath *validation.PropertyPath,
		lang language.Tag,
	) validation.Violation {
		return &mockViolation{
			err:             err,
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

type mockTranslator struct {
	translate func(tag language.Tag, message string, pluralCount int) string
}

func (m mockTranslator) Translate(tag language.Tag, message string, pluralCount int) string {
	return m.translate(tag, message, pluralCount)
}

type asyncConstraint func(value *string, scope validation.Scope) error

func (f asyncConstraint) ValidateString(value *string, scope validation.Scope) error {
	return f(value, scope)
}
