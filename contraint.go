package validation

import (
	"github.com/muonsoft/validation/generic"

	"time"
)

type Constraint interface {
	Option
	GetName() string
}

// NilConstraint is used for constraints that needs to check value for nil. In common case
// you have no need to implement it in your constraints because nil values should be ignored.
type NilConstraint interface {
	Constraint
	ValidateNil(scope Scope) error
}

type BoolConstraint interface {
	Constraint
	ValidateBool(value *bool, scope Scope) error
}

type NumberConstraint interface {
	Constraint
	ValidateNumber(value generic.Number, scope Scope) error
}

type StringConstraint interface {
	Constraint
	ValidateString(value *string, scope Scope) error
}

type IterableConstraint interface {
	Constraint
	ValidateIterable(value generic.Iterable, scope Scope) error
}

type CountableConstraint interface {
	Constraint
	ValidateCountable(count int, scope Scope) error
}

type TimeConstraint interface {
	Constraint
	ValidateTime(value *time.Time, scope Scope) error
}
