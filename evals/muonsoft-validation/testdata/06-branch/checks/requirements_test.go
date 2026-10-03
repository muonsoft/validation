package scenario_test

import (
	"context"
	s "example.com/validation-eval/scenario"
	"github.com/muonsoft/validation"
	"testing"
)

func TestBranchBarrier(t *testing.T) {
	calls := 0
	check := validation.Func[string](func(context.Context, *validation.Validator, string) error { calls++; return nil })
	err := s.Validate(t.Context(), newValidator(t), s.Input{}, check)
	wantViolations(t, err, map[string]error{"name": validation.ErrIsBlank, "reference": validation.ErrIsBlank})
	if calls != 0 {
		t.Fatal("eager argument execution")
	}
}
func TestSiblingIndependence(t *testing.T) {
	calls := 0
	check := validation.Func[string](func(ctx context.Context, v *validation.Validator, _ string) error {
		calls++
		return v.CreateViolation(ctx, validation.ErrNotValid, validation.ErrNotValid.Message())
	})
	err := s.Validate(t.Context(), newValidator(t), s.Input{Reference: "A"}, check)
	wantViolations(t, err, map[string]error{"name": validation.ErrIsBlank, "reference": validation.ErrNotValid})
	if calls != 1 {
		t.Fatal("dependent branch improperly skipped")
	}
}
