package it

import (
	"context"
	"time"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/validate"
)

// DateTimeConstraint checks that the string value is a valid date and time value specified by a specific layout.
// The layout can be redefined using the [DateTimeConstraint.WithLayout] method.
type DateTimeConstraint struct {
	isIgnored         bool
	groups            []string
	err               error
	layout            string
	messageTemplate   string
	messageParameters validation.TemplateParameterList
}

// IsDateTime checks that the string value is a valid date and time. By default, it uses [time.RFC3339] layout.
// The layout can be redefined using the [DateTimeConstraint.WithLayout] method.
func IsDateTime() DateTimeConstraint {
	return DateTimeConstraint{
		layout:          time.RFC3339,
		err:             validation.ErrInvalidDateTime,
		messageTemplate: validation.ErrInvalidDateTime.Message(),
	}
}

// IsDate checks that the string value is a valid date. It uses "2006-01-02" layout.
// The layout can be redefined using the [DateTimeConstraint.WithLayout] method.
func IsDate() DateTimeConstraint {
	return DateTimeConstraint{
		layout:          "2006-01-02",
		err:             validation.ErrInvalidDate,
		messageTemplate: validation.ErrInvalidDate.Message(),
	}
}

// IsTime checks that the string value is a valid time. It uses "15:04:05" layout.
// The layout can be redefined using the WithLayout method.
func IsTime() DateTimeConstraint {
	return DateTimeConstraint{
		layout:          "15:04:05",
		err:             validation.ErrInvalidTime,
		messageTemplate: validation.ErrInvalidTime.Message(),
	}
}

// WithLayout specifies the layout to be used for datetime parsing.
func (c DateTimeConstraint) WithLayout(layout string) DateTimeConstraint {
	c.layout = layout
	return c
}

// WithError overrides default error for produced violation.
func (c DateTimeConstraint) WithError(err error) DateTimeConstraint {
	c.err = err
	return c
}

// WithMessage sets the violation message template. You can set custom template parameters
// for injecting its values into the final message. Also, you can use default parameters:
//
//	{{ layout }} - date time layout used for parsing;
//	{{ value }} - the current (invalid) value.
func (c DateTimeConstraint) WithMessage(template string, parameters ...validation.TemplateParameter) DateTimeConstraint {
	c.messageTemplate = template
	c.messageParameters = parameters
	return c
}

// When enables conditional validation of this constraint. If the expression evaluates to false,
// then the constraint will be ignored.
func (c DateTimeConstraint) When(condition bool) DateTimeConstraint {
	c.isIgnored = !condition
	return c
}

// WhenGroups enables conditional validation of the constraint by using the validation groups.
func (c DateTimeConstraint) WhenGroups(groups ...string) DateTimeConstraint {
	c.groups = groups
	return c
}

func (c DateTimeConstraint) ValidateString(ctx context.Context, validator *validation.Validator, value *string) error {
	if c.isIgnored || validator.IsIgnoredForGroups(c.groups...) || value == nil || *value == "" {
		return nil
	}
	if _, err := time.Parse(c.layout, *value); err == nil {
		return nil
	}

	return validator.BuildViolation(ctx, c.err, c.messageTemplate).
		WithParameters(
			c.messageParameters.Prepend(
				validation.TemplateParameter{Key: "{{ layout }}", Value: c.layout},
				validation.TemplateParameter{Key: "{{ value }}", Value: *value},
			)...,
		).
		WithParameter("{{ value }}", *value).Create()
}

// Validate implements [validation.Constraint][string] so the constraint can be used with [validation.Each] and [validation.This].
func (c DateTimeConstraint) Validate(ctx context.Context, validator *validation.Validator, v string) error {
	return c.ValidateString(ctx, validator, &v)
}

// TimezoneConstraint validates whether the string value is a known IANA timezone identifier.
// See [validate.Timezone] for the supported identifiers.
// Use [TimezoneConstraint.WithZone] to restrict identifiers to a geographical region.
type TimezoneConstraint struct {
	isIgnored         bool
	groups            []string
	zone              validate.TimezoneZone
	err               error
	messageTemplate   string
	messageParameters validation.TemplateParameterList
}

// IsTimezone validates whether the string value is a known IANA timezone identifier.
// See [TimezoneConstraint] for configuration options.
func IsTimezone() TimezoneConstraint {
	return TimezoneConstraint{
		err:             validation.ErrInvalidTimezone,
		messageTemplate: validation.ErrInvalidTimezone.Message(),
	}
}

// WithZone restricts valid timezone identifiers to the given geographical region
// (default accepts any IANA zone). Allowed values are [validate.TimezoneZoneAll],
// [validate.TimezoneZoneAfrica], [validate.TimezoneZoneAmerica], [validate.TimezoneZoneAntarctica],
// [validate.TimezoneZoneArctic], [validate.TimezoneZoneAsia], [validate.TimezoneZoneAtlantic],
// [validate.TimezoneZoneAustralia], [validate.TimezoneZoneEurope], [validate.TimezoneZoneIndian],
// and [validate.TimezoneZonePacific].
func (c TimezoneConstraint) WithZone(zone validate.TimezoneZone) TimezoneConstraint {
	c.zone = zone
	return c
}

// WithError overrides default error for produced violation.
func (c TimezoneConstraint) WithError(err error) TimezoneConstraint {
	c.err = err
	return c
}

// WithMessage sets the violation message template. You can set custom template parameters
// for injecting its values into the final message. Also, you can use default parameters:
//
//	{{ value }} - the current (invalid) value.
func (c TimezoneConstraint) WithMessage(template string, parameters ...validation.TemplateParameter) TimezoneConstraint {
	c.messageTemplate = template
	c.messageParameters = parameters
	return c
}

// When enables conditional validation of this constraint. If the expression evaluates to false,
// then the constraint will be ignored.
func (c TimezoneConstraint) When(condition bool) TimezoneConstraint {
	c.isIgnored = !condition
	return c
}

// WhenGroups enables conditional validation of the constraint by using the validation groups.
func (c TimezoneConstraint) WhenGroups(groups ...string) TimezoneConstraint {
	c.groups = groups
	return c
}

func (c TimezoneConstraint) ValidateString(ctx context.Context, validator *validation.Validator, value *string) error {
	if c.isIgnored || validator.IsIgnoredForGroups(c.groups...) || value == nil || *value == "" {
		return nil
	}
	if validate.Timezone(*value, validate.WithTimezoneZone(c.zone)) == nil {
		return nil
	}

	return validator.BuildViolation(ctx, c.err, c.messageTemplate).
		WithParameters(
			c.messageParameters.Prepend(
				validation.TemplateParameter{Key: "{{ value }}", Value: *value},
			)...,
		).
		Create()
}

// Validate implements [validation.Constraint][string] so the constraint can be used with [validation.Each] and [validation.This].
func (c TimezoneConstraint) Validate(ctx context.Context, validator *validation.Validator, v string) error {
	return c.ValidateString(ctx, validator, &v)
}
