package scenario

import (
	"context"
	"io"

	"github.com/muonsoft/validation"
)

func (x Upload) Validate(ctx context.Context, v *validation.Validator, body io.Reader) error {
	return nil
}
