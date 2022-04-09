package validation

import (
	"time"

	"github.com/muonsoft/validation/code"
	"github.com/muonsoft/validation/message"
	"golang.org/x/text/language"
)

// Argument used to set up the validation process. It is used to set up the current validation scope and to pass
// arguments for validation values.
type Argument interface {
	setUp(ctx *executionContext) error
}

type argumentFunc func(ctx *executionContext) error

func (f argumentFunc) setUp(ctx *executionContext) error {
	return f(ctx)
}

// Nil argument is used to validate nil values of any nillable types.
func Nil(isNil bool, constraints ...NilConstraint) NilArgument {
	return NilArgument{isNil: isNil, constraints: constraints}
}

// NilProperty argument is an alias for Nil that automatically adds property name to the current scope.
func NilProperty(name string, isNil bool, constraints ...NilConstraint) NilArgument {
	return NilArgument{isNil: isNil, constraints: constraints, options: []Option{PropertyName(name)}}
}

// Bool argument is used to validate boolean values.
func Bool(value bool, constraints ...BoolConstraint) BoolArgument {
	return BoolArgument{value: &value, constraints: constraints}
}

// BoolProperty argument is an alias for Bool that automatically adds property name to the current scope.
func BoolProperty(name string, value bool, constraints ...BoolConstraint) BoolArgument {
	return BoolArgument{value: &value, constraints: constraints, options: []Option{PropertyName(name)}}
}

// NilBool argument is used to validate nillable boolean values.
func NilBool(value *bool, constraints ...BoolConstraint) BoolArgument {
	return BoolArgument{value: value, constraints: constraints}
}

// NilBoolProperty argument is an alias for NilBool that automatically adds property name to the current scope.
func NilBoolProperty(name string, value *bool, constraints ...BoolConstraint) BoolArgument {
	return BoolArgument{value: value, constraints: constraints, options: []Option{PropertyName(name)}}
}

// Number argument is used to validate numbers.
func Number[T Numeric](value T, constraints ...NumberConstraint[T]) NumberArgument[T] {
	return NumberArgument[T]{value: &value, constraints: constraints}
}

// NumberProperty argument is an alias for Number that automatically adds property name to the current scope.
func NumberProperty[T Numeric](name string, value T, constraints ...NumberConstraint[T]) NumberArgument[T] {
	return NumberArgument[T]{value: &value, constraints: constraints, options: []Option{PropertyName(name)}}
}

// NilNumber argument is used to validate nillable numbers.
func NilNumber[T Numeric](value *T, constraints ...NumberConstraint[T]) NumberArgument[T] {
	return NumberArgument[T]{value: value, constraints: constraints}
}

// NilNumberProperty argument is an alias for NilNumber that automatically adds property name to the current scope.
func NilNumberProperty[T Numeric](name string, value *T, constraints ...NumberConstraint[T]) NumberArgument[T] {
	return NumberArgument[T]{value: value, constraints: constraints, options: []Option{PropertyName(name)}}
}

// String argument is used to validate strings.
func String(value string, constraints ...StringConstraint) StringArgument {
	return StringArgument{value: &value, constraints: constraints}
}

// StringProperty argument is an alias for String that automatically adds property name to the current scope.
func StringProperty(name string, value string, constraints ...StringConstraint) StringArgument {
	return StringArgument{value: &value, constraints: constraints, options: []Option{PropertyName(name)}}
}

// NilString argument is used to validate nillable strings.
func NilString(value *string, constraints ...StringConstraint) StringArgument {
	return StringArgument{value: value, constraints: constraints}
}

// NilStringProperty argument is an alias for NilString that automatically adds property name to the current scope.
func NilStringProperty(name string, value *string, constraints ...StringConstraint) StringArgument {
	return StringArgument{value: value, constraints: constraints, options: []Option{PropertyName(name)}}
}

// Countable argument can be used to validate size of an array, slice, or map. You can pass result of len()
// function as an argument.
func Countable(count int, constraints ...CountableConstraint) CountableArgument {
	return CountableArgument{count: count, constraints: constraints}
}

