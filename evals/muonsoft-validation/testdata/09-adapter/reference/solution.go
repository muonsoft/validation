package scenario

import (
	"context"
	"github.com/muonsoft/validation"
)

func ValidateDetail(ctx context.Context, v *validation.Validator, x Legacy, limit int) error {
	return v.Validate(ctx, validation.ValidProperty("detail", validation.ValidatableFunc(func(ctx context.Context, child *validation.Validator) error { return x.Validate(ctx, child, limit) })))
}
