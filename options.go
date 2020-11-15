package validation

import "context"

type Option interface {
	Set(options *Options) error
}

type OptionFunc func(options *Options) error

func (f OptionFunc) Set(options *Options) error {
	return f(options)
}

type Options struct {
	Context      context.Context
	PropertyPath PropertyPath
	Constraints  []Constraint
	NewViolation NewViolationFunc
}

func (o Options) NewConstraintViolation(c Constraint) Violation {
	return o.NewViolation(c.GetCode(), c.GetMessageTemplate(), c.GetParameters(), o.PropertyPath)
}

func (o *Options) apply(options ...Option) error {
	for _, option := range options {
		err := option.Set(o)
		if err != nil {
			return err
		}
	}

	return nil
}

func (o *Options) applyNonConstraints(options ...Option) error {
	for _, option := range options {
		if _, isConstraint := option.(Constraint); isConstraint {
			continue
		}

		err := option.Set(o)
		if err != nil {
			return err
		}
	}

	return nil
}

func Context(ctx context.Context) Option {
	return OptionFunc(func(options *Options) error {
		options.Context = ctx

		return nil
	})
}

func PropertyName(propertyName string) Option {
	return OptionFunc(func(options *Options) error {
		options.PropertyPath = append(options.PropertyPath, PropertyNameElement{propertyName})

		return nil
	})
}

func ArrayIndex(index int) Option {
	return OptionFunc(func(options *Options) error {
		options.PropertyPath = append(options.PropertyPath, ArrayIndexElement{index})

		return nil
	})
}

func PassOptions(passedOptions []Option) Option {
	return OptionFunc(func(options *Options) error {
		return options.applyNonConstraints(passedOptions...)
	})
}

func extendAndPassOptions(extendedOptions *Options, passedOptions ...Option) Option {
	return OptionFunc(func(options *Options) error {
		options.Context = extendedOptions.Context
		options.PropertyPath = append(options.PropertyPath, extendedOptions.PropertyPath...)
		options.NewViolation = extendedOptions.NewViolation

		return options.applyNonConstraints(passedOptions...)
	})
}
