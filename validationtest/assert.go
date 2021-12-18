package validationtest

import (
	"testing"

	"github.com/muonsoft/validation"
)

type AssertViolationFunc func(t *testing.T, violation validation.Violation) bool

type AssertViolationListFunc func(t *testing.T, violations []validation.Violation) bool

// Deprecated: use Assert instead.
func AssertIsViolation(t *testing.T, err error, assert AssertViolationFunc) bool {
	t.Helper()
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

// Deprecated: use Assert instead.
func AssertIsViolationList(t *testing.T, err error, assert AssertViolationListFunc) bool {
	t.Helper()
	if err == nil {
		t.Errorf("err is nil, but expected to be a ViolationList")
		return false
	}
	violations, is := validation.UnwrapViolationList(err)
	if !is {
		t.Errorf("failed asserting that error '%s' is ViolationList", err)
		return false
	}

	return assert(t, violations.AsSlice())
}

// Deprecated: use Assert instead.
func AssertOneViolationInList(t *testing.T, err error, assert AssertViolationFunc) bool {
	t.Helper()
	return AssertIsViolationList(t, err, func(t *testing.T, violations []validation.Violation) bool {
		t.Helper()
		if len(violations) != 1 {
			t.Errorf("failed asserting that violations list contains exactly one violation, actual count is %d", len(violations))
			return false
		}
		return assert(t, violations[0])
	})
}
