package scenario

import (
	"context"
	"github.com/muonsoft/validation"
)

type Input struct{ Name, First, Second string }

var ErrMissing = validation.NewError("missing_reference", "Reference is missing.")

type Repository interface {
	Exists(context.Context, string) (bool, error)
	Save(context.Context, Input) error
}
type Service struct{ Repository Repository }
