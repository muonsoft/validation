package it

import (
	"strconv"
	"time"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/code"
	"github.com/muonsoft/validation/generic"
	"github.com/muonsoft/validation/is"
	"github.com/muonsoft/validation/message"
)

// NumberComparisonConstraint is used for various numeric comparisons between integer and float values.
// Values are compared as integers if the compared and specified values are integers.
// Otherwise, numbers are always compared as floating point numbers.
type NumberComparisonConstraint struct {
	isIgnored         bool
	groups            []string
	code              string
	messageTemplate   string
	messageParameters validation.TemplateParameterList
	comparedValue     string
	isValid           func(value generic.Number) bool
}

// IsEqualToInteger checks that the number (integer or float) is equal to the specified integer value.
// Values are compared as integers if the compared and specified values are integers.
// Otherwise, numbers are always compared as floating point numbers.
func IsEqualToInteger(value int64) NumberComparisonConstraint {
	v := generic.NewNumberFromInt(value)

	return NumberComparisonConstraint{
		code:            code.Equal,
		messageTemplate: message.Templates[code.Equal],
		comparedValue:   v.String(),
		isValid: func(n generic.Number) bool {
			return n.IsEqualTo(v)
		},
	}
}

// IsEqualToFloat checks that the number (integer or float) is equal to the specified float value.
// Values are compared as integers if the compared and specified values are integers.
// Otherwise, numbers are always compared as floating point numbers.
func IsEqualToFloat(value float64) NumberComparisonConstraint {
	v := generic.NewNumberFromFloat(value)

	return NumberComparisonConstraint{
		code:            code.Equal,
		messageTemplate: message.Templates[code.Equal],
		comparedValue:   v.String(),
		isValid: func(n generic.Number) bool {
			return n.IsEqualTo(v)
		},
	}
}

// IsNotEqualToInteger checks that the number (integer or float) is not equal to the specified integer value.
// Values are compared as integers if the compared and specified values are integers.
// Otherwise, numbers are always compared as floating point numbers.
func IsNotEqualToInteger(value int64) NumberComparisonConstraint {
	v := generic.NewNumberFromInt(value)

	return NumberComparisonConstraint{
		code:            code.NotEqual,
		messageTemplate: message.Templates[code.NotEqual],
		comparedValue:   v.String(),
		isValid: func(n generic.Number) bool {
			return !n.IsEqualTo(v)
		},
	}
}

// IsNotEqualToFloat checks that the number (integer or float) is not equal to the specified float value.
// Values are compared as integers if the compared and specified values are integers.
// Otherwise, numbers are always compared as floating point numbers.
func IsNotEqualToFloat(value float64) NumberComparisonConstraint {
	v := generic.NewNumberFromFloat(value)

	return NumberComparisonConstraint{
		code:            code.NotEqual,
		messageTemplate: message.Templates[code.NotEqual],
		comparedValue:   v.String(),
		isValid: func(n generic.Number) bool {
			return !n.IsEqualTo(v)
		},
	}
}

// IsLessThanInteger checks that the number (integer or float) is less than the specified integer value.
// Values are compared as integers if the compared and specified values are integers.
// Otherwise, numbers are always compared as floating point numbers.
func IsLessThanInteger(value int64) NumberComparisonConstraint {
	v := generic.NewNumberFromInt(value)

	return NumberComparisonConstraint{
		code:            code.TooHigh,
		messageTemplate: message.Templates[code.TooHigh],
		comparedValue:   v.String(),
		isValid: func(n generic.Number) bool {
			return n.IsLessThan(v)
		},
	}
}

// IsLessThanFloat checks that the number (integer or float) is less than the specified float value.
// Values are compared as integers if the compared and specified values are integers.
// Otherwise, numbers are always compared as floating point numbers.
func IsLessThanFloat(value float64) NumberComparisonConstraint {
	v := generic.NewNumberFromFloat(value)

	return NumberComparisonConstraint{
		code:            code.TooHigh,
		messageTemplate: message.Templates[code.TooHigh],
		comparedValue:   v.String(),
		isValid: func(n generic.Number) bool {
			return n.IsLessThan(v)
		},
	}
}

