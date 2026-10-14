package scenario_test

import (
	"testing"

	s "example.com/validation-eval/scenario"
	"github.com/muonsoft/validation"
)

func TestPaymentIndependentRules(t *testing.T) {
	want(t, (s.Payment{Card: "4111111111111112", IBAN: "DE89370400440532013001", BIC: "deutdeff"}).Validate(t.Context(), newValidator(t).AtProperty("payment")),
		expected{"payment.card", validation.ErrInvalidLUHN}, expected{"payment.iban", validation.ErrInvalidIBAN}, expected{"payment.bic", validation.ErrInvalidBIC})
	want(t, (s.Payment{Card: "378282246310005"}).Validate(t.Context(), newValidator(t)), expected{"card", validation.ErrInvalidCardScheme})
	want(t, (s.Payment{Card: "x"}).Validate(t.Context(), newValidator(t)), expected{"card", validation.ErrInvalidCardScheme}, expected{"card", validation.ErrInvalidLUHN})
}
func TestPaymentPresenceAndFormatting(t *testing.T) {
	want(t, (s.Payment{}).Validate(t.Context(), newValidator(t)), expected{"card", validation.ErrIsBlank})
	for _, card := range []string{"4111 1111 1111 1111", "4111-1111-1111-1111"} {
		want(t, (s.Payment{Card: card}).Validate(t.Context(), newValidator(t)), expected{"card", validation.ErrInvalidCardScheme}, expected{"card", validation.ErrInvalidLUHN})
	}
}
func TestPaymentValid(t *testing.T) {
	for _, card := range []string{"4111111111111111", "5555555555554444", "2223000048400011"} {
		want(t, (s.Payment{Card: card, IBAN: "DE89370400440532013000", BIC: "DEUTDEFF"}).Validate(t.Context(), newValidator(t)))
		want(t, (s.Payment{Card: card}).Validate(t.Context(), newValidator(t)))
	}
}
