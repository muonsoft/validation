package it

import (
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/code"
	"github.com/muonsoft/validation/generic"
	"github.com/muonsoft/validation/message"

	"time"
)

// NotBlankConstraint checks that a value is not blank: not equal to zero, an empty string, an empty
// slice/array, an empty map, false or nil. Nil behavior is configurable via AllowNil() method.
// To check that a value is not nil only use NotNilConstraint.
type NotBlankConstraint struct {
	isIgnored         bool
	allowNil          bool
	code              string
	messageTemplate   string
	messageParameters validation.TemplateParameterList
}

// IsNotBlank creates a NotBlankConstraint for checking that value is not empty.
//
// Example
//  s := ""
//  err := validator.ValidateString(&s, it.IsNotBlank())
func IsNotBlank() NotBlankConstraint {
	return NotBlankConstraint{
		code:            code.NotBlank,
		messageTemplate: message.NotBlank,
	}
}

// SetUp always returns no error.
func (c NotBlankConstraint) SetUp() error {
	return nil
}

// Name is the constraint name.
func (c NotBlankConstraint) Name() string {
	return "NotBlankConstraint"
}

// AllowNil makes nil values valid.
func (c NotBlankConstraint) AllowNil() NotBlankConstraint {
	c.allowNil = true
	return c
}

// When enables conditional validation of this constraint. If the expression evaluates to false,
// then the constraint will be ignored.
func (c NotBlankConstraint) When(condition bool) NotBlankConstraint {
	c.isIgnored = !condition
	return c
}

// Code overrides default code for produced violation.
func (c NotBlankConstraint) Code(code string) NotBlankConstraint {
	c.code = code
	return c
}

// Message sets the violation message template. You can set custom template parameters
// for injecting its values into the final message.
func (c NotBlankConstraint) Message(template string, parameters ...validation.TemplateParameter) NotBlankConstraint {
	c.messageTemplate = template
	c.messageParameters = parameters
	return c
}

func (c NotBlankConstraint) ValidateNil(scope validation.Scope) error {
	if c.isIgnored || c.allowNil {
		return nil
	}

	return c.newViolation(scope)
}

