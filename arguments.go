package validation

import (
	"context"
	"fmt"
	"time"

	"github.com/muonsoft/validation/generic"
	"golang.org/x/text/language"
)

// Argument used to set up the validation process. It is used to set up the current validation scope and to pass
// arguments for validation values.
type Argument interface {
	set(arguments *Arguments) error
}

type argumentFunc func(arguments *Arguments) error

func (f argumentFunc) set(arguments *Arguments) error {
	return f(arguments)
}

type Arguments struct {
	scope      Scope
	validators []validateFunc
}

func (args *Arguments) addValidator(validator validateFunc) {
	args.validators = append(args.validators, validator)
}

// Value argument is used to validate any supported value. It uses reflection to detect the type of the argument
// and pass it to a specific validation method.
//
// If the validator cannot determine the value or it is not supported, then NotValidatableError will be returned
// when calling the validator.Validate method.
//
// Example
//  v := 0
//  err := validator.Validate(Value(v, it.IsNotBlank()))
func Value(value interface{}, options ...Option) Argument {
	return argumentFunc(func(arguments *Arguments) error {
		v, err := newValueValidator(value, options)
		if err != nil {
			return err
		}

		arguments.addValidator(v)

		return nil
	})
}

// PropertyValue argument is an alias for Value that automatically adds property name to the current scope.
func PropertyValue(name string, value interface{}, options ...Option) Argument {
	return Value(value, append([]Option{PropertyName(name)}, options...)...)
}

// Bool argument is used to validate boolean values.
//
// Example
//  v := false
//  err := validator.Validate(Bool(&v, it.IsNotBlank()))
func Bool(value *bool, options ...Option) Argument {
	return argumentFunc(func(arguments *Arguments) error {
		arguments.addValidator(newBoolValidator(value, options))

		return nil
	})
}

// BoolProperty argument is an alias for Bool that automatically adds property name to the current scope.
func BoolProperty(name string, value *bool, options ...Option) Argument {
	return Bool(value, append([]Option{PropertyName(name)}, options...)...)
}

// Number argument is used to validate numbers (any types of integers or floats). At the moment it uses
// reflection to detect numeric value. Given value is internally converted into int64 or float64 to make comparisons.
//
// Warning! This method will be changed after generics implementation in Go.
//
// Example
//  v := 0
//  err := validator.Validate(Number(&v, it.IsNotBlank()))
func Number(value interface{}, options ...Option) Argument {
	return argumentFunc(func(arguments *Arguments) error {
		number, err := generic.NewNumber(value)
		if err != nil {
			return fmt.Errorf("cannot convert value '%v' to number: %w", value, err)
		}

		arguments.addValidator(newNumberValidator(*number, options))

		return nil
	})
}

// NumberProperty argument is an alias for Number that automatically adds property name to the current scope.
func NumberProperty(name string, value interface{}, options ...Option) Argument {
	return Number(value, append([]Option{PropertyName(name)}, options...)...)
}

// String argument is used to validate strings.
//
// Example
//  v := ""
//  err := validator.Validate(String(&v, it.IsNotBlank()))
func String(value *string, options ...Option) Argument {
	return argumentFunc(func(arguments *Arguments) error {
		arguments.addValidator(newStringValidator(value, options))

		return nil
	})
}

// StringProperty argument is an alias for String that automatically adds property name to the current scope.
func StringProperty(name string, value *string, options ...Option) Argument {
	return String(value, append([]Option{PropertyName(name)}, options...)...)
}

// Iterable argument is used to validate arrays, slices, or maps. At the moment it uses reflection
// to iterate over values. So you can expect a performance hit using this method. For better performance
// it is recommended to make a custom type that implements the Validatable interface.
//
// Warning! This argument is subject to change in the final versions of the library.
//
// Example
//  v := []string{}
//  err := validator.Validate(Iterable(&v, it.IsNotBlank()))
func Iterable(value interface{}, options ...Option) Argument {
	return argumentFunc(func(arguments *Arguments) error {
		iterable, err := generic.NewIterable(value)
		if err != nil {
			return fmt.Errorf("cannot convert value '%v' to iterable: %w", value, err)
		}

		arguments.addValidator(newIterableValidator(iterable, options))

		return nil
	})
}

