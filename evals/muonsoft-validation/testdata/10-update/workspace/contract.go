package scenario

import (
	"context"
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
)

type Record struct{ ID, Name string }

func (r Record) Validate(ctx context.Context, v *validation.Validator) error {
	return v.Validate(ctx, validation.StringProperty("name", r.Name, it.IsNotBlank()))
}

var ErrDuplicate = validation.NewError("duplicate_name", "Name is already used.")

type Repository interface {
	Load(context.Context, string) (Record, error)
	NameTaken(context.Context, string, string) (bool, error)
	Save(context.Context, Record) error
}
type Service struct{ Repository Repository }
