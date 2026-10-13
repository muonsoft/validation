package scenario

import (
	"context"

	"github.com/muonsoft/validation"
)

// ValidateString delegates to the application-owned tenant policy.
func (p EmailPolicy) ValidateString(ctx context.Context, v *validation.Validator, value *string) error {
	if value == nil || *value == "" {
		return nil
	}
	return p.Check(ctx, v, *value)
}
