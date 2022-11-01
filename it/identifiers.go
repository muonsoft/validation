package it

import (
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/is"
)

// IsULID validates whether the value is a valid ULID (Universally Unique Lexicographically Sortable Identifier).
// See https://github.com/ulid/spec for ULID specifications.
func IsULID() validation.StringFuncConstraint {
	return validation.OfStringBy(is.ULID).
		WithError(validation.ErrInvalidULID).
		WithMessage(validation.ErrInvalidULID.Message())
}
