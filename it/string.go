package it

import (
	"context"
	"regexp"
	"strconv"
	"unicode/utf8"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/is"
)

// LengthConstraint checks that a given string length is between some minimum and maximum value.
// If you want to check the length of the array, slice or a map use CountConstraint.
type LengthConstraint struct {
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

func newLengthConstraint(min int, max int, checkMin bool, checkMax bool) LengthConstraint {
	return LengthConstraint{
		min:                  min,
		max:                  max,
		checkMin:             checkMin,
		checkMax:             checkMax,
		minErr:               validation.ErrTooShort,
		maxErr:               validation.ErrTooLong,
		exactErr:             validation.ErrNotExactLength,
		minMessageTemplate:   validation.ErrTooShort.Template(),
		maxMessageTemplate:   validation.ErrTooLong.Template(),
		exactMessageTemplate: validation.ErrNotExactLength.Template(),
	}
}

// HasMinLength creates a LengthConstraint that checks the length of the string
// is greater than the minimum value.
func HasMinLength(min int) LengthConstraint {
	return newLengthConstraint(min, 0, true, false)
}

// HasMaxLength creates a LengthConstraint that checks the length of the string
// is less than the maximum value.
func HasMaxLength(max int) LengthConstraint {
	return newLengthConstraint(0, max, false, true)
}

// HasLengthBetween creates a LengthConstraint that checks the length of the string
// is between some minimum and maximum value.
func HasLengthBetween(min int, max int) LengthConstraint {
	return newLengthConstraint(min, max, true, true)
}

// HasExactLength creates a LengthConstraint that checks the length of the string
// has exact value.
func HasExactLength(count int) LengthConstraint {
	return newLengthConstraint(count, count, true, true)
}

// When enables conditional validation of this constraint. If the expression evaluates to false,
// then the constraint will be ignored.
func (c LengthConstraint) When(condition bool) LengthConstraint {
	c.isIgnored = !condition
	return c
}

// WhenGroups enables conditional validation of the constraint by using the validation groups.
func (c LengthConstraint) WhenGroups(groups ...string) LengthConstraint {
	c.groups = groups
	return c
}

// WithMinError overrides default underlying error for violation that will be shown if the string length
// is less than the minimum value.
func (c LengthConstraint) WithMinError(err error) LengthConstraint {
	c.minErr = err
	return c
}

// WithMaxError overrides default underlying error for violation that will be shown if the string length
// is greater than the maximum value.
func (c LengthConstraint) WithMaxError(err error) LengthConstraint {
	c.maxErr = err
	return c
}

// WithExactError overrides default underlying error for violation that will be shown if minimum and maximum values
// are equal and the length of the string is not exactly this value.
func (c LengthConstraint) WithExactError(err error) LengthConstraint {
	c.exactErr = err
	return c
}

// WithMinMessage sets the violation message that will be shown if the string length is less than
// the minimum value. You can set custom template parameters for injecting its values
// into the final message. Also, you can use default parameters:
//
//	{{ length }} - the current string length;
//	{{ limit }} - the lower limit;
//	{{ value }} - the current (invalid) value.
func (c LengthConstraint) WithMinMessage(template string, parameters ...validation.TemplateParameter) LengthConstraint {
	c.minMessageTemplate = template
	c.minMessageParameters = parameters
	return c
}

// WithMaxMessage sets the violation message that will be shown if the string length is greater than
// the maximum value. You can set custom template parameters for injecting its values
// into the final message. Also, you can use default parameters:
//
//	{{ length }} - the current string length;
//	{{ limit }} - the lower limit;
//	{{ value }} - the current (invalid) value.
func (c LengthConstraint) WithMaxMessage(template string, parameters ...validation.TemplateParameter) LengthConstraint {
	c.maxMessageTemplate = template
	c.maxMessageParameters = parameters
	return c
}

// WithExactMessage sets the violation message that will be shown if minimum and maximum values are equal and
// the length of the string is not exactly this value. You can set custom template parameters
// for injecting its values into the final message. Also, you can use default parameters:
//
//	{{ length }} - the current string length;
//	{{ limit }} - the lower limit;
//	{{ value }} - the current (invalid) value.
func (c LengthConstraint) WithExactMessage(template string, parameters ...validation.TemplateParameter) LengthConstraint {
	c.exactMessageTemplate = template
	c.exactMessageParameters = parameters
	return c
}

func (c LengthConstraint) ValidateString(ctx context.Context, validator *validation.Validator, value *string) error {
	if c.isIgnored || validator.IsIgnoredForGroups(c.groups...) || value == nil || *value == "" {
		return nil
	}

	count := utf8.RuneCountInString(*value)

	if c.checkMax && count > c.max {
		return c.newViolation(ctx, validator, count, c.max, *value, c.maxErr, c.maxMessageTemplate, c.maxMessageParameters)
	}
	if c.checkMin && count < c.min {
		return c.newViolation(ctx, validator, count, c.min, *value, c.minErr, c.minMessageTemplate, c.minMessageParameters)
	}

	return nil
}

func (c LengthConstraint) newViolation(
	ctx context.Context,
	validator *validation.Validator,
	count, limit int,
	value string,
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
				validation.TemplateParameter{Key: "{{ value }}", Value: strconv.Quote(value)},
				validation.TemplateParameter{Key: "{{ length }}", Value: strconv.Itoa(count)},
				validation.TemplateParameter{Key: "{{ limit }}", Value: strconv.Itoa(limit)},
			)...,
		).
		Create()
}

