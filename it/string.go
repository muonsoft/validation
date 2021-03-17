package it

import (
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/code"
	"github.com/muonsoft/validation/message"

	"regexp"
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

// RegexConstraint is used to ensure that the given value corresponds to regex pattern.
type RegexConstraint struct {
	isIgnored       bool
	match           bool
	messageTemplate string
	regex           *regexp.Regexp
}

// Matches creates a RegexConstraint for checking whether a value match.
//
// Example
//	err := validator.ValidateString(&s, it.Matches(regexp.MustCompile("^[a-z]+$")))
func Matches(regex *regexp.Regexp) RegexConstraint {
	return RegexConstraint{
		regex:           regex,
		match:           true,
		messageTemplate: message.NotValid,
	}
}

// Matches creates a RegexConstraint for checking whether a value not match.
//
// Example
//	err := validator.ValidateString(&s, it.DoesNotMatch(regexp.MustCompile("^[a-z]+$")))
func DoesNotMatch(regex *regexp.Regexp) RegexConstraint {
	return RegexConstraint{
		regex:           regex,
		match:           false,
		messageTemplate: message.NotValid,
	}
}

// SetUp will return an error if the pattern is empty.
func (c RegexConstraint) SetUp() error {
	if c.regex == nil {
		return errEmptyRegex
	}

	return nil
}

// Name is the constraint name.
func (c RegexConstraint) Name() string {
	return "RegexConstraint"
}

// Message sets the violation message template. You can use template parameters.
// for injecting its values into the final message:
//
//	{{ value }} - the current (invalid) value.
func (c RegexConstraint) Message(message string) RegexConstraint {
	c.messageTemplate = message
	return c
}

// When enables conditional validation of this constraint. If the expression evaluates to false,
// then the constraint will be ignored.
func (c RegexConstraint) When(condition bool) RegexConstraint {
	c.isIgnored = !condition
	return c
}

func (c RegexConstraint) ValidateString(value *string, scope validation.Scope) error {
	if c.isIgnored || value == nil || *value == "" {
		return nil
	}

	if !(c.match != c.regex.MatchString(*value)) {
		return nil
	}

	return scope.
		BuildViolation(code.MatchingFailed, c.messageTemplate).
		SetParameter("{{ value }}", *value).
		CreateViolation()
}
