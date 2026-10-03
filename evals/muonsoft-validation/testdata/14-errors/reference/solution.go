package scenario

import "errors"

func ErrorText(err error) string {
	if err == nil {
		return ""
	}
	if errors.Is(err, ErrInvalid) {
		return "Invalid input"
	}
	return err.Error()
}
