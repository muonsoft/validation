package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/muonsoft/validation"
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
		{"Slice", validation.Slice([]string{"a"}, mockFailingSliceConstraint{})},
		{"EachString", validation.EachString([]string{""}, it.IsNotBlank())},
		{"EachNumber", validation.EachNumber[int]([]int{0}, it.IsNotBlankNumber[int]())},
		{"EachComparable", validation.EachComparable[int]([]int{1}, it.IsOneOf(2))},
		{"Each", validation.Each([]string{""}, it.IsNotBlank())},
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
				At(validation.PropertyName("internal")),
			expectedPath: "property.internal",
		},
		{
			name: "BoolProperty",
			argument: validation.BoolProperty("property", false, it.IsNotBlank()).
				At(validation.PropertyName("internal")),
			expectedPath: "property.internal",
		},
		{
			name: "NilBoolProperty",
			argument: validation.NilBoolProperty("property", boolValue(false), it.IsNotBlank()).
				At(validation.PropertyName("internal")),
			expectedPath: "property.internal",
		},
		{
			name: "NumberProperty",
			argument: validation.NumberProperty[int]("property", 0, it.IsNotBlankNumber[int]()).
				At(validation.PropertyName("internal")),
			expectedPath: "property.internal",
		},
		{
			name: "NilNumberProperty",
			argument: validation.NilNumberProperty[int]("property", intValue(0), it.IsNotBlankNumber[int]()).
				At(validation.PropertyName("internal")),
			expectedPath: "property.internal",
		},
		{
			name: "StringProperty",
			argument: validation.StringProperty("property", "", it.IsNotBlank()).
				At(validation.PropertyName("internal")),
			expectedPath: "property.internal",
		},
		{
			name: "NilStringProperty",
			argument: validation.NilStringProperty("property", stringValue(""), it.IsNotBlank()).
				At(validation.PropertyName("internal")),
			expectedPath: "property.internal",
		},
		{
			name: "CountableProperty",
			argument: validation.CountableProperty("property", 0, it.IsNotBlank()).
				At(validation.PropertyName("internal")),
			expectedPath: "property.internal",
		},
		{
			name: "TimeProperty",
			argument: validation.TimeProperty("property", time.Time{}, it.IsNotBlank()).
				At(validation.PropertyName("internal")),
			expectedPath: "property.internal",
		},
		{
			name: "NilTimeProperty",
			argument: validation.NilTimeProperty("property", nilTime, it.IsNotBlank()).
				At(validation.PropertyName("internal")),
			expectedPath: "property.internal",
		},
		{
			name: "ComparableProperty",
			argument: validation.ComparableProperty[string]("property", "foo", it.IsOneOf("bar")).
				At(validation.PropertyName("internal")),
			expectedPath: "property.internal",
		},
		{
			name: "NilComparableProperty",
			argument: validation.NilComparableProperty[string]("property", stringValue("foo"), it.IsOneOf("bar")).
				At(validation.PropertyName("internal")),
			expectedPath: "property.internal",
		},
		{
			name: "ComparablesProperty",
			argument: validation.ComparablesProperty[string]("property", []string{"foo", "foo"}, it.HasUniqueValues[string]()).
				At(validation.PropertyName("internal")),
			expectedPath: "property.internal",
		},
		{
			name: "SliceProperty",
			argument: validation.SliceProperty("property", []string{"a"}, mockFailingSliceConstraint{}).
				At(validation.PropertyName("internal")),
			expectedPath: "property.internal",
		},
		{
			name: "EachStringProperty",
			argument: validation.EachStringProperty("property", []string{""}, it.IsNotBlank()).
				At(validation.PropertyName("internal")),
			expectedPath: "property.internal[0]",
		},
		{
			name: "EachNumberProperty",
			argument: validation.EachNumberProperty[int]("property", []int{0}, it.IsNotBlankNumber[int]()).
				At(validation.PropertyName("internal")),
			expectedPath: "property.internal[0]",
		},
		{
			name: "EachComparableProperty",
			argument: validation.EachComparableProperty[string]("property", []string{"foo"}, it.IsOneOf("bar")).
				At(validation.PropertyName("internal")),
			expectedPath: "property.internal[0]",
		},
		{
			name: "EachProperty",
			argument: validation.EachProperty("property", []string{""}, it.IsNotBlank()).
				At(validation.PropertyName("internal")),
			expectedPath: "property.internal[0]",
		},
		{
			name: "ValidProperty",
			argument: validation.ValidProperty("property", mockValidatableString{""}).
				At(validation.PropertyName("internal")),
			expectedPath: "property.internal.value",
		},
		{
			name: "ValidSliceProperty",
			argument: validation.ValidSliceProperty("property", []mockValidatableString{{""}}).
				At(validation.PropertyName("internal")),
			expectedPath: "property.internal[0].value",
		},
		{
			name: "ValidMapProperty",
			argument: validation.ValidMapProperty("property", map[string]mockValidatableString{"key": {""}}).
				At(validation.PropertyName("internal")),
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
		WithError(validation.ErrNotValid).
		WithMessage(message.NotValid).
		WithPropertyPath("")
}

func TestCheck_WhenFalseWithPath_ExpectViolationWithPath(t *testing.T) {
	err := validator.Validate(
		context.Background(),
		validation.Check(false).At(
			validation.PropertyName("properties"),
			validation.ArrayIndex(0),
			validation.PropertyName("property"),
		),
	)

	validationtest.Assert(t, err).IsViolationList().
		WithOneViolation().
		WithError(validation.ErrNotValid).
		WithMessage(message.NotValid).
		WithPropertyPath("properties[0].property")
}

