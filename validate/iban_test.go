package validate_test

import (
	"errors"
	"testing"

	"github.com/muonsoft/validation/validate"
)

func TestIBAN(t *testing.T) {
	tests := []struct {
		name          string
		value         string
		expectedError error
	}{
		{name: "empty", value: ""},
		{name: "DE valid compact", value: "DE89370400440532013000"},
		{name: "DE valid lowercase", value: "de89370400440532013000"},
		{name: "CH valid spaced", value: "CH93 0076 2011 6238 5295 7"},
		{name: "GB valid spaced", value: "GB82 WEST 1234 5698 7654 32"},
		{name: "GB valid compact", value: "GB82WEST12345698765432"},
		{name: "NL valid", value: "NL91ABNA0417164300"},
		{name: "FR valid with letter in BBAN", value: "FR1420041010050500013M02606"},
		{name: "PL valid", value: "PL61109010140000071219812874"},
		{name: "ES valid", value: "ES9121000418450200051332"},
		{name: "BE valid", value: "BE68539007547034"},
		{name: "AT valid", value: "AT611904300234573201"},
		{name: "IT valid with X in BBAN", value: "IT60X0542811101000000123456"},
		{name: "SE valid", value: "SE4550000000058398257466"},
		{name: "BR valid trailing bank branch letters", value: "BR1300000000140581368018290C1"},
		{name: "MU valid complex BBAN", value: "MU17BOMM0101101030300200000MUR"},
		{name: "SC valid bank code and currency suffix", value: "SC62SSCB11011000000000000000USD"},
		{name: "unsupported US", value: "US64SVBX1101057138", expectedError: validate.ErrInvalidIBAN},
		{name: "unknown country XX", value: "XX89370400440532013000", expectedError: validate.ErrInvalidIBAN},
		{name: "DE too short", value: "DE8937040044053201300", expectedError: validate.ErrInvalidIBAN},
		{name: "DE too long", value: "DE893704004405320130000", expectedError: validate.ErrInvalidIBAN},
		{name: "check digits 01", value: "DE01370400440532013000", expectedError: validate.ErrInvalidIBAN},
		{name: "check digits 99", value: "DE99370400440532013000", expectedError: validate.ErrInvalidIBAN},
		{name: "wrong mod 97", value: "DE89370400440532013001", expectedError: validate.ErrInvalidIBAN},
		{name: "BR format matches Symfony wrong mod 97", value: "BR1800000000140581368018290C1", expectedError: validate.ErrInvalidIBAN},
		{name: "SC wrong length for country pattern", value: "SC18SSCB110110000000000000000USD", expectedError: validate.ErrInvalidIBAN},
		{name: "MU wrong bank code width", value: "MU17BO1M0101101030300200000MUR", expectedError: validate.ErrInvalidIBAN},
		{name: "NBSP between groups", value: "DE89\u00a0370400440532013000"},
		{name: "NNBSP between groups", value: "DE89\u202f370400440532013000"},
		{name: "hyphen", value: "DE89-370400440532013000", expectedError: validate.ErrInvalidIBAN},
		{name: "non ascii cyrillic", value: "DE89з70400440532013000", expectedError: validate.ErrInvalidIBAN},
		{name: "only spaces", value: "   ", expectedError: validate.ErrInvalidIBAN},
		{name: "only NBSP", value: "\u00a0\u00a0", expectedError: validate.ErrInvalidIBAN},
		{name: "too short after strip", value: "DE8", expectedError: validate.ErrInvalidIBAN},
		{name: "country code not letters", value: "12DE89370400440532013000", expectedError: validate.ErrInvalidIBAN},
		{name: "check digits not numeric", value: "DEAB370400440532013000", expectedError: validate.ErrInvalidIBAN},
		{name: "DE format letter where digit", value: "DE8937040A440532013000", expectedError: validate.ErrInvalidIBAN},
		{name: "tab", value: "DE89\t370400440532013000", expectedError: validate.ErrInvalidIBAN},
		{name: "incomplete UTF-8 NBSP", value: "DE89\xc2", expectedError: validate.ErrInvalidIBAN},
		{name: "incomplete UTF-8 NNBSP", value: "DE89\xe2\x80", expectedError: validate.ErrInvalidIBAN},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.IBAN(tt.value)
			if tt.expectedError == nil {
				if err != nil {
					t.Fatalf("IBAN(%q): %v", tt.value, err)
				}
				return
			}
			if !errors.Is(err, tt.expectedError) {
				t.Fatalf("IBAN(%q): got %v, want %v", tt.value, err, tt.expectedError)
			}
		})
	}
}
