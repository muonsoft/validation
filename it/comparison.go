package it

import (
	"strconv"
	"time"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/code"
	"github.com/muonsoft/validation/generic"
	"github.com/muonsoft/validation/message"
)

// NumberComparisonConstraint is used for various numeric comparisons between integer and float values.
// Values are compared as integers if the compared and specified values are integers.
// Otherwise, numbers are always compared as floating point numbers.
type NumberComparisonConstraint struct {
	isIgnored       bool
	code            string
	messageTemplate string
	comparedValue   string
	isValid         func(value generic.Number) bool
}

// IsEqualToInteger checks that the number (integer or float) is equal to the specified integer value.
// Values are compared as integers if the compared and specified values are integers.
// Otherwise, numbers are always compared as floating point numbers.
//
// Example
//  v := 1
//  err := validator.ValidateNumber(&v, it.IsEqualToInteger(2))
func IsEqualToInteger(value int64) NumberComparisonConstraint {
	v := generic.NewNumberFromInt(value)

	return NumberComparisonConstraint{
		code:            code.Equal,
		messageTemplate: message.Equal,
		comparedValue:   v.String(),
		isValid: func(n generic.Number) bool {
			return n.IsEqualTo(v)
		},
	}
}

// IsEqualToFloat checks that the number (integer or float) is equal to the specified float value.
// Values are compared as integers if the compared and specified values are integers.
// Otherwise, numbers are always compared as floating point numbers.
//
// Example
//  v := 1.1
//  err := validator.ValidateNumber(&v, it.IsEqualToFloat(1.2))
func IsEqualToFloat(value float64) NumberComparisonConstraint {
	v := generic.NewNumberFromFloat(value)

	return NumberComparisonConstraint{
		code:            code.Equal,
		messageTemplate: message.Equal,
		comparedValue:   v.String(),
		isValid: func(n generic.Number) bool {
			return n.IsEqualTo(v)
		},
	}
}

// IsNotEqualToInteger checks that the number (integer or float) is not equal to the specified integer value.
// Values are compared as integers if the compared and specified values are integers.
// Otherwise, numbers are always compared as floating point numbers.
//
// Example
//  v := 1
//  err := validator.ValidateNumber(&v, it.IsNotEqualToInteger(1))
func IsNotEqualToInteger(value int64) NumberComparisonConstraint {
	v := generic.NewNumberFromInt(value)

	return NumberComparisonConstraint{
		code:            code.NotEqual,
		messageTemplate: message.NotEqual,
		comparedValue:   v.String(),
		isValid: func(n generic.Number) bool {
			return !n.IsEqualTo(v)
		},
	}
}

// IsNotEqualToFloat checks that the number (integer or float) is not equal to the specified float value.
// Values are compared as integers if the compared and specified values are integers.
// Otherwise, numbers are always compared as floating point numbers.
//
// Example
//  v := 1.1
//  err := validator.ValidateNumber(&v, it.IsNotEqualToFloat(1.1))
func IsNotEqualToFloat(value float64) NumberComparisonConstraint {
	v := generic.NewNumberFromFloat(value)

	return NumberComparisonConstraint{
		code:            code.NotEqual,
		messageTemplate: message.NotEqual,
		comparedValue:   v.String(),
		isValid: func(n generic.Number) bool {
			return !n.IsEqualTo(v)
		},
	}
}

// IsLessThanInteger checks that the number (integer or float) is less than the specified integer value.
// Values are compared as integers if the compared and specified values are integers.
// Otherwise, numbers are always compared as floating point numbers.
//
// Example
//  v := 1
//  err := validator.ValidateNumber(&v, it.IsLessThanInteger(1))
func IsLessThanInteger(value int64) NumberComparisonConstraint {
	v := generic.NewNumberFromInt(value)

	return NumberComparisonConstraint{
		code:            code.TooHigh,
		messageTemplate: message.TooHigh,
		comparedValue:   v.String(),
		isValid: func(n generic.Number) bool {
			return n.IsLessThan(v)
		},
	}
}

// IsLessThanFloat checks that the number (integer or float) is less than the specified float value.
// Values are compared as integers if the compared and specified values are integers.
// Otherwise, numbers are always compared as floating point numbers.
//
// Example
//  v := 1.1
//  err := validator.ValidateNumber(&v, it.IsLessThanFloat(1.1))
func IsLessThanFloat(value float64) NumberComparisonConstraint {
	v := generic.NewNumberFromFloat(value)

	return NumberComparisonConstraint{
		code:            code.TooHigh,
		messageTemplate: message.TooHigh,
		comparedValue:   v.String(),
		isValid: func(n generic.Number) bool {
			return n.IsLessThan(v)
		},
	}
}

