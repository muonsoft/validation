package it

import (
	"context"
	"slices"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/validate"
)

// CSSColorConstraint checks whether a CSS color matches a selected format.
// Supported syntaxes follow Symfony 8.0; see [validate.CSSColor].
// Nil and empty values are skipped. Use [IsNotBlank] to require a value.
type CSSColorConstraint struct {
	isIgnored         bool
	invalidFormats    bool
	groups            []string
	formats           []validate.CSSColorFormat
	err               error
	messageTemplate   string
	messageParameters validation.TemplateParameterList
}

// IsCSSColor creates a CSS color constraint. No formats selects all supported
// formats; use validate.CSSColor* constants to restrict them. Unknown formats
// produce a configuration error during validation, including for nil/empty input.
func IsCSSColor(formats ...validate.CSSColorFormat) CSSColorConstraint {
	return CSSColorConstraint{
		formats:         slices.Clone(formats),
		invalidFormats:  validate.CSSColor("", formats...) != nil,
		err:             validation.ErrInvalidCSSColor,
		messageTemplate: validation.ErrInvalidCSSColor.Message(),
	}
}

// WithError overrides the underlying error for a nonmatching CSS color.
func (c CSSColorConstraint) WithError(err error) CSSColorConstraint {
	c.err = err
	return c
}

// WithMessage overrides the violation message; {{ value }} contains the input.
func (c CSSColorConstraint) WithMessage(template string, parameters ...validation.TemplateParameter) CSSColorConstraint {
	c.messageTemplate = template
	c.messageParameters = parameters
	return c
}

// When disables validation when condition is false.
func (c CSSColorConstraint) When(condition bool) CSSColorConstraint {
	c.isIgnored = !condition
	return c
}

// WhenGroups limits validation to the specified groups.
func (c CSSColorConstraint) WhenGroups(groups ...string) CSSColorConstraint {
	c.groups = groups
	return c
}

// ValidateString validates an optional CSS color against the configured formats.
func (c CSSColorConstraint) ValidateString(ctx context.Context, validator *validation.Validator, value *string) error {
	if c.isIgnored || validator.IsIgnoredForGroups(c.groups...) {
		return nil
	}
	if c.invalidFormats {
		return validator.CreateConstraintError("CSSColorConstraint", "unknown CSS color format")
	}
	if value == nil || *value == "" {
		return nil
	}
	if validate.CSSColor(*value, c.formats...) == nil {
		return nil
	}
	return validator.BuildViolation(ctx, c.err, c.messageTemplate).WithParameters(
		c.messageParameters.Prepend(validation.TemplateParameter{Key: "{{ value }}", Value: *value})...,
	).Create()
}

// Validate implements [validation.Constraint][string] for [validation.This] and [validation.Each].
func (c CSSColorConstraint) Validate(ctx context.Context, validator *validation.Validator, value string) error {
	return c.ValidateString(ctx, validator, &value)
}
