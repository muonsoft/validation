package scenario

import (
	"context"
	"github.com/muonsoft/validation"
)

func (c PrefixConstraint) ValidateString(ctx context.Context, v *validation.Validator, value *string) error {
	return nil
}
func NewValidator() (*validation.Validator, error) { return validation.NewValidator() }
