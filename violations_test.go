package validation_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/muonsoft/validation"
	"github.com/stretchr/testify/assert"
)

func TestViolation_Error_MessageOnly_ErrorWithMessage(t *testing.T) {
	validator := newValidator(t)

	violation := validator.BuildViolation(context.Background(), "", "message").CreateViolation()

	assert.Equal(t, "violation: message", violation.Error())
}

func TestNewViolationList(t *testing.T) {
	first := newViolationWithCode(t, "first")
	last := newViolationWithCode(t, "last")

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
		newViolationWithCode(t, "first"),
		newViolationWithCode(t, "second"),
		newViolationWithCode(t, "third"),
	)
	iterated := make([]string, 0)
	indices := make([]int, 0)

	err := violations.Each(func(i int, violation validation.Violation) error {
		iterated = append(iterated, violation.Code())
		indices = append(indices, i)
		return nil
	})

	assert.NoError(t, err)
	assert.Equal(t, []int{0, 1, 2}, indices)
	assert.Equal(t, []string{"first", "second", "third"}, iterated)
}

func TestViolationList_Each_WhenErrorReturned_ExpectLoopBreak(t *testing.T) {
	violations := validation.NewViolationList(
		newViolationWithCode(t, "first"),
		newViolationWithCode(t, "second"),
		newViolationWithCode(t, "third"),
	)
	iterated := make([]string, 0)
	indices := make([]int, 0)

	err := violations.Each(func(i int, violation validation.Violation) error {
		iterated = append(iterated, violation.Code())
		indices = append(indices, i)
		return fmt.Errorf("error at %d", i)
	})

	assert.EqualError(t, err, "error at 0")
	assert.Equal(t, []int{0}, indices)
	assert.Equal(t, []string{"first"}, iterated)
}

func TestViolationList_Join(t *testing.T) {
	tests := []struct {
		name          string
		list          *validation.ViolationList
		joined        *validation.ViolationList
		expectedCodes []string
	}{
		{
			name:          "nil joined list",
			list:          newViolationList(t),
			joined:        nil,
			expectedCodes: nil,
		},
		{
			name:          "empty joined list",
			list:          newViolationList(t),
			joined:        newViolationList(t),
			expectedCodes: nil,
		},
		{
			name:          "joined list with one element",
			list:          newViolationList(t),
			joined:        newViolationList(t, "code"),
			expectedCodes: []string{"code"},
		},
		{
			name:          "1 + 1",
			list:          newViolationList(t, "first"),
			joined:        newViolationList(t, "second"),
			expectedCodes: []string{"first", "second"},
		},
		{
			name:          "2 + 1",
			list:          newViolationList(t, "first", "second"),
			joined:        newViolationList(t, "third"),
			expectedCodes: []string{"first", "second", "third"},
		},
		{
			name:          "1 + 0",
			list:          newViolationList(t, "first"),
			joined:        newViolationList(t),
			expectedCodes: []string{"first"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.list.Join(test.joined)

			if assert.Equal(t, len(test.expectedCodes), test.list.Len()) {
				i := 0
				for e := test.list.First(); e != nil; e = e.Next() {
					assert.Equal(t, test.expectedCodes[i], e.Violation().Code())
					i++
				}
			}
		})
	}
}

func TestViolationList_Join_WhenEmptyListJoinedCoupleOfTimes_ExpectJoinedList(t *testing.T) {
	list := newViolationList(t, "first")
	list.Join(newViolationList(t))
	list.Join(newViolationList(t))

	assert.Equal(t, 1, list.Len())
}

func TestInternalViolation_Is(t *testing.T) {
	tests := []struct {
		name       string
		codes      []string
		expectedIs bool
	}{
		{
			name:       "empty list",
			expectedIs: false,
		},
		{
			name:       "no matches",
			codes:      []string{"alpha", "beta"},
			expectedIs: false,
		},
		{
			name:       "one of the codes is matching",
			codes:      []string{"alpha", "beta", "code"},
			expectedIs: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			violation := newViolationWithCode(t, "code")

			is := violation.Is(test.codes...)

			assert.Equal(t, test.expectedIs, is)
		})
	}
}

func TestViolationList_Has(t *testing.T) {
	tests := []struct {
		name       string
		codes      []string
		expectedIs bool
	}{
		{
			name:       "empty list",
			expectedIs: false,
		},
		{
			name:       "no matches",
			codes:      []string{"alpha", "beta"},
			expectedIs: false,
		},
		{
			name:       "one of the codes is matching",
			codes:      []string{"alpha", "beta", "code"},
			expectedIs: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			violations := validation.NewViolationList(newViolationWithCode(t, "code"))

			has := violations.Has(test.codes...)

			assert.Equal(t, test.expectedIs, has)
		})
	}
}

func TestViolationList_Filter_ViolationsWithCodes_FilteredList(t *testing.T) {
	violations := newViolationList(t, "alpha", "beta", "gamma", "delta")

	filtered := violations.Filter("delta", "beta").AsSlice()

	if assert.Len(t, filtered, 2) {
		assert.Equal(t, "beta", filtered[0].Code())
		assert.Equal(t, "delta", filtered[1].Code())
	}
}

