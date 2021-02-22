package it

import (
	"strconv"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/code"
	"github.com/muonsoft/validation/generic"
	"github.com/muonsoft/validation/message"
)

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

func HasMinCount(min int) CountConstraint {
	return newCountConstraint(min, 0, true, false)
}

func HasMaxCount(max int) CountConstraint {
	return newCountConstraint(0, max, false, true)
}

func HasCountBetween(min int, max int) CountConstraint {
	return newCountConstraint(min, max, true, true)
}

func HasExactCount(count int) CountConstraint {
	return newCountConstraint(count, count, true, true)
}

func (c CountConstraint) When(condition bool) CountConstraint {
	c.isIgnored = !condition
	return c
}

func (c CountConstraint) MinMessage(message string) CountConstraint {
	c.minMessageTemplate = message
	return c
}

func (c CountConstraint) MaxMessage(message string) CountConstraint {
	c.maxMessageTemplate = message
	return c
}

func (c CountConstraint) ExactMessage(message string) CountConstraint {
	c.exactMessageTemplate = message
	return c
}

func (c CountConstraint) ValidateIterable(value generic.Iterable, options validation.Options) error {
	return c.ValidateCountable(value.Count(), options)
}

func (c CountConstraint) ValidateCountable(count int, options validation.Options) error {
	if c.isIgnored {
		return nil
	}
	if c.checkMax && count > c.max {
		return c.newViolation(count, c.max, code.CountTooMany, c.maxMessageTemplate, options)
	}
	if c.checkMin && count < c.min {
		return c.newViolation(count, c.min, code.CountTooFew, c.minMessageTemplate, options)
	}

	return nil
}

func (c CountConstraint) Set(options *validation.Options) error {
	options.Constraints = append(options.Constraints, c)

	return nil
}

func (c CountConstraint) GetName() string {
	return "CountConstraint"
}

func (c CountConstraint) newViolation(
	count, limit int,
	violationCode, message string,
	options validation.Options,
) validation.Violation {
	if c.checkMin && c.checkMax && c.min == c.max {
		message = c.exactMessageTemplate
		violationCode = code.CountExact
	}

	return options.BuildViolation(violationCode, message).
		SetParameters(map[string]string{
			"{{ count }}": strconv.Itoa(count),
			"{{ limit }}": strconv.Itoa(limit),
		}).
		GetViolation()
}
