package test

import (
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/message"
)

var barcodeConstraintsTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsEAN8 passes on valid value",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("42345671"),
		constraint:      it.IsEAN8(),
		assert:          assertNoError,
	},
	{
		name:            "IsEAN8 violation on invalid value",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("42345670"),
		constraint:      it.IsEAN8(),
		assert:          assertHasOneViolation(validation.ErrInvalidEAN8, message.InvalidEAN8),
	},
	{
		name:            "IsEAN13 passes on valid value",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("4719512002889"),
		constraint:      it.IsEAN13(),
		assert:          assertNoError,
	},
	{
		name:            "IsEAN13 violation on invalid value",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("4006381333932"),
		constraint:      it.IsEAN13(),
		assert:          assertHasOneViolation(validation.ErrInvalidEAN13, message.InvalidEAN13),
	},
	{
		name:            "IsUPCA passes on valid value",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("614141000036"),
		constraint:      it.IsUPCA(),
		assert:          assertNoError,
	},
	{
		name:            "IsUPCA violation on invalid value",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("614141000037"),
		constraint:      it.IsUPCA(),
		assert:          assertHasOneViolation(validation.ErrInvalidUPCA, message.InvalidUPCA),
	},
	{
		name:            "IsUPCE passes on valid value",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("01234505"),
		constraint:      it.IsUPCE(),
		assert:          assertNoError,
	},
	{
		name:            "IsUPCE violation on invalid value",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("01234501"),
		constraint:      it.IsUPCE(),
		assert:          assertHasOneViolation(validation.ErrInvalidUPCE, message.InvalidUPCE),
	},
}
