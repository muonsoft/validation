package is

import "github.com/muonsoft/validation/validate"

// NoSuspiciousCharacters reports whether value passes the same checks as [validate.NoSuspiciousCharacters].
func NoSuspiciousCharacters(value string, options ...validate.NoSuspiciousCharactersOption) bool {
	return validate.NoSuspiciousCharacters(value, options...) == nil
}
