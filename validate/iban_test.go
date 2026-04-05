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

		// Territory alias codes: the pattern uses the parent country prefix,
		// so a valid IBAN for the territory must start with the parent prefix.
		{name: "FI valid (also covers AX alias)", value: "FI2112345600000785"},
		{name: "AX territory always rejected (pattern expects FI prefix)", value: "AX2112345600000785", expectedError: validate.ErrInvalidIBAN},
		{name: "GG territory always rejected (pattern expects GB prefix)", value: "GG29NWBK60161331926819", expectedError: validate.ErrInvalidIBAN},
		{name: "JE territory always rejected (pattern expects GB prefix)", value: "JE29NWBK60161331926819", expectedError: validate.ErrInvalidIBAN},
		{name: "IM territory always rejected (pattern expects GB prefix)", value: "IM29NWBK60161331926819", expectedError: validate.ErrInvalidIBAN},
		{name: "GF territory always rejected (pattern expects FR prefix)", value: "GF7630006000011234567890189", expectedError: validate.ErrInvalidIBAN},

		// Mixed-case with letters in BBAN
		{name: "IT valid lowercase with letter in BBAN", value: "it60x0542811101000000123456"},
		{name: "GB valid mixed case", value: "gb82West12345698765432"},

		{name: "unsupported US", value: "US64SVBX1101057138", expectedError: validate.ErrInvalidIBAN},
		{name: "unknown country XX", value: "XX89370400440532013000", expectedError: validate.ErrInvalidIBAN},
		{name: "DE too short", value: "DE8937040044053201300", expectedError: validate.ErrInvalidIBAN},
		{name: "DE too long", value: "DE893704004405320130000", expectedError: validate.ErrInvalidIBAN},
		{name: "check digits 00", value: "DE00370400440532013000", expectedError: validate.ErrInvalidIBAN},
		{name: "check digits 01", value: "DE01370400440532013000", expectedError: validate.ErrInvalidIBAN},
		{name: "check digits 02 wrong mod 97", value: "DE02370400440532013000", expectedError: validate.ErrInvalidIBAN},
		{name: "check digits 98 wrong mod 97", value: "DE98370400440532013000", expectedError: validate.ErrInvalidIBAN},
		{name: "check digits 99", value: "DE99370400440532013000", expectedError: validate.ErrInvalidIBAN},
		{name: "wrong mod 97", value: "DE89370400440532013001", expectedError: validate.ErrInvalidIBAN},
		{name: "BR format matches Symfony wrong mod 97", value: "BR1800000000140581368018290C1", expectedError: validate.ErrInvalidIBAN},
		{name: "SC wrong length for country pattern", value: "SC18SSCB110110000000000000000USD", expectedError: validate.ErrInvalidIBAN},
		{name: "MU wrong bank code width", value: "MU17BO1M0101101030300200000MUR", expectedError: validate.ErrInvalidIBAN},
		{name: "NBSP between groups", value: "DE89\u00a0370400440532013000"},
		{name: "NNBSP between groups", value: "DE89\u202f370400440532013000"},
		{name: "multiple NBSP", value: "DE89\u00a03704\u00a00044\u00a00532\u00a0013000"},
		{name: "multiple NNBSP", value: "DE89\u202f3704\u202f0044\u202f0532\u202f013000"},
		{name: "mixed spaces", value: "DE89 \u00a03704\u202f00440532013000"},
		{name: "hyphen", value: "DE89-370400440532013000", expectedError: validate.ErrInvalidIBAN},
		{name: "non ascii cyrillic", value: "DE89з70400440532013000", expectedError: validate.ErrInvalidIBAN},
		{name: "only spaces", value: "   ", expectedError: validate.ErrInvalidIBAN},
		{name: "only NBSP", value: "\u00a0\u00a0", expectedError: validate.ErrInvalidIBAN},
		{name: "only NNBSP", value: "\u202f\u202f", expectedError: validate.ErrInvalidIBAN},
		{name: "too short after strip", value: "DE8", expectedError: validate.ErrInvalidIBAN},
		{name: "single char", value: "D", expectedError: validate.ErrInvalidIBAN},
		{name: "two chars", value: "DE", expectedError: validate.ErrInvalidIBAN},
		{name: "three chars", value: "DE8", expectedError: validate.ErrInvalidIBAN},
		{name: "country code not letters", value: "12DE89370400440532013000", expectedError: validate.ErrInvalidIBAN},
		{name: "country lowercase digits", value: "1289370400440532013000", expectedError: validate.ErrInvalidIBAN},
		{name: "check digits not numeric", value: "DEAB370400440532013000", expectedError: validate.ErrInvalidIBAN},
		{name: "DE format letter where digit", value: "DE8937040A440532013000", expectedError: validate.ErrInvalidIBAN},
		{name: "tab", value: "DE89\t370400440532013000", expectedError: validate.ErrInvalidIBAN},
		{name: "newline", value: "DE89\n370400440532013000", expectedError: validate.ErrInvalidIBAN},
		{name: "incomplete UTF-8 NBSP", value: "DE89\xc2", expectedError: validate.ErrInvalidIBAN},
		{name: "incomplete UTF-8 NNBSP", value: "DE89\xe2\x80", expectedError: validate.ErrInvalidIBAN},
		{name: "NNBSP single byte only", value: "DE89\xe2", expectedError: validate.ErrInvalidIBAN},
		{name: "trailing space", value: "DE89370400440532013000 "},
		{name: "leading space", value: " DE89370400440532013000"},
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
