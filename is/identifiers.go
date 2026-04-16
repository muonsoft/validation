package is

import "github.com/muonsoft/validation/validate"

// ULID validates whether the value is a valid ULID (Universally Unique Lexicographically Sortable Identifier).
// See https://github.com/ulid/spec for ULID specifications.
func ULID(value string) bool {
	return validate.ULID(value) == nil
}

// IBAN validates whether the value is a valid International Bank Account Number (IBAN).
// See [github.com/muonsoft/validation/validate.IBAN] for validation rules.
//
// See https://en.wikipedia.org/wiki/International_Bank_Account_Number.
func IBAN(value string) bool {
	return validate.IBAN(value) == nil
}

// BIC validates whether the value is a valid Business Identifier Code (BIC / SWIFT).
// See [github.com/muonsoft/validation/validate.BIC] for validation rules and options.
//
// See https://en.wikipedia.org/wiki/ISO_9362.
func BIC(value string, options ...func(*validate.BICOptions)) bool {
	return validate.BIC(value, options...) == nil
}

// ISIN validates whether the value is a valid International Securities Identification Number (ISIN).
// See [github.com/muonsoft/validation/validate.ISIN] for validation rules and possible errors.
//
// See https://en.wikipedia.org/wiki/International_Securities_Identification_Number.
func ISIN(value string) bool {
	return validate.ISIN(value) == nil
}

// ISSN validates whether the value is a valid International Standard Serial Number (ISSN).
// See [github.com/muonsoft/validation/validate.ISSN] for validation rules and possible errors.
//
// See https://www.issn.org/understanding-the-issn/what-is-an-issn/.
func ISSN(value string) bool {
	return validate.ISSN(value) == nil
}

// ISBN validates whether the value is a valid ISBN-10 or ISBN-13.
// See [github.com/muonsoft/validation/validate.ISBN] for options and possible errors.
//
// See https://en.wikipedia.org/wiki/ISBN.
func ISBN(value string, options ...func(*validate.ISBNOptions)) bool {
	return validate.ISBN(value, options...) == nil
}

// LUHN validates whether the value passes the Luhn (mod 10) checksum.
// See [github.com/muonsoft/validation/validate.LUHN] for validation rules and possible errors.
//
// See https://en.wikipedia.org/wiki/Luhn_algorithm.
func LUHN(value string) bool {
	return validate.LUHN(value) == nil
}

// Currency validates whether the value is a recognized ISO 4217 alphabetic currency code.
// See [github.com/muonsoft/validation/validate.Currency] for rules and possible errors.
//
// See https://www.iso.org/iso-4217-currency-codes.html.
func Currency(value string) bool {
	return validate.Currency(value) == nil
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
