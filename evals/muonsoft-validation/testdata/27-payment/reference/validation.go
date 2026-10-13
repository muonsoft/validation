package scenario

import (
	"context"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/validate"
)

func (x Payment) Validate(ctx context.Context, v *validation.Validator) error {
	return v.Validate(ctx,
		validation.StringProperty("card", x.Card, it.IsNotBlank(), it.IsCardScheme(validate.CardSchemeVisa, validate.CardSchemeMastercard), it.IsLUHN()),
		validation.StringProperty("iban", x.IBAN, it.IsIBAN()),
		validation.StringProperty("bic", x.BIC, it.IsBIC()),
	)
}