// IsLessThanOrEqualInteger checks that the number (integer or float) is less than or
// equal to the specified integer value. Values are compared as integers if the compared
// and specified values are integers. Otherwise, numbers are always compared as floating point numbers.
func IsLessThanOrEqualInteger(value int64) NumberComparisonConstraint {
	v := generic.NewNumberFromInt(value)

	return NumberComparisonConstraint{
		code:            code.TooHighOrEqual,
		messageTemplate: message.Templates[code.TooHighOrEqual],
		comparedValue:   v.String(),
		isValid: func(n generic.Number) bool {
			return n.IsLessThan(v) || n.IsEqualTo(v)
		},
	}
}

// IsLessThanOrEqualFloat checks that the number (integer or float) is less than or
// equal to the specified float value. Values are compared as integers if the compared
// and specified values are integers. Otherwise, numbers are always compared as floating point numbers.
func IsLessThanOrEqualFloat(value float64) NumberComparisonConstraint {
	v := generic.NewNumberFromFloat(value)

	return NumberComparisonConstraint{
		code:            code.TooHighOrEqual,
		messageTemplate: message.Templates[code.TooHighOrEqual],
		comparedValue:   v.String(),
		isValid: func(n generic.Number) bool {
			return n.IsLessThan(v) || n.IsEqualTo(v)
		},
	}
}

// IsGreaterThanInteger checks that the number (integer or float) is greater than the specified integer value.
// Values are compared as integers if the compared and specified values are integers.
// Otherwise, numbers are always compared as floating point numbers.
func IsGreaterThanInteger(value int64) NumberComparisonConstraint {
	v := generic.NewNumberFromInt(value)

	return NumberComparisonConstraint{
		code:            code.TooLow,
		messageTemplate: message.Templates[code.TooLow],
		comparedValue:   v.String(),
		isValid: func(n generic.Number) bool {
			return n.IsGreaterThan(v)
		},
	}
}

// IsGreaterThanFloat checks that the number (integer or float) is greater than the specified float value.
// Values are compared as integers if the compared and specified values are integers.
// Otherwise, numbers are always compared as floating point numbers.
func IsGreaterThanFloat(value float64) NumberComparisonConstraint {
	v := generic.NewNumberFromFloat(value)

	return NumberComparisonConstraint{
		code:            code.TooLow,
		messageTemplate: message.Templates[code.TooLow],
		comparedValue:   v.String(),
		isValid: func(n generic.Number) bool {
			return n.IsGreaterThan(v)
		},
	}
}

// IsGreaterThanOrEqualInteger checks that the number (integer or float) is greater than or
// equal to the specified integer value. Values are compared as integers if the compared
// and specified values are integers. Otherwise, numbers are always compared as floating point numbers.
func IsGreaterThanOrEqualInteger(value int64) NumberComparisonConstraint {
	v := generic.NewNumberFromInt(value)

	return NumberComparisonConstraint{
		code:            code.TooLowOrEqual,
		messageTemplate: message.Templates[code.TooLowOrEqual],
		comparedValue:   v.String(),
		isValid: func(n generic.Number) bool {
			return n.IsGreaterThan(v) || n.IsEqualTo(v)
		},
	}
}

// IsGreaterThanOrEqualFloat checks that the number (integer or float) is greater than or
// equal to the specified float value. Values are compared as integers if the compared
// and specified values are integers. Otherwise, numbers are always compared as floating point numbers.
func IsGreaterThanOrEqualFloat(value float64) NumberComparisonConstraint {
	v := generic.NewNumberFromFloat(value)

	return NumberComparisonConstraint{
		code:            code.TooLowOrEqual,
		messageTemplate: message.Templates[code.TooLowOrEqual],
		comparedValue:   v.String(),
		isValid: func(n generic.Number) bool {
			return n.IsGreaterThan(v) || n.IsEqualTo(v)
		},
	}
}

