package is

import "github.com/muonsoft/validation/validate"

// ULID validates whether the value is a valid ULID (Universally Unique Lexicographically Sortable Identifier).
// See https://github.com/ulid/spec for ULID specifications.
func ULID(value string) bool {
	return validate.ULID(value) == nil
}
