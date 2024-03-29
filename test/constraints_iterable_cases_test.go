package test

import (
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
)

var countConstraintTestCases = []ConstraintValidationTestCase{
	{
		name:            "HasMinCount violation on nil",
		isApplicableFor: specificValueTypes(iterableType, countableType),
		constraint:      it.HasMinCount(1),
		assert: assertHasOneViolation(
			validation.ErrTooFewElements,
			"This collection should contain 1 element or more.",
		),
	},
	{
		name:            "HasMinCount violation on nil ignored when condition false",
		isApplicableFor: specificValueTypes(iterableType, countableType),
		constraint:      it.HasMinCount(1).When(false),
		assert:          assertNoError,
	},
	{
		name:            "HasMinCount violation on nil ignored when groups not match",
		isApplicableFor: specificValueTypes(iterableType, countableType),
		constraint:      it.HasMinCount(1).WhenGroups(testGroup),
		assert:          assertNoError,
	},
	{
		name:            "HasMinCount violation on nil when condition true",
		isApplicableFor: specificValueTypes(iterableType, countableType),
		constraint:      it.HasMinCount(1).When(true),
		assert: assertHasOneViolation(
			validation.ErrTooFewElements,
			"This collection should contain 1 element or more.",
		),
	},
	{
		name:            "HasMinCount violation with custom message",
		isApplicableFor: specificValueTypes(iterableType, countableType),
		constraint: it.HasMinCount(1).
			WithMinError(ErrMin).
			WithMinMessage(
				"Unexpected count {{ count }} at {{ custom }}, should not be less than {{ limit }}.",
				validation.TemplateParameter{Key: "{{ custom }}", Value: "parameter"},
			),
		assert: assertHasOneViolation(
			ErrMin,
			"Unexpected count 0 at parameter, should not be less than 1.",
		),
	},
	{
		name:            "HasMinCount violation on small collection",
		isApplicableFor: specificValueTypes(iterableType, countableType),
		sliceValue:      []string{"a", "b"},
		mapValue:        map[string]string{"a": "a", "b": "b"},
		constraint:      it.HasMinCount(3),
		assert: assertHasOneViolation(
			validation.ErrTooFewElements,
			"This collection should contain 3 elements or more.",
		),
	},
	{
		name:            "HasMinCount passes on equal collection",
		isApplicableFor: specificValueTypes(iterableType, countableType),
		sliceValue:      []string{"a", "b"},
		mapValue:        map[string]string{"a": "a", "b": "b"},
		constraint:      it.HasMinCount(2),
		assert:          assertNoError,
	},
	{
		name:            "HasMaxCount violation on max",
		isApplicableFor: specificValueTypes(iterableType, countableType),
		sliceValue:      []string{"a", "b"},
		mapValue:        map[string]string{"a": "a", "b": "b"},
		constraint:      it.HasMaxCount(1),
		assert: assertHasOneViolation(
			validation.ErrTooManyElements,
			"This collection should contain 1 element or less.",
		),
	},
	{
		name:            "HasMaxCount passes on equal collection",
		isApplicableFor: specificValueTypes(iterableType, countableType),
		sliceValue:      []string{"a", "b"},
		mapValue:        map[string]string{"a": "a", "b": "b"},
		constraint:      it.HasMaxCount(2),
		assert:          assertNoError,
	},
	{
		name:            "HasMaxCount violation on max with custom message",
		isApplicableFor: specificValueTypes(iterableType, countableType),
		sliceValue:      []string{"a", "b"},
		mapValue:        map[string]string{"a": "a", "b": "b"},
		constraint: it.HasMaxCount(1).
			WithMaxError(ErrMax).
			WithMaxMessage(
				"Unexpected count {{ count }} at {{ custom }}, should not be greater than {{ limit }}.",
				validation.TemplateParameter{Key: "{{ custom }}", Value: "parameter"},
			),
		assert: assertHasOneViolation(
			ErrMax,
			"Unexpected count 2 at parameter, should not be greater than 1.",
		),
	},
	{
		name:            "HasCountBetween passes on valid collection",
		isApplicableFor: specificValueTypes(iterableType, countableType),
		sliceValue:      []string{"a", "b"},
		mapValue:        map[string]string{"a": "a", "b": "b"},
		constraint:      it.HasCountBetween(1, 3),
		assert:          assertNoError,
	},
	{
		name:            "HasExactCount passes on valid collection",
		isApplicableFor: specificValueTypes(iterableType, countableType),
		sliceValue:      []string{"a", "b"},
		mapValue:        map[string]string{"a": "a", "b": "b"},
		constraint:      it.HasExactCount(2),
		assert:          assertNoError,
	},
	{
		name:            "HasExactCount violation on nil with exact message",
		isApplicableFor: specificValueTypes(iterableType, countableType),
		constraint:      it.HasExactCount(1),
		assert: assertHasOneViolation(
			validation.ErrNotExactCount,
			"This collection should contain exactly 1 element.",
		),
	},
	{
		name:            "HasExactCount violation on nil with custom message",
		isApplicableFor: specificValueTypes(iterableType, countableType),
		constraint: it.HasExactCount(1).
			WithExactError(ErrExact).
			WithExactMessage(
				"Unexpected count {{ count }} at {{ custom }}, should be exactly {{ limit }}.",
				validation.TemplateParameter{Key: "{{ custom }}", Value: "parameter"},
			),
		assert: assertHasOneViolation(
			ErrExact,
			"Unexpected count 0 at parameter, should be exactly 1.",
		),
	},
	{
		name:            "HasCountDivisibleBy passes on valid collection",
		isApplicableFor: specificValueTypes(iterableType, countableType),
		sliceValue:      []string{"a", "b"},
		mapValue:        map[string]string{"a": "a", "b": "b"},
		constraint:      it.HasCountDivisibleBy(2),
		assert:          assertNoError,
	},
	{
		name:            "HasCountDivisibleBy violation on invalid collection",
		isApplicableFor: specificValueTypes(iterableType, countableType),
		sliceValue:      []string{"a", "b"},
		mapValue:        map[string]string{"a": "a", "b": "b"},
		constraint:      it.HasCountDivisibleBy(3),
		assert: assertHasOneViolation(
			validation.ErrNotDivisibleCount,
			"The number of elements in this collection should be a multiple of 3.",
		),
	},
	{
		name:            "HasCountDivisibleBy violation with custom message",
		isApplicableFor: specificValueTypes(iterableType, countableType),
		sliceValue:      []string{"a", "b"},
		mapValue:        map[string]string{"a": "a", "b": "b"},
		constraint: it.HasCountDivisibleBy(3).
			WithDivisibleError(ErrCustom).
			WithDivisibleMessage(
				"Unexpected count {{ count }} at {{ custom }}, should be divisible by {{ divisibleBy }}.",
				validation.TemplateParameter{Key: "{{ custom }}", Value: "parameter"},
			),
		assert: assertHasOneViolation(
			ErrCustom,
			"Unexpected count 2 at parameter, should be divisible by 3.",
		),
	},
}
