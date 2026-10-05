package scenario

import (
	"context"

	"github.com/muonsoft/validation"
)

type Document struct {
	Title    string
	Schedule *Schedule
}
type Schedule struct {
	Note         *string
	Mode         string
	Lower, Upper *int
	Limits       Limits
}

var ErrRange = validation.NewError("range", "Upper must exceed lower.")
var ErrLimit = validation.NewError("limit", "Value exceeds {{ maximum }}.")

type Limits interface {
	Validate(context.Context, *validation.Validator, int) error
}
