package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/code"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/message"
	"github.com/muonsoft/validation/validationtest"
	"github.com/muonsoft/validation/validator"
	"github.com/stretchr/testify/assert"
)

func TestValidate_WhenArgumentForGivenType_ExpectValidationExecuted(t *testing.T) {
	tests := []struct {
		name     string
		argument validation.Argument
	}{
		{"Nil", validation.Nil(true, it.IsNotNil())},
		{"Bool", validation.Bool(false, it.IsNotBlank())},
		{"NilBool", validation.NilBool(boolValue(false), it.IsNotBlank())},
		{"Number", validation.Number[int](0, it.IsNotBlankNumber[int]())},
		{"NilNumber", validation.NilNumber[int](intValue(0), it.IsNotBlankNumber[int]())},
		{"String", validation.String("", it.IsNotBlank())},
		{"NilString", validation.NilString(stringValue(""), it.IsNotBlank())},
		{"Countable", validation.Countable(0, it.IsNotBlank())},
		{"Time", validation.Time(time.Time{}, it.IsNotBlank())},
		{"NilTime", validation.NilTime(nilTime, it.IsNotBlank())},
		{"Comparable", validation.Comparable[string]("foo", it.IsOneOf("bar"))},
		{"NilComparable", validation.NilComparable[string](stringValue("foo"), it.IsOneOf("bar"))},
		{"Comparables", validation.Comparables[string]([]string{"foo", "foo"}, it.HasUniqueValues[string]())},
		{"EachString", validation.EachString([]string{""}, it.IsNotBlank())},
		{"EachNumber", validation.EachNumber[int]([]int{0}, it.IsNotBlankNumber[int]())},
		{"Valid", validation.Valid(mockValidatableString{""})},
		{"ValidSlice", validation.ValidSlice([]mockValidatableString{{""}})},
		{"ValidMap", validation.ValidMap(map[string]mockValidatableString{"key": {""}})},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validator.Validate(context.Background(), test.argument)

			validationtest.Assert(t, err).IsViolationList().WithOneViolation()
		})
	}
}

func TestValidate_WhenPropertyArgument_ExpectValidPathInViolation(t *testing.T) {
	tests := []struct {
		name         string
		argument     validation.Argument
		expectedPath string
	}{
		{
			name: "NilProperty",
			argument: validation.NilProperty("property", true, it.IsNotNil()).
				With(validation.PropertyName("internal")),
			expectedPath: "property.internal",
		},
		{
			name: "BoolProperty",
			argument: validation.BoolProperty("property", false, it.IsNotBlank()).
				With(validation.PropertyName("internal")),
			expectedPath: "property.internal",
		},
		{
			name: "NilBoolProperty",
			argument: validation.NilBoolProperty("property", boolValue(false), it.IsNotBlank()).
				With(validation.PropertyName("internal")),
			expectedPath: "property.internal",
		},
		{
			name: "NumberProperty",
			argument: validation.NumberProperty[int]("property", 0, it.IsNotBlankNumber[int]()).
				With(validation.PropertyName("internal")),
			expectedPath: "property.internal",
		},
		{
			name: "NilNumberProperty",
			argument: validation.NilNumberProperty[int]("property", intValue(0), it.IsNotBlankNumber[int]()).
				With(validation.PropertyName("internal")),
			expectedPath: "property.internal",
		},
		{
			name: "StringProperty",
			argument: validation.StringProperty("property", "", it.IsNotBlank()).
				With(validation.PropertyName("internal")),
			expectedPath: "property.internal",
		},
		{
			name: "NilStringProperty",
			argument: validation.NilStringProperty("property", stringValue(""), it.IsNotBlank()).
				With(validation.PropertyName("internal")),
			expectedPath: "property.internal",
		},
		{
			name: "CountableProperty",
			argument: validation.CountableProperty("property", 0, it.IsNotBlank()).
				With(validation.PropertyName("internal")),
			expectedPath: "property.internal",
		},
		{
			name: "TimeProperty",
			argument: validation.TimeProperty("property", time.Time{}, it.IsNotBlank()).
				With(validation.PropertyName("internal")),
			expectedPath: "property.internal",
		},
		{
			name: "NilTimeProperty",
			argument: validation.NilTimeProperty("property", nilTime, it.IsNotBlank()).
				With(validation.PropertyName("internal")),
			expectedPath: "property.internal",
		},
		{
			name: "ComparableProperty",
			argument: validation.ComparableProperty[string]("property", "foo", it.IsOneOf("bar")).
				With(validation.PropertyName("internal")),
			expectedPath: "property.internal",
		},
		{
			name: "NilComparableProperty",
			argument: validation.NilComparableProperty[string]("property", stringValue("foo"), it.IsOneOf("bar")).
				With(validation.PropertyName("internal")),
			expectedPath: "property.internal",
		},
		{
			name: "ComparablesProperty",
			argument: validation.ComparablesProperty[string]("property", []string{"foo", "foo"}, it.HasUniqueValues[string]()).
				With(validation.PropertyName("internal")),
			expectedPath: "property.internal",
		},
		{
			name: "EachStringProperty",
			argument: validation.EachStringProperty("property", []string{""}, it.IsNotBlank()).
				With(validation.PropertyName("internal")),
			expectedPath: "property.internal[0]",
		},
		{
			name: "EachNumberProperty",
			argument: validation.EachNumberProperty[int]("property", []int{0}, it.IsNotBlankNumber[int]()).
				With(validation.PropertyName("internal")),
			expectedPath: "property.internal[0]",
		},
		{
			name: "ValidProperty",
			argument: validation.ValidProperty("property", mockValidatableString{""}).
				With(validation.PropertyName("internal")),
			expectedPath: "property.internal.value",
		},
		{
			name: "ValidSliceProperty",
			argument: validation.ValidSliceProperty("property", []mockValidatableString{{""}}).
				With(validation.PropertyName("internal")),
			expectedPath: "property.internal[0].value",
		},
		{
			name: "ValidMapProperty",
			argument: validation.ValidMapProperty("property", map[string]mockValidatableString{"key": {""}}).
				With(validation.PropertyName("internal")),
			expectedPath: "property.internal.key.value",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validator.Validate(context.Background(), test.argument)

			validationtest.Assert(t, err).IsViolationList().WithOneViolation().WithPropertyPath(test.expectedPath)
		})
	}
}

