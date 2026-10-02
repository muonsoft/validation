package it

import (
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/is"
)

// IsCountry checks for a valid two-letter country or territory code.
// Nil and empty values are ignored; use [IsNotBlank] to require a value.
// Case is ignored and whitespace is not trimmed. The input is not modified.
// Uses x/text country data, including CLDR additions and historical codes;
// acceptance does not imply service availability.
// See [github.com/muonsoft/validation/validate.Country] for the full contract.
func IsCountry() validation.StringFuncConstraint {
	return validation.OfStringBy(is.Country).
		WithError(validation.ErrInvalidCountry).
		WithMessage(validation.ErrInvalidCountry.Message())
}

// IsLanguage checks for a valid standalone two- or three-letter language code.
// Nil and empty values are ignored; use [IsNotBlank] to require a value.
// Case is ignored and whitespace is not trimmed. The input is not modified.
// Composite tags require [IsLocale]. Special and private-use codes recognized
// by x/text are accepted; acceptance does not imply translation availability.
// See [github.com/muonsoft/validation/validate.Language] for the full contract.
func IsLanguage() validation.StringFuncConstraint {
	return validation.OfStringBy(is.Language).
		WithError(validation.ErrInvalidLanguage).
		WithMessage(validation.ErrInvalidLanguage.Message())
}

// IsLocale checks for a valid locale identifier.
// Nil and empty values are ignored; use [IsNotBlank] to require a value.
// Case is ignored and whitespace is not trimmed. The input is not modified.
// Accepts composite tags, underscore separators, special and private-use tags,
// and extensions supported by x/text; acceptance does not imply application support.
// See [github.com/muonsoft/validation/validate.Locale] for the full contract.
func IsLocale() validation.StringFuncConstraint {
	return validation.OfStringBy(is.Locale).
		WithError(validation.ErrInvalidLocale).
		WithMessage(validation.ErrInvalidLocale.Message())
}
