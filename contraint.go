package validation

import (
	"context"
	"time"
)

// Numeric is used as a type parameter for numeric values.
type Numeric interface {
	~float32 | ~float64 |
		~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

// Constraint is a generic interface for client-side typed constraints.
type Constraint[T any] interface {
	Validate(ctx context.Context, validator *Validator, v T) error
}

// NilConstraint is used for a special cases to check a value for nil.
type NilConstraint interface {
	ValidateNil(ctx context.Context, validator *Validator, isNil bool) error
}

// BoolConstraint is used to build constraints for boolean values validation.
type BoolConstraint interface {
	ValidateBool(ctx context.Context, validator *Validator, value *bool) error
}

// NumberConstraint is used to build constraints for numeric values validation.
type NumberConstraint[T Numeric] interface {
	ValidateNumber(ctx context.Context, validator *Validator, value *T) error
}

// StringConstraint is used to build constraints for string values validation.
type StringConstraint interface {
	ValidateString(ctx context.Context, validator *Validator, value *string) error
}

// ComparableConstraint is used to build constraints for generic comparable value validation.
type ComparableConstraint[T comparable] interface {
	ValidateComparable(ctx context.Context, validator *Validator, value *T) error
}

// ComparablesConstraint is used to build constraints for generic comparable values validation.
type ComparablesConstraint[T comparable] interface {
	ValidateComparables(ctx context.Context, validator *Validator, values []T) error
}

// CountableConstraint is used to build constraints for simpler validation of iterable elements count.
type CountableConstraint interface {
	ValidateCountable(ctx context.Context, validator *Validator, count int) error
}

// TimeConstraint is used to build constraints for date/time validation.
type TimeConstraint interface {
	ValidateTime(ctx context.Context, validator *Validator, value *time.Time) error
}

// StringFuncConstraint can be used to create constraints for validating string values
// based on function with signature func(string) bool.
type StringFuncConstraint struct {
	isIgnored         bool
	isValid           func(string) bool
	groups            []string
	err               error
	messageTemplate   string
	messageParameters TemplateParameterList
}

// OfStringBy creates a new string constraint from a function with signature func(string) bool.
func OfStringBy(isValid func(string) bool) StringFuncConstraint {
	return StringFuncConstraint{
		isValid:         isValid,
		err:             ErrNotValid,
		messageTemplate: ErrNotValid.Message(),
	}
}

// WithError overrides default error for produced violation.
func (c StringFuncConstraint) WithError(err error) StringFuncConstraint {
	c.err = err
	return c
}

// WithMessage sets the violation message template. You can set custom template parameters
// for injecting its values into the final message. Also, you can use default parameters:
//
//	{{ value }} - the current (invalid) value.
func (c StringFuncConstraint) WithMessage(template string, parameters ...TemplateParameter) StringFuncConstraint {
	c.messageTemplate = template
	c.messageParameters = parameters
	return c
}

// When enables conditional validation of this constraint. If the expression evaluates to false,
// then the constraint will be ignored.
func (c StringFuncConstraint) When(condition bool) StringFuncConstraint {
	c.isIgnored = !condition
	return c
}

// WhenGroups enables conditional validation of the constraint by using the validation groups.
func (c StringFuncConstraint) WhenGroups(groups ...string) StringFuncConstraint {
	c.groups = groups
	return c
}

func (c StringFuncConstraint) ValidateString(ctx context.Context, validator *Validator, value *string) error {
	if c.isIgnored || validator.IsIgnoredForGroups(c.groups...) || value == nil || *value == "" || c.isValid(*value) {
		return nil
	}

	return validator.BuildViolation(ctx, c.err, c.messageTemplate).
		WithParameters(
			c.messageParameters.Prepend(
				TemplateParameter{Key: "{{ value }}", Value: *value},
			)...,
		).
		WithParameter("{{ value }}", *value).
		Create()
}
