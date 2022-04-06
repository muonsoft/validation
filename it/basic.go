package it

import (
	"time"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/code"
	"github.com/muonsoft/validation/generic"
	"github.com/muonsoft/validation/message"
)

// NotBlankConstraint checks that a value is not blank: not equal to zero, an empty string, an empty
// slice/array, an empty map, false or nil. Nil behavior is configurable via AllowNil() method.
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
func IsNotBlank[T comparable]() NotBlankConstraint[T] {
	return NotBlankConstraint[T]{
		code:            code.NotBlank,
		messageTemplate: message.Templates[code.NotBlank],
	}
}

// SetUp always returns no error.
func (c NotBlankConstraint[T]) SetUp() error {
	return nil
}

// Name is the constraint name.
func (c NotBlankConstraint[T]) Name() string {
	return "NotBlankConstraint"
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
func (c NotBlankConstraint) WhenGroups(groups ...string) NotBlankConstraint {
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

func (c NotBlankConstraint[T]) ValidateNil(scope validation.Scope) error {
	if c.isIgnored || c.allowNil || scope.IsIgnored(c.groups...) {
		return nil
	}

	return c.newViolation(scope)
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

func (c NotBlankConstraint[T]) ValidateStrings(values []string, scope validation.Scope) error {
	if c.isIgnored || scope.IsIgnored(c.groups...) {
		return nil
	}
	if c.allowNil && values == nil {
		return nil
	}
	if len(values) > 0 {
		return nil
	}

	return c.newViolation(scope)
}

func (c NotBlankConstraint[T]) ValidateIterable(value generic.Iterable, scope validation.Scope) error {
	if c.isIgnored || scope.IsIgnored(c.groups...) {
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
type BlankConstraint struct {
	isIgnored         bool
	groups            []string
	code              string
	messageTemplate   string
	messageParameters validation.TemplateParameterList
}

// IsBlank creates a BlankConstraint for checking that value is empty.
func IsBlank() BlankConstraint {
	return BlankConstraint{
		code:            code.Blank,
		messageTemplate: message.Templates[code.Blank],
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

// WhenGroups enables conditional validation of the constraint by using the validation groups.
func (c BlankConstraint) WhenGroups(groups ...string) BlankConstraint {
	c.groups = groups
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
	if c.isIgnored || scope.IsIgnored(c.groups...) || value == nil || !*value {
		return nil
	}

	return c.newViolation(scope)
}

func (c BlankConstraint) ValidateNumber(value generic.Number, scope validation.Scope) error {
	if c.isIgnored || scope.IsIgnored(c.groups...) || value.IsNil() || value.IsZero() {
		return nil
	}

	return c.newViolation(scope)
}

func (c BlankConstraint) ValidateString(value *string, scope validation.Scope) error {
	if c.isIgnored || scope.IsIgnored(c.groups...) || value == nil || *value == "" {
		return nil
	}

	return c.newViolation(scope)
}

func (c BlankConstraint) ValidateStrings(values []string, scope validation.Scope) error {
	if c.isIgnored || scope.IsIgnored(c.groups...) || len(values) == 0 {
		return nil
	}

	return c.newViolation(scope)
}

func (c BlankConstraint) ValidateIterable(value generic.Iterable, scope validation.Scope) error {
	if c.isIgnored || scope.IsIgnored(c.groups...) || value.Count() == 0 {
		return nil
	}

	return c.newViolation(scope)
}

func (c BlankConstraint) ValidateCountable(count int, scope validation.Scope) error {
	if c.isIgnored || scope.IsIgnored(c.groups...) || count == 0 {
		return nil
	}

	return c.newViolation(scope)
}

func (c BlankConstraint) ValidateTime(value *time.Time, scope validation.Scope) error {
	if c.isIgnored || scope.IsIgnored(c.groups...) || value == nil || value.IsZero() {
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
	groups            []string
	code              string
	messageTemplate   string
	messageParameters validation.TemplateParameterList
}

// IsNotNil creates a NotNilConstraint to check that a value is not strictly equal to nil.
func IsNotNil() NotNilConstraint {
	return NotNilConstraint{
		code:            code.NotNil,
		messageTemplate: message.Templates[code.NotNil],
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

// WhenGroups enables conditional validation of the constraint by using the validation groups.
func (c NotNilConstraint) WhenGroups(groups ...string) NotNilConstraint {
	c.groups = groups
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
	if c.isIgnored || scope.IsIgnored(c.groups...) {
		return nil
	}

	return c.newViolation(scope)
}

func (c NotNilConstraint) ValidateBool(value *bool, scope validation.Scope) error {
	if c.isIgnored || scope.IsIgnored(c.groups...) || value != nil {
		return nil
	}

	return c.newViolation(scope)
}

func (c NotNilConstraint) ValidateNumber(value generic.Number, scope validation.Scope) error {
	if c.isIgnored || scope.IsIgnored(c.groups...) || !value.IsNil() {
		return nil
	}

	return c.newViolation(scope)
}

func (c NotNilConstraint) ValidateString(value *string, scope validation.Scope) error {
	if c.isIgnored || scope.IsIgnored(c.groups...) || value != nil {
		return nil
	}

	return c.newViolation(scope)
}

func (c NotNilConstraint) ValidateStrings(values []string, scope validation.Scope) error {
	if c.isIgnored || scope.IsIgnored(c.groups...) || values != nil {
		return nil
	}

	return c.newViolation(scope)
}

func (c NotNilConstraint) ValidateTime(value *time.Time, scope validation.Scope) error {
	if c.isIgnored || scope.IsIgnored(c.groups...) || value != nil {
		return nil
	}

	return c.newViolation(scope)
}

func (c NotNilConstraint) ValidateIterable(value generic.Iterable, scope validation.Scope) error {
	if c.isIgnored || scope.IsIgnored(c.groups...) || !value.IsNil() {
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
	groups            []string
	code              string
	messageTemplate   string
	messageParameters validation.TemplateParameterList
}

// IsNil creates a NilConstraint to check that a value is strictly equal to nil.
func IsNil() NilConstraint {
	return NilConstraint{
		code:            code.Nil,
		messageTemplate: message.Templates[code.Nil],
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

// WhenGroups enables conditional validation of the constraint by using the validation groups.
func (c NilConstraint) WhenGroups(groups ...string) NilConstraint {
	c.groups = groups
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

func (c NilConstraint) ValidateBool(value *bool, scope validation.Scope) error {
	if c.isIgnored || scope.IsIgnored(c.groups...) || value == nil {
		return nil
	}

	return c.newViolation(scope)
}

func (c NilConstraint) ValidateNumber(value generic.Number, scope validation.Scope) error {
	if c.isIgnored || scope.IsIgnored(c.groups...) || value.IsNil() {
		return nil
	}

	return c.newViolation(scope)
}

func (c NilConstraint) ValidateString(value *string, scope validation.Scope) error {
	if c.isIgnored || scope.IsIgnored(c.groups...) || value == nil {
		return nil
	}

	return c.newViolation(scope)
}

func (c NilConstraint) ValidateStrings(values []string, scope validation.Scope) error {
	if c.isIgnored || scope.IsIgnored(c.groups...) || values == nil {
		return nil
	}

	return c.newViolation(scope)
}

func (c NilConstraint) ValidateTime(value *time.Time, scope validation.Scope) error {
	if c.isIgnored || scope.IsIgnored(c.groups...) || value == nil {
		return nil
	}

	return c.newViolation(scope)
}

func (c NilConstraint) ValidateIterable(value generic.Iterable, scope validation.Scope) error {
	if c.isIgnored || scope.IsIgnored(c.groups...) || value.IsNil() {
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
