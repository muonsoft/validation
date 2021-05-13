package validate

import (
	"errors"
	"net/url"
)

var (
	ErrUnexpectedSchema = errors.New("unexpected schema")
)

// URL is used to validate that value is a valid URL string. By default (if no schemas are passed),
// the function checks only for the http:// and https:// schemas. Use the schemas argument
// to configure the list of expected schemas. If an empty string is passed as a schema, then
// URL value may be treated as relative (without schema, e.g. "//example.com").
//
// If value is not a valid URL the function will return one of the errors:
//	• parsing error from url.Parse method if value cannot be parsed as an URL;
//	• ErrUnexpectedSchema if schema is not matching one of the listed schemas;
//	• ErrInvalid if value is not matching the regular expression.
func URL(value string, schemas ...string) error {
	if len(schemas) == 0 {
		schemas = []string{"http", "https"}
	}
	u, err := url.Parse(value)
	if err != nil {
		return err
	}

	err = validateSchema(u, schemas)
	if err != nil {
		return err
	}

	if !urlRegex.MatchString(value) {
		return ErrInvalid
	}

	return nil
}

func validateSchema(u *url.URL, schemas []string) error {
	for _, schema := range schemas {
		if schema == u.Scheme {
			return nil
		}
	}

	return ErrUnexpectedSchema
}
