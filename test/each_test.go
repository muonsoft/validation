package test

import (
	"context"
	"errors"
	"testing"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/validationtest"
	"github.com/muonsoft/validation/validator"
	"github.com/stretchr/testify/assert"
)

func TestEach_WhenInvalidElement_ExpectViolationWithIndexInPath(t *testing.T) {
	items := []string{"ok", "", "valid"}
	err := validator.Validate(context.Background(), validation.Each(items, eachStringNotBlankConstraint))

	validationtest.Assert(t, err).IsViolationList().WithLen(1)
	validationtest.Assert(t, err).IsViolationList().HasViolationAt(0).WithPropertyPath("[1]")
}

func TestEach_WhenMultipleInvalidElements_ExpectAllViolationsCollected(t *testing.T) {
	items := []string{"", ""}
	err := validator.Validate(context.Background(), validation.Each(items, eachStringNotBlankConstraint))

	validationtest.Assert(t, err).IsViolationList().WithLen(2)
	validationtest.Assert(t, err).IsViolationList().HasViolationAt(0).WithPropertyPath("[0]")
	validationtest.Assert(t, err).IsViolationList().HasViolationAt(1).WithPropertyPath("[1]")
}

func TestEach_WhenValidElements_ExpectNoViolations(t *testing.T) {
	items := []string{"a", "b", "c"}
	err := validator.Validate(context.Background(), validation.Each(items, eachStringNotBlankConstraint))

	assert.NoError(t, err)
}

func TestEach_WhenEmptySlice_ExpectNoViolations(t *testing.T) {
	var items []string
	err := validator.Validate(context.Background(), validation.Each(items, eachStringNotBlankConstraint))

	assert.NoError(t, err)
}

func TestEachProperty_WhenInvalidElement_ExpectPropertyNameInPath(t *testing.T) {
	items := []string{""}
	err := validator.Validate(context.Background(), validation.EachProperty("tags", items, eachStringNotBlankConstraint))

	validationtest.Assert(t, err).IsViolationList().WithOneViolation().WithPropertyPath("tags[0]")
}

func TestEachProperty_WhenMultipleInvalid_ExpectAllViolationsWithPath(t *testing.T) {
	items := []string{"", ""}
	err := validator.Validate(context.Background(), validation.EachProperty("items", items, eachStringNotBlankConstraint))

	validationtest.Assert(t, err).IsViolationList().WithLen(2)
	validationtest.Assert(t, err).IsViolationList().HasViolationAt(0).WithPropertyPath("items[0]")
	validationtest.Assert(t, err).IsViolationList().HasViolationAt(1).WithPropertyPath("items[1]")
}

func TestFunc_ImplementsConstraint_WhenUsedInEach(t *testing.T) {
	customConstraint := validation.Func[string](func(ctx context.Context, v *validation.Validator, s string) error {
		return v.Validate(ctx, validation.String(s, it.HasMinLength(3)))
	})
	items := []string{"ab"}
	err := validator.Validate(context.Background(), validation.Each(items, customConstraint))

	validationtest.Assert(t, err).IsViolationList().WithOneViolation().WithPropertyPath("[0]")
}

func TestEach_WhenConstraintReturnsFatalError_ExpectErrorPropagated(t *testing.T) {
	fatalErr := errors.New("fatal")
	fatalConstraint := validation.Func[string](func(context.Context, *validation.Validator, string) error {
		return fatalErr
	})
	items := []string{"x"}
	err := validator.Validate(context.Background(), validation.Each(items, fatalConstraint))

	assert.ErrorIs(t, err, fatalErr)
}

func TestEach_WhenUsedWithThis_ExpectSameConstraintType(t *testing.T) {
	// Each accepts Constraint[E], same as This; Func[E] implements Constraint[E]
	single := ""
	errSingle := validator.Validate(context.Background(), validation.This(single, eachStringNotBlankConstraint))
	validationtest.Assert(t, errSingle).IsViolationList().WithOneViolation()

	slice := []string{""}
	errSlice := validator.Validate(context.Background(), validation.Each(slice, eachStringNotBlankConstraint))
	validationtest.Assert(t, errSlice).IsViolationList().WithOneViolation().WithPropertyPath("[0]")
}
