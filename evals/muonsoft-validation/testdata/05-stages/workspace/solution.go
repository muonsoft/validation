package scenario

import (
	"context"
	"github.com/muonsoft/validation"
)

func (x Input) Validate(ctx context.Context, v *validation.Validator) error          { return nil }
func (s Service) Handle(ctx context.Context, v *validation.Validator, x Input) error { return nil }
