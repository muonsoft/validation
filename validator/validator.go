package validator

import (
	"time"

	"github.com/muonsoft/validation"
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

func GetScope() validation.Scope {
	return validator.GetScope()
}

func Validate(arguments ...validation.Argument) error {
	return validator.Validate(arguments...)
}

func InScope(scope validation.Scope) *validation.Validator {
	return validator.InScope(scope)
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
