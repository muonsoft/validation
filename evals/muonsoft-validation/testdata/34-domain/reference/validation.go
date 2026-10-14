package scenario

import (
	"context"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
)

func (x Signup) Validate(ctx context.Context, v *validation.Validator, policy EmailPolicy) error {
	return v.Validate(ctx,
		validation.StringProperty("name", x.Name, it.IsNotBlank()),
		validation.Sequentially(
			validation.StringProperty("email", x.Email, it.IsNotBlank(), it.IsEmail()),
			validation.StringProperty("email", x.Email, policy)))
}
