package validate

import "errors"

var (
	ErrTooShort          = errors.New("too short")
	ErrTooLong           = errors.New("too long")
	ErrTooLarge          = errors.New("too large")
	ErrInvalidCharacters = errors.New("invalid characters")
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
