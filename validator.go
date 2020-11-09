package validation

type Validator struct {
	defaultOptions Options
}

type ValidatorOption func(options *Options)

func NewValidator(options ...ValidatorOption) *Validator {
	defaultOptions := Options{
		NewViolation: NewViolation,
	}

	for _, setOption := range options {
		setOption(&defaultOptions)
	}

	return &Validator{defaultOptions: defaultOptions}
}

func OverrideNewViolation(violationFunc NewViolationFunc) ValidatorOption {
	return func(options *Options) {
		options.NewViolation = violationFunc
	}
}

var validator = NewValidator()

type validateByConstraintFunc func(constraint Constraint, options Options) error

func (validator *Validator) Validate(violations ...error) error {
	filteredViolations := make(ViolationList, 0, len(violations))

	for _, err := range violations {
		addErr := filteredViolations.AddFromError(err)
		if addErr != nil {
			return addErr
		}
	}

	return filteredViolations.AsError()
}

func (validator *Validator) ValidateString(value *string, options ...Option) error {
	return validator.executeValidation(options, func(constraint Constraint, options Options) (err error) {
		if constraintValidator, ok := constraint.(StringConstraint); ok {
			err = constraintValidator.ValidateString(value, options)
		} else {
			err = &ErrInapplicableConstraint{Code: constraint.GetCode(), Type: "string"}
		}

		return err
	})
}

func (validator *Validator) ValidateInt(value *int, options ...Option) error {
	return validator.executeValidation(options, func(constraint Constraint, options Options) (err error) {
		if constraintValidator, ok := constraint.(IntConstraint); ok {
			err = constraintValidator.ValidateInt(value, options)
		} else {
			err = &ErrInapplicableConstraint{Code: constraint.GetCode(), Type: "int"}
		}

		return err
	})
}

func (validator *Validator) ValidateFloat(value *float64, options ...Option) error {
	return validator.executeValidation(options, func(constraint Constraint, options Options) (err error) {
		if constraintValidator, ok := constraint.(FloatConstraint); ok {
			err = constraintValidator.ValidateFloat(value, options)
		} else {
			err = &ErrInapplicableConstraint{Code: constraint.GetCode(), Type: "float"}
		}

		return err
	})
}

func (validator *Validator) executeValidation(options []Option, validate validateByConstraintFunc) error {
	opts, err := validator.createOptionsFromDefaults(options)
	if err != nil {
		return err
	}

	violations := make(ViolationList, 0, len(opts.Constraints))

	for _, constraint := range opts.Constraints {
		err := violations.AddFromError(validate(constraint, *opts))
		if err != nil {
			return err
		}
	}

	return violations.AsError()
}

func (validator *Validator) createOptionsFromDefaults(options []Option) (*Options, error) {
	opts := validator.defaultOptions

	for _, option := range options {
		err := option.Set(&opts)
		if err != nil {
			return nil, err
		}
	}

	return &opts, nil
}
