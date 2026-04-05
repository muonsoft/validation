package test

import (
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/message"
	"github.com/muonsoft/validation/validate"
)

var suspiciousCharactersConstraintTestCases = []ConstraintValidationTestCase{
	{
		name:            "NoSuspiciousCharacters passes on nil",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.NoSuspiciousCharacters(),
		assert:          assertNoError,
	},
	{
		name:            "NoSuspiciousCharacters passes on empty",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue(""),
		constraint:      it.NoSuspiciousCharacters(),
		assert:          assertNoError,
	},
	{
		name:            "NoSuspiciousCharacters passes on ASCII",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("symfony"),
		constraint:      it.NoSuspiciousCharacters(),
		assert:          assertNoError,
	},
	{
		name:            "NoSuspiciousCharacters violation on zero-width space",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("a\u200b"),
		constraint:      it.NoSuspiciousCharacters(),
		assert:          assertHasOneViolation(validation.ErrSuspiciousInvisible, message.SuspiciousInvisible),
	},
	{
		name:            "NoSuspiciousCharacters violation on mixed decimal digits",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("8৪"),
		constraint:      it.NoSuspiciousCharacters(),
		assert:          assertHasOneViolation(validation.ErrSuspiciousMixedNumbers, message.SuspiciousMixedNumbers),
	},
	{
		name:            "NoSuspiciousCharacters violation on hidden overlay",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("i\u0307"),
		constraint:      it.NoSuspiciousCharacters(),
		assert:          assertHasOneViolation(validation.ErrSuspiciousHiddenOverlay, message.SuspiciousHiddenOverlay),
	},
	{
		name:            "NoSuspiciousCharacters single-script restriction",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("a\u0430"),
		constraint: it.NoSuspiciousCharacters().
			WithSuspiciousRestriction(validate.SuspiciousRestrictionSingleScript),
		assert: assertHasOneViolation(validation.ErrSuspiciousCharactersRestriction, message.SuspiciousCharactersRestriction),
	},
	{
		name:            "NoSuspiciousCharacters locale restriction",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("πει"),
		constraint: it.NoSuspiciousCharacters().
			WithSuspiciousRestriction(validate.SuspiciousRestrictionLocales).
			WithSuspiciousLocales("en"),
		assert: assertHasOneViolation(validation.ErrSuspiciousCharactersRestriction, message.SuspiciousCharactersRestriction),
	},
	{
		name:            "NoSuspiciousCharacters When false ignores violation",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("a\u200b"),
		constraint:      it.NoSuspiciousCharacters().When(false),
		assert:          assertNoError,
	},
}
