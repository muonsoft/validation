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

type NewViolationFunc func(
	code,
	messageTemplate string,
	parameters map[string]string,
	propertyPath PropertyPath,
) Violation

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

func NewViolation(
	code,
	messageTemplate string,
	parameters map[string]string,
	propertyPath PropertyPath,
) Violation {
	return &internalViolation{
		Code:            code,
		Message:         renderMessage(messageTemplate, parameters),
		MessageTemplate: messageTemplate,
		Parameters:      parameters,
		PropertyPath:    propertyPath,
	}
}

type internalViolation struct {
	Code            string
	Message         string
	MessageTemplate string
	Parameters      map[string]string
	PropertyPath    PropertyPath
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
