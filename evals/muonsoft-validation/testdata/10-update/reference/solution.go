package scenario

import (
	"context"
	"github.com/muonsoft/validation"
	"strings"
)

func (s Service) Update(ctx context.Context, v *validation.Validator, id, name string) error {
	r, err := s.Repository.Load(ctx, id)
	if err != nil {
		return err
	}
	r.Name = strings.TrimSpace(name)
	if err := v.ValidateIt(ctx, r); err != nil {
		return err
	}
	taken, err := s.Repository.NameTaken(ctx, r.Name, r.ID)
	if err != nil {
		return err
	}
	if taken {
		return v.AtProperty("name").CreateViolation(ctx, ErrDuplicate, ErrDuplicate.Message())
	}
	return s.Repository.Save(ctx, r)
}
