package scenario

import (
	"context"
	"github.com/muonsoft/validation"
)

type Legacy struct{ Value int }

var ErrLimit = validation.NewError("limit", "Value exceeds the limit.")

func (x Legacy) Validate(ctx context.Context, v *validation.Validator, limit int) error {
	if v.IsIgnoredForGroups("strict") {
		return nil
	}
	if x.Value > limit {
		return v.AtProperty("value").CreateViolation(ctx, ErrLimit, ErrLimit.Message())
	}
	return nil
}
