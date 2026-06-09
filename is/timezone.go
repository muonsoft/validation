package is

import "github.com/muonsoft/validation/validate"

// Timezone validates whether the value is a known IANA timezone identifier.
// See [github.com/muonsoft/validation/validate.Timezone] for rules and options.
func Timezone(value string, options ...func(*validate.TimezoneOptions)) bool {
	return validate.Timezone(value, options...) == nil
}
