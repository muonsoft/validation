package test

import (
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/code"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/validationtest"
	"github.com/muonsoft/validation/validator"
	"github.com/stretchr/testify/assert"

	"regexp"
	"testing"
	"unicode/utf8"
)

func TestValidate_WhenInvalidValueForElseBranchOfConditionConstraint_ExpectViolations(t *testing.T) {
	value := "name"

	err := validator.ValidateString(
		&value,
		validation.When(utf8.RuneCountInString(value) <= 3).
			Then(
				it.Matches(regexp.MustCompile(`^\\w$`)),
			).
			Else(
				it.Matches(regexp.MustCompile(`^\\d$`)),
			),
	)

	validationtest.AssertIsViolationList(t, err, func(t *testing.T, violations validation.ViolationList) bool {
		t.Helper()
		return assert.Len(t, violations, 1) &&
			assert.Equal(t, code.MatchingFailed, violations[0].Code())
	})
}

func TestValidate_WhenInvalidValueForThenBranchOfConditionConstraint_ExpectViolations(t *testing.T) {
	value := "name"

	err := validator.ValidateString(
		&value,
		validation.When(utf8.RuneCountInString(value) <= 4).
			Then(
				it.Matches(regexp.MustCompile(`^\\d$`)),
			).
			Else(
				it.Matches(regexp.MustCompile(`^\\w$`)),
			),
	)

	validationtest.AssertIsViolationList(t, err, func(t *testing.T, violations validation.ViolationList) bool {
		t.Helper()
		return assert.Len(t, violations, 1) &&
			assert.Equal(t, code.MatchingFailed, violations[0].Code())
	})
}
