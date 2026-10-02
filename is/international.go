package is

import "github.com/muonsoft/validation/validate"

// Country reports whether value is a valid two-letter country or territory code.
// Empty strings pass. Case is ignored and whitespace is not trimmed.
// Uses x/text country data, including CLDR additions and historical codes;
// acceptance does not imply service availability.
// See [validate.Country] for the full contract and data-version limitations.
func Country(value string) bool {
	return validate.Country(value) == nil
}

// Language reports whether value is a valid standalone two- or three-letter language code.
// Empty strings pass. Case is ignored and whitespace is not trimmed.
// Composite tags require [Locale]. Special and private-use codes recognized
// by x/text are accepted; acceptance does not imply translation availability.
// See [validate.Language] for the full contract and data-version limitations.
func Language(value string) bool {
	return validate.Language(value) == nil
}

// Locale reports whether value is a valid locale identifier.
// Empty strings pass. Case is ignored and whitespace is not trimmed.
// Accepts composite tags, underscore separators, special and private-use tags,
// and extensions supported by x/text; acceptance does not imply application support.
// See [validate.Locale] for the full contract and data-version limitations.
func Locale(value string) bool {
	return validate.Locale(value) == nil
}
