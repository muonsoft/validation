package test

import (
	"context"
	"testing"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/code"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/validationtest"
	"github.com/stretchr/testify/assert"
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
			Then(validation.String("", it.IsNotBlank().Code("then"))).
			Else(validation.String("", it.IsNotBlank().Code("else"))),
	)

	validationtest.Assert(t, err).IsViolationList().WithOneViolation().WithCode("then")
}

func TestWhenArgument_WhenConditionIsFalse_ExpectAllConstraintsOfElseBranchApplied(t *testing.T) {
	err := newValidator(t).Validate(
		context.Background(),
		validation.When(false).
			Then(validation.String("", it.IsNotBlank().Code("then"))).
			Else(validation.String("", it.IsNotBlank().Code("else"))),
	)

	validationtest.Assert(t, err).IsViolationList().WithOneViolation().WithCode("else")
}

func TestWhenArgument_WhenConditionIsFalseAndNoElseBranch_ExpectNoViolations(t *testing.T) {
	err := newValidator(t).Validate(
		context.Background(),
		validation.When(false).
			Then(validation.String("", it.IsNotBlank().Code("then"))),
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
			Then(validation.String("", it.IsNotBlank().Code("then"))),
	)

	validationtest.Assert(t, err).IsViolationList().WithOneViolation().
		WithCode("then").
		WithPropertyPath("properties[0].property")
}

func TestWhenGroupsArgument_WhenGroupMatches_ExpectViolation(t *testing.T) {
	err := newValidator(t).WithGroups(testGroup).Validate(
		context.Background(),
		validation.WhenGroups(testGroup).
			Then(validation.String("", it.IsNotBlank().WhenGroups(testGroup).Code("then"))).
			Else(validation.String("", it.IsNotBlank().WhenGroups(testGroup).Code("else"))),
	)

	validationtest.Assert(t, err).IsViolationList().WithOneViolation().WithCode("then")
}

func TestWhenGroupsArgument_WhenGroupNotMatches_ExpectNoError(t *testing.T) {
	err := newValidator(t).Validate(
		context.Background(),
		validation.WhenGroups(testGroup).
			Then(validation.String("", it.IsNotBlank().Code("then"))).
			Else(validation.String("", it.IsNotBlank().Code("else"))),
	)

	validationtest.Assert(t, err).IsViolationList().WithOneViolation().WithCode("else")
}

func TestWhenGroupsArgument_WhenGroupNotMatchesNoElseBranch_ExpectNoViolations(t *testing.T) {
	err := newValidator(t).Validate(
		context.Background(),
		validation.WhenGroups(testGroup).
			Then(validation.String("", it.IsNotBlank().Code("then"))),
	)

	assertNoError(t, err)
}

func TestWhenGroupsArgument_WhenPathIsSet_ExpectViolationWithPath(t *testing.T) {
	err := newValidator(t).WithGroups(testGroup).Validate(
		context.Background(),
		validation.WhenGroups(testGroup).
			Then(validation.String("", it.IsNotBlank().WhenGroups(testGroup).Code("then"))).
			With(
				validation.PropertyName("properties"),
				validation.ArrayIndex(0),
				validation.PropertyName("property"),
			),
	)

	validationtest.Assert(t, err).IsViolationList().WithOneViolation().
		WithCode("then").
		WithPropertyPath("properties[0].property")
}

func TestSequentialArgument_WhenInvalidValueAtFirstConstraint_ExpectOneViolation(t *testing.T) {
	err := newValidator(t).Validate(
		context.Background(),
		validation.Sequentially(
			validation.String("", it.IsNotBlank().Code("first")),
			validation.String("", it.IsNotBlank().Code("second")),
		),
	)

	validationtest.Assert(t, err).IsViolationList().WithOneViolation().WithCode("first")
}

func TestSequentialArgument_WhenPathIsSet_ExpectOneViolationWithPath(t *testing.T) {
	err := newValidator(t).Validate(
		context.Background(),
		validation.Sequentially(
			validation.String("", it.IsNotBlank().Code("first")),
			validation.String("", it.IsNotBlank().Code("second")),
		).With(
			validation.PropertyName("properties"),
			validation.ArrayIndex(0),
			validation.PropertyName("property"),
		),
	)

	validationtest.Assert(t, err).IsViolationList().WithOneViolation().
		WithCode("first").
		WithPropertyPath("properties[0].property")
}

func TestSequentialArgument_WhenValidationIsDisabled_ExpectNoErrors(t *testing.T) {
	err := newValidator(t).Validate(
		context.Background(),
		validation.Sequentially(
			validation.String("", it.IsNotBlank().Code("first")),
			validation.String("", it.IsNotBlank().Code("second")),
		).When(false),
	)

	assert.NoError(t, err)
}

func TestAtLeastOneOfArgument_WhenInvalidValueAtFirstConstraint_ExpectAllViolations(t *testing.T) {
	err := newValidator(t).Validate(
		context.Background(),
		validation.AtLeastOneOf(
			validation.String("", it.IsNotBlank().Code("first")),
			validation.String("", it.IsNotBlank().Code("second")),
		),
	)

	validationtest.Assert(t, err).IsViolationList().WithCodes("first", "second")
}

func TestAtLeastOneOfArgument_WhenInvalidValueAtSecondConstraint_ExpectNoViolation(t *testing.T) {
	err := newValidator(t).Validate(
		context.Background(),
		validation.AtLeastOneOf(
			validation.String("", it.IsNotBlank().Code("first")),
			validation.String("foo", it.IsNotBlank().Code("second")),
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
		WithCode(code.NotBlank).
		WithPropertyPath("properties[0].property")
}

func TestAtLeastOneOfArgument_WhenValidationIsDisabled_ExpectNoError(t *testing.T) {
	err := newValidator(t).Validate(
		context.Background(),
		validation.AtLeastOneOf(
			validation.String("", it.IsNotBlank().Code("first")),
			validation.String("", it.IsNotBlank().Code("second")),
		).When(false),
	)

	assert.NoError(t, err)
}
