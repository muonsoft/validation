package validation

import (
	"reflect"
)

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

func OverrideDefaults(options ...ValidatorOption) {
	for _, setOption := range options {
		setOption(&validator.defaultOptions)
	}
}

var validator = NewValidator()

type validateByConstraintFunc func(constraint Constraint, options Options) error

func (validator *Validator) Validate(value interface{}, options ...Option) error {
	v := reflect.ValueOf(value)

	switch v.Kind() {
	case reflect.Ptr:
		return validator.validatePointer(v, options)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i := v.Int()
		return validator.ValidateInt(&i, options...)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u := v.Uint()
		return validator.ValidateUint(&u, options...)
	case reflect.Float32, reflect.Float64:
		f := v.Float()
		return validator.ValidateFloat(&f, options...)
	case reflect.String:
		s := v.String()
		return validator.ValidateString(&s, options...)
	}

	return &ErrNotValidatable{Value: v}
}

func (validator *Validator) validatePointer(v reflect.Value, options []Option) error {
	p := v.Elem()
	if v.IsNil() {
		return validator.validateNil(options...)
	}

	switch p.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i := p.Int()
		return validator.ValidateInt(&i, options...)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u := p.Uint()
		return validator.ValidateUint(&u, options...)
	case reflect.Float32, reflect.Float64:
		f := p.Float()
		return validator.ValidateFloat(&f, options...)
	case reflect.String:
		s := p.String()
		return validator.ValidateString(&s, options...)
	}

	return &ErrNotValidatable{Value: v}
}

func (validator *Validator) ValidateInt(value *int64, options ...Option) error {
	return validator.executeValidation(options, func(constraint Constraint, options Options) (err error) {
		if constraintValidator, ok := constraint.(IntConstraint); ok {
			err = constraintValidator.ValidateInt(value, options)
		} else {
			err = &ErrInapplicableConstraint{Code: constraint.GetCode(), Type: "int"}
		}

		return err
	})
}

func (validator *Validator) ValidateUint(value *uint64, options ...Option) error {
	return validator.executeValidation(options, func(constraint Constraint, options Options) (err error) {
		if constraintValidator, ok := constraint.(UintConstraint); ok {
			err = constraintValidator.ValidateUint(value, options)
		} else {
			err = &ErrInapplicableConstraint{Code: constraint.GetCode(), Type: "uint"}
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

func (validator *Validator) validateNil(options ...Option) error {
	return validator.executeValidation(options, func(constraint Constraint, options Options) error {
		if constraintValidator, ok := constraint.(NilConstraint); ok {
			return constraintValidator.ValidateNil(options)
		}

		return nil
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
