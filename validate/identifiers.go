package validate

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/muonsoft/validation/internal/uuid"
)

var (
	ErrTooShort               = errors.New("too short")
	ErrTooLong                = errors.New("too long")
	ErrTooLarge               = errors.New("too large")
	ErrInvalidCharacters      = errors.New("invalid characters")
	ErrInvalidHyphenPlacement = errors.New("invalid hyphen placement")
	ErrInvalidVersion         = errors.New("invalid version")
	ErrIsNil                  = errors.New("is nil")
)

var ulidChars = newCharSet("0123456789ABCDEFGHJKMNPQRSTVWXYZabcdefghjkmnpqrstvwxyz")

// Same pattern as Symfony Isin::VALIDATION_PATTERN (length is checked separately).
var isinPattern = regexp.MustCompile(`^[A-Z]{2}[A-Z0-9]{9}[0-9]$`)

// Same pattern as Symfony Issn::PATTERN (optional hyphen between the two groups).
var issnPattern = regexp.MustCompile(`^[0-9]{4}-?[0-9]{3}[0-9Xx]$`)

// ULID validates whether the value is a valid ULID (Universally Unique Lexicographically Sortable Identifier).
// See https://github.com/ulid/spec for ULID specifications.
//
// Possible errors:
//   - [ErrTooShort] on values with length less than 26;
//   - [ErrTooLong] on values with length greater than 26;
//   - [ErrInvalidCharacters] on values with unexpected characters;
//   - [ErrTooLarge] on too big value (larger than '7ZZZZZZZZZZZZZZZZZZZZZZZZZ').
func ULID(value string) error {
	if len(value) < 26 {
		return ErrTooShort
	}
	if len(value) > 26 {
		return ErrTooLong
	}

	for _, c := range value {
		if !ulidChars.Contains(c) {
			return ErrInvalidCharacters
		}
	}

	// Largest valid ULID is '7ZZZZZZZZZZZZZZZZZZZZZZZZZ'
	// See https://github.com/ulid/spec#overflow-errors-when-parsing-base32-strings
	if value[0] > '7' {
		return ErrTooLarge
	}

	return nil
}

// ISIN validates whether the value is a valid International Securities Identification Number (ISIN).
// Letters are treated case-insensitively (normalized to upper case for validation), matching
// Symfony\Component\Validator\Constraints\Isin.
//
// Possible errors:
//   - [ErrTooShort] on values with length less than 12;
//   - [ErrTooLong] on values with length greater than 12;
//   - [ErrInvalidCharacters] on values that do not match the ISIN pattern;
//   - [ErrInvalidChecksum] when the check digit is wrong;
//
// See https://en.wikipedia.org/wiki/International_Securities_Identification_Number.
func ISIN(value string) error {
	if len(value) < 12 {
		return ErrTooShort
	}
	if len(value) > 12 {
		return ErrTooLong
	}

	s := strings.ToUpper(value)
	if !isinPattern.MatchString(s) {
		return ErrInvalidCharacters
	}
	if !isinLUHN(s) {
		return ErrInvalidChecksum
	}

	return nil
}

// ISSN validates whether the value is a valid International Standard Serial Number (ISSN).
// The check digit uses the ISO 3297 mod 11 algorithm (weights 8…2 on the first seven digits).
// An optional hyphen is allowed between the two four-character groups, matching
// Symfony\Component\Validator\Constraints\Issn.
//
// Possible errors:
//   - [ErrTooShort] when fewer than eight digits/check characters remain after removing hyphens;
//   - [ErrTooLong] when more than eight remain;
//   - [ErrInvalidCharacters] when the format is not NNNN-NNNC (with optional hyphen) or digits are invalid;
//   - [ErrInvalidChecksum] when the check character is wrong;
//
// See https://www.issn.org/understanding-the-issn/what-is-an-issn/ and ISO 3297.
func ISSN(value string) error {
	if issnPattern.MatchString(value) {
		return issnValidateBody(strings.ToUpper(strings.ReplaceAll(value, "-", "")))
	}
	s := strings.ReplaceAll(value, "-", "")
	switch {
	case len(s) < 8:
		return ErrTooShort
	case len(s) > 8:
		return ErrTooLong
	default:
		return ErrInvalidCharacters
	}
}

func issnValidateBody(canonical string) error {
	var sum int
	for i := 0; i < 7; i++ {
		c := canonical[i]
		if c < '0' || c > '9' {
			return ErrInvalidCharacters
		}
		sum += int(c-'0') * (8 - i)
	}
	check := 11 - (sum % 11)
	if check == 11 {
		check = 0
	}
	var expected byte
	if check < 10 {
		expected = byte('0' + check)
	} else {
		expected = 'X'
	}
	if canonical[7] != expected {
		return ErrInvalidChecksum
	}
	return nil
}

