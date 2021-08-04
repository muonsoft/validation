package test

import (
	"context"
	"testing"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/code"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/validationtest"
	"github.com/muonsoft/validation/validator"
	"github.com/stretchr/testify/assert"
)

func TestValidate_Value_WhenSliceOfValidatable_ExpectViolationsWithValidPaths(t *testing.T) {
	strings := []mockValidatableString{{value: ""}}

	err := validator.Validate(context.Background(), validation.Value(strings))

	validationtest.AssertIsViolationList(t, err, func(t *testing.T, violations []validation.Violation) bool {
		t.Helper()
		if assert.Len(t, violations, 1) {
			assert.Equal(t, code.NotBlank, violations[0].Code())
			assert.Equal(t, "[0].value", violations[0].PropertyPath().String())
		}
		return true
	})
}

func TestValidate_Value_WhenSliceOfValidatableWithConstraints_ExpectCollectionViolationsWithValidPaths(t *testing.T) {
	strings := []mockValidatableString{{value: ""}}

	err := validator.Validate(context.Background(), validation.Value(strings, it.HasMinCount(2)))

	validationtest.AssertIsViolationList(t, err, func(t *testing.T, violations []validation.Violation) bool {
		t.Helper()
		if assert.Len(t, violations, 2) {
			assert.Equal(t, code.CountTooFew, violations[0].Code())
			assert.Equal(t, "", violations[0].PropertyPath().String())
			assert.Equal(t, code.NotBlank, violations[1].Code())
			assert.Equal(t, "[0].value", violations[1].PropertyPath().String())
		}
		return true
	})
}

func TestValidate_Value_WhenMapOfValidatable_ExpectViolationsWithValidPaths(t *testing.T) {
	strings := map[string]mockValidatableString{"key": {value: ""}}

	err := validator.Validate(context.Background(), validation.Value(strings))

	validationtest.AssertIsViolationList(t, err, func(t *testing.T, violations []validation.Violation) bool {
		t.Helper()
		if assert.Len(t, violations, 1) {
			assert.Equal(t, code.NotBlank, violations[0].Code())
			assert.Equal(t, "key.value", violations[0].PropertyPath().String())
		}
		return true
	})
}

func TestValidate_Value_WhenMapOfValidatableWithConstraints_ExpectCollectionViolationsWithValidPaths(t *testing.T) {
	strings := map[string]mockValidatableString{"key": {value: ""}}

	err := validator.Validate(context.Background(), validation.Value(strings, it.HasMinCount(2)))

	validationtest.AssertIsViolationList(t, err, func(t *testing.T, violations []validation.Violation) bool {
		t.Helper()
		if assert.Len(t, violations, 2) {
			assert.Equal(t, code.CountTooFew, violations[0].Code())
			assert.Equal(t, "", violations[0].PropertyPath().String())
			assert.Equal(t, code.NotBlank, violations[1].Code())
			assert.Equal(t, "key.value", violations[1].PropertyPath().String())
		}
		return true
	})
}
