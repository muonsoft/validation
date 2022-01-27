package it

import (
	"strings"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/code"
	"github.com/muonsoft/validation/message"
)

// ChoiceConstraint is used to ensure that the given value corresponds to one of the expected choices.
type ChoiceConstraint struct {
	choices           map[string]bool
	choicesValue      string
	groups            []string
	code              string
	messageTemplate   string
	messageParameters validation.TemplateParameterList
	isIgnored         bool
}

// IsOneOfStrings creates a ChoiceConstraint for checking that values are in the expected list of strings.
func IsOneOfStrings(values ...string) ChoiceConstraint {
	choices := make(map[string]bool, len(values))
	for _, value := range values {
		choices[value] = true
	}

	return ChoiceConstraint{
		choices:         choices,
		choicesValue:    strings.Join(values, ", "),
		code:            code.NoSuchChoice,
		messageTemplate: message.Templates[code.NoSuchChoice],
	}
}

// SetUp will return an error if the list of choices is empty.
func (c ChoiceConstraint) SetUp() error {
	if len(c.choices) == 0 {
		return errEmptyChoices
	}

	return nil
}

// Name is the constraint name.
func (c ChoiceConstraint) Name() string {
	return "ChoiceConstraint"
}

// Code overrides default code for produced violation.
func (c ChoiceConstraint) Code(code string) ChoiceConstraint {
	c.code = code
	return c
}

// Message sets the violation message template. You can set custom template parameters
// for injecting its values into the final message. Also, you can use default parameters:
//
//	{{ choices }} - a comma-separated list of available choices;
//	{{ value }} - the current (invalid) value.
func (c ChoiceConstraint) Message(template string, parameters ...validation.TemplateParameter) ChoiceConstraint {
	c.messageTemplate = template
	c.messageParameters = parameters
	return c
}

// When enables conditional validation of this constraint. If the expression evaluates to false,
// then the constraint will be ignored.
func (c ChoiceConstraint) When(condition bool) ChoiceConstraint {
	c.isIgnored = !condition
	return c
}

// WhenGroups enables conditional validation of the constraint by using the validation groups.
func (c ChoiceConstraint) WhenGroups(groups ...string) ChoiceConstraint {
	c.groups = groups
	return c
}

func (c ChoiceConstraint) ValidateString(value *string, scope validation.Scope) error {
	if c.isIgnored || scope.IsIgnored(c.groups...) || value == nil || *value == "" {
		return nil
	}
	if c.choices[*value] {
		return nil
	}

	return scope.
		BuildViolation(c.code, c.messageTemplate).
		SetParameters(
			c.messageParameters.Prepend(
				validation.TemplateParameter{Key: "{{ value }}", Value: *value},
				validation.TemplateParameter{Key: "{{ choices }}", Value: c.choicesValue},
			)...,
		).
		CreateViolation()
}
