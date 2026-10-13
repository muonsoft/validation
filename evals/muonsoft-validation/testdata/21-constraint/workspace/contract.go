package scenario

import (
	"context"

	"github.com/muonsoft/validation"
)

type Entry struct{ Category string }
type Catalog interface {
	Lookup(context.Context, string) (Entry, bool, error)
}

var ErrUnknown = validation.NewError("unknown_code", "Code was not found.")
var ErrCategory = validation.NewError("wrong_category", "Code must belong to {{ category }}.")

type Signup struct {
	Code  *string
	Label string
}
type Transfer struct {
	Destination string
	Label       string
}
