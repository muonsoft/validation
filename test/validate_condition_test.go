package test

import (
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/code"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/validationtest"
	"github.com/muonsoft/validation/validator"
	"github.com/stretchr/testify/assert"

	"testing"
)

type Brand struct {
	Name string
	Tags []string
}

func (b Brand) Validate(validator *validation.Validator) error {
	return validator.Validate(
		validation.String(
			&b.Name,
			validation.PropertyName("name"),
			it.IsNotBlank(),
			validation.When(len(b.Tags) <= 2).
				Then(
					it.HasMinLength(3),
				).
				Else(
					it.HasMaxLength(10),
				),
		),
		validation.Iterable(
			b.Tags,
			validation.PropertyName("tags"),
			it.HasMinCount(1),
		),
	)
}

func TestValidateValue_WithConditionConstraints_ExpectViolationsByElseBranch(t *testing.T) {
	b := Brand{
		Name: "a",
		Tags: []string{
			"tag 1",
			"tag 2",
		},
	}

	err := validator.ValidateValue(b)

	validationtest.AssertIsViolationList(t, err, func(t *testing.T, violations validation.ViolationList) bool {
		t.Helper()
		if assert.Len(t, violations, 1) {
			assert.Equal(t, code.LengthTooFew, violations[0].Code())
			assert.Equal(t, "name", violations[0].PropertyPath().String())
		}
		return true
	})
}

func TestValidateValue_WithConditionConstraints_ExpectViolationsByThenBranch(t *testing.T) {
	b := Brand{
		Name: "name length more than 10",
		Tags: []string{
			"tag 1",
			"tag 2",
			"tag 3",
		},
	}

	err := validator.ValidateValue(b)

	validationtest.AssertIsViolationList(t, err, func(t *testing.T, violations validation.ViolationList) bool {
		t.Helper()
		if assert.Len(t, violations, 1) {
			assert.Equal(t, code.LengthTooMany, violations[0].Code())
			assert.Equal(t, "name", violations[0].PropertyPath().String())
		}
		return true
	})
}
