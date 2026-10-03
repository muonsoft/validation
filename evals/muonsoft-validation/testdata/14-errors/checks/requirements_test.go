package scenario_test

import (
	"errors"
	"fmt"
	"testing"

	s "example.com/validation-eval/scenario"
)

func TestErrorPresentation(t *testing.T) {
	for _, tc := range []struct {
		err  error
		want string
	}{
		{nil, ""}, {s.ErrInvalid, "Invalid input"},
		{fmt.Errorf("field: %w", s.ErrInvalid), "Invalid input"},
		{errors.New("connection lost"), "connection lost"},
		{errors.New(s.ErrInvalid.Error()), s.ErrInvalid.Error()},
	} {
		if got := s.ErrorText(tc.err); got != tc.want {
			t.Fatalf("got %q, want %q", got, tc.want)
		}
	}
}
func TestExistingCodeValidation(t *testing.T) {
	if !errors.Is(s.ValidateCode(""), s.ErrInvalid) || s.ValidateCode(" ") != nil {
		t.Fatal("existing validation changed")
	}
}
