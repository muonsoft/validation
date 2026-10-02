package it

import (
	"context"
	"regexp"
	"strconv"

	"github.com/muonsoft/validation"
)

// WordCountConstraint checks inclusive bounds on the number of words in a string.
// A word starts with a Unicode letter or number and continues with letters,
// numbers, or combining marks. Internal ASCII and curly apostrophes join words;
// other punctuation (including hyphens and underscores) separates them.
// Symbols and emoji do not count as words. Adjacent letters in scripts without
// spaces count as one word: this is not locale-aware ICU word segmentation.
// Nil and empty strings are ignored; whitespace-only strings have zero words.
// Use [IsNotBlank] as well when a value is required.
// Bounds must be non-negative, and minimum must not exceed maximum;
// invalid bounds produce a constraint configuration error during validation.
type WordCountConstraint struct {
	isIgnored            bool
	checkMin             bool
	checkMax             bool
	min                  int
	max                  int
	groups               []string
	minErr               error
	maxErr               error
	minMessageTemplate   string
	minMessageParameters validation.TemplateParameterList
	maxMessageTemplate   string
	maxMessageParameters validation.TemplateParameterList
}

var wordPattern = regexp.MustCompile(`[\pL\pN][\pL\pN\pM]*(?:['’][\pL\pN][\pL\pN\pM]*)*`)

func newWordCountConstraint(vMin int, vMax int, checkMin bool, checkMax bool) WordCountConstraint {
	return WordCountConstraint{
		min:                vMin,
		max:                vMax,
		checkMin:           checkMin,
		checkMax:           checkMax,
		minErr:             validation.ErrTooFewWords,
		maxErr:             validation.ErrTooManyWords,
		minMessageTemplate: validation.ErrTooFewWords.Message(),
		maxMessageTemplate: validation.ErrTooManyWords.Message(),
	}
}

// HasMinWordCount creates a [WordCountConstraint] that checks the number of words
// is at least the minimum value.
func HasMinWordCount(vMin int) WordCountConstraint {
	return newWordCountConstraint(vMin, 0, true, false)
}

// HasMaxWordCount creates a [WordCountConstraint] that checks the number of words
// is at most the maximum value.
func HasMaxWordCount(vMax int) WordCountConstraint {
	return newWordCountConstraint(0, vMax, false, true)
}

// HasWordCountBetween creates a [WordCountConstraint] that checks the number of words
// is between the inclusive minimum and maximum values.
func HasWordCountBetween(vMin int, vMax int) WordCountConstraint {
	return newWordCountConstraint(vMin, vMax, true, true)
}

// When enables conditional validation of this constraint. If the expression evaluates to false,
// then the constraint will be ignored.
func (c WordCountConstraint) When(condition bool) WordCountConstraint {
	c.isIgnored = !condition
	return c
}

// WhenGroups enables conditional validation of the constraint by using the validation groups.
func (c WordCountConstraint) WhenGroups(groups ...string) WordCountConstraint {
	c.groups = groups
	return c
}

// WithMinError overrides default underlying error for violation that will be shown if the word count
// is less than the minimum value.
func (c WordCountConstraint) WithMinError(err error) WordCountConstraint {
	c.minErr = err
	return c
}

// WithMaxError overrides default underlying error for violation that will be shown if the word count
// is greater than the maximum value.
func (c WordCountConstraint) WithMaxError(err error) WordCountConstraint {
	c.maxErr = err
	return c
}

// WithMinMessage sets the violation message that will be shown if the word count is less than
// the minimum value. You can set custom template parameters for injecting its values
// into the final message. Also, you can use default parameters:
//
//	{{ count }} - the current word count;
//	{{ limit }} - the bound;
//	{{ value }} - the current (invalid) value.
func (c WordCountConstraint) WithMinMessage(template string, parameters ...validation.TemplateParameter) WordCountConstraint {
	c.minMessageTemplate = template
	c.minMessageParameters = parameters
	return c
}

// WithMaxMessage sets the violation message that will be shown if the word count is greater than
// the maximum value. You can set custom template parameters for injecting its values
// into the final message. Also, you can use default parameters:
//
//	{{ count }} - the current word count;
//	{{ limit }} - the bound;
//	{{ value }} - the current (invalid) value.
func (c WordCountConstraint) WithMaxMessage(template string, parameters ...validation.TemplateParameter) WordCountConstraint {
	c.maxMessageTemplate = template
	c.maxMessageParameters = parameters
	return c
}

// ValidateString validates a string pointer against the configured word count bounds.
func (c WordCountConstraint) ValidateString(ctx context.Context, validator *validation.Validator, value *string) error {
	if c.hasInvalidBounds() {
		return validator.CreateConstraintError("WordCountConstraint", "bounds must be non-negative and minimum must not exceed maximum")
	}
	if c.isIgnored || validator.IsIgnoredForGroups(c.groups...) || value == nil || *value == "" {
		return nil
	}

	count := len(wordPattern.FindAllStringIndex(*value, -1))

	if c.checkMax && count > c.max {
		return c.newViolation(ctx, validator, count, c.max, *value, c.maxErr, c.maxMessageTemplate, c.maxMessageParameters)
	}
	if c.checkMin && count < c.min {
		return c.newViolation(ctx, validator, count, c.min, *value, c.minErr, c.minMessageTemplate, c.minMessageParameters)
	}

	return nil
}

func (c WordCountConstraint) hasInvalidBounds() bool {
	return c.min < 0 || c.max < 0 || (c.checkMin && c.checkMax && c.min > c.max)
}

func (c WordCountConstraint) newViolation(
	ctx context.Context,
	validator *validation.Validator,
	count, limit int,
	value string,
	err error,
	template string,
	parameters validation.TemplateParameterList,
) validation.Violation {
	return validator.BuildViolation(ctx, err, template).
		WithPluralCount(limit).
		WithParameters(
			parameters.Prepend(
				validation.TemplateParameter{Key: "{{ value }}", Value: strconv.Quote(value)},
				validation.TemplateParameter{Key: "{{ count }}", Value: strconv.Itoa(count)},
				validation.TemplateParameter{Key: "{{ limit }}", Value: strconv.Itoa(limit)},
			)...,
		).
		Create()
}

// Validate implements [validation.Constraint][string] so the constraint can be used with [validation.Each] and [validation.This].
func (c WordCountConstraint) Validate(ctx context.Context, validator *validation.Validator, v string) error {
	return c.ValidateString(ctx, validator, &v)
}
