package scenario

import (
	"context"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
)

func (x ProfilePatch) Validate(ctx context.Context, v *validation.Validator) error {
	var args []validation.Argument
	if x.Country.Present && x.Country.Value != nil {
		args = append(args, validation.Sequentially(
			validation.StringProperty("country", *x.Country.Value, it.IsNotBlank()),
			validation.StringProperty("country", *x.Country.Value, it.IsCountry()),
			validation.StringProperty("country", *x.Country.Value, it.IsOneOf("DE", "US", "GB").WithoutBlank().WithError(ErrCountryPolicy).WithMessage(ErrCountryPolicy.Message()))))
	}
	if x.Locale.Present && x.Locale.Value != nil {
		args = append(args, validation.Sequentially(
			validation.StringProperty("locale", *x.Locale.Value, it.IsNotBlank()),
			validation.StringProperty("locale", *x.Locale.Value, it.IsLocale()),
			validation.StringProperty("locale", *x.Locale.Value, it.IsOneOf("en-US", "de-DE").WithoutBlank().WithError(ErrLocalePolicy).WithMessage(ErrLocalePolicy.Message()))))
	}
	return v.Validate(ctx, args...)
}
