package validate_test

import (
	"errors"
	"testing"

	"github.com/muonsoft/validation/is"
	"github.com/muonsoft/validation/validate"
)

func TestUTF8(t *testing.T) {
	cases := []struct {
		name  string
		value string
		valid bool
	}{
		{"empty", "", true},
		{"ASCII", "Hello", true},
		{"multilingual", "Привет 日本語 🙂", true},
		{"replacement character", "\ufffd", true},
		{"NUL and invisible characters", "\x00\u200b", true},
		{"HTML", "<script>alert(1)</script>", true},
		{"mojibake", "Ã©", true},
		{"combining marks", "e\u0301", true},
		{"maximum code point", "\U0010ffff", true},
		{"invalid byte", "\xff", false},
		{"continuation byte", "\x80", false},
		{"truncated two bytes", "\xc2", false},
		{"truncated three bytes", "text\xe2\x82", false},
		{"truncated four bytes", "\xf0\x9f\x99", false},
		{"invalid continuation", "\xc2A", false},
		{"overlong encoding", "\xc0\xaf", false},
		{"surrogate", "\xed\xa0\x80", false},
		{"above Unicode maximum", "\xf4\x90\x80\x80", false},
		{"Latin-1 e acute", "caf\xe9", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validate.UTF8(tc.value)
			if tc.valid && err != nil {
				t.Fatalf("UTF8() = %v, want nil", err)
			}
			if !tc.valid && !errors.Is(err, validate.ErrInvalidUTF8) {
				t.Fatalf("UTF8() = %v, want ErrInvalidUTF8", err)
			}
			if got := is.UTF8(tc.value); got != tc.valid {
				t.Errorf("is.UTF8() = %v, want %v", got, tc.valid)
			}
		})
	}
}