// IsPositive checks that the value is a positive number (integer or float). Zero is neither
// positive nor negative. If you want to allow zero use IsPositiveOrZero comparison.
func IsPositive() NumberComparisonConstraint {
	v := generic.NewNumberFromInt(0)

	return NumberComparisonConstraint{
		code:            code.NotPositive,
		messageTemplate: message.Templates[code.NotPositive],
		comparedValue:   v.String(),
		isValid: func(n generic.Number) bool {
			return n.IsGreaterThan(v)
		},
	}
}

// IsPositiveOrZero checks that the value is a positive number (integer or float) or equal to zero.
// If you don't want to allow zero as a valid value, use IsPositive comparison.
func IsPositiveOrZero() NumberComparisonConstraint {
	v := generic.NewNumberFromInt(0)

	return NumberComparisonConstraint{
		code:            code.NotPositiveOrZero,
		messageTemplate: message.Templates[code.NotPositiveOrZero],
		comparedValue:   v.String(),
		isValid: func(n generic.Number) bool {
			return n.IsGreaterThan(v) || n.IsEqualTo(v)
		},
	}
}

// IsNegative checks that the value is a negative number (integer or float). Zero is neither
// positive nor negative. If you want to allow zero use IsNegativeOrZero comparison.
func IsNegative() NumberComparisonConstraint {
	v := generic.NewNumberFromInt(0)

	return NumberComparisonConstraint{
		code:            code.NotNegative,
		messageTemplate: message.Templates[code.NotNegative],
		comparedValue:   v.String(),
		isValid: func(n generic.Number) bool {
			return n.IsLessThan(v)
		},
	}
}

// IsNegativeOrZero checks that the value is a negative number (integer or float) or equal to zero.
// If you don't want to allow zero as a valid value, use IsNegative comparison.
func IsNegativeOrZero() NumberComparisonConstraint {
	v := generic.NewNumberFromInt(0)

	return NumberComparisonConstraint{
		code:            code.NotNegativeOrZero,
		messageTemplate: message.Templates[code.NotNegativeOrZero],
		comparedValue:   v.String(),
		isValid: func(n generic.Number) bool {
			return n.IsLessThan(v) || n.IsEqualTo(v)
		},
	}
}

// SetUp always returns no error.
func (c NumberComparisonConstraint) SetUp() error {
	return nil
}

// Name is the constraint name.
func (c NumberComparisonConstraint) Name() string {
	return "NumberComparisonConstraint"
}

// Code overrides default code for produced violation.
func (c NumberComparisonConstraint) Code(code string) NumberComparisonConstraint {
	c.code = code
	return c
}

// Message sets the violation message template. You can set custom template parameters
// for injecting its values into the final message. Also, you can use default parameters:
//
//  {{ comparedValue }} - the expected value;
//  {{ value }} - the current (invalid) value.
func (c NumberComparisonConstraint) Message(
	template string,
	parameters ...validation.TemplateParameter,
) NumberComparisonConstraint {
	c.messageTemplate = template
	c.messageParameters = parameters
	return c
}

// When enables conditional validation of this constraint. If the expression evaluates to false,
// then the constraint will be ignored.
func (c NumberComparisonConstraint) When(condition bool) NumberComparisonConstraint {
	c.isIgnored = !condition
	return c
}

// WhenGroups enables conditional validation of the constraint by using the validation groups.
func (c NumberComparisonConstraint) WhenGroups(groups ...string) NumberComparisonConstraint {
	c.groups = groups
	return c
}

func (c NumberComparisonConstraint) ValidateNumber(value generic.Number, scope validation.Scope) error {
	if c.isIgnored || scope.IsIgnored(c.groups...) || value.IsNil() || c.isValid(value) {
		return nil
	}

	return scope.BuildViolation(c.code, c.messageTemplate).
		SetParameters(
			c.messageParameters.Prepend(
				validation.TemplateParameter{Key: "{{ comparedValue }}", Value: c.comparedValue},
				validation.TemplateParameter{Key: "{{ value }}", Value: value.String()},
			)...,
		).
		CreateViolation()
}

