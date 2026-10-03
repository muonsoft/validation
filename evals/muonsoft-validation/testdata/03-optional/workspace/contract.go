package scenario

import (
	"github.com/muonsoft/validation"
)

type Input struct {
	Note         *string
	Status       string
	Lower, Upper *int
}

var ErrRange = validation.NewError("invalid_range", "Upper must exceed lower.")
