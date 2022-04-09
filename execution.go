package validation

import (
	"time"
)

type executionContext struct {
	scope      Scope
	validators []validateFunc
}

func (ctx *executionContext) addValidator(validator validateFunc) {
	ctx.validators = append(ctx.validators, validator)
}

type BoolArgument struct {
	value       *bool
	constraints []BoolConstraint
	options     []Option
}

func (arg BoolArgument) With(options ...Option) BoolArgument {
	arg.options = append(arg.options, options...)
	return arg
}

func (arg BoolArgument) setUp(ctx *executionContext) error {
	ctx.addValidator(func(scope Scope) (*ViolationList, error) {
		violations := NewViolationList()

		err := scope.applyOptions(arg.options...)
		if err != nil {
			return nil, err
		}

		for i := range arg.constraints {
			err := violations.AppendFromError(arg.constraints[i].ValidateBool(arg.value, scope))
			if err != nil {
				return nil, err
			}
		}

		return violations, nil
	})

	return nil
}

type NumberArgument[T Numeric] struct {
	value       *T
	constraints []NumberConstraint[T]
	options     []Option
}

func (arg NumberArgument[T]) With(options ...Option) NumberArgument[T] {
	arg.options = append(arg.options, options...)
	return arg
}

func (arg NumberArgument[T]) setUp(ctx *executionContext) error {
	ctx.addValidator(func(scope Scope) (*ViolationList, error) {
		violations := NewViolationList()

		err := scope.applyOptions(arg.options...)
		if err != nil {
			return nil, err
		}

		for i := range arg.constraints {
			err := violations.AppendFromError(arg.constraints[i].ValidateNumber(arg.value, scope))
			if err != nil {
				return nil, err
			}
		}

		return violations, nil
	})

	return nil
}

type StringArgument struct {
	value       *string
	constraints []StringConstraint
	options     []Option
}

func (arg StringArgument) With(options ...Option) StringArgument {
	arg.options = append(arg.options, options...)
	return arg
}

func (arg StringArgument) setUp(ctx *executionContext) error {
	ctx.addValidator(func(scope Scope) (*ViolationList, error) {
		violations := NewViolationList()

		err := scope.applyOptions(arg.options...)
		if err != nil {
			return nil, err
		}

		for i := range arg.constraints {
			err := violations.AppendFromError(arg.constraints[i].ValidateString(arg.value, scope))
			if err != nil {
				return nil, err
			}
		}

		return violations, nil
	})

	return nil
}

type CountableArgument struct {
	count       int
	constraints []CountableConstraint
	options     []Option
}

func (arg CountableArgument) With(options ...Option) CountableArgument {
	arg.options = append(arg.options, options...)
	return arg
}

func (arg CountableArgument) setUp(ctx *executionContext) error {
	ctx.addValidator(func(scope Scope) (*ViolationList, error) {
		violations := NewViolationList()

		err := scope.applyOptions(arg.options...)
		if err != nil {
			return nil, err
		}

		for i := range arg.constraints {
			err := violations.AppendFromError(arg.constraints[i].ValidateCountable(arg.count, scope))
			if err != nil {
				return nil, err
			}
		}

		return violations, nil
	})

	return nil
}

type TimeArgument struct {
	value       *time.Time
	constraints []TimeConstraint
	options     []Option
}

func (arg TimeArgument) With(options ...Option) TimeArgument {
	arg.options = append(arg.options, options...)
	return arg
}

func (arg TimeArgument) setUp(ctx *executionContext) error {
	ctx.addValidator(func(scope Scope) (*ViolationList, error) {
		violations := NewViolationList()

		err := scope.applyOptions(arg.options...)
		if err != nil {
			return nil, err
		}

		for i := range arg.constraints {
			err := violations.AppendFromError(arg.constraints[i].ValidateTime(arg.value, scope))
			if err != nil {
				return nil, err
			}
		}

		return violations, nil
	})

	return nil
}

type EachStringArgument struct {
	values      []string
	constraints []StringConstraint
	options     []Option
}

func (arg EachStringArgument) With(options ...Option) EachStringArgument {
	arg.options = append(arg.options, options...)
	return arg
}

func (arg EachStringArgument) setUp(ctx *executionContext) error {
	ctx.addValidator(func(scope Scope) (*ViolationList, error) {
		violations := NewViolationList()

		err := scope.applyOptions(arg.options...)
		if err != nil {
			return nil, err
		}

		for i := range arg.values {
			for j := range arg.constraints {
				err := violations.AppendFromError(arg.constraints[j].ValidateString(&arg.values[i], scope.AtIndex(i)))
				if err != nil {
					return nil, err
				}
			}
		}

		return violations, nil
	})

	return nil
}

