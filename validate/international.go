package validate

import (
	"errors"

	"golang.org/x/text/language"
)

var (
	// ErrInvalidCountry indicates an invalid two-letter country or territory code.
	ErrInvalidCountry = errors.New("invalid country")
	// ErrInvalidLanguage indicates an invalid standalone language code.
	ErrInvalidLanguage = errors.New("invalid language")
	// ErrInvalidLocale indicates an invalid locale identifier.
	ErrInvalidLocale = errors.New("invalid locale")
)

// Country validates a two-letter country or territory code using
// [language.ParseRegion] and [language.Region.IsCountry]. It returns nil for an
// empty or valid value, otherwise [ErrInvalidCountry]. Case is ignored;
// whitespace, three-letter codes, numeric regions, and region groups are rejected.
//
// Use it for country fields in addresses or APIs. The data comes from x/text,
// including CLDR additions and recognized historical codes, rather than a strict
// list of currently assigned ISO 3166-1 codes. For example, XK, UK, SU, and
// the special code UN pass; EU and ZZ do not. Acceptance does not imply delivery
// or service availability.
// Values are not modified; recognized codes may change with the x/text version.
func Country(value string) error {
	if value == "" {
		return nil
	}
	if len(value) != 2 {
		return ErrInvalidCountry
	}
	region, err := language.ParseRegion(value)
	if err != nil || !region.IsCountry() || region.IsGroup() {
		return ErrInvalidCountry
	}
	return nil
}

// Language validates a standalone two- or three-letter language code using
// [language.ParseBase]. It returns nil for an empty or valid value, otherwise
// [ErrInvalidLanguage]. Case is ignored; whitespace and composite tags such as
// en-US or zh-Hant are rejected. Use [Locale] for composite tags.
//
// Use it for language metadata, not to assert that a translation is available.
// All codes recognized by x/text are accepted, including legacy codes, special
// codes such as und (undetermined), mul (multiple languages), zxx (no linguistic
// content), and the reserved private-use range qaa-qtz. It does not restrict input
// to living languages. Values are not modified; data follows the x/text version.
func Language(value string) error {
	if value == "" {
		return nil
	}
	if _, err := language.ParseBase(value); err != nil {
		return ErrInvalidLanguage
	}
	return nil
}

// Locale validates a locale identifier using [language.Parse]. It returns nil
// for an empty value or a successful parse, otherwise [ErrInvalidLocale], even
// if the parser recovered a partial tag. It accepts BCP 47 tags and the Unicode
// locale forms supported by x/text, including underscore separators. Case is
// ignored; surrounding whitespace is not trimmed. Values are not modified.
//
// Use it for locale preferences such as en, en-US, or zh-Hant-TW. A region is
// optional; numeric regions, recognized legacy tags, und, root, private-use tags,
// and extensions are accepted. Unlike [Language], it accepts composite tags;
// unlike [Country], its region need not identify a country.
//
// Acceptance does not guarantee translations, formatting resources, or meaningful
// combinations of subtags. Extension values follow the parser's checks, not an
// application allowlist (for example, en-u-ca-foobar passes). Accepted identifiers
// follow the x/text version, not Symfony's locale list. Use a separate allowlist
// to enforce application support.
func Locale(value string) error {
	if value == "" {
		return nil
	}
	if _, err := language.Parse(value); err != nil {
		return ErrInvalidLocale
	}
	return nil
}