// IsLessThanOrEqualInteger checks that the number (integer or float) is less than or
// equal to the specified integer value. Values are compared as integers if the compared
// and specified values are integers. Otherwise, numbers are always compared as floating point numbers.
//
// Example
//  v := 1
//  err := validator.ValidateNumber(&v, it.IsLessThanOrEqualInteger(2))
func IsLessThanOrEqualInteger(value int64) NumberComparisonConstraint {
	v := generic.NewNumberFromInt(value)

	return NumberComparisonConstraint{
		code:            code.TooHighOrEqual,
		messageTemplate: message.TooHighOrEqual,
		comparedValue:   v.String(),
		isValid: func(n generic.Number) bool {
			return n.IsLessThan(v) || n.IsEqualTo(v)
		},
	}
}

// IsLessThanOrEqualFloat checks that the number (integer or float) is less than or
// equal to the specified float value. Values are compared as integers if the compared
// and specified values are integers. Otherwise, numbers are always compared as floating point numbers.
//
// Example
//  v := 1.1
//  err := validator.ValidateNumber(&v, it.IsLessThanOrEqualFloat(1.2))
func IsLessThanOrEqualFloat(value float64) NumberComparisonConstraint {
	v := generic.NewNumberFromFloat(value)

	return NumberComparisonConstraint{
		code:            code.TooHighOrEqual,
		messageTemplate: message.TooHighOrEqual,
		comparedValue:   v.String(),
		isValid: func(n generic.Number) bool {
			return n.IsLessThan(v) || n.IsEqualTo(v)
		},
	}
}

// IsGreaterThanInteger checks that the number (integer or float) is greater than the specified integer value.
// Values are compared as integers if the compared and specified values are integers.
// Otherwise, numbers are always compared as floating point numbers.
//
// Example
//  v := 1
//  err := validator.ValidateNumber(&v, it.IsGreaterThanInteger(1))
func IsGreaterThanInteger(value int64) NumberComparisonConstraint {
	v := generic.NewNumberFromInt(value)

	return NumberComparisonConstraint{
		code:            code.TooLow,
		messageTemplate: message.TooLow,
		comparedValue:   v.String(),
		isValid: func(n generic.Number) bool {
			return n.IsGreaterThan(v)
		},
	}
}

// IsGreaterThanFloat checks that the number (integer or float) is greater than the specified float value.
// Values are compared as integers if the compared and specified values are integers.
// Otherwise, numbers are always compared as floating point numbers.
//
// Example
//  v := 1.1
//  err := validator.ValidateNumber(&v, it.IsGreaterThanFloat(1.1))
func IsGreaterThanFloat(value float64) NumberComparisonConstraint {
	v := generic.NewNumberFromFloat(value)

	return NumberComparisonConstraint{
		code:            code.TooLow,
		messageTemplate: message.TooLow,
		comparedValue:   v.String(),
		isValid: func(n generic.Number) bool {
			return n.IsGreaterThan(v)
		},
	}
}

// IsGreaterThanOrEqualInteger checks that the number (integer or float) is greater than or
// equal to the specified integer value. Values are compared as integers if the compared
// and specified values are integers. Otherwise, numbers are always compared as floating point numbers.
//
// Example
//  v := 1
//  err := validator.ValidateNumber(&v, it.IsGreaterThanOrEqualInteger(2))
func IsGreaterThanOrEqualInteger(value int64) NumberComparisonConstraint {
	v := generic.NewNumberFromInt(value)

	return NumberComparisonConstraint{
		code:            code.TooLowOrEqual,
		messageTemplate: message.TooLowOrEqual,
		comparedValue:   v.String(),
		isValid: func(n generic.Number) bool {
			return n.IsGreaterThan(v) || n.IsEqualTo(v)
		},
	}
}

// IsGreaterThanOrEqualFloat checks that the number (integer or float) is greater than or
// equal to the specified float value. Values are compared as integers if the compared
// and specified values are integers. Otherwise, numbers are always compared as floating point numbers.
//
// Example
//  v := 1.1
//  err := validator.ValidateNumber(&v, it.IsGreaterThanOrEqualFloat(1.2))
func IsGreaterThanOrEqualFloat(value float64) NumberComparisonConstraint {
	v := generic.NewNumberFromFloat(value)

	return NumberComparisonConstraint{
		code:            code.TooLowOrEqual,
		messageTemplate: message.TooLowOrEqual,
		comparedValue:   v.String(),
		isValid: func(n generic.Number) bool {
			return n.IsGreaterThan(v) || n.IsEqualTo(v)
		},
	}
}

