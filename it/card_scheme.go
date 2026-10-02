package it

import (
	"context"
	"errors"
	"slices"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/validate"
)

// CardSchemeConstraint checks whether a card number matches a selected scheme.
// It checks prefixes and lengths only; combine with [IsLUHN] for a checksum.
// Nil and empty values are skipped. Use [IsNotBlank] to require a value.
type CardSchemeConstraint struct {
	isIgnored         bool
	groups            []string
	schemes           []validate.CardSchemeName
	err               error
	messageTemplate   string
	messageParameters validation.TemplateParameterList
}

// IsCardScheme creates a constraint for one or more explicit payment card schemes.
// Use the validate.CardScheme* constants. Empty or unknown scheme selections
// produce a configuration error during validation, including for nil/empty input.
// See [validate.CardScheme] for the numbering rules.
func IsCardScheme(schemes ...validate.CardSchemeName) CardSchemeConstraint {
	return CardSchemeConstraint{
		schemes:         slices.Clone(schemes),
		err:             validation.ErrInvalidCardScheme,
		messageTemplate: validation.ErrInvalidCardScheme.Message(),
	}
}

// WithError overrides the underlying error for a nonmatching card number.
func (c CardSchemeConstraint) WithError(err error) CardSchemeConstraint {
	c.err = err
	return c
}

// WithMessage overrides the violation message; {{ value }} contains the input.
func (c CardSchemeConstraint) WithMessage(template string, parameters ...validation.TemplateParameter) CardSchemeConstraint {
	c.messageTemplate = template
	c.messageParameters = parameters
	return c
}

// When disables validation when condition is false.
func (c CardSchemeConstraint) When(condition bool) CardSchemeConstraint {
	c.isIgnored = !condition
	return c
}

// WhenGroups limits validation to the specified groups.
func (c CardSchemeConstraint) WhenGroups(groups ...string) CardSchemeConstraint {
	c.groups = groups
	return c
}

// ValidateString validates an optional card number against the configured schemes.
func (c CardSchemeConstraint) ValidateString(ctx context.Context, validator *validation.Validator, value *string) error {
	if c.isIgnored || validator.IsIgnoredForGroups(c.groups...) {
		return nil
	}
	input := ""
	if value != nil {
		input = *value
	}
	err := validate.CardScheme(input, c.schemes...)
	if err == nil {
		return nil
	}
	if errors.Is(err, validate.ErrInvalidCardSchemes) {
		return validator.CreateConstraintError("CardSchemeConstraint", "at least one known card scheme is required")
	}
	return validator.BuildViolation(ctx, c.err, c.messageTemplate).WithParameters(
		c.messageParameters.Prepend(validation.TemplateParameter{Key: "{{ value }}", Value: input})...,
	).Create()
}

// Validate implements [validation.Constraint][string] for [validation.This] and [validation.Each].
func (c CardSchemeConstraint) Validate(ctx context.Context, validator *validation.Validator, value string) error {
	return c.ValidateString(ctx, validator, &value)
}
