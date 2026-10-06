package scenario

import (
	"context"

	"github.com/muonsoft/validation"
)

var ErrMissing = validation.NewError("branch_missing", "Item is unavailable.")

type Item struct {
	ID       string
	Quantity int
}
type Request struct{ Items []Item }
type Lookup interface {
	Exists(context.Context, string) (bool, error)
}