// IsPositive checks that the value is a positive number (integer or float). Zero is neither
// positive nor negative. If you want to allow zero use IsPositiveOrZero comparison.
//
// Example
//  v := -1
//  err := validator.ValidateNumber(&v, it.IsPositive())
func IsPositive() NumberComparisonConstraint {
	v := generic.NewNumberFromInt(0)

	return NumberComparisonConstraint{
		code:            code.NotPositive,
		messageTemplate: message.NotPositive,
		comparedValue:   v.String(),
		isValid: func(n generic.Number) bool {
			return n.IsGreaterThan(v)
		},
	}
}

// IsPositiveOrZero checks that the value is a positive number (integer or float) or equal to zero.
// If you don't want to allow zero as a valid value, use IsPositive comparison.
//
// Example
//  v := -1
//  err := validator.ValidateNumber(&v, it.IsPositiveOrZero())
func IsPositiveOrZero() NumberComparisonConstraint {
	v := generic.NewNumberFromInt(0)

	return NumberComparisonConstraint{
		code:            code.NotPositiveOrZero,
		messageTemplate: message.NotPositiveOrZero,
		comparedValue:   v.String(),
		isValid: func(n generic.Number) bool {
			return n.IsGreaterThan(v) || n.IsEqualTo(v)
		},
	}
}

// IsNegative checks that the value is a negative number (integer or float). Zero is neither
// positive nor negative. If you want to allow zero use IsNegativeOrZero comparison.
//
// Example
//  v := 1
//  err := validator.ValidateNumber(&v, it.IsNegative())
func IsNegative() NumberComparisonConstraint {
	v := generic.NewNumberFromInt(0)

	return NumberComparisonConstraint{
		code:            code.NotNegative,
		messageTemplate: message.NotNegative,
		comparedValue:   v.String(),
		isValid: func(n generic.Number) bool {
			return n.IsLessThan(v)
		},
	}
}

