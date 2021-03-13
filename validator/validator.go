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
func Validate(arguments ...validation.Argument) error {
	return validator.Validate(arguments...)
}

// ValidateValue is an alias for validating a single value of any supported type.
func ValidateValue(value interface{}, options ...validation.Option) error {
	return validator.ValidateValue(value, options...)
}

// ValidateBool is an alias for validating a single boolean value.
func ValidateBool(value *bool, options ...validation.Option) error {
	return validator.ValidateBool(value, options...)
}

// ValidateNumber is an alias for validating a single numeric value (integer or float).
func ValidateNumber(value interface{}, options ...validation.Option) error {
	return validator.ValidateNumber(value, options...)
}

// ValidateString is an alias for validating a single string value.
func ValidateString(value *string, options ...validation.Option) error {
	return validator.ValidateString(value, options...)
}

// ValidateIterable is an alias for validating a single iterable value (an array, slice, or map).
func ValidateIterable(value interface{}, options ...validation.Option) error {
	return validator.ValidateIterable(value, options...)
}

// ValidateCountable is an alias for validating a single countable value (an array, slice, or map).
func ValidateCountable(count int, options ...validation.Option) error {
	return validator.ValidateCountable(count, options...)
}

// ValidateTime is an alias for validating a single time value.
func ValidateTime(value *time.Time, options ...validation.Option) error {
	return validator.ValidateTime(value, options...)
}

// ValidateEach is an alias for validating each value of an iterable (an array, slice, or map).
func ValidateEach(value interface{}, options ...validation.Option) error {
	return validator.ValidateEach(value, options...)
}

// ValidateEachString is an alias for validating each value of a strings slice.
func ValidateEachString(strings []string, options ...validation.Option) error {
	return validator.ValidateEachString(strings, options...)
}

// ValidateValidatable is an alias for validating value that implements the Validatable interface.
func ValidateValidatable(validatable validation.Validatable, options ...validation.Option) error {
	return validator.ValidateValidatable(validatable, options...)
}

// WithContext method creates a new scoped validator with a given context. You can use this method to pass
// a context value to all used constraints.

// Example
//  err := validator.WithContext(request.Context()).Validate(
//      String(&s, it.IsNotBlank()), // now all called constraints will use passed context in their methods
//  )
func WithContext(ctx context.Context) *validation.Validator {
	return validator.WithContext(ctx)
}

// WithLanguage method creates a new scoped validator with a given language tag. All created violations
// will be translated into this language.
//
// Example
//  err := validator.WithLanguage(language.Russian).Validate(
//      validation.ValidateString(&s, it.IsNotBlank()), // violation from this constraint will be translated
//  )
func WithLanguage(tag language.Tag) *validation.Validator {
	return validator.WithLanguage(tag)
}

// AtProperty mthod creates a new scoped validator with injected property name element to scope property path.
func AtProperty(name string) *validation.Validator {
	return validator.AtProperty(name)
}

// AtIndex method creates a new scoped validator with injected array index element to scope property path.
func AtIndex(index int) *validation.Validator {
	return validator.AtIndex(index)
}

// BuildViolation can be used to build a custom violation on the client-side.
//
// Example
//  err := validator.BuildViolation("", "").
//      SetParameter("key", "value").
//      GetViolation()
func BuildViolation(code, message string) *validation.ViolationBuilder {
	return validator.BuildViolation(code, message)
}