func (c NotBlankConstraint) ValidateBool(value *bool, scope validation.Scope) error {
	if c.isIgnored {
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

func (c NotBlankConstraint) ValidateNumber(value generic.Number, scope validation.Scope) error {
	if c.isIgnored {
		return nil
	}
	if c.allowNil && value.IsNil() {
		return nil
	}
	if !value.IsNil() && !value.IsZero() {
		return nil
	}

	return c.newViolation(scope)
}

func (c NotBlankConstraint) ValidateString(value *string, scope validation.Scope) error {
	if c.isIgnored {
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

func (c NotBlankConstraint) ValidateIterable(value generic.Iterable, scope validation.Scope) error {
	if c.isIgnored {
		return nil
	}
	if c.allowNil && value.IsNil() {
		return nil
	}
	if value.Count() > 0 {
		return nil
	}

	return c.newViolation(scope)
}

func (c NotBlankConstraint) ValidateCountable(count int, scope validation.Scope) error {
	if c.isIgnored || count > 0 {
		return nil
	}

	return c.newViolation(scope)
}

func (c NotBlankConstraint) ValidateTime(value *time.Time, scope validation.Scope) error {
	if c.isIgnored {
		return nil
	}
	if c.allowNil && value == nil {
		return nil
	}

	var empty time.Time
	if value != nil && *value != empty {
		return nil
	}

	return c.newViolation(scope)
}

func (c NotBlankConstraint) newViolation(scope validation.Scope) validation.Violation {
	return scope.BuildViolation(c.code, c.messageTemplate).
		SetParameters(c.messageParameters...).
		CreateViolation()
}

// BlankConstraint checks that a value is blank: equal to false, nil, zero, an empty string, an empty
// slice, array, or a map.
type BlankConstraint struct {
	isIgnored         bool
	code              string
	messageTemplate   string
	messageParameters validation.TemplateParameterList
}

// IsBlank creates a BlankConstraint for checking that value is empty.
//
// Example
//  s := "foo"
//  err := validator.ValidateString(&s, it.IsBlank())
func IsBlank() BlankConstraint {
	return BlankConstraint{
		code:            code.Blank,
		messageTemplate: message.Blank,
	}
}

// SetUp always returns no error.
func (c BlankConstraint) SetUp() error {
	return nil
}

// Name is the constraint name.
func (c BlankConstraint) Name() string {
	return "BlankConstraint"
}

// When enables conditional validation of this constraint. If the expression evaluates to false,
// then the constraint will be ignored.
func (c BlankConstraint) When(condition bool) BlankConstraint {
	c.isIgnored = !condition
	return c
}

// Code overrides default code for produced violation.
func (c BlankConstraint) Code(code string) BlankConstraint {
	c.code = code
	return c
}

// Message sets the violation message template. You can set custom template parameters
// for injecting its values into the final message.
func (c BlankConstraint) Message(template string, parameters ...validation.TemplateParameter) BlankConstraint {
	c.messageTemplate = template
	c.messageParameters = parameters
	return c
}

func (c BlankConstraint) ValidateNil(scope validation.Scope) error {
	return nil
}

func (c BlankConstraint) ValidateBool(value *bool, scope validation.Scope) error {
	if c.isIgnored || value == nil || !*value {
		return nil
	}

	return c.newViolation(scope)
}

func (c BlankConstraint) ValidateNumber(value generic.Number, scope validation.Scope) error {
	if c.isIgnored || value.IsNil() || value.IsZero() {
		return nil
	}

	return c.newViolation(scope)
}

func (c BlankConstraint) ValidateString(value *string, scope validation.Scope) error {
	if c.isIgnored || value == nil || *value == "" {
		return nil
	}

	return c.newViolation(scope)
}

func (c BlankConstraint) ValidateIterable(value generic.Iterable, scope validation.Scope) error {
	if c.isIgnored || value.Count() == 0 {
		return nil
	}

	return c.newViolation(scope)
}

func (c BlankConstraint) ValidateCountable(count int, scope validation.Scope) error {
	if c.isIgnored || count == 0 {
		return nil
	}

	return c.newViolation(scope)
}

func (c BlankConstraint) ValidateTime(value *time.Time, scope validation.Scope) error {
	var empty time.Time
	if c.isIgnored || value == nil || *value == empty {
		return nil
	}

	return c.newViolation(scope)
}

func (c BlankConstraint) newViolation(scope validation.Scope) validation.Violation {
	return scope.BuildViolation(c.code, c.messageTemplate).
		SetParameters(c.messageParameters...).
		CreateViolation()
}

// NotNilConstraint checks that a value in not strictly equal to nil. To check that values in not blank use
// NotBlankConstraint.
type NotNilConstraint struct {
	isIgnored         bool
	code              string
	messageTemplate   string
	messageParameters validation.TemplateParameterList
}

// IsNotNil creates a NotNilConstraint to check that a value is not strictly equal to nil.
//
// Example
//  var s *string
//  err := validator.ValidateString(s, it.IsNotNil())
func IsNotNil() NotNilConstraint {
	return NotNilConstraint{
		code:            code.NotNil,
		messageTemplate: message.NotNil,
	}
}

// SetUp always returns no error.
func (c NotNilConstraint) SetUp() error {
	return nil
}

// Name is the constraint name.
func (c NotNilConstraint) Name() string {
	return "NotNilConstraint"
}

// When enables conditional validation of this constraint. If the expression evaluates to false,
// then the constraint will be ignored.
func (c NotNilConstraint) When(condition bool) NotNilConstraint {
	c.isIgnored = !condition
	return c
}

// Code overrides default code for produced violation.
func (c NotNilConstraint) Code(code string) NotNilConstraint {
	c.code = code
	return c
}

// Message sets the violation message template. You can set custom template parameters
// for injecting its values into the final message.
func (c NotNilConstraint) Message(template string, parameters ...validation.TemplateParameter) NotNilConstraint {
	c.messageTemplate = template
	c.messageParameters = parameters
	return c
}

func (c NotNilConstraint) ValidateNil(scope validation.Scope) error {
	if c.isIgnored {
		return nil
	}

	return c.newViolation(scope)
}

func (c NotNilConstraint) ValidateNumber(value generic.Number, scope validation.Scope) error {
	if c.isIgnored {
		return nil
	}
	if !value.IsNil() {
		return nil
	}

	return c.newViolation(scope)
}

func (c NotNilConstraint) ValidateString(value *string, scope validation.Scope) error {
	if c.isIgnored {
		return nil
	}
	if value != nil {
		return nil
	}

	return c.newViolation(scope)
}

func (c NotNilConstraint) ValidateTime(value *time.Time, scope validation.Scope) error {
	if c.isIgnored {
		return nil
	}
	if value != nil {
		return nil
	}

	return c.newViolation(scope)
}

func (c NotNilConstraint) ValidateIterable(value generic.Iterable, scope validation.Scope) error {
	if c.isIgnored {
		return nil
	}
	if !value.IsNil() {
		return nil
	}

	return c.newViolation(scope)
}

func (c NotNilConstraint) newViolation(scope validation.Scope) validation.Violation {
	return scope.BuildViolation(c.code, c.messageTemplate).
		SetParameters(c.messageParameters...).
		CreateViolation()
}

// NilConstraint checks that a value in strictly equal to nil. To check that values in blank use
// BlankConstraint.
type NilConstraint struct {
	isIgnored         bool
	code              string
	messageTemplate   string
	messageParameters validation.TemplateParameterList
}

// IsNil creates a NilConstraint to check that a value is strictly equal to nil.
//
// Example
//  var s *string
//  err := validator.ValidateString(s, it.IsNil())
func IsNil() NilConstraint {
	return NilConstraint{
		code:            code.Nil,
		messageTemplate: message.Nil,
	}
}

// SetUp always returns no error.
func (c NilConstraint) SetUp() error {
	return nil
}

// Name is the constraint name.
func (c NilConstraint) Name() string {
	return "NilConstraint"
}

// When enables conditional validation of this constraint. If the expression evaluates to false,
// then the constraint will be ignored.
func (c NilConstraint) When(condition bool) NilConstraint {
	c.isIgnored = !condition
	return c
}

// Code overrides default code for produced violation.
func (c NilConstraint) Code(code string) NilConstraint {
	c.code = code
	return c
}

// Message sets the violation message template. You can set custom template parameters
// for injecting its values into the final message.
func (c NilConstraint) Message(template string, parameters ...validation.TemplateParameter) NilConstraint {
	c.messageTemplate = template
	c.messageParameters = parameters
	return c
}

func (c NilConstraint) ValidateNil(scope validation.Scope) error {
	return nil
}

func (c NilConstraint) ValidateNumber(value generic.Number, scope validation.Scope) error {
	if c.isIgnored || value.IsNil() {
		return nil
	}

	return c.newViolation(scope)
}

func (c NilConstraint) ValidateString(value *string, scope validation.Scope) error {
	if c.isIgnored || value == nil {
		return nil
	}

	return c.newViolation(scope)
}

func (c NilConstraint) ValidateTime(value *time.Time, scope validation.Scope) error {
	if c.isIgnored || value == nil {
		return nil
	}

	return c.newViolation(scope)
}

func (c NilConstraint) ValidateIterable(value generic.Iterable, scope validation.Scope) error {
	if c.isIgnored || value.IsNil() {
		return nil
	}

	return c.newViolation(scope)
}

func (c NilConstraint) newViolation(scope validation.Scope) validation.Violation {
	return scope.BuildViolation(c.code, c.messageTemplate).
		SetParameters(c.messageParameters...).
		CreateViolation()
}

// BoolConstraint checks that a bool value in strictly equal to expected bool value.
type BoolConstraint struct {
	isIgnored         bool
	expected          bool
	code              string
	messageTemplate   string
	messageParameters validation.TemplateParameterList
}

// IsTrue creates a BoolConstraint to check that a value is not strictly equal to true.
//
// Example
//  var b *bool
//  err := validator.ValidateBool(b, it.IsTrue())
func IsTrue() BoolConstraint {
	return BoolConstraint{
		expected:        true,
		code:            code.True,
		messageTemplate: message.True,
	}
}

// IsFalse creates a BoolConstraint to check that a value is not strictly equal to false.
//
// Example
//  var b *bool
//  err := validator.ValidateBool(b, it.IsFalse())
func IsFalse() BoolConstraint {
	return BoolConstraint{
		expected:        false,
		code:            code.False,
		messageTemplate: message.False,
	}
}

// SetUp always returns no error.
func (c BoolConstraint) SetUp() error {
	return nil
}

// Name is the constraint name.
func (c BoolConstraint) Name() string {
	return "BoolConstraint"
}

// When enables conditional validation of this constraint. If the expression evaluates to false,
// then the constraint will be ignored.
func (c BoolConstraint) When(condition bool) BoolConstraint {
	c.isIgnored = !condition
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
	if c.isIgnored || value == nil || *value == c.expected {
		return nil
	}

	return scope.BuildViolation(c.code, c.messageTemplate).
		SetParameters(c.messageParameters...).
		CreateViolation()
}