type EachNumberArgument[T Numeric] struct {
	values      []T
	constraints []NumberConstraint[T]
	options     []Option
}

func (arg EachNumberArgument[T]) With(options ...Option) EachNumberArgument[T] {
	arg.options = append(arg.options, options...)
	return arg
}

func (arg EachNumberArgument[T]) setUp(ctx *executionContext) error {
	ctx.addValidator(func(scope Scope) (*ViolationList, error) {
		violations := NewViolationList()

		err := scope.applyOptions(arg.options...)
		if err != nil {
			return nil, err
		}

		for i := range arg.values {
			for j := range arg.constraints {
				err := violations.AppendFromError(arg.constraints[j].ValidateNumber(&arg.values[i], scope.AtIndex(i)))
				if err != nil {
					return nil, err
				}
			}
		}

		return violations, nil
	})

	return nil
}

type ValidArgument struct {
	value   Validatable
	options []Option
}

func (arg ValidArgument) With(options ...Option) ValidArgument {
	arg.options = append(arg.options, options...)
	return arg
}

func (arg ValidArgument) setUp(ctx *executionContext) error {
	ctx.addValidator(func(scope Scope) (*ViolationList, error) {
		err := scope.applyOptions(arg.options...)
		if err != nil {
			return nil, err
		}

		err = arg.value.Validate(scope.context, scope.Validator())
		violations, ok := UnwrapViolationList(err)
		if ok {
			return violations, nil
		}

		return nil, err
	})

	return nil
}

type ValidSliceArgument[T Validatable] struct {
	values  []T
	options []Option
}

func (arg ValidSliceArgument[T]) With(options ...Option) ValidSliceArgument[T] {
	arg.options = append(arg.options, options...)
	return arg
}

func (arg ValidSliceArgument[T]) setUp(ctx *executionContext) error {
	ctx.addValidator(func(scope Scope) (*ViolationList, error) {
		violations := NewViolationList()

		err := scope.applyOptions(arg.options...)
		if err != nil {
			return nil, err
		}

		for i, value := range arg.values {
			s := scope.AtIndex(i)
			err := violations.AppendFromError(value.Validate(s.context, s.Validator()))
			if err != nil {
				return nil, err
			}
		}

		return violations, nil
	})

	return nil
}

type ValidMapArgument[T Validatable] struct {
	values  map[string]T
	options []Option
}

func (arg ValidMapArgument[T]) With(options ...Option) ValidMapArgument[T] {
	arg.options = append(arg.options, options...)
	return arg
}

func (arg ValidMapArgument[T]) setUp(ctx *executionContext) error {
	ctx.addValidator(func(scope Scope) (*ViolationList, error) {
		violations := NewViolationList()

		err := scope.applyOptions(arg.options...)
		if err != nil {
			return nil, err
		}

		for key, value := range arg.values {
			s := scope.AtProperty(key)
			err := violations.AppendFromError(value.Validate(s.context, s.Validator()))
			if err != nil {
				return nil, err
			}
		}

		return violations, nil
	})

	return nil
}

type ComparableArgument[T comparable] struct {
	value       *T
	constraints []ComparableConstraint[T]
	options     []Option
}

func (arg ComparableArgument[T]) With(options ...Option) ComparableArgument[T] {
	arg.options = append(arg.options, options...)
	return arg
}

func (arg ComparableArgument[T]) setUp(ctx *executionContext) error {
	ctx.addValidator(func(scope Scope) (*ViolationList, error) {
		violations := NewViolationList()

		err := scope.applyOptions(arg.options...)
		if err != nil {
			return nil, err
		}

		for i := range arg.constraints {
			err := violations.AppendFromError(arg.constraints[i].ValidateComparable(arg.value, scope))
			if err != nil {
				return nil, err
			}
		}

		return violations, nil
	})

	return nil
}

type ComparablesArgument[T comparable] struct {
	values      []T
	constraints []ComparablesConstraint[T]
	options     []Option
}

func (arg ComparablesArgument[T]) With(options ...Option) ComparablesArgument[T] {
	arg.options = append(arg.options, options...)
	return arg
}

func (arg ComparablesArgument[T]) setUp(ctx *executionContext) error {
	ctx.addValidator(func(scope Scope) (*ViolationList, error) {
		violations := NewViolationList()

		err := scope.applyOptions(arg.options...)
		if err != nil {
			return nil, err
		}

		for i := range arg.constraints {
			err := violations.AppendFromError(arg.constraints[i].ValidateComparables(arg.values, scope))
			if err != nil {
				return nil, err
			}
		}

		return violations, nil
	})

	return nil
}
