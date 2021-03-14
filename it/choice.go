package it

import (
	"strings"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/code"
	"github.com/muonsoft/validation/message"
)

// ChoiceConstraint is used to ensure that the given value corresponds to one of the expected choices.
type ChoiceConstraint struct {
	choices         map[string]bool
	choicesValue    string
	messageTemplate string
	isIgnored       bool
}

// IsOneOfStrings creates a ChoiceConstraint for checking that values are in the expected list of strings.
//
// Example
//	err := validator.ValidateString(&s, it.IsOneOfStrings("one", "two", "three))
func IsOneOfStrings(values ...string) ChoiceConstraint {
	choices := make(map[string]bool, len(values))
	for _, value := range values {
		choices[value] = true
	}

	return ChoiceConstraint{
		choices:         choices,
		choicesValue:    strings.Join(values, ", "),
		messageTemplate: message.NoSuchChoice,
	}
}

// SetUp will return an error if the list of choices is empty.
func (c ChoiceConstraint) SetUp(scope *validation.Scope) error {
	if len(c.choices) == 0 {
		return errEmptyChoices
	}

	return nil
}

// Name is the constraint name.
func (c ChoiceConstraint) Name() string {
	return "ChoiceConstraint"
}

// Message sets the violation message template. You can use template parameters
// for injecting its values into the final message:
//
//	{{ choices }} - a comma-separated list of available choices;
//	{{ value }} - the current (invalid) value.
func (c ChoiceConstraint) Message(message string) ChoiceConstraint {
	c.messageTemplate = message
	return c
}

// When enables conditional validation of this constraint. If the expression evaluates to false,
// then the constraint will be ignored.
func (c ChoiceConstraint) When(condition bool) ChoiceConstraint {
	c.isIgnored = !condition
	return c
}

func (c ChoiceConstraint) ValidateString(value *string, scope validation.Scope) error {
	if c.isIgnored || value == nil || *value == "" {
		return nil
	}
	if c.choices[*value] {
		return nil
	}

	return scope.
		BuildViolation(code.NoSuchChoice, c.messageTemplate).
		SetParameters(map[string]string{
			"{{ value }}":   *value,
			"{{ choices }}": c.choicesValue,
		}).
		GetViolation()
}
