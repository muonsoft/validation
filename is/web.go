package is

import "github.com/muonsoft/validation/validate"

// Email is used for simplified validation of an email address. It allows all values
// with an "@" symbol in, and a "." in the second host part of the email address.
func Email(value string) bool {
	return looseEmailRegex.MatchString(value)
}

// HTML5Email is used for validation of an email address based on pattern for HTML5
// (see https://html.spec.whatwg.org/multipage/input.html#valid-e-mail-address).
func HTML5Email(value string) bool {
	return html5EmailRegex.MatchString(value)
}

// URL is used to check that value is a valid URL string. By default (if no schemas are passed),
// the function checks only for the http:// and https:// schemas. Use the schemas argument
// to configure the list of expected schemas. If an empty string is passed as a schema, then
// URL value may be treated as relative (without schema, e.g. "//example.com").
func URL(value string, schemas ...string) bool {
	return validate.URL(value, schemas...) == nil
}
