package it

import (
	"fmt"
	"reflect"
	"time"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/is"
)

// ComparisonConstraint is used for comparisons between comparable generic types.
type ComparisonConstraint[T comparable] struct {
	isIgnored         bool
	value             T
	groups            []string
	err               error
	messageTemplate   string
	messageParameters validation.TemplateParameterList
	comparedValue     string
	isValid           func(value T) bool
}

// IsEqualTo checks that the value is equal to the specified value.
func IsEqualTo[T comparable](value T) ComparisonConstraint[T] {
	return ComparisonConstraint[T]{
		err:             validation.ErrNotEqual,
		value:           value,
		messageTemplate: validation.ErrNotEqual.Template(),
		comparedValue:   formatComparable(value),
		isValid:         func(v T) bool { return v == value },
	}
}

// IsNotEqualTo checks that the value is not equal to the specified value.
func IsNotEqualTo[T comparable](value T) ComparisonConstraint[T] {
	return ComparisonConstraint[T]{
		err:             validation.ErrIsEqual,
		value:           value,
		messageTemplate: validation.ErrIsEqual.Template(),
		comparedValue:   formatComparable(value),
		isValid:         func(v T) bool { return v != value },
	}
}

// WithError overrides default error for produced violation.
func (c ComparisonConstraint[T]) WithError(err error) ComparisonConstraint[T] {
	c.err = err
	return c
}

// WithMessage sets the violation message template. You can set custom template parameters
// for injecting its values into the final message. Also, you can use default parameters:
//
//  {{ comparedValue }} - the expected value;
//  {{ value }} - the current (invalid) value.
func (c ComparisonConstraint[T]) WithMessage(
	template string,
	parameters ...validation.TemplateParameter,
) ComparisonConstraint[T] {
	c.messageTemplate = template
	c.messageParameters = parameters
	return c
}

// When enables conditional validation of this constraint. If the expression evaluates to false,
// then the constraint will be ignored.
func (c ComparisonConstraint[T]) When(condition bool) ComparisonConstraint[T] {
	c.isIgnored = !condition
	return c
}

// WhenGroups enables conditional validation of the constraint by using the validation groups.
func (c ComparisonConstraint[T]) WhenGroups(groups ...string) ComparisonConstraint[T] {
	c.groups = groups
	return c
}

func (c ComparisonConstraint[T]) ValidateNumber(value *T, scope validation.Scope) error {
	return c.ValidateComparable(value, scope)
}

func (c ComparisonConstraint[T]) ValidateString(value *T, scope validation.Scope) error {
	return c.ValidateComparable(value, scope)
}

func (c ComparisonConstraint[T]) ValidateComparable(value *T, scope validation.Scope) error {
	if c.isIgnored || scope.IsIgnored(c.groups...) || value == nil || c.isValid(*value) {
		return nil
	}

	return scope.BuildViolation(c.err, c.messageTemplate).
		WithParameters(
			c.messageParameters.Prepend(
				validation.TemplateParameter{Key: "{{ comparedValue }}", Value: c.comparedValue},
				validation.TemplateParameter{Key: "{{ value }}", Value: formatComparable(*value)},
			)...,
		).
		Create()
}

// NumberComparisonConstraint is used for various numeric comparisons between integer and float values.
type NumberComparisonConstraint[T validation.Numeric] struct {
	isIgnored         bool
	value             T
	groups            []string
	err               error
	messageTemplate   string
	messageParameters validation.TemplateParameterList
	comparedValue     string
	isValid           func(value T) bool
}

// IsLessThan checks that the number is less than the specified value.
func IsLessThan[T validation.Numeric](value T) NumberComparisonConstraint[T] {
	return NumberComparisonConstraint[T]{
		err:             validation.ErrTooHigh,
		value:           value,
		messageTemplate: validation.ErrTooHigh.Template(),
		comparedValue:   fmt.Sprint(value),
		isValid:         func(n T) bool { return n < value },
	}
}

// IsLessThanOrEqual checks that the number is less than or equal to the specified value.
func IsLessThanOrEqual[T validation.Numeric](value T) NumberComparisonConstraint[T] {
	return NumberComparisonConstraint[T]{
		err:             validation.ErrTooHighOrEqual,
		value:           value,
		messageTemplate: validation.ErrTooHighOrEqual.Template(),
		comparedValue:   fmt.Sprint(value),
		isValid:         func(n T) bool { return n <= value },
	}
}

// IsGreaterThan checks that the number is greater than the specified value.
func IsGreaterThan[T validation.Numeric](value T) NumberComparisonConstraint[T] {
	return NumberComparisonConstraint[T]{
		err:             validation.ErrTooLow,
		value:           value,
		messageTemplate: validation.ErrTooLow.Template(),
		comparedValue:   fmt.Sprint(value),
		isValid:         func(n T) bool { return n > value },
	}
}

