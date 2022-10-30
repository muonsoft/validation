package it

import (
	"context"
	"time"

	"github.com/muonsoft/validation"
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
