package validator

import (
	"context"
	"time"

	"github.com/muonsoft/validation"
	"golang.org/x/text/language"
)

var validator, _ = validation.NewValidator()

func SetOptions(options ...validation.ValidatorOption) error {
	for _, setOption := range options {
		err := setOption(validator)
		if err != nil {
			return err
		}
	}

	return nil
}

func Reset() {
	validator, _ = validation.NewValidator()
}

func Validate(arguments ...validation.Argument) error {
	return validator.Validate(arguments...)
}

func ValidateValue(value interface{}, options ...validation.Option) error {
	return validator.ValidateValue(value, options...)
}

func ValidateBool(value *bool, options ...validation.Option) error {
	return validator.ValidateBool(value, options...)
}

func ValidateNumber(value interface{}, options ...validation.Option) error {
	return validator.ValidateNumber(value, options...)
}

func ValidateString(value *string, options ...validation.Option) error {
	return validator.ValidateString(value, options...)
}

func ValidateIterable(value interface{}, options ...validation.Option) error {
	return validator.ValidateIterable(value, options...)
}

func ValidateCountable(count int, options ...validation.Option) error {
	return validator.ValidateCountable(count, options...)
}

func ValidateTime(value *time.Time, options ...validation.Option) error {
	return validator.ValidateTime(value, options...)
}

func ValidateValidatable(validatable validation.Validatable, options ...validation.Option) error {
	return validator.ValidateValidatable(validatable, options...)
}

func ValidateEach(value interface{}, options ...validation.Option) error {
	return validator.ValidateEach(value, options...)
}

func ValidateEachString(strings []string, options ...validation.Option) error {
	return validator.ValidateEachString(strings, options...)
}

func WithContext(ctx context.Context) *validation.Validator {
	return validator.WithContext(ctx)
}

func WithLanguage(tag language.Tag) *validation.Validator {
	return validator.WithLanguage(tag)
}

func AtProperty(name string) *validation.Validator {
	return validator.AtProperty(name)
}

func AtIndex(index int) *validation.Validator {
	return validator.AtIndex(index)
}

func BuildViolation(code, message string) *validation.ViolationBuilder {
	return validator.BuildViolation(code, message)
}
