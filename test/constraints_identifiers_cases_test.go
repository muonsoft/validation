package test

import (
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/message"
)

var identifierConstraintsTestCases = mergeTestCases(
	ulidConstraintTestCases,
	uuidConstraintTestCases,
)

var ulidConstraintTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsULID passes on valid value",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("01ARZ3NDEKTSV4RRFFQ69G5FAV"),
		constraint:      it.IsULID(),
		assert:          assertNoError,
	},
	{
		name:            "IsULID violation on invalid value",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("81ARZ3NDEKTSV4RRFFQ69G5FAV"),
		constraint:      it.IsULID(),
		assert:          assertHasOneViolation(validation.ErrInvalidULID, message.InvalidULID),
	},
}

var uuidConstraintTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsUUID passes on nil",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsUUID(),
		assert:          assertNoError,
	},
	{
		name:            "IsUUID passes on empty value",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsUUID(),
		stringValue:     stringValue(""),
		assert:          assertNoError,
	},
	{
		name:            "IsUUID passes on valid value",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsUUID(),
		stringValue:     stringValue("661eeca0-bc27-4ecc-8f69-6ffb7b1d5a92"),
		assert:          assertNoError,
	},
	{
		name:            "IsUUID violation on invalid value",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsUUID(),
		stringValue:     stringValue("invalid"),
		assert:          assertHasOneViolation(validation.ErrInvalidUUID, message.InvalidUUID),
	},
	{
		name:            "IsUUID passes on nil value",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsUUID(),
		stringValue:     stringValue("00000000-0000-0000-0000-000000000000"),
		assert:          assertNoError,
	},
	{
		name:            "IsUUID violation on nil value when not nil",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsUUID().NotNil(),
		stringValue:     stringValue("00000000-0000-0000-0000-000000000000"),
		assert:          assertHasOneViolation(validation.ErrInvalidUUID, message.InvalidUUID),
	},
	{
		name:            "IsUUID violation on non-canonical value",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsUUID(),
		stringValue:     stringValue("661eeca0bc274ecc8f696ffb7b1d5a92"),
		assert:          assertHasOneViolation(validation.ErrInvalidUUID, message.InvalidUUID),
	},
	{
		name:            "IsUUID passes on non-canonical value when allowed",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsUUID().NonCanonical(),
		stringValue:     stringValue("661eeca0bc274ecc8f696ffb7b1d5a92"),
		assert:          assertNoError,
	},
	{
		name:            "IsUUID violation on not allowed version",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsUUID().WithVersions(1),
		stringValue:     stringValue("661eeca0-bc27-4ecc-8f69-6ffb7b1d5a92"),
		assert:          assertHasOneViolation(validation.ErrInvalidUUID, message.InvalidUUID),
	},
	{
		name:            "IsUUID violation with given error and message",
		isApplicableFor: specificValueTypes(stringType),
		constraint: it.IsUUID().
			WithError(ErrCustom).
			WithMessage(
				`Invalid value "{{ value }}" for {{ custom }}.`,
				validation.TemplateParameter{Key: "{{ custom }}", Value: "parameter"},
			),
		stringValue: stringValue("invalid"),
		assert:      assertHasOneViolation(ErrCustom, `Invalid value "invalid" for parameter.`),
	},
	{
		name:            "IsUUID passes when condition is false",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsUUID().When(false),
		stringValue:     stringValue("invalid"),
		assert:          assertNoError,
	},
	{
		name:            "IsUUID violation when condition is true",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsUUID().When(true),
		stringValue:     stringValue("invalid"),
		assert:          assertHasOneViolation(validation.ErrInvalidUUID, message.InvalidUUID),
	},
	{
		name:            "IsUUID passes when groups not match",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsUUID().WhenGroups(testGroup),
		stringValue:     stringValue("invalid"),
		assert:          assertNoError,
	},
}