// IsNegativeOrZero checks that the value is a negative number (integer or float) or equal to zero.
// If you don't want to allow zero as a valid value, use IsNegative comparison.
//
// Example
//  v := -1
//  err := validator.ValidateNumber(&v, it.IsNegativeOrZero())
func IsNegativeOrZero() NumberComparisonConstraint {
	v := generic.NewNumberFromInt(0)

	return NumberComparisonConstraint{
		code:            code.NotNegativeOrZero,
		messageTemplate: message.NotNegativeOrZero,
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

// Message sets the violation message template. You can use template parameters
// for injecting its values into the final message:
//
//  {{ comparedValue }} - the expected value;
//  {{ value }} - the current (invalid) value.
func (c NumberComparisonConstraint) Message(message string) NumberComparisonConstraint {
	c.messageTemplate = message
	return c
}

// When enables conditional validation of this constraint. If the expression evaluates to false,
// then the constraint will be ignored.
func (c NumberComparisonConstraint) When(condition bool) NumberComparisonConstraint {
	c.isIgnored = !condition
	return c
}

func (c NumberComparisonConstraint) ValidateNumber(value generic.Number, scope validation.Scope) error {
	if c.isIgnored || value.IsNil() || c.isValid(value) {
		return nil
	}

	return scope.BuildViolation(c.code, c.messageTemplate).
		SetParameters([]validation.TemplateParameter{
			{Key: "{{ comparedValue }}", Value: c.comparedValue},
			{Key: "{{ value }}", Value: value.String()},
		}).
		CreateViolation()
}

// StringComparisonConstraint is used to compare strings.
type StringComparisonConstraint struct {
	isIgnored       bool
	code            string
	messageTemplate string
	comparedValue   string
	isValid         func(value string) bool
}

// IsEqualToString checks that the string value is equal to the specified string value.
//
// Example
//  v := "actual"
//  err := validator.ValidateString(&v, it.IsEqualToString("expected"))
func IsEqualToString(value string) StringComparisonConstraint {
	return StringComparisonConstraint{
		code:            code.Equal,
		messageTemplate: message.Equal,
		comparedValue:   value,
		isValid: func(actualValue string) bool {
			return value == actualValue
		},
	}
}

// IsNotEqualToString checks that the string value is not equal to the specified string value.
//
// Example
//  v := "expected"
//  err := validator.ValidateString(&v, it.IsNotEqualToString("expected"))
func IsNotEqualToString(value string) StringComparisonConstraint {
	return StringComparisonConstraint{
		code:            code.NotEqual,
		messageTemplate: message.NotEqual,
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

// Message sets the violation message template. You can use template parameters
// for injecting its values into the final message:
//
//  {{ comparedValue }} - the expected value;
//  {{ value }} - the current (invalid) value.
//
// All string values are quoted strings.
func (c StringComparisonConstraint) Message(message string) StringComparisonConstraint {
	c.messageTemplate = message
	return c
}

// When enables conditional validation of this constraint. If the expression evaluates to false,
// then the constraint will be ignored.
func (c StringComparisonConstraint) When(condition bool) StringComparisonConstraint {
	c.isIgnored = !condition
	return c
}

func (c StringComparisonConstraint) ValidateString(value *string, scope validation.Scope) error {
	if c.isIgnored || value == nil || c.isValid(*value) {
		return nil
	}

	return scope.BuildViolation(c.code, c.messageTemplate).
		SetParameters([]validation.TemplateParameter{
			{Key: "{{ comparedValue }}", Value: strconv.Quote(c.comparedValue)},
			{Key: "{{ value }}", Value: strconv.Quote(*value)},
		}).
		CreateViolation()
}

// TimeComparisonConstraint is used to compare time values.
type TimeComparisonConstraint struct {
	isIgnored       bool
	code            string
	messageTemplate string
	comparedValue   time.Time
	layout          string
	isValid         func(value time.Time) bool
}

// IsEarlierThan checks that the given time is earlier than the specified value.
//
// Example
//  t := time.Now()
//  err := validator.ValidateTime(&t, it.IsEarlierThan(time.Now().Add(time.Hour)))
func IsEarlierThan(value time.Time) TimeComparisonConstraint {
	return TimeComparisonConstraint{
		code:            code.TooLate,
		messageTemplate: message.TooLate,
		comparedValue:   value,
		layout:          time.RFC3339,
		isValid: func(actualValue time.Time) bool {
			return actualValue.Before(value)
		},
	}
}

// IsEarlierThanOrEqual checks that the given time is earlier or equal to the specified value.
//
// Example
//  t := time.Now()
//  err := validator.ValidateTime(&t, it.IsEarlierThanOrEqual(time.Now().Add(time.Hour)))
func IsEarlierThanOrEqual(value time.Time) TimeComparisonConstraint {
	return TimeComparisonConstraint{
		code:            code.TooLateOrEqual,
		messageTemplate: message.TooLateOrEqual,
		comparedValue:   value,
		layout:          time.RFC3339,
		isValid: func(actualValue time.Time) bool {
			return actualValue.Before(value) || actualValue.Equal(value)
		},
	}
}

// IsLaterThan checks that the given time is later than the specified value.
//
// Example
//  t := time.Now()
//  err := validator.ValidateTime(&t, it.IsLaterThan(time.Now().Sub(time.Hour)))
func IsLaterThan(value time.Time) TimeComparisonConstraint {
	return TimeComparisonConstraint{
		code:            code.TooEarly,
		messageTemplate: message.TooEarly,
		comparedValue:   value,
		layout:          time.RFC3339,
		isValid: func(actualValue time.Time) bool {
			return actualValue.After(value)
		},
	}
}

// IsLaterThanOrEqual checks that the given time is later or equal to the specified value.
//
// Example
//  t := time.Now()
//  err := validator.ValidateTime(&t, it.IsLaterThanOrEqual(time.Now().Sub(time.Hour)))
func IsLaterThanOrEqual(value time.Time) TimeComparisonConstraint {
	return TimeComparisonConstraint{
		code:            code.TooEarlyOrEqual,
		messageTemplate: message.TooEarlyOrEqual,
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

// Message sets the violation message template. You can use template parameters
// for injecting its values into the final message:
//
//  {{ comparedValue }} - the expected value;
//  {{ value }} - the current (invalid) value.
//
// All values are formatted by the layout that can be defined by the Layout method.
// Default layout is time.RFC3339.
func (c TimeComparisonConstraint) Message(message string) TimeComparisonConstraint {
	c.messageTemplate = message
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

func (c TimeComparisonConstraint) ValidateTime(value *time.Time, scope validation.Scope) error {
	if c.isIgnored || value == nil || c.isValid(*value) {
		return nil
	}

	return scope.BuildViolation(c.code, c.messageTemplate).
		SetParameters([]validation.TemplateParameter{
			{Key: "{{ comparedValue }}", Value: c.comparedValue.Format(c.layout)},
			{Key: "{{ value }}", Value: value.Format(c.layout)},
		}).
		CreateViolation()
}
