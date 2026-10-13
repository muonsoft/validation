package scenario

import "github.com/muonsoft/validation/validate"

type PasswordUpdate struct {
	Password *string
	Estimate func(string) validate.PasswordStrengthScore
}
