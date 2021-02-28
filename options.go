package validation

import (
	"context"

	"golang.org/x/text/language"
)

type Option interface {
	Set(scope *Scope) error
}

type OptionFunc func(scope *Scope) error

func (f OptionFunc) Set(scope *Scope) error {
	return f(scope)
}

func Context(ctx context.Context) Option {
	return OptionFunc(func(scope *Scope) error {
		scope.context = ctx

		return nil
	})
}

func PropertyName(propertyName string) Option {
	return OptionFunc(func(scope *Scope) error {
		scope.propertyPath = append(scope.propertyPath, PropertyNameElement{propertyName})

		return nil
	})
}

func ArrayIndex(index int) Option {
	return OptionFunc(func(scope *Scope) error {
		scope.propertyPath = append(scope.propertyPath, ArrayIndexElement{index})

		return nil
	})
}

// Language option sets current language for translation of violation message.
func Language(tag language.Tag) Option {
	return OptionFunc(func(scope *Scope) error {
		scope.language = tag

		return nil
	})
}

func PassOptions(passedOptions []Option) Option {
	return OptionFunc(func(scope *Scope) error {
		return scope.applyNonConstraints(passedOptions...)
	})
}