func TestCheck_WhenFalse_ExpectViolation(t *testing.T) {
	err := validator.Validate(context.Background(), validation.Check(false))

	validationtest.Assert(t, err).IsViolationList().
		WithOneViolation().
		WithCode(code.NotValid).
		WithMessage(message.Templates[code.NotValid]).
		WithPropertyPath("")
}

func TestCheck_WhenFalseWithPath_ExpectViolationWithPath(t *testing.T) {
	err := validator.Validate(
		context.Background(),
		validation.Check(false).With(
			validation.PropertyName("properties"),
			validation.ArrayIndex(0),
			validation.PropertyName("property"),
		),
	)

	validationtest.Assert(t, err).IsViolationList().
		WithOneViolation().
		WithCode(code.NotValid).
		WithMessage(message.Templates[code.NotValid]).
		WithPropertyPath("properties[0].property")
}

func TestCheck_WhenCustomCodeAndTemplate_ExpectCodeAndTemplateInViolation(t *testing.T) {
	err := validator.Validate(
		context.Background(),
		validation.Check(false).
			Code("custom").
			Message("message with {{ value }}", validation.TemplateParameter{Key: "{{ value }}", Value: "value"}),
	)

	validationtest.Assert(t, err).IsViolationList().
		WithOneViolation().
		WithCode("custom").
		WithMessage("message with value").
		WithPropertyPath("")
}

func TestCheck_WhenTrue_ExpectNoViolation(t *testing.T) {
	err := validator.Validate(context.Background(), validation.Check(true))

	assert.NoError(t, err)
}

func TestCheck_When_WhenConditionIsFalse_ExpectNoViolation(t *testing.T) {
	err := validator.Validate(context.Background(), validation.Check(false).When(false))

	assert.NoError(t, err)
}

func TestCheck_When_WhenGroupsNotMatch_ExpectNoViolation(t *testing.T) {
	err := validator.Validate(context.Background(), validation.Check(false).WhenGroups(testGroup))

	assert.NoError(t, err)
}

func TestCheckProperty_WhenFalse_ExpectPropertyNameInViolation(t *testing.T) {
	err := validator.Validate(context.Background(), validation.CheckProperty("propertyName", false))

	validationtest.Assert(t, err).IsViolationList().WithAttributes(
		validationtest.ViolationAttributes{
			Code:         code.NotValid,
			Message:      message.Templates[code.NotValid],
			PropertyPath: "propertyName",
		},
	)
}

func TestCheckNoViolations_WhenThereAreViolations_ExpectAppendedViolationsReturned(t *testing.T) {
	violations := validator.ValidateString(context.Background(), "", it.IsNotBlank())

	err := validator.Validate(
		context.Background(),
		validation.CheckNoViolations(violations),
		validation.Check(false),
	)

	validationtest.Assert(t, err).IsViolationList().WithCodes(code.NotBlank, code.NotValid)
}

func TestCheckNoViolations_WhenThereIsAnError_ExpectError(t *testing.T) {
	err := validator.Validate(
		context.Background(),
		validation.CheckNoViolations(fmt.Errorf("error")),
		validation.Check(false),
	)

	assert.EqualError(t, err, "error")
}
