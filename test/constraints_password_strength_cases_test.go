package test

import (
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/message"
)

var passwordStrengthConstraintTestCases = []ConstraintValidationTestCase{
	{name: "PasswordStrength nil", isApplicableFor: specificValueTypes(stringType), constraint: it.HasPasswordStrength(), assert: assertNoError},
	{name: "PasswordStrength empty", isApplicableFor: specificValueTypes(stringType), stringValue: stringValue(""), constraint: it.HasPasswordStrength(), assert: assertHasOneViolation(validation.ErrPasswordTooWeak, message.PasswordTooWeak)},
	{name: "PasswordStrength weak", isApplicableFor: specificValueTypes(stringType), stringValue: stringValue("password"), constraint: it.HasPasswordStrength(), assert: assertHasOneViolation(validation.ErrPasswordTooWeak, message.PasswordTooWeak)},
	{name: "PasswordStrength valid", isApplicableFor: specificValueTypes(stringType), stringValue: stringValue("Reasonable-pwd"), constraint: it.HasPasswordStrength(), assert: assertNoError},
}
