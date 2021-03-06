package validation

import (
	"reflect"
)

type Validatable interface {
	Validate(validator *Validator) error
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
