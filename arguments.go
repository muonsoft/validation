package validation

import (
	"context"
	"time"
)

// Argument used to set up the validation process. It is used to set up the current validation context and to pass
// arguments for validation values.
type Argument interface {
	setUp(ctx *executionContext)
}

// Nil argument is used to validate nil values of any nillable types.
func Nil(isNil bool, constraints ...NilConstraint) ValidatorArgument {
	return NewArgument(validateNil(isNil, constraints))
}

// NilProperty argument is an alias for [Nil] that automatically adds property name to the current validation context.
func NilProperty(name string, isNil bool, constraints ...NilConstraint) ValidatorArgument {
	return NewArgument(validateNil(isNil, constraints)).At(PropertyName(name))
}

// Bool argument is used to validate boolean values.
func Bool(value bool, constraints ...BoolConstraint) ValidatorArgument {
	return NewArgument(validateBool(&value, constraints))
}

// BoolProperty argument is an alias for [Bool] that automatically adds property name to the current validation context.
func BoolProperty(name string, value bool, constraints ...BoolConstraint) ValidatorArgument {
	return NewArgument(validateBool(&value, constraints)).At(PropertyName(name))
}

// NilBool argument is used to validate nillable boolean values.
func NilBool(value *bool, constraints ...BoolConstraint) ValidatorArgument {
	return NewArgument(validateBool(value, constraints))
}

// NilBoolProperty argument is an alias for [NilBool] that automatically adds property name to the current validation context.
func NilBoolProperty(name string, value *bool, constraints ...BoolConstraint) ValidatorArgument {
	return NewArgument(validateBool(value, constraints)).At(PropertyName(name))
}

// Number argument is used to validate numbers.
func Number[T Numeric](value T, constraints ...NumberConstraint[T]) ValidatorArgument {
	return NewArgument(validateNumber(&value, constraints))
}

// NumberProperty argument is an alias for [Number] that automatically adds property name to the current validation context.
func NumberProperty[T Numeric](name string, value T, constraints ...NumberConstraint[T]) ValidatorArgument {
	return NewArgument(validateNumber(&value, constraints)).At(PropertyName(name))
}

// NilNumber argument is used to validate nillable numbers.
func NilNumber[T Numeric](value *T, constraints ...NumberConstraint[T]) ValidatorArgument {
	return NewArgument(validateNumber(value, constraints))
}

// NilNumberProperty argument is an alias for [NilNumber] that automatically adds property name to the current validation context.
func NilNumberProperty[T Numeric](name string, value *T, constraints ...NumberConstraint[T]) ValidatorArgument {
	return NewArgument(validateNumber(value, constraints)).At(PropertyName(name))
}

// String argument is used to validate strings.
func String(value string, constraints ...StringConstraint) ValidatorArgument {
	return NewArgument(validateString(&value, constraints))
}

// StringProperty argument is an alias for [String] that automatically adds property name to the current validation context.
func StringProperty(name string, value string, constraints ...StringConstraint) ValidatorArgument {
	return NewArgument(validateString(&value, constraints)).At(PropertyName(name))
}

// NilString argument is used to validate nillable strings.
func NilString(value *string, constraints ...StringConstraint) ValidatorArgument {
	return NewArgument(validateString(value, constraints))
}

// NilStringProperty argument is an alias for [NilString] that automatically adds property name to the current validation context.
func NilStringProperty(name string, value *string, constraints ...StringConstraint) ValidatorArgument {
	return NewArgument(validateString(value, constraints)).At(PropertyName(name))
}

// Countable argument can be used to validate size of an array, slice, or map. You can pass result of len()
// function as an argument.
func Countable(count int, constraints ...CountableConstraint) ValidatorArgument {
	return NewArgument(validateCountable(count, constraints))
}

// CountableProperty argument is an alias for [Countable] that automatically adds property name to the current validation context.
func CountableProperty(name string, count int, constraints ...CountableConstraint) ValidatorArgument {
	return NewArgument(validateCountable(count, constraints)).At(PropertyName(name))
}

// Time argument is used to validate [time.Time] value.
func Time(value time.Time, constraints ...TimeConstraint) ValidatorArgument {
	return NewArgument(validateTime(&value, constraints))
}

// TimeProperty argument is an alias for [Time] that automatically adds property name to the current validation context.
func TimeProperty(name string, value time.Time, constraints ...TimeConstraint) ValidatorArgument {
	return NewArgument(validateTime(&value, constraints)).At(PropertyName(name))
}

// NilTime argument is used to validate nillable [time.Time] value.
func NilTime(value *time.Time, constraints ...TimeConstraint) ValidatorArgument {
	return NewArgument(validateTime(value, constraints))
}

