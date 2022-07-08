package validation

// Option is used to set up validation process of a value.
type Option interface {
	// SetUp is called when the validation process is initialized and can be used to modify validation context.
	SetUp(validator *Validator) *Validator
}

// optionFunc is an adapter that allows to use ordinary functions as validation options.
type optionFunc func(validator *Validator) *Validator

func (f optionFunc) SetUp(validator *Validator) *Validator {
	return f(validator)
}

// PropertyName option adds name of the given property to current validation path.
func PropertyName(propertyName string) Option {
	return optionFunc(func(validator *Validator) *Validator {
		v := validator.copy()
		v.propertyPath = v.propertyPath.WithProperty(propertyName)

		return v
	})
}

// ArrayIndex option adds index of the given array to current validation path.
func ArrayIndex(index int) Option {
	return optionFunc(func(validator *Validator) *Validator {
		v := validator.copy()
		v.propertyPath = v.propertyPath.WithIndex(index)

		return v
	})
}
