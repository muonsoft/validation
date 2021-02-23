package internal

import (
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/code"
	"github.com/muonsoft/validation/it"
)

var countTestCases = []ValidateTestCase{
	{
		name:            "HasMinCount violation on nil",
		isApplicableFor: specificValueTypes(iterableType, countableType),
		options:         []validation.Option{it.HasMinCount(1)},
		assert: assertHasOneViolation(
			code.CountTooFew,
			"This collection should contain 1 element or more.",
			"",
		),
	},
	{
		name:            "HasMinCount violation on nil ignored when condition false",
		isApplicableFor: specificValueTypes(iterableType, countableType),
		options:         []validation.Option{it.HasMinCount(1).When(false)},
		assert:          assertNoError,
	},
	{
		name:            "HasMinCount violation on nil when condition true",
		isApplicableFor: specificValueTypes(iterableType, countableType),
		options:         []validation.Option{it.HasMinCount(1).When(true)},
		assert: assertHasOneViolation(
			code.CountTooFew,
			"This collection should contain 1 element or more.",
			"",
		),
	},
	{
		name:            "HasMinCount violation with custom property path",
		isApplicableFor: specificValueTypes(iterableType, countableType),
		options: []validation.Option{
			validation.PropertyName("properties"),
			validation.ArrayIndex(0),
			validation.PropertyName("value"),
			it.HasMinCount(1),
		},
		assert: assertHasOneViolation(
			code.CountTooFew,
			"This collection should contain 1 element or more.",
			customPath,
		),
	},
	{
		name:            "HasMinCount violation with custom message",
		isApplicableFor: specificValueTypes(iterableType, countableType),
		options: []validation.Option{
			it.HasMinCount(1).MinMessage(customMessage),
		},
		assert: assertHasOneViolation(code.CountTooFew, customMessage, ""),
	},
	{
		name:            "HasMinCount violation on small collection",
		isApplicableFor: specificValueTypes(iterableType, countableType),
		sliceValue:      []string{"a", "b"},
		mapValue:        map[string]string{"a": "a", "b": "b"},
		options:         []validation.Option{it.HasMinCount(3)},
		assert: assertHasOneViolation(
			code.CountTooFew,
			"This collection should contain 3 elements or more.",
			"",
		),
	},
	{
		name:            "HasMinCount passes on equal collection",
		isApplicableFor: specificValueTypes(iterableType, countableType),
		sliceValue:      []string{"a", "b"},
		mapValue:        map[string]string{"a": "a", "b": "b"},
		options:         []validation.Option{it.HasMinCount(2)},
		assert:          assertNoError,
	},
	{
		name:            "HasMaxCount violation on max",
		isApplicableFor: specificValueTypes(iterableType, countableType),
		sliceValue:      []string{"a", "b"},
		mapValue:        map[string]string{"a": "a", "b": "b"},
		options:         []validation.Option{it.HasMaxCount(1)},
		assert: assertHasOneViolation(
			code.CountTooMany,
			"This collection should contain 1 element or less.",
			"",
		),
	},
	{
		name:            "HasMaxCount passes on equal collection",
		isApplicableFor: specificValueTypes(iterableType, countableType),
		sliceValue:      []string{"a", "b"},
		mapValue:        map[string]string{"a": "a", "b": "b"},
		options:         []validation.Option{it.HasMaxCount(2)},
		assert:          assertNoError,
	},
	{
		name:            "HasMaxCount violation on max with custom message",
		isApplicableFor: specificValueTypes(iterableType, countableType),
		sliceValue:      []string{"a", "b"},
		mapValue:        map[string]string{"a": "a", "b": "b"},
		options:         []validation.Option{it.HasMaxCount(1).MaxMessage(customMessage)},
		assert:          assertHasOneViolation(code.CountTooMany, customMessage, ""),
	},
	{
		name:            "HasCountBetween passes on valid collection",
		isApplicableFor: specificValueTypes(iterableType, countableType),
		sliceValue:      []string{"a", "b"},
		mapValue:        map[string]string{"a": "a", "b": "b"},
		options:         []validation.Option{it.HasCountBetween(1, 3)},
		assert:          assertNoError,
	},
	{
		name:            "HasExactCount passes on valid collection",
		isApplicableFor: specificValueTypes(iterableType, countableType),
		sliceValue:      []string{"a", "b"},
		mapValue:        map[string]string{"a": "a", "b": "b"},
		options:         []validation.Option{it.HasExactCount(2)},
		assert:          assertNoError,
	},
	{
		name:            "HasExactCount violation on nil with exact message",
		isApplicableFor: specificValueTypes(iterableType, countableType),
		options:         []validation.Option{it.HasExactCount(1)},
		assert: assertHasOneViolation(
			code.CountExact,
			"This collection should contain exactly 1 element.",
			"",
		),
	},
	{
		name:            "HasExactCount violation on nil with custom message",
		isApplicableFor: specificValueTypes(iterableType, countableType),
		options:         []validation.Option{it.HasExactCount(1).ExactMessage(customMessage)},
		assert:          assertHasOneViolation(code.CountExact, customMessage, ""),
	},
}
