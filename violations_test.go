package validation_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/validationtest"
	"github.com/stretchr/testify/assert"
)

var ErrTest = errors.New("test")

func TestViolation_Error_MessageOnly_ErrorWithMessage(t *testing.T) {
	validator := newValidator(t)

	violation := validator.BuildViolation(context.Background(), ErrTest, "message").Create()

	assert.Equal(t, "violation: message", violation.Error())
}

func TestNewViolationList(t *testing.T) {
	first := newViolationWithError(t, errors.New("first"))
	last := newViolationWithError(t, errors.New("last"))

	violations := validation.NewViolationList(first, last)

	assert.Equal(t, 2, violations.Len())
	if assert.NotNil(t, violations.First()) {
		assert.Equal(t, first, violations.First().Violation())
	}
	if assert.NotNil(t, violations.Last()) {
		assert.Equal(t, last, violations.Last().Violation())
	}
}

func TestViolationList_Len_WhenNil_ExpectZero(t *testing.T) {
	var violations *validation.ViolationList

	length := violations.Len()

	assert.Equal(t, 0, length)
}

func TestViolationList_Each_WhenMultipleViolations_ExpectAllIterated(t *testing.T) {
	violations := validation.NewViolationList(
		newViolationWithError(t, errors.New("first")),
		newViolationWithError(t, errors.New("second")),
		newViolationWithError(t, errors.New("third")),
	)
	iterated := make([]string, 0)
	indices := make([]int, 0)

	err := violations.ForEach(func(i int, violation validation.Violation) error {
		iterated = append(iterated, violation.Unwrap().Error())
		indices = append(indices, i)
		return nil
	})

	assert.NoError(t, err)
	assert.Equal(t, []int{0, 1, 2}, indices)
	assert.Equal(t, []string{"first", "second", "third"}, iterated)
}

func TestViolationList_Each_WhenErrorReturned_ExpectLoopBreak(t *testing.T) {
	violations := validation.NewViolationList(
		newViolationWithError(t, errors.New("first")),
		newViolationWithError(t, errors.New("second")),
		newViolationWithError(t, errors.New("third")),
	)
	iterated := make([]string, 0)
	indices := make([]int, 0)

	err := violations.ForEach(func(i int, violation validation.Violation) error {
		iterated = append(iterated, violation.Unwrap().Error())
		indices = append(indices, i)
		return fmt.Errorf("error at %d", i)
	})

	assert.EqualError(t, err, "error at 0")
	assert.Equal(t, []int{0}, indices)
	assert.Equal(t, []string{"first"}, iterated)
}

func TestViolationList_Join(t *testing.T) {
	tests := []struct {
		name           string
		list           *validation.ViolationList
		joined         *validation.ViolationList
		expectedErrors []string
	}{
		{
			name:           "nil joined list",
			list:           newViolationList(t),
			joined:         nil,
			expectedErrors: nil,
		},
		{
			name:           "empty joined list",
			list:           newViolationList(t),
			joined:         newViolationList(t),
			expectedErrors: nil,
		},
		{
			name:           "joined list with one element",
			list:           newViolationList(t),
			joined:         newViolationList(t, errors.New("code")),
			expectedErrors: []string{"code"},
		},
		{
			name:           "1 + 1",
			list:           newViolationList(t, errors.New("first")),
			joined:         newViolationList(t, errors.New("second")),
			expectedErrors: []string{"first", "second"},
		},
		{
			name:           "2 + 1",
			list:           newViolationList(t, errors.New("first"), errors.New("second")),
			joined:         newViolationList(t, errors.New("third")),
			expectedErrors: []string{"first", "second", "third"},
		},
		{
			name:           "1 + 0",
			list:           newViolationList(t, errors.New("first")),
			joined:         newViolationList(t),
			expectedErrors: []string{"first"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.list.Join(test.joined)

			if assert.Equal(t, len(test.expectedErrors), test.list.Len()) {
				i := 0
				for e := test.list.First(); e != nil; e = e.Next() {
					assert.Equal(t, test.expectedErrors[i], e.Violation().Unwrap().Error())
					i++
				}
			}
		})
	}
}

func TestViolationList_Join_WhenEmptyListJoinedCoupleOfTimes_ExpectJoinedList(t *testing.T) {
	list := newViolationList(t, errors.New("first"))
	list.Join(newViolationList(t))
	list.Join(newViolationList(t))

	assert.Equal(t, 1, list.Len())
}

func TestInternalViolation_Is(t *testing.T) {
	tests := []struct {
		name       string
		violation  validation.Violation
		err        error
		expectedIs bool
	}{
		{
			name:       "nil error",
			violation:  newViolationWithError(t, ErrTest),
			expectedIs: false,
		},
		{
			name:       "matching error",
			violation:  newViolationWithError(t, ErrTest),
			err:        ErrTest,
			expectedIs: true,
		},
		{
			name:       "not equal by value",
			violation:  newViolationWithError(t, ErrTest),
			err:        errors.New(ErrTest.Error()),
			expectedIs: false,
		},
		{
			name:       "wrapped error",
			violation:  newViolationWithError(t, fmt.Errorf("%w", ErrTest)),
			err:        ErrTest,
			expectedIs: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			is := errors.Is(test.violation, test.err)

			assert.Equal(t, test.expectedIs, is)
		})
	}
}