// IsGreaterThanOrEqual checks that the number is greater than or equal to the specified value.
func IsGreaterThanOrEqual[T validation.Numeric](value T) NumberComparisonConstraint[T] {
	return NumberComparisonConstraint[T]{
		err:             validation.ErrTooLowOrEqual,
		value:           value,
		messageTemplate: validation.ErrTooLowOrEqual.Template(),
		comparedValue:   fmt.Sprint(value),
		isValid:         func(n T) bool { return n >= value },
	}
}

// IsPositive checks that the value is a positive number. Zero is neither positive nor negative.
// If you want to allow zero use IsPositiveOrZero comparison.
func IsPositive[T validation.Numeric]() NumberComparisonConstraint[T] {
	return NumberComparisonConstraint[T]{
		err:             validation.ErrNotPositive,
		value:           0,
		messageTemplate: validation.ErrNotPositive.Template(),
		comparedValue:   "0",
		isValid:         func(n T) bool { return n > 0 },
	}
}

// IsPositiveOrZero checks that the value is a positive number or equal to zero.
// If you don't want to allow zero as a valid value, use IsPositive comparison.
func IsPositiveOrZero[T validation.Numeric]() NumberComparisonConstraint[T] {
	return NumberComparisonConstraint[T]{
		err:             validation.ErrNotPositiveOrZero,
		value:           0,
		messageTemplate: validation.ErrNotPositiveOrZero.Template(),
		comparedValue:   "0",
		isValid:         func(n T) bool { return n >= 0 },
	}
}

// IsNegative checks that the value is a negative number. Zero is neither positive nor negative.
// If you want to allow zero use IsNegativeOrZero comparison.
func IsNegative[T validation.Numeric]() NumberComparisonConstraint[T] {
	return NumberComparisonConstraint[T]{
		err:             validation.ErrNotNegative,
		value:           0,
		messageTemplate: validation.ErrNotNegative.Template(),
		comparedValue:   "0",
		isValid:         func(n T) bool { return n < 0 },
	}
}

// IsNegativeOrZero checks that the value is a negative number or equal to zero.
// If you don't want to allow zero as a valid value, use IsNegative comparison.
func IsNegativeOrZero[T validation.Numeric]() NumberComparisonConstraint[T] {
	return NumberComparisonConstraint[T]{
		err:             validation.ErrNotNegativeOrZero,
		value:           0,
		messageTemplate: validation.ErrNotNegativeOrZero.Template(),
		comparedValue:   "0",
		isValid:         func(n T) bool { return n <= 0 },
	}
}

// WithError overrides default error for produced violation.
func (c NumberComparisonConstraint[T]) WithError(err error) NumberComparisonConstraint[T] {
	c.err = err
	return c
}

// WithMessage sets the violation message template. You can set custom template parameters
// for injecting its values into the final message. Also, you can use default parameters:
//
//  {{ comparedValue }} - the expected value;
//  {{ value }} - the current (invalid) value.
func (c NumberComparisonConstraint[T]) WithMessage(
	template string,
	parameters ...validation.TemplateParameter,
) NumberComparisonConstraint[T] {
	c.messageTemplate = template
	c.messageParameters = parameters
	return c
}

// When enables conditional validation of this constraint. If the expression evaluates to false,
// then the constraint will be ignored.
func (c NumberComparisonConstraint[T]) When(condition bool) NumberComparisonConstraint[T] {
	c.isIgnored = !condition
	return c
}

// WhenGroups enables conditional validation of the constraint by using the validation groups.
func (c NumberComparisonConstraint[T]) WhenGroups(groups ...string) NumberComparisonConstraint[T] {
	c.groups = groups
	return c
}

func (c NumberComparisonConstraint[T]) ValidateNumber(value *T, scope validation.Scope) error {
	if c.isIgnored || scope.IsIgnored(c.groups...) || value == nil || c.isValid(*value) {
		return nil
	}

	return scope.BuildViolation(c.err, c.messageTemplate).
		WithParameters(
			c.messageParameters.Prepend(
				validation.TemplateParameter{Key: "{{ comparedValue }}", Value: c.comparedValue},
				validation.TemplateParameter{Key: "{{ value }}", Value: fmt.Sprint(*value)},
			)...,
		).
		Create()
}

// RangeConstraint is used to check that a given number value is between some minimum and maximum.
type RangeConstraint[T validation.Numeric] struct {
	isIgnored         bool
	groups            []string
	err               error
	messageTemplate   string
	messageParameters validation.TemplateParameterList
	min               T
	max               T
}

