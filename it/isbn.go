package it

import (
	"context"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/validate"
)

type isbnMode int

const (
	isbnModeAny isbnMode = iota
	isbnMode10
	isbnMode13
)

// ISBNConstraint validates whether a string value is a valid ISBN-10 or ISBN-13.
// Hyphens (U+002D) are stripped before validation, matching Symfony Isbn.
//
// By default, either ISBN-10 or ISBN-13 is accepted. Use [ISBNConstraint.Only10] or
// [ISBNConstraint.Only13] to restrict the format.
//
// Empty values are skipped; combine with [NotBlank] or similar to reject empty strings.
//
// Behavior is aligned with Symfony\Component\Validator\Constraints\Isbn.
//
// See https://en.wikipedia.org/wiki/ISBN.
type ISBNConstraint struct {
	isIgnored     bool
	groups        []string
	mode          isbnMode
	options       []func(*validate.ISBNOptions)
	err10         error
	err13         error
	errBoth       error
	message10     string
	message13     string
	messageBoth   string
	messageParams validation.TemplateParameterList
}

// IsISBN validates whether the value is a valid ISBN-10 or ISBN-13.
func IsISBN() ISBNConstraint {
	return ISBNConstraint{
		err10:       validation.ErrInvalidISBN10,
		err13:       validation.ErrInvalidISBN13,
		errBoth:     validation.ErrInvalidISBN,
		message10:   validation.ErrInvalidISBN10.Message(),
		message13:   validation.ErrInvalidISBN13.Message(),
		messageBoth: validation.ErrInvalidISBN.Message(),
	}
}

// Only10 restricts validation to ISBN-10 (Symfony Isbn::ISBN_10).
func (c ISBNConstraint) Only10() ISBNConstraint {
	c.mode = isbnMode10
	c.options = append(c.options, validate.ISBNOnly10())
	return c
}

// Only13 restricts validation to ISBN-13 (Symfony Isbn::ISBN_13).
func (c ISBNConstraint) Only13() ISBNConstraint {
	c.mode = isbnMode13
	c.options = append(c.options, validate.ISBNOnly13())
	return c
}

// WithError overrides the default errors for all ISBN modes.
func (c ISBNConstraint) WithError(err error) ISBNConstraint {
	c.err10 = err
	c.err13 = err
	c.errBoth = err
	return c
}

// WithISBN10Error overrides the error used when validating as ISBN-10 only.
func (c ISBNConstraint) WithISBN10Error(err error) ISBNConstraint {
	c.err10 = err
	return c
}

// WithISBN13Error overrides the error used when validating as ISBN-13 only.
func (c ISBNConstraint) WithISBN13Error(err error) ISBNConstraint {
	c.err13 = err
	return c
}

// WithBothISBNError overrides the error used when either ISBN-10 or ISBN-13 is accepted.
func (c ISBNConstraint) WithBothISBNError(err error) ISBNConstraint {
	c.errBoth = err
	return c
}

// WithMessage sets the message template for all ISBN modes.
func (c ISBNConstraint) WithMessage(template string, parameters ...validation.TemplateParameter) ISBNConstraint {
	c.message10 = template
	c.message13 = template
	c.messageBoth = template
	c.messageParams = parameters
	return c
}

// WithISBN10Message sets the message template for ISBN-10-only validation.
func (c ISBNConstraint) WithISBN10Message(template string, parameters ...validation.TemplateParameter) ISBNConstraint {
	c.message10 = template
	c.messageParams = parameters
	return c
}

// WithISBN13Message sets the message template for ISBN-13-only validation.
func (c ISBNConstraint) WithISBN13Message(template string, parameters ...validation.TemplateParameter) ISBNConstraint {
	c.message13 = template
	c.messageParams = parameters
	return c
}

// WithBothISBNMessage sets the message template when either format is accepted.
func (c ISBNConstraint) WithBothISBNMessage(template string, parameters ...validation.TemplateParameter) ISBNConstraint {
	c.messageBoth = template
	c.messageParams = parameters
	return c
}

// When enables conditional validation of this constraint.
func (c ISBNConstraint) When(condition bool) ISBNConstraint {
	c.isIgnored = !condition
	return c
}

// WhenGroups enables conditional validation by validation groups.
func (c ISBNConstraint) WhenGroups(groups ...string) ISBNConstraint {
	c.groups = groups
	return c
}

func (c ISBNConstraint) ValidateString(ctx context.Context, validator *validation.Validator, value *string) error {
	if c.isIgnored || validator.IsIgnoredForGroups(c.groups...) || value == nil || *value == "" {
		return nil
	}

	err := validate.ISBN(*value, c.options...)
	if err == nil {
		return nil
	}

	var vErr error
	var msg string

	switch c.mode {
	case isbnMode10:
		vErr = c.err10
		msg = c.message10
	case isbnMode13:
		vErr = c.err13
		msg = c.message13
	default:
		vErr = c.errBoth
		msg = c.messageBoth
	}

	return validator.BuildViolation(ctx, vErr, msg).
		WithParameters(
			c.messageParams.Prepend(
				validation.TemplateParameter{Key: "{{ value }}", Value: *value},
			)...,
		).
		Create()
}

// Validate implements [validation.Constraint][string].
func (c ISBNConstraint) Validate(ctx context.Context, validator *validation.Validator, v string) error {
	return c.ValidateString(ctx, validator, &v)
}
