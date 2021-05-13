package validation

import (
	"github.com/muonsoft/validation/code"
	"github.com/muonsoft/validation/generic"
	"github.com/muonsoft/validation/message"

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

// CustomStringConstraint can be used to create custom constraints for validating string values
// based on function with signature func(string) bool.
type CustomStringConstraint struct {
	isIgnored       bool
	isValid         func(string) bool
	name            string
	code            string
	messageTemplate string
}

// NewCustomStringConstraint creates a new string constraint from a function with signature func(string) bool.
// Optional parameters can be used to set up constraint name (first parameter), violation code (second),
// message template (third). All other parameters are ignored.
func NewCustomStringConstraint(isValid func(string) bool, parameters ...string) CustomStringConstraint {
	constraint := CustomStringConstraint{
		isValid:         isValid,
		name:            "CustomStringConstraint",
		code:            code.NotValid,
		messageTemplate: message.NotValid,
	}

	if len(parameters) > 0 {
		constraint.name = parameters[0]
	}
	if len(parameters) > 1 {
		constraint.code = parameters[1]
	}
	if len(parameters) > 2 {
		constraint.messageTemplate = parameters[2]
	}

	return constraint
}

// SetUp always returns no error.
func (c CustomStringConstraint) SetUp() error {
	return nil
}

// Name is the constraint name. It can be set via first parameter of function NewCustomStringConstraint.
func (c CustomStringConstraint) Name() string {
	return c.name
}

// Message sets the violation message template. You can use template parameters
// for injecting its values into the final message:
//
//	{{ value }} - the current (invalid) value.
func (c CustomStringConstraint) Message(message string) CustomStringConstraint {
	c.messageTemplate = message
	return c
}

// When enables conditional validation of this constraint. If the expression evaluates to false,
// then the constraint will be ignored.
func (c CustomStringConstraint) When(condition bool) CustomStringConstraint {
	c.isIgnored = !condition
	return c
}

func (c CustomStringConstraint) ValidateString(value *string, scope Scope) error {
	if c.isIgnored || value == nil || *value == "" || c.isValid(*value) {
		return nil
	}

	return scope.BuildViolation(c.code, c.messageTemplate).
		AddParameter("{{ value }}", *value).
		CreateViolation()
}

// ConditionalConstraint is used for conditional validation.
// Use the When function to initiate a conditional check.
// If the condition is true, then the constraints passed through the Then function will be applied.
// Otherwise, the constraints passed through the Else function will apply.
type ConditionalConstraint struct {
	condition       bool
	thenConstraints []Constraint
	elseConstraints []Constraint
}

// When function using to initiate a conditional check.
func When(condition bool) ConditionalConstraint {
	return ConditionalConstraint{
		condition: condition,
	}
}

// Then function applied the constraints if the condition is true.
//
// Example
//  v := "foo"
//	err := validator.ValidateString(
//		&value,
//		validation.When(true).
//		Then(
//			it.Matches(regexp.MustCompile(`^\\w$`)),
//		),
//	)
func (c ConditionalConstraint) Then(constraints ...Constraint) ConditionalConstraint {
	c.thenConstraints = constraints
	return c
}

// Else function applied the constraints if the condition is false.
//
// Example
//  v := "foo"
//	err := validator.ValidateString(
//		&value,
//		validation.When(false).
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

// SetUp will return an error if the constraints to apply is empty.
func (c ConditionalConstraint) SetUp() error {
	if len(c.thenConstraints) == 0 {
		return errThenBranchNotSet
	}

	return nil
}

func (c ConditionalConstraint) validate(
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

// SequentiallyConstraint is used to set constraints allowing to interrupt the validation once
// the first violation is raised.
type SequentiallyConstraint struct {
	constraints []Constraint
}

// Sequentially function using to set of constraints that should be validated step-by-step.
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

func (c SequentiallyConstraint) validate(
	scope Scope,
	violations *ViolationList,
	validate ValidateByConstraintFunc,
) error {
	for _, constraint := range c.constraints {
		err := validate(constraint, scope)
		if err != nil {
			return violations.AppendFromError(err)
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
