package test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/validationtest"
	"github.com/stretchr/testify/assert"
)

var (
	ErrThen   = errors.New("then")
	ErrElse   = errors.New("else")
	ErrFirst  = errors.New("first")
	ErrSecond = errors.New("second")
)

func TestValidatorArgument_WhenConditionIsFalse_ExpectNoErrors(t *testing.T) {
	err := newValidator(t).Validate(
		context.Background(),
		validation.String("", it.IsNotBlank()).When(false),
	)

	assertNoError(t, err)
}

func TestWhenArgument_WhenConditionIsTrue_ExpectAllConstraintsOfThenBranchApplied(t *testing.T) {
	err := newValidator(t).Validate(
		context.Background(),
		validation.When(true).
			Then(validation.String("", it.IsNotBlank().WithError(ErrThen))).
			Else(validation.String("", it.IsNotBlank().WithError(ErrElse))),
	)

	validationtest.Assert(t, err).IsViolationList().WithOneViolation().WithError(ErrThen)
}

func TestWhenArgument_WhenConditionIsFalse_ExpectAllConstraintsOfElseBranchApplied(t *testing.T) {
	err := newValidator(t).Validate(
		context.Background(),
		validation.When(false).
			Then(validation.String("", it.IsNotBlank().WithError(ErrThen))).
			Else(validation.String("", it.IsNotBlank().WithError(ErrElse))),
	)

	validationtest.Assert(t, err).IsViolationList().WithOneViolation().WithError(ErrElse)
}

func TestWhenArgument_WhenConditionIsFalseAndNoElseBranch_ExpectNoViolations(t *testing.T) {
	err := newValidator(t).Validate(
		context.Background(),
		validation.When(false).
			Then(validation.String("", it.IsNotBlank().WithError(ErrThen))),
	)

	assertNoError(t, err)
}

func TestWhenArgument_WhenPathIsSet_ExpectViolationWithPath(t *testing.T) {
	err := newValidator(t).Validate(
		context.Background(),
		validation.When(true).
			With(
				validation.PropertyName("properties"),
				validation.ArrayIndex(0),
				validation.PropertyName("property"),
			).
			Then(validation.String("", it.IsNotBlank().WithError(ErrThen))),
	)

	validationtest.Assert(t, err).IsViolationList().WithOneViolation().
		WithError(ErrThen).
		WithPropertyPath("properties[0].property")
}

func TestWhenGroupsArgument_WhenGroupMatches_ExpectViolation(t *testing.T) {
	err := newValidator(t).WithGroups(testGroup).Validate(
		context.Background(),
		validation.WhenGroups(testGroup).
			Then(validation.String("", it.IsNotBlank().WhenGroups(testGroup).WithError(ErrThen))).
			Else(validation.String("", it.IsNotBlank().WhenGroups(testGroup).WithError(ErrElse))),
	)

	validationtest.Assert(t, err).IsViolationList().WithOneViolation().WithError(ErrThen)
}

func TestWhenGroupsArgument_WhenGroupNotMatches_ExpectNoError(t *testing.T) {
	err := newValidator(t).Validate(
		context.Background(),
		validation.WhenGroups(testGroup).
			Then(validation.String("", it.IsNotBlank().WithError(ErrThen))).
			Else(validation.String("", it.IsNotBlank().WithError(ErrElse))),
	)

	validationtest.Assert(t, err).IsViolationList().WithOneViolation().WithError(ErrElse)
}

func TestWhenGroupsArgument_WhenGroupNotMatchesNoElseBranch_ExpectNoViolations(t *testing.T) {
	err := newValidator(t).Validate(
		context.Background(),
		validation.WhenGroups(testGroup).
			Then(validation.String("", it.IsNotBlank().WithError(ErrThen))),
	)

	assertNoError(t, err)
}