// CountableProperty argument is an alias for Countable that automatically adds property name to the current scope.
func CountableProperty(name string, count int, constraints ...CountableConstraint) CountableArgument {
	return CountableArgument{count: count, constraints: constraints, options: []Option{PropertyName(name)}}
}

// Time argument is used to validate time.Time value.
func Time(value time.Time, constraints ...TimeConstraint) TimeArgument {
	return TimeArgument{value: &value, constraints: constraints}
}

// TimeProperty argument is an alias for Time that automatically adds property name to the current scope.
func TimeProperty(name string, value time.Time, constraints ...TimeConstraint) TimeArgument {
	return TimeArgument{value: &value, constraints: constraints, options: []Option{PropertyName(name)}}
}

// NilTime argument is used to validate nillable time.Time value.
func NilTime(value *time.Time, constraints ...TimeConstraint) TimeArgument {
	return TimeArgument{value: value, constraints: constraints}
}

// NilTimeProperty argument is an alias for NilTime that automatically adds property name to the current scope.
func NilTimeProperty(name string, value *time.Time, constraints ...TimeConstraint) TimeArgument {
	return TimeArgument{value: value, constraints: constraints, options: []Option{PropertyName(name)}}
}

// EachString is used to validate a slice of strings.
func EachString(values []string, constraints ...StringConstraint) EachStringArgument {
	return EachStringArgument{values: values, constraints: constraints}
}

// EachStringProperty argument is an alias for EachString that automatically adds property name to the current scope.
func EachStringProperty(name string, values []string, constraints ...StringConstraint) EachStringArgument {
	return EachStringArgument{values: values, constraints: constraints, options: []Option{PropertyName(name)}}
}

// EachNumber is used to validate a slice of numbers.
func EachNumber[T Numeric](values []T, constraints ...NumberConstraint[T]) EachNumberArgument[T] {
	return EachNumberArgument[T]{values: values, constraints: constraints}
}

// EachNumberProperty argument is an alias for EachString that automatically adds property name to the current scope.
func EachNumberProperty[T Numeric](name string, values []T, constraints ...NumberConstraint[T]) EachNumberArgument[T] {
	return EachNumberArgument[T]{values: values, constraints: constraints, options: []Option{PropertyName(name)}}
}

// Valid is used to run validation on the Validatable type. This method is recommended
// to build a complex validation process.
func Valid(value Validatable) ValidArgument {
	return ValidArgument{value: value}
}

// ValidProperty argument is an alias for Valid that automatically adds property name to the current scope.
func ValidProperty(name string, value Validatable) ValidArgument {
	return ValidArgument{value: value, options: []Option{PropertyName(name)}}
}

// ValidSlice is a generic argument used to run validation on the slice of Validatable types.
// This method is recommended to build a complex validation process.
func ValidSlice[T Validatable](values []T) ValidSliceArgument[T] {
	return ValidSliceArgument[T]{values: values}
}

// ValidSliceProperty argument is an alias for ValidSlice that automatically adds property name to the current scope.
func ValidSliceProperty[T Validatable](name string, values []T) ValidSliceArgument[T] {
	return ValidSliceArgument[T]{values: values, options: []Option{PropertyName(name)}}
}

// ValidMap is a generic argument used to run validation on the map of Validatable types.
// This method is recommended to build a complex validation process.
func ValidMap[T Validatable](values map[string]T) ValidMapArgument[T] {
	return ValidMapArgument[T]{values: values}
}

// ValidMapProperty argument is an alias for ValidSlice that automatically adds property name to the current scope.
func ValidMapProperty[T Validatable](name string, values map[string]T) ValidMapArgument[T] {
	return ValidMapArgument[T]{values: values, options: []Option{PropertyName(name)}}
}

// Comparable argument is used to validate generic comparable value.
func Comparable[T comparable](value T, constraints ...ComparableConstraint[T]) ComparableArgument[T] {
	return ComparableArgument[T]{value: &value, constraints: constraints}
}

// ComparableProperty argument is an alias for Comparable that automatically adds property name to the current scope.
func ComparableProperty[T comparable](name string, value T, constraints ...ComparableConstraint[T]) ComparableArgument[T] {
	return ComparableArgument[T]{value: &value, constraints: constraints, options: []Option{PropertyName(name)}}
}

