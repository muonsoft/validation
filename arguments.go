package validation

import (
	"time"

	"golang.org/x/text/language"
)

// Argument used to set up the validation process. It is used to set up the current validation scope and to pass
// arguments for validation values.
type Argument interface {
	setUp(ctx *executionContext)
}

// Nil argument is used to validate nil values of any nillable types.
func Nil(isNil bool, constraints ...NilConstraint) ValidatorArgument {
	return NewArgument(validateNil(isNil, constraints))
}

// NilProperty argument is an alias for Nil that automatically adds property name to the current scope.
func NilProperty(name string, isNil bool, constraints ...NilConstraint) ValidatorArgument {
	return NewArgument(validateNil(isNil, constraints)).With(PropertyName(name))
}

// Bool argument is used to validate boolean values.
func Bool(value bool, constraints ...BoolConstraint) ValidatorArgument {
	return NewArgument(validateBool(&value, constraints))
}

// BoolProperty argument is an alias for Bool that automatically adds property name to the current scope.
func BoolProperty(name string, value bool, constraints ...BoolConstraint) ValidatorArgument {
	return NewArgument(validateBool(&value, constraints)).With(PropertyName(name))
}

// NilBool argument is used to validate nillable boolean values.
func NilBool(value *bool, constraints ...BoolConstraint) ValidatorArgument {
	return NewArgument(validateBool(value, constraints))
}

// NilBoolProperty argument is an alias for NilBool that automatically adds property name to the current scope.
func NilBoolProperty(name string, value *bool, constraints ...BoolConstraint) ValidatorArgument {
	return NewArgument(validateBool(value, constraints)).With(PropertyName(name))
}

// Number argument is used to validate numbers.
func Number[T Numeric](value T, constraints ...NumberConstraint[T]) ValidatorArgument {
	return NewArgument(validateNumber(&value, constraints))
}

// NumberProperty argument is an alias for Number that automatically adds property name to the current scope.
func NumberProperty[T Numeric](name string, value T, constraints ...NumberConstraint[T]) ValidatorArgument {
	return NewArgument(validateNumber(&value, constraints)).With(PropertyName(name))
}

// NilNumber argument is used to validate nillable numbers.
func NilNumber[T Numeric](value *T, constraints ...NumberConstraint[T]) ValidatorArgument {
	return NewArgument(validateNumber(value, constraints))
}

// NilNumberProperty argument is an alias for NilNumber that automatically adds property name to the current scope.
func NilNumberProperty[T Numeric](name string, value *T, constraints ...NumberConstraint[T]) ValidatorArgument {
	return NewArgument(validateNumber(value, constraints)).With(PropertyName(name))
}

// String argument is used to validate strings.
func String(value string, constraints ...StringConstraint) ValidatorArgument {
	return NewArgument(validateString(&value, constraints))
}

// StringProperty argument is an alias for String that automatically adds property name to the current scope.
func StringProperty(name string, value string, constraints ...StringConstraint) ValidatorArgument {
	return NewArgument(validateString(&value, constraints)).With(PropertyName(name))
}

// NilString argument is used to validate nillable strings.
func NilString(value *string, constraints ...StringConstraint) ValidatorArgument {
	return NewArgument(validateString(value, constraints))
}

// NilStringProperty argument is an alias for NilString that automatically adds property name to the current scope.
func NilStringProperty(name string, value *string, constraints ...StringConstraint) ValidatorArgument {
	return NewArgument(validateString(value, constraints)).With(PropertyName(name))
}

// Countable argument can be used to validate size of an array, slice, or map. You can pass result of len()
// function as an argument.
func Countable(count int, constraints ...CountableConstraint) ValidatorArgument {
	return NewArgument(validateCountable(count, constraints))
}

// CountableProperty argument is an alias for Countable that automatically adds property name to the current scope.
func CountableProperty(name string, count int, constraints ...CountableConstraint) ValidatorArgument {
	return NewArgument(validateCountable(count, constraints)).With(PropertyName(name))
}

// Time argument is used to validate time.Time value.
func Time(value time.Time, constraints ...TimeConstraint) ValidatorArgument {
	return NewArgument(validateTime(&value, constraints))
}

