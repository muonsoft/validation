package scenario

import (
	"context"

	"github.com/muonsoft/validation"
)

func (x Article) Validate(ctx context.Context, v *validation.Validator) error {
	if x.Title == "" {
		return v.AtProperty("title").CreateViolation(ctx, ErrTitleRequired, ErrTitleRequired.Message())
	}
	if len(x.Title) < 3 || len(x.Title) > 12 {
		return v.AtProperty("title").CreateViolation(ctx, ErrTitleLength, ErrTitleLength.Message())
	}
	if x.Contact == "" {
		return v.AtProperty("contact").CreateViolation(ctx, ErrContactRequired, ErrContactRequired.Message())
	}
	return nil
}
