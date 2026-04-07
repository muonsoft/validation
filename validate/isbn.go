package validate

import (
	"errors"
	"strings"
)

// ISBN validation errors, aligned with Symfony\Component\Validator\Constraints\Isbn codes.
var (
	ErrISBNInvalidCharacters = errors.New("ISBN invalid characters")
	ErrISBNTooShort          = errors.New("ISBN too short")
	ErrISBNTooLong           = errors.New("ISBN too long")
	ErrISBNChecksumFailed    = errors.New("ISBN checksum failed")
	ErrISBNTypeNotRecognized = errors.New("ISBN type not recognized")
)

// ISBNType selects which ISBN format to validate, matching Symfony Isbn "type".
type ISBNType int

const (
	// ISBNTypeAny accepts ISBN-10 or ISBN-13 (default).
	ISBNTypeAny ISBNType = iota
	// ISBNType10 accepts only ISBN-10 after removing hyphens.
	ISBNType10
	// ISBNType13 accepts only ISBN-13 after removing hyphens.
	ISBNType13
)

// ISBN validates whether the value is a valid ISBN-10 or ISBN-13.
// Hyphens (U+002D) are stripped; other characters are not allowed in the canonical form.
// Behavior is aligned with Symfony\Component\Validator\Constraints\Isbn and IsbnValidator.
//
// Empty string is considered valid (use [NotBlank] or similar to reject empty values).
//
// Possible errors:
//   - [ErrISBNInvalidCharacters], [ErrISBNTooShort], [ErrISBNTooLong], [ErrISBNChecksumFailed]
//     for the selected type;
//   - [ErrISBNTypeNotRecognized] in "any" mode when the length is between 11 and 12 digits.
//
// See https://en.wikipedia.org/wiki/ISBN.
func ISBN(value string, options ...func(*ISBNOptions)) error {
	if value == "" {
		return nil
	}

	opts := &ISBNOptions{typ: ISBNTypeAny}
	for _, set := range options {
		set(opts)
	}

	canonical := strings.ReplaceAll(value, "-", "")

	switch opts.typ {
	case ISBNType10:
		return validateISBN10Body(canonical)
	case ISBNType13:
		return validateISBN13Body(canonical)
	default:
		return validateISBNAny(canonical)
	}
}

// ISBNOptions configures [ISBN].
type ISBNOptions struct {
	typ ISBNType
}

// ISBNOnly10 restricts validation to ISBN-10 (Symfony Isbn::ISBN_10).
func ISBNOnly10() func(*ISBNOptions) {
	return func(o *ISBNOptions) {
		o.typ = ISBNType10
	}
}

// ISBNOnly13 restricts validation to ISBN-13 (Symfony Isbn::ISBN_13).
func ISBNOnly13() func(*ISBNOptions) {
	return func(o *ISBNOptions) {
		o.typ = ISBNType13
	}
}

func validateISBNAny(canonical string) error {
	code := validateISBN10Body(canonical)
	if errors.Is(code, ErrISBNTooLong) {
		code = validateISBN13Body(canonical)
		if errors.Is(code, ErrISBNTooShort) {
			return ErrISBNTypeNotRecognized
		}
	}
	return code
}

func parseISBN10Char(c byte) (digit int, err error) {
	switch {
	case c == 'X':
		return 10, nil
	case c == 'x':
		// Lowercase is invalid in Symfony (only uppercase X is accepted as check digit letter).
		return 0, ErrISBNInvalidCharacters
	case c >= '0' && c <= '9':
		return int(c - '0'), nil
	default:
		return 0, ErrISBNInvalidCharacters
	}
}

// validateISBN10Body mirrors Symfony IsbnValidator::validateIsbn10 (error priority:
// invalid characters, too short/long, checksum).
func validateISBN10Body(isbn string) error {
	var checkSum int
	for i := 0; i < 10; i++ {
		if i >= len(isbn) {
			return ErrISBNTooShort
		}
		digit, err := parseISBN10Char(isbn[i])
		if err != nil {
			return err
		}
		checkSum += digit * (10 - i)
	}
	if len(isbn) > 10 {
		return ErrISBNTooLong
	}
	if checkSum%11 != 0 {
		return ErrISBNChecksumFailed
	}
	return nil
}

// validateISBN13Body mirrors Symfony IsbnValidator::validateIsbn13.
func validateISBN13Body(isbn string) error {
	for i := 0; i < len(isbn); i++ {
		if isbn[i] < '0' || isbn[i] > '9' {
			return ErrISBNInvalidCharacters
		}
	}
	switch {
	case len(isbn) < 13:
		return ErrISBNTooShort
	case len(isbn) > 13:
		return ErrISBNTooLong
	}

	var checkSum int
	for i := 0; i < 13; i += 2 {
		checkSum += int(isbn[i] - '0')
	}
	for i := 1; i < 12; i += 2 {
		checkSum += int(isbn[i]-'0') * 3
	}
	if checkSum%10 != 0 {
		return ErrISBNChecksumFailed
	}
	return nil
}
