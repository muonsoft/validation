package validate_test

import (
	"errors"
	"testing"

	"github.com/muonsoft/validation/validate"
)

func TestLuhn(t *testing.T) {
	t.Parallel()

	tests := []struct {
		value         string
		expectedError error
	}{
		{value: "", expectedError: nil},
		{value: "79927398713", expectedError: nil},
		{value: "79927398710", expectedError: validate.ErrInvalidChecksum},
		{value: "0000000000000000", expectedError: validate.ErrInvalidChecksum},
		{value: "12345a", expectedError: validate.ErrContainsNonDigit},
		{value: " 79927398713", expectedError: validate.ErrContainsNonDigit},
	}

	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			t.Parallel()

			err := validate.Luhn(tt.value)
			if !errors.Is(err, tt.expectedError) {
				t.Fatalf("Luhn(%q): got error %v, want %v", tt.value, err, tt.expectedError)
			}
		})
	}
}