// IsBetween checks that the number is between specified minimum and maximum numeric values.
func IsBetween[T validation.Numeric](min, max T) RangeConstraint[T] {
	return RangeConstraint[T]{
		min:             min,
		max:             max,
		err:             validation.ErrNotInRange,
		messageTemplate: validation.ErrNotInRange.Template(),
	}
}

// Name is the constraint name.
func (c RangeConstraint[T]) Name() string {
	return fmt.Sprintf("RangeConstraint[%s]", reflect.TypeOf(c.min).String())
}

// WithError overrides default error for produced violation.
func (c RangeConstraint[T]) WithError(err error) RangeConstraint[T] {
	c.err = err
	return c
}

// WithMessage sets the violation message template. You can set custom template parameters
// for injecting its values into the final message. Also, you can use default parameters:
//
//  {{ max }} - the upper limit;
//  {{ min }} - the lower limit;
//  {{ value }} - the current (invalid) value.
func (c RangeConstraint[T]) WithMessage(template string, parameters ...validation.TemplateParameter) RangeConstraint[T] {
	c.messageTemplate = template
	c.messageParameters = parameters
	return c
}

// When enables conditional validation of this constraint. If the expression evaluates to false,
// then the constraint will be ignored.
func (c RangeConstraint[T]) When(condition bool) RangeConstraint[T] {
	c.isIgnored = !condition
	return c
}

// WhenGroups enables conditional validation of the constraint by using the validation groups.
func (c RangeConstraint[T]) WhenGroups(groups ...string) RangeConstraint[T] {
	c.groups = groups
	return c
}

func (c RangeConstraint[T]) ValidateNumber(value *T, scope validation.Scope) error {
	if c.min >= c.max {
		return scope.NewConstraintError(c.Name(), "invalid range")
	}
	if c.isIgnored || value == nil || scope.IsIgnored(c.groups...) {
		return nil
	}
	if *value < c.min || *value > c.max {
		return c.newViolation(*value, scope)
	}

	return nil
}

func (c RangeConstraint[T]) newViolation(value T, scope validation.Scope) error {
	return scope.BuildViolation(c.err, c.messageTemplate).
		WithParameters(
			c.messageParameters.Prepend(
				validation.TemplateParameter{Key: "{{ min }}", Value: fmt.Sprint(c.min)},
				validation.TemplateParameter{Key: "{{ max }}", Value: fmt.Sprint(c.max)},
				validation.TemplateParameter{Key: "{{ value }}", Value: fmt.Sprint(value)},
			)...,
		).
		Create()
}

// TimeComparisonConstraint is used to compare time values.
type TimeComparisonConstraint struct {
	isIgnored         bool
	groups            []string
	err               error
	messageTemplate   string
	messageParameters validation.TemplateParameterList
	comparedValue     time.Time
	layout            string
	isValid           func(value time.Time) bool
}

// IsEarlierThan checks that the given time is earlier than the specified value.
func IsEarlierThan(value time.Time) TimeComparisonConstraint {
	return TimeComparisonConstraint{
		err:             validation.ErrTooLate,
		messageTemplate: validation.ErrTooLate.Template(),
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
		err:             validation.ErrTooLateOrEqual,
		messageTemplate: validation.ErrTooLateOrEqual.Template(),
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
		err:             validation.ErrTooEarly,
		messageTemplate: validation.ErrTooEarly.Template(),
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
		err:             validation.ErrTooEarlyOrEqual,
		messageTemplate: validation.ErrTooEarlyOrEqual.Template(),
		comparedValue:   value,
		layout:          time.RFC3339,
		isValid: func(actualValue time.Time) bool {
			return actualValue.After(value) || actualValue.Equal(value)
		},
	}
}

// WithError overrides default error for produced violation.
func (c TimeComparisonConstraint) WithError(err error) TimeComparisonConstraint {
	c.err = err
	return c
}

// WithMessage sets the violation message template. You can set custom template parameters
// for injecting its values into the final message. Also, you can use default parameters:
//
//  {{ comparedValue }} - the expected value;
//  {{ value }} - the current (invalid) value.
//
// All values are formatted by the layout that can be defined by the WithLayout method.
// Default layout is time.RFC3339.
func (c TimeComparisonConstraint) WithMessage(
	template string,
	parameters ...validation.TemplateParameter,
) TimeComparisonConstraint {
	c.messageTemplate = template
	c.messageParameters = parameters
	return c
}

// WithLayout can be used to set the layout that is used to format time values.
func (c TimeComparisonConstraint) WithLayout(layout string) TimeComparisonConstraint {
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

	return scope.BuildViolation(c.err, c.messageTemplate).
		WithParameters(
			c.messageParameters.Prepend(
				validation.TemplateParameter{Key: "{{ comparedValue }}", Value: c.comparedValue.Format(c.layout)},
				validation.TemplateParameter{Key: "{{ value }}", Value: value.Format(c.layout)},
			)...,
		).
		Create()
}

