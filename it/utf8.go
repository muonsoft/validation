package it

import (
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/is"
)

// IsUTF8 checks that a string contains valid UTF-8 bytes. Go strings do not
// guarantee this. Use it to reject malformed or truncated text from imports
// or integrations before storage or further processing.
//
// Nil and empty values are ignored; use [IsNotBlank] to require a value.
// Only UTF-8 is supported, unlike Symfony's configurable Charset constraint.
// This does not detect or convert encodings, normalize text, or check content
// safety. Validly encoded mojibake, U+FFFD, invisible characters, and HTML pass.
// Check original input before a decoder replaces malformed bytes with U+FFFD.
func IsUTF8() validation.StringFuncConstraint {
	return validation.OfStringBy(is.UTF8).
		WithError(validation.ErrInvalidUTF8).
		WithMessage(validation.ErrInvalidUTF8.Message())
}
