package validation

import "context"

type Options struct {
	Context      context.Context
	PropertyPath PropertyPath
	Constraints  []Constraint
	NewViolation NewViolationFunc
}

func (o Options) NewConstraintViolation(c Constraint) Violation {
	return o.NewViolation(c.GetCode(), c.GetMessageTemplate(), c.GetParameters(), o.PropertyPath)
}

type Option interface {
	Set(options *Options) error
}

type OptionFunc func(options *Options) error

func (f OptionFunc) Set(options *Options) error {
	return f(options)
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
		for _, setOption := range passedOptions {
			if _, isConstraint := setOption.(Constraint); isConstraint {
				continue
			}

			err := setOption.Set(options)
			if err != nil {
				return err
			}
		}

		return nil
	})
}
