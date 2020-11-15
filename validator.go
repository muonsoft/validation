package validation

import (
	"fmt"
	"reflect"

	"github.com/muonsoft/validation/generic"
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
	if validatable, ok := value.(Validatable); ok {
		return validatable.Validate(
			extendAndPassOptions(&validator.defaultOptions, options...),
		)
	}

	v := reflect.ValueOf(value)

	switch v.Kind() {
	case reflect.Ptr:
		return validator.validatePointer(v, options)
	case reflect.Bool:
		b := v.Bool()
		return validator.ValidateBool(&b, options...)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return validator.ValidateNumber(value, options...)
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
	case reflect.Bool:
		b := p.Bool()
		return validator.ValidateBool(&b, options...)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return validator.ValidateNumber(p.Interface(), options...)
	case reflect.String:
		s := p.String()
		return validator.ValidateString(&s, options...)
	}

	return &ErrNotValidatable{Value: v}
}

func (validator *Validator) ValidateBool(value *bool, options ...Option) error {
	return validator.executeValidation(options, func(constraint Constraint, options Options) (err error) {
		if constraintValidator, ok := constraint.(BoolConstraint); ok {
			err = constraintValidator.ValidateBool(value, options)
		} else {
			err = &ErrInapplicableConstraint{Code: constraint.GetCode(), Type: "bool"}
		}

		return err
	})
}

func (validator *Validator) ValidateNumber(value interface{}, options ...Option) error {
	number, err := generic.NewNumber(value)
	if err != nil {
		return fmt.Errorf("cannot convert value '%v' to number: %w", value, err)
	}

	return validator.executeValidation(options, func(constraint Constraint, options Options) (err error) {
		if constraintValidator, ok := constraint.(NumberConstraint); ok {
			err = constraintValidator.ValidateNumber(*number, options)
		} else {
			err = &ErrInapplicableConstraint{Code: constraint.GetCode(), Type: "number"}
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

func (validator *Validator) WithOptions(options ...Option) (*Validator, error) {
	newOptions := validator.defaultOptions
	err := newOptions.applyNonConstraints(options...)
	if err != nil {
		return nil, err
	}

	return &Validator{defaultOptions: newOptions}, nil
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
	err := opts.apply(options...)
	if err != nil {
		return nil, err
	}

	return &opts, nil
}