// TimeProperty argument is an alias for Time that automatically adds property name to the current scope.
func TimeProperty(name string, value time.Time, constraints ...TimeConstraint) ValidatorArgument {
	return NewArgument(validateTime(&value, constraints)).With(PropertyName(name))
}

// NilTime argument is used to validate nillable time.Time value.
func NilTime(value *time.Time, constraints ...TimeConstraint) ValidatorArgument {
	return NewArgument(validateTime(value, constraints))
}

// NilTimeProperty argument is an alias for NilTime that automatically adds property name to the current scope.
func NilTimeProperty(name string, value *time.Time, constraints ...TimeConstraint) ValidatorArgument {
	return NewArgument(validateTime(value, constraints)).With(PropertyName(name))
}

// Valid is used to run validation on the Validatable type. This method is recommended
// to build a complex validation process.
func Valid(value Validatable) ValidatorArgument {
	return NewArgument(validateIt(value))
}

// ValidProperty argument is an alias for Valid that automatically adds property name to the current scope.
func ValidProperty(name string, value Validatable) ValidatorArgument {
	return NewArgument(validateIt(value)).With(PropertyName(name))
}

// ValidSlice is a generic argument used to run validation on the slice of Validatable types.
// This method is recommended to build a complex validation process.
func ValidSlice[T Validatable](values []T) ValidatorArgument {
	return NewArgument(validateSlice(values))
}

// ValidSliceProperty argument is an alias for ValidSlice that automatically adds property name to the current scope.
func ValidSliceProperty[T Validatable](name string, values []T) ValidatorArgument {
	return NewArgument(validateSlice(values)).With(PropertyName(name))
}

// ValidMap is a generic argument used to run validation on the map of Validatable types.
// This method is recommended to build a complex validation process.
func ValidMap[T Validatable](values map[string]T) ValidatorArgument {
	return NewArgument(validateMap(values))
}

// ValidMapProperty argument is an alias for ValidSlice that automatically adds property name to the current scope.
func ValidMapProperty[T Validatable](name string, values map[string]T) ValidatorArgument {
	return NewArgument(validateMap(values)).With(PropertyName(name))
}

// Comparable argument is used to validate generic comparable value.
func Comparable[T comparable](value T, constraints ...ComparableConstraint[T]) ValidatorArgument {
	return NewArgument(validateComparable(&value, constraints))
}

// ComparableProperty argument is an alias for Comparable that automatically adds property name to the current scope.
func ComparableProperty[T comparable](name string, value T, constraints ...ComparableConstraint[T]) ValidatorArgument {
	return NewArgument(validateComparable(&value, constraints)).With(PropertyName(name))
}

// NilComparable argument is used to validate nillable generic comparable value.
func NilComparable[T comparable](value *T, constraints ...ComparableConstraint[T]) ValidatorArgument {
	return NewArgument(validateComparable(value, constraints))
}

// NilComparableProperty argument is an alias for NilComparable that automatically adds property name to the current scope.
func NilComparableProperty[T comparable](name string, value *T, constraints ...ComparableConstraint[T]) ValidatorArgument {
	return NewArgument(validateComparable(value, constraints)).With(PropertyName(name))
}

// Comparables argument is used to validate generic comparable types.
func Comparables[T comparable](values []T, constraints ...ComparablesConstraint[T]) ValidatorArgument {
	return NewArgument(validateComparables(values, constraints))
}

// ComparablesProperty argument is an alias for Comparables that automatically adds property name to the current scope.
func ComparablesProperty[T comparable](name string, values []T, constraints ...ComparablesConstraint[T]) ValidatorArgument {
	return NewArgument(validateComparables(values, constraints)).With(PropertyName(name))
}

// EachString is used to validate a slice of strings.
func EachString(values []string, constraints ...StringConstraint) ValidatorArgument {
	return NewArgument(validateEachString(values, constraints))
}

// EachStringProperty argument is an alias for EachString that automatically adds property name to the current scope.
func EachStringProperty(name string, values []string, constraints ...StringConstraint) ValidatorArgument {
	return NewArgument(validateEachString(values, constraints)).With(PropertyName(name))
}

// EachNumber is used to validate a slice of numbers.
func EachNumber[T Numeric](values []T, constraints ...NumberConstraint[T]) ValidatorArgument {
	return NewArgument(validateEachNumber(values, constraints))
}