// NilTimeProperty argument is an alias for [NilTime] that automatically adds property name to the current validation context.
func NilTimeProperty(name string, value *time.Time, constraints ...TimeConstraint) ValidatorArgument {
	return NewArgument(validateTime(value, constraints)).At(PropertyName(name))
}

// Valid is used to run validation on the [Validatable] type. This method is recommended
// to build a complex validation process.
func Valid(value Validatable) ValidatorArgument {
	return NewArgument(validateIt(value))
}

// ValidProperty argument is an alias for [Valid] that automatically adds property name to the current validation context.
func ValidProperty(name string, value Validatable) ValidatorArgument {
	return NewArgument(validateIt(value)).At(PropertyName(name))
}

// ValidSlice is a generic argument used to run validation on the slice of [Validatable] types.
// This method is recommended to build a complex validation process.
func ValidSlice[T Validatable](values []T) ValidatorArgument {
	return NewArgument(validateSlice(values))
}

// ValidSliceProperty argument is an alias for [ValidSlice] that automatically adds property name to the current validation context.
func ValidSliceProperty[T Validatable](name string, values []T) ValidatorArgument {
	return NewArgument(validateSlice(values)).At(PropertyName(name))
}

// ValidMap is a generic argument used to run validation on the map of [Validatable] types.
// This method is recommended to build a complex validation process.
func ValidMap[T Validatable](values map[string]T) ValidatorArgument {
	return NewArgument(validateMap(values))
}

// ValidMapProperty argument is an alias for [ValidSlice] that automatically adds property name to the current validation context.
func ValidMapProperty[T Validatable](name string, values map[string]T) ValidatorArgument {
	return NewArgument(validateMap(values)).At(PropertyName(name))
}

// Comparable argument is used to validate generic comparable value.
func Comparable[T comparable](value T, constraints ...ComparableConstraint[T]) ValidatorArgument {
	return NewArgument(validateComparable(&value, constraints))
}

// ComparableProperty argument is an alias for [Comparable] that automatically adds property name to the current validation context.
func ComparableProperty[T comparable](name string, value T, constraints ...ComparableConstraint[T]) ValidatorArgument {
	return NewArgument(validateComparable(&value, constraints)).At(PropertyName(name))
}

// NilComparable argument is used to validate nillable generic comparable value.
func NilComparable[T comparable](value *T, constraints ...ComparableConstraint[T]) ValidatorArgument {
	return NewArgument(validateComparable(value, constraints))
}

// NilComparableProperty argument is an alias for [NilComparable] that automatically adds property name to the current validation context.
func NilComparableProperty[T comparable](name string, value *T, constraints ...ComparableConstraint[T]) ValidatorArgument {
	return NewArgument(validateComparable(value, constraints)).At(PropertyName(name))
}

// Comparables argument is used to validate generic comparable types.
func Comparables[T comparable](values []T, constraints ...ComparablesConstraint[T]) ValidatorArgument {
	return NewArgument(validateComparables(values, constraints))
}

// ComparablesProperty argument is an alias for [Comparables] that automatically adds property name to the current validation context.
func ComparablesProperty[T comparable](name string, values []T, constraints ...ComparablesConstraint[T]) ValidatorArgument {
	return NewArgument(validateComparables(values, constraints)).At(PropertyName(name))
}

// EachString is used to validate a slice of strings.
func EachString(values []string, constraints ...StringConstraint) ValidatorArgument {
	return NewArgument(validateEachString(values, constraints))
}

// EachStringProperty argument is an alias for [EachString] that automatically adds property name to the current validation context.
func EachStringProperty(name string, values []string, constraints ...StringConstraint) ValidatorArgument {
	return NewArgument(validateEachString(values, constraints)).At(PropertyName(name))
}

// EachNumber is used to validate a slice of numbers.
func EachNumber[T Numeric](values []T, constraints ...NumberConstraint[T]) ValidatorArgument {
	return NewArgument(validateEachNumber(values, constraints))
}

// EachNumberProperty argument is an alias for [EachNumber] that automatically adds property name to the current validation context.
func EachNumberProperty[T Numeric](name string, values []T, constraints ...NumberConstraint[T]) ValidatorArgument {
	return NewArgument(validateEachNumber(values, constraints)).At(PropertyName(name))
}

// EachComparable is used to validate a slice of generic comparables.
func EachComparable[T comparable](values []T, constraints ...ComparableConstraint[T]) ValidatorArgument {
	return NewArgument(validateEachComparable(values, constraints))
}

