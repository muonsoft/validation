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
	messageTemplate string
	isIgnored       bool
	allowNil        bool
}

// IsNotBlank creates a NotBlankConstraint for checking that value is not empty.
//
// Example
//  s := ""
//  err := validator.ValidateString(&s, it.IsNotBlank())
func IsNotBlank() NotBlankConstraint {
	return NotBlankConstraint{
		messageTemplate: message.NotBlank,
	}
}

// Name is the constraint name.
func (c NotBlankConstraint) Name() string {
	return "NotBlankConstraint"
}

// SetUp always returns no error.
func (c NotBlankConstraint) SetUp() error {
	return nil
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

// Message sets the violation message template.
func (c NotBlankConstraint) Message(message string) NotBlankConstraint {
	c.messageTemplate = message
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
	return scope.BuildViolation(code.NotBlank, c.messageTemplate).CreateViolation()
}

// BlankConstraint checks that a value is blank: equal to false, nil, zero, an empty string, an empty
// slice, array, or a map.
type BlankConstraint struct {
	messageTemplate string
	isIgnored       bool
}

// IsBlank creates a BlankConstraint for checking that value is empty.
//
// Example
//  s := "foo"
//  err := validator.ValidateString(&s, it.IsBlank())
func IsBlank() BlankConstraint {
	return BlankConstraint{
		messageTemplate: message.Blank,
	}
}

// Name is the constraint name.
func (c BlankConstraint) Name() string {
	return "BlankConstraint"
}

// SetUp always returns no error.
func (c BlankConstraint) SetUp() error {
	return nil
}

// When enables conditional validation of this constraint. If the expression evaluates to false,
// then the constraint will be ignored.
func (c BlankConstraint) When(condition bool) BlankConstraint {
	c.isIgnored = !condition
	return c
}

// Message sets the violation message template.
func (c BlankConstraint) Message(message string) BlankConstraint {
	c.messageTemplate = message
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
	return scope.BuildViolation(code.Blank, c.messageTemplate).CreateViolation()
}

// NotNilConstraint checks that a value in not strictly equal to nil. To check that values in not blank use
// NotBlankConstraint.
type NotNilConstraint struct {
	messageTemplate string
	isIgnored       bool
}

// IsNotNil creates a NotNilConstraint to check that a value is not strictly equal to nil.
//
// Example
//  var s *string
//  err := validator.ValidateString(s, it.IsNotNil())
func IsNotNil() NotNilConstraint {
	return NotNilConstraint{
		messageTemplate: message.NotNil,
	}
}

// Name is the constraint name.
func (c NotNilConstraint) Name() string {
	return "NotNilConstraint"
}

// SetUp always returns no error.
func (c NotNilConstraint) SetUp() error {
	return nil
}

// When enables conditional validation of this constraint. If the expression evaluates to false,
// then the constraint will be ignored.
func (c NotNilConstraint) When(condition bool) NotNilConstraint {
	c.isIgnored = !condition
	return c
}

// Message sets the violation message template.
func (c NotNilConstraint) Message(message string) NotNilConstraint {
	c.messageTemplate = message
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
	return scope.BuildViolation(code.NotNil, c.messageTemplate).CreateViolation()
}

// NilConstraint checks that a value in strictly equal to nil. To check that values in blank use
// BlankConstraint.
type NilConstraint struct {
	messageTemplate string
	isIgnored       bool
}

// IsNil creates a NilConstraint to check that a value is strictly equal to nil.
//
// Example
//  var s *string
//  err := validator.ValidateString(s, it.IsNil())
func IsNil() NilConstraint {
	return NilConstraint{
		messageTemplate: message.Nil,
	}
}

func (c NilConstraint) When(condition bool) NilConstraint {
	c.isIgnored = !condition
	return c
}

func (c NilConstraint) Message(message string) NilConstraint {
	c.messageTemplate = message
	return c
}

func (c NilConstraint) SetUp() error {
	return nil
}

func (c NilConstraint) Name() string {
	return "NilConstraint"
}

func (c NilConstraint) ValidateNil(scope validation.Scope) error {
	if c.isIgnored {
		return nil
	}

	return c.newViolation(scope)
}

func (c NilConstraint) ValidateNumber(value generic.Number, scope validation.Scope) error {
	if c.isIgnored {
		return nil
	}
	if value.IsNil() {
		return nil
	}

	return c.newViolation(scope)
}

func (c NilConstraint) ValidateString(value *string, scope validation.Scope) error {
	if c.isIgnored {
		return nil
	}
	if value == nil {
		return nil
	}

	return c.newViolation(scope)
}

func (c NilConstraint) ValidateTime(value *time.Time, scope validation.Scope) error {
	if c.isIgnored {
		return nil
	}
	if value == nil {
		return nil
	}

	return c.newViolation(scope)
}

func (c NilConstraint) ValidateIterable(value generic.Iterable, scope validation.Scope) error {
	if c.isIgnored {
		return nil
	}
	if value.IsNil() {
		return nil
	}

	return c.newViolation(scope)
}

func (c NilConstraint) newViolation(scope validation.Scope) validation.Violation {
	return scope.BuildViolation(code.Nil, c.messageTemplate).CreateViolation()
}

// BoolConstraint checks that a bool value in strictly equal to expected bool value.
type BoolConstraint struct {
	isIgnored            bool
	checkFalse           bool
	checkTrue            bool
	falseMessageTemplate string
	trueMessageTemplate  string
}

func newBoolConstraint(checkFalse, checkTrue bool) BoolConstraint {
	return BoolConstraint{
		checkFalse:           checkFalse,
		checkTrue:            checkTrue,
		falseMessageTemplate: message.False,
		trueMessageTemplate:  message.True,
	}
}

// IsTrue creates a BoolConstraint to check that a value is not strictly equal to true.
//
// Example
//  var b *bool
//  err := validator.ValidateBool(b, it.IsTrue())
func IsTrue() BoolConstraint {
	return newBoolConstraint(false, true)
}

// IsFalse creates a BoolConstraint to check that a value is not strictly equal to false.
//
// Example
//  var b *bool
//  err := validator.ValidateBool(b, it.IsFalse())
func IsFalse() BoolConstraint {
	return newBoolConstraint(true, false)
}

func (c BoolConstraint) When(condition bool) BoolConstraint {
	c.isIgnored = !condition
	return c
}

func (c BoolConstraint) FalseMessage(message string) BoolConstraint {
	c.falseMessageTemplate = message
	return c
}

func (c BoolConstraint) TrueMessage(message string) BoolConstraint {
	c.trueMessageTemplate = message
	return c
}

func (c BoolConstraint) SetUp() error {
	return nil
}

func (c BoolConstraint) Name() string {
	return "BoolConstraint"
}

func (c BoolConstraint) ValidateBool(value *bool, scope validation.Scope) error {
	if c.isIgnored || value == nil {
		return nil
	}

	if c.checkFalse && *value {
		return c.newViolation(code.False, c.falseMessageTemplate, scope)
	}
	if c.checkTrue && !*value {
		return c.newViolation(code.True, c.trueMessageTemplate, scope)
	}

	return nil
}

func (c BoolConstraint) newViolation(violationCode, message string, scope validation.Scope) validation.Violation {
	return scope.BuildViolation(violationCode, message).CreateViolation()
}
