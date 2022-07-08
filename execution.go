package validation

import (
	"context"
	"time"
)

type executionContext struct {
	validations []ValidateFunc
}

func (ctx *executionContext) addValidation(options []Option, validate ValidateFunc) {
	ctx.validations = append(ctx.validations, func(ctx context.Context, validator *Validator) (*ViolationList, error) {
		return validate(ctx, validator.withOptions(options...))
	})
}

func validateNil(isNil bool, constraints []NilConstraint) ValidateFunc {
	return func(ctx context.Context, validator *Validator) (*ViolationList, error) {
		violations := NewViolationList()

		for i := range constraints {
			err := violations.AppendFromError(constraints[i].ValidateNil(ctx, validator, isNil))
			if err != nil {
				return nil, err
			}
		}

		return violations, nil
	}
}

func validateBool(value *bool, constraints []BoolConstraint) ValidateFunc {
	return func(ctx context.Context, validator *Validator) (*ViolationList, error) {
		violations := NewViolationList()

		for i := range constraints {
			err := violations.AppendFromError(constraints[i].ValidateBool(ctx, validator, value))
			if err != nil {
				return nil, err
			}
		}

		return violations, nil
	}
}

func validateNumber[T Numeric](value *T, constraints []NumberConstraint[T]) ValidateFunc {
	return func(ctx context.Context, validator *Validator) (*ViolationList, error) {
		violations := NewViolationList()

		for i := range constraints {
			err := violations.AppendFromError(constraints[i].ValidateNumber(ctx, validator, value))
			if err != nil {
				return nil, err
			}
		}

		return violations, nil
	}
}

func validateString(value *string, constraints []StringConstraint) ValidateFunc {
	return func(ctx context.Context, validator *Validator) (*ViolationList, error) {
		violations := NewViolationList()

		for i := range constraints {
			err := violations.AppendFromError(constraints[i].ValidateString(ctx, validator, value))
			if err != nil {
				return nil, err
			}
		}

		return violations, nil
	}
}

func validateCountable(count int, constraints []CountableConstraint) ValidateFunc {
	return func(ctx context.Context, validator *Validator) (*ViolationList, error) {
		violations := NewViolationList()

		for i := range constraints {
			err := violations.AppendFromError(constraints[i].ValidateCountable(ctx, validator, count))
			if err != nil {
				return nil, err
			}
		}

		return violations, nil
	}
}

func validateTime(value *time.Time, constraints []TimeConstraint) ValidateFunc {
	return func(ctx context.Context, validator *Validator) (*ViolationList, error) {
		violations := NewViolationList()

		for i := range constraints {
			err := violations.AppendFromError(constraints[i].ValidateTime(ctx, validator, value))
			if err != nil {
				return nil, err
			}
		}

		return violations, nil
	}
}

func validateEachString(values []string, constraints []StringConstraint) ValidateFunc {
	return func(ctx context.Context, validator *Validator) (*ViolationList, error) {
		violations := NewViolationList()

		for i := range values {
			for _, constraint := range constraints {
				err := violations.AppendFromError(constraint.ValidateString(ctx, validator.AtIndex(i), &values[i]))
				if err != nil {
					return nil, err
				}
			}
		}

		return violations, nil
	}
}

func validateEachNumber[T Numeric](values []T, constraints []NumberConstraint[T]) ValidateFunc {
	return func(ctx context.Context, validator *Validator) (*ViolationList, error) {
		violations := NewViolationList()

		for i := range values {
			for _, constraint := range constraints {
				err := violations.AppendFromError(constraint.ValidateNumber(ctx, validator.AtIndex(i), &values[i]))
				if err != nil {
					return nil, err
				}
			}
		}

		return violations, nil
	}
}

func validateEachComparable[T comparable](values []T, constraints []ComparableConstraint[T]) ValidateFunc {
	return func(ctx context.Context, validator *Validator) (*ViolationList, error) {
		violations := NewViolationList()

		for i := range values {
			for _, constraint := range constraints {
				err := violations.AppendFromError(constraint.ValidateComparable(ctx, validator.AtIndex(i), &values[i]))
				if err != nil {
					return nil, err
				}
			}
		}

		return violations, nil
	}
}

func validateIt(value Validatable) ValidateFunc {
	return func(ctx context.Context, validator *Validator) (*ViolationList, error) {
		err := value.Validate(ctx, validator)
		violations, ok := UnwrapViolationList(err)
		if ok {
			return violations, nil
		}

		return nil, err
	}
}

func validateSlice[T Validatable](values []T) ValidateFunc {
	return func(ctx context.Context, validator *Validator) (*ViolationList, error) {
		violations := NewViolationList()

		for i, value := range values {
			err := violations.AppendFromError(value.Validate(ctx, validator.AtIndex(i)))
			if err != nil {
				return nil, err
			}
		}

		return violations, nil
	}
}

func validateMap[T Validatable](values map[string]T) ValidateFunc {
	return func(ctx context.Context, validator *Validator) (*ViolationList, error) {
		violations := NewViolationList()

		for key, value := range values {
			err := violations.AppendFromError(value.Validate(ctx, validator.AtProperty(key)))
			if err != nil {
				return nil, err
			}
		}

		return violations, nil
	}
}

func validateComparable[T comparable](value *T, constraints []ComparableConstraint[T]) ValidateFunc {
	return func(ctx context.Context, validator *Validator) (*ViolationList, error) {
		violations := NewViolationList()

		for i := range constraints {
			err := violations.AppendFromError(constraints[i].ValidateComparable(ctx, validator, value))
			if err != nil {
				return nil, err
			}
		}

		return violations, nil
	}
}

func validateComparables[T comparable](values []T, constraints []ComparablesConstraint[T]) ValidateFunc {
	return func(ctx context.Context, validator *Validator) (*ViolationList, error) {
		violations := NewViolationList()

		for i := range constraints {
			err := violations.AppendFromError(constraints[i].ValidateComparables(ctx, validator, values))
			if err != nil {
				return nil, err
			}
		}

		return violations, nil
	}
}
