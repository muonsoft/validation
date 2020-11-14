package validationtest

import (
	"testing"

	"github.com/muonsoft/validation"
)

type AssertViolationFunc func(t *testing.T, violation validation.Violation) bool

type AssertViolationListFunc func(t *testing.T, violations validation.ViolationList) bool

func AssertIsViolation(t *testing.T, err error, assert AssertViolationFunc) bool {
	if err == nil {
		t.Errorf("err is nil, but expected to be a Violation")
		return false
	}
	violation, is := validation.UnwrapViolation(err)
	if !is {
		t.Errorf("failed asserting that error '%s' is Violation", err)
		return false
	}

	return assert(t, violation)
}

func AssertIsViolationList(t *testing.T, err error, assert AssertViolationListFunc) bool {
	if err == nil {
		t.Errorf("err is nil, but expected to be a ViolationList")
		return false
	}
	violations, is := validation.UnwrapViolationList(err)
	if !is {
		t.Errorf("failed asserting that error '%s' is ViolationList", err)
		return false
	}

	return assert(t, violations)
}