// NilComparable argument is used to validate nillable generic comparable value.
func NilComparable[T comparable](value *T, constraints ...ComparableConstraint[T]) ComparableArgument[T] {
	return ComparableArgument[T]{value: value, constraints: constraints}
}

// NilComparableProperty argument is an alias for NilComparable that automatically adds property name to the current scope.
func NilComparableProperty[T comparable](name string, value *T, constraints ...ComparableConstraint[T]) ComparableArgument[T] {
	return ComparableArgument[T]{value: value, constraints: constraints, options: []Option{PropertyName(name)}}
}

// Comparables argument is used to validate generic comparable types.
func Comparables[T comparable](values []T, constraints ...ComparablesConstraint[T]) ComparablesArgument[T] {
	return ComparablesArgument[T]{values: values, constraints: constraints}
}

// ComparablesProperty argument is an alias for Comparables that automatically adds property name to the current scope.
func ComparablesProperty[T comparable](name string, values []T, constraints ...ComparablesConstraint[T]) ComparablesArgument[T] {
	return ComparablesArgument[T]{values: values, constraints: constraints, options: []Option{PropertyName(name)}}
}

// CheckNoViolations is a special argument that checks err for violations. If err contains Violation or ViolationList
// then these violations will be appended into returned violation list from the validator. If err contains an error
// that does not implement an error interface, then the validation process will be terminated and
// this error will be returned.
func CheckNoViolations(err error) Argument {
	return argumentFunc(func(arguments *executionContext) error {
		arguments.addValidator(func(scope Scope) (*ViolationList, error) {
			violations := NewViolationList()
			fatal := violations.AppendFromError(err)
			if fatal != nil {
				return nil, fatal
			}

			return violations, nil
		})

		return nil
	})
}

// Check argument can be useful for quickly checking the result of some simple expression
// that returns a boolean value.
func Check(isValid bool) Checker {
	return Checker{
		isValid:         isValid,
		code:            code.NotValid,
		messageTemplate: message.Templates[code.NotValid],
	}
}

// CheckProperty argument is an alias for Check that automatically adds property name to the current scope.
// It is useful to apply a simple checks on structs.
func CheckProperty(name string, isValid bool) Checker {
	return Checker{
		propertyName:    name,
		isValid:         isValid,
		code:            code.NotValid,
		messageTemplate: message.Templates[code.NotValid],
	}
}

// Language argument sets the current language for translation of a violation message.
func Language(tag language.Tag) Argument {
	return argumentFunc(func(arguments *executionContext) error {
		arguments.scope.language = tag

		return nil
	})
}

// NewArgument can be used to implement your own validation arguments for the specific types.
// See example for more details.
func NewArgument(options []Option, validate ValidateByConstraintFunc) Argument {
	return argumentFunc(func(arguments *executionContext) error {
		arguments.addValidator(newValidator(options, validate))

		return nil
	})
}

// Checker is an argument that can be useful for quickly checking the result of
// some simple expression that returns a boolean value.
type Checker struct {
	isIgnored         bool
	isValid           bool
	propertyName      string
	groups            []string
	code              string
	messageTemplate   string
	messageParameters TemplateParameterList
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

// Code overrides default code for produced violation.
func (c Checker) Code(code string) Checker {
	c.code = code
	return c
}

// Message sets the violation message template. You can set custom template parameters
// for injecting its values into the final message.
func (c Checker) Message(template string, parameters ...TemplateParameter) Checker {
	c.messageTemplate = template
	c.messageParameters = parameters
	return c
}

func (c Checker) setUp(arguments *executionContext) error {
	arguments.addValidator(c.validate)
	return nil
}

func (c Checker) validate(scope Scope) (*ViolationList, error) {
	if c.isValid || c.isIgnored || scope.IsIgnored(c.groups...) {
		return nil, nil
	}
	if c.propertyName != "" {
		scope = scope.AtProperty(c.propertyName)
	}

	violation := scope.BuildViolation(c.code, c.messageTemplate).
		SetParameters(c.messageParameters...).
		CreateViolation()

	return NewViolationList(violation), nil
}
