package test

import (
	"context"
	"testing"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/code"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/validationtest"
	"github.com/muonsoft/validation/validator"
)

func TestValidate_Value_WhenSliceOfValidatable_ExpectViolationsWithValidPaths(t *testing.T) {
	strings := []mockValidatableString{{value: ""}}

	err := validator.Validate(context.Background(), validation.Value(strings))

	validationtest.Assert(t, err).IsViolationList().
		WithOneViolation().
		WithCode(code.NotBlank).
		WithPropertyPath("[0].value")
}

func TestValidate_Value_WhenSliceOfValidatableWithConstraints_ExpectCollectionViolationsWithValidPaths(t *testing.T) {
	strings := []mockValidatableString{{value: ""}}

	err := validator.Validate(context.Background(), validation.Value(strings, it.HasMinCount(2)))

	list := validationtest.Assert(t, err).IsViolationList()
	list.HasViolationAt(0).WithCode(code.CountTooFew).WithPropertyPath("")
	list.HasViolationAt(1).WithCode(code.NotBlank).WithPropertyPath("[0].value")
}

func TestValidate_Value_WhenMapOfValidatable_ExpectViolationsWithValidPaths(t *testing.T) {
	strings := map[string]mockValidatableString{"key": {value: ""}}

	err := validator.Validate(context.Background(), validation.Value(strings))

	validationtest.Assert(t, err).IsViolationList().
		WithOneViolation().
		WithCode(code.NotBlank).
		WithPropertyPath("key.value")
}

func TestValidate_Value_WhenMapOfValidatableWithConstraints_ExpectCollectionViolationsWithValidPaths(t *testing.T) {
	strings := map[string]mockValidatableString{"key": {value: ""}}

	err := validator.Validate(context.Background(), validation.Value(strings, it.HasMinCount(2)))

	list := validationtest.Assert(t, err).IsViolationList()
	list.HasViolationAt(0).WithCode(code.CountTooFew).WithPropertyPath("")
	list.HasViolationAt(1).WithCode(code.NotBlank).WithPropertyPath("key.value")
}
