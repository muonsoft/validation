package validationtest

import (
	"testing"

	"github.com/muonsoft/validation"
)

type AssertViolationFunc func(t *testing.T, violation validation.Violation) bool

type AssertViolationListFunc func(t *testing.T, violations validation.ViolationList) bool

func AssertIsViolation(t *testing.T, err error, assert AssertViolationFunc) bool {
	violation, is := validation.UnwrapViolation(err)
	if !is {
		t.Errorf("failed asserting that error '%s' is Violation", err)
		return false
	}

	return assert(t, violation)
}

func AssertIsViolationList(t *testing.T, err error, assert AssertViolationListFunc) bool {
	violations, is := validation.UnwrapViolationList(err)
	if !is {
		t.Errorf("failed asserting that error '%s' is ViolationList", err)
		return false
	}

	return assert(t, violations)
}
