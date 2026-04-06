package validate

import (
	"errors"
	"strings"
	"unicode"

	"golang.org/x/text/language"
)

// ErrInvalidBIC is returned by [BIC] when the value is not a valid Business Identifier Code (BIC / SWIFT).
// Behavior is aligned with Symfony\Component\Validator\Constraints\Bic and BicValidator (country check via
// ISO 3166-1 alpha-2 regions known to [golang.org/x/text/language] plus Symfony's BIC-to-IBAN territory map).
var ErrInvalidBIC = errors.New("invalid BIC")

// ErrBICIBANCountryMismatch is returned by [BIC] when an associated IBAN is set and its country code
// does not match the BIC's country/territory code (using the same territory mapping as Symfony BicValidator).
var ErrBICIBANCountryMismatch = errors.New("BIC country does not match IBAN")

// bicTerritoryToIBANCountry maps BIC country/territory codes to the parent IBAN country code, matching
// Symfony\Component\Validator\Constraints\BicValidator::BIC_COUNTRY_TO_IBAN_COUNTRY_MAP.
var bicTerritoryToIBANCountry = map[string]string{
	"GF": "FR", "PF": "FR", "TF": "FR", "GP": "FR", "MQ": "FR", "YT": "FR",
	"NC": "FR", "RE": "FR", "BL": "FR", "MF": "FR", "PM": "FR", "WF": "FR",
	"JE": "GB", "IM": "GB", "GG": "GB", "VG": "GB",
	"AX": "FI",
	"IC": "ES", "EA": "ES",
}

const (
	bicModeStrict          = "strict"
	bicModeCaseInsensitive = "case-insensitive"
	defaultBICMode         = bicModeStrict
)

// BICOptions configures [BIC] validation.
type BICOptions struct {
	mode string
	iban string
}

// BICCaseInsensitive enables case-insensitive validation (Symfony Bic::VALIDATION_MODE_CASE_INSENSITIVE):
// lowercase letters are allowed, and the BIC country/territory code is matched case-insensitively.
// The default is strict mode (uppercase ASCII only).
func BICCaseInsensitive() func(*BICOptions) {
	return func(o *BICOptions) {
		o.mode = bicModeCaseInsensitive
	}
}

// BICWithIBAN sets an IBAN value to assert that its country code matches the BIC's territory/country
// (same rules as Symfony Bic constraint "iban" option).
func BICWithIBAN(iban string) func(*BICOptions) {
	return func(o *BICOptions) {
		o.iban = iban
	}
}

func newBICOptions() BICOptions {
	return BICOptions{mode: defaultBICMode}
}

// BIC validates whether the value is a valid Business Identifier Code (BIC / SWIFT), aligned with
// Symfony\Component\Validator\Constraints\Bic and BicValidator.
//
// Spaces (U+0020) are stripped before validation, matching Symfony str_replace(' ', ”, $value).
//
// Empty string is considered valid (use [NotBlank] or similar to reject empty values).
//
// Possible errors:
//   - [ErrInvalidBIC] when length, characters, country/territory, or strict-case rules fail;
//   - [ErrBICIBANCountryMismatch] when [BICWithIBAN] is set with a non-empty IBAN whose first two letters
//     (after IBAN canonicalization) form an alpha country code that does not match the BIC territory.
//
// See https://en.wikipedia.org/wiki/ISO_9362.
func BIC(value string, options ...func(*BICOptions)) error {
	if value == "" {
		return nil
	}

	opts := newBICOptions()
	for _, opt := range options {
		opt(&opts)
	}

	s := stripBICSpaces(value)
	if err := bicValidateStructure(s); err != nil {
		return err
	}

	bicCC := bicCountryCodeFromBIC(s, opts.mode)
	if !bicCountryOrTerritoryKnown(bicCC) {
		return ErrInvalidBIC
	}

	if opts.mode == bicModeStrict && strings.ToUpper(s) != s {
		return ErrInvalidBIC
	}

	return bicValidateAgainstIBAN(bicCC, opts.iban)
}

func bicValidateStructure(s string) error {
	if len(s) != 8 && len(s) != 11 {
		return ErrInvalidBIC
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= utf8ASCIIUpperBound || !unicode.IsLetter(rune(c)) && !unicode.IsDigit(rune(c)) {
			return ErrInvalidBIC
		}
	}
	return nil
}

func bicCountryCodeFromBIC(s string, mode string) string {
	cc := s[4:6]
	if mode == bicModeCaseInsensitive {
		return strings.ToUpper(cc)
	}
	return cc
}

func bicValidateAgainstIBAN(bicCC, iban string) error {
	if iban == "" {
		return nil
	}
	ibanCanon, ok := canonicalizeIBAN(iban)
	if !ok || len(ibanCanon) < 2 {
		return nil
	}
	ibanCC := ibanCanon[:2]
	if !isAlpha2(ibanCC) {
		return nil
	}
	if bicMatchesIBANCountry(bicCC, ibanCC) {
		return nil
	}
	return ErrBICIBANCountryMismatch
}

func stripBICSpaces(value string) string {
	if !strings.Contains(value, " ") {
		return value
	}
	return strings.ReplaceAll(value, " ", "")
}

func bicCountryOrTerritoryKnown(bicCC string) bool {
	if len(bicCC) != 2 {
		return false
	}
	if _, ok := bicTerritoryToIBANCountry[bicCC]; ok {
		return true
	}
	r, err := language.ParseRegion(bicCC)
	if err != nil {
		return false
	}
	return r.IsCountry() && !r.IsGroup()
}

func bicMatchesIBANCountry(bicCountryCode, ibanCountryCode string) bool {
	if ibanCountryCode == bicCountryCode {
		return true
	}
	if parent, ok := bicTerritoryToIBANCountry[bicCountryCode]; ok && ibanCountryCode == parent {
		return true
	}
	return false
}
