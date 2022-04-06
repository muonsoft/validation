package it

import (
	"strconv"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/code"
	"github.com/muonsoft/validation/generic"
	"github.com/muonsoft/validation/message"
)

// CountConstraint checks that a given collection's (array, slice or a map) length is between some minimum and
// maximum value.
type CountConstraint struct {
	isIgnored              bool
	checkMin               bool
	checkMax               bool
	min                    int
	max                    int
	groups                 []string
	minCode                string
	maxCode                string
	exactCode              string
	minMessageTemplate     string
	minMessageParameters   validation.TemplateParameterList
	maxMessageTemplate     string
	maxMessageParameters   validation.TemplateParameterList
	exactMessageTemplate   string
	exactMessageParameters validation.TemplateParameterList
}

func newCountConstraint(min int, max int, checkMin bool, checkMax bool) CountConstraint {
	return CountConstraint{
		min:                  min,
		max:                  max,
		checkMin:             checkMin,
		checkMax:             checkMax,
		minCode:              code.CountTooFew,
		maxCode:              code.CountTooMany,
		exactCode:            code.CountExact,
		minMessageTemplate:   message.Templates[code.CountTooFew],
		maxMessageTemplate:   message.Templates[code.CountTooMany],
		exactMessageTemplate: message.Templates[code.CountExact],
	}
}

// HasMinCount creates a CountConstraint that checks the length of the iterable (slice, array, or map)
// is greater than the minimum value.
func HasMinCount(min int) CountConstraint {
	return newCountConstraint(min, 0, true, false)
}

// HasMaxCount creates a CountConstraint that checks the length of the iterable (slice, array, or map)
// is less than the maximum value.
func HasMaxCount(max int) CountConstraint {
	return newCountConstraint(0, max, false, true)
}

// HasCountBetween creates a CountConstraint that checks the length of the iterable (slice, array, or map)
// is between some minimum and maximum value.
func HasCountBetween(min int, max int) CountConstraint {
	return newCountConstraint(min, max, true, true)
}

// HasExactCount creates a CountConstraint that checks the length of the iterable (slice, array, or map)
// has exact value.
func HasExactCount(count int) CountConstraint {
	return newCountConstraint(count, count, true, true)
}

// SetUp always returns no error.
func (c CountConstraint) SetUp() error {
	return nil
}

// Name is the constraint name.
func (c CountConstraint) Name() string {
	return "CountConstraint"
}

// When enables conditional validation of this constraint. If the expression evaluates to false,
// then the constraint will be ignored.
func (c CountConstraint) When(condition bool) CountConstraint {
	c.isIgnored = !condition
	return c
}

// WhenGroups enables conditional validation of the constraint by using the validation groups.
func (c CountConstraint) WhenGroups(groups ...string) CountConstraint {
	c.groups = groups
	return c
}

// MinCode overrides default code for violation that will be shown if the
// collection length is less than the minimum value.
func (c CountConstraint) MinCode(code string) CountConstraint {
	c.minCode = code
	return c
}

// MaxCode overrides default code for violation that will be shown if the
// collection length is greater than the maximum value.
func (c CountConstraint) MaxCode(code string) CountConstraint {
	c.maxCode = code
	return c
}

// ExactCode overrides default code for violation that will be shown if minimum and
// maximum values are equal and the length of the collection is not exactly this value.
func (c CountConstraint) ExactCode(code string) CountConstraint {
	c.exactCode = code
	return c
}

// MinMessage sets the violation message that will be shown if the collection length is less than
// the minimum value. You can set custom template parameters for injecting its values
// into the final message. Also, you can use default parameters:
//
//	{{ count }} - the current collection size;
//	{{ limit }} - the lower limit.
func (c CountConstraint) MinMessage(template string, parameters ...validation.TemplateParameter) CountConstraint {
	c.minMessageTemplate = template
	c.minMessageParameters = parameters
	return c
}

// MaxMessage sets the violation message that will be shown if the collection length is greater than
// the maximum value. You can set custom template parameters for injecting its values
// into the final message. Also, you can use default parameters:
//
//	{{ count }} - the current collection size;
//	{{ limit }} - the upper limit.
func (c CountConstraint) MaxMessage(template string, parameters ...validation.TemplateParameter) CountConstraint {
	c.maxMessageTemplate = template
	c.maxMessageParameters = parameters
	return c
}

// ExactMessage sets the violation message that will be shown if minimum and maximum values are equal and
// the length of the collection is not exactly this value. You can set custom template parameters
// for injecting its values into the final message. Also, you can use default parameters:
//
//	{{ count }} - the current collection size;
//	{{ limit }} - the exact expected collection size.
func (c CountConstraint) ExactMessage(template string, parameters ...validation.TemplateParameter) CountConstraint {
	c.exactMessageTemplate = template
	c.exactMessageParameters = parameters
	return c
}

func (c CountConstraint) ValidateIterable(value generic.Iterable, scope validation.Scope) error {
	return c.ValidateCountable(value.Count(), scope)
}

func (c CountConstraint) ValidateCountable(count int, scope validation.Scope) error {
	if c.isIgnored || scope.IsIgnored(c.groups...) {
		return nil
	}
	if c.checkMax && count > c.max {
		return c.newViolation(count, c.max, c.maxCode, c.maxMessageTemplate, c.maxMessageParameters, scope)
	}
	if c.checkMin && count < c.min {
		return c.newViolation(count, c.min, c.minCode, c.minMessageTemplate, c.minMessageParameters, scope)
	}

	return nil
}

func (c CountConstraint) newViolation(
	count, limit int,
	violationCode, template string,
	parameters validation.TemplateParameterList,
	scope validation.Scope,
) validation.Violation {
	if c.checkMin && c.checkMax && c.min == c.max {
		template = c.exactMessageTemplate
		parameters = c.exactMessageParameters
		violationCode = c.exactCode
	}

	return scope.BuildViolation(violationCode, template).
		SetPluralCount(limit).
		SetParameters(
			parameters.Prepend(
				validation.TemplateParameter{Key: "{{ count }}", Value: strconv.Itoa(count)},
				validation.TemplateParameter{Key: "{{ limit }}", Value: strconv.Itoa(limit)},
			)...,
		).
		CreateViolation()
}
