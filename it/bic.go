package it

import (
	"context"
	"errors"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/validate"
)

// BICConstraint validates whether a string value is a valid Business Identifier Code (BIC / SWIFT).
//
// By default, validation is in strict mode (uppercase ASCII only), matching Symfony Bic::VALIDATION_MODE_STRICT.
// Use [BICConstraint.CaseInsensitive] for Symfony's case-insensitive mode.
//
// Use [BICConstraint.WithIBAN] to assert that the BIC's country/territory matches the given IBAN's country code
// (Symfony Bic "iban" option).
//
// Empty values are skipped; combine with [NotBlank] or similar to reject empty strings.
//
// See https://en.wikipedia.org/wiki/ISO_9362.
type BICConstraint struct {
	isIgnored             bool
	groups                []string
	caseInsensitive       bool
	iban                  string
	err                   error
	ibanErr               error
	messageTemplate       string
	ibanMessageTemplate   string
	messageParameters     validation.TemplateParameterList
	ibanMessageParameters validation.TemplateParameterList
}

// IsBIC validates whether the value is a valid Business Identifier Code (BIC / SWIFT).
// Behavior is aligned with Symfony\Component\Validator\Constraints\Bic.
func IsBIC() BICConstraint {
	return BICConstraint{
		err:                 validation.ErrInvalidBIC,
		ibanErr:             validation.ErrBICIBANCountryMismatch,
		messageTemplate:     validation.ErrInvalidBIC.Message(),
		ibanMessageTemplate: validation.ErrBICIBANCountryMismatch.Message(),
	}
}

// CaseInsensitive enables Symfony Bic::VALIDATION_MODE_CASE_INSENSITIVE (lowercase letters allowed).
func (c BICConstraint) CaseInsensitive() BICConstraint {
	c.caseInsensitive = true
	return c
}

// WithIBAN sets an IBAN to check that its country code matches the BIC (Symfony Bic "iban" option).
func (c BICConstraint) WithIBAN(iban string) BICConstraint {
	c.iban = iban
	return c
}

// WithError overrides the default error for the main BIC format violation.
// IBAN country mismatch still uses [validation.ErrBICIBANCountryMismatch] unless overridden via [BICConstraint.WithIBANError].
func (c BICConstraint) WithError(err error) BICConstraint {
	c.err = err
	return c
}

// WithIBANError overrides the default error for BIC vs IBAN country mismatch.
func (c BICConstraint) WithIBANError(err error) BICConstraint {
	c.ibanErr = err
	return c
}

// WithMessage sets the violation message template for the main BIC format violation.
func (c BICConstraint) WithMessage(template string, parameters ...validation.TemplateParameter) BICConstraint {
	c.messageTemplate = template
	c.messageParameters = parameters
	return c
}

// WithIBANMessage sets the violation message template when the BIC does not match the configured IBAN.
// Default parameters: {{ value }}, {{ iban }}.
func (c BICConstraint) WithIBANMessage(template string, parameters ...validation.TemplateParameter) BICConstraint {
	c.ibanMessageTemplate = template
	c.ibanMessageParameters = parameters
	return c
}

// When enables conditional validation of this constraint.
func (c BICConstraint) When(condition bool) BICConstraint {
	c.isIgnored = !condition
	return c
}

// WhenGroups enables conditional validation by validation groups.
func (c BICConstraint) WhenGroups(groups ...string) BICConstraint {
	c.groups = groups
	return c
}

func (c BICConstraint) bicOptions() []func(*validate.BICOptions) {
	var opts []func(*validate.BICOptions)
	if c.caseInsensitive {
		opts = append(opts, validate.BICCaseInsensitive())
	}
	if c.iban != "" {
		opts = append(opts, validate.BICWithIBAN(c.iban))
	}
	return opts
}

func (c BICConstraint) ValidateString(ctx context.Context, validator *validation.Validator, value *string) error {
	if c.isIgnored || validator.IsIgnoredForGroups(c.groups...) || value == nil || *value == "" {
		return nil
	}

	err := validate.BIC(*value, c.bicOptions()...)
	if err == nil {
		return nil
	}

	if errors.Is(err, validate.ErrBICIBANCountryMismatch) {
		return validator.BuildViolation(ctx, c.ibanErr, c.ibanMessageTemplate).
			WithParameters(
				c.ibanMessageParameters.Prepend(
					validation.TemplateParameter{Key: "{{ iban }}", Value: c.iban},
					validation.TemplateParameter{Key: "{{ value }}", Value: *value},
				)...,
			).
			Create()
	}

	return validator.BuildViolation(ctx, c.err, c.messageTemplate).
		WithParameters(
			c.messageParameters.Prepend(
				validation.TemplateParameter{Key: "{{ value }}", Value: *value},
			)...,
		).
		Create()
}

// Validate implements [validation.Constraint][string].
func (c BICConstraint) Validate(ctx context.Context, validator *validation.Validator, v string) error {
	return c.ValidateString(ctx, validator, &v)
}
