package scenario

import (
	"context"

	"github.com/muonsoft/validation"
)

type Reference struct {
	ID     string
	Weight int
}
type Bundle struct {
	ID, Name, PreviousName string
	References             []Reference
}
type Command struct {
	ID, Name   string
	References []Reference
}
type Repository interface {
	Load(context.Context, string) (*Bundle, error)
	// NameTaken checks name availability, excluding the bundle with excludedID.
	NameTaken(ctx context.Context, name, excludedID string) (bool, error)
	Find(context.Context, []string) (map[string]bool, error)
	Save(context.Context, *Bundle) error
}
type Service struct {
	Repository Repository
	Validator  *validation.Validator
}

var ErrDuplicate = validation.NewError("duplicate", "Name is already used.")
var ErrMissing = validation.NewError("missing", "Reference does not exist.")
