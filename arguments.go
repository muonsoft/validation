package validation

import (
	"fmt"
	"time"

	"github.com/muonsoft/validation/generic"
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

func Bool(value *bool, options ...Option) Argument {
	return ArgumentFunc(func(arguments *Arguments) error {
		arguments.addValidator(newBoolValidator(value, options))

		return nil
	})
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

func String(value *string, options ...Option) Argument {
	return ArgumentFunc(func(arguments *Arguments) error {
		arguments.addValidator(newStringValidator(value, options))

		return nil
	})
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

func Countable(count int, options ...Option) Argument {
	return ArgumentFunc(func(arguments *Arguments) error {
		arguments.addValidator(newCountableValidator(count, options))

		return nil
	})
}

func Time(value *time.Time, options ...Option) Argument {
	return ArgumentFunc(func(arguments *Arguments) error {
		arguments.addValidator(newTimeValidator(value, options))

		return nil
	})
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

func EachString(values []string, options ...Option) Argument {
	return ArgumentFunc(func(arguments *Arguments) error {
		arguments.addValidator(newEachStringValidator(values, options))

		return nil
	})
}

func Valid(value Validatable, options ...Option) Argument {
	return ArgumentFunc(func(arguments *Arguments) error {
		arguments.addValidator(newValidValidator(value, options))

		return nil
	})
}