func TestViolation_Error_MessageAndPropertyPath_ErrorWithPropertyPathAndMessage(t *testing.T) {
	validator := newValidator(t)
	violation := validator.BuildViolation(context.Background(), "", "message").
		SetPropertyPath(validation.NewPropertyPath(validation.PropertyNameElement("propertyPath"))).
		CreateViolation()

	err := violation.Error()

	assert.Equal(t, "violation at 'propertyPath': message", err)
}

func TestViolationList_Error_CoupleOfViolations_JoinedMessage(t *testing.T) {
	validator := newValidator(t)
	violations := validation.NewViolationList(
		validator.BuildViolation(context.Background(), "", "first message").
			SetPropertyPath(
				validation.NewPropertyPath(
					validation.PropertyNameElement("path"),
					validation.ArrayIndexElement(0)),
			).
			CreateViolation(),
		validator.BuildViolation(context.Background(), "", "second message").
			SetPropertyPath(
				validation.NewPropertyPath(
					validation.PropertyNameElement("path"),
					validation.ArrayIndexElement(1)),
			).
			CreateViolation(),
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
	err := fmt.Errorf("%w", newViolationWithCode(t, "code"))

	is := validation.IsViolation(err)

	assert.True(t, is)
}

func TestIsViolationList_CustomError_False(t *testing.T) {
	err := errors.New("error")

	is := validation.IsViolationList(err)

	assert.False(t, is)
}

func TestIsViolationList_Violation_True(t *testing.T) {
	err := fmt.Errorf("%w", validation.NewViolationList(newViolationWithCode(t, "code")))

	is := validation.IsViolationList(err)

	assert.True(t, is)
}

func TestUnwrapViolation_WrappedViolation_UnwrappedViolation(t *testing.T) {
	wrapped := newViolationWithCode(t, "code")
	err := fmt.Errorf("error: %w", wrapped)

	unwrapped, ok := validation.UnwrapViolation(err)

	assert.True(t, ok)
	assert.Equal(t, wrapped, unwrapped)
}

func TestUnwrapViolationList_WrappedViolationList_UnwrappedViolationList(t *testing.T) {
	wrapped := validation.NewViolationList(newViolationWithCode(t, "code"))
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
			violation: validator.BuildViolation(context.Background(), "code", "message").
				SetParameters(validation.TemplateParameter{Key: "key", Value: "value"}).
				SetPropertyPath(
					validation.NewPropertyPath(
						validation.PropertyNameElement("properties"),
						validation.ArrayIndexElement(1),
						validation.PropertyNameElement("name"),
					),
				).CreateViolation(),
			expectedJSON: `{
				"code": "code",
				"message": "message",
				"propertyPath": "properties[1].name"
			}`,
		},
		{
			name:         "empty data",
			violation:    validator.BuildViolation(context.Background(), "", "").CreateViolation(),
			expectedJSON: `{"code": "", "message": ""}`,
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
				validator.BuildViolation(context.Background(), "", "").CreateViolation(),
			),
			expectedJSON: `[{"code": "", "message": ""}]`,
		},
		{
			name: "one full violation",
			list: validation.NewViolationList(
				validator.BuildViolation(context.Background(), "code", "message").
					SetParameters(validation.TemplateParameter{Key: "key", Value: "value"}).
					SetPropertyPath(
						validation.NewPropertyPath(
							validation.PropertyNameElement("properties"),
							validation.ArrayIndexElement(1),
							validation.PropertyNameElement("name"),
						),
					).CreateViolation(),
			),
			expectedJSON: `[
				{
					"code": "code",
					"message": "message",
					"propertyPath": "properties[1].name"
				}
			]`,
		},
		{
			name: "two violations",
			list: validation.NewViolationList(
				validator.BuildViolation(context.Background(), "code", "message").
					SetParameters(validation.TemplateParameter{Key: "key", Value: "value"}).
					SetPropertyPath(
						validation.NewPropertyPath(
							validation.PropertyNameElement("properties"),
							validation.ArrayIndexElement(1),
							validation.PropertyNameElement("name"),
						),
					).CreateViolation(),
				validator.BuildViolation(context.Background(), "code", "message").
					SetParameters(validation.TemplateParameter{Key: "key", Value: "value"}).
					SetPropertyPath(
						validation.NewPropertyPath(
							validation.PropertyNameElement("properties"),
							validation.ArrayIndexElement(1),
							validation.PropertyNameElement("name"),
						),
					).CreateViolation(),
			),
			expectedJSON: `[
				{
					"code": "code",
					"message": "message",
					"propertyPath": "properties[1].name"
				},
				{
					"code": "code",
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

func newViolationWithCode(t *testing.T, code string) validation.Violation {
	t.Helper()
	validator := newValidator(t)
	violation := validator.BuildViolation(context.Background(), code, "").CreateViolation()
	return violation
}

func newViolationList(t *testing.T, codes ...string) *validation.ViolationList {
	t.Helper()
	validator := newValidator(t)
	violations := validation.NewViolationList()
	for _, code := range codes {
		violation := validator.BuildViolation(context.Background(), code, "").CreateViolation()
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
