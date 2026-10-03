package scenario_test

import (
	"errors"
	s "example.com/validation-eval/scenario"
	"testing"
)

func TestFormatting(t *testing.T) {
	for _, c := range []struct{ first, last, want string }{{" A ", " B ", "A B"}, {" ", " B ", "B"}, {" A ", " ", "A"}, {" ", " ", "Anonymous"}} {
		if got := s.DisplayName(c.first, c.last); got != c.want {
			t.Fatalf("got %q, want %q", got, c.want)
		}
	}
}
func TestExistingValidation(t *testing.T) {
	if !errors.Is(s.ValidateName(""), s.ErrEmpty) || s.ValidateName(" ") != nil {
		t.Fatal("existing contract changed")
	}
}
