package validation

import (
	"reflect"
	"time"

	"github.com/muonsoft/validation/generic"
)

// ValidateByConstraintFunc is used for building validation functions for the values of specific types.
type ValidateByConstraintFunc func(constraint Constraint, scope Scope) error

type validateFunc func(scope Scope) (ViolationList, error)

func newValueValidator(value interface{}, options []Option) (validateFunc, error) {
	switch v := value.(type) {
	case Validatable:
		return newValidValidator(v, options), nil
	case time.Time:
		return newTimeValidator(&v, options), nil
	case *time.Time:
		return newTimeValidator(v, options), nil
	}

	v := reflect.ValueOf(value)

	switch v.Kind() {
	case reflect.Ptr:
		return newValuePointerValidator(v, options)
	case reflect.Bool:
		b := v.Bool()
		return newBoolValidator(&b, options), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		n, err := generic.NewNumber(value)
		if err != nil {
			return nil, err
		}

		return newNumberValidator(*n, options), nil
	case reflect.String:
		s := v.String()
		return newStringValidator(&s, options), nil
	case reflect.Array, reflect.Slice, reflect.Map:
		i, err := generic.NewIterable(value)
		if err != nil {
			return nil, err
		}

		return newIterableValidator(i, options), nil
	}

	return nil, &NotValidatableError{Value: v}
}

func newValuePointerValidator(value reflect.Value, options []Option) (validateFunc, error) {
	p := value.Elem()
	if value.IsNil() {
		return newNilValidator(options), nil
	}

	switch p.Kind() {
	case reflect.Bool:
		b := p.Bool()
		return newBoolValidator(&b, options), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		n, err := generic.NewNumber(p.Interface())
		if err != nil {
			return nil, err
		}

		return newNumberValidator(*n, options), nil
	case reflect.String:
		s := p.String()
		return newStringValidator(&s, options), nil
	case reflect.Array, reflect.Slice, reflect.Map:
		i, err := generic.NewIterable(p.Interface())
		if err != nil {
			return nil, err
		}

		return newIterableValidator(i, options), nil
	}

	return nil, &NotValidatableError{Value: value}
}

func newNilValidator(options []Option) validateFunc {
	return newValidator(options, func(constraint Constraint, scope Scope) error {
		if constraintValidator, ok := constraint.(NilConstraint); ok {
			return constraintValidator.ValidateNil(scope)
		}

		return nil
	})
}

func newBoolValidator(value *bool, options []Option) validateFunc {
	return newValidator(options, func(constraint Constraint, scope Scope) error {
		if c, ok := constraint.(BoolConstraint); ok {
			return c.ValidateBool(value, scope)
		}

		return NewInapplicableConstraintError(constraint, "bool")
	})
}

func newNumberValidator(value generic.Number, options []Option) validateFunc {
	return newValidator(options, func(constraint Constraint, scope Scope) error {
		if c, ok := constraint.(NumberConstraint); ok {
			return c.ValidateNumber(value, scope)
		}

		return NewInapplicableConstraintError(constraint, "number")
	})
}

func newStringValidator(value *string, options []Option) validateFunc {
	return newValidator(options, func(constraint Constraint, scope Scope) error {
		if c, ok := constraint.(StringConstraint); ok {
			return c.ValidateString(value, scope)
		}

		return NewInapplicableConstraintError(constraint, "string")
	})
}

func newIterableValidator(iterable generic.Iterable, options []Option) validateFunc {
	return func(scope Scope) (ViolationList, error) {
		err := scope.applyOptions(options...)
		if err != nil {
			return nil, err
		}

		violations, err := validateOnScope(scope, options, func(constraint Constraint, scope Scope) error {
			if c, ok := constraint.(IterableConstraint); ok {
				return c.ValidateIterable(iterable, scope)
			}

			return NewInapplicableConstraintError(constraint, "iterable")
		})
		if err != nil {
			return nil, err
		}

		if iterable.IsElementImplements(validatableType) {
			vs, err := validateIterableOfValidatables(scope, iterable)
			if err != nil {
				return nil, err
			}
			violations = append(violations, vs...)
		}

		return violations, nil
	}
}

