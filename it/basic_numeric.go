package it

import (
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/code"
	"github.com/muonsoft/validation/message"
)

// NotBlankNumberConstraint checks that a numeric value is not blank: not equal to zero or nil.
// Nil behavior is configurable via AllowNil() method.
// To check that a value is not nil only use NotNilNumberConstraint.
type NotBlankNumberConstraint[T validation.Numeric] struct {
	blank             T
	isIgnored         bool
	allowNil          bool
	groups            []string
	code              string
	messageTemplate   string
	messageParameters validation.TemplateParameterList
}

// IsNotBlankNumber creates a NotBlankNumberConstraint for checking that numeric value is not empty.
func IsNotBlankNumber[T validation.Numeric]() NotBlankNumberConstraint[T] {
	return NotBlankNumberConstraint[T]{
		code:            code.NotBlank,
		messageTemplate: message.Templates[code.NotBlank],
	}
}

// AllowNil makes nil values valid.
func (c NotBlankNumberConstraint[T]) AllowNil() NotBlankNumberConstraint[T] {
	c.allowNil = true
	return c
}

// When enables conditional validation of this constraint. If the expression evaluates to false,
// then the constraint will be ignored.
func (c NotBlankNumberConstraint[T]) When(condition bool) NotBlankNumberConstraint[T] {
	c.isIgnored = !condition
	return c
}

// WhenGroups enables conditional validation of the constraint by using the validation groups.
func (c NotBlankNumberConstraint[T]) WhenGroups(groups ...string) NotBlankNumberConstraint[T] {
	c.groups = groups
	return c
}

// Code overrides default code for produced violation.
func (c NotBlankNumberConstraint[T]) Code(code string) NotBlankNumberConstraint[T] {
	c.code = code
	return c
}

// Message sets the violation message template. You can set custom template parameters
// for injecting its values into the final message.
func (c NotBlankNumberConstraint[T]) Message(template string, parameters ...validation.TemplateParameter) NotBlankNumberConstraint[T] {
	c.messageTemplate = template
	c.messageParameters = parameters
	return c
}

func (c NotBlankNumberConstraint[T]) ValidateNumber(value *T, scope validation.Scope) error {
	if c.isIgnored || scope.IsIgnored(c.groups...) {
		return nil
	}
	if c.allowNil && value == nil {
		return nil
	}
	if value != nil && *value != c.blank {
		return nil
	}

	return c.newViolation(scope)
}

func (c NotBlankNumberConstraint[T]) newViolation(scope validation.Scope) validation.Violation {
	return scope.BuildViolation(c.code, c.messageTemplate).
		SetParameters(c.messageParameters...).
		CreateViolation()
}
