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
	isIgnored            bool
	checkMin             bool
	checkMax             bool
	min                  int
	max                  int
	minMessageTemplate   string
	maxMessageTemplate   string
	exactMessageTemplate string
}

func newCountConstraint(min int, max int, checkMin bool, checkMax bool) CountConstraint {
	return CountConstraint{
		min:                  min,
		max:                  max,
		checkMin:             checkMin,
		checkMax:             checkMax,
		minMessageTemplate:   message.CountTooFew,
		maxMessageTemplate:   message.CountTooMany,
		exactMessageTemplate: message.CountExact,
	}
}

// HasMinCount creates a CountConstraint that checks the length of the iterable (slice, array, or map)
// is greater than the minimum value.
//
// Example
//  v := []int{1, 2}
//  err := validator.ValidateIterable(v, it.HasMinCount(3))
func HasMinCount(min int) CountConstraint {
	return newCountConstraint(min, 0, true, false)
}

// HasMaxCount creates a CountConstraint that checks the length of the iterable (slice, array, or map)
// is less than the maximum value.
//
// Example
//  v := []int{1, 2}
//  err := validator.ValidateIterable(v, it.HasMaxCount(1))
func HasMaxCount(max int) CountConstraint {
	return newCountConstraint(0, max, false, true)
}

// HasCountBetween creates a CountConstraint that checks the length of the iterable (slice, array, or map)
// is between some minimum and maximum value.
//
// Example
//  v := []int{1, 2}
//  err := validator.ValidateIterable(v, it.HasCountBetween(3, 10))
func HasCountBetween(min int, max int) CountConstraint {
	return newCountConstraint(min, max, true, true)
}

// HasExactCount creates a CountConstraint that checks the length of the iterable (slice, array, or map)
// has exact value.
//
// Example
//  v := []int{1, 2}
//  err := validator.ValidateIterable(v, it.HasExactCount(3))
func HasExactCount(count int) CountConstraint {
	return newCountConstraint(count, count, true, true)
}

// Name is the constraint name.
func (c CountConstraint) Name() string {
	return "CountConstraint"
}

// SetUp always returns no error.
func (c CountConstraint) SetUp() error {
	return nil
}

// When enables conditional validation of this constraint. If the expression evaluates to false,
// then the constraint will be ignored.
func (c CountConstraint) When(condition bool) CountConstraint {
	c.isIgnored = !condition
	return c
}

// MinMessage sets the violation message that will be shown if the collection length is less than
// the minimum value.
func (c CountConstraint) MinMessage(message string) CountConstraint {
	c.minMessageTemplate = message
	return c
}

// MaxMessage sets the violation message that will be shown if the collection length is greater than
// the maximum value.
func (c CountConstraint) MaxMessage(message string) CountConstraint {
	c.maxMessageTemplate = message
	return c
}

// ExactMessage sets the violation message that will be shown if minimum and maximum values are equal and
// the length of the collection is not exactly this value.
func (c CountConstraint) ExactMessage(message string) CountConstraint {
	c.exactMessageTemplate = message
	return c
}

func (c CountConstraint) ValidateIterable(value generic.Iterable, scope validation.Scope) error {
	return c.ValidateCountable(value.Count(), scope)
}

func (c CountConstraint) ValidateCountable(count int, scope validation.Scope) error {
	if c.isIgnored {
		return nil
	}
	if c.checkMax && count > c.max {
		return c.newViolation(count, c.max, code.CountTooMany, c.maxMessageTemplate, scope)
	}
	if c.checkMin && count < c.min {
		return c.newViolation(count, c.min, code.CountTooFew, c.minMessageTemplate, scope)
	}

	return nil
}

func (c CountConstraint) newViolation(
	count, limit int,
	violationCode, message string,
	scope validation.Scope,
) validation.Violation {
	if c.checkMin && c.checkMax && c.min == c.max {
		message = c.exactMessageTemplate
		violationCode = code.CountExact
	}

	return scope.BuildViolation(violationCode, message).
		SetPluralCount(limit).
		SetParameters(map[string]string{
			"{{ count }}": strconv.Itoa(count),
			"{{ limit }}": strconv.Itoa(limit),
		}).
		CreateViolation()
}
