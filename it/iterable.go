package it

import (
	"context"
	"strconv"

	"github.com/muonsoft/validation"
)

// CountConstraint checks that a given collection's (array, slice or a map) length is between some minimum and
// maximum value.
type CountConstraint struct {
	isIgnored                    bool
	checkMin                     bool
	checkMax                     bool
	checkDivisible               bool
	min                          int
	max                          int
	divisibleBy                  int
	groups                       []string
	minErr                       error
	maxErr                       error
	exactErr                     error
	divisibleErr                 error
	minMessageTemplate           string
	minMessageParameters         validation.TemplateParameterList
	maxMessageTemplate           string
	maxMessageParameters         validation.TemplateParameterList
	exactMessageTemplate         string
	exactMessageParameters       validation.TemplateParameterList
	divisibleByMessageTemplate   string
	divisibleByMessageParameters validation.TemplateParameterList
}

func newCountConstraint() CountConstraint {
	return CountConstraint{
		minErr:                     validation.ErrTooFewElements,
		maxErr:                     validation.ErrTooManyElements,
		exactErr:                   validation.ErrNotExactCount,
		divisibleErr:               validation.ErrNotDivisibleCount,
		minMessageTemplate:         validation.ErrTooFewElements.Message(),
		maxMessageTemplate:         validation.ErrTooManyElements.Message(),
		exactMessageTemplate:       validation.ErrNotExactCount.Message(),
		divisibleByMessageTemplate: validation.ErrNotDivisibleCount.Message(),
	}
}

func newCountComparison(min int, max int, checkMin bool, checkMax bool) CountConstraint {
	c := newCountConstraint()
	c.min = min
	c.max = max
	c.checkMin = checkMin
	c.checkMax = checkMax

	return c
}

// HasMinCount creates a [CountConstraint] that checks the length of the iterable (slice, array, or map)
// is greater than the minimum value.
func HasMinCount(min int) CountConstraint {
	return newCountComparison(min, 0, true, false)
}

// HasMaxCount creates a [CountConstraint] that checks the length of the iterable (slice, array, or map)
// is less than the maximum value.
func HasMaxCount(max int) CountConstraint {
	return newCountComparison(0, max, false, true)
}

// HasCountBetween creates a [CountConstraint] that checks the length of the iterable (slice, array, or map)
// is between some minimum and maximum value.
func HasCountBetween(min int, max int) CountConstraint {
	return newCountComparison(min, max, true, true)
}

// HasExactCount creates a [CountConstraint] that checks the length of the iterable (slice, array, or map)
// has exact value.
func HasExactCount(count int) CountConstraint {
	return newCountComparison(count, count, true, true)
}

// HasCountDivisibleBy creates a [CountConstraint] that checks the length of the iterable (slice, array, or map)
// is divisible by the specific value.
func HasCountDivisibleBy(divisor int) CountConstraint {
	c := newCountConstraint()
	c.checkDivisible = true
	c.divisibleBy = divisor

	return c
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

// WithDivisibleError overrides default underlying error for violation that will be shown
// the length of the collection is not divisible by specific value.
func (c CountConstraint) WithDivisibleError(err error) CountConstraint {
	c.divisibleErr = err
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

// WithDivisibleMessage sets the violation message that will be shown if the length of the collection
// is not divisible by specific value. You can set custom template parameters
// for injecting its values into the final message. Also, you can use default parameters:
//
//	{{ count }} - the current collection size;
//	{{ divisibleBy }} - the divisor for the collection size.
func (c CountConstraint) WithDivisibleMessage(template string, parameters ...validation.TemplateParameter) CountConstraint {
	c.divisibleByMessageTemplate = template
	c.divisibleByMessageParameters = parameters
	return c
}

func (c CountConstraint) ValidateCountable(ctx context.Context, validator *validation.Validator, count int) error {
	if c.isIgnored || validator.IsIgnoredForGroups(c.groups...) {
		return nil
	}
	if c.checkDivisible && count%c.divisibleBy != 0 {
		return c.newNotDivisibleViolation(ctx, validator, count)
	}
	if c.checkMax && count > c.max {
		return c.newViolation(ctx, validator, count, c.max, c.maxErr, c.maxMessageTemplate, c.maxMessageParameters)
	}
	if c.checkMin && count < c.min {
		return c.newViolation(ctx, validator, count, c.min, c.minErr, c.minMessageTemplate, c.minMessageParameters)
	}

	return nil
}

func (c CountConstraint) newViolation(
	ctx context.Context,
	validator *validation.Validator,
	count, limit int,
	err error,
	template string,
	parameters validation.TemplateParameterList,
) validation.Violation {
	if c.checkMin && c.checkMax && c.min == c.max {
		template = c.exactMessageTemplate
		parameters = c.exactMessageParameters
		err = c.exactErr
	}

	return validator.BuildViolation(ctx, err, template).
		WithPluralCount(limit).
		WithParameters(
			parameters.Prepend(
				validation.TemplateParameter{Key: "{{ count }}", Value: strconv.Itoa(count)},
				validation.TemplateParameter{Key: "{{ limit }}", Value: strconv.Itoa(limit)},
			)...,
		).
		Create()
}

func (c CountConstraint) newNotDivisibleViolation(
	ctx context.Context,
	validator *validation.Validator,
	count int,
) validation.Violation {
	return validator.BuildViolation(ctx, c.divisibleErr, c.divisibleByMessageTemplate).
		WithPluralCount(c.divisibleBy).
		WithParameters(
			c.divisibleByMessageParameters.Prepend(
				validation.TemplateParameter{Key: "{{ count }}", Value: strconv.Itoa(count)},
				validation.TemplateParameter{Key: "{{ divisibleBy }}", Value: strconv.Itoa(c.divisibleBy)},
			)...,
		).
		Create()
}
