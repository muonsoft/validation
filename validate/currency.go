package validate

import (
	"errors"

	"golang.org/x/text/currency"
)

// ErrInvalidCurrency is returned by [Currency] when the value is not a valid ISO 4217 currency code.
var ErrInvalidCurrency = errors.New("invalid currency")

// Currency validates whether the value is a recognized ISO 4217 alphabetic currency code (three letters).
// Letter case is normalized the same way as [golang.org/x/text/currency.ParseISO] (all upper or all lower).
//
// Empty string is considered valid (use [NotBlank] or similar to reject empty values).
//
// Possible errors:
//   - [ErrInvalidCurrency] when the string is not exactly three letters, is malformed, or is not a known code.
//
// See https://www.iso.org/iso-4217-currency-codes.html and [golang.org/x/text/currency.ParseISO].
func Currency(value string) error {
	if value == "" {
		return nil
	}
	if _, err := currency.ParseISO(value); err != nil {
		return ErrInvalidCurrency
	}
	return nil
}
