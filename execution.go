package validation

import (
	"time"
)

type executionContext struct {
	scope      Scope
	validators []ValidateOnScopeFunc
}

func (ctx *executionContext) addValidator(options []Option, validator ValidateOnScopeFunc) {
	ctx.validators = append(ctx.validators, func(scope Scope) (*ViolationList, error) {
		err := scope.applyOptions(options...)
		if err != nil {
			return nil, err
		}

		return validator(scope)
	})
}

func validateNil(isNil bool, constraints []NilConstraint) ValidateOnScopeFunc {
	return func(scope Scope) (*ViolationList, error) {
		violations := NewViolationList()

		for i := range constraints {
			err := violations.AppendFromError(constraints[i].ValidateNil(isNil, scope))
			if err != nil {
				return nil, err
			}
		}

		return violations, nil
	}
}

func validateBool(value *bool, constraints []BoolConstraint) ValidateOnScopeFunc {
	return func(scope Scope) (*ViolationList, error) {
		violations := NewViolationList()

		for i := range constraints {
			err := violations.AppendFromError(constraints[i].ValidateBool(value, scope))
			if err != nil {
				return nil, err
			}
		}

		return violations, nil
	}
}

func validateNumber[T Numeric](value *T, constraints []NumberConstraint[T]) ValidateOnScopeFunc {
	return func(scope Scope) (*ViolationList, error) {
		violations := NewViolationList()

		for i := range constraints {
			err := violations.AppendFromError(constraints[i].ValidateNumber(value, scope))
			if err != nil {
				return nil, err
			}
		}

		return violations, nil
	}
}

func validateString(value *string, constraints []StringConstraint) ValidateOnScopeFunc {
	return func(scope Scope) (*ViolationList, error) {
		violations := NewViolationList()

		for i := range constraints {
			err := violations.AppendFromError(constraints[i].ValidateString(value, scope))
			if err != nil {
				return nil, err
			}
		}

		return violations, nil
	}
}

func validateCountable(count int, constraints []CountableConstraint) ValidateOnScopeFunc {
	return func(scope Scope) (*ViolationList, error) {
		violations := NewViolationList()

		for i := range constraints {
			err := violations.AppendFromError(constraints[i].ValidateCountable(count, scope))
			if err != nil {
				return nil, err
			}
		}

		return violations, nil
	}
}

func validateTime(value *time.Time, constraints []TimeConstraint) ValidateOnScopeFunc {
	return func(scope Scope) (*ViolationList, error) {
		violations := NewViolationList()

		for i := range constraints {
			err := violations.AppendFromError(constraints[i].ValidateTime(value, scope))
			if err != nil {
				return nil, err
			}
		}

		return violations, nil
	}
}

func validateEachString(values []string, constraints []StringConstraint) ValidateOnScopeFunc {
	return func(scope Scope) (*ViolationList, error) {
		violations := NewViolationList()

		for i := range values {
			for j := range constraints {
				err := violations.AppendFromError(constraints[j].ValidateString(&values[i], scope.AtIndex(i)))
				if err != nil {
					return nil, err
				}
			}
		}

		return violations, nil
	}
}

func validateEachNumber[T Numeric](values []T, constraints []NumberConstraint[T]) ValidateOnScopeFunc {
	return func(scope Scope) (*ViolationList, error) {
		violations := NewViolationList()

		for i := range values {
			for j := range constraints {
				err := violations.AppendFromError(constraints[j].ValidateNumber(&values[i], scope.AtIndex(i)))
				if err != nil {
					return nil, err
				}
			}
		}

		return violations, nil
	}
}
func validateIt(value Validatable) ValidateOnScopeFunc {
	return func(scope Scope) (*ViolationList, error) {
		err := value.Validate(scope.context, scope.Validator())
		violations, ok := UnwrapViolationList(err)
		if ok {
			return violations, nil
		}

		return nil, err
	}
}

func validateSlice[T Validatable](values []T) ValidateOnScopeFunc {
	return func(scope Scope) (*ViolationList, error) {
		violations := NewViolationList()

		for i, value := range values {
			s := scope.AtIndex(i)
			err := violations.AppendFromError(value.Validate(s.context, s.Validator()))
			if err != nil {
				return nil, err
			}
		}

		return violations, nil
	}
}

func validateMap[T Validatable](values map[string]T) ValidateOnScopeFunc {
	return func(scope Scope) (*ViolationList, error) {
		violations := NewViolationList()

		for key, value := range values {
			s := scope.AtProperty(key)
			err := violations.AppendFromError(value.Validate(s.context, s.Validator()))
			if err != nil {
				return nil, err
			}
		}

		return violations, nil
	}
}

func validateComparable[T comparable](value *T, constraints []ComparableConstraint[T]) ValidateOnScopeFunc {
	return func(scope Scope) (*ViolationList, error) {
		violations := NewViolationList()

		for i := range constraints {
			err := violations.AppendFromError(constraints[i].ValidateComparable(value, scope))
			if err != nil {
				return nil, err
			}
		}

		return violations, nil
	}
}

func validateComparables[T comparable](values []T, constraints []ComparablesConstraint[T]) ValidateOnScopeFunc {
	return func(scope Scope) (*ViolationList, error) {
		violations := NewViolationList()

		for i := range constraints {
			err := violations.AppendFromError(constraints[i].ValidateComparables(values, scope))
			if err != nil {
				return nil, err
			}
		}

		return violations, nil
	}
}
