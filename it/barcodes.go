package it

import (
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/code"
	"github.com/muonsoft/validation/is"
	"github.com/muonsoft/validation/message"
)

// IsEAN8 is used to validate EAN-8 value.
//
// See https://en.wikipedia.org/wiki/EAN-8.
func IsEAN8() validation.CustomStringConstraint {
	return validation.NewCustomStringConstraint(
		is.EAN8,
		"EAN8Constraint",
		code.InvalidEAN8,
		message.Templates[code.InvalidEAN8],
	)
}

// IsEAN13 is used to validate EAN-13 value.
//
// See https://en.wikipedia.org/wiki/International_Article_Number.
func IsEAN13() validation.CustomStringConstraint {
	return validation.NewCustomStringConstraint(
		is.EAN13,
		"EAN13Constraint",
		code.InvalidEAN13,
		message.Templates[code.InvalidEAN13],
	)
}

// IsUPCA is used to validate UPC-A value.
//
// See https://en.wikipedia.org/wiki/Universal_Product_Code.
func IsUPCA() validation.CustomStringConstraint {
	return validation.NewCustomStringConstraint(
		is.UPCA,
		"UPCAConstraint",
		code.InvalidUPCA,
		message.Templates[code.InvalidUPCA],
	)
}

// IsUPCE is used to validate UPC-E value.
//
// See https://en.wikipedia.org/wiki/Universal_Product_Code#UPC-E.
func IsUPCE() validation.CustomStringConstraint {
	return validation.NewCustomStringConstraint(
		is.UPCE,
		"UPCEConstraint",
		code.InvalidUPCE,
		message.Templates[code.InvalidUPCE],
	)
}
