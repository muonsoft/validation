package validation

import (
	"fmt"
	"reflect"
)

type InapplicableConstraintError struct {
	Code string
	Type string
}

func (err *InapplicableConstraintError) Error() string {
	return fmt.Sprintf("constraint with code '%s' cannot be applied to %s", err.Code, err.Type)
}

func newInapplicableConstraintError(constraint Constraint, valueType string) *InapplicableConstraintError {
	return &InapplicableConstraintError{
		Code: constraint.GetCode(),
		Type: valueType,
	}
}

type NotValidatableError struct {
	Value reflect.Value
}

func (err *NotValidatableError) Error() string {
	return fmt.Sprintf("value of type %v is not validatable", err.Value.Kind())
}
