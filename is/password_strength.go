package is

import "github.com/muonsoft/validation/validate"

// PasswordStrength reports whether a password meets the configured minimum score.
// See [validate.PasswordStrength] for options and empty-value behavior.
func PasswordStrength(password string, options ...func(*validate.PasswordStrengthOptions)) bool {
	return validate.PasswordStrength(password, options...) == nil
}
