package it

import (
	"strconv"

	"github.com/muonsoft/validation"
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
	minErr                 error
	maxErr                 error
	exactErr               error
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
		minErr:               validation.ErrTooFewElements,
		maxErr:               validation.ErrTooManyElements,
		exactErr:             validation.ErrNotExactCount,
		minMessageTemplate:   validation.ErrTooFewElements.Template(),
		maxMessageTemplate:   validation.ErrTooManyElements.Template(),
		exactMessageTemplate: validation.ErrNotExactCount.Template(),
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

// WithMinError overrides default underlying error for violation that will be shown if the
// collection length is less than the minimum value.
func (c CountConstraint) WithMinError(err error) CountConstraint {
	c.minErr = err
	return c
}

// WithMaxError overrides default underlying error for violation that will be shown if the
// collection length is greater than the maximum value.
func (c CountConstraint) WithMaxError(err error) CountConstraint {
	c.maxErr = err
	return c
}

// WithExactError overrides default underlying error for violation that will be shown if minimum and
// maximum values are equal and the length of the collection is not exactly this value.
func (c CountConstraint) WithExactError(err error) CountConstraint {
	c.exactErr = err
	return c
}

// WithMinMessage sets the violation message that will be shown if the collection length is less than
// the minimum value. You can set custom template parameters for injecting its values
// into the final message. Also, you can use default parameters:
//
//	{{ count }} - the current collection size;
//	{{ limit }} - the lower limit.
func (c CountConstraint) WithMinMessage(template string, parameters ...validation.TemplateParameter) CountConstraint {
	c.minMessageTemplate = template
	c.minMessageParameters = parameters
	return c
}

// WithMaxMessage sets the violation message that will be shown if the collection length is greater than
// the maximum value. You can set custom template parameters for injecting its values
// into the final message. Also, you can use default parameters:
//
//	{{ count }} - the current collection size;
//	{{ limit }} - the upper limit.
func (c CountConstraint) WithMaxMessage(template string, parameters ...validation.TemplateParameter) CountConstraint {
	c.maxMessageTemplate = template
	c.maxMessageParameters = parameters
	return c
}

// WithExactMessage sets the violation message that will be shown if minimum and maximum values are equal and
// the length of the collection is not exactly this value. You can set custom template parameters
// for injecting its values into the final message. Also, you can use default parameters:
//
//	{{ count }} - the current collection size;
//	{{ limit }} - the exact expected collection size.
func (c CountConstraint) WithExactMessage(template string, parameters ...validation.TemplateParameter) CountConstraint {
	c.exactMessageTemplate = template
	c.exactMessageParameters = parameters
	return c
}

func (c CountConstraint) ValidateCountable(count int, scope validation.Scope) error {
	if c.isIgnored || scope.IsIgnored(c.groups...) {
		return nil
	}
	if c.checkMax && count > c.max {
		return c.newViolation(count, c.max, c.maxErr, c.maxMessageTemplate, c.maxMessageParameters, scope)
	}
	if c.checkMin && count < c.min {
		return c.newViolation(count, c.min, c.minErr, c.minMessageTemplate, c.minMessageParameters, scope)
	}

	return nil
}

func (c CountConstraint) newViolation(
	count, limit int,
	err error,
	template string,
	parameters validation.TemplateParameterList,
	scope validation.Scope,
) validation.Violation {
	if c.checkMin && c.checkMax && c.min == c.max {
		template = c.exactMessageTemplate
		parameters = c.exactMessageParameters
		err = c.exactErr
	}

	return scope.BuildViolation(err, template).
		WithPluralCount(limit).
		WithParameters(
			parameters.Prepend(
				validation.TemplateParameter{Key: "{{ count }}", Value: strconv.Itoa(count)},
				validation.TemplateParameter{Key: "{{ limit }}", Value: strconv.Itoa(limit)},
			)...,
		).
		Create()
}
