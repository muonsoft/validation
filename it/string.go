package it

import (
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/code"
	"github.com/muonsoft/validation/message"

	"strconv"
	"unicode/utf8"
)

// LengthConstraint checks that a given string length is between some minimum and maximum value.
type LengthConstraint struct {
	isIgnored            bool
	checkMin             bool
	checkMax             bool
	min                  int
	max                  int
	minMessageTemplate   string
	maxMessageTemplate   string
	exactMessageTemplate string
}

func newLengthConstraint(min int, max int, checkMin bool, checkMax bool) LengthConstraint {
	return LengthConstraint{
		min:                  min,
		max:                  max,
		checkMin:             checkMin,
		checkMax:             checkMax,
		minMessageTemplate:   message.LengthTooFew,
		maxMessageTemplate:   message.LengthTooMany,
		exactMessageTemplate: message.LengthExact,
	}
}

func HasMinLength(min int) LengthConstraint {
	return newLengthConstraint(min, 0, true, false)
}

func HasMaxLength(max int) LengthConstraint {
	return newLengthConstraint(0, max, false, true)
}

func HasLengthBetween(min int, max int) LengthConstraint {
	return newLengthConstraint(min, max, true, true)
}

func HasExactLength(count int) LengthConstraint {
	return newLengthConstraint(count, count, true, true)
}

func (c LengthConstraint) SetUp() error {
	return nil
}

func (c LengthConstraint) Name() string {
	return "LengthConstraint"
}

func (c LengthConstraint) When(condition bool) LengthConstraint {
	c.isIgnored = !condition
	return c
}

func (c LengthConstraint) MinMessage(message string) LengthConstraint {
	c.minMessageTemplate = message
	return c
}

func (c LengthConstraint) MaxMessage(message string) LengthConstraint {
	c.maxMessageTemplate = message
	return c
}

func (c LengthConstraint) ExactMessage(message string) LengthConstraint {
	c.exactMessageTemplate = message
	return c
}

func (c LengthConstraint) ValidateString(value *string, scope validation.Scope) error {
	if c.isIgnored || value == nil || *value == "" {
		return nil
	}

	count := utf8.RuneCountInString(*value)

	if c.checkMax && count > c.max {
		return c.newViolation(count, c.max, code.LengthTooMany, c.maxMessageTemplate, scope)
	}
	if c.checkMin && count < c.min {
		return c.newViolation(count, c.min, code.LengthTooFew, c.minMessageTemplate, scope)
	}

	return nil
}

func (c LengthConstraint) newViolation(
	count, limit int,
	violationCode, message string,
	scope validation.Scope,
) validation.Violation {
	if c.checkMin && c.checkMax && c.min == c.max {
		message = c.exactMessageTemplate
		violationCode = code.LengthExact
	}

	return scope.BuildViolation(violationCode, message).
		SetPluralCount(limit).
		SetParameters(map[string]string{
			"{{ count }}": strconv.Itoa(count),
			"{{ limit }}": strconv.Itoa(limit),
		}).
		CreateViolation()
}
