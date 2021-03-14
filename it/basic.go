package it

import (
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/code"
	"github.com/muonsoft/validation/generic"
	"github.com/muonsoft/validation/message"

	"time"
)

// NotBlankConstraint checks that a value is not blank: not equal to zero, an empty string, an empty
// slice/array, an empty map, `false` or `nil`. Nil behavior is configurable via AllowNil() method.
// To check that a value is not nil only use NotNilConstraint.
type NotBlankConstraint struct {
	messageTemplate string
	isIgnored       bool
	allowNil        bool
}

func IsNotBlank() NotBlankConstraint {
	return NotBlankConstraint{
		messageTemplate: message.NotBlank,
	}
}

func (c NotBlankConstraint) AllowNil() NotBlankConstraint {
	c.allowNil = true
	return c
}

func (c NotBlankConstraint) When(condition bool) NotBlankConstraint {
	c.isIgnored = !condition
	return c
}

func (c NotBlankConstraint) Message(message string) NotBlankConstraint {
	c.messageTemplate = message
	return c
}

func (c NotBlankConstraint) SetUp(scope *validation.Scope) error {
	return nil
}

func (c NotBlankConstraint) Name() string {
	return "NotBlankConstraint"
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
	return scope.BuildViolation(code.NotBlank, c.messageTemplate).GetViolation()
}

// BlankConstraint checks that a value is blank: equal to `false`, `nil`, zero, an empty string, an empty
// slice, array, or a map.
type BlankConstraint struct {
	messageTemplate string
	isIgnored       bool
}

func IsBlank() BlankConstraint {
	return BlankConstraint{
		messageTemplate: message.Blank,
	}
}

func (c BlankConstraint) SetUp(scope *validation.Scope) error {
	return nil
}

func (c BlankConstraint) Name() string {
	return "BlankConstraint"
}

func (c BlankConstraint) When(condition bool) BlankConstraint {
	c.isIgnored = !condition
	return c
}

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
	return scope.BuildViolation(code.Blank, c.messageTemplate).GetViolation()
}

// NotNilConstraint checks that a value in not strictly equal to `nil`. To check that values in not blank use
// NotBlankConstraint.
type NotNilConstraint struct {
	messageTemplate string
	isIgnored       bool
}

func IsNotNil() NotNilConstraint {
	return NotNilConstraint{
		messageTemplate: message.NotNil,
	}
}

func (c NotNilConstraint) When(condition bool) NotNilConstraint {
	c.isIgnored = !condition
	return c
}

func (c NotNilConstraint) Message(message string) NotNilConstraint {
	c.messageTemplate = message
	return c
}

func (c NotNilConstraint) SetUp(scope *validation.Scope) error {
	return nil
}

func (c NotNilConstraint) Name() string {
	return "NotNilConstraint"
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
	return scope.BuildViolation(code.NotNil, c.messageTemplate).GetViolation()
}