// RangeConstraint is used to check that a given number value is between some minimum and maximum.
// Values are compared as integers if the compared and specified values are integers.
// Otherwise, numbers are always compared as floating point numbers.
type RangeConstraint struct {
	isIgnored         bool
	groups            []string
	code              string
	messageTemplate   string
	messageParameters validation.TemplateParameterList
	min               generic.Number
	max               generic.Number
}

// IsBetweenIntegers checks that the number (integer or float) is between specified minimum and
// maximum integer values. Values are compared as integers if the compared and specified
// values are integers. Otherwise, numbers are always compared as floating point numbers.
func IsBetweenIntegers(min, max int64) RangeConstraint {
	return RangeConstraint{
		min:             generic.NewNumberFromInt(min),
		max:             generic.NewNumberFromInt(max),
		code:            code.NotInRange,
		messageTemplate: message.Templates[code.NotInRange],
	}
}

// IsBetweenFloats checks that the number (integer or float) is between specified minimum and
// maximum float values. Values are compared as integers if the compared and specified
// values are integers. Otherwise, numbers are always compared as floating point numbers.
func IsBetweenFloats(min, max float64) RangeConstraint {
	return RangeConstraint{
		min:             generic.NewNumberFromFloat(min),
		max:             generic.NewNumberFromFloat(max),
		code:            code.NotInRange,
		messageTemplate: message.Templates[code.NotInRange],
	}
}

// SetUp returns an error if min is greater than or equal to max.
func (c RangeConstraint) SetUp() error {
	if c.min.IsGreaterThan(c.max) || c.min.IsEqualTo(c.max) {
		return errInvalidRange
	}

	return nil
}

// Name is the constraint name.
func (c RangeConstraint) Name() string {
	return "RangeConstraint"
}

// Code overrides default code for produced violation.
func (c RangeConstraint) Code(code string) RangeConstraint {
	c.code = code
	return c
}

// Message sets the violation message template. You can set custom template parameters
// for injecting its values into the final message. Also, you can use default parameters:
//
//  {{ max }} - the upper limit;
//  {{ min }} - the lower limit;
//  {{ value }} - the current (invalid) value.
func (c RangeConstraint) Message(template string, parameters ...validation.TemplateParameter) RangeConstraint {
	c.messageTemplate = template
	c.messageParameters = parameters
	return c
}

// When enables conditional validation of this constraint. If the expression evaluates to false,
// then the constraint will be ignored.
func (c RangeConstraint) When(condition bool) RangeConstraint {
	c.isIgnored = !condition
	return c
}

// WhenGroups enables conditional validation of the constraint by using the validation groups.
func (c RangeConstraint) WhenGroups(groups ...string) RangeConstraint {
	c.groups = groups
	return c
}

func (c RangeConstraint) ValidateNumber(value generic.Number, scope validation.Scope) error {
	if c.isIgnored || scope.IsIgnored(c.groups...) {
		return nil
	}
	if value.IsLessThan(c.min) || value.IsGreaterThan(c.max) {
		return c.newViolation(value, scope)
	}

	return nil
}

func (c RangeConstraint) newViolation(value generic.Number, scope validation.Scope) error {
	return scope.BuildViolation(c.code, c.messageTemplate).
		SetParameters(
			c.messageParameters.Prepend(
				validation.TemplateParameter{Key: "{{ min }}", Value: c.min.String()},
				validation.TemplateParameter{Key: "{{ max }}", Value: c.max.String()},
				validation.TemplateParameter{Key: "{{ value }}", Value: value.String()},
			)...,
		).
		CreateViolation()
}

// StringComparisonConstraint is used to compare strings.
type StringComparisonConstraint struct {
	isIgnored         bool
	groups            []string
	code              string
	messageTemplate   string
	messageParameters validation.TemplateParameterList
	comparedValue     string
	isValid           func(value string) bool
}

