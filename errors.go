package validation

import (
	"errors"
	"fmt"
	"reflect"
)

// InapplicableConstraintError occurs when trying to use constraint on not applicable values.
// For example, if you are trying to compare slice with a number.
type InapplicableConstraintError struct {
	Constraint Constraint
	ValueType  string
}

func (err InapplicableConstraintError) Error() string {
	return fmt.Sprintf("%s cannot be applied to %s", err.Constraint.Name(), err.ValueType)
}

// NewInapplicableConstraintError helps to create a error on trying to use constraint on not applicable values.
func NewInapplicableConstraintError(constraint Constraint, valueType string) InapplicableConstraintError {
	return InapplicableConstraintError{
		Constraint: constraint,
		ValueType:  valueType,
	}
}

// NotValidatableError occurs when validator cannot determine type by reflection or it is not supported
// by validator.
type NotValidatableError struct {
	Value reflect.Value
}

func (err NotValidatableError) Error() string {
	return fmt.Sprintf("value of type %v is not validatable", err.Value.Kind())
}

var errDefaultLanguageNotLoaded = errors.New("default language is not loaded")