func TestCheck_WhenCustomCodeAndTemplate_ExpectCodeAndTemplateInViolation(t *testing.T) {
	err := validator.Validate(
		context.Background(),
		validation.Check(false).
			WithError(ErrCustom).
			WithMessage("message with {{ value }}", validation.TemplateParameter{Key: "{{ value }}", Value: "value"}),
	)

	validationtest.Assert(t, err).IsViolationList().
		WithOneViolation().
		WithError(ErrCustom).
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
			Error:        validation.ErrNotValid,
			Message:      message.NotValid,
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

	validationtest.Assert(t, err).IsViolationList().WithErrors(validation.ErrIsBlank, validation.ErrNotValid)
}

func TestCheckNoViolations_WhenThereIsAnError_ExpectError(t *testing.T) {
	err := validator.Validate(
		context.Background(),
		validation.CheckNoViolations(fmt.Errorf("error")),
		validation.Check(false),
	)

	assert.EqualError(t, err, "error")
}

type itemWithID struct {
	ID string
}

func TestSlice_HasUniqueValuesBy_WhenUniqueKeys_ExpectNoError(t *testing.T) {
	items := []itemWithID{{ID: "a"}, {ID: "b"}, {ID: "c"}}

	err := validator.Validate(
		context.Background(),
		validation.Slice(items, it.HasUniqueValuesBy(func(x itemWithID) string { return x.ID })),
	)

	assert.NoError(t, err)
}

func TestSlice_HasUniqueValuesBy_WhenDuplicateKeys_ExpectViolationsAtIndex(t *testing.T) {
	items := []itemWithID{{ID: "a"}, {ID: "b"}, {ID: "a"}}

	err := validator.Validate(
		context.Background(),
		validation.Slice(items, it.HasUniqueValuesBy(func(x itemWithID) string { return x.ID })),
	)

	validationtest.Assert(t, err).IsViolationList().WithLen(2).
		WithErrors(validation.ErrNotUnique, validation.ErrNotUnique)
	// Both indices 0 and 2 have duplicate key "a"
	list, _ := validation.UnwrapViolations(err)
	paths := make([]string, 0, 2)
	list.ForEach(func(i int, v validation.Violation) error {
		paths = append(paths, v.PropertyPath().String())
		return nil
	})
	assert.Contains(t, paths, "[0]")
	assert.Contains(t, paths, "[2]")
}

func TestSlice_HasUniqueValuesBy_WhenEmptySlice_ExpectNoError(t *testing.T) {
	items := []itemWithID{}

	err := validator.Validate(
		context.Background(),
		validation.Slice(items, it.HasUniqueValuesBy(func(x itemWithID) string { return x.ID })),
	)

	assert.NoError(t, err)
}

func TestSlice_HasUniqueValuesBy_WhenNilSlice_ExpectNoError(t *testing.T) {
	var items []itemWithID

	err := validator.Validate(
		context.Background(),
		validation.Slice(items, it.HasUniqueValuesBy(func(x itemWithID) string { return x.ID })),
	)

	assert.NoError(t, err)
}

func TestSlice_HasUniqueValuesBy_WhenSkipWhen_ExpectSkippedItemsNotChecked(t *testing.T) {
	items := []itemWithID{{ID: ""}, {ID: "a"}, {ID: ""}} // empty ID skipped, so only "a" counted once

	err := validator.Validate(
		context.Background(),
		validation.Slice(items, it.HasUniqueValuesBy(func(x itemWithID) string { return x.ID }).
			SkipWhen(func(x itemWithID) bool { return x.ID == "" })),
	)

	assert.NoError(t, err)
}

func TestSlice_HasUniqueValuesBy_WhenSkipEmptyKeys_ExpectEmptyKeysIgnored(t *testing.T) {
	items := []itemWithID{{ID: ""}, {ID: "a"}, {ID: "b"}, {ID: ""}} // empty keys skipped, "a" and "b" unique

	err := validator.Validate(
		context.Background(),
		validation.Slice(items, it.HasUniqueValuesBy(func(x itemWithID) string { return x.ID }).SkipEmptyKeys()),
	)

	assert.NoError(t, err)
}

func TestSlice_HasUniqueValuesBy_WhenSkipEmptyKeys_AndDuplicateNonEmptyKeys_ExpectViolations(t *testing.T) {
	items := []itemWithID{{ID: ""}, {ID: "a"}, {ID: "a"}, {ID: ""}}

	err := validator.Validate(
		context.Background(),
		validation.Slice(items, it.HasUniqueValuesBy(func(x itemWithID) string { return x.ID }).SkipEmptyKeys()),
	)

	validationtest.Assert(t, err).IsViolationList().WithLen(2).WithErrors(validation.ErrNotUnique, validation.ErrNotUnique)
}

func TestSliceProperty_HasUniqueValuesBy_WhenDuplicateKeys_ExpectPropertyPathWithIndex(t *testing.T) {
	items := []itemWithID{{ID: "x"}, {ID: "x"}}

	err := validator.Validate(
		context.Background(),
		validation.SliceProperty("items", items, it.HasUniqueValuesBy(func(x itemWithID) string { return x.ID })),
	)

	validationtest.Assert(t, err).IsViolationList().WithLen(2)
	list, _ := validation.UnwrapViolations(err)
	paths := make([]string, 0, 2)
	list.ForEach(func(i int, v validation.Violation) error {
		paths = append(paths, v.PropertyPath().String())
		return nil
	})
	assert.Contains(t, paths, "items[0]")
	assert.Contains(t, paths, "items[1]")
}
