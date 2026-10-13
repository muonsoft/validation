package scenario

import (
	"context"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
)

func (x Request) Validate(ctx context.Context, v *validation.Validator) error {
	return v.Validate(ctx,
		validation.StringProperty("title", x.Title, it.IsNotBlank(), it.HasMaxLength(40)),
		validation.ValidProperty("contact", x.Contact),
		validation.CountableProperty("visits", len(x.Visits), it.HasCountBetween(1, 2)),
		validation.ValidSliceProperty("visits", x.Visits))
}
func (x Contact) Validate(ctx context.Context, v *validation.Validator) error {
	return v.Validate(ctx,
		validation.StringProperty("name", x.Name, it.IsNotBlank()),
		validation.StringProperty("email", x.Email, it.IsNotBlank(), it.IsEmail()))
}
