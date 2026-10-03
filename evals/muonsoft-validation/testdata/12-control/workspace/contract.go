package scenario

import (
	"errors"
)

var ErrEmpty = errors.New("empty name")

func ValidateName(name string) error {
	if name == "" {
		return ErrEmpty
	}
	return nil
}
