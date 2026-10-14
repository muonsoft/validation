package scenario_test

import (
	"testing"

	s "example.com/validation-eval/scenario"
	"github.com/muonsoft/validation"
)

func patch(value string) s.TextPatch { return s.TextPatch{Present: true, Value: &value} }
func TestProfilePresence(t *testing.T) {
	bad := "invalid"
	for _, x := range []s.ProfilePatch{{}, {Country: s.TextPatch{Present: true}, Locale: s.TextPatch{Present: true}}, {Country: s.TextPatch{Value: &bad}, Locale: s.TextPatch{Value: &bad}}} {
		want(t, x.Validate(t.Context(), newValidator(t)))
	}
	want(t, (s.ProfilePatch{Country: patch(""), Locale: patch("")}).Validate(t.Context(), newValidator(t)), expected{"country", validation.ErrIsBlank}, expected{"locale", validation.ErrIsBlank})
}
func TestProfileFormatBeforePolicy(t *testing.T) {
	for _, tc := range []struct {
		country, locale       string
		countryErr, localeErr error
	}{
		{"EU", "en--US", validation.ErrInvalidCountry, validation.ErrInvalidLocale},
		{"DEU", "en_US.UTF-8", validation.ErrInvalidCountry, validation.ErrInvalidLocale},
		{"FR", "zh-Hant-TW", s.ErrCountryPolicy, s.ErrLocalePolicy},
		{"de", "en_US", s.ErrCountryPolicy, s.ErrLocalePolicy},
	} {
		err := (s.ProfilePatch{Country: patch(tc.country), Locale: patch(tc.locale)}).Validate(t.Context(), newValidator(t).AtProperty("patch"))
		want(t, err, expected{"patch.country", tc.countryErr}, expected{"patch.locale", tc.localeErr})
	}
}
func TestProfileValidAndUnchanged(t *testing.T) {
	for _, country := range []string{"DE", "US", "GB"} {
		for _, locale := range []string{"en-US", "de-DE"} {
			x := s.ProfilePatch{Country: patch(country), Locale: patch(locale)}
			want(t, x.Validate(t.Context(), newValidator(t)))
			if *x.Country.Value != country || *x.Locale.Value != locale {
				t.Fatal("mutated input")
			}
		}
	}
	x := s.ProfilePatch{Country: patch("de"), Locale: patch("en-US")}
	want(t, x.Validate(t.Context(), newValidator(t)), expected{"country", s.ErrCountryPolicy})
	if *x.Country.Value != "de" {
		t.Fatal("normalized country")
	}
}
