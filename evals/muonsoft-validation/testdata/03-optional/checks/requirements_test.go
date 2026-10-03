package scenario_test

import (
	s "example.com/validation-eval/scenario"
	"github.com/muonsoft/validation"
	"testing"
)

func TestOptional(t *testing.T) {
	n := 0
	v := newValidator(t)
	for _, x := range []s.Input{{Status: "draft"}, {Status: "ready", Lower: &n}, {Status: "ready", Upper: &n}} {
		wantViolations(t, v.ValidateIt(t.Context(), x), nil)
	}
}
func TestIndependentConditions(t *testing.T) {
	blank := ""
	low, high := 3, 2
	x := s.Input{Note: &blank, Lower: &low, Upper: &high}
	wantViolations(t, newValidator(t).ValidateIt(t.Context(), x), map[string]error{"note": validation.ErrIsBlank, "status": validation.ErrNoSuchChoice, "upper": s.ErrRange})
	if blank != "" || low != 3 || high != 2 {
		t.Fatal("input mutated")
	}
	x.Status = "unknown"
	x.Note = nil
	x.Lower = nil
	wantViolations(t, newValidator(t).ValidateIt(t.Context(), x), map[string]error{"status": validation.ErrNoSuchChoice})
}
