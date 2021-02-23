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
	defaultOptions := newDefaultOptions()

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

func ResetDefaults() {
	validator.defaultOptions = newDefaultOptions()
}

var validator = NewValidator()

type validateByConstraintFunc func(constraint Constraint, options Options) error

func (validator *Validator) Validate(value interface{}, options ...Option) error {
	if validatable, ok := value.(Validatable); ok {
		return validator.ValidateValidatable(validatable, options...)
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
	case reflect.Array, reflect.Slice, reflect.Map:
		return validator.ValidateIterable(value, options...)
	}

	return &NotValidatableError{Value: v}
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
	case reflect.Array, reflect.Slice, reflect.Map:
		return validator.ValidateIterable(p.Interface(), options...)
	}

	return &NotValidatableError{Value: v}
}

func (validator *Validator) ValidateBool(value *bool, options ...Option) error {
	return validator.executeValidationAndHandleError(options, func(constraint Constraint, options Options) (err error) {
		if constraintValidator, ok := constraint.(BoolConstraint); ok {
			err = constraintValidator.ValidateBool(value, options)
		} else {
			err = newInapplicableConstraintError(constraint, "bool")
		}

		return err
	})
}

func (validator *Validator) ValidateNumber(value interface{}, options ...Option) error {
	number, err := generic.NewNumber(value)
	if err != nil {
		return fmt.Errorf("cannot convert value '%v' to number: %w", value, err)
	}

	return validator.executeValidationAndHandleError(options, func(constraint Constraint, options Options) (err error) {
		if constraintValidator, ok := constraint.(NumberConstraint); ok {
			err = constraintValidator.ValidateNumber(*number, options)
		} else {
			err = newInapplicableConstraintError(constraint, "number")
		}

		return err
	})
}

func (validator *Validator) ValidateString(value *string, options ...Option) error {
	return validator.executeValidationAndHandleError(options, func(constraint Constraint, options Options) (err error) {
		if constraintValidator, ok := constraint.(StringConstraint); ok {
			err = constraintValidator.ValidateString(value, options)
		} else {
			err = newInapplicableConstraintError(constraint, "string")
		}

		return err
	})
}

func (validator *Validator) ValidateIterable(value interface{}, options ...Option) error {
	iterable, err := generic.NewIterable(value)
	if err != nil {
		return fmt.Errorf("cannot convert value '%v' to iterable: %w", value, err)
	}

	violations, err := validator.executeValidation(options, func(constraint Constraint, options Options) (err error) {
		if constraintValidator, ok := constraint.(IterableConstraint); ok {
			err = constraintValidator.ValidateIterable(iterable, options)
		} else {
			err = newInapplicableConstraintError(constraint, "iterable")
		}

		return err
	})
	if err != nil {
		return err
	}

	if iterable.IsElementImplements(validatableType) {
		elementViolations, err := validator.validateIterableOfValidatables(iterable, options)
		if err != nil {
			return err
		}
		violations = append(violations, elementViolations...)
	}

	return violations.AsError()
}

func (validator *Validator) ValidateCountable(count int, options ...Option) error {
	return validator.executeValidationAndHandleError(options, func(constraint Constraint, options Options) (err error) {
		if constraintValidator, ok := constraint.(CountableConstraint); ok {
			err = constraintValidator.ValidateCountable(count, options)
		} else {
			err = newInapplicableConstraintError(constraint, "countable")
		}

		return err
	})
}

func (validator *Validator) ValidateValidatable(validatable Validatable, options ...Option) error {
	return validatable.Validate(extendAndPassOptions(&validator.defaultOptions, options...))
}

func (validator *Validator) ValidateEach(value interface{}, options ...Option) error {
	iterable, err := generic.NewIterable(value)
	if err != nil {
		return fmt.Errorf("cannot convert value '%v' to iterable: %w", value, err)
	}

	violations := make(ViolationList, 0)

	err = iterable.Iterate(func(key generic.Key, value interface{}) error {
		opts := options
		if key.IsIndex() {
			opts = append(opts, ArrayIndex(key.Index()))
		} else {
			opts = append(opts, PropertyName(key.String()))
		}

		err := validator.Validate(value, opts...)
		return violations.AddFromError(err)
	})
	if err != nil {
		return err
	}

	return violations.AsError()
}

func (validator *Validator) ValidateEachString(strings []string, options ...Option) error {
	violations := make(ViolationList, 0)

	for i := range strings {
		opts := append(options, ArrayIndex(i))
		err := violations.AddFromError(validator.ValidateString(&strings[i], opts...))
		if err != nil {
			return err
		}
	}

	return violations.AsError()
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
	return validator.executeValidationAndHandleError(options, func(constraint Constraint, options Options) error {
		if constraintValidator, ok := constraint.(NilConstraint); ok {
			return constraintValidator.ValidateNil(options)
		}

		return nil
	})
}

func (validator *Validator) validateIterableOfValidatables(
	iterable generic.Iterable,
	options []Option,
) (ViolationList, error) {
	violations := make(ViolationList, 0)

	err := iterable.Iterate(func(key generic.Key, value interface{}) error {
		opts := options
		if key.IsIndex() {
			opts = append(opts, ArrayIndex(key.Index()))
		} else {
			opts = append(opts, PropertyName(key.String()))
		}

		elementValidator, err := validator.WithOptions(opts...)
		if err != nil {
			return err
		}

		err = elementValidator.ValidateValidatable(value.(Validatable))
		return violations.AddFromError(err)
	})

	return violations, err
}

func (validator *Validator) executeValidationAndHandleError(options []Option, validate validateByConstraintFunc) error {
	violations, err := validator.executeValidation(options, validate)
	if err != nil {
		return err
	}
	return violations.AsError()
}

func (validator *Validator) executeValidation(
	options []Option,
	validate validateByConstraintFunc,
) (ViolationList, error) {
	opts, err := validator.createOptionsFromDefaults(options)
	if err != nil {
		return nil, err
	}

	violations := make(ViolationList, 0, len(opts.Constraints))

	for _, constraint := range opts.Constraints {
		err := violations.AddFromError(validate(constraint, *opts))
		if err != nil {
			return nil, err
		}
	}

	return violations, nil
}

func (validator *Validator) createOptionsFromDefaults(options []Option) (*Options, error) {
	opts := validator.defaultOptions
	err := opts.apply(options...)
	if err != nil {
		return nil, err
	}

	return &opts, nil
}