func isinCharValue(c byte) int {
	switch {
	case c >= '0' && c <= '9':
		return int(c - '0')
	case c >= 'A' && c <= 'Z':
		return int(c - 'A' + 10)
	default:
		return -1
	}
}

func isinLUHN(s string) bool {
	var b strings.Builder
	b.Grow(24)
	for i := 0; i < 12; i++ {
		v := isinCharValue(s[i])
		if v < 0 {
			return false
		}
		b.WriteString(strconv.Itoa(v))
	}

	return luhnValidDigits(b.String())
}

// UUID validates whether a string value is a valid UUID (also known as GUID).
//
// By default, it uses strict mode and checks the UUID as specified in RFC 4122.
// To parse additional formats, use the [AllowNonCanonicalUUIDFormats] option.
//
// In addition, it checks if the UUID version matches one of
// the registered versions: 1, 2, 3, 4, 5, 6 or 7.
// Use [AllowUUIDVersions] to validate for a specific set of versions.
//
// Nil UUID ("00000000-0000-0000-0000-000000000000") values are considered as valid.
// Use [DenyNilUUID] to disallow nil value.
//
// Possible errors:
//   - [ErrTooShort] on values with length less than 36 (or 32 for non-canonical formats);
//   - [ErrTooLong] on values with length greater than 36 (or 45 for non-canonical formats);
//   - [ErrInvalidCharacters] on values with unexpected characters;
//   - [ErrInvalidHyphenPlacement] on invalid placements of hyphens;
//   - [ErrInvalidVersion] on a restricted versions;
//   - [ErrIsNil] on nil value (if [DenyNilUUID] options is enabled);
//   - [ErrInvalid] on other cases;
//
// See http://tools.ietf.org/html/rfc4122.
func UUID(value string, options ...func(o *UUIDOptions)) error {
	opts := &UUIDOptions{}
	for _, set := range options {
		set(opts)
	}

	var (
		u   uuid.UUID
		err error
	)

	if opts.isNonStrict {
		u, err = uuid.FromString(value)
	} else {
		u, err = uuid.CanonicalFromString(value)
	}
	if err != nil {
		return convertUUIDError(err)
	}

	if !u.IsNil() && !u.ValidVersion(opts.versions...) {
		return ErrInvalidVersion
	}

	if opts.isNotNil && u.IsNil() {
		return ErrIsNil
	}

	return nil
}

// AllowNonCanonicalUUIDFormats used to enable parsing UUID value from non-canonical formats.
//
// Following formats are supported:
//   - "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
//   - "{6ba7b810-9dad-11d1-80b4-00c04fd430c8}",
//   - "urn:uuid:6ba7b810-9dad-11d1-80b4-00c04fd430c8"
//   - "6ba7b8109dad11d180b400c04fd430c8"
//   - "{6ba7b8109dad11d180b400c04fd430c8}",
//   - "urn:uuid:6ba7b8109dad11d180b400c04fd430c8".
func AllowNonCanonicalUUIDFormats() func(o *UUIDOptions) {
	return func(o *UUIDOptions) {
		o.isNonStrict = true
	}
}

// AllowUUIDVersions used to set valid versions of the UUID value.
// If the versions are empty, the UUID will be checked for compliance with the default
// registered versions: 1, 2, 3, 4, 5, 6 or 7.
func AllowUUIDVersions(versions ...byte) func(o *UUIDOptions) {
	return func(o *UUIDOptions) {
		o.versions = versions
	}
}

// DenyNilUUID used to treat nil UUID ("00000000-0000-0000-0000-000000000000") value as invalid.
func DenyNilUUID() func(o *UUIDOptions) {
	return func(o *UUIDOptions) {
		o.isNotNil = true
	}
}

// UUIDOptions are used to set up validation process of the [UUID] function.
type UUIDOptions struct {
	versions    []byte
	isNonStrict bool
	isNotNil    bool
}

//nolint:errorlint
func convertUUIDError(err error) error {
	switch err {
	case uuid.ErrTooShort:
		return ErrTooShort
	case uuid.ErrTooLong:
		return ErrTooLong
	case uuid.ErrInvalidHyphenPlacement:
		return ErrInvalidHyphenPlacement
	}

	return ErrInvalid
}

type charSet map[rune]struct{}

func (s charSet) Contains(c rune) bool {
	_, exist := s[c]

	return exist
}

func newCharSet(s string) charSet {
	chars := make(charSet, len(s))

	for _, c := range s {
		chars[c] = struct{}{}
	}

	return chars
}
