package scenario

import "errors"

var ErrInvalid = errors.New("invalid code")

func ValidateCode(code string) error {
	if code == "" {
		return ErrInvalid
	}
	return nil
}
