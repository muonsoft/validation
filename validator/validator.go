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

// SetOptions can be used to set up a singleton validator. Make sure you call this function once
// at the initialization of your application.
func SetOptions(options ...validation.ValidatorOption) error {
	for _, setOption := range options {
		err := setOption(validator)
		if err != nil {
			return err
		}
	}

	return nil
}

// Reset function recreates singleton validator. Generally, it can be used in tests.
func Reset() {
	validator, _ = validation.NewValidator()
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

// ValidateNumber is an alias for validating a single numeric value (integer or float).
func ValidateNumber(ctx context.Context, value interface{}, options ...validation.Option) error {
	return validator.ValidateNumber(ctx, value, options...)
}

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
