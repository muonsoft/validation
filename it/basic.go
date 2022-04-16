package it

import (
	"time"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/code"
	"github.com/muonsoft/validation/message"
)

// NotBlankConstraint checks that a value is not blank: an empty string, an empty countable (slice/array/map),
// an empty generic number, generic comparable, false or nil. Nil behavior is configurable via AllowNil() method.
// To check that a value is not nil only use NotNilConstraint.
type NotBlankConstraint[T comparable] struct {
	blank             T
	isIgnored         bool
	allowNil          bool
	groups            []string
	code              string
	messageTemplate   string
	messageParameters validation.TemplateParameterList
}

// IsNotBlank creates a NotBlankConstraint for checking that value is not empty.
func IsNotBlank() NotBlankConstraint[string] {
	return IsNotBlankComparable[string]()
}

// IsNotBlankNumber creates a NotBlankConstraint for checking that numeric value is not empty.
func IsNotBlankNumber[T validation.Numeric]() NotBlankConstraint[T] {
	return IsNotBlankComparable[T]()
}

// IsNotBlankComparable creates a NotBlankConstraint for checking that comparable value is not empty.
func IsNotBlankComparable[T comparable]() NotBlankConstraint[T] {
	return NotBlankConstraint[T]{
		code:            code.NotBlank,
		messageTemplate: message.Templates[code.NotBlank],
	}
}

// AllowNil makes nil values valid.
func (c NotBlankConstraint[T]) AllowNil() NotBlankConstraint[T] {
	c.allowNil = true
	return c
}

// When enables conditional validation of this constraint. If the expression evaluates to false,
// then the constraint will be ignored.
func (c NotBlankConstraint[T]) When(condition bool) NotBlankConstraint[T] {
	c.isIgnored = !condition
	return c
}

// WhenGroups enables conditional validation of the constraint by using the validation groups.
func (c NotBlankConstraint[T]) WhenGroups(groups ...string) NotBlankConstraint[T] {
	c.groups = groups
	return c
}

// Code overrides default code for produced violation.
func (c NotBlankConstraint[T]) Code(code string) NotBlankConstraint[T] {
	c.code = code
	return c
}

// Message sets the violation message template. You can set custom template parameters
// for injecting its values into the final message.
func (c NotBlankConstraint[T]) Message(template string, parameters ...validation.TemplateParameter) NotBlankConstraint[T] {
	c.messageTemplate = template
	c.messageParameters = parameters
	return c
}

func (c NotBlankConstraint[T]) ValidateBool(value *bool, scope validation.Scope) error {
	if c.isIgnored || scope.IsIgnored(c.groups...) {
		return nil
	}
	if c.allowNil && value == nil {
		return nil
	}
	if value != nil && *value {
		return nil
	}

	return c.newViolation(scope)
}

func (c NotBlankConstraint[T]) ValidateNumber(value *T, scope validation.Scope) error {
	return c.ValidateComparable(value, scope)
}

func (c NotBlankConstraint[T]) ValidateString(value *string, scope validation.Scope) error {
	if c.isIgnored || scope.IsIgnored(c.groups...) {
		return nil
	}
	if c.allowNil && value == nil {
		return nil
	}
	if value != nil && *value != "" {
		return nil
	}

	return c.newViolation(scope)
}

