package scenario

import (
	"context"
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
)

func (x Document) Validate(ctx context.Context, v *validation.Validator) error {
	return v.Validate(ctx, validation.StringProperty("code", x.Code, it.IsNotBlank()), validation.StringProperty("title", x.Title, it.IsNotBlank().WhenGroups("publish")))
}
