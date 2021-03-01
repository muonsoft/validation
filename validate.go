package validation

import (
	"reflect"
	"time"
)

type Validatable interface {
	Validate(scope Scope) error
}

func Validate(arguments ...Argument) error {
	return validator.Validate(arguments...)
}

func InScope(scope Scope) *Validator {
	return validator.InScope(scope)
}

func ValidateValue(value interface{}, options ...Option) error {
	return validator.ValidateValue(value, options...)
}

func ValidateBool(value *bool, options ...Option) error {
	return validator.ValidateBool(value, options...)
}

func ValidateNumber(value interface{}, options ...Option) error {
	return validator.ValidateNumber(value, options...)
}

func ValidateString(value *string, options ...Option) error {
	return validator.ValidateString(value, options...)
}

func ValidateIterable(value interface{}, options ...Option) error {
	return validator.ValidateIterable(value, options...)
}

func ValidateCountable(count int, options ...Option) error {
	return validator.ValidateCountable(count, options...)
}

func ValidateTime(value *time.Time, options ...Option) error {
	return validator.ValidateTime(value, options...)
}

func ValidateValidatable(validatable Validatable, options ...Option) error {
	return validator.ValidateValidatable(validatable, options...)
}

func ValidateEach(value interface{}, options ...Option) error {
	return validator.ValidateEach(value, options...)
}

func ValidateEachString(strings []string, options ...Option) error {
	return validator.ValidateEachString(strings, options...)
}

func Filter(violations ...error) error {
	filteredViolations := make(ViolationList, 0, len(violations))

	for _, err := range violations {
		addErr := filteredViolations.AddFromError(err)
		if addErr != nil {
			return addErr
		}
	}

	return filteredViolations.AsError()
}

var validatableType = reflect.TypeOf((*Validatable)(nil)).Elem()
