package validation

type Validator struct {
}

func NewValidator() *Validator {
	return &Validator{}
}

var validator = NewValidator()

type validateByConstraintFunc func(constraint Constraint, options Options) error

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
	return validator.executeValidation(options, func(constraint Constraint, options Options) (err error) {
		if constraintValidator, ok := constraint.(StringConstraint); ok {
			err = constraintValidator.ValidateString(value, options)
		} else {
			err = &ErrInapplicableConstraint{Code: constraint.Code(), Type: "string"}
		}

		return err
	})
}

func (validator *Validator) ValidateInt(value *int, options ...Option) error {
	return validator.executeValidation(options, func(constraint Constraint, options Options) (err error) {
		if constraintValidator, ok := constraint.(IntConstraint); ok {
			err = constraintValidator.ValidateInt(value, options)
		} else {
			err = &ErrInapplicableConstraint{Code: constraint.Code(), Type: "int"}
		}

		return err
	})
}

func (validator *Validator) ValidateFloat(value *float64, options ...Option) error {
	return validator.executeValidation(options, func(constraint Constraint, options Options) (err error) {
		if constraintValidator, ok := constraint.(FloatConstraint); ok {
			err = constraintValidator.ValidateFloat(value, options)
		} else {
			err = &ErrInapplicableConstraint{Code: constraint.Code(), Type: "float"}
		}

		return err
	})
}

func (validator *Validator) executeValidation(options []Option, validate validateByConstraintFunc) error {
	opts, err := collectOptions(options)
	if err != nil {
		return err
	}

	violations := make(ViolationList, 0, len(opts.Constraints))

	for _, constraint := range opts.Constraints {
		err := validate(constraint, *opts)

		if violation, ok := UnwrapViolation(err); ok {
			violations = append(violations, violation)
		} else if violationList, ok := UnwrapViolationList(err); ok {
			violations = append(violations, violationList...)
		} else if err != nil {
			return err
		}
	}

	if len(violations) == 0 {
		return nil
	}

	return violations
}
