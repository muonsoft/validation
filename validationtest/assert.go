// Copyright 2021 Igor Lazarev. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package validationtest contains helper functions for testing purposes.
package validationtest

import (
	"testing"

	"github.com/muonsoft/validation"
)

type AssertViolationFunc func(t *testing.T, violation validation.Violation) bool

type AssertViolationListFunc func(t *testing.T, violations []validation.Violation) bool

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