// RegexpConstraint is used to ensure that the given value corresponds to regex pattern.
type RegexpConstraint struct {
	isIgnored         bool
	match             bool
	groups            []string
	err               error
	messageTemplate   string
	messageParameters validation.TemplateParameterList
	regex             *regexp.Regexp
}

// Matches creates a RegexpConstraint for checking whether a value matches a regular expression.
func Matches(regex *regexp.Regexp) RegexpConstraint {
	return RegexpConstraint{
		regex:           regex,
		match:           true,
		err:             validation.ErrNotValid,
		messageTemplate: validation.ErrNotValid.Template(),
	}
}

// DoesNotMatch creates a RegexpConstraint for checking whether a value does not match a regular expression.
func DoesNotMatch(regex *regexp.Regexp) RegexpConstraint {
	return RegexpConstraint{
		regex:           regex,
		match:           false,
		err:             validation.ErrNotValid,
		messageTemplate: validation.ErrNotValid.Template(),
	}
}

// WithError overrides default error for produced violation.
func (c RegexpConstraint) WithError(err error) RegexpConstraint {
	c.err = err
	return c
}

// WithMessage sets the violation message template. You can set custom template parameters
// for injecting its values into the final message. Also, you can use default parameters:
//
//	{{ value }} - the current (invalid) value.
func (c RegexpConstraint) WithMessage(template string, parameters ...validation.TemplateParameter) RegexpConstraint {
	c.messageTemplate = template
	c.messageParameters = parameters
	return c
}

// When enables conditional validation of this constraint. If the expression evaluates to false,
// then the constraint will be ignored.
func (c RegexpConstraint) When(condition bool) RegexpConstraint {
	c.isIgnored = !condition
	return c
}

// WhenGroups enables conditional validation of the constraint by using the validation groups.
func (c RegexpConstraint) WhenGroups(groups ...string) RegexpConstraint {
	c.groups = groups
	return c
}

func (c RegexpConstraint) ValidateString(ctx context.Context, validator *validation.Validator, value *string) error {
	if c.regex == nil {
		return validator.CreateConstraintError("RegexpConstraint", "nil regex")
	}
	if c.isIgnored || validator.IsIgnoredForGroups(c.groups...) || value == nil || *value == "" {
		return nil
	}
	if c.match == c.regex.MatchString(*value) {
		return nil
	}

	return validator.
		BuildViolation(ctx, c.err, c.messageTemplate).
		WithParameters(
			c.messageParameters.Prepend(
				validation.TemplateParameter{Key: "{{ value }}", Value: *value},
			)...,
		).
		Create()
}

// IsJSON validates that a value is a valid JSON.
func IsJSON() validation.StringFuncConstraint {
	return validation.OfStringBy(is.JSON).
		WithError(validation.ErrInvalidJSON).
		WithMessage(validation.ErrInvalidJSON.Template())
}

// IsInteger checks that string value is an integer.
func IsInteger() validation.StringFuncConstraint {
	return validation.OfStringBy(is.Integer).
		WithError(validation.ErrNotInteger).
		WithMessage(validation.ErrNotInteger.Template())
}

// IsNumeric checks that string value is a valid numeric (integer or float).
func IsNumeric() validation.StringFuncConstraint {
	return validation.OfStringBy(is.Number).
		WithError(validation.ErrNotNumeric).
		WithMessage(validation.ErrNotNumeric.Template())
}
