package scenario

import (
	"context"

	"github.com/muonsoft/validation"
)

func (x Request) Validate(ctx context.Context, v *validation.Validator) error { return nil }
func (x Contact) Validate(ctx context.Context, v *validation.Validator) error { return nil }
