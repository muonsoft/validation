package it

import (
	"context"
	"errors"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/validate"
)

// WeekConstraint validates ISO weeks with optional inclusive bounds.
// See [validate.Week] for the format and literal year range 0000–9999.
// Nil and empty values are ignored. Use [IsNotBlank] to require a value.
type WeekConstraint struct {
	isIgnored        bool
	groups           []string
	min, max         string
	formatErr        error
	formatMessage    string
	formatParameters validation.TemplateParameterList
	numberErr        error
	numberMessage    string
	numberParameters validation.TemplateParameterList
	minErr           error
	minMessage       string
	minParameters    validation.TemplateParameterList
	maxErr           error
	maxMessage       string
	maxParameters    validation.TemplateParameterList
}

// IsWeek creates a constraint for ISO weeks in YYYY-Www format.
// WithMin and WithMax restrict the range; malformed or reversed bounds
// produce a constraint configuration error during validation.
func IsWeek() WeekConstraint {
	return WeekConstraint{
		formatErr:     validation.ErrInvalidWeekFormat,
		formatMessage: validation.ErrInvalidWeekFormat.Message(),
		numberErr:     validation.ErrInvalidWeekNumber,
		numberMessage: validation.ErrInvalidWeekNumber.Message(),
		minErr:        validation.ErrWeekTooEarly,
		minMessage:    validation.ErrWeekTooEarly.Message(),
		maxErr:        validation.ErrWeekTooLate,
		maxMessage:    validation.ErrWeekTooLate.Message(),
	}
}

// WithMin sets the inclusive minimum ISO week; empty removes the bound.
func (c WeekConstraint) WithMin(value string) WeekConstraint {
	c.min = value
	return c
}

// WithMax sets the inclusive maximum ISO week; empty removes the bound.
func (c WeekConstraint) WithMax(value string) WeekConstraint {
	c.max = value
	return c
}

// WithFormatError overrides the underlying error for the format check.
func (c WeekConstraint) WithFormatError(err error) WeekConstraint {
	c.formatErr = err
	return c
}

// WithFormatMessage overrides the message for the format check.
// Available parameters are {{ value }}, {{ min }}, and {{ max }}.
func (c WeekConstraint) WithFormatMessage(template string, parameters ...validation.TemplateParameter) WeekConstraint {
	c.formatMessage = template
	c.formatParameters = parameters
	return c
}

// WithWeekNumberError overrides the underlying error for the number check.
func (c WeekConstraint) WithWeekNumberError(err error) WeekConstraint {
	c.numberErr = err
	return c
}

// WithWeekNumberMessage overrides the message for the number check.
// Available parameters are {{ value }}, {{ min }}, and {{ max }}.
func (c WeekConstraint) WithWeekNumberMessage(template string, parameters ...validation.TemplateParameter) WeekConstraint {
	c.numberMessage = template
	c.numberParameters = parameters
	return c
}

// WithMinError overrides the underlying error for the min check.
func (c WeekConstraint) WithMinError(err error) WeekConstraint {
	c.minErr = err
	return c
}

// WithMinMessage overrides the message for the min check.
// Available parameters are {{ value }}, {{ min }}, and {{ max }}.
func (c WeekConstraint) WithMinMessage(template string, parameters ...validation.TemplateParameter) WeekConstraint {
	c.minMessage = template
	c.minParameters = parameters
	return c
}

// WithMaxError overrides the underlying error for the max check.
func (c WeekConstraint) WithMaxError(err error) WeekConstraint {
	c.maxErr = err
	return c
}

// WithMaxMessage overrides the message for the max check.
// Available parameters are {{ value }}, {{ min }}, and {{ max }}.
func (c WeekConstraint) WithMaxMessage(template string, parameters ...validation.TemplateParameter) WeekConstraint {
	c.maxMessage = template
	c.maxParameters = parameters
	return c
}

// When disables validation when condition is false.
func (c WeekConstraint) When(condition bool) WeekConstraint {
	c.isIgnored = !condition
	return c
}

// WhenGroups limits validation to the specified groups.
func (c WeekConstraint) WhenGroups(groups ...string) WeekConstraint {
	c.groups = groups
	return c
}

// ValidateString validates an optional ISO week string.
func (c WeekConstraint) ValidateString(ctx context.Context, validator *validation.Validator, value *string) error {
	if c.isIgnored || validator.IsIgnoredForGroups(c.groups...) {
		return nil
	}
	input := ""
	if value != nil {
		input = *value
	}
	err := validate.Week(input, validate.WithMinWeek(c.min), validate.WithMaxWeek(c.max))
	if err == nil {
		return nil
	}
	if errors.Is(err, validate.ErrInvalidWeekBounds) {
		return validator.CreateConstraintError("WeekConstraint", "bounds must be valid ISO weeks and minimum must not exceed maximum")
	}
	template, parameters, cause := c.violation(err)
	return validator.BuildViolation(ctx, cause, template).WithParameters(parameters.Prepend(
		validation.TemplateParameter{Key: "{{ value }}", Value: input},
		validation.TemplateParameter{Key: "{{ min }}", Value: c.min},
		validation.TemplateParameter{Key: "{{ max }}", Value: c.max},
	)...).Create()
}

func (c WeekConstraint) violation(err error) (string, validation.TemplateParameterList, error) {
	switch {
	case errors.Is(err, validate.ErrInvalidWeekNumber):
		return c.numberMessage, c.numberParameters, c.numberErr
	case errors.Is(err, validate.ErrWeekTooEarly):
		return c.minMessage, c.minParameters, c.minErr
	case errors.Is(err, validate.ErrWeekTooLate):
		return c.maxMessage, c.maxParameters, c.maxErr
	default:
		return c.formatMessage, c.formatParameters, c.formatErr
	}
}

// Validate implements [validation.Constraint][string] for [validation.This] and [validation.Each].
func (c WeekConstraint) Validate(ctx context.Context, validator *validation.Validator, value string) error {
	return c.ValidateString(ctx, validator, &value)
}
