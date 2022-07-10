package validationtest_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/validationtest"
	"github.com/muonsoft/validation/validator"
	"github.com/stretchr/testify/assert"
)

func TestAssertion_IsViolation(t *testing.T) {
	tester := &Tester{}

	validationtest.Assert(tester, fmt.Errorf("error")).IsViolation()

	tester.AssertOneMessage(t, "failed asserting that err is a Violation")
}

func TestAssertion_IsViolationList(t *testing.T) {
	tester := &Tester{}

	validationtest.Assert(tester, fmt.Errorf("error")).IsViolationList()

	tester.AssertOneMessage(t, "failed asserting that err is a ViolationList")
}

func TestViolationListAssertion_WithLen(t *testing.T) {
	tester := &Tester{}
	violations := validation.NewViolationList()

	validationtest.Assert(tester, violations).IsViolationList().WithLen(1)

	tester.AssertOneMessage(t, "failed asserting that violation list length is equal to 1, actual is 0")
}

func TestViolationListAssertion_WithOneViolation(t *testing.T) {
	tester := &Tester{}
	violations := validation.NewViolationList()

	validationtest.Assert(tester, violations).IsViolationList().WithOneViolation()

	tester.AssertOneMessage(t, "failed asserting that violation list contains exactly one violation")
}

func TestViolationListAssertion_HasViolationAt(t *testing.T) {
	tester := &Tester{}
	violations := validation.NewViolationList()

	validationtest.Assert(tester, violations).IsViolationList().HasViolationAt(5)

	tester.AssertOneMessage(t, "failed asserting that violation list contains violation at index 5")
}

func TestViolationListAssertion_WithErrors_WhenEmptyList_ExpectError(t *testing.T) {
	tester := &Tester{}
	violations := validation.NewViolationList()

	validationtest.Assert(tester, violations).IsViolationList().WithErrors(errors.New("one"))

	tester.AssertOneMessage(t, "failed asserting that violation list length is equal to 1, actual is 0")
}

func TestViolationListAssertion_WithErrors_WhenInvalidError_ExpectError(t *testing.T) {
	tester := &Tester{}
	violations := validation.NewViolationList(
		validator.BuildViolation(context.Background(), errors.New("error"), "message").Create(),
	)

	validationtest.Assert(tester, violations).IsViolationList().WithErrors(errors.New("expected"))

	tester.AssertOneMessage(t, `failed asserting that violation at 0 has error "expected", actual is "error"`)
}

func TestViolationListAssertion_WithAttributes(t *testing.T) {
	violation := validator.BuildViolation(context.Background(), errors.New("error"), "message").
		SetPropertyPath(validation.NewPropertyPath(validation.PropertyName("path"))).
		Create()

	tests := []struct {
		name            string
		violations      *validation.ViolationList
		attributes      validationtest.ViolationAttributes
		expectedMessage string
	}{
		{
			name:            "empty list",
			expectedMessage: "failed asserting that violation list length is equal to 1, actual is 0",
		},
		{
			name:            "invalid error",
			violations:      validation.NewViolationList(violation),
			attributes:      validationtest.ViolationAttributes{Error: errors.New("expected")},
			expectedMessage: `failed asserting that violation at 0 has error "expected", actual is "error"`,
		},
		{
			name:            "invalid message",
			violations:      validation.NewViolationList(violation),
			attributes:      validationtest.ViolationAttributes{Message: "expected"},
			expectedMessage: `failed asserting that violation at 0 has message "expected", actual is "message"`,
		},
		{
			name:            "invalid path",
			violations:      validation.NewViolationList(violation),
			attributes:      validationtest.ViolationAttributes{PropertyPath: "expected"},
			expectedMessage: `failed asserting that violation at 0 has property path "expected", actual is "path"`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tester := &Tester{}

			validationtest.Assert(tester, test.violations).IsViolationList().WithAttributes(test.attributes)

			tester.AssertOneMessage(t, test.expectedMessage)
		})
	}
}