// IsEqualToString checks that the string value is equal to the specified string value.
func IsEqualToString(value string) StringComparisonConstraint {
	return StringComparisonConstraint{
		code:            code.Equal,
		messageTemplate: message.Templates[code.Equal],
		comparedValue:   value,
		isValid: func(actualValue string) bool {
			return value == actualValue
		},
	}
}

// IsNotEqualToString checks that the string value is not equal to the specified string value.
func IsNotEqualToString(value string) StringComparisonConstraint {
	return StringComparisonConstraint{
		code:            code.NotEqual,
		messageTemplate: message.Templates[code.NotEqual],
		comparedValue:   value,
		isValid: func(actualValue string) bool {
			return value != actualValue
		},
	}
}

// SetUp always returns no error.
func (c StringComparisonConstraint) SetUp() error {
	return nil
}

// Name is the constraint name.
func (c StringComparisonConstraint) Name() string {
	return "StringComparisonConstraint"
}

// Code overrides default code for produced violation.
func (c StringComparisonConstraint) Code(code string) StringComparisonConstraint {
	c.code = code
	return c
}

// Message sets the violation message template. You can set custom template parameters
// for injecting its values into the final message. Also, you can use default parameters:
//
//  {{ comparedValue }} - the expected value;
//  {{ value }} - the current (invalid) value.
//
// All string values are quoted strings.
func (c StringComparisonConstraint) Message(
	template string,
	parameters ...validation.TemplateParameter,
) StringComparisonConstraint {
	c.messageTemplate = template
	c.messageParameters = parameters
	return c
}

// When enables conditional validation of this constraint. If the expression evaluates to false,
// then the constraint will be ignored.
func (c StringComparisonConstraint) When(condition bool) StringComparisonConstraint {
	c.isIgnored = !condition
	return c
}

// WhenGroups enables conditional validation of the constraint by using the validation groups.
func (c StringComparisonConstraint) WhenGroups(groups ...string) StringComparisonConstraint {
	c.groups = groups
	return c
}

func (c StringComparisonConstraint) ValidateString(value *string, scope validation.Scope) error {
	if c.isIgnored || scope.IsIgnored(c.groups...) || value == nil || c.isValid(*value) {
		return nil
	}

	return scope.BuildViolation(c.code, c.messageTemplate).
		SetParameters(
			c.messageParameters.Prepend(
				validation.TemplateParameter{Key: "{{ comparedValue }}", Value: strconv.Quote(c.comparedValue)},
				validation.TemplateParameter{Key: "{{ value }}", Value: strconv.Quote(*value)},
			)...,
		).
		CreateViolation()
}

// TimeComparisonConstraint is used to compare time values.
type TimeComparisonConstraint struct {
	isIgnored         bool
	groups            []string
	code              string
	messageTemplate   string
	messageParameters validation.TemplateParameterList
	comparedValue     time.Time
	layout            string
	isValid           func(value time.Time) bool
}

// IsEarlierThan checks that the given time is earlier than the specified value.
func IsEarlierThan(value time.Time) TimeComparisonConstraint {
	return TimeComparisonConstraint{
		code:            code.TooLate,
		messageTemplate: message.Templates[code.TooLate],
		comparedValue:   value,
		layout:          time.RFC3339,
		isValid: func(actualValue time.Time) bool {
			return actualValue.Before(value)
		},
	}
}

// IsEarlierThanOrEqual checks that the given time is earlier or equal to the specified value.
func IsEarlierThanOrEqual(value time.Time) TimeComparisonConstraint {
	return TimeComparisonConstraint{
		code:            code.TooLateOrEqual,
		messageTemplate: message.Templates[code.TooLateOrEqual],
		comparedValue:   value,
		layout:          time.RFC3339,
		isValid: func(actualValue time.Time) bool {
			return actualValue.Before(value) || actualValue.Equal(value)
		},
	}
}

// IsLaterThan checks that the given time is later than the specified value.
func IsLaterThan(value time.Time) TimeComparisonConstraint {
	return TimeComparisonConstraint{
		code:            code.TooEarly,
		messageTemplate: message.Templates[code.TooEarly],
		comparedValue:   value,
		layout:          time.RFC3339,
		isValid: func(actualValue time.Time) bool {
			return actualValue.After(value)
		},
	}
}

