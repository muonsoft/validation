package validation

import (
	"errors"
	"strings"
)

type Violation interface {
	error
	GetCode() string
	GetMessage() string
	GetMessageTemplate() string
	GetParameters() map[string]string
	GetPropertyPath() PropertyPath
}

type ViolationList []Violation

func (violations ViolationList) Error() string {
	var s strings.Builder
	s.Grow(32 * len(violations))

	for i, v := range violations {
		if i > 0 {
			s.WriteString("; ")
		}
		if iv, ok := v.(internalViolation); ok {
			iv.writeToBuilder(&s)
		} else {
			s.WriteString(v.Error())
		}
	}

	return s.String()
}

func NewViolation(
	code,
	messageTemplate string,
	parameters map[string]string,
	propertyPath PropertyPath,
) Violation {
	return &internalViolation{
		code:            code,
		message:         renderMessage(messageTemplate, parameters),
		messageTemplate: messageTemplate,
		parameters:      parameters,
		propertyPath:    propertyPath,
	}
}

type internalViolation struct {
	code            string
	message         string
	messageTemplate string
	parameters      map[string]string
	propertyPath    PropertyPath
}

func (v internalViolation) Error() string {
	var s strings.Builder
	s.Grow(32)
	v.writeToBuilder(&s)

	return s.String()
}

func (v internalViolation) writeToBuilder(s *strings.Builder) {
	s.WriteString("violation")
	if len(v.propertyPath) > 0 {
		s.WriteString(" at '" + v.propertyPath.Format() + "'")
	}
	s.WriteString(": " + v.message)
}

func (v internalViolation) GetCode() string {
	return v.code
}

func (v internalViolation) GetMessage() string {
	return v.message
}

func (v internalViolation) GetMessageTemplate() string {
	return v.messageTemplate
}

func (v internalViolation) GetParameters() map[string]string {
	return v.parameters
}

func (v internalViolation) GetPropertyPath() PropertyPath {
	return v.propertyPath
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
