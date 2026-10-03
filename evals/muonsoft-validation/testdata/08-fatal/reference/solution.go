package scenario

import (
	"context"
	"github.com/muonsoft/validation"
)

func ValidateAndSave(ctx context.Context, v *validation.Validator, first, second, third validation.Validatable, save Save) error {
	if err := v.Validate(ctx, validation.Valid(first), validation.Valid(second), validation.Valid(third)); err != nil {
		return err
	}
	return save(ctx)
}
