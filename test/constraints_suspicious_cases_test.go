package test

import (
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/message"
	"github.com/muonsoft/validation/validate"
)

var suspiciousCharactersConstraintTestCases = []ConstraintValidationTestCase{
	{
		name:            "HasNoSuspiciousCharacters passes on nil",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.HasNoSuspiciousCharacters(),
		assert:          assertNoError,
	},
	{
		name:            "HasNoSuspiciousCharacters passes on empty",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue(""),
		constraint:      it.HasNoSuspiciousCharacters(),
		assert:          assertNoError,
	},
	{
		name:            "HasNoSuspiciousCharacters passes on ASCII",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("symfony"),
		constraint:      it.HasNoSuspiciousCharacters(),
		assert:          assertNoError,
	},
	{
		name:            "HasNoSuspiciousCharacters violation on zero-width space",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("a\u200b"),
		constraint:      it.HasNoSuspiciousCharacters(),
		assert:          assertHasOneViolation(validation.ErrSuspiciousInvisible, message.SuspiciousInvisible),
	},
	{
		name:            "HasNoSuspiciousCharacters violation on mixed decimal digits",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("8৪"),
		constraint:      it.HasNoSuspiciousCharacters(),
		assert:          assertHasOneViolation(validation.ErrSuspiciousMixedNumbers, message.SuspiciousMixedNumbers),
	},
	{
		name:            "HasNoSuspiciousCharacters violation on hidden overlay",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("i\u0307"),
		constraint:      it.HasNoSuspiciousCharacters(),
		assert:          assertHasOneViolation(validation.ErrSuspiciousHiddenOverlay, message.SuspiciousHiddenOverlay),
	},
	{
		name:            "HasNoSuspiciousCharacters single-script restriction",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("a\u0430"),
		constraint: it.HasNoSuspiciousCharacters().
			WithSuspiciousRestriction(validate.SuspiciousRestrictionSingleScript),
		assert: assertHasOneViolation(validation.ErrSuspiciousCharactersRestriction, message.SuspiciousCharactersRestriction),
	},
	{
		name:            "HasNoSuspiciousCharacters locale restriction",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("πει"),
		constraint: it.HasNoSuspiciousCharacters().
			WithSuspiciousRestriction(validate.SuspiciousRestrictionLocales).
			WithSuspiciousLocales("en"),
		assert: assertHasOneViolation(validation.ErrSuspiciousCharactersRestriction, message.SuspiciousCharactersRestriction),
	},
	{
		name:            "HasNoSuspiciousCharacters When false ignores violation",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("a\u200b"),
		constraint:      it.HasNoSuspiciousCharacters().When(false),
		assert:          assertNoError,
	},
}
