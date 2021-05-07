package test

import (
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/code"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/message"
)

var urlConstraintTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsURL passes on nil",
		isApplicableFor: specificValueTypes(stringType),
		options:         []validation.Option{it.IsURL()},
		assert:          assertNoError,
	},
	{
		name:            "IsURL passes on empty value",
		isApplicableFor: specificValueTypes(stringType),
		options:         []validation.Option{it.IsURL()},
		stringValue:     stringValue(""),
		assert:          assertNoError,
	},
	{
		name:            "IsURL passes on valid URL",
		isApplicableFor: specificValueTypes(stringType),
		options:         []validation.Option{it.IsURL()},
		stringValue:     stringValue("http://example.com"),
		assert:          assertNoError,
	},
	{
		name:            "IsURL violation on invalid URL",
		isApplicableFor: specificValueTypes(stringType),
		options:         []validation.Option{it.IsURL()},
		stringValue:     stringValue("example.com"),
		assert:          assertHasOneViolation(code.InvalidURL, message.InvalidURL, ""),
	},
	{
		name:            "IsURL error on empty protocols",
		isApplicableFor: specificValueTypes(stringType),
		options:         []validation.Option{it.IsURL().Protocols()},
		stringValue:     stringValue(""),
		assert:          assertError(`failed to set up constraint "URLConstraint": empty list of protocols`),
	},
	{
		name:            "IsURL passes on valid URL with custom protocol",
		isApplicableFor: specificValueTypes(stringType),
		options:         []validation.Option{it.IsURL().Protocols("ftp")},
		stringValue:     stringValue("ftp://example.com"),
		assert:          assertNoError,
	},
	{
		name:            "IsRelativeURL passes on valid relative URL",
		isApplicableFor: specificValueTypes(stringType),
		options:         []validation.Option{it.IsRelativeURL()},
		stringValue:     stringValue("//example.com"),
		assert:          assertNoError,
	},
	{
		name:            "IsRelativeURL passes on valid absolute URL",
		isApplicableFor: specificValueTypes(stringType),
		options:         []validation.Option{it.IsRelativeURL()},
		stringValue:     stringValue("https://example.com"),
		assert:          assertNoError,
	},
	{
		name:            "IsURL violation on invalid URL with custom message",
		isApplicableFor: specificValueTypes(stringType),
		options:         []validation.Option{it.IsURL().Message(`Unexpected URL "{{ value }}"`)},
		stringValue:     stringValue("example.com"),
		assert:          assertHasOneViolation(code.InvalidURL, `Unexpected URL "example.com"`, ""),
	},
	{
		name:            "IsURL passes when condition is false",
		isApplicableFor: specificValueTypes(stringType),
		options:         []validation.Option{it.IsURL().When(false)},
		stringValue:     stringValue("example.com"),
		assert:          assertNoError,
	},
	{
		name:            "IsURL violation when condition is true",
		isApplicableFor: specificValueTypes(stringType),
		options:         []validation.Option{it.IsURL().When(true)},
		stringValue:     stringValue("example.com"),
		assert:          assertHasOneViolation(code.InvalidURL, message.InvalidURL, ""),
	},
}

var emailConstraintTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsEmail passes on valid email",
		isApplicableFor: specificValueTypes(stringType),
		options:         []validation.Option{it.IsEmail()},
		stringValue:     stringValue("user@example.com"),
		assert:          assertNoError,
	},
	{
		name:            "IsEmail violation on invalid email",
		isApplicableFor: specificValueTypes(stringType),
		options:         []validation.Option{it.IsEmail()},
		stringValue:     stringValue("invalid"),
		assert:          assertHasOneViolation(code.InvalidEmail, message.InvalidEmail, ""),
	},
}
