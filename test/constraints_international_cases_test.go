package test

import (
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/message"
)

var internationalConstraintTestCases = []ConstraintValidationTestCase{
	{
		name: "IsCountry nil", isApplicableFor: specificValueTypes(stringType),
		constraint: it.IsCountry(), assert: assertNoError,
	},
	{
		name: "IsCountry empty", isApplicableFor: specificValueTypes(stringType),
		constraint: it.IsCountry(), assert: assertNoError,
		stringValue: stringValue(""),
	},
	{
		name: "IsCountry valid", isApplicableFor: specificValueTypes(stringType),
		constraint: it.IsCountry(), assert: assertNoError,
		stringValue: stringValue("de"),
	},
	{
		name: "IsCountry invalid", isApplicableFor: specificValueTypes(stringType),
		constraint: it.IsCountry(), assert: assertHasOneViolation(validation.ErrInvalidCountry, message.InvalidCountry),
		stringValue: stringValue("EU"),
	},
	{
		name: "IsCountry ignored condition", isApplicableFor: specificValueTypes(stringType),
		constraint: it.IsCountry().When(false), assert: assertNoError,
		stringValue: stringValue("EU"),
	},
	{
		name: "IsCountry enabled condition", isApplicableFor: specificValueTypes(stringType),
		constraint: it.IsCountry().When(true), assert: assertHasOneViolation(validation.ErrInvalidCountry, message.InvalidCountry),
		stringValue: stringValue("EU"),
	},
	{
		name: "IsCountry ignored group", isApplicableFor: specificValueTypes(stringType),
		constraint: it.IsCountry().WhenGroups(testGroup), assert: assertNoError,
		stringValue: stringValue("EU"),
	},
	{
		name: "IsCountry custom error and message", isApplicableFor: specificValueTypes(stringType),
		constraint: it.IsCountry().WithError(ErrCustom).WithMessage(customMessage, validation.TemplateParameter{Key: "{{ custom }}", Value: "parameter"}), assert: assertHasOneViolation(ErrCustom, renderedCustomMessage),
		stringValue: stringValue("EU"),
	},
	{
		name: "IsLanguage nil", isApplicableFor: specificValueTypes(stringType),
		constraint: it.IsLanguage(), assert: assertNoError,
	},
	{
		name: "IsLanguage empty", isApplicableFor: specificValueTypes(stringType),
		constraint: it.IsLanguage(), assert: assertNoError,
		stringValue: stringValue(""),
	},
	{
		name: "IsLanguage valid", isApplicableFor: specificValueTypes(stringType),
		constraint: it.IsLanguage(), assert: assertNoError,
		stringValue: stringValue("eng"),
	},
	{
		name: "IsLanguage invalid", isApplicableFor: specificValueTypes(stringType),
		constraint: it.IsLanguage(), assert: assertHasOneViolation(validation.ErrInvalidLanguage, message.InvalidLanguage),
		stringValue: stringValue("en-US"),
	},
	{
		name: "IsLanguage ignored condition", isApplicableFor: specificValueTypes(stringType),
		constraint: it.IsLanguage().When(false), assert: assertNoError,
		stringValue: stringValue("en-US"),
	},
	{
		name: "IsLanguage enabled condition", isApplicableFor: specificValueTypes(stringType),
		constraint: it.IsLanguage().When(true), assert: assertHasOneViolation(validation.ErrInvalidLanguage, message.InvalidLanguage),
		stringValue: stringValue("en-US"),
	},
	{
		name: "IsLanguage ignored group", isApplicableFor: specificValueTypes(stringType),
		constraint: it.IsLanguage().WhenGroups(testGroup), assert: assertNoError,
		stringValue: stringValue("en-US"),
	},
	{
		name: "IsLanguage custom error and message", isApplicableFor: specificValueTypes(stringType),
		constraint: it.IsLanguage().WithError(ErrCustom).WithMessage(customMessage, validation.TemplateParameter{Key: "{{ custom }}", Value: "parameter"}), assert: assertHasOneViolation(ErrCustom, renderedCustomMessage),
		stringValue: stringValue("en-US"),
	},
	{
		name: "IsLocale nil", isApplicableFor: specificValueTypes(stringType),
		constraint: it.IsLocale(), assert: assertNoError,
	},
	{
		name: "IsLocale empty", isApplicableFor: specificValueTypes(stringType),
		constraint: it.IsLocale(), assert: assertNoError,
		stringValue: stringValue(""),
	},
	{
		name: "IsLocale valid", isApplicableFor: specificValueTypes(stringType),
		constraint: it.IsLocale(), assert: assertNoError,
		stringValue: stringValue("zh_Hant_TW"),
	},
	{
		name: "IsLocale invalid", isApplicableFor: specificValueTypes(stringType),
		constraint: it.IsLocale(), assert: assertHasOneViolation(validation.ErrInvalidLocale, message.InvalidLocale),
		stringValue: stringValue("en--US"),
	},
	{
		name: "IsLocale ignored condition", isApplicableFor: specificValueTypes(stringType),
		constraint: it.IsLocale().When(false), assert: assertNoError,
		stringValue: stringValue("en--US"),
	},
	{
		name: "IsLocale enabled condition", isApplicableFor: specificValueTypes(stringType),
		constraint: it.IsLocale().When(true), assert: assertHasOneViolation(validation.ErrInvalidLocale, message.InvalidLocale),
		stringValue: stringValue("en--US"),
	},
	{
		name: "IsLocale ignored group", isApplicableFor: specificValueTypes(stringType),
		constraint: it.IsLocale().WhenGroups(testGroup), assert: assertNoError,
		stringValue: stringValue("en--US"),
	},
	{
		name: "IsLocale custom error and message", isApplicableFor: specificValueTypes(stringType),
		constraint: it.IsLocale().WithError(ErrCustom).WithMessage(customMessage, validation.TemplateParameter{Key: "{{ custom }}", Value: "parameter"}), assert: assertHasOneViolation(ErrCustom, renderedCustomMessage),
		stringValue: stringValue("en--US"),
	},
}
