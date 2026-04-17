package validate_test

import (
	"errors"
	"testing"

	"github.com/muonsoft/validation/validate"
)

func TestCurrency(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr error
	}{
		{name: "empty", value: "", wantErr: nil},
		{name: "EUR upper", value: "EUR", wantErr: nil},
		{name: "usd lower", value: "usd", wantErr: nil},
		{name: "chf mixed case normalizes", value: "CHf", wantErr: nil},
		{name: "too short", value: "EU", wantErr: validate.ErrInvalidCurrency},
		{name: "too long", value: "EURO", wantErr: validate.ErrInvalidCurrency},
		{name: "unknown", value: "ZZZ", wantErr: validate.ErrInvalidCurrency},
		{name: "digits", value: "123", wantErr: validate.ErrInvalidCurrency},
		{name: "space", value: "EUR ", wantErr: validate.ErrInvalidCurrency},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Currency(tt.value)
			if tt.wantErr == nil {
				if err != nil {
					t.Fatalf("Currency(%q): %v", tt.value, err)
				}
				return
			}
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Currency(%q): got %v, want %v", tt.value, err, tt.wantErr)
			}
		})
	}
}
