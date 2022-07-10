package it

import (
	"context"
	"time"

	"github.com/muonsoft/validation"
)

// NotBlankConstraint checks that a value is not blank: an empty string, an empty countable (slice/array/map),
// an empty generic number, generic comparable, false or nil. Nil behavior is configurable via WithAllowedNil() method.
// To check that a value is not nil only use NotNilConstraint.
type NotBlankConstraint[T comparable] struct {
	blank             T
	isIgnored         bool
	allowNil          bool
	groups            []string
	err               error
	messageTemplate   string
	messageParameters validation.TemplateParameterList
}

// IsNotBlank creates a NotBlankConstraint for checking that value is not empty.
func IsNotBlank() NotBlankConstraint[string] {
	return IsNotBlankComparable[string]()
}

// IsNotBlankNumber creates a NotBlankConstraint for checking that numeric value is not empty.
func IsNotBlankNumber[T validation.Numeric]() NotBlankConstraint[T] {
	return IsNotBlankComparable[T]()
}

// IsNotBlankComparable creates a NotBlankConstraint for checking that comparable value is not empty.
func IsNotBlankComparable[T comparable]() NotBlankConstraint[T] {
	return NotBlankConstraint[T]{
		err:             validation.ErrIsBlank,
		messageTemplate: validation.ErrIsBlank.Message(),
	}
}

// WithAllowedNil makes nil values valid.
func (c NotBlankConstraint[T]) WithAllowedNil() NotBlankConstraint[T] {
	c.allowNil = true
	return c
}

// When enables conditional validation of this constraint. If the expression evaluates to false,
// then the constraint will be ignored.
func (c NotBlankConstraint[T]) When(condition bool) NotBlankConstraint[T] {
	c.isIgnored = !condition
	return c
}

// WhenGroups enables conditional validation of the constraint by using the validation groups.
func (c NotBlankConstraint[T]) WhenGroups(groups ...string) NotBlankConstraint[T] {
	c.groups = groups
	return c
}

// WithError overrides default error for produced violation.
func (c NotBlankConstraint[T]) WithError(err error) NotBlankConstraint[T] {
	c.err = err
	return c
}

// WithMessage sets the violation message template. You can set custom template parameters
// for injecting its values into the final message.
func (c NotBlankConstraint[T]) WithMessage(template string, parameters ...validation.TemplateParameter) NotBlankConstraint[T] {
	c.messageTemplate = template
	c.messageParameters = parameters
	return c
}

func (c NotBlankConstraint[T]) ValidateBool(ctx context.Context, validator *validation.Validator, value *bool) error {
	if c.isIgnored || validator.IsIgnoredForGroups(c.groups...) {
		return nil
	}
	if c.allowNil && value == nil {
		return nil
	}
	if value != nil && *value {
		return nil
	}

	return c.newViolation(ctx, validator)
}

func (c NotBlankConstraint[T]) ValidateNumber(ctx context.Context, validator *validation.Validator, value *T) error {
	return c.ValidateComparable(ctx, validator, value)
}

func (c NotBlankConstraint[T]) ValidateString(ctx context.Context, validator *validation.Validator, value *string) error {
	if c.isIgnored || validator.IsIgnoredForGroups(c.groups...) {
		return nil
	}
	if c.allowNil && value == nil {
		return nil
	}
	if value != nil && *value != "" {
		return nil
	}

	return c.newViolation(ctx, validator)
}

func (c NotBlankConstraint[T]) ValidateComparable(ctx context.Context, validator *validation.Validator, value *T) error {
	if c.isIgnored || validator.IsIgnoredForGroups(c.groups...) {
		return nil
	}
	if c.allowNil && value == nil {
		return nil
	}
	if value != nil && *value != c.blank {
		return nil
	}

	return c.newViolation(ctx, validator)
}

func (c NotBlankConstraint[T]) ValidateCountable(ctx context.Context, validator *validation.Validator, count int) error {
	if c.isIgnored || validator.IsIgnoredForGroups(c.groups...) || count > 0 {
		return nil
	}

	return c.newViolation(ctx, validator)
}

func (c NotBlankConstraint[T]) ValidateTime(ctx context.Context, validator *validation.Validator, value *time.Time) error {
	if c.isIgnored || validator.IsIgnoredForGroups(c.groups...) {
		return nil
	}
	if c.allowNil && value == nil {
		return nil
	}

	if value != nil && !value.IsZero() {
		return nil
	}

	return c.newViolation(ctx, validator)
}

