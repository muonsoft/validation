package it

import (
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/code"
	"github.com/muonsoft/validation/generic"
	"github.com/muonsoft/validation/message"

	"time"
)

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

func (c NotBlankConstraint) Set(scope *validation.Scope) error {
	scope.AddConstraint(c)

	return nil
}

func (c NotBlankConstraint) GetName() string {
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

func (c NotBlankConstraint) ValidateTime(time *time.Time, scope validation.Scope) error {
	if c.isIgnored {
		return nil
	}
	if c.allowNil && time == nil {
		return nil
	}
	if time != nil {
		return nil
	}

	return c.newViolation(scope)
}

func (c NotBlankConstraint) newViolation(scope validation.Scope) validation.Violation {
	return scope.BuildViolation(code.NotBlank, c.messageTemplate).GetViolation()
}

type BlankConstraint struct {
	messageTemplate string
	isIgnored       bool
}

func IsBlank() BlankConstraint {
	return BlankConstraint{
		messageTemplate: message.Blank,
	}
}

func (c BlankConstraint) Set(scope *validation.Scope) error {
	scope.AddConstraint(c)

	return nil
}

func (c BlankConstraint) GetName() string {
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

func (c BlankConstraint) ValidateTime(time *time.Time, scope validation.Scope) error {
	if c.isIgnored || time == nil {
		return nil
	}

	return c.newViolation(scope)
}

func (c BlankConstraint) newViolation(scope validation.Scope) validation.Violation {
	return scope.BuildViolation(code.Blank, c.messageTemplate).GetViolation()
}
