package is

import "github.com/muonsoft/validation/validate"

// CSSColor reports whether value is a CSS color in any selected format.
// No formats selects all supported formats.
// See [validate.CSSColor] for supported formats and empty-value behavior.
func CSSColor(value string, formats ...validate.CSSColorFormat) bool {
	return validate.CSSColor(value, formats...) == nil
}