// IsLaterThanOrEqual checks that the given time is later or equal to the specified value.
func IsLaterThanOrEqual(value time.Time) TimeComparisonConstraint {
	return TimeComparisonConstraint{
		code:            code.TooEarlyOrEqual,
		messageTemplate: message.Templates[code.TooEarlyOrEqual],
		comparedValue:   value,
		layout:          time.RFC3339,
		isValid: func(actualValue time.Time) bool {
			return actualValue.After(value) || actualValue.Equal(value)
		},
	}
}

// SetUp always returns no error.
func (c TimeComparisonConstraint) SetUp() error {
	return nil
}

// Name is the constraint name.
func (c TimeComparisonConstraint) Name() string {
	return "TimeComparisonConstraint"
}

// Code overrides default code for produced violation.
func (c TimeComparisonConstraint) Code(code string) TimeComparisonConstraint {
	c.code = code
	return c
}

// Message sets the violation message template. You can set custom template parameters
// for injecting its values into the final message. Also, you can use default parameters:
//
//  {{ comparedValue }} - the expected value;
//  {{ value }} - the current (invalid) value.
//
// All values are formatted by the layout that can be defined by the Layout method.
// Default layout is time.RFC3339.
func (c TimeComparisonConstraint) Message(
	template string,
	parameters ...validation.TemplateParameter,
) TimeComparisonConstraint {
	c.messageTemplate = template
	c.messageParameters = parameters
	return c
}

// Layout can be used to set the layout that is used to format time values.
func (c TimeComparisonConstraint) Layout(layout string) TimeComparisonConstraint {
	c.layout = layout
	return c
}

// When enables conditional validation of this constraint. If the expression evaluates to false,
// then the constraint will be ignored.
func (c TimeComparisonConstraint) When(condition bool) TimeComparisonConstraint {
	c.isIgnored = !condition
	return c
}

// WhenGroups enables conditional validation of the constraint by using the validation groups.
func (c TimeComparisonConstraint) WhenGroups(groups ...string) TimeComparisonConstraint {
	c.groups = groups
	return c
}

func (c TimeComparisonConstraint) ValidateTime(value *time.Time, scope validation.Scope) error {
	if c.isIgnored || scope.IsIgnored(c.groups...) || value == nil || c.isValid(*value) {
		return nil
	}

	return scope.BuildViolation(c.code, c.messageTemplate).
		SetParameters(
			c.messageParameters.Prepend(
				validation.TemplateParameter{Key: "{{ comparedValue }}", Value: c.comparedValue.Format(c.layout)},
				validation.TemplateParameter{Key: "{{ value }}", Value: value.Format(c.layout)},
			)...,
		).
		CreateViolation()
}

// TimeRangeConstraint is used to check that a given time value is between some minimum and maximum.
type TimeRangeConstraint struct {
	isIgnored         bool
	groups            []string
	code              string
	messageTemplate   string
	messageParameters validation.TemplateParameterList
	layout            string
	min               time.Time
	max               time.Time
}

// IsBetweenTime checks that the time is between specified minimum and maximum time values.
func IsBetweenTime(min, max time.Time) TimeRangeConstraint {
	return TimeRangeConstraint{
		code:            code.NotInRange,
		messageTemplate: message.Templates[code.NotInRange],
		layout:          time.RFC3339,
		min:             min,
		max:             max,
	}
}

// SetUp returns an error if min is greater than or equal to max.
func (c TimeRangeConstraint) SetUp() error {
	if c.min.After(c.max) || c.min.Equal(c.max) {
		return errInvalidRange
	}

	return nil
}

// Name is the constraint name.
func (c TimeRangeConstraint) Name() string {
	return "TimeRangeConstraint"
}

// Code overrides default code for produced violation.
func (c TimeRangeConstraint) Code(code string) TimeRangeConstraint {
	c.code = code
	return c
}

