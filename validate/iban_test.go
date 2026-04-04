package validate_test

import (
	"errors"
	"testing"

	"github.com/muonsoft/validation/validate"
)

func TestIBAN(t *testing.T) {
	tests := []struct {
		value         string
		expectedError error
	}{
		{value: ""},
		{value: "DE89370400440532013000"},
		{value: "de89370400440532013000"},
		{value: "CH93 0076 2011 6238 5295 7"},
		{value: "GB82 WEST 1234 5698 7654 32"},
		{value: "GB82WEST12345698765432"},
		{value: "NL91ABNA0417164300"},
		{value: "FR1420041010050500013M02606"},
		{value: "US64SVBX1101057138", expectedError: validate.ErrInvalidIBAN},
		{value: "XX89370400440532013000", expectedError: validate.ErrInvalidIBAN},
		{value: "DE8937040044053201300", expectedError: validate.ErrInvalidIBAN},
		{value: "DE893704004405320130000", expectedError: validate.ErrInvalidIBAN},
		{value: "DE01370400440532013000", expectedError: validate.ErrInvalidIBAN},
		{value: "DE99370400440532013000", expectedError: validate.ErrInvalidIBAN},
		{value: "DE89370400440532013001", expectedError: validate.ErrInvalidIBAN},
		{value: "DE89\u00a0370400440532013000"},
		{value: "DE89\u202f370400440532013000"},
		{value: "DE89-370400440532013000", expectedError: validate.ErrInvalidIBAN},
		{value: "DE89з70400440532013000", expectedError: validate.ErrInvalidIBAN},
	}
	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
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
