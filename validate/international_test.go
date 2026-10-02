package validate_test

import (
	"errors"
	"testing"

	"github.com/muonsoft/validation/is"
	"github.com/muonsoft/validation/validate"
)

func TestInternational(t *testing.T) {
	checks := []struct {
		name         string
		check        func(string) error
		predicate    func(string) bool
		invalidError error
		valid        []string
		invalid      []string
	}{
		{
			name: "Country", check: validate.Country, predicate: is.Country, invalidError: validate.ErrInvalidCountry,
			valid:   []string{"", "DE", "de", "dE", "FR", "US", "AX", "XK", "UK", "SU", "AN", "UN"},
			invalid: []string{"D", "DEU", "276", "001", "EU", "ZZ", "AA", "QQ", "12", "!@", "DE-DE", " DE", "DE ", "\t", "ДЕ", "�", "\x00"},
		},
		{
			name: "Language", check: validate.Language, predicate: is.Language, invalidError: validate.ErrInvalidLanguage,
			valid:   []string{"", "en", "EN", "eN", "eng", "ENG", "de", "deu", "ger", "iw", "und", "UND", "mul", "zxx", "qaa", "qtz"},
			invalid: []string{"e", "engl", "zzz", "qzz", "en-US", "zh-Hant", "en_US", "English", "123", " en", "en ", "\n", "中文", "\x00"},
		},
		{
			name: "Locale", check: validate.Locale, predicate: is.Locale, invalidError: validate.ErrInvalidLocale,
			valid:   []string{"", "en", "eng", "en-US", "EN-us", "en_US", "zh-Hant-TW", "sr-Latn", "de-CH-1901", "es-419", "en-001", "en-999", "und", "root", "und-Latn", "und-US", "qaa", "x-private", "en-x-private", "en-u-ca-gregory", "en-u-ca-foobar", "i-klingon", "iw-IL", "en-ZZ", "en-XX"},
			invalid: []string{"!", "zzz", "en-Abcd", "en-foobar", "en--US", "en_", "en-", "en-u", "x", "x-", "en-US.UTF-8", "en@calendar=gregorian", "en,fr;q=0.8", " en", "en ", "\n", "русский", "\x00"},
		},
	}
	for _, check := range checks {
		t.Run(check.name, func(t *testing.T) {
			for _, value := range check.valid {
				t.Run("valid/"+value, func(t *testing.T) {
					if err := check.check(value); err != nil {
						t.Errorf("check(%q) = %v, want nil", value, err)
					}
					if !check.predicate(value) {
						t.Errorf("predicate(%q) = false, want true", value)
					}
				})
			}
			for _, value := range check.invalid {
				t.Run("invalid/"+value, func(t *testing.T) {
					if err := check.check(value); !errors.Is(err, check.invalidError) {
						t.Errorf("check(%q) = %v, want %v", value, err, check.invalidError)
					}
					if check.predicate(value) {
						t.Errorf("predicate(%q) = true, want false", value)
					}
				})
			}
		})
	}
}