func TestViolationAssertion_WithError(t *testing.T) {
	tester := &Tester{}
	violation := validator.BuildViolation(context.Background(), errors.New("error"), "message").
		SetPropertyPath(validation.NewPropertyPath(validation.PropertyName("path"))).
		Create()

	validationtest.Assert(tester, violation).IsViolation().WithError(errors.New("expected"))

	tester.AssertOneMessage(t, `failed asserting that violation has error "expected", actual is "error"`)
}

func TestViolationAssertion_WithError_AtIndex(t *testing.T) {
	tester := &Tester{}
	violation := validator.BuildViolation(context.Background(), errors.New("error"), "message").
		SetPropertyPath(validation.NewPropertyPath(validation.PropertyName("path"))).
		Create()
	violations := validation.NewViolationList(violation)

	validationtest.Assert(tester, violations).IsViolationList().WithOneViolation().WithError(errors.New("expected"))

	tester.AssertOneMessage(t, `failed asserting that violation #0 has error "expected", actual is "error"`)
}

func TestViolationAssertion_WithMessage(t *testing.T) {
	tester := &Tester{}
	violation := validator.BuildViolation(context.Background(), errors.New("error"), "message").
		SetPropertyPath(validation.NewPropertyPath(validation.PropertyName("path"))).
		Create()

	validationtest.Assert(tester, violation).IsViolation().WithMessage("expected")

	tester.AssertOneMessage(t, `failed asserting that violation has message "expected", actual is "message"`)
}

func TestViolationAssertion_WithPropertyPath(t *testing.T) {
	tester := &Tester{}
	violation := validator.BuildViolation(context.Background(), errors.New("error"), "message").
		SetPropertyPath(validation.NewPropertyPath(validation.PropertyName("path"))).
		Create()

	validationtest.Assert(tester, violation).IsViolation().WithPropertyPath("expected")

	tester.AssertOneMessage(t, `failed asserting that violation has property path "expected", actual is "path"`)
}

func TestViolationAssertion_EqualTo(t *testing.T) {
	tester := &Tester{}
	violation := validator.BuildViolation(context.Background(), errors.New("error"), "message").
		SetPropertyPath(validation.NewPropertyPath(validation.PropertyName("path"))).
		Create()
	differentViolation := validator.BuildViolation(context.Background(), errors.New("error"), "message").Create()

	validationtest.Assert(tester, violation).IsViolation().EqualTo(differentViolation)

	assert.Len(t, tester.messages, 1)
}

func TestViolationAssertion_EqualError(t *testing.T) {
	tester := &Tester{}
	violation := validator.BuildViolation(context.Background(), errors.New("error"), "message").
		SetPropertyPath(validation.NewPropertyPath(validation.PropertyName("path"))).
		Create()

	validationtest.Assert(tester, violation).IsViolation().EqualToError("expected")

	tester.AssertOneMessage(t, `failed asserting that violation error is equal to "expected", actual is "violation at 'path': message"`)
}

type Tester struct {
	messages []string
}

func (tester *Tester) Helper() {
}

func (tester *Tester) Error(args ...interface{}) {
	tester.messages = append(tester.messages, fmt.Sprint(args...))
}

func (tester *Tester) Errorf(format string, args ...interface{}) {
	tester.messages = append(tester.messages, fmt.Sprintf(format, args...))
}

func (tester *Tester) Fatal(args ...interface{}) {
	tester.messages = append(tester.messages, fmt.Sprint(args...))
}

func (tester *Tester) AssertOneMessage(t *testing.T, message string) {
	t.Helper()
	if len(tester.messages) != 1 {
		t.Errorf("failed asserting that tester has exactly one message, actual count is %d", len(tester.messages))
		return
	}
	if !strings.Contains(tester.messages[0], message) {
		t.Errorf(`failed asserting that tester message "%s" contains "%s"`, tester.messages[0], message)
	}
}
