package scenario

import (
	"context"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
)

func (x Publication) Validate(ctx context.Context, v *validation.Validator) error {
	return v.Validate(ctx,
		validation.StringProperty("country", x.Country, it.IsCountry()),
		validation.StringProperty("language", x.Language, it.IsLanguage()),
		validation.StringProperty("locale", x.Locale, it.IsLocale()),
		validation.StringProperty("currency", x.Currency, it.IsCurrency()),
		validation.StringProperty("week", x.Week, it.IsWeek().WithMin("2020-W53").WithMax("2021-W02")),
		validation.StringProperty("summary", x.Summary, it.HasWordCountBetween(2, 4)),
	)
}
