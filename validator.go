package validation

type Validator struct {
}

func NewValidator() *Validator {
	return &Validator{}
}

var validator = NewValidator()

func (validator *Validator) Validate(violations ...error) error {
	filteredViolations := make(ViolationList, 0, len(violations))

	for _, err := range violations {
		if violation, ok := UnwrapViolation(err); ok {
			filteredViolations = append(filteredViolations, violation)
		} else if violationList, ok := UnwrapViolationList(err); ok {
			filteredViolations = append(filteredViolations, violationList...)
		} else if err != nil {
			return err
		}
	}

	if len(filteredViolations) == 0 {
		return nil
	}

	return filteredViolations
}

func (validator *Validator) ValidateString(value *string, options ...Option) error {
	opts, err := collectOptions(options)
	if err != nil {
		return err
	}

	for _, constraint := range opts.Constraints {
		if constraintValidator, ok := constraint.(StringConstraint); ok {
			err := constraintValidator.ValidateString(value, *opts)
			if err != nil {
				return err
			}
		} else {
			return &ErrInapplicableConstraint{Code: constraint.Code(), Type: "string"}
		}
	}

	return nil
}

func (validator *Validator) ValidateInt(value *int, options ...Option) error {
	opts, err := collectOptions(options)
	if err != nil {
		return err
	}

	for _, constraint := range opts.Constraints {
		if constraintValidator, ok := constraint.(IntConstraint); ok {
			err := constraintValidator.ValidateInt(value, *opts)
			if err != nil {
				return err
			}
		} else {
			return &ErrInapplicableConstraint{Code: constraint.Code(), Type: "int"}
		}
	}

	return nil
}

func (validator *Validator) ValidateFloat(value *float64, options ...Option) error {
	opts, err := collectOptions(options)
	if err != nil {
		return err
	}

	for _, constraint := range opts.Constraints {
		if constraintValidator, ok := constraint.(FloatConstraint); ok {
			err := constraintValidator.ValidateFloat(value, *opts)
			if err != nil {
				return err
			}
		} else {
			return &ErrInapplicableConstraint{Code: constraint.Code(), Type: "float"}
		}
	}

	return nil
}
