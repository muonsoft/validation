package validation

func Validate(violations ...error) error {
	return validator.Validate(violations...)
}

func ValidateString(value *string, options ...Option) error {
	return validator.ValidateString(value, options...)
}

func ValidateInt(value *int, options ...Option) error {
	return validator.ValidateInt(value, options...)
}

func ValidateFloat(value *float64, options ...Option) error {
	return validator.ValidateFloat(value, options...)
}
