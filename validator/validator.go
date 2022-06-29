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

// Validate is the main validation method. It accepts validation arguments. executionContext can be
// used to tune up the validation process or to pass values of a specific type.
func Validate(ctx context.Context, arguments ...validation.Argument) error {
	return validator.Validate(ctx, arguments...)
}

// ValidateBool is an alias for validating a single boolean value.
func ValidateBool(ctx context.Context, value bool, constraints ...validation.BoolConstraint) error {
	return validator.ValidateBool(ctx, value, constraints...)
}

// ValidateInt is an alias for validating a single integer value.
func ValidateInt(ctx context.Context, value int, constraints ...validation.NumberConstraint[int]) error {
	return validator.Validate(ctx, validation.Number(value, constraints...))
}

// ValidateFloat is an alias for validating a single float value.
func ValidateFloat(ctx context.Context, value float64, constraints ...validation.NumberConstraint[float64]) error {
	return validator.Validate(ctx, validation.Number(value, constraints...))
}

// ValidateString is an alias for validating a single string value.
func ValidateString(ctx context.Context, value string, constraints ...validation.StringConstraint) error {
	return validator.ValidateString(ctx, value, constraints...)
}

// ValidateStrings is an alias for validating slice of strings.
func ValidateStrings(ctx context.Context, values []string, constraints ...validation.ComparablesConstraint[string]) error {
	return validator.ValidateStrings(ctx, values, constraints...)
}

// ValidateCountable is an alias for validating a single countable value (an array, slice, or map).
func ValidateCountable(ctx context.Context, count int, constraints ...validation.CountableConstraint) error {
	return validator.ValidateCountable(ctx, count, constraints...)
}

// ValidateTime is an alias for validating a single time value.
func ValidateTime(ctx context.Context, value time.Time, constraints ...validation.TimeConstraint) error {
	return validator.ValidateTime(ctx, value, constraints...)
}

// ValidateEachString is an alias for validating each value of a strings slice.
func ValidateEachString(ctx context.Context, strings []string, constraints ...validation.StringConstraint) error {
	return validator.ValidateEachString(ctx, strings, constraints...)
}

// ValidateIt is an alias for validating value that implements the Validatable interface.
func ValidateIt(ctx context.Context, validatable validation.Validatable) error {
	return validator.ValidateIt(ctx, validatable)
}

// GetConstraint is used to get the constraint from the internal validator store.
// If the constraint does not exist, then the validator will return nil.
// For storing a constraint you should use the StoredConstraint option.
//
// Experimental. This feature is experimental and may be changed in future versions.
func GetConstraint(key string) interface{} {
	return validator.GetConstraint(key)
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

// CreateViolation can be used to quickly create a custom violation on the client-side.
func CreateViolation(ctx context.Context, err error, message string, path ...validation.PropertyPathElement) validation.Violation {
	return validator.CreateViolation(ctx, err, message, path...)
}

// BuildViolation can be used to build a custom violation on the client-side.
func BuildViolation(ctx context.Context, err error, message string) *validation.ViolationBuilder {
	return validator.BuildViolation(ctx, err, message)
}

// BuildViolationList can be used to build a custom violation list on the client-side.
func BuildViolationList(ctx context.Context) *validation.ViolationListBuilder {
	return validator.BuildViolationList(ctx)
}
