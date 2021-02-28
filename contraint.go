package validation

import (
	"github.com/muonsoft/validation/generic"
)

type Constraint interface {
	Option
	GetName() string
}

// NilConstraint is used for constraints that needs to check value for nil. In common case
// you have no need to implement it in your constraints because nil values should be ignored.
type NilConstraint interface {
	ValidateNil(scope Scope) error
}

type BoolConstraint interface {
	ValidateBool(value *bool, scope Scope) error
}

type NumberConstraint interface {
	ValidateNumber(value generic.Number, scope Scope) error
}

type StringConstraint interface {
	ValidateString(value *string, scope Scope) error
}

type IterableConstraint interface {
	ValidateIterable(value generic.Iterable, scope Scope) error
}

type CountableConstraint interface {
	ValidateCountable(count int, scope Scope) error
}
