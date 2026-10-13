package scenario_test

import (
	"testing"

	s "example.com/validation-eval/scenario"
	"github.com/muonsoft/validation"
)

func TestPublicationFormats(t *testing.T) {
	x := s.Publication{Country: "EU", Language: "en-US", Locale: "en--US", Currency: "XYZ", Week: "2021-W53", Summary: "!!!"}
	want(t, x.Validate(t.Context(), newValidator(t).AtProperty("payload")),
		expected{"payload.country", validation.ErrInvalidCountry}, expected{"payload.language", validation.ErrInvalidLanguage},
		expected{"payload.locale", validation.ErrInvalidLocale}, expected{"payload.currency", validation.ErrInvalidCurrency},
		expected{"payload.week", validation.ErrInvalidWeekNumber}, expected{"payload.summary", validation.ErrTooFewWords})
}
func TestPublicationBoundaries(t *testing.T) {
	for _, tc := range []struct {
		value string
		code  error
	}{
		{"2020-W52", validation.ErrWeekTooEarly}, {"2021-W03", validation.ErrWeekTooLate},
		{"2020-w53", validation.ErrInvalidWeekFormat}, {"2021-W00", validation.ErrInvalidWeekFormat},
	} {
		t.Run(tc.value, func(t *testing.T) {
			want(t, (s.Publication{Week: tc.value}).Validate(t.Context(), newValidator(t)), expected{"week", tc.code})
		})
	}
	for _, text := range []string{"don't", "’", "   "} {
		want(t, (s.Publication{Summary: text}).Validate(t.Context(), newValidator(t)), expected{"summary", validation.ErrTooFewWords})
	}
	want(t, (s.Publication{Summary: "one two three four five"}).Validate(t.Context(), newValidator(t)), expected{"summary", validation.ErrTooManyWords})
}
func TestPublicationValid(t *testing.T) {
	want(t, (s.Publication{}).Validate(t.Context(), newValidator(t)))
	for _, week := range []string{"2020-W53", "2021-W01", "2021-W02"} {
		for _, summary := range []string{"well-known", "don't panic", "Привет мир", "one two three four", "e\u0301lan vital"} {
			x := s.Publication{Country: "de", Language: "eng", Locale: "zh-Hant-TW", Currency: "EUR", Week: week, Summary: summary}
			want(t, x.Validate(t.Context(), newValidator(t)))
		}
	}
	want(t, (s.Publication{Country: "UK", Language: "und", Locale: "en_US"}).Validate(t.Context(), newValidator(t)))
}
