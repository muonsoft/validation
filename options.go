package validation

// Option is used to set up validation process of a value.
type Option interface {
	// SetUp commonly used to tune validation Scope. Also, it can be used to gracefully handle errors
	// while initializing constraints.
	SetUp(scope *Scope) error
}

// optionFunc is an adapter that allows you to use ordinary functions as validation options.
type optionFunc func(scope *Scope) error

func (f optionFunc) SetUp(scope *Scope) error {
	return f(scope)
}

// PropertyName option adds name of the given property to current validation Scope.
func PropertyName(propertyName string) Option {
	return optionFunc(func(scope *Scope) error {
		scope.propertyPath = append(scope.propertyPath, PropertyNameElement(propertyName))

		return nil
	})
}

// ArrayIndex option adds index of the given array to current validation Scope.
func ArrayIndex(index int) Option {
	return optionFunc(func(scope *Scope) error {
		scope.propertyPath = append(scope.propertyPath, ArrayIndexElement(index))

		return nil
	})
}
