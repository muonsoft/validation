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

// ConditionalConstraint is used to construct constraints using a conditional construct.
type ConditionalConstraint struct {
	condition       bool
	thenConstraints []Constraint
	elseConstraints []Constraint
}

// When creates a ConditionalConstraint that set condition.
//
// Example
//  v := "name"
//	err := validator.ValidateString(
//		&value,
//		validation.When(utf8.RuneCountInString(value) <= 3).
//		Then(
//			it.Matches(regexp.MustCompile(`^\\w$`)),
//		),
//	)
func When(condition bool) ConditionalConstraint {
	return ConditionalConstraint{
		condition: condition,
	}
}

// Then creates a ConditionalConstraint that set a list of then branch constraints.
//
// Example
//  v := "name"
//	err := validator.ValidateString(
//		&value,
//		validation.When(utf8.RuneCountInString(value) <= 3).
//		Then(
//			it.Matches(regexp.MustCompile(`^\\w$`)),
//		),
//	)
func (c ConditionalConstraint) Then(constraints ...Constraint) ConditionalConstraint {
	c.thenConstraints = constraints
	return c
}

// Else creates a ConditionalConstraint that set a list of else branch constraints.
//
// Example
//  v := "name"
//	err := validator.ValidateString(
//		&value,
//		validation.When(utf8.RuneCountInString(value) <= 3).
//		Then(
//			it.Matches(regexp.MustCompile(`^\\w$`)),
//		).
//		Else(
//			it.Matches(regexp.MustCompile(`^\\d$`)),
//		),
//	)
func (c ConditionalConstraint) Else(constraints ...Constraint) ConditionalConstraint {
	c.elseConstraints = constraints
	return c
}

// Name is the constraint name.
func (c ConditionalConstraint) Name() string {
	return "ConditionalConstraint"
}

// SetUp will return an error if the list of then branch constraints is empty.
func (c ConditionalConstraint) SetUp() error {
	if len(c.thenConstraints) == 0 {
		return errThenBranchNotSet
	}

	return nil
}

func (c *ConditionalConstraint) validateConditionConstraints(
	scope Scope,
	violations *ViolationList,
	validate ValidateByConstraintFunc,
) error {
	var constraints []Constraint
	if c.condition {
		constraints = c.thenConstraints
	} else {
		constraints = c.elseConstraints
	}

	for _, constraint := range constraints {
		err := violations.AppendFromError(validate(constraint, scope))
		if err != nil {
			return err
		}
	}

	return nil
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

type SequentiallyConstraint struct {
	constraints []Constraint
}

func Sequentially(constraints ...Constraint) SequentiallyConstraint {
	return SequentiallyConstraint{
		constraints: constraints,
	}
}

// Name is the constraint name.
func (c SequentiallyConstraint) Name() string {
	return "SequentiallyConstraint"
}

// SetUp will return an error if the list of constraints is empty.
func (c SequentiallyConstraint) SetUp() error {
	if len(c.constraints) == 0 {
		return errSequentiallyConstraintsNotSet
	}
	return nil
}

func (c *SequentiallyConstraint) validateSequentiallyConstraints(
	scope Scope,
	violations *ViolationList,
	validate ValidateByConstraintFunc,
) error {
	var isViolation bool
	for _, constraint := range c.constraints {
		err := validate(constraint, scope)
		if err != nil {
			isViolation = true
		}
		err = violations.AppendFromError(err)
		if err != nil {
			return err
		}

		if isViolation {
			return nil
		}
	}

	return nil
}