// EachNumberProperty argument is an alias for EachString that automatically adds property name to the current scope.
func EachNumberProperty[T Numeric](name string, values []T, constraints ...NumberConstraint[T]) ValidatorArgument {
	return NewArgument(validateEachNumber(values, constraints)).With(PropertyName(name))
}

// EachComparable is used to validate a slice of generic comparables.
func EachComparable[T comparable](values []T, constraints ...ComparableConstraint[T]) ValidatorArgument {
	return NewArgument(validateEachComparable(values, constraints))
}

// EachComparableProperty argument is an alias for EachComparable that automatically adds property name to the current scope.
func EachComparableProperty[T comparable](name string, values []T, constraints ...ComparableConstraint[T]) ValidatorArgument {
	return NewArgument(validateEachComparable(values, constraints)).With(PropertyName(name))
}

// CheckNoViolations is a special argument that checks err for violations. If err contains Violation or ViolationList
// then these violations will be appended into returned violation list from the validator. If err contains an error
// that does not implement an error interface, then the validation process will be terminated and
// this error will be returned.
func CheckNoViolations(err error) ValidatorArgument {
	return NewArgument(func(scope Scope) (*ViolationList, error) {
		return unwrapViolationList(err)
	})
}

// Check argument can be useful for quickly checking the result of some simple expression
// that returns a boolean value.
func Check(isValid bool) Checker {
	return Checker{
		isValid:         isValid,
		err:             ErrNotValid,
		messageTemplate: ErrNotValid.Template(),
	}
}

// CheckProperty argument is an alias for Check that automatically adds property name to the current scope.
// It is useful to apply a simple checks on structs.
func CheckProperty(name string, isValid bool) Checker {
	return Checker{
		propertyName:    name,
		isValid:         isValid,
		err:             ErrNotValid,
		messageTemplate: ErrNotValid.Template(),
	}
}

// Language argument sets the current language for translation of a violation message.
func Language(tag language.Tag) Argument {
	return argumentFunc(func(ctx *executionContext) {
		ctx.scope.language = tag
	})
}

type ValidateOnScopeFunc func(scope Scope) (*ViolationList, error)

// NewArgument can be used to implement validation functional arguments for the specific types.
func NewArgument(validate ValidateOnScopeFunc) ValidatorArgument {
	return ValidatorArgument{validate: validate}
}

// NewTypedArgument creates a generic validation argument that can help implement the validation
// argument for client-side types.
func NewTypedArgument[T any](v T, constraints ...Constraint[T]) ValidatorArgument {
	return NewArgument(func(scope Scope) (*ViolationList, error) {
		violations := NewViolationList()

		for _, constraint := range constraints {
			err := violations.AppendFromError(constraint.Validate(scope.context, scope.Validator(), v))
			if err != nil {
				return nil, err
			}
		}

		return violations, nil
	})
}

// ValidatorArgument is common implementation of Argument that is used to run validation
// process on given argument.
type ValidatorArgument struct {
	isIgnored bool
	validate  ValidateOnScopeFunc
	options   []Option
}

// With returns a copy of ValidatorArgument with appended options.
func (arg ValidatorArgument) With(options ...Option) ValidatorArgument {
	arg.options = append(arg.options, options...)
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
		ctx.addValidator(arg.options, arg.validate)
	}
}

type argumentFunc func(ctx *executionContext)

func (f argumentFunc) setUp(ctx *executionContext) {
	f(ctx)
}

// Checker is an argument that can be useful for quickly checking the result of
// some simple expression that returns a boolean value.
type Checker struct {
	options           []Option
	isIgnored         bool
	isValid           bool
	propertyName      string
	groups            []string
	err               error
	messageTemplate   string
	messageParameters TemplateParameterList
}

// With returns a copy of Checker with appended options.
func (c Checker) With(options ...Option) Checker {
	c.options = append(c.options, options...)
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
	arguments.addValidator(c.options, c.validate)
}

func (c Checker) validate(scope Scope) (*ViolationList, error) {
	if c.isValid || c.isIgnored || scope.IsIgnored(c.groups...) {
		return nil, nil
	}
	if c.propertyName != "" {
		scope = scope.AtProperty(c.propertyName)
	}

	violation := scope.BuildViolation(c.err, c.messageTemplate).
		WithParameters(c.messageParameters...).
		Create()

	return NewViolationList(violation), nil
}
