package it

import (
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/is"
)

// IsEAN8 is used to validate EAN-8 value.
//
// See https://en.wikipedia.org/wiki/EAN-8.
func IsEAN8() validation.StringFuncConstraint {
	return validation.OfStringBy(is.EAN8).
		WithError(validation.ErrInvalidEAN8).
		WithMessage(validation.ErrInvalidEAN8.Message())
}

// IsEAN13 is used to validate EAN-13 value.
//
// See https://en.wikipedia.org/wiki/International_Article_Number.
func IsEAN13() validation.StringFuncConstraint {
	return validation.OfStringBy(is.EAN13).
		WithError(validation.ErrInvalidEAN13).
		WithMessage(validation.ErrInvalidEAN13.Message())
}

// IsUPCA is used to validate UPC-A value.
//
// See https://en.wikipedia.org/wiki/Universal_Product_Code.
func IsUPCA() validation.StringFuncConstraint {
	return validation.OfStringBy(is.UPCA).
		WithError(validation.ErrInvalidUPCA).
		WithMessage(validation.ErrInvalidUPCA.Message())
}

// IsUPCE is used to validate UPC-E value.
//
// See https://en.wikipedia.org/wiki/Universal_Product_Code#UPC-E.
func IsUPCE() validation.StringFuncConstraint {
	return validation.OfStringBy(is.UPCE).
		WithError(validation.ErrInvalidUPCE).
		WithMessage(validation.ErrInvalidUPCE.Message())
}
