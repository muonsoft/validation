package it

import (
	"context"
	"fmt"
	"strings"

	"github.com/muonsoft/validation"
)

// ChoiceConstraint is used to ensure that the given value corresponds to one of the expected choices.
// Zero values (zero numbers or empty strings) are also compared with the given choices.
// In order for a blank value to be valid, use the WithAllowedBlank method.
type ChoiceConstraint[T comparable] struct {
	blank             T
	choices           map[T]bool
	choicesValue      string
	groups            []string
	err               error
	messageTemplate   string
	messageParameters validation.TemplateParameterList
	allowBlank        bool
	isIgnored         bool
}

// IsOneOf creates a ChoiceConstraint for checking that values are in the expected list of values.
// Zero values (zero numbers or empty strings) are also compared with the given choices.
// In order for a blank value to be valid, use the WithAllowedBlank method.
func IsOneOf[T comparable](values ...T) ChoiceConstraint[T] {
	choices := make(map[T]bool, len(values))
	for _, value := range values {
		choices[value] = true
	}

	s := strings.Builder{}
	for i, value := range values {
		if i > 0 {
			s.WriteString(", ")
		}
		s.WriteString(fmt.Sprint(value))
	}

	return ChoiceConstraint[T]{
		choices:         choices,
		choicesValue:    s.String(),
		err:             validation.ErrNoSuchChoice,
		messageTemplate: validation.ErrNoSuchChoice.Message(),
	}
}

// WithAllowedBlank makes zero values valid.
func (c ChoiceConstraint[T]) WithAllowedBlank() ChoiceConstraint[T] {
	c.allowBlank = true
	return c
}

// WithError overrides default error for produced violation.
func (c ChoiceConstraint[T]) WithError(err error) ChoiceConstraint[T] {
	c.err = err
	return c
}

// WithMessage sets the violation message template. You can set custom template parameters
// for injecting its values into the final message. Also, you can use default parameters:
//
//	{{ choices }} - a comma-separated list of available choices;
//	{{ value }} - the current (invalid) value.
func (c ChoiceConstraint[T]) WithMessage(template string, parameters ...validation.TemplateParameter) ChoiceConstraint[T] {
	c.messageTemplate = template
	c.messageParameters = parameters
	return c
}

// When enables conditional validation of this constraint. If the expression evaluates to false,
// then the constraint will be ignored.
func (c ChoiceConstraint[T]) When(condition bool) ChoiceConstraint[T] {
	c.isIgnored = !condition
	return c
}

// WhenGroups enables conditional validation of the constraint by using the validation groups.
func (c ChoiceConstraint[T]) WhenGroups(groups ...string) ChoiceConstraint[T] {
	c.groups = groups
	return c
}

func (c ChoiceConstraint[T]) ValidateNumber(ctx context.Context, validator *validation.Validator, value *T) error {
	return c.ValidateComparable(ctx, validator, value)
}

func (c ChoiceConstraint[T]) ValidateString(ctx context.Context, validator *validation.Validator, value *T) error {
	return c.ValidateComparable(ctx, validator, value)
}

func (c ChoiceConstraint[T]) ValidateComparable(ctx context.Context, validator *validation.Validator, value *T) error {
	if len(c.choices) == 0 {
		return validator.CreateConstraintError("ChoiceConstraint", "empty list of choices")
	}
	if c.isIgnored || validator.IsIgnoredForGroups(c.groups...) || value == nil || c.allowBlank && *value == c.blank {
		return nil
	}
	if c.choices[*value] {
		return nil
	}

	return validator.
		BuildViolation(ctx, c.err, c.messageTemplate).
		WithParameters(
			c.messageParameters.Prepend(
				validation.TemplateParameter{Key: "{{ value }}", Value: fmt.Sprint(*value)},
				validation.TemplateParameter{Key: "{{ choices }}", Value: c.choicesValue},
			)...,
		).
		Create()
}
