package scenario_test

import (
	"errors"
	"testing"

	"github.com/muonsoft/validation"
)

func newValidator(t *testing.T) *validation.Validator {
	t.Helper()
	v, err := validation.NewValidator()
	if err != nil {
		t.Fatal(err)
	}
	return v
}

func wantViolations(t *testing.T, err error, want map[string]error) {
	t.Helper()
	if len(want) == 0 {
		if err != nil {
			t.Fatalf("expected success, got %v", err)
		}
		return
	}
	list, ok := validation.UnwrapViolations(err)
	if !ok || list.Len() != len(want) {
		t.Fatalf("expected %d violations, got %v", len(want), err)
	}
	seen := map[string]bool{}
	for _, violation := range list.AsSlice() {
		path := violation.PropertyPath().String()
		code, exists := want[path]
		if !exists || seen[path] || !errors.Is(violation, code) {
			t.Fatalf("unexpected or duplicate violation: %v", violation)
		}
		seen[path] = true
	}
}
