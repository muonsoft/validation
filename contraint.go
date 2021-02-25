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
	ValidateNil(options Options) error
}

type BoolConstraint interface {
	ValidateBool(value *bool, options Options) error
}

type NumberConstraint interface {
	ValidateNumber(value generic.Number, options Options) error
}

type StringConstraint interface {
	ValidateString(value *string, options Options) error
}

type IterableConstraint interface {
	ValidateIterable(value generic.Iterable, options Options) error
}

type CountableConstraint interface {
	ValidateCountable(count int, options Options) error
}

type TimeConstraint interface {
	ValidateTime(time *time.Time, options Options) error
}