// TimeRangeConstraint is used to check that a given time value is between some minimum and maximum.
type TimeRangeConstraint struct {
	isIgnored         bool
	groups            []string
	err               error
	messageTemplate   string
	messageParameters validation.TemplateParameterList
	layout            string
	min               time.Time
	max               time.Time
}

// IsBetweenTime checks that the time is between specified minimum and maximum time values.
func IsBetweenTime(min, max time.Time) TimeRangeConstraint {
	return TimeRangeConstraint{
		err:             validation.ErrNotInRange,
		messageTemplate: validation.ErrNotInRange.Template(),
		layout:          time.RFC3339,
		min:             min,
		max:             max,
	}
}

// WithError overrides default error for produced violation.
func (c TimeRangeConstraint) WithError(err error) TimeRangeConstraint {
	c.err = err
	return c
}

// WithMessage sets the violation message template. You can set custom template parameters
// for injecting its values into the final message. Also, you can use default parameters:
//
//  {{ max }} - the upper limit;
//  {{ min }} - the lower limit;
//  {{ value }} - the current (invalid) value.
//
// All values are formatted by the layout that can be defined by the WithLayout method.
// Default layout is time.RFC3339.
func (c TimeRangeConstraint) WithMessage(template string, parameters ...validation.TemplateParameter) TimeRangeConstraint {
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

// WithLayout can be used to set the layout that is used to format time values.
func (c TimeRangeConstraint) WithLayout(layout string) TimeRangeConstraint {
	c.layout = layout
	return c
}

func (c TimeRangeConstraint) ValidateTime(value *time.Time, scope validation.Scope) error {
	if c.min.After(c.max) || c.min.Equal(c.max) {
		return scope.NewConstraintError("TimeRangeConstraint", "invalid range")
	}
	if c.isIgnored || scope.IsIgnored(c.groups...) || value == nil {
		return nil
	}
	if value.Before(c.min) || value.After(c.max) {
		return c.newViolation(value, scope)
	}

	return nil
}

func (c TimeRangeConstraint) newViolation(value *time.Time, scope validation.Scope) validation.Violation {
	return scope.BuildViolation(c.err, c.messageTemplate).
		WithParameters(
			c.messageParameters.Prepend(
				validation.TemplateParameter{Key: "{{ min }}", Value: c.min.Format(c.layout)},
				validation.TemplateParameter{Key: "{{ max }}", Value: c.max.Format(c.layout)},
				validation.TemplateParameter{Key: "{{ value }}", Value: value.Format(c.layout)},
			)...,
		).
		Create()
}

// UniqueConstraint is used to check that all elements of the given collection are unique.
type UniqueConstraint[T comparable] struct {
	isIgnored         bool
	groups            []string
	err               error
	messageTemplate   string
	messageParameters validation.TemplateParameterList
}

// HasUniqueValues checks that all elements of the given collection are unique
// (none of them is present more than once).
func HasUniqueValues[T comparable]() UniqueConstraint[T] {
	return UniqueConstraint[T]{
		err:             validation.ErrNotUnique,
		messageTemplate: validation.ErrNotUnique.Template(),
	}
}

// WithError overrides default error for produced violation.
func (c UniqueConstraint[T]) WithError(err error) UniqueConstraint[T] {
	c.err = err
	return c
}

// WithMessage sets the violation message template. You can set custom template parameters
// for injecting its values into the final message.
func (c UniqueConstraint[T]) WithMessage(template string, parameters ...validation.TemplateParameter) UniqueConstraint[T] {
	c.messageTemplate = template
	c.messageParameters = parameters
	return c
}

// When enables conditional validation of this constraint. If the expression evaluates to false,
// then the constraint will be ignored.
func (c UniqueConstraint[T]) When(condition bool) UniqueConstraint[T] {
	c.isIgnored = !condition
	return c
}

// WhenGroups enables conditional validation of the constraint by using the validation groups.
func (c UniqueConstraint[T]) WhenGroups(groups ...string) UniqueConstraint[T] {
	c.groups = groups
	return c
}

func (c UniqueConstraint[T]) ValidateComparables(values []T, scope validation.Scope) error {
	if c.isIgnored || scope.IsIgnored(c.groups...) || is.Unique(values) {
		return nil
	}

	return scope.BuildViolation(c.err, c.messageTemplate).
		WithParameters(c.messageParameters...).
		Create()
}

func formatComparable[T comparable](value T) string {
	if s, ok := any(value).(string); ok {
		return `"` + s + `"`
	}

	return fmt.Sprint(value)
}
