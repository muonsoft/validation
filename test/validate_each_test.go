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

func TestValidateEach_WhenSliceOfStrings_ExpectViolationOnEachElement(t *testing.T) {
	strings := []string{"", ""}

	err := validator.Validate(context.Background(), validation.Each(strings, it.IsNotBlank()))

	validationtest.AssertIsViolationList(t, err, func(t *testing.T, violations []validation.Violation) bool {
		t.Helper()
		if assert.Len(t, violations, 2) {
			assert.Equal(t, code.NotBlank, violations[0].Code())
			assert.Equal(t, "[0]", violations[0].PropertyPath().String())
			assert.Equal(t, code.NotBlank, violations[1].Code())
			assert.Equal(t, "[1]", violations[1].PropertyPath().String())
		}
		return true
	})
}

func TestValidateEach_WhenMapOfStrings_ExpectViolationOnEachElement(t *testing.T) {
	strings := map[string]string{"key1": "", "key2": ""}

	err := validator.Validate(context.Background(), validation.Each(strings, it.IsNotBlank()))

	validationtest.AssertIsViolationList(t, err, func(t *testing.T, violations []validation.Violation) bool {
		t.Helper()
		if assert.Len(t, violations, 2) {
			assert.Equal(t, code.NotBlank, violations[0].Code())
			assert.Contains(t, violations[0].PropertyPath().String(), "key")
			assert.Equal(t, code.NotBlank, violations[1].Code())
			assert.Contains(t, violations[1].PropertyPath().String(), "key")
		}
		return true
	})
}

func TestValidateEachString_WhenSliceOfStrings_ExpectViolationOnEachElement(t *testing.T) {
	strings := []string{"", ""}

	err := validator.Validate(context.Background(), validation.EachString(strings, it.IsNotBlank()))

	validationtest.AssertIsViolationList(t, err, func(t *testing.T, violations []validation.Violation) bool {
		t.Helper()
		if assert.Len(t, violations, 2) {
			assert.Equal(t, code.NotBlank, violations[0].Code())
			assert.Equal(t, "[0]", violations[0].PropertyPath().String())
			assert.Equal(t, code.NotBlank, violations[1].Code())
			assert.Equal(t, "[1]", violations[1].PropertyPath().String())
		}
		return true
	})
}
