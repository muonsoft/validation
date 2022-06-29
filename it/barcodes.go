package it

import (
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/is"
)

// IsEAN8 is used to validate EAN-8 value.
//
// See https://en.wikipedia.org/wiki/EAN-8.
func IsEAN8() validation.CustomStringConstraint {
	return validation.NewCustomStringConstraint(is.EAN8).
		WithError(validation.ErrInvalidEAN8).
		WithMessage(validation.ErrInvalidEAN8.Template())
}

// IsEAN13 is used to validate EAN-13 value.
//
// See https://en.wikipedia.org/wiki/International_Article_Number.
func IsEAN13() validation.CustomStringConstraint {
	return validation.NewCustomStringConstraint(is.EAN13).
		WithError(validation.ErrInvalidEAN13).
		WithMessage(validation.ErrInvalidEAN13.Template())
}

// IsUPCA is used to validate UPC-A value.
//
// See https://en.wikipedia.org/wiki/Universal_Product_Code.
func IsUPCA() validation.CustomStringConstraint {
	return validation.NewCustomStringConstraint(is.UPCA).
		WithError(validation.ErrInvalidUPCA).
		WithMessage(validation.ErrInvalidUPCA.Template())
}

// IsUPCE is used to validate UPC-E value.
//
// See https://en.wikipedia.org/wiki/Universal_Product_Code#UPC-E.
func IsUPCE() validation.CustomStringConstraint {
	return validation.NewCustomStringConstraint(is.UPCE).
		WithError(validation.ErrInvalidUPCE).
		WithMessage(validation.ErrInvalidUPCE.Template())
}
