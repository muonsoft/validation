package validation

import (
	"errors"
	"strings"

	"golang.org/x/text/language"
)

type Violation interface {
	error
	GetCode() string
	GetMessage() string
	GetMessageTemplate() string
	GetParameters() map[string]string
	GetPropertyPath() PropertyPath
}

type ViolationFactory interface {
	BuildViolation(code, message string) *ViolationBuilder
	CreateViolation(
		code,
		messageTemplate string,
		pluralCount int,
		parameters map[string]string,
		propertyPath PropertyPath,
		lang language.Tag,
	) Violation
}

type ViolationList []Violation

type NewViolationFunc func(
	code,
	messageTemplate string,
	pluralCount int,
	parameters map[string]string,
	propertyPath PropertyPath,
	lang language.Tag,
) Violation

func (f NewViolationFunc) BuildViolation(code, message string) *ViolationBuilder {
	return newViolationBuilder(f, code, message)
}

func (f NewViolationFunc) CreateViolation(
	code,
	messageTemplate string,
	pluralCount int,
	parameters map[string]string,
	propertyPath PropertyPath,
	lang language.Tag,
) Violation {
	return f(code, messageTemplate, pluralCount, parameters, propertyPath, lang)
}

func (violations ViolationList) Error() string {
	var s strings.Builder
	s.Grow(32 * len(violations))

	for i, v := range violations {
		if i > 0 {
			s.WriteString("; ")
		}
		if iv, ok := v.(*internalViolation); ok {
			iv.writeToBuilder(&s)
		} else {
			s.WriteString(v.Error())
		}
	}

	return s.String()
}

func (violations *ViolationList) AddFromError(err error) error {
	if violation, ok := UnwrapViolation(err); ok {
		*violations = append(*violations, violation)
	} else if violationList, ok := UnwrapViolationList(err); ok {
		*violations = append(*violations, violationList...)
	} else if err != nil {
		return err
	}

	return nil
}

func (violations ViolationList) AsError() error {
	if len(violations) == 0 {
		return nil
	}

	return violations
}

func IsViolation(err error) bool {
	var violation Violation

	return errors.As(err, &violation)
}

func IsViolationList(err error) bool {
	var violations ViolationList

	return errors.As(err, &violations)
}

func UnwrapViolation(err error) (Violation, bool) {
	var violation Violation

	as := errors.As(err, &violation)

	return violation, as
}

func UnwrapViolationList(err error) (ViolationList, bool) {
	var violation ViolationList

	as := errors.As(err, &violation)

	return violation, as
}

type internalViolation struct {
	Code            string            `json:"code"`
	Message         string            `json:"message"`
	MessageTemplate string            `json:"-"`
	Parameters      map[string]string `json:"-"`
	PropertyPath    PropertyPath      `json:"propertyPath,omitempty"`
}

func (v internalViolation) Error() string {
	var s strings.Builder
	s.Grow(32)
	v.writeToBuilder(&s)

	return s.String()
}

func (v internalViolation) writeToBuilder(s *strings.Builder) {
	s.WriteString("violation")
	if len(v.PropertyPath) > 0 {
		s.WriteString(" at '" + v.PropertyPath.Format() + "'")
	}
	s.WriteString(": " + v.Message)
}

func (v internalViolation) GetCode() string {
	return v.Code
}

func (v internalViolation) GetMessage() string {
	return v.Message
}

func (v internalViolation) GetMessageTemplate() string {
	return v.MessageTemplate
}

func (v internalViolation) GetParameters() map[string]string {
	return v.Parameters
}

func (v internalViolation) GetPropertyPath() PropertyPath {
	return v.PropertyPath
}

type internalViolationFactory struct {
	translator *Translator
}

func newViolationFactory(translator *Translator) *internalViolationFactory {
	return &internalViolationFactory{translator: translator}
}

func (factory *internalViolationFactory) CreateViolation(
	code,
	messageTemplate string,
	pluralCount int,
	parameters map[string]string,
	propertyPath PropertyPath,
	lang language.Tag,
) Violation {
	message := factory.translator.translate(lang, messageTemplate, pluralCount)

	return &internalViolation{
		Code:            code,
		Message:         renderMessage(message, parameters),
		MessageTemplate: messageTemplate,
		Parameters:      parameters,
		PropertyPath:    propertyPath,
	}
}

func (factory *internalViolationFactory) BuildViolation(code, message string) *ViolationBuilder {
	return newViolationBuilder(factory, code, message)
}

type ViolationBuilder struct {
	code            string
	messageTemplate string
	pluralCount     int
	parameters      map[string]string
	propertyPath    PropertyPath
	language        language.Tag

	violationFactory ViolationFactory
}

func newViolationBuilder(factory ViolationFactory, code, message string) *ViolationBuilder {
	return &ViolationBuilder{
		code:             code,
		messageTemplate:  message,
		violationFactory: factory,
	}
}

func (b *ViolationBuilder) SetParameters(parameters map[string]string) *ViolationBuilder {
	b.parameters = parameters

	return b
}

func (b *ViolationBuilder) SetParameter(name, value string) *ViolationBuilder {
	if b.parameters == nil {
		b.parameters = make(map[string]string)
	}
	b.parameters[name] = value

	return b
}

func (b *ViolationBuilder) SetPropertyPath(path PropertyPath) *ViolationBuilder {
	b.propertyPath = path

	return b
}

func (b *ViolationBuilder) SetPluralCount(pluralCount int) *ViolationBuilder {
	b.pluralCount = pluralCount

	return b
}

func (b *ViolationBuilder) SetLanguage(tag language.Tag) *ViolationBuilder {
	b.language = tag

	return b
}

func (b *ViolationBuilder) GetViolation() Violation {
	return b.violationFactory.CreateViolation(
		b.code,
		b.messageTemplate,
		b.pluralCount,
		b.parameters,
		b.propertyPath,
		b.language,
	)
}
