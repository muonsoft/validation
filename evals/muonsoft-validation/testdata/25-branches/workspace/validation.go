package scenario

import (
	"context"

	"github.com/muonsoft/validation"
)

func (x Request) Validate(ctx context.Context, v *validation.Validator, lookup Lookup) error {
	return nil
}