func TestViolationList_Is(t *testing.T) {
	tests := []struct {
		name       string
		violation  validation.Violation
		err        error
		expectedIs bool
	}{
		{
			name:       "nil error",
			violation:  newViolationWithError(t, ErrTest),
			expectedIs: false,
		},
		{
			name:       "matching error",
			violation:  newViolationWithError(t, ErrTest),
			err:        ErrTest,
			expectedIs: true,
		},
		{
			name:       "not equal by value",
			violation:  newViolationWithError(t, ErrTest),
			err:        errors.New(ErrTest.Error()),
			expectedIs: false,
		},
		{
			name:       "wrapped error",
			violation:  newViolationWithError(t, fmt.Errorf("%w", ErrTest)),
			err:        ErrTest,
			expectedIs: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			violations := validation.NewViolationList(test.violation)

			is := errors.Is(violations, test.err)

			assert.Equal(t, test.expectedIs, is)
		})
	}
}

func TestViolationList_Filter_ViolationsWithCodes_FilteredList(t *testing.T) {
	errA := errors.New("alpha")
	errB := errors.New("beta")
	errC := errors.New("gamma")
	errD := errors.New("delta")
	violations := newViolationList(t, errA, errB, errC, errD)

	filtered := violations.Filter(errD, errB).AsSlice()

	if assert.Len(t, filtered, 2) {
		assert.Equal(t, errB, filtered[0].Unwrap())
		assert.Equal(t, errD, filtered[1].Unwrap())
	}
}

func TestViolation_Error_MessageAndPropertyPath_ErrorWithPropertyPathAndMessage(t *testing.T) {
	validator := newValidator(t)
	violation := validator.BuildViolation(context.Background(), ErrTest, "message").
		At(validation.PropertyName("propertyPath")).
		Create()

	err := violation.Error()

	assert.Equal(t, "violation at 'propertyPath': message", err)
}

func TestViolationList_Error_CoupleOfViolations_JoinedMessage(t *testing.T) {
	validator := newValidator(t)
	violations := validation.NewViolationList(
		validator.BuildViolation(context.Background(), ErrTest, "first message").
			At(
				validation.PropertyName("path"),
				validation.ArrayIndex(0),
			).
			Create(),
		validator.BuildViolation(context.Background(), ErrTest, "second message").
			At(
				validation.PropertyName("path"),
				validation.ArrayIndex(1),
			).
			Create(),
	)

	err := violations.Error()

	assert.Equal(t, "violation at 'path[0]': first message; violation at 'path[1]': second message", err)
}

func TestViolationList_Error_EmptyList_ErrorWithHelpMessage(t *testing.T) {
	violations := validation.NewViolationList()

	err := violations.Error()

	assert.Equal(t, "the list of violations is empty, it looks like you forgot to use the AsError method somewhere", err)
}

func TestIsViolation_CustomError_False(t *testing.T) {
	err := errors.New("error")

	is := validation.IsViolation(err)

	assert.False(t, is)
}

func TestIsViolation_Violation_True(t *testing.T) {
	err := fmt.Errorf("%w", newViolationWithError(t, ErrTest))

	is := validation.IsViolation(err)

	assert.True(t, is)
}

func TestIsViolationList_CustomError_False(t *testing.T) {
	err := errors.New("error")

	is := validation.IsViolationList(err)

	assert.False(t, is)
}

func TestIsViolationList_Violation_True(t *testing.T) {
	err := fmt.Errorf("%w", validation.NewViolationList(newViolationWithError(t, ErrTest)))

	is := validation.IsViolationList(err)

	assert.True(t, is)
}

func TestUnwrapViolation_WrappedViolation_UnwrappedViolation(t *testing.T) {
	wrapped := newViolationWithError(t, ErrTest)
	err := fmt.Errorf("error: %w", wrapped)

	unwrapped, ok := validation.UnwrapViolation(err)

	assert.True(t, ok)
	assert.Equal(t, wrapped, unwrapped)
}

func TestUnwrapViolationList_WrappedViolationList_UnwrappedViolationList(t *testing.T) {
	wrapped := validation.NewViolationList(newViolationWithError(t, ErrTest))
	err := fmt.Errorf("error: %w", wrapped)

	unwrapped, ok := validation.UnwrapViolationList(err)

	assert.True(t, ok)
	assert.Equal(t, wrapped, unwrapped)
}

func TestUnwrapViolationList_NoError_NoListAndFalse(t *testing.T) {
	unwrapped, ok := validation.UnwrapViolationList(nil)

	assert.Nil(t, unwrapped)
	assert.False(t, ok)
}

