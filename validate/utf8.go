package validate

import (
	"errors"
	"unicode/utf8"
)

// ErrInvalidUTF8 indicates a malformed UTF-8 byte sequence.
var ErrInvalidUTF8 = errors.New("invalid UTF-8")

// UTF8 returns ErrInvalidUTF8 if value contains malformed UTF-8, otherwise nil.
// The empty string is valid.
//
// Go strings can contain arbitrary bytes. Use this check at input boundaries,
// such as file imports or legacy integrations, to reject truncated or malformed
// text before storage or processing that might reject or replace invalid bytes.
//
// This checks UTF-8 only, unlike Symfony's configurable Charset constraint.
// It does not detect or convert other encodings, normalize text, or validate
// its meaning or safety. Validly encoded mojibake, U+FFFD, invisible characters,
// and HTML are accepted. Validate the original bytes before decoding replaces
// malformed sequences; such replacements cannot be detected by this check.
func UTF8(value string) error {
	if !utf8.ValidString(value) {
		return ErrInvalidUTF8
	}
	return nil
}