func (c NotBlankConstraint[T]) newViolation(ctx context.Context, validator *validation.Validator) validation.Violation {
	return validator.BuildViolation(ctx, c.err, c.messageTemplate).
		WithParameters(c.messageParameters...).
		Create()
}

// BlankConstraint checks that a value is blank: equal to false, nil, zero, an empty string, an empty
// slice, array, or a map.
type BlankConstraint[T comparable] struct {
	blank             T
	isIgnored         bool
	groups            []string
	err               error
	messageTemplate   string
	messageParameters validation.TemplateParameterList
}

// IsBlank creates a BlankConstraint for checking that value is empty.
func IsBlank() BlankConstraint[string] {
	return IsBlankComparable[string]()
}

// IsBlankNumber creates a BlankConstraint for checking that numeric value is nil or zero.
func IsBlankNumber[T validation.Numeric]() BlankConstraint[T] {
	return IsBlankComparable[T]()
}

// IsBlankComparable creates a BlankConstraint for checking that comparable value is not empty.
func IsBlankComparable[T comparable]() BlankConstraint[T] {
	return BlankConstraint[T]{
		err:             validation.ErrNotBlank,
		messageTemplate: validation.ErrNotBlank.Message(),
	}
}

// When enables conditional validation of this constraint. If the expression evaluates to false,
// then the constraint will be ignored.
func (c BlankConstraint[T]) When(condition bool) BlankConstraint[T] {
	c.isIgnored = !condition
	return c
}

// WhenGroups enables conditional validation of the constraint by using the validation groups.
func (c BlankConstraint[T]) WhenGroups(groups ...string) BlankConstraint[T] {
	c.groups = groups
	return c
}

// WithError overrides default error for produced violation.
func (c BlankConstraint[T]) WithError(err error) BlankConstraint[T] {
	c.err = err
	return c
}

// WithMessage sets the violation message template. You can set custom template parameters
// for injecting its values into the final message.
func (c BlankConstraint[T]) WithMessage(template string, parameters ...validation.TemplateParameter) BlankConstraint[T] {
	c.messageTemplate = template
	c.messageParameters = parameters
	return c
}

func (c BlankConstraint[T]) ValidateBool(ctx context.Context, validator *validation.Validator, value *bool) error {
	if c.isIgnored || validator.IsIgnoredForGroups(c.groups...) || value == nil || !*value {
		return nil
	}

	return c.newViolation(ctx, validator)
}

func (c BlankConstraint[T]) ValidateNumber(ctx context.Context, validator *validation.Validator, value *T) error {
	return c.ValidateComparable(ctx, validator, value)
}

func (c BlankConstraint[T]) ValidateString(ctx context.Context, validator *validation.Validator, value *string) error {
	if c.isIgnored || validator.IsIgnoredForGroups(c.groups...) || value == nil || *value == "" {
		return nil
	}

	return c.newViolation(ctx, validator)
}

func (c BlankConstraint[T]) ValidateComparable(ctx context.Context, validator *validation.Validator, value *T) error {
	if c.isIgnored || validator.IsIgnoredForGroups(c.groups...) || value == nil || *value == c.blank {
		return nil
	}

	return c.newViolation(ctx, validator)
}

func (c BlankConstraint[T]) ValidateCountable(ctx context.Context, validator *validation.Validator, count int) error {
	if c.isIgnored || validator.IsIgnoredForGroups(c.groups...) || count == 0 {
		return nil
	}

	return c.newViolation(ctx, validator)
}

func (c BlankConstraint[T]) ValidateTime(ctx context.Context, validator *validation.Validator, value *time.Time) error {
	if c.isIgnored || validator.IsIgnoredForGroups(c.groups...) || value == nil || value.IsZero() {
		return nil
	}

	return c.newViolation(ctx, validator)
}

func (c BlankConstraint[T]) newViolation(ctx context.Context, validator *validation.Validator) validation.Violation {
	return validator.BuildViolation(ctx, c.err, c.messageTemplate).
		WithParameters(c.messageParameters...).
		Create()
}

