// Copyright 2021 Igor Lazarev. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package validator contains Validator service singleton.
// It can be used in a custom application to perform the validation process.
package validator

import (
	"context"
	"time"

	"github.com/muonsoft/validation"
	"golang.org/x/text/language"
)

var validator, _ = validation.NewValidator()

// SetUp can be used to set up a new instance of singleton validator. Make sure you call this function once
// at the initialization of your application because it totally replaces validator instance.
func SetUp(options ...validation.ValidatorOption) (err error) {
	validator, err = validation.NewValidator(options...)

	return err
}

// SetOptions can be used to set up a new instance of singleton validator. Make sure you call this function once
// at the initialization of your application because it totally replaces validator instance.
//
// Deprecated: use SetUp function instead.
func SetOptions(options ...validation.ValidatorOption) (err error) {
	validator, err = validation.NewValidator(options...)

	return err
}

// Reset function recreates singleton validator. Generally, it can be used in tests.
//
// Deprecated: use SetUp function instead.
func Reset() {
	_ = SetUp()
}

// Validate is the main validation method. It accepts validation arguments. Arguments can be
// used to tune up the validation process or to pass values of a specific type.
func Validate(ctx context.Context, arguments ...validation.Argument) error {
	return validator.Validate(ctx, arguments...)
}

// ValidateValue is an alias for validating a single value of any supported type.
func ValidateValue(ctx context.Context, value interface{}, options ...validation.Option) error {
	return validator.ValidateValue(ctx, value, options...)
}

// ValidateBool is an alias for validating a single boolean value.
func ValidateBool(ctx context.Context, value bool, options ...validation.Option) error {
	return validator.ValidateBool(ctx, value, options...)
}

//
// // ValidateNumber is an alias for validating a single numeric value (integer or float).
// func ValidateNumber(ctx context.Context, value interface{}, options ...validation.Option) error {
// 	return validator.ValidateNumber(ctx, value, options...)
// }

// ValidateString is an alias for validating a single string value.
func ValidateString(ctx context.Context, value string, options ...validation.Option) error {
	return validator.ValidateString(ctx, value, options...)
}

// ValidateStrings is an alias for validating slice of strings.
func ValidateStrings(ctx context.Context, values []string, options ...validation.Option) error {
	return validator.ValidateStrings(ctx, values, options...)
}

// ValidateIterable is an alias for validating a single iterable value (an array, slice, or map).
func ValidateIterable(ctx context.Context, value interface{}, options ...validation.Option) error {
	return validator.ValidateIterable(ctx, value, options...)
}

// ValidateCountable is an alias for validating a single countable value (an array, slice, or map).
func ValidateCountable(ctx context.Context, count int, options ...validation.Option) error {
	return validator.ValidateCountable(ctx, count, options...)
}

// ValidateTime is an alias for validating a single time value.
func ValidateTime(ctx context.Context, value time.Time, options ...validation.Option) error {
	return validator.ValidateTime(ctx, value, options...)
}

// ValidateEach is an alias for validating each value of an iterable (an array, slice, or map).
func ValidateEach(ctx context.Context, value interface{}, options ...validation.Option) error {
	return validator.ValidateEach(ctx, value, options...)
}

// ValidateEachString is an alias for validating each value of a strings slice.
func ValidateEachString(ctx context.Context, strings []string, options ...validation.Option) error {
	return validator.ValidateEachString(ctx, strings, options...)
}

// ValidateValidatable is an alias for validating value that implements the Validatable interface.
func ValidateValidatable(ctx context.Context, validatable validation.Validatable, options ...validation.Option) error {
	return validator.ValidateValidatable(ctx, validatable, options...)
}

// ValidateBy is used to get the constraint from the internal validator store.
// If the constraint does not exist, then the validator will
// return a ConstraintNotFoundError during the validation process.
// For storing a constraint you should use the validation.StoredConstraint option.
func ValidateBy(constraintKey string) validation.Constraint {
	return validator.ValidateBy(constraintKey)
}

// WithGroups is used to execute conditional validation based on validation groups. It creates
// a new scoped validation with a given set of groups.
//
// By default, when validating an object all constraints of it will be checked whether or not
// they pass. In some cases, however, you will need to validate an object against
// only some specific group of constraints. To do this, you can organize each constraint
// into one or more validation groups and then apply validation against one group of constraints.
//
// Validation groups are working together only with validation groups passed
// to a constraint by WhenGroups() method. This method is implemented in all built-in constraints.
// If you want to use validation groups for your own constraints do not forget to implement
// this method in your constraint.
//
// Be careful, empty groups are considered as the default group.
// Its value is equal to the validation.DefaultGroup ("default").
func WithGroups(groups ...string) *validation.Validator {
	return validator.WithGroups(groups...)
}

// WithLanguage method creates a new scoped validator with a given language tag. All created violations
// will be translated into this language.
func WithLanguage(tag language.Tag) *validation.Validator {
	return validator.WithLanguage(tag)
}

// AtProperty method creates a new scoped validator with injected property name element to scope property path.
func AtProperty(name string) *validation.Validator {
	return validator.AtProperty(name)
}

// AtIndex method creates a new scoped validator with injected array index element to scope property path.
func AtIndex(index int) *validation.Validator {
	return validator.AtIndex(index)
}

// BuildViolation can be used to build a custom violation on the client-side.
func BuildViolation(ctx context.Context, code, message string) *validation.ViolationBuilder {
	return validator.BuildViolation(ctx, code, message)
}
