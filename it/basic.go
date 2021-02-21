package it

import (
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/code"
	"github.com/muonsoft/validation/generic"
	"github.com/muonsoft/validation/message"
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

func (c NotBlankConstraint) Set(options *validation.Options) error {
	options.Constraints = append(options.Constraints, c)

	return nil
}

func (c NotBlankConstraint) GetName() string {
	return "NotBlankConstraint"
}

func (c NotBlankConstraint) ValidateNil(options validation.Options) error {
	if c.isIgnored || c.allowNil {
		return nil
	}

	return c.newViolation(options)
}

func (c NotBlankConstraint) ValidateBool(value *bool, options validation.Options) error {
	if c.isIgnored {
		return nil
	}
	if c.allowNil && value == nil {
		return nil
	}
	if value != nil && *value {
		return nil
	}

	return c.newViolation(options)
}

func (c NotBlankConstraint) ValidateNumber(value generic.Number, options validation.Options) error {
	if c.isIgnored {
		return nil
	}
	if c.allowNil && value.IsNil() {
		return nil
	}
	if !value.IsNil() && !value.IsZero() {
		return nil
	}

	return c.newViolation(options)
}

func (c NotBlankConstraint) ValidateString(value *string, options validation.Options) error {
	if c.isIgnored {
		return nil
	}
	if c.allowNil && value == nil {
		return nil
	}
	if value != nil && *value != "" {
		return nil
	}

	return c.newViolation(options)
}

func (c NotBlankConstraint) ValidateIterable(value generic.Iterable, options validation.Options) error {
	if c.isIgnored {
		return nil
	}
	if c.allowNil && value.IsNil() {
		return nil
	}
	if value.Count() > 0 {
		return nil
	}

	return c.newViolation(options)
}

func (c NotBlankConstraint) ValidateCountable(count int, options validation.Options) error {
	if c.isIgnored || count > 0 {
		return nil
	}

	return c.newViolation(options)
}

func (c NotBlankConstraint) newViolation(options validation.Options) validation.Violation {
	return options.NewViolation(code.NotBlank, c.messageTemplate, nil, options.PropertyPath)
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

func (c BlankConstraint) Set(options *validation.Options) error {
	options.Constraints = append(options.Constraints, c)

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

func (c BlankConstraint) ValidateNil(options validation.Options) error {
	return nil
}

func (c BlankConstraint) ValidateBool(value *bool, options validation.Options) error {
	if c.isIgnored || value == nil || !*value {
		return nil
	}

	return c.newViolation(options)
}

func (c BlankConstraint) ValidateNumber(value generic.Number, options validation.Options) error {
	if c.isIgnored || value.IsNil() || value.IsZero() {
		return nil
	}

	return c.newViolation(options)
}

func (c BlankConstraint) ValidateString(value *string, options validation.Options) error {
	if c.isIgnored || value == nil || *value == "" {
		return nil
	}

	return c.newViolation(options)
}

func (c BlankConstraint) ValidateIterable(value generic.Iterable, options validation.Options) error {
	if c.isIgnored || value.Count() == 0 {
		return nil
	}

	return c.newViolation(options)
}

func (c BlankConstraint) ValidateCountable(count int, options validation.Options) error {
	if c.isIgnored || count == 0 {
		return nil
	}

	return c.newViolation(options)
}

func (c BlankConstraint) newViolation(options validation.Options) validation.Violation {
	return options.NewViolation(code.Blank, c.messageTemplate, nil, options.PropertyPath)
}
