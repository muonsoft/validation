package is

import "github.com/muonsoft/validation/validate"

// CardScheme reports whether value matches any of the explicitly selected schemes.
// See [validate.CardScheme] for supported schemes and empty-value behavior.
func CardScheme(value string, schemes ...validate.CardSchemeName) bool {
	return validate.CardScheme(value, schemes...) == nil
}
