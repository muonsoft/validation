package validation

import (
	"reflect"
	"time"
)

type Validatable interface {
	Validate(options ...Option) error
}

func Validate(value interface{}, options ...Option) error {
	return validator.Validate(value, options...)
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

func ValidateTime(time *time.Time, options ...Option) error {
	return validator.ValidateTime(time, options...)
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

func WithOptions(options ...Option) (*Validator, error) {
	return validator.WithOptions(options...)
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
