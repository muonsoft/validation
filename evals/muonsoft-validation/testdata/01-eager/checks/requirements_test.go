package scenario_test

import (
	s "example.com/validation-eval/scenario"
	"github.com/muonsoft/validation"
	"testing"
)

func TestEager(t *testing.T) {
	var _ validation.Validatable = s.Input{}
	wantViolations(t, newValidator(t).ValidateIt(t.Context(), s.Input{Count: -1}), map[string]error{"name": validation.ErrIsBlank, "count": validation.ErrNotPositive})
}
func TestValid(t *testing.T) {
	wantViolations(t, newValidator(t).ValidateIt(t.Context(), s.Input{Name: "A", Count: 1}), nil)
}
