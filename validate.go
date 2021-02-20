package validation

type Validatable interface {
	Validate(options ...Option) error
}

func Validate(value interface{}, options ...Option) error {
	return validator.Validate(value, options...)
}

func ValidateBool(value *bool, options ...Option) error {
	return validator.ValidateBool(value, options...)
}

func ValidateNumber(value interface{}, options ...Option) error {
	return validator.ValidateNumber(value, options...)
}

func ValidateString(value *string, options ...Option) error {
	return validator.ValidateString(value, options...)
}

func ValidateIterable(value interface{}, options ...Option) error {
	return validator.ValidateIterable(value, options...)
}

func ValidateCountable(count int, options ...Option) error {
	return validator.ValidateCountable(count, options...)
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