// Message sets the violation message template. You can set custom template parameters
// for injecting its values into the final message. Also, you can use default parameters:
//
//  {{ max }} - the upper limit;
//  {{ min }} - the lower limit;
//  {{ value }} - the current (invalid) value.
//
// All values are formatted by the layout that can be defined by the Layout method.
// Default layout is time.RFC3339.
func (c TimeRangeConstraint) Message(template string, parameters ...validation.TemplateParameter) TimeRangeConstraint {
	c.messageTemplate = template
	c.messageParameters = parameters
	return c
}

// When enables conditional validation of this constraint. If the expression evaluates to false,
// then the constraint will be ignored.
func (c TimeRangeConstraint) When(condition bool) TimeRangeConstraint {
	c.isIgnored = !condition
	return c
}

// WhenGroups enables conditional validation of the constraint by using the validation groups.
func (c TimeRangeConstraint) WhenGroups(groups ...string) TimeRangeConstraint {
	c.groups = groups
	return c
}

// Layout can be used to set the layout that is used to format time values.
func (c TimeRangeConstraint) Layout(layout string) TimeRangeConstraint {
	c.layout = layout
	return c
}

func (c TimeRangeConstraint) ValidateTime(value *time.Time, scope validation.Scope) error {
	if c.isIgnored || scope.IsIgnored(c.groups...) || value == nil {
		return nil
	}
	if value.Before(c.min) || value.After(c.max) {
		return c.newViolation(value, scope)
	}

	return nil
}

func (c TimeRangeConstraint) newViolation(value *time.Time, scope validation.Scope) validation.Violation {
	return scope.BuildViolation(c.code, c.messageTemplate).
		SetParameters(
			c.messageParameters.Prepend(
				validation.TemplateParameter{Key: "{{ min }}", Value: c.min.Format(c.layout)},
				validation.TemplateParameter{Key: "{{ max }}", Value: c.max.Format(c.layout)},
				validation.TemplateParameter{Key: "{{ value }}", Value: value.Format(c.layout)},
			)...,
		).
		CreateViolation()
}

// UniqueConstraint is used to check that all elements of the given collection are unique.
type UniqueConstraint struct {
	isIgnored         bool
	groups            []string
	code              string
	messageTemplate   string
	messageParameters validation.TemplateParameterList
}

// HasUniqueValues checks that all elements of the given collection are unique
// (none of them is present more than once).
func HasUniqueValues() UniqueConstraint {
	return UniqueConstraint{
		code:            code.NotUnique,
		messageTemplate: message.Templates[code.NotUnique],
	}
}

// SetUp always returns no error.
func (c UniqueConstraint) SetUp() error {
	return nil
}

// Name is the constraint name.
func (c UniqueConstraint) Name() string {
	return "UniqueConstraint"
}

// Code overrides default code for produced violation.
func (c UniqueConstraint) Code(code string) UniqueConstraint {
	c.code = code
	return c
}

// Message sets the violation message template. You can set custom template parameters
// for injecting its values into the final message.
func (c UniqueConstraint) Message(template string, parameters ...validation.TemplateParameter) UniqueConstraint {
	c.messageTemplate = template
	c.messageParameters = parameters
	return c
}

// When enables conditional validation of this constraint. If the expression evaluates to false,
// then the constraint will be ignored.
func (c UniqueConstraint) When(condition bool) UniqueConstraint {
	c.isIgnored = !condition
	return c
}

// WhenGroups enables conditional validation of the constraint by using the validation groups.
func (c UniqueConstraint) WhenGroups(groups ...string) UniqueConstraint {
	c.groups = groups
	return c
}

func (c UniqueConstraint) ValidateStrings(values []string, scope validation.Scope) error {
	if c.isIgnored || scope.IsIgnored(c.groups...) || is.UniqueStrings(values) {
		return nil
	}

	return scope.BuildViolation(c.code, c.messageTemplate).
		SetParameters(c.messageParameters...).
		CreateViolation()
}
