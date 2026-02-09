// Copyright 2021 Igor Lazarev. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package validator contains Validator service singleton.
// It can be used in a custom application to perform the validation process.
package validator

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/muonsoft/validation"
	"golang.org/x/text/language"
)

var validatorPtr atomic.Pointer[validation.Validator]

func init() {
	v, _ := validation.NewValidator()
	validatorPtr.Store(v)
}

// Default returns the default validator instance.
func Default() *validation.Validator {
	return validatorPtr.Load()
}

// SetDefault sets the default validator. Call it once at application initialization.
// It completely replaces the current default validator.
func SetDefault(validator *validation.Validator) {
	validatorPtr.Store(validator)
}

// Instance returns the instance of the singleton validator.
//
// Deprecated: Use Default() instead. Instance will be removed in stable v1.
func Instance() *validation.Validator {
	return Default()
}

// SetUp can be used to set up a new instance of singleton validator. Make sure you call this function once
// at the initialization of your application because it totally replaces validator instance.
//
// Deprecated: Use SetDefault() instead. SetUp will be removed in stable v1.
func SetUp(options ...validation.ValidatorOption) error {
	v, err := validation.NewValidator(options...)
	if err != nil {
		return err
	}

	validatorPtr.Store(v)

	return nil
}

// Validate is the main validation method. It accepts validation arguments. executionContext can be
// used to tune up the validation process or to pass values of a specific type.
func Validate(ctx context.Context, arguments ...validation.Argument) error {
	return Default().Validate(ctx, arguments...)
}

// ValidateBool is an alias for validating a single boolean value.
func ValidateBool(ctx context.Context, value bool, constraints ...validation.BoolConstraint) error {
	return Default().ValidateBool(ctx, value, constraints...)
}

// ValidateInt is an alias for validating a single integer value.
func ValidateInt(ctx context.Context, value int, constraints ...validation.NumberConstraint[int]) error {
	return Default().Validate(ctx, validation.Number(value, constraints...))
}

// ValidateFloat is an alias for validating a single float value.
func ValidateFloat(ctx context.Context, value float64, constraints ...validation.NumberConstraint[float64]) error {
	return Default().Validate(ctx, validation.Number(value, constraints...))
}

// ValidateString is an alias for validating a single string value.
func ValidateString(ctx context.Context, value string, constraints ...validation.StringConstraint) error {
	return Default().ValidateString(ctx, value, constraints...)
}

// ValidateStrings is an alias for validating slice of strings.
func ValidateStrings(ctx context.Context, values []string, constraints ...validation.ComparablesConstraint[string]) error {
	return Default().ValidateStrings(ctx, values, constraints...)
}

// ValidateCountable is an alias for validating a single countable value (an array, slice, or map).
func ValidateCountable(ctx context.Context, count int, constraints ...validation.CountableConstraint) error {
	return Default().ValidateCountable(ctx, count, constraints...)
}

// ValidateTime is an alias for validating a single time value.
func ValidateTime(ctx context.Context, value time.Time, constraints ...validation.TimeConstraint) error {
	return Default().ValidateTime(ctx, value, constraints...)
}

// ValidateEachString is an alias for validating each value of a strings slice.
func ValidateEachString(ctx context.Context, strings []string, constraints ...validation.StringConstraint) error {
	return Default().ValidateEachString(ctx, strings, constraints...)
}

// ValidateIt is an alias for validating value that implements the Validatable interface.
func ValidateIt(ctx context.Context, validatable validation.Validatable) error {
	return Default().ValidateIt(ctx, validatable)
}

// WithGroups is used to execute conditional validation based on validation groups. It creates
// a new context validator with a given set of groups.
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
	return Default().WithGroups(groups...)
}

// WithLanguage method creates a new context validator with a given language tag. All created violations
// will be translated into this language.
func WithLanguage(tag language.Tag) *validation.Validator {
	return Default().WithLanguage(tag)
}

// At method creates a new context validator with appended property path.
func At(path ...validation.PropertyPathElement) *validation.Validator {
	return Default().At(path...)
}

// AtProperty method creates a new context validator with appended property name to the property path.
func AtProperty(name string) *validation.Validator {
	return Default().AtProperty(name)
}

// AtIndex method creates a new context validator with appended array index to the property path.
func AtIndex(index int) *validation.Validator {
	return Default().AtIndex(index)
}

// CreateViolation can be used to quickly create a custom violation on the client-side.
func CreateViolation(ctx context.Context, err error, message string, path ...validation.PropertyPathElement) validation.Violation {
	return Default().CreateViolation(ctx, err, message, path...)
}

// BuildViolation can be used to build a custom violation on the client-side.
func BuildViolation(ctx context.Context, err error, message string) *validation.ViolationBuilder {
	return Default().BuildViolation(ctx, err, message)
}

// BuildViolationList can be used to build a custom violation list on the client-side.
func BuildViolationList(ctx context.Context) *validation.ViolationListBuilder {
	return Default().BuildViolationList(ctx)
}
