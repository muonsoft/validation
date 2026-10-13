package scenario

import (
	"context"

	"github.com/muonsoft/validation"
)

type CodeConstraint struct{}

func NewCodeConstraint(catalog Catalog, category string) CodeConstraint { return CodeConstraint{} }
func (c CodeConstraint) ValidateString(ctx context.Context, v *validation.Validator, value *string) error {
	return nil
}
func ValidatorOptions() []validation.ValidatorOption { return nil }
