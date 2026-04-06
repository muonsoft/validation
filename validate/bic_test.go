package validate_test

import (
	"errors"
	"testing"

	"github.com/muonsoft/validation/validate"
)

func TestBIC(t *testing.T) {
	t.Parallel()

	deIBAN := "DE89370400440532013000"
	gbIBAN := "GB29NWBK60161331926819"

	cases := []struct {
		name        string
		value       string
		options     []func(*validate.BICOptions)
		expectedErr error
	}{
		{name: "empty", value: "", expectedErr: nil},
		{name: "valid 8", value: "DEUTDEFF", expectedErr: nil},
		{name: "valid 11", value: "DEUTDEFF500", expectedErr: nil},
		{name: "spaces stripped", value: "DEUT DE FF", expectedErr: nil},
		{name: "territory JE maps to GB IBAN", value: "AAAAJE22", options: []func(*validate.BICOptions){validate.BICWithIBAN(gbIBAN)}, expectedErr: nil},
		{name: "territory JE mismatch DE IBAN", value: "AAAAJE22", options: []func(*validate.BICOptions){validate.BICWithIBAN(deIBAN)}, expectedErr: validate.ErrBICIBANCountryMismatch},
		{name: "DE BIC matches DE IBAN", value: "DEUTDEFF", options: []func(*validate.BICOptions){validate.BICWithIBAN(deIBAN)}, expectedErr: nil},
		{name: "DE BIC matches lowercase IBAN", value: "DEUTDEFF", options: []func(*validate.BICOptions){validate.BICWithIBAN("de89370400440532013000")}, expectedErr: nil},
		{name: "DE BIC matches spaced lowercase IBAN", value: "DEUTDEFF", options: []func(*validate.BICOptions){validate.BICWithIBAN("de89 3704 0044 0532 0130 00")}, expectedErr: nil},
		{name: "DE BIC mismatch GB IBAN", value: "DEUTDEFF", options: []func(*validate.BICOptions){validate.BICWithIBAN(gbIBAN)}, expectedErr: validate.ErrBICIBANCountryMismatch},
		{name: "DE BIC mismatch spaced lowercase GB IBAN", value: "DEUTDEFF", options: []func(*validate.BICOptions){validate.BICWithIBAN("gb29 nwbk 6016 1331 9268 19")}, expectedErr: validate.ErrBICIBANCountryMismatch},
		{name: "empty IBAN option skips cross-check", value: "DEUTDEFF", options: []func(*validate.BICOptions){validate.BICWithIBAN("")}, expectedErr: nil},
		{name: "too short", value: "DEUTDEF", expectedErr: validate.ErrInvalidBIC},
		{name: "too long", value: "DEUTDEFF5000", expectedErr: validate.ErrInvalidBIC},
		{name: "non alphanumeric", value: "DEUTDE#F", expectedErr: validate.ErrInvalidBIC},
		{name: "unknown country ZZ", value: "DEUTZZFF", expectedErr: validate.ErrInvalidBIC},
		{name: "group region EU", value: "DEUTEUXX", expectedErr: validate.ErrInvalidBIC},
		{name: "strict lowercase", value: "deutdeff", expectedErr: validate.ErrInvalidBIC},
		{name: "case insensitive lowercase", value: "deutdeff", options: []func(*validate.BICOptions){validate.BICCaseInsensitive()}, expectedErr: nil},
		{name: "case insensitive unknown country", value: "deutzzff", options: []func(*validate.BICOptions){validate.BICCaseInsensitive()}, expectedErr: validate.ErrInvalidBIC},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := validate.BIC(tc.value, tc.options...)
			if tc.expectedErr == nil {
				if err != nil {
					t.Fatalf("BIC(%q): got %v, want nil", tc.value, err)
				}
				return
			}
			if !errors.Is(err, tc.expectedErr) {
				t.Fatalf("BIC(%q): got %v, want %v", tc.value, err, tc.expectedErr)
			}
		})
	}
}
