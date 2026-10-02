package test

import (
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/message"
)

var fileConstraintTestCases = []ConstraintValidationTestCase{
	{name: "FileName nil", isApplicableFor: specificValueTypes(stringType), stringValue: nil, constraint: it.IsFileName(), assert: assertNoError},
	{name: "FileName empty", isApplicableFor: specificValueTypes(stringType), stringValue: stringValue(""), constraint: it.IsFileName(), assert: assertNoError},
	{name: "FileName valid", isApplicableFor: specificValueTypes(stringType), stringValue: stringValue("report.pdf"), constraint: it.IsFileName(), assert: assertNoError},
	{name: "FileName invalid", isApplicableFor: specificValueTypes(stringType), stringValue: stringValue("../report.pdf"), constraint: it.IsFileName(), assert: assertHasOneViolation(validation.ErrInvalidFileName, message.InvalidFileName)},
	{name: "FileExtension nil", isApplicableFor: specificValueTypes(stringType), stringValue: nil, constraint: it.HasFileExtension("pdf"), assert: assertNoError},
	{name: "FileExtension empty", isApplicableFor: specificValueTypes(stringType), stringValue: stringValue(""), constraint: it.HasFileExtension("pdf"), assert: assertNoError},
	{name: "FileExtension valid", isApplicableFor: specificValueTypes(stringType), stringValue: stringValue("report.pdf"), constraint: it.HasFileExtension("pdf"), assert: assertNoError},
	{name: "FileExtension invalid", isApplicableFor: specificValueTypes(stringType), stringValue: stringValue("report.exe"), constraint: it.HasFileExtension("pdf"), assert: assertHasOneViolation(validation.ErrInvalidFileExtension, message.InvalidFileExtension)},
	{name: "MIMEType nil", isApplicableFor: specificValueTypes(stringType), stringValue: nil, constraint: it.IsMIMEType("application/pdf"), assert: assertNoError},
	{name: "MIMEType empty", isApplicableFor: specificValueTypes(stringType), stringValue: stringValue(""), constraint: it.IsMIMEType("application/pdf"), assert: assertNoError},
	{name: "MIMEType valid", isApplicableFor: specificValueTypes(stringType), stringValue: stringValue("application/pdf"), constraint: it.IsMIMEType("application/pdf"), assert: assertNoError},
	{name: "MIMEType invalid", isApplicableFor: specificValueTypes(stringType), stringValue: stringValue("image/png"), constraint: it.IsMIMEType("application/pdf"), assert: assertHasOneViolation(validation.ErrInvalidMIMEType, message.InvalidMIMEType)},
}
