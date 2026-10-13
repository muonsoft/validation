package scenario

import "github.com/muonsoft/validation"

var ErrCountryPolicy = validation.NewError("country_policy", "Country is not served.")
var ErrLocalePolicy = validation.NewError("locale_policy", "Locale is not supported.")

type TextPatch struct {
	Present bool
	Value   *string
}
type ProfilePatch struct{ Country, Locale TextPatch }
