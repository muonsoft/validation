package scenario

import (
	"context"
	"github.com/muonsoft/validation"
)

func ValidateAndSave(ctx context.Context, v *validation.Validator, first, second, third validation.Validatable, save Save) error {
	if err := validation.Filter(first.Validate(ctx, v), second.Validate(ctx, v), third.Validate(ctx, v)); err != nil {
		return err
	}
	return save(ctx)
}
