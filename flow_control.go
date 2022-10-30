package validation

import (
	"context"
	"sync"
)

// WhenArgument is used to build conditional validation. Use the [When] function to initiate a conditional check.
// If the condition is true, then the arguments passed through the [WhenArgument.Then] function will be processed.
// Otherwise, the arguments passed through the [WhenArgument.Else] function will be processed.
type WhenArgument struct {
	isTrue        bool
	path          []PropertyPathElement
	thenArguments []Argument
	elseArguments []Argument
}

// When function is used to initiate conditional validation.
// If the condition is true, then the arguments passed through the [WhenArgument.Then] function will be processed.
// Otherwise, the arguments passed through the [WhenArgument.Else] function will be processed.
func When(isTrue bool) WhenArgument {
	return WhenArgument{isTrue: isTrue}
}

// Then function is used to set a sequence of arguments to be processed if the condition is true.
func (arg WhenArgument) Then(arguments ...Argument) WhenArgument {
	arg.thenArguments = arguments
	return arg
}

// Else function is used to set a sequence of arguments to be processed if a condition is false.
func (arg WhenArgument) Else(arguments ...Argument) WhenArgument {
	arg.elseArguments = arguments
	return arg
}

// At returns a copy of [WhenArgument] with appended property path suffix.
func (arg WhenArgument) At(path ...PropertyPathElement) WhenArgument {
	arg.path = append(arg.path, path...)
	return arg
}

func (arg WhenArgument) setUp(ctx *executionContext) {
	ctx.addValidation(arg.validate, arg.path...)
}

func (arg WhenArgument) validate(ctx context.Context, validator *Validator) (*ViolationList, error) {
	var err error
	if arg.isTrue {
		err = validator.Validate(ctx, arg.thenArguments...)
	} else {
		err = validator.Validate(ctx, arg.elseArguments...)
	}

	return unwrapViolationList(err)
}

// WhenGroupsArgument is used to build conditional validation based on groups. Use the [WhenGroups] function
// to initiate a conditional check. If validation group matches to the validator one,
// then the arguments passed through the [WhenGroupsArgument.Then] function will be processed.
// Otherwise, the arguments passed through the [WhenGroupsArgument.Else] function will be processed.
type WhenGroupsArgument struct {
	groups        []string
	path          []PropertyPathElement
	thenArguments []Argument
	elseArguments []Argument
}

// WhenGroups is used to build conditional validation based on groups. If validation group matches to
// the validator one, then the arguments passed through the [WhenGroupsArgument.Then] function will be processed.
// Otherwise, the arguments passed through the [WhenGroupsArgument.Else] function will be processed.
func WhenGroups(groups ...string) WhenGroupsArgument {
	return WhenGroupsArgument{groups: groups}
}

// Then function is used to set a sequence of arguments to be processed if the validation group is active.
func (arg WhenGroupsArgument) Then(arguments ...Argument) WhenGroupsArgument {
	arg.thenArguments = arguments
	return arg
}

// Else function is used to set a sequence of arguments to be processed if the validation group is active.
func (arg WhenGroupsArgument) Else(arguments ...Argument) WhenGroupsArgument {
	arg.elseArguments = arguments
	return arg
}

// At returns a copy of [WhenGroupsArgument] with appended property path suffix.
func (arg WhenGroupsArgument) At(path ...PropertyPathElement) WhenGroupsArgument {
	arg.path = append(arg.path, path...)
	return arg
}

func (arg WhenGroupsArgument) setUp(ctx *executionContext) {
	ctx.addValidation(arg.validate, arg.path...)
}

func (arg WhenGroupsArgument) validate(ctx context.Context, validator *Validator) (*ViolationList, error) {
	var err error
	if validator.IsIgnoredForGroups(arg.groups...) {
		err = validator.Validate(ctx, arg.elseArguments...)
	} else {
		err = validator.Validate(ctx, arg.thenArguments...)
	}

	return unwrapViolationList(err)
}

// SequentialArgument can be used to interrupt validation process when the first violation is raised.
type SequentialArgument struct {
	isIgnored bool
	path      []PropertyPathElement
	arguments []Argument
}

// Sequentially function used to run validation process step-by-step.
func Sequentially(arguments ...Argument) SequentialArgument {
	return SequentialArgument{arguments: arguments}
}

// At returns a copy of [SequentialArgument] with appended property path suffix.
func (arg SequentialArgument) At(path ...PropertyPathElement) SequentialArgument {
	arg.path = append(arg.path, path...)
	return arg
}

// When enables conditional validation of this argument. If the expression evaluates to false,
// then the argument will be ignored.
func (arg SequentialArgument) When(condition bool) SequentialArgument {
	arg.isIgnored = !condition
	return arg
}

func (arg SequentialArgument) setUp(ctx *executionContext) {
	ctx.addValidation(arg.validate, arg.path...)
}