func TestWhenGroupsArgument_WhenPathIsSet_ExpectViolationWithPath(t *testing.T) {
	err := newValidator(t).WithGroups(testGroup).Validate(
		context.Background(),
		validation.WhenGroups(testGroup).
			Then(validation.String("", it.IsNotBlank().WhenGroups(testGroup).WithError(ErrThen))).
			With(
				validation.PropertyName("properties"),
				validation.ArrayIndex(0),
				validation.PropertyName("property"),
			),
	)

	validationtest.Assert(t, err).IsViolationList().WithOneViolation().
		WithError(ErrThen).
		WithPropertyPath("properties[0].property")
}

func TestSequentialArgument_WhenInvalidValueAtFirstConstraint_ExpectOneViolation(t *testing.T) {
	err := newValidator(t).Validate(
		context.Background(),
		validation.Sequentially(
			validation.String("", it.IsNotBlank().WithError(ErrFirst)),
			validation.String("", it.IsNotBlank().WithError(ErrSecond)),
		),
	)

	validationtest.Assert(t, err).IsViolationList().WithOneViolation().WithError(ErrFirst)
}

func TestSequentialArgument_WhenPathIsSet_ExpectOneViolationWithPath(t *testing.T) {
	err := newValidator(t).Validate(
		context.Background(),
		validation.Sequentially(
			validation.String("", it.IsNotBlank().WithError(ErrFirst)),
			validation.String("", it.IsNotBlank().WithError(ErrSecond)),
		).With(
			validation.PropertyName("properties"),
			validation.ArrayIndex(0),
			validation.PropertyName("property"),
		),
	)

	validationtest.Assert(t, err).IsViolationList().WithOneViolation().
		WithError(ErrFirst).
		WithPropertyPath("properties[0].property")
}

func TestSequentialArgument_WhenValidationIsDisabled_ExpectNoErrors(t *testing.T) {
	err := newValidator(t).Validate(
		context.Background(),
		validation.Sequentially(
			validation.String("", it.IsNotBlank().WithError(ErrFirst)),
			validation.String("", it.IsNotBlank().WithError(ErrSecond)),
		).When(false),
	)

	assert.NoError(t, err)
}

func TestAtLeastOneOfArgument_WhenInvalidValueAtFirstConstraint_ExpectAllViolations(t *testing.T) {
	err := newValidator(t).Validate(
		context.Background(),
		validation.AtLeastOneOf(
			validation.String("", it.IsNotBlank().WithError(ErrFirst)),
			validation.String("", it.IsNotBlank().WithError(ErrSecond)),
		),
	)

	validationtest.Assert(t, err).IsViolationList().WithErrors(ErrFirst, ErrSecond)
}

func TestAtLeastOneOfArgument_WhenInvalidValueAtSecondConstraint_ExpectNoViolation(t *testing.T) {
	err := newValidator(t).Validate(
		context.Background(),
		validation.AtLeastOneOf(
			validation.String("", it.IsNotBlank().WithError(ErrFirst)),
			validation.String("foo", it.IsNotBlank().WithError(ErrSecond)),
		),
	)

	assertNoError(t, err)
}

func TestAtLeastOneOfArgument_WhenPathIsSet_ExpectOneViolationWithPath(t *testing.T) {
	err := newValidator(t).Validate(
		context.Background(),
		validation.AtLeastOneOf(
			validation.String("", it.IsNotBlank()),
		).With(
			validation.PropertyName("properties"),
			validation.ArrayIndex(0),
			validation.PropertyName("property"),
		),
	)

	validationtest.Assert(t, err).IsViolationList().WithOneViolation().
		WithError(validation.ErrIsBlank).
		WithPropertyPath("properties[0].property")
}

func TestAtLeastOneOfArgument_WhenValidationIsDisabled_ExpectNoError(t *testing.T) {
	err := newValidator(t).Validate(
		context.Background(),
		validation.AtLeastOneOf(
			validation.String("", it.IsNotBlank().WithError(ErrFirst)),
			validation.String("", it.IsNotBlank().WithError(ErrSecond)),
		).When(false),
	)

	assert.NoError(t, err)
}

