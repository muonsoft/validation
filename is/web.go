package is

import "github.com/muonsoft/validation/validate"

// URL is used to check that value is a valid URL string. By default (if no schemas are passed),
// the function checks only for the http:// and https:// schemas. Use the schemas argument
// to configure the list of expected schemas. If an empty string is passed as a schema, then
// URL value may be treated as relative (without schema, e.g. "//example.com").
func URL(value string, schemas ...string) bool {
	return validate.URL(value, schemas...) == nil
}

func Email(value string) bool {
	return looseEmailRegex.MatchString(value)
}

func HTML5Email(value string) bool {
	return html5EmailRegex.MatchString(value)
}