// NotNilConstraint checks that a value in not strictly equal to nil. To check that values in not blank use
// NotBlankConstraint.
type NotNilConstraint[T comparable] struct {
	isIgnored         bool
	groups            []string
	err               error
	messageTemplate   string
	messageParameters validation.TemplateParameterList
}

// IsNotNil creates a NotNilConstraint to check that a value is not strictly equal to nil.
func IsNotNil() NotNilConstraint[string] {
	return IsNotNilComparable[string]()
}

// IsNotNilNumber creates a NotNilConstraint to check that a numeric value is not strictly equal to nil.
func IsNotNilNumber[T validation.Numeric]() NotNilConstraint[T] {
	return IsNotNilComparable[T]()
}

// IsNotNilComparable creates a NotNilConstraint to check that a comparable value is not strictly equal to nil.
func IsNotNilComparable[T comparable]() NotNilConstraint[T] {
	return NotNilConstraint[T]{
		err:             validation.ErrIsNil,
		messageTemplate: validation.ErrIsNil.Message(),
	}
}

// When enables conditional validation of this constraint. If the expression evaluates to false,
// then the constraint will be ignored.
func (c NotNilConstraint[T]) When(condition bool) NotNilConstraint[T] {
	c.isIgnored = !condition
	return c
}

// WhenGroups enables conditional validation of the constraint by using the validation groups.
func (c NotNilConstraint[T]) WhenGroups(groups ...string) NotNilConstraint[T] {
	c.groups = groups
	return c
}

// WithError overrides default error for produced violation.
func (c NotNilConstraint[T]) WithError(err error) NotNilConstraint[T] {
	c.err = err
	return c
}

// WithMessage sets the violation message template. You can set custom template parameters
// for injecting its values into the final message.
func (c NotNilConstraint[T]) WithMessage(template string, parameters ...validation.TemplateParameter) NotNilConstraint[T] {
	c.messageTemplate = template
	c.messageParameters = parameters
	return c
}

func (c NotNilConstraint[T]) ValidateNil(ctx context.Context, validator *validation.Validator, isNil bool) error {
	if c.isIgnored || validator.IsIgnoredForGroups(c.groups...) || !isNil {
		return nil
	}

	return c.newViolation(ctx, validator)
}

func (c NotNilConstraint[T]) ValidateBool(ctx context.Context, validator *validation.Validator, value *bool) error {
	return c.ValidateNil(ctx, validator, value == nil)
}

func (c NotNilConstraint[T]) ValidateNumber(ctx context.Context, validator *validation.Validator, value *T) error {
	return c.ValidateNil(ctx, validator, value == nil)
}

func (c NotNilConstraint[T]) ValidateString(ctx context.Context, validator *validation.Validator, value *string) error {
	return c.ValidateNil(ctx, validator, value == nil)
}

func (c NotNilConstraint[T]) ValidateComparable(ctx context.Context, validator *validation.Validator, value *T) error {
	return c.ValidateNil(ctx, validator, value == nil)
}

func (c NotNilConstraint[T]) ValidateTime(ctx context.Context, validator *validation.Validator, value *time.Time) error {
	return c.ValidateNil(ctx, validator, value == nil)
}

func (c NotNilConstraint[T]) newViolation(ctx context.Context, validator *validation.Validator) validation.Violation {
	return validator.BuildViolation(ctx, c.err, c.messageTemplate).
		WithParameters(c.messageParameters...).
		Create()
}

// NilConstraint checks that a value in strictly equal to nil. To check that values in blank use
// BlankConstraint.
type NilConstraint[T comparable] struct {
	isIgnored         bool
	groups            []string
	err               error
	messageTemplate   string
	messageParameters validation.TemplateParameterList
}

// IsNil creates a NilConstraint to check that a value is strictly equal to nil.
func IsNil() NilConstraint[string] {
	return IsNilComparable[string]()
}

// IsNilNumber creates a NilConstraint to check that a numeric value is strictly equal to nil.
func IsNilNumber[T validation.Numeric]() NilConstraint[T] {
	return IsNilComparable[T]()
}

// IsNilComparable creates a NilConstraint to check that a comparable value is strictly equal to nil.
func IsNilComparable[T comparable]() NilConstraint[T] {
	return NilConstraint[T]{
		err:             validation.ErrNotNil,
		messageTemplate: validation.ErrNotNil.Message(),
	}
}