// IterableProperty argument is an alias for Iterable that automatically adds property name to the current scope.
func IterableProperty(name string, value interface{}, options ...Option) Argument {
	return Iterable(value, append([]Option{PropertyName(name)}, options...)...)
}

// Countable argument can be used to validate size of an array, slice, or map. You can pass result of len()
// function as an argument.
//
// Example
//  s := []string{"a", "b"}
//  err := validator.Validate(Countable(len(s), it.HasMinCount(3)))
func Countable(count int, options ...Option) Argument {
	return argumentFunc(func(arguments *Arguments) error {
		arguments.addValidator(newCountableValidator(count, options))

		return nil
	})
}

// CountableProperty argument is an alias for Countable that automatically adds property name to the current scope.
func CountableProperty(name string, count int, options ...Option) Argument {
	return Countable(count, append([]Option{PropertyName(name)}, options...)...)
}

// Time argument is used to validate time.Time value.
//
// Example
//  t := time.Now()
//  err := validator.Validate(Time(&t, it.IsNotBlank()))
func Time(value *time.Time, options ...Option) Argument {
	return argumentFunc(func(arguments *Arguments) error {
		arguments.addValidator(newTimeValidator(value, options))

		return nil
	})
}

// TimeProperty argument is an alias for Time that automatically adds property name to the current scope.
func TimeProperty(name string, value *time.Time, options ...Option) Argument {
	return Time(value, append([]Option{PropertyName(name)}, options...)...)
}

// Each is used to validate each value of iterable (array, slice, or map). At the moment it uses reflection
// to iterate over values. So you can expect a performance hit using this method. For better performance
// it is recommended to make a custom type that implements the Validatable interface. Also, you can use
// EachString argument to validate slice of strings.
//
// Warning! This argument is subject to change in the final versions of the library.
//
// Example
//  v := []string{""}
//  err := validator.Validate(Each(&v, it.IsNotBlank()))
func Each(value interface{}, options ...Option) Argument {
	return argumentFunc(func(arguments *Arguments) error {
		iterable, err := generic.NewIterable(value)
		if err != nil {
			return fmt.Errorf("cannot convert value '%v' to iterable: %w", value, err)
		}

		arguments.addValidator(newEachValidator(iterable, options))

		return nil
	})
}

// EachProperty argument is an alias for Each that automatically adds property name to the current scope.
func EachProperty(name string, value interface{}, options ...Option) Argument {
	return Each(value, append([]Option{PropertyName(name)}, options...)...)
}

// EachString is used to validate a slice of strings. This is a more performant version of Each argument.
//
// Example
//  v := []string{""}
//  err := validator.Validate(EachString(&v, it.IsNotBlank()))
func EachString(values []string, options ...Option) Argument {
	return argumentFunc(func(arguments *Arguments) error {
		arguments.addValidator(newEachStringValidator(values, options))

		return nil
	})
}

// EachStringProperty argument is an alias for EachString that automatically adds property name to the current scope.
func EachStringProperty(name string, values []string, options ...Option) Argument {
	return EachString(values, append([]Option{PropertyName(name)}, options...)...)
}

// Valid is used to run validation on the Validatable type. This method is recommended
// to run a complex validation process.
func Valid(value Validatable, options ...Option) Argument {
	return argumentFunc(func(arguments *Arguments) error {
		arguments.addValidator(newValidValidator(value, options))

		return nil
	})
}

// ValidProperty argument is an alias for Valid that automatically adds property name to the current scope.
func ValidProperty(name string, value Validatable, options ...Option) Argument {
	return Valid(value, append([]Option{PropertyName(name)}, options...)...)
}

// Context can be used to pass context to validation constraints via scope.
//
// Example
//  err := validator.Validate(
//      Context(request.Context()),
//      String(&s, it.IsNotBlank()), // now all called constraints will use passed context in their methods
//  )
func Context(ctx context.Context) Argument {
	return argumentFunc(func(arguments *Arguments) error {
		arguments.scope.context = ctx

		return nil
	})
}

// Language argument sets the current language for translation of a violation message.
//
// Example
//  err := validator.Validate(
//      Language(language.Russian),
//      String(&s, it.IsNotBlank()), // all violations created in scope will be translated into Russian
//  )
func Language(tag language.Tag) Argument {
	return argumentFunc(func(arguments *Arguments) error {
		arguments.scope.language = tag

		return nil
	})
}
