package validation

import (
	"time"

	"github.com/muonsoft/validation/code"
	"github.com/muonsoft/validation/message"
)

// Numeric is used as a type parameter for numeric values.
type Numeric interface {
	~float32 | ~float64 |
		~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

// NilConstraint is used for a special cases to check a value for nil.
type NilConstraint interface {
	ValidateNil(isNil bool, scope Scope) error
}

// BoolConstraint is used to build constraints for boolean values validation.
type BoolConstraint interface {
	ValidateBool(value *bool, scope Scope) error
}

// NumberConstraint is used to build constraints for numeric values validation.
type NumberConstraint[T Numeric] interface {
	ValidateNumber(value *T, scope Scope) error
}

// StringConstraint is used to build constraints for string values validation.
type StringConstraint interface {
	ValidateString(value *string, scope Scope) error
}

// ComparableConstraint is used to build constraints for generic comparable value validation.
type ComparableConstraint[T comparable] interface {
	ValidateComparable(value *T, scope Scope) error
}

// ComparablesConstraint is used to build constraints for generic comparable values validation.
type ComparablesConstraint[T comparable] interface {
	ValidateComparables(values []T, scope Scope) error
}

// CountableConstraint is used to build constraints for simpler validation of iterable elements count.
type CountableConstraint interface {
	ValidateCountable(count int, scope Scope) error
}

// TimeConstraint is used to build constraints for date/time validation.
type TimeConstraint interface {
	ValidateTime(value *time.Time, scope Scope) error
}

// CustomStringConstraint can be used to create custom constraints for validating string values
// based on function with signature func(string) bool.
type CustomStringConstraint struct {
	isIgnored         bool
	isValid           func(string) bool
	groups            []string
	code              string
	messageTemplate   string
	messageParameters TemplateParameterList
}

// NewCustomStringConstraint creates a new string constraint from a function with signature func(string) bool.
// Optional parameters can be used to set up violation code (first), message template (second).
// All other parameters are ignored.
func NewCustomStringConstraint(isValid func(string) bool, parameters ...string) CustomStringConstraint {
	constraint := CustomStringConstraint{
		isValid:         isValid,
		code:            code.NotValid,
		messageTemplate: message.Templates[code.NotValid],
	}

	if len(parameters) > 0 {
		constraint.code = parameters[0]
	}
	if len(parameters) > 1 {
		constraint.messageTemplate = parameters[1]
	}

	return constraint
}

// Code overrides default code for produced violation.
func (c CustomStringConstraint) Code(code string) CustomStringConstraint {
	c.code = code
	return c
}

// Message sets the violation message template. You can set custom template parameters
// for injecting its values into the final message. Also, you can use default parameters:
//
//	{{ value }} - the current (invalid) value.
func (c CustomStringConstraint) Message(template string, parameters ...TemplateParameter) CustomStringConstraint {
	c.messageTemplate = template
	c.messageParameters = parameters
	return c
}

// When enables conditional validation of this constraint. If the expression evaluates to false,
// then the constraint will be ignored.
func (c CustomStringConstraint) When(condition bool) CustomStringConstraint {
	c.isIgnored = !condition
	return c
}

// WhenGroups enables conditional validation of the constraint by using the validation groups.
func (c CustomStringConstraint) WhenGroups(groups ...string) CustomStringConstraint {
	c.groups = groups
	return c
}

func (c CustomStringConstraint) ValidateString(value *string, scope Scope) error {
	if c.isIgnored || scope.IsIgnored(c.groups...) || value == nil || *value == "" || c.isValid(*value) {
		return nil
	}

	return scope.BuildViolation(c.code, c.messageTemplate).
		WithParameters(
			c.messageParameters.Prepend(
				TemplateParameter{Key: "{{ value }}", Value: *value},
			)...,
		).
		WithParameter("{{ value }}", *value).
		Create()
}