func TestAllArgument_WhenInvalidValueAtFirstConstraint_ExpectAllViolations(t *testing.T) {
	err := newValidator(t).Validate(
		context.Background(),
		validation.All(
			validation.String("", it.IsNotBlank().WithError(ErrFirst)),
			validation.String("", it.IsNotBlank().WithError(ErrSecond)),
		),
	)

	validationtest.Assert(t, err).IsViolationList().WithErrors(ErrFirst, ErrSecond)
}

func TestAllArgument_WhenPathIsSet_ExpectOneViolationWithPath(t *testing.T) {
	err := newValidator(t).Validate(
		context.Background(),
		validation.All(
			validation.String("", it.IsNotBlank().WithError(ErrFirst)),
			validation.String("", it.IsNotBlank().WithError(ErrSecond)),
		).With(
			validation.PropertyName("properties"),
			validation.ArrayIndex(0),
			validation.PropertyName("property"),
		),
	)

	violations := validationtest.Assert(t, err).IsViolationList()
	violations.HasViolationAt(0).WithError(ErrFirst).WithPropertyPath("properties[0].property")
	violations.HasViolationAt(1).WithError(ErrSecond).WithPropertyPath("properties[0].property")
}

func TestAllArgument_WhenValidationIsDisabled_ExpectNoErrors(t *testing.T) {
	err := newValidator(t).Validate(
		context.Background(),
		validation.All(
			validation.String("", it.IsNotBlank().WithError(ErrFirst)),
			validation.String("", it.IsNotBlank().WithError(ErrSecond)),
		).When(false),
	)

	assert.NoError(t, err)
}

func TestAsyncArgument_WhenInvalidValueAtFirstConstraint_ExpectAllViolations(t *testing.T) {
	err := newValidator(t).Validate(
		context.Background(),
		validation.Async(
			validation.String("", it.IsNotBlank().WithError(ErrFirst)),
		),
	)

	validationtest.Assert(t, err).IsViolationList().WithErrors(ErrFirst)
}

func TestAsyncArgument_WhenPathIsSet_ExpectOneViolationWithPath(t *testing.T) {
	err := newValidator(t).Validate(
		context.Background(),
		validation.Async(
			validation.String("", it.IsNotBlank().WithError(ErrFirst)),
		).With(
			validation.PropertyName("properties"),
			validation.ArrayIndex(0),
			validation.PropertyName("property"),
		),
	)

	violations := validationtest.Assert(t, err).IsViolationList()
	violations.HasViolationAt(0).WithError(ErrFirst).WithPropertyPath("properties[0].property")
}

func TestAsyncArgument_WhenValidationIsDisabled_ExpectNoErrors(t *testing.T) {
	err := newValidator(t).Validate(
		context.Background(),
		validation.Async(
			validation.String("", it.IsNotBlank().WithError(ErrFirst)),
			validation.String("", it.IsNotBlank().WithError(ErrSecond)),
		).When(false),
	)

	assert.NoError(t, err)
}

func TestAsyncArgument_WhenFatalError_ExpectContextCanceled(t *testing.T) {
	cancellation := make(chan bool, 1)
	fatal := fmt.Errorf("fatal")

	err := newValidator(t).Validate(
		context.Background(),
		validation.Async(
			validation.String("", asyncConstraint(func(ctx context.Context, validator *validation.Validator, value *string) error {
				return fatal
			})),
			validation.String("", asyncConstraint(func(ctx context.Context, validator *validation.Validator, value *string) error {
				select {
				case <-time.After(10 * time.Millisecond):
					cancellation <- false
				case <-ctx.Done():
					cancellation <- true
				}
				return nil
			})),
		),
	)

	assert.ErrorIs(t, err, fatal)
	if isCanceled, ok := <-cancellation; !isCanceled || !ok {
		assert.Fail(t, "context is expected to be canceled")
	}
}
