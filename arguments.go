package validation

import (
	"context"
	"fmt"
	"time"

	"github.com/muonsoft/validation/generic"
	"golang.org/x/text/language"
)

type Argument interface {
	set(arguments *Arguments) error
}

type ArgumentFunc func(arguments *Arguments) error

func (f ArgumentFunc) set(arguments *Arguments) error {
	return f(arguments)
}

type validateFunc func(scope Scope) (ViolationList, error)

type Arguments struct {
	scope      Scope
	validators []validateFunc
}

func (args *Arguments) addValidator(validator validateFunc) {
	args.validators = append(args.validators, validator)
}

func Value(value interface{}, options ...Option) Argument {
	return ArgumentFunc(func(arguments *Arguments) error {
		v, err := newValueValidator(value, options)
		if err != nil {
			return err
		}

		arguments.addValidator(v)

		return nil
	})
}

func PropertyValue(name string, value interface{}, options ...Option) Argument {
	return Value(value, append([]Option{PropertyName(name)}, options...)...)
}

func Bool(value *bool, options ...Option) Argument {
	return ArgumentFunc(func(arguments *Arguments) error {
		arguments.addValidator(newBoolValidator(value, options))

		return nil
	})
}

func BoolProperty(name string, value *bool, options ...Option) Argument {
	return Bool(value, append([]Option{PropertyName(name)}, options...)...)
}

func Number(value interface{}, options ...Option) Argument {
	return ArgumentFunc(func(arguments *Arguments) error {
		number, err := generic.NewNumber(value)
		if err != nil {
			return fmt.Errorf("cannot convert value '%v' to number: %w", value, err)
		}

		arguments.addValidator(newNumberValidator(*number, options))

		return nil
	})
}

func NumberProperty(name string, value interface{}, options ...Option) Argument {
	return Number(value, append([]Option{PropertyName(name)}, options...)...)
}

func String(value *string, options ...Option) Argument {
	return ArgumentFunc(func(arguments *Arguments) error {
		arguments.addValidator(newStringValidator(value, options))

		return nil
	})
}

func StringProperty(name string, value *string, options ...Option) Argument {
	return String(value, append([]Option{PropertyName(name)}, options...)...)
}

func Iterable(value interface{}, options ...Option) Argument {
	return ArgumentFunc(func(arguments *Arguments) error {
		iterable, err := generic.NewIterable(value)
		if err != nil {
			return fmt.Errorf("cannot convert value '%v' to iterable: %w", value, err)
		}

		arguments.addValidator(newIterableValidator(iterable, options))

		return nil
	})
}

func IterableProperty(name string, value interface{}, options ...Option) Argument {
	return Iterable(value, append([]Option{PropertyName(name)}, options...)...)
}

func Countable(count int, options ...Option) Argument {
	return ArgumentFunc(func(arguments *Arguments) error {
		arguments.addValidator(newCountableValidator(count, options))

		return nil
	})
}

func CountableProperty(name string, count int, options ...Option) Argument {
	return Countable(count, append([]Option{PropertyName(name)}, options...)...)
}

func Time(value *time.Time, options ...Option) Argument {
	return ArgumentFunc(func(arguments *Arguments) error {
		arguments.addValidator(newTimeValidator(value, options))

		return nil
	})
}

func TimeProperty(name string, value *time.Time, options ...Option) Argument {
	return Time(value, append([]Option{PropertyName(name)}, options...)...)
}

func Each(value interface{}, options ...Option) Argument {
	return ArgumentFunc(func(arguments *Arguments) error {
		iterable, err := generic.NewIterable(value)
		if err != nil {
			return fmt.Errorf("cannot convert value '%v' to iterable: %w", value, err)
		}

		arguments.addValidator(newEachValidator(iterable, options))

		return nil
	})
}

func EachProperty(name string, value interface{}, options ...Option) Argument {
	return Each(value, append([]Option{PropertyName(name)}, options...)...)
}

func EachString(values []string, options ...Option) Argument {
	return ArgumentFunc(func(arguments *Arguments) error {
		arguments.addValidator(newEachStringValidator(values, options))

		return nil
	})
}

func EachStringProperty(name string, values []string, options ...Option) Argument {
	return EachString(values, append([]Option{PropertyName(name)}, options...)...)
}

func Valid(value Validatable, options ...Option) Argument {
	return ArgumentFunc(func(arguments *Arguments) error {
		arguments.addValidator(newValidValidator(value, options))

		return nil
	})
}

func ValidProperty(name string, value Validatable, options ...Option) Argument {
	return Valid(value, append([]Option{PropertyName(name)}, options...)...)
}

// Context can be used to pass context to validation constraints via Scope.
func Context(ctx context.Context) Argument {
	return ArgumentFunc(func(arguments *Arguments) error {
		arguments.scope.context = ctx

		return nil
	})
}

// Language argument sets the current language for translation of violation message.
func Language(tag language.Tag) Argument {
	return ArgumentFunc(func(arguments *Arguments) error {
		arguments.scope.language = tag

		return nil
	})
}
