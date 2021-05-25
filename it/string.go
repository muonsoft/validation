package it

import (
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/code"
	"github.com/muonsoft/validation/is"
	"github.com/muonsoft/validation/message"

	"regexp"
	"strconv"
	"unicode/utf8"
)

// LengthConstraint checks that a given string length is between some minimum and maximum value.
// If you want to check the length of the array, slice or a map use CountConstraint.
type LengthConstraint struct {
	isIgnored              bool
	checkMin               bool
	checkMax               bool
	min                    int
	max                    int
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

func newLengthConstraint(min int, max int, checkMin bool, checkMax bool) LengthConstraint {
	return LengthConstraint{
		min:                  min,
		max:                  max,
		checkMin:             checkMin,
		checkMax:             checkMax,
		minCode:              code.LengthTooFew,
		maxCode:              code.LengthTooMany,
		exactCode:            code.LengthExact,
		minMessageTemplate:   message.LengthTooFew,
		maxMessageTemplate:   message.LengthTooMany,
		exactMessageTemplate: message.LengthExact,
	}
}

// HasMinLength creates a LengthConstraint that checks the length of the string
// is greater than the minimum value.
//
// Example
//  v := "foo"
//  err := validator.ValidateString(&v, it.HasMinLength(5))
func HasMinLength(min int) LengthConstraint {
	return newLengthConstraint(min, 0, true, false)
}

// HasMaxLength creates a LengthConstraint that checks the length of the string
// is less than the maximum value.
//
// Example
//  v := "foo"
//  err := validator.ValidateString(&v, it.HasMaxLength(2))
func HasMaxLength(max int) LengthConstraint {
	return newLengthConstraint(0, max, false, true)
}

// HasLengthBetween creates a LengthConstraint that checks the length of the string
// is between some minimum and maximum value.
//
// Example
//  v := "foo"
//  err := validator.ValidateString(&v, it.HasLengthBetween(5, 10))
func HasLengthBetween(min int, max int) LengthConstraint {
	return newLengthConstraint(min, max, true, true)
}

// HasExactLength creates a LengthConstraint that checks the length of the string
// has exact value.
//
// Example
//  v := "foo"
//  err := validator.ValidateString(&v, it.HasExactLength(5))
func HasExactLength(count int) LengthConstraint {
	return newLengthConstraint(count, count, true, true)
}

// Name is the constraint name.
func (c LengthConstraint) Name() string {
	return "LengthConstraint"
}

// SetUp always returns no error.
func (c LengthConstraint) SetUp() error {
	return nil
}

// When enables conditional validation of this constraint. If the expression evaluates to false,
// then the constraint will be ignored.
func (c LengthConstraint) When(condition bool) LengthConstraint {
	c.isIgnored = !condition
	return c
}

// MinCode overrides default code for violation that will be shown if the string length
// is less than the minimum value.
func (c LengthConstraint) MinCode(code string) LengthConstraint {
	c.minCode = code
	return c
}

// MaxCode overrides default code for violation that will be shown if the string length
// is greater than the maximum value.
func (c LengthConstraint) MaxCode(code string) LengthConstraint {
	c.maxCode = code
	return c
}

// ExactCode overrides default code for violation that will be shown if minimum and maximum values
// are equal and the length of the string is not exactly this value.
func (c LengthConstraint) ExactCode(code string) LengthConstraint {
	c.exactCode = code
	return c
}

// MinMessage sets the violation message that will be shown if the string length is less than
// the minimum value. You can set custom template parameters for injecting its values
// into the final message. Also, you can use default parameters:
//
//	{{ length }} - the current string length;
//	{{ limit }} - the lower limit;
//	{{ value }} - the current (invalid) value.
func (c LengthConstraint) MinMessage(template string, parameters ...validation.TemplateParameter) LengthConstraint {
	c.minMessageTemplate = template
	c.minMessageParameters = parameters
	return c
}

// MaxMessage sets the violation message that will be shown if the string length is greater than
// the maximum value. You can set custom template parameters for injecting its values
// into the final message. Also, you can use default parameters:
//
//	{{ length }} - the current string length;
//	{{ limit }} - the lower limit;
//	{{ value }} - the current (invalid) value.
func (c LengthConstraint) MaxMessage(template string, parameters ...validation.TemplateParameter) LengthConstraint {
	c.maxMessageTemplate = template
	c.maxMessageParameters = parameters
	return c
}

// ExactMessage sets the violation message that will be shown if minimum and maximum values are equal and
// the length of the string is not exactly this value. You can set custom template parameters
// for injecting its values into the final message. Also, you can use default parameters:
//
//	{{ length }} - the current string length;
//	{{ limit }} - the lower limit;
//	{{ value }} - the current (invalid) value.
func (c LengthConstraint) ExactMessage(template string, parameters ...validation.TemplateParameter) LengthConstraint {
	c.exactMessageTemplate = template
	c.exactMessageParameters = parameters
	return c
}

func (c LengthConstraint) ValidateString(value *string, scope validation.Scope) error {
	if c.isIgnored || value == nil || *value == "" {
		return nil
	}

	count := utf8.RuneCountInString(*value)

	if c.checkMax && count > c.max {
		return c.newViolation(count, c.max, *value, c.maxCode, c.maxMessageTemplate, c.maxMessageParameters, scope)
	}
	if c.checkMin && count < c.min {
		return c.newViolation(count, c.min, *value, c.minCode, c.minMessageTemplate, c.minMessageParameters, scope)
	}

	return nil
}

func (c LengthConstraint) newViolation(
	count, limit int,
	value, violationCode, template string,
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
				validation.TemplateParameter{Key: "{{ value }}", Value: strconv.Quote(value)},
				validation.TemplateParameter{Key: "{{ length }}", Value: strconv.Itoa(count)},
				validation.TemplateParameter{Key: "{{ limit }}", Value: strconv.Itoa(limit)},
			)...,
		).
		CreateViolation()
}

