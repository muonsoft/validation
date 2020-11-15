package validation

type Validatable interface {
	Validate(options ...Option) error
}

func Validate(value interface{}, options ...Option) error {
	return validator.Validate(value, options...)
}

func ValidateString(value *string, options ...Option) error {
	return validator.ValidateString(value, options...)
}

func ValidateInt(value *int64, options ...Option) error {
	return validator.ValidateInt(value, options...)
}

func ValidateUint(value *uint64, options ...Option) error {
	return validator.ValidateUint(value, options...)
}

func ValidateFloat(value *float64, options ...Option) error {
	return validator.ValidateFloat(value, options...)
}

func WithOptions(options ...Option) (*Validator, error) {
	return validator.WithOptions(options...)
}

func Filter(violations ...error) error {
	filteredViolations := make(ViolationList, 0, len(violations))

	for _, err := range violations {
		addErr := filteredViolations.AddFromError(err)
		if addErr != nil {
			return addErr
		}
	}

	return filteredViolations.AsError()
}