func newCountableValidator(count int, options []Option) validateFunc {
	return newValidator(options, func(constraint Constraint, scope Scope) error {
		if c, ok := constraint.(CountableConstraint); ok {
			return c.ValidateCountable(count, scope)
		}

		return NewInapplicableConstraintError(constraint, "countable")
	})
}

func newTimeValidator(value *time.Time, options []Option) validateFunc {
	return newValidator(options, func(constraint Constraint, scope Scope) error {
		if c, ok := constraint.(TimeConstraint); ok {
			return c.ValidateTime(value, scope)
		}

		return NewInapplicableConstraintError(constraint, "time")
	})
}

func newEachValidator(iterable generic.Iterable, options []Option) validateFunc {
	return func(scope Scope) (ViolationList, error) {
		violations := make(ViolationList, 0)

		err := iterable.Iterate(func(key generic.Key, value interface{}) error {
			opts := options
			if key.IsIndex() {
				opts = append(opts, ArrayIndex(key.Index()))
			} else {
				opts = append(opts, PropertyName(key.String()))
			}

			validate, err := newValueValidator(value, opts)
			if err != nil {
				return err
			}

			vs, err := validate(scope)
			if err != nil {
				return err
			}
			violations = append(violations, vs...)

			return nil
		})

		return violations, err
	}
}

func newEachStringValidator(values []string, options []Option) validateFunc {
	return func(scope Scope) (ViolationList, error) {
		violations := make(ViolationList, 0)

		for i := range values {
			opts := append(options, ArrayIndex(i))
			validate := newStringValidator(&values[i], opts)
			vs, err := validate(scope)
			if err != nil {
				return nil, err
			}
			violations = append(violations, vs...)
		}

		return violations, nil
	}
}

func newValidValidator(value Validatable, options []Option) validateFunc {
	return func(scope Scope) (ViolationList, error) {
		err := scope.applyOptions(options...)
		if err != nil {
			return nil, err
		}

		err = value.Validate(newScopedValidator(scope))
		violations, ok := UnwrapViolationList(err)
		if ok {
			return violations, nil
		}

		return nil, err
	}
}

func newValidator(options []Option, validate ValidateByConstraintFunc) validateFunc {
	return func(scope Scope) (ViolationList, error) {
		err := scope.applyOptions(options...)
		if err != nil {
			return nil, err
		}

		return validateOnScope(scope, options, validate)
	}
}

func validateOnScope(scope Scope, options []Option, validate ValidateByConstraintFunc) (ViolationList, error) {
	violations := make(ViolationList, 0)

	for _, option := range options {
		if constraint, ok := option.(ConditionalConstraint); ok {
			err := validateConditionConstraints(scope, constraint, &violations, validate)
			if err != nil {
				return nil, err
			}
		} else if constraint, ok := option.(Constraint); ok {
			err := violations.AppendFromError(validate(constraint, scope))
			if err != nil {
				return nil, err
			}
		}
	}

	return violations, nil
}

func validateConditionConstraints(
	scope Scope,
	constraint ConditionalConstraint,
	violations *ViolationList,
	validate ValidateByConstraintFunc,
) error {
	if constraint.condition {
		for _, constraint := range constraint.thenConstraints {
			err := violations.AppendFromError(validate(constraint, scope))
			if err != nil {
				return err
			}
		}
	} else {
		for _, constraint := range constraint.elseConstraints {
			err := violations.AppendFromError(validate(constraint, scope))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func validateIterableOfValidatables(scope Scope, iterable generic.Iterable) (ViolationList, error) {
	violations := make(ViolationList, 0)

	err := iterable.Iterate(func(key generic.Key, value interface{}) error {
		s := scope
		if key.IsIndex() {
			s = s.atIndex(key.Index())
		} else {
			s = s.atProperty(key.String())
		}

		validate := newValidValidator(value.(Validatable), nil)
		vs, err := validate(s)
		if err != nil {
			return err
		}

		violations = append(violations, vs...)

		return nil
	})

	return violations, err
}
