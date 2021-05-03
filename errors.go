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

// ConstraintAlreadyStoredError is returned when trying to put a constraint
// in the validator store using an existing key.
type ConstraintAlreadyStoredError struct {
	Key string
}

func (err ConstraintAlreadyStoredError) Error() string {
	return fmt.Sprintf(`constraint with key "%s" already stored`, err.Key)
}

// ConstraintNotFoundError is returned when trying to get a constraint
// from the validator store using a non-existent key.
type ConstraintNotFoundError struct {
	Key string
}

func (err ConstraintNotFoundError) Error() string {
	return fmt.Sprintf(`constraint with key "%s" is not stored in the validator`, err.Key)
}

var (
	errDefaultLanguageNotLoaded      = errors.New("default language is not loaded")
	errThenBranchNotSet              = errors.New("then branch of conditional constraint not set")
	errSequentiallyConstraintsNotSet = errors.New("then branch of conditional constraint not set")
)
