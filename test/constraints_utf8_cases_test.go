package test

import (
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/message"
)

var utf8ConstraintTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsUTF8 nil",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsUTF8(),
		assert:          assertNoError,
	},
	{
		name:            "IsUTF8 empty",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsUTF8(),
		assert:          assertNoError,
		stringValue:     stringValue(""),
	},
	{
		name:            "IsUTF8 multilingual",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsUTF8(),
		assert:          assertNoError,
		stringValue:     stringValue("Привет 日本語 🙂"),
	},
	{
		name:            "IsUTF8 replacement character",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsUTF8(),
		assert:          assertNoError,
		stringValue:     stringValue("\ufffd"),
	},
	{
		name:            "IsUTF8 invalid bytes",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsUTF8(),
		assert:          assertHasOneViolation(validation.ErrInvalidUTF8, message.InvalidUTF8),
		stringValue:     stringValue("\xff"),
	},
	{
		name:            "IsUTF8 truncated input",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsUTF8(),
		assert:          assertHasOneViolation(validation.ErrInvalidUTF8, message.InvalidUTF8),
		stringValue:     stringValue("\xe2\x82"),
	},
	{
		name:            "IsUTF8 ignored condition",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsUTF8().When(false),
		assert:          assertNoError,
		stringValue:     stringValue("\xff"),
	},
	{
		name:            "IsUTF8 enabled condition",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsUTF8().When(true),
		assert:          assertHasOneViolation(validation.ErrInvalidUTF8, message.InvalidUTF8),
		stringValue:     stringValue("\xff"),
	},
	{
		name:            "IsUTF8 ignored group",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsUTF8().WhenGroups(testGroup),
		assert:          assertNoError,
		stringValue:     stringValue("\xff"),
	},
	{
		name:            "IsUTF8 custom error and message",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsUTF8().WithError(ErrCustom).WithMessage(customMessage, validation.TemplateParameter{Key: "{{ custom }}", Value: "parameter"}),
		assert:          assertHasOneViolation(ErrCustom, renderedCustomMessage),
		stringValue:     stringValue("\xff"),
	},
}