// When enables conditional validation of this constraint. If the expression evaluates to false,
// then the constraint will be ignored.
func (c NilConstraint[T]) When(condition bool) NilConstraint[T] {
	c.isIgnored = !condition
	return c
}

// WhenGroups enables conditional validation of the constraint by using the validation groups.
func (c NilConstraint[T]) WhenGroups(groups ...string) NilConstraint[T] {
	c.groups = groups
	return c
}

// WithError overrides default error for produced violation.
func (c NilConstraint[T]) WithError(err error) NilConstraint[T] {
	c.err = err
	return c
}

// WithMessage sets the violation message template. You can set custom template parameters
// for injecting its values into the final message.
func (c NilConstraint[T]) WithMessage(template string, parameters ...validation.TemplateParameter) NilConstraint[T] {
	c.messageTemplate = template
	c.messageParameters = parameters
	return c
}

func (c NilConstraint[T]) ValidateNil(ctx context.Context, validator *validation.Validator, isNil bool) error {
	if c.isIgnored || validator.IsIgnoredForGroups(c.groups...) || isNil {
		return nil
	}

	return c.newViolation(ctx, validator)
}

func (c NilConstraint[T]) ValidateBool(ctx context.Context, validator *validation.Validator, value *bool) error {
	return c.ValidateNil(ctx, validator, value == nil)
}

func (c NilConstraint[T]) ValidateNumber(ctx context.Context, validator *validation.Validator, value *T) error {
	return c.ValidateNil(ctx, validator, value == nil)
}

func (c NilConstraint[T]) ValidateString(ctx context.Context, validator *validation.Validator, value *string) error {
	return c.ValidateNil(ctx, validator, value == nil)
}

func (c NilConstraint[T]) ValidateComparable(ctx context.Context, validator *validation.Validator, value *T) error {
	return c.ValidateNil(ctx, validator, value == nil)
}

func (c NilConstraint[T]) ValidateTime(ctx context.Context, validator *validation.Validator, value *time.Time) error {
	return c.ValidateNil(ctx, validator, value == nil)
}

func (c NilConstraint[T]) newViolation(ctx context.Context, validator *validation.Validator) validation.Violation {
	return validator.BuildViolation(ctx, c.err, c.messageTemplate).
		WithParameters(c.messageParameters...).
		Create()
}

// BoolConstraint checks that a bool value in strictly equal to expected bool value.
type BoolConstraint struct {
	isIgnored         bool
	expected          bool
	groups            []string
	err               error
	messageTemplate   string
	messageParameters validation.TemplateParameterList
}

// IsTrue creates a BoolConstraint to check that a value is not strictly equal to true.
func IsTrue() BoolConstraint {
	return BoolConstraint{
		expected:        true,
		err:             validation.ErrNotTrue,
		messageTemplate: validation.ErrNotTrue.Message(),
	}
}

// IsFalse creates a BoolConstraint to check that a value is not strictly equal to false.
func IsFalse() BoolConstraint {
	return BoolConstraint{
		expected:        false,
		err:             validation.ErrNotFalse,
		messageTemplate: validation.ErrNotFalse.Message(),
	}
}

// When enables conditional validation of this constraint. If the expression evaluates to false,
// then the constraint will be ignored.
func (c BoolConstraint) When(condition bool) BoolConstraint {
	c.isIgnored = !condition
	return c
}

// WhenGroups enables conditional validation of the constraint by using the validation groups.
func (c BoolConstraint) WhenGroups(groups ...string) BoolConstraint {
	c.groups = groups
	return c
}

// WithError overrides default error for produced violation.
func (c BoolConstraint) WithError(err error) BoolConstraint {
	c.err = err
	return c
}

// WithMessage sets the violation message template. You can set custom template parameters
// for injecting its values into the final message.
func (c BoolConstraint) WithMessage(template string, parameters ...validation.TemplateParameter) BoolConstraint {
	c.messageTemplate = template
	c.messageParameters = parameters
	return c
}

func (c BoolConstraint) ValidateBool(ctx context.Context, validator *validation.Validator, value *bool) error {
	if c.isIgnored || validator.IsIgnoredForGroups(c.groups...) || value == nil || *value == c.expected {
		return nil
	}

	return validator.BuildViolation(ctx, c.err, c.messageTemplate).
		WithParameters(c.messageParameters...).
		Create()
}
