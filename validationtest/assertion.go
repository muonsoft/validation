// Copyright 2021 Igor Lazarev. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package validationtest contains helper functions for testing purposes.
package validationtest

import (
	"strconv"
	"testing"

	"github.com/muonsoft/validation"
	"github.com/stretchr/testify/assert"
)

// TestingT is an interface wrapper around *testing.T.
type TestingT interface {
	Helper()
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
}

// ViolationAttributes are used to compare violation against expected values. An empty value is not compared.
type ViolationAttributes struct {
	Code         string
	Message      string
	PropertyPath string
}

// Assertion is a structure for testing an error for implementing validator violations.
type Assertion struct {
	t   TestingT
	err error
}

// Assert creates a new Assertion for the error.
func Assert(t TestingT, err error) *Assertion {
	t.Helper()

	return &Assertion{t: t, err: err}
}

// IsViolation checks that err implements validation.Violation and returns ViolationAssertion
// for attributes assertions.
func (a *Assertion) IsViolation() *ViolationAssertion {
	a.t.Helper()

	violation, ok := validation.UnwrapViolation(a.err)
	if !ok {
		a.t.Error("failed asserting that err is a Violation")

		return nil
	}

	return newViolationAssertion(a.t, violation)
}

// IsViolationList checks that err implements validation.IsViolationList and returns ViolationListAssertion
// for attributes assertions.
func (a *Assertion) IsViolationList() *ViolationListAssertion {
	a.t.Helper()

	violations, ok := validation.UnwrapViolationList(a.err)
	if !ok {
		a.t.Error("failed asserting that err is a ViolationList")

		return nil
	}

	return &ViolationListAssertion{t: a.t, violations: violations}
}

// ViolationListAssertion is a structure for testing violation list attributes.
type ViolationListAssertion struct {
	t          TestingT
	violations *validation.ViolationList
}

// Assert is used for the client-side side assertion of violations by a callback function.
func (a *ViolationListAssertion) Assert(
	assert func(tb testing.TB, violations []validation.Violation),
) *ViolationListAssertion {
	if a == nil {
		return nil
	}
	a.t.Helper()

	if tb, ok := a.t.(testing.TB); ok {
		assert(tb, a.violations.AsSlice())
	} else {
		a.t.Fatal("t must implement testing.TB")
	}

	return a
}

// WithLen checks that the violation list has exact length.
func (a *ViolationListAssertion) WithLen(length int) *ViolationListAssertion {
	if a == nil {
		return nil
	}
	a.t.Helper()

	actual := a.violations.Len()
	if actual != length {
		a.t.Errorf(
			"failed asserting that violation list length is equal to %d, actual is %d",
			length,
			actual,
		)
	}

	return a
}

// WithOneViolation checks that the violation list contains exactly one violation and returns
// a ViolationAssertion to test it.
func (a *ViolationListAssertion) WithOneViolation() *ViolationAssertion {
	if a == nil {
		return nil
	}
	a.t.Helper()

	if a.violations.Len() != 1 {
		a.t.Error("failed asserting that violation list contains exactly one violation")
		return nil
	}

	return newViolationAssertionAt(a.t, a.violations.First(), 0)
}

// HasViolationAt checks that the violation list contains element at specific index and returns
// a ViolationAssertion to test it.
func (a *ViolationListAssertion) HasViolationAt(index int) *ViolationAssertion {
	if a == nil {
		return nil
	}
	a.t.Helper()

	violations := a.violations.AsSlice()
	if index >= len(violations) {
		a.t.Errorf("failed asserting that violation list contains violation at index %d", index)
		return nil
	}

	return newViolationAssertionAt(a.t, violations[index], 0)
}

// WithCodes checks that the violation list contains violations with specific codes in a given order.
func (a *ViolationListAssertion) WithCodes(codes ...string) *ViolationListAssertion {
	if a == nil {
		return nil
	}
	a.t.Helper()

	length := a.violations.Len()
	if length != len(codes) {
		a.t.Errorf(
			"failed asserting that violation list length is equal to %d, actual is %d",
			len(codes),
			length,
		)
		return a
	}

	a.violations.Each(func(i int, violation validation.Violation) error {
		if violation.Code() != codes[i] {
			a.t.Errorf(
				`failed asserting that violation at %d has code "%s", actual is "%s"`,
				i,
				codes[i],
				violation.Code(),
			)
		}
		return nil
	})

	return a
}