func (c NotBlankConstraint[T]) ValidateComparable(value *T, scope validation.Scope) error {
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

func (c NotBlankConstraint[T]) ValidateCountable(count int, scope validation.Scope) error {
	if c.isIgnored || scope.IsIgnored(c.groups...) || count > 0 {
		return nil
	}

	return c.newViolation(scope)
}

func (c NotBlankConstraint[T]) ValidateTime(value *time.Time, scope validation.Scope) error {
	if c.isIgnored || scope.IsIgnored(c.groups...) {
		return nil
	}
	if c.allowNil && value == nil {
		return nil
	}

	if value != nil && !value.IsZero() {
		return nil
	}

	return c.newViolation(scope)
}

func (c NotBlankConstraint[T]) newViolation(scope validation.Scope) validation.Violation {
	return scope.BuildViolation(c.code, c.messageTemplate).
		SetParameters(c.messageParameters...).
		CreateViolation()
}

// BlankConstraint checks that a value is blank: equal to false, nil, zero, an empty string, an empty
// slice, array, or a map.
type BlankConstraint[T comparable] struct {
	blank             T
	isIgnored         bool
	groups            []string
	code              string
	messageTemplate   string
	messageParameters validation.TemplateParameterList
}

// IsBlank creates a BlankConstraint for checking that value is empty.
func IsBlank() BlankConstraint[string] {
	return IsBlankComparable[string]()
}

// IsBlankNumber creates a BlankConstraint for checking that numeric value is nil or zero.
func IsBlankNumber[T validation.Numeric]() BlankConstraint[T] {
	return IsBlankComparable[T]()
}

// IsBlankComparable creates a BlankConstraint for checking that comparable value is not empty.
func IsBlankComparable[T comparable]() BlankConstraint[T] {
	return BlankConstraint[T]{
		code:            code.Blank,
		messageTemplate: message.Templates[code.Blank],
	}
}

// When enables conditional validation of this constraint. If the expression evaluates to false,
// then the constraint will be ignored.
func (c BlankConstraint[T]) When(condition bool) BlankConstraint[T] {
	c.isIgnored = !condition
	return c
}

// WhenGroups enables conditional validation of the constraint by using the validation groups.
func (c BlankConstraint[T]) WhenGroups(groups ...string) BlankConstraint[T] {
	c.groups = groups
	return c
}

// Code overrides default code for produced violation.
func (c BlankConstraint[T]) Code(code string) BlankConstraint[T] {
	c.code = code
	return c
}

// Message sets the violation message template. You can set custom template parameters
// for injecting its values into the final message.
func (c BlankConstraint[T]) Message(template string, parameters ...validation.TemplateParameter) BlankConstraint[T] {
	c.messageTemplate = template
	c.messageParameters = parameters
	return c
}

func (c BlankConstraint[T]) ValidateBool(value *bool, scope validation.Scope) error {
	if c.isIgnored || scope.IsIgnored(c.groups...) || value == nil || !*value {
		return nil
	}

	return c.newViolation(scope)
}

func (c BlankConstraint[T]) ValidateNumber(value *T, scope validation.Scope) error {
	return c.ValidateComparable(value, scope)
}

func (c BlankConstraint[T]) ValidateString(value *string, scope validation.Scope) error {
	if c.isIgnored || scope.IsIgnored(c.groups...) || value == nil || *value == "" {
		return nil
	}

	return c.newViolation(scope)
}

func (c BlankConstraint[T]) ValidateComparable(value *T, scope validation.Scope) error {
	if c.isIgnored || scope.IsIgnored(c.groups...) || value == nil || *value == c.blank {
		return nil
	}

	return c.newViolation(scope)
}

func (c BlankConstraint[T]) ValidateCountable(count int, scope validation.Scope) error {
	if c.isIgnored || scope.IsIgnored(c.groups...) || count == 0 {
		return nil
	}

	return c.newViolation(scope)
}

func (c BlankConstraint[T]) ValidateTime(value *time.Time, scope validation.Scope) error {
	if c.isIgnored || scope.IsIgnored(c.groups...) || value == nil || value.IsZero() {
		return nil
	}

	return c.newViolation(scope)
}

func (c BlankConstraint[T]) newViolation(scope validation.Scope) validation.Violation {
	return scope.BuildViolation(c.code, c.messageTemplate).
		SetParameters(c.messageParameters...).
		CreateViolation()
}

// NotNilConstraint checks that a value in not strictly equal to nil. To check that values in not blank use
// NotBlankConstraint.
type NotNilConstraint[T comparable] struct {
	isIgnored         bool
	groups            []string
	code              string
	messageTemplate   string
	messageParameters validation.TemplateParameterList
}

// IsNotNil creates a NotNilConstraint to check that a value is not strictly equal to nil.
func IsNotNil() NotNilConstraint[string] {
	return IsNotNilComparable[string]()
}

// IsNotNilNumber creates a NotNilConstraint to check that a numeric value is not strictly equal to nil.
func IsNotNilNumber[T validation.Numeric]() NotNilConstraint[T] {
	return IsNotNilComparable[T]()
}

// IsNotNilComparable creates a NotNilConstraint to check that a comparable value is not strictly equal to nil.
func IsNotNilComparable[T comparable]() NotNilConstraint[T] {
	return NotNilConstraint[T]{
		code:            code.NotNil,
		messageTemplate: message.Templates[code.NotNil],
	}
}

// When enables conditional validation of this constraint. If the expression evaluates to false,
// then the constraint will be ignored.
func (c NotNilConstraint[T]) When(condition bool) NotNilConstraint[T] {
	c.isIgnored = !condition
	return c
}

// WhenGroups enables conditional validation of the constraint by using the validation groups.
func (c NotNilConstraint[T]) WhenGroups(groups ...string) NotNilConstraint[T] {
	c.groups = groups
	return c
}

// Code overrides default code for produced violation.
func (c NotNilConstraint[T]) Code(code string) NotNilConstraint[T] {
	c.code = code
	return c
}

// Message sets the violation message template. You can set custom template parameters
// for injecting its values into the final message.
func (c NotNilConstraint[T]) Message(template string, parameters ...validation.TemplateParameter) NotNilConstraint[T] {
	c.messageTemplate = template
	c.messageParameters = parameters
	return c
}

func (c NotNilConstraint[T]) ValidateNil(isNil bool, scope validation.Scope) error {
	if c.isIgnored || scope.IsIgnored(c.groups...) || !isNil {
		return nil
	}

	return c.newViolation(scope)
}

func (c NotNilConstraint[T]) ValidateBool(value *bool, scope validation.Scope) error {
	return c.ValidateNil(value == nil, scope)
}

func (c NotNilConstraint[T]) ValidateNumber(value *T, scope validation.Scope) error {
	return c.ValidateNil(value == nil, scope)
}

func (c NotNilConstraint[T]) ValidateString(value *string, scope validation.Scope) error {
	return c.ValidateNil(value == nil, scope)
}

func (c NotNilConstraint[T]) ValidateComparable(value *T, scope validation.Scope) error {
	return c.ValidateNil(value == nil, scope)
}

func (c NotNilConstraint[T]) ValidateTime(value *time.Time, scope validation.Scope) error {
	return c.ValidateNil(value == nil, scope)
}

func (c NotNilConstraint[T]) newViolation(scope validation.Scope) validation.Violation {
	return scope.BuildViolation(c.code, c.messageTemplate).
		SetParameters(c.messageParameters...).
		CreateViolation()
}

// NilConstraint checks that a value in strictly equal to nil. To check that values in blank use
// BlankConstraint.
type NilConstraint[T comparable] struct {
	isIgnored         bool
	groups            []string
	code              string
	messageTemplate   string
	messageParameters validation.TemplateParameterList
}

// IsNil creates a NilConstraint to check that a value is strictly equal to nil.
func IsNil() NilConstraint[string] {
	return IsNilComparable[string]()
}

// IsNilNumber creates a NilConstraint to check that a numeric value is strictly equal to nil.
func IsNilNumber[T validation.Numeric]() NilConstraint[T] {
	return IsNilComparable[T]()
}

// IsNilComparable creates a NilConstraint to check that a comparable value is strictly equal to nil.
func IsNilComparable[T comparable]() NilConstraint[T] {
	return NilConstraint[T]{
		code:            code.Nil,
		messageTemplate: message.Templates[code.Nil],
	}
}

// When enables conditional validation of this constraint. If the expression evaluates to false,
// then the constraint will be ignored.
func (c NilConstraint[T]) When(condition bool) NilConstraint[T] {
	c.isIgnored = !condition
	return c
}

// WhenGroups enables conditional validation of the constraint by using the validation groups.
func (c NilConstraint[T]) WhenGroups(groups ...string) NilConstraint[T] {
	c.groups = groups
	return c
}

// Code overrides default code for produced violation.
func (c NilConstraint[T]) Code(code string) NilConstraint[T] {
	c.code = code
	return c
}

// Message sets the violation message template. You can set custom template parameters
// for injecting its values into the final message.
func (c NilConstraint[T]) Message(template string, parameters ...validation.TemplateParameter) NilConstraint[T] {
	c.messageTemplate = template
	c.messageParameters = parameters
	return c
}

func (c NilConstraint[T]) ValidateNil(isNil bool, scope validation.Scope) error {
	if c.isIgnored || scope.IsIgnored(c.groups...) || isNil {
		return nil
	}

	return c.newViolation(scope)
}

func (c NilConstraint[T]) ValidateBool(value *bool, scope validation.Scope) error {
	return c.ValidateNil(value == nil, scope)
}

func (c NilConstraint[T]) ValidateNumber(value *T, scope validation.Scope) error {
	return c.ValidateNil(value == nil, scope)
}

func (c NilConstraint[T]) ValidateString(value *string, scope validation.Scope) error {
	return c.ValidateNil(value == nil, scope)
}

func (c NilConstraint[T]) ValidateComparable(value *T, scope validation.Scope) error {
	return c.ValidateNil(value == nil, scope)
}

func (c NilConstraint[T]) ValidateTime(value *time.Time, scope validation.Scope) error {
	return c.ValidateNil(value == nil, scope)
}

func (c NilConstraint[T]) newViolation(scope validation.Scope) validation.Violation {
	return scope.BuildViolation(c.code, c.messageTemplate).
		SetParameters(c.messageParameters...).
		CreateViolation()
}

// BoolConstraint checks that a bool value in strictly equal to expected bool value.
type BoolConstraint struct {
	isIgnored         bool
	expected          bool
	groups            []string
	code              string
	messageTemplate   string
	messageParameters validation.TemplateParameterList
}

// IsTrue creates a BoolConstraint to check that a value is not strictly equal to true.
func IsTrue() BoolConstraint {
	return BoolConstraint{
		expected:        true,
		code:            code.True,
		messageTemplate: message.Templates[code.True],
	}
}

// IsFalse creates a BoolConstraint to check that a value is not strictly equal to false.
func IsFalse() BoolConstraint {
	return BoolConstraint{
		expected:        false,
		code:            code.False,
		messageTemplate: message.Templates[code.False],
	}
}

// When enables conditional validation of this constraint. If the expression evaluates to false,
// then the constraint will be ignored.
func (c BoolConstraint) When(condition bool) BoolConstraint {
	c.isIgnored = !condition
	return c
}

// WhenGroups enables conditional validation of the constraint by using the validation groups.
func (c BoolConstraint) WhenGroups(groups ...string) BoolConstraint {
	c.groups = groups
	return c
}

// Code overrides default code for produced violation.
func (c BoolConstraint) Code(code string) BoolConstraint {
	c.code = code
	return c
}

// Message sets the violation message template. You can set custom template parameters
// for injecting its values into the final message.
func (c BoolConstraint) Message(template string, parameters ...validation.TemplateParameter) BoolConstraint {
	c.messageTemplate = template
	c.messageParameters = parameters
	return c
}

func (c BoolConstraint) ValidateBool(value *bool, scope validation.Scope) error {
	if c.isIgnored || scope.IsIgnored(c.groups...) || value == nil || *value == c.expected {
		return nil
	}

	return scope.BuildViolation(c.code, c.messageTemplate).
		SetParameters(c.messageParameters...).
		CreateViolation()
}