// RegexConstraint is used to ensure that the given value corresponds to regex pattern.
type RegexConstraint struct {
	isIgnored         bool
	match             bool
	code              string
	messageTemplate   string
	messageParameters validation.TemplateParameterList
	regex             *regexp.Regexp
}

// Matches creates a RegexConstraint for checking whether a value matches a regular expression.
//
// Example
//	err := validator.ValidateString(&s, it.Matches(regexp.MustCompile("^[a-z]+$")))
func Matches(regex *regexp.Regexp) RegexConstraint {
	return RegexConstraint{
		regex:           regex,
		match:           true,
		code:            code.MatchingFailed,
		messageTemplate: message.NotValid,
	}
}

// DoesNotMatch creates a RegexConstraint for checking whether a value does not match a regular expression.
//
// Example
//	err := validator.ValidateString(&s, it.DoesNotMatch(regexp.MustCompile("^[a-z]+$")))
func DoesNotMatch(regex *regexp.Regexp) RegexConstraint {
	return RegexConstraint{
		regex:           regex,
		match:           false,
		code:            code.MatchingFailed,
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

// Code overrides default code for produced violation.
func (c RegexConstraint) Code(code string) RegexConstraint {
	c.code = code
	return c
}

// Message sets the violation message template. You can set custom template parameters
// for injecting its values into the final message. Also, you can use default parameters:
//
//	{{ value }} - the current (invalid) value.
func (c RegexConstraint) Message(template string, parameters ...validation.TemplateParameter) RegexConstraint {
	c.messageTemplate = template
	c.messageParameters = parameters
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
	if c.match == c.regex.MatchString(*value) {
		return nil
	}

	return scope.
		BuildViolation(c.code, c.messageTemplate).
		SetParameters(
			c.messageParameters.Prepend(
				validation.TemplateParameter{Key: "{{ value }}", Value: *value},
			)...,
		).
		CreateViolation()
}

// IsJSON validates that a value is a valid JSON.
func IsJSON() validation.CustomStringConstraint {
	return validation.NewCustomStringConstraint(
		is.JSON,
		"JSONConstraint",
		code.InvalidJSON,
		message.InvalidJSON,
	)
}