func TestMarshalViolationToJSON(t *testing.T) {
	validator := newValidator(t)

	tests := []struct {
		name         string
		violation    validation.Violation
		expectedJSON string
	}{
		{
			name: "full data",
			violation: validator.BuildViolation(context.Background(), ErrTest, "message").
				WithParameters(validation.TemplateParameter{Key: "key", Value: "value"}).
				At(
					validation.PropertyName("properties"),
					validation.ArrayIndex(1),
					validation.PropertyName("name"),
				).Create(),
			expectedJSON: `{
				"error": "test",
				"message": "message",
				"propertyPath": "properties[1].name"
			}`,
		},
		{
			name:         "empty data",
			violation:    validator.BuildViolation(context.Background(), nil, "").Create(),
			expectedJSON: `{"message": ""}`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			data, err := json.Marshal(test.violation)

			if assert.NoError(t, err) {
				assert.JSONEq(t, test.expectedJSON, string(data))
			}
		})
	}
}

func TestMarshalViolationListToJSON(t *testing.T) {
	validator := newValidator(t)

	tests := []struct {
		name         string
		list         *validation.ViolationList
		expectedJSON string
	}{
		{
			name:         "empty list",
			list:         validation.NewViolationList(),
			expectedJSON: `[]`,
		},
		{
			name: "empty data",
			list: validation.NewViolationList(
				validator.BuildViolation(context.Background(), nil, "").Create(),
			),
			expectedJSON: `[{"message": ""}]`,
		},
		{
			name: "one full violation",
			list: validation.NewViolationList(
				validator.BuildViolation(context.Background(), ErrTest, "message").
					WithParameters(validation.TemplateParameter{Key: "key", Value: "value"}).
					At(
						validation.PropertyName("properties"),
						validation.ArrayIndex(1),
						validation.PropertyName("name"),
					).Create(),
			),
			expectedJSON: `[
				{
					"error": "test",
					"message": "message",
					"propertyPath": "properties[1].name"
				}
			]`,
		},
		{
			name: "two violations",
			list: validation.NewViolationList(
				validator.BuildViolation(context.Background(), ErrTest, "message").
					WithParameters(validation.TemplateParameter{Key: "key", Value: "value"}).
					At(
						validation.PropertyName("properties"),
						validation.ArrayIndex(1),
						validation.PropertyName("name"),
					).Create(),
				validator.BuildViolation(context.Background(), ErrTest, "message").
					WithParameters(validation.TemplateParameter{Key: "key", Value: "value"}).
					At(
						validation.PropertyName("properties"),
						validation.ArrayIndex(1),
						validation.PropertyName("name"),
					).Create(),
			),
			expectedJSON: `[
				{
					"error": "test",
					"message": "message",
					"propertyPath": "properties[1].name"
				},
				{
					"error": "test",
					"message": "message",
					"propertyPath": "properties[1].name"
				}
			]`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			data, err := json.Marshal(test.list)

			if assert.NoError(t, err) {
				assert.JSONEq(t, test.expectedJSON, string(data))
			}
		})
	}
}

func TestValidator_BuildViolationList_WhenBasePath_ExpectViolationListBuiltWithPath(t *testing.T) {
	validator := newValidator(t)
	err1 := errors.New("error 1")
	err2 := errors.New("error 2")
	err3 := errors.New("error 3")

	err := validator.
		AtProperty("base").AtIndex(1).
		BuildViolationList(context.Background()).
		At(validation.PropertyName("listPath")).
		AtProperty("list").AtIndex(0).
		BuildViolation(err1, "message with {{ parameter 1 }}").
		WithParameter("{{ parameter 1 }}", "value 1").
		At(validation.PropertyName("propertyPath")).
		AtProperty("properties").AtIndex(0).
		Add().
		BuildViolation(err2, "message with {{ parameter 2 }}").
		WithParameter("{{ parameter 2 }}", "value 2").
		At(validation.PropertyName("propertyPath")).
		AtProperty("properties").AtIndex(1).
		Add().
		AddViolation(err3, "message 3", validation.PropertyName("singleProperty")).
		Create().
		AsError()

	validationtest.Assert(t, err).IsViolationList().WithAttributes(
		validationtest.ViolationAttributes{
			Error:        err1,
			Message:      "message with value 1",
			PropertyPath: "base[1].listPath.list[0].propertyPath.properties[0]",
		},
		validationtest.ViolationAttributes{
			Error:        err2,
			Message:      "message with value 2",
			PropertyPath: "base[1].listPath.list[0].propertyPath.properties[1]",
		},
		validationtest.ViolationAttributes{
			Error:        err3,
			Message:      "message 3",
			PropertyPath: "base[1].listPath.list[0].singleProperty",
		},
	)
}

func newViolationWithError(t *testing.T, err error) validation.Violation {
	t.Helper()
	validator := newValidator(t)
	violation := validator.BuildViolation(context.Background(), err, "").Create()
	return violation
}

func newViolationList(t *testing.T, errs ...error) *validation.ViolationList {
	t.Helper()
	validator := newValidator(t)
	violations := validation.NewViolationList()
	for _, err := range errs {
		violation := validator.BuildViolation(context.Background(), err, "").Create()
		violations.Append(violation)
	}
	return violations
}

func newValidator(t *testing.T) *validation.Validator {
	t.Helper()
	validator, err := validation.NewValidator()
	if err != nil {
		t.Fatal(err)
	}
	return validator
}
