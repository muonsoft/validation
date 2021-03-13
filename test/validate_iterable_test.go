package test

import (
	"testing"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/code"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/validationtest"
	"github.com/muonsoft/validation/validator"
	"github.com/stretchr/testify/assert"
)

func TestValidateIterable_WhenSliceOfValidatable_ExpectViolationsWithValidPaths(t *testing.T) {
	strings := []mockValidatableString{{value: ""}}

	err := validator.ValidateValue(strings)

	validationtest.AssertIsViolationList(t, err, func(t *testing.T, violations validation.ViolationList) bool {
		t.Helper()
		if assert.Len(t, violations, 1) {
			assert.Equal(t, code.NotBlank, violations[0].GetCode())
			assert.Equal(t, "[0].value", violations[0].GetPropertyPath().String())
		}
		return true
	})
}

func TestValidateIterable_WhenSliceOfValidatableWithConstraints_ExpectCollectionViolationsWithValidPaths(t *testing.T) {
	strings := []mockValidatableString{{value: ""}}

	err := validator.ValidateValue(strings, it.HasMinCount(2))

	validationtest.AssertIsViolationList(t, err, func(t *testing.T, violations validation.ViolationList) bool {
		t.Helper()
		if assert.Len(t, violations, 2) {
			assert.Equal(t, code.CountTooFew, violations[0].GetCode())
			assert.Equal(t, "", violations[0].GetPropertyPath().String())
			assert.Equal(t, code.NotBlank, violations[1].GetCode())
			assert.Equal(t, "[0].value", violations[1].GetPropertyPath().String())
		}
		return true
	})
}

func TestValidateIterable_WhenMapOfValidatable_ExpectViolationsWithValidPaths(t *testing.T) {
	strings := map[string]mockValidatableString{"key": {value: ""}}

	err := validator.ValidateValue(strings)

	validationtest.AssertIsViolationList(t, err, func(t *testing.T, violations validation.ViolationList) bool {
		t.Helper()
		if assert.Len(t, violations, 1) {
			assert.Equal(t, code.NotBlank, violations[0].GetCode())
			assert.Equal(t, "key.value", violations[0].GetPropertyPath().String())
		}
		return true
	})
}

func TestValidateIterable_WhenMapOfValidatableWithConstraints_ExpectCollectionViolationsWithValidPaths(t *testing.T) {
	strings := map[string]mockValidatableString{"key": {value: ""}}

	err := validator.ValidateValue(strings, it.HasMinCount(2))

	validationtest.AssertIsViolationList(t, err, func(t *testing.T, violations validation.ViolationList) bool {
		t.Helper()
		if assert.Len(t, violations, 2) {
			assert.Equal(t, code.CountTooFew, violations[0].GetCode())
			assert.Equal(t, "", violations[0].GetPropertyPath().String())
			assert.Equal(t, code.NotBlank, violations[1].GetCode())
			assert.Equal(t, "key.value", violations[1].GetPropertyPath().String())
		}
		return true
	})
}
