package test

import (
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/message"
	"github.com/muonsoft/validation/validate"
)

var cssColorConstraintTestCases = []ConstraintValidationTestCase{
	{
		name:            "CSSColor nil",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsCSSColor(),
		assert:          assertNoError,
	},
	{
		name:            "CSSColor empty",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue(""),
		constraint:      it.IsCSSColor(),
		assert:          assertNoError,
	},
	{
		name:            "CSSColor valid",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("#AbCd"),
		constraint:      it.IsCSSColor(),
		assert:          assertNoError,
	},
	{
		name:            "CSSColor invalid",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("bad"),
		constraint:      it.IsCSSColor(),
		assert:          assertHasOneViolation(validation.ErrInvalidCSSColor, message.InvalidCSSColor),
	},
	{
		name:            "CSSColor restricted",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("red"),
		constraint:      it.IsCSSColor(validate.CSSColorRGB),
		assert:          assertHasOneViolation(validation.ErrInvalidCSSColor, message.InvalidCSSColor),
	},
	{
		name:            "CSSColor multiple",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("red"),
		constraint:      it.IsCSSColor(validate.CSSColorRGB, validate.CSSColorBasicNamedColors),
		assert:          assertNoError,
	},
	{
		name:            "CSSColor disabled",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("bad"),
		constraint:      it.IsCSSColor().When(false),
		assert:          assertNoError,
	},
	{
		name:            "CSSColor enabled",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("bad"),
		constraint:      it.IsCSSColor().When(true),
		assert:          assertHasOneViolation(validation.ErrInvalidCSSColor, message.InvalidCSSColor),
	},
	{
		name:            "CSSColor groups",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("bad"),
		constraint:      it.IsCSSColor().WhenGroups(testGroup),
		assert:          assertNoError,
	},
	{
		name:            "CSSColor unknown",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsCSSColor("unknown"),
		assert:          assertError("validate by CSSColorConstraint: unknown CSS color format"),
	},
	{
		name:            "CSSColor custom",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("bad"),
		constraint:      it.IsCSSColor().WithError(ErrCustom).WithMessage("{{ value }} {{ custom }}", validation.TemplateParameter{Key: "{{ custom }}", Value: "parameter"}),
		assert:          assertHasOneViolation(ErrCustom, "bad parameter"),
	},
}