func (arg SequentialArgument) validate(ctx context.Context, validator *Validator) (*ViolationList, error) {
	if arg.isIgnored {
		return nil, nil
	}

	violations := &ViolationList{}

	for _, argument := range arg.arguments {
		err := violations.AppendFromError(validator.Validate(ctx, argument))
		if err != nil {
			return nil, err
		}
		if violations.len > 0 {
			return violations, nil
		}
	}

	return violations, nil
}

// AtLeastOneOfArgument can be used to set up validation process to check that the value satisfies
// at least one of the given constraints. The validation stops as soon as one constraint is satisfied.
type AtLeastOneOfArgument struct {
	isIgnored bool
	path      []PropertyPathElement
	arguments []Argument
}

// AtLeastOneOf can be used to set up validation process to check that the value satisfies
// at least one of the given constraints. The validation stops as soon as one constraint is satisfied.
func AtLeastOneOf(arguments ...Argument) AtLeastOneOfArgument {
	return AtLeastOneOfArgument{arguments: arguments}
}

// At returns a copy of [AtLeastOneOfArgument] with appended property path suffix.
func (arg AtLeastOneOfArgument) At(path ...PropertyPathElement) AtLeastOneOfArgument {
	arg.path = append(arg.path, path...)
	return arg
}

// When enables conditional validation of this argument. If the expression evaluates to false,
// then the argument will be ignored.
func (arg AtLeastOneOfArgument) When(condition bool) AtLeastOneOfArgument {
	arg.isIgnored = !condition
	return arg
}

func (arg AtLeastOneOfArgument) setUp(ctx *executionContext) {
	ctx.addValidation(arg.validate, arg.path...)
}

func (arg AtLeastOneOfArgument) validate(ctx context.Context, validator *Validator) (*ViolationList, error) {
	if arg.isIgnored {
		return nil, nil
	}

	violations := &ViolationList{}

	for _, argument := range arg.arguments {
		violation := validator.Validate(ctx, argument)
		if violation == nil {
			return nil, nil
		}

		err := violations.AppendFromError(violation)
		if err != nil {
			return nil, err
		}
	}

	return violations, nil
}

// AllArgument can be used to interrupt validation process when the first violation is raised.
type AllArgument struct {
	isIgnored bool
	path      []PropertyPathElement
	arguments []Argument
}

// All runs validation for each argument. It works exactly as [Validator.Validate] method.
// It can be helpful to build complex validation process.
func All(arguments ...Argument) AllArgument {
	return AllArgument{arguments: arguments}
}

// At returns a copy of [AllArgument] with appended property path suffix.
func (arg AllArgument) At(path ...PropertyPathElement) AllArgument {
	arg.path = append(arg.path, path...)
	return arg
}

// When enables conditional validation of this argument. If the expression evaluates to false,
// then the argument will be ignored.
func (arg AllArgument) When(condition bool) AllArgument {
	arg.isIgnored = !condition
	return arg
}

func (arg AllArgument) setUp(ctx *executionContext) {
	ctx.addValidation(arg.validate, arg.path...)
}

func (arg AllArgument) validate(ctx context.Context, validator *Validator) (*ViolationList, error) {
	if arg.isIgnored {
		return nil, nil
	}

	violations := &ViolationList{}

	for _, argument := range arg.arguments {
		err := violations.AppendFromError(validator.Validate(ctx, argument))
		if err != nil {
			return nil, err
		}
	}

	return violations, nil
}

// AsyncArgument can be used to interrupt validation process when the first violation is raised.
type AsyncArgument struct {
	isIgnored bool
	path      []PropertyPathElement
	arguments []Argument
}

// Async implements async/await pattern and runs validation for each argument in a separate goroutine.
func Async(arguments ...Argument) AsyncArgument {
	return AsyncArgument{arguments: arguments}
}

// At returns a copy of [AsyncArgument] with appended property path suffix.
func (arg AsyncArgument) At(path ...PropertyPathElement) AsyncArgument {
	arg.path = append(arg.path, path...)
	return arg
}

// When enables conditional validation of this argument. If the expression evaluates to false,
// then the argument will be ignored.
func (arg AsyncArgument) When(condition bool) AsyncArgument {
	arg.isIgnored = !condition
	return arg
}

func (arg AsyncArgument) setUp(ctx *executionContext) {
	ctx.addValidation(arg.validate, arg.path...)
}

func (arg AsyncArgument) validate(ctx context.Context, validator *Validator) (*ViolationList, error) {
	if arg.isIgnored {
		return nil, nil
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	waiter := &sync.WaitGroup{}
	waiter.Add(len(arg.arguments))
	errs := make(chan error)
	for _, argument := range arg.arguments {
		go func(argument Argument) {
			defer waiter.Done()
			errs <- validator.Validate(ctx, argument)
		}(argument)
	}

	go func() {
		waiter.Wait()
		close(errs)
	}()

	violations := &ViolationList{}

	for violation := range errs {
		err := violations.AppendFromError(violation)
		if err != nil {
			return nil, err
		}
	}

	return violations, nil
}
