package validation

import (
	"fmt"
	"reflect"
)

type InapplicableConstraintError struct {
	Constraint Constraint
	ValueType  string
}

func (err *InapplicableConstraintError) Error() string {
	return fmt.Sprintf("%s cannot be applied to %s", err.Constraint.GetName(), err.ValueType)
}

func newInapplicableConstraintError(constraint Constraint, valueType string) *InapplicableConstraintError {
	return &InapplicableConstraintError{
		Constraint: constraint,
		ValueType:  valueType,
	}
}

type NotValidatableError struct {
	Value reflect.Value
}

func (err *NotValidatableError) Error() string {
	return fmt.Sprintf("value of type %v is not validatable", err.Value.Kind())
}