// EachComparableProperty argument is an alias for [EachComparable] that automatically adds property name to the current validation context.
func EachComparableProperty[T comparable](name string, values []T, constraints ...ComparableConstraint[T]) ValidatorArgument {
	return NewArgument(validateEachComparable(values, constraints)).At(PropertyName(name))
}

// CheckNoViolations is a special argument that checks err for violations. If err contains [Violation] or [ViolationList]
// then these violations will be appended into returned violation list from the validator. If err contains an error
// that does not implement an error interface, then the validation process will be terminated and
// this error will be returned.
func CheckNoViolations(err error) ValidatorArgument {
	return NewArgument(func(ctx context.Context, validator *Validator) (*ViolationList, error) {
		return unwrapViolationList(err)
	})
}

// Check argument can be useful for quickly checking the result of some simple expression
// that returns a boolean value.
func Check(isValid bool) Checker {
	return Checker{
		isValid:         isValid,
		err:             ErrNotValid,
		messageTemplate: ErrNotValid.Message(),
	}
}

// CheckProperty argument is an alias for [Check] that automatically adds property name to the current validation context.
// It is useful to apply a simple checks on structs.
func CheckProperty(name string, isValid bool) Checker {
	return Check(isValid).At(PropertyName(name))
}

type ValidateFunc func(ctx context.Context, validator *Validator) (*ViolationList, error)

// NewArgument can be used to implement validation functional arguments for the specific types.
func NewArgument(validate ValidateFunc) ValidatorArgument {
	return ValidatorArgument{validate: validate}
}

// This creates a generic validation argument that can help implement the validation
// argument for client-side types.
func This[T any](v T, constraints ...Constraint[T]) ValidatorArgument {
	return NewArgument(func(ctx context.Context, validator *Validator) (*ViolationList, error) {
		violations := NewViolationList()

		for _, constraint := range constraints {
			err := violations.AppendFromError(constraint.Validate(ctx, validator, v))
			if err != nil {
				return nil, err
			}
		}

		return violations, nil
	})
}

// ValidatorArgument is common implementation of [Argument] that is used to run validation
// process on given argument.
type ValidatorArgument struct {
	isIgnored bool
	validate  ValidateFunc
	path      []PropertyPathElement
}

// At returns a copy of [ValidatorArgument] with appended property path suffix.
func (arg ValidatorArgument) At(path ...PropertyPathElement) ValidatorArgument {
	arg.path = append(arg.path, path...)
	return arg
}

// When enables conditional validation of this argument. If the expression evaluates to false,
// then the argument will be ignored.
func (arg ValidatorArgument) When(condition bool) ValidatorArgument {
	arg.isIgnored = !condition
	return arg
}

func (arg ValidatorArgument) setUp(ctx *executionContext) {
	if !arg.isIgnored {
		ctx.addValidation(arg.validate, arg.path...)
	}
}

// Checker is an argument that can be useful for quickly checking the result of
// some simple expression that returns a boolean value.
type Checker struct {
	isIgnored         bool
	isValid           bool
	path              []PropertyPathElement
	groups            []string
	err               error
	messageTemplate   string
	messageParameters TemplateParameterList
}

// At returns a copy of [Checker] with appended property path suffix.
func (c Checker) At(path ...PropertyPathElement) Checker {
	c.path = append(c.path, path...)
	return c
}

// When enables conditional validation of this constraint. If the expression evaluates to false,
// then the constraint will be ignored.
func (c Checker) When(condition bool) Checker {
	c.isIgnored = !condition
	return c
}

// WhenGroups enables conditional validation of the constraint by using the validation groups.
func (c Checker) WhenGroups(groups ...string) Checker {
	c.groups = groups
	return c
}

// WithError overrides default code for produced violation.
func (c Checker) WithError(err error) Checker {
	c.err = err
	return c
}

// WithMessage sets the violation message template. You can set custom template parameters
// for injecting its values into the final message.
func (c Checker) WithMessage(template string, parameters ...TemplateParameter) Checker {
	c.messageTemplate = template
	c.messageParameters = parameters
	return c
}

func (c Checker) setUp(arguments *executionContext) {
	arguments.addValidation(c.validate, c.path...)
}

func (c Checker) validate(ctx context.Context, validator *Validator) (*ViolationList, error) {
	if c.isValid || c.isIgnored || validator.IsIgnoredForGroups(c.groups...) {
		return nil, nil
	}

	violation := validator.BuildViolation(ctx, c.err, c.messageTemplate).
		WithParameters(c.messageParameters...).
		Create()

	return NewViolationList(violation), nil
}
