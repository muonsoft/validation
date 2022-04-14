package validation

// Option is used to set up validation process of a value.
type Option interface {
	// SetUp is called when the validation process is initialized and can be used to modify scope values.
	SetUp(scope *Scope) error
}

// optionFunc is an adapter that allows to use ordinary functions as validation options.
type optionFunc func(scope *Scope) error

func (f optionFunc) SetUp(scope *Scope) error {
	return f(scope)
}

// PropertyName option adds name of the given property to current validation scope.
func PropertyName(propertyName string) Option {
	return optionFunc(func(scope *Scope) error {
		scope.propertyPath = scope.propertyPath.WithProperty(propertyName)

		return nil
	})
}

// ArrayIndex option adds index of the given array to current validation scope.
func ArrayIndex(index int) Option {
	return optionFunc(func(scope *Scope) error {
		scope.propertyPath = scope.propertyPath.WithIndex(index)

		return nil
	})
}
