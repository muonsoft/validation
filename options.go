package validation

// Option is used to set up validation process of a value.
type Option interface {
	// SetUp commonly used to tune validation Scope. Also, it can be used to gracefully handle errors
	// while initializing constraints.
	SetUp(scope *Scope) error
}

// OptionFunc is an adapter to allow use the use of ordinary functions as validation options.
type OptionFunc func(scope *Scope) error

func (f OptionFunc) SetUp(scope *Scope) error {
	return f(scope)
}

func PropertyName(propertyName string) Option {
	return OptionFunc(func(scope *Scope) error {
		scope.propertyPath = append(scope.propertyPath, PropertyNameElement(propertyName))

		return nil
	})
}

func ArrayIndex(index int) Option {
	return OptionFunc(func(scope *Scope) error {
		scope.propertyPath = append(scope.propertyPath, ArrayIndexElement(index))

		return nil
	})
}
