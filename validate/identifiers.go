package validate

import (
	"errors"

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
