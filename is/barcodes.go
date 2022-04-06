package is

import "github.com/muonsoft/validation/validate"

// EAN8 checks that string contains valid EAN-8 code.
//
// See https://en.wikipedia.org/wiki/EAN-8.
func EAN8(value string) bool {
	return validate.EAN8(value) == nil
}

// EAN13 checks that string contains valid EAN-13 code.
//
// See https://en.wikipedia.org/wiki/International_Article_Number.
func EAN13(value string) bool {
	return validate.EAN13(value) == nil
}

// UPCA checks that string contains valid UPC-A code.
//
// See https://en.wikipedia.org/wiki/Universal_Product_Code.
func UPCA(value string) bool {
	return validate.UPCA(value) == nil
}

// UPCE checks that string contains valid UPC-E code.
//
// See https://en.wikipedia.org/wiki/Universal_Product_Code#UPC-E.
func UPCE(value string) bool {
	return validate.UPCE(value) == nil
}
