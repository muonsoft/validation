package scenario

import (
	"context"

	"github.com/muonsoft/validation"
)

func (x Bundle) Validate(ctx context.Context, v *validation.Validator) error { return nil }
func CheckReferences(ctx context.Context, v *validation.Validator, r Repository, refs []Reference) error {
	return nil
}
