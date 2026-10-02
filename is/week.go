package is

import "github.com/muonsoft/validation/validate"

// Week reports whether value is a valid ISO week within optional inclusive bounds.
// See [validate.Week] for format, year range, options, and empty-value behavior.
func Week(value string, options ...func(*validate.WeekOptions)) bool {
	return validate.Week(value, options...) == nil
}
