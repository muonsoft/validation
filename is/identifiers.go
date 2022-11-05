package is

import "github.com/muonsoft/validation/validate"

// ULID validates whether the value is a valid ULID (Universally Unique Lexicographically Sortable Identifier).
// See https://github.com/ulid/spec for ULID specifications.
func ULID(value string) bool {
	return validate.ULID(value) == nil
}

// UUID validates whether a string value is a valid UUID (also known as GUID).
//
// By default, it uses strict mode and checks the UUID as specified in RFC 4122.
// To parse additional formats, use the
// [github.com/muonsoft/validation/validate.AllowNonCanonicalUUIDFormats] option.
//
// In addition, it checks if the UUID version matches one of
// the registered versions: 1, 2, 3, 4, 5, 6 or 7.
// Use [github.com/muonsoft/validation/validate.AllowUUIDVersions] to validate
// for a specific set of versions.
//
// Nil UUID ("00000000-0000-0000-0000-000000000000") values are considered as valid.
// Use [github.com/muonsoft/validation/validate.DenyNilUUID] to disallow nil value.
//
// See http://tools.ietf.org/html/rfc4122.
func UUID(value string, options ...func(o *validate.UUIDOptions)) bool {
	return validate.UUID(value, options...) == nil
}
