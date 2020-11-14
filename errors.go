package validation

import (
	"fmt"
	"reflect"
)

type ErrInapplicableConstraint struct {
	Code string
	Type string
}

func (err *ErrInapplicableConstraint) Error() string {
	return fmt.Sprintf("constraint with GetCode '%s' cannot be applied to %s", err.Code, err.Type)
}

type ErrNotValidatable struct {
	Value reflect.Value
}

func (err *ErrNotValidatable) Error() string {
	return fmt.Sprintf("value of type %v is not validatable", err.Value.Kind())
}
