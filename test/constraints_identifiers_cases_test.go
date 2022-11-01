package test

import (
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/message"
)

var identifierConstraintsTestCases = []ConstraintValidationTestCase{
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
