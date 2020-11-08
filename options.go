package validation

import "context"

type Options struct {
	Context      context.Context
	PropertyPath PropertyPath
	Constraints  []Constraint
}

type Option interface {
	Set(options *Options) error
}

type OptionFunc func(options *Options) error

func Context(ctx context.Context) Option {
	return optionSetter{option: func(options *Options) error {
		options.Context = ctx

		return nil
	}}
}

func PropertyName(propertyName string) Option {
	return optionSetter{option: func(options *Options) error {
		options.PropertyPath = append(options.PropertyPath, PropertyNameElement{propertyName})

		return nil
	}}
}

func ArrayIndex(index int) Option {
	return optionSetter{option: func(options *Options) error {
		options.PropertyPath = append(options.PropertyPath, ArrayIndexElement{index})

		return nil
	}}
}

type optionSetter struct {
	option OptionFunc
}

func (set optionSetter) Set(options *Options) error {
	return set.option(options)
}

func collectOptions(options []Option) (*Options, error) {
	opts := &Options{}

	for _, option := range options {
		err := option.Set(opts)
		if err != nil {
			return nil, err
		}
	}

	return opts, nil
}
