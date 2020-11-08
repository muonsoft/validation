package it

import (
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/code"
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

func (c NotBlankConstraint) Code() string {
	return c.code
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

	return c.createViolation(options)
}

func (c NotBlankConstraint) ValidateInt(value *int, options validation.Options) error {
	if c.isIgnored {
		return nil
	}

	if value == nil {
		if c.allowNil {
			return nil
		}
	} else if *value != 0 {
		return nil
	}

	return c.createViolation(options)
}

func (c NotBlankConstraint) ValidateFloat(value *float64, options validation.Options) error {
	if c.isIgnored {
		return nil
	}

	if value == nil {
		if c.allowNil {
			return nil
		}
	} else if *value != 0 {
		return nil
	}

	return c.createViolation(options)
}

func (c NotBlankConstraint) createViolation(options validation.Options) validation.Violation {
	return validation.NewViolation(
		c.code,
		c.messageTemplate,
		nil,
		options.PropertyPath,
	)
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

func (c BlankConstraint) Code() string {
	return c.code
}

func (c BlankConstraint) When(condition bool) BlankConstraint {
	c.isIgnored = !condition
	return c
}

func (c BlankConstraint) Message(message string) BlankConstraint {
	c.messageTemplate = message
	return c
}

func (c BlankConstraint) ValidateString(value *string, options validation.Options) error {
	if c.isIgnored || value == nil || *value == "" {
		return nil
	}

	return c.createViolation(options)
}

func (c BlankConstraint) ValidateInt(value *int, options validation.Options) error {
	if c.isIgnored || value == nil || *value == 0 {
		return nil
	}

	return c.createViolation(options)
}

func (c BlankConstraint) ValidateFloat(value *float64, options validation.Options) error {
	if c.isIgnored || value == nil || *value == 0 {
		return nil
	}

	return c.createViolation(options)
}

func (c BlankConstraint) createViolation(options validation.Options) validation.Violation {
	return validation.NewViolation(
		c.code,
		c.messageTemplate,
		nil,
		options.PropertyPath,
	)
}
