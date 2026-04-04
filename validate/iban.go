package validate

import (
	"errors"
	"strconv"
	"strings"
)

// ErrInvalidIBAN is returned by [IBAN] when the value is not a valid International Bank Account Number.
// Behavior is aligned with Symfony\Component\Validator\Constraints\Iban and IbanValidator.
var ErrInvalidIBAN = errors.New("invalid IBAN")

// IBAN validates whether the value is a valid International Bank Account Number.
// Spaces (including U+00A0 and U+202F), ASCII letters, and digits are accepted; letters are normalized to upper case.
//
// Empty string is considered valid (use [NotBlank] or similar to reject empty values).
//
// Possible errors:
//   - [ErrInvalidIBAN] when the value fails any IBAN rule (invalid characters, country, format, check digits, or mod-97).
//
// See https://en.wikipedia.org/wiki/International_Bank_Account_Number.
func IBAN(value string) error {
	if value == "" {
		return nil
	}

	s, ok := canonicalizeIBAN(value)
	if !ok {
		return ErrInvalidIBAN
	}

	return validateCanonicalIBAN(s)
}

func validateCanonicalIBAN(s string) error {
	if len(s) < 4 {
		return ErrInvalidIBAN
	}
	if ibanHasNonAlphanumeric(s) {
		return ErrInvalidIBAN
	}

	cc := s[:2]
	if !isAlpha2(cc) {
		return ErrInvalidIBAN
	}

	pat, ok := ibanCountryPatterns[cc]
	if !ok || pat == nil || !pat.MatchString(s) {
		return ErrInvalidIBAN
	}

	if err := ibanCheckDigitRange(s[2:4]); err != nil {
		return err
	}

	if ibanMod97(s) != 1 {
		return ErrInvalidIBAN
	}

	return nil
}

func ibanHasNonAlphanumeric(s string) bool {
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' {
			continue
		}
		return true
	}
	return false
}

func ibanCheckDigitRange(two string) error {
	cd, err := strconv.Atoi(two)
	if err != nil || cd < 2 || cd > 98 {
		return ErrInvalidIBAN
	}
	return nil
}

func isAlpha2(cc string) bool {
	if len(cc) != 2 {
		return false
	}
	return cc[0] >= 'A' && cc[0] <= 'Z' && cc[1] >= 'A' && cc[1] <= 'Z'
}

// canonicalizeIBAN strips IBAN grouping spaces and returns upper-case ASCII; ok is false on disallowed bytes/UTF-8.
func canonicalizeIBAN(value string) (string, bool) {
	var b strings.Builder
	b.Grow(len(value))

	for i := 0; i < len(value); {
		c := value[i]
		if c == ' ' {
			i++
			continue
		}
		if n := ibanNBSPPrefix(value, i); n > 0 {
			i += n
			continue
		}
		if n := ibanNNBSPPrefix(value, i); n > 0 {
			i += n
			continue
		}
		if c >= utf8ASCIIUpperBound {
			return "", false
		}
		if c >= 'a' && c <= 'z' {
			c -= 'a' - 'A'
		}
		b.WriteByte(c)
		i++
	}

	return b.String(), true
}

const utf8ASCIIUpperBound = 0x80

func ibanNBSPPrefix(value string, i int) int {
	if value[i] == 0xc2 && i+1 < len(value) && value[i+1] == 0xa0 {
		return 2
	}
	return 0
}

func ibanNNBSPPrefix(value string, i int) int {
	if i+3 <= len(value) && value[i] == 0xe2 && value[i+1] == 0x80 && value[i+2] == 0xaf {
		return 3
	}
	return 0
}

func ibanExpandedNumeric(s string) string {
	rearr := s[4:] + s[:4]
	var b strings.Builder
	b.Grow(len(rearr) * 2)
	for i := 0; i < len(rearr); i++ {
		c := rearr[i]
		switch {
		case c >= '0' && c <= '9':
			b.WriteByte(c)
		case c >= 'A' && c <= 'Z':
			b.WriteString(strconv.Itoa(int(c - 'A' + 10)))
		default:
			return ""
		}
	}
	return b.String()
}

func ibanMod97(s string) int {
	digits := ibanExpandedNumeric(s)
	if digits == "" {
		return 0
	}
	rest := 0
	for start := 0; start < len(digits); start += 7 {
		end := start + 7
		if end > len(digits) {
			end = len(digits)
		}
		for j := start; j < end; j++ {
			rest = (rest*10 + int(digits[j]-'0')) % 97
		}
	}
	return rest
}
