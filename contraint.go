package validation

import (
	"github.com/muonsoft/validation/generic"

	"time"
)

// Constraint is the base interface to build validation constraints.
type Constraint interface {
	Option
	// Name is a constraint name that can be used in internal errors.
	Name() string
}

// NilConstraint is used for constraints that need to check value for nil. In common case
// you do not need to implement it in your constraints because nil values should be ignored.
type NilConstraint interface {
	Constraint
	ValidateNil(scope Scope) error
}

// BoolConstraint is used to build constraints for boolean values validation.
type BoolConstraint interface {
	Constraint
	ValidateBool(value *bool, scope Scope) error
}

// NumberConstraint is used to build constraints for numeric values validation.
//
// At this moment working with numbers is based on reflection.
// Be aware. This constraint is subject to be changed after generics implementation in Go.
type NumberConstraint interface {
	Constraint
	ValidateNumber(value generic.Number, scope Scope) error
}

// StringConstraint is used to build constraints for string values validation.
type StringConstraint interface {
	Constraint
	ValidateString(value *string, scope Scope) error
}

// IterableConstraint is used to build constraints for validation of iterables (arrays, slices, or maps).
//
// At this moment working with numbers is based on reflection.
// Be aware. This constraint is subject to be changed after generics implementation in Go.
type IterableConstraint interface {
	Constraint
	ValidateIterable(value generic.Iterable, scope Scope) error
}

// CountableConstraint is used to build constraints for simpler validation of iterable elements count.
type CountableConstraint interface {
	Constraint
	ValidateCountable(count int, scope Scope) error
}

// TimeConstraint is used to build constraints for date/time validation.
type TimeConstraint interface {
	Constraint
	ValidateTime(value *time.Time, scope Scope) error
}

type ConditionalConstraint struct {
	condition       bool
	thenConstraints []Constraint
	elseConstraints []Constraint
}

func When(condition bool) ConditionalConstraint {
	return ConditionalConstraint{
		condition: condition,
	}
}

func (c ConditionalConstraint) Then(constraints ...Constraint) ConditionalConstraint {
	c.thenConstraints = constraints
	return c
}

func (c ConditionalConstraint) Else(constraints ...Constraint) ConditionalConstraint {
	c.elseConstraints = constraints
	return c
}

func (c ConditionalConstraint) SetUp() error {
	if len(c.thenConstraints) == 0 {
		return errThenBranchNotSet
	}

	return nil
}

func (c ConditionalConstraint) Name() string {
	return "ConditionalConstraint"
}

type notFoundConstraint struct {
	key string
}

func (c notFoundConstraint) SetUp() error {
	return ConstraintNotFoundError{Key: c.key}
}

func (c notFoundConstraint) Name() string {
	return "notFoundConstraint"
}
