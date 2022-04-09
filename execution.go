package validation

type NumberArgument[T Numeric] struct {
	value       *T
	constraints []ComparableConstraint[T]
	options     []Option
}

func (arg NumberArgument[T]) WithOptions(options ...Option) NumberArgument[T] {
	arg.options = options
	return arg
}

func (arg NumberArgument[T]) setUp(ctx *executionContext) error {
	ctx.addValidator(func(scope Scope) (*ViolationList, error) {
		violations := NewViolationList()

		// todo: apply options

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