// WithAttributes checks that the violation list contains violations with the expected attributes in a given order.
// Empty values are not compared.
func (a *ViolationListAssertion) WithAttributes(violations ...ViolationAttributes) *ViolationListAssertion {
	if a == nil {
		return nil
	}
	a.t.Helper()

	length := a.violations.Len()
	if length != len(violations) {
		a.t.Errorf(
			"failed asserting that violation list length is equal to %d, actual is %d",
			len(violations),
			length,
		)
		return a
	}

	a.violations.Each(func(i int, violation validation.Violation) error {
		expected := violations[i]

		if expected.Code != "" && violation.Code() != expected.Code {
			a.t.Errorf(
				`failed asserting that violation at %d has code "%s", actual is "%s"`,
				i,
				expected.Code,
				violation.Code(),
			)
		}

		if expected.Message != "" && violation.Message() != expected.Message {
			a.t.Errorf(
				`failed asserting that violation at %d has message "%s", actual is "%s"`,
				i,
				expected.Message,
				violation.Message(),
			)
		}

		if expected.PropertyPath != "" && violation.PropertyPath().String() != expected.PropertyPath {
			a.t.Errorf(
				`failed asserting that violation at %d has property path "%s", actual is "%s"`,
				i,
				expected.PropertyPath,
				violation.PropertyPath().String(),
			)
		}

		return nil
	})

	return a
}

// ViolationAssertion is a structure for testing violation attributes.
type ViolationAssertion struct {
	t         TestingT
	violation validation.Violation
	index     int
}

func newViolationAssertion(t TestingT, violation validation.Violation) *ViolationAssertion {
	return &ViolationAssertion{t: t, violation: violation, index: -1}
}

func newViolationAssertionAt(t TestingT, violation validation.Violation, index int) *ViolationAssertion {
	return &ViolationAssertion{t: t, violation: violation, index: index}
}

// Assert is used for the client-side assertion of the violation by a callback function.
func (a *ViolationAssertion) Assert(
	assert func(tb testing.TB, violation validation.Violation),
) *ViolationAssertion {
	if a == nil {
		return nil
	}
	a.t.Helper()

	if tb, ok := a.t.(testing.TB); ok {
		assert(tb, a.violation)
	} else {
		a.t.Fatal("t must implement testing.TB")
	}

	return a
}

// WithCode checks that violation has expected code.
func (a *ViolationAssertion) WithCode(code string) *ViolationAssertion {
	if a == nil {
		return nil
	}
	a.t.Helper()

	actual := a.violation.Code()
	if actual != code {
		a.t.Errorf(
			`failed asserting that violation%s has code "%s", actual is "%s"`,
			a.atIndex(),
			code,
			actual,
		)
	}

	return a
}

// WithMessage checks that violation has expected message.
func (a *ViolationAssertion) WithMessage(message string) *ViolationAssertion {
	if a == nil {
		return nil
	}
	a.t.Helper()

	actual := a.violation.Message()
	if actual != message {
		a.t.Errorf(
			`failed asserting that violation%s has message "%s", actual is "%s"`,
			a.atIndex(),
			message,
			actual,
		)
	}

	return a
}

// WithPropertyPath checks that the tested violation has an expected property path.
func (a *ViolationAssertion) WithPropertyPath(path string) *ViolationAssertion {
	if a == nil {
		return nil
	}
	a.t.Helper()

	actual := a.violation.PropertyPath().String()
	if actual != path {
		a.t.Errorf(
			`failed asserting that violation%s has property path "%s", actual is "%s"`,
			a.atIndex(),
			path,
			actual,
		)
	}

	return a
}

// EqualTo checks that the tested assertion is equal to the expected one.
func (a *ViolationAssertion) EqualTo(violation validation.Violation) *ViolationAssertion {
	if a == nil {
		return nil
	}
	a.t.Helper()

	assert.Equal(a.t, violation, a.violation, "failed asserting that violations are equal")

	return a
}

// EqualToError checks that violation rendered to an error is equal to the expected one.
func (a *ViolationAssertion) EqualToError(errString string) *ViolationAssertion {
	if a == nil {
		return nil
	}
	a.t.Helper()

	actual := a.violation.Error()
	if actual != errString {
		a.t.Errorf(
			`failed asserting that violation%s error is equal to "%s", actual is "%s"`,
			a.atIndex(),
			errString,
			actual,
		)
	}

	return a
}

func (a *ViolationAssertion) atIndex() string {
	if a.index < 0 {
		return ""
	}
	return " #" + strconv.Itoa(a.index)
}