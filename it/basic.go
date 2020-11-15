package it

import (
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/code"
	"github.com/muonsoft/validation/generic"
	"github.com/muonsoft/validation/message"
)

type NotBlankConstraint struct {
	code            string
	messageTemplate string
	isIgnored       bool
	allowNil        bool
}

func IsNotBlank() NotBlankConstraint {
	return NotBlankConstraint{
		code:            code.NotBlank,
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

func (c NotBlankConstraint) GetCode() string {
	return c.code
}

func (c NotBlankConstraint) GetMessageTemplate() string {
	return c.messageTemplate
}

func (c NotBlankConstraint) GetParameters() map[string]string {
	return nil
}

func (c NotBlankConstraint) ValidateNil(options validation.Options) error {
	if c.isIgnored || c.allowNil {
		return nil
	}

	return options.NewConstraintViolation(c)
}

func (c NotBlankConstraint) ValidateNumber(value generic.Number, options validation.Options) error {
	if c.isIgnored {
		return nil
	}

	if value.IsNil() {
		if c.allowNil {
			return nil
		}
	} else if !value.IsZero() {
		return nil
	}

	return options.NewConstraintViolation(c)
}

func (c NotBlankConstraint) ValidateString(value *string, options validation.Options) error {
	if c.isIgnored {
		return nil
	}

	if value == nil {
		if c.allowNil {
			return nil
		}
	} else if *value != "" {
		return nil
	}

	return options.NewConstraintViolation(c)
}

type BlankConstraint struct {
	code            string
	messageTemplate string
	isIgnored       bool
}

func IsBlank() BlankConstraint {
	return BlankConstraint{
		code:            code.Blank,
		messageTemplate: message.Blank,
	}
}

func (c BlankConstraint) Set(options *validation.Options) error {
	options.Constraints = append(options.Constraints, c)

	return nil
}

func (c BlankConstraint) GetCode() string {
	return c.code
}

func (c BlankConstraint) GetMessageTemplate() string {
	return c.messageTemplate
}

func (c BlankConstraint) GetParameters() map[string]string {
	return nil
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

func (c BlankConstraint) ValidateNumber(value generic.Number, options validation.Options) error {
	if c.isIgnored || value.IsNil() || value.IsZero() {
		return nil
	}

	return options.NewConstraintViolation(c)
}

func (c BlankConstraint) ValidateString(value *string, options validation.Options) error {
	if c.isIgnored || value == nil || *value == "" {
		return nil
	}

	return options.NewConstraintViolation(c)
}
