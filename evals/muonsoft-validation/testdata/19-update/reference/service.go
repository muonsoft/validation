package scenario

import (
	"context"
	"strings"

	"github.com/muonsoft/validation"
)

func (s Service) Update(ctx context.Context, c Command) error {
	loaded, err := s.Repository.Load(ctx, c.ID)
	if err != nil {
		return err
	}
	next := *loaded
	next.Name = strings.TrimSpace(c.Name)
	next.References = append([]Reference(nil), c.References...)
	for i := range next.References {
		next.References[i].ID = strings.TrimSpace(next.References[i].ID)
	}
	if err := next.Validate(ctx, s.Validator); err != nil {
		return err
	}
	err = s.Validator.Validate(ctx,
		validation.ValidProperty("name", validation.ValidatableFunc(func(ctx context.Context, v *validation.Validator) error {
			taken, err := s.Repository.NameTaken(ctx, next.Name, next.ID)
			if err != nil {
				return err
			}
			if taken {
				return v.CreateViolation(ctx, ErrDuplicate, ErrDuplicate.Message())
			}
			return nil
		})),
		validation.ValidProperty("references", validation.ValidatableFunc(func(ctx context.Context, v *validation.Validator) error {
			return CheckReferences(ctx, v, s.Repository, next.References)
		})))
	if err != nil {
		return err
	}
	return s.Repository.Save(ctx, &next)
}
