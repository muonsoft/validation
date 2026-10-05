package scenario_test

import (
	"errors"
	"github.com/muonsoft/validation"
	"testing"
)

type expected struct {
	path string
	code error
}

func want(t *testing.T, err error, expected ...expected) {
	t.Helper()
	if len(expected) == 0 {
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		return
	}
	list, ok := validation.UnwrapViolations(err)
	if !ok || list.Len() != len(expected) {
		t.Fatalf("want %d violations, got %v", len(expected), err)
	}
	used := make([]bool, len(expected))
	for _, item := range list.AsSlice() {
		found := false
		for i, e := range expected {
			if !used[i] && item.PropertyPath().String() == e.path && errors.Is(item, e.code) {
				used[i] = true
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("unexpected violation %s: %v", item.PropertyPath().String(), item)
		}
	}
}
func wantPath(t *testing.T, err error, text string, parts ...validation.PropertyPathElement) {
	t.Helper()
	list, ok := validation.UnwrapViolations(err)
	if !ok {
		t.Fatal(err)
	}
	for _, v := range list.AsSlice() {
		if v.PropertyPath().String() != text {
			continue
		}
		got := v.PropertyPath().Elements()
		if len(got) != len(parts) {
			t.Fatalf("path elements: %v", got)
		}
		for i, p := range parts {
			if got[i].IsIndex() != p.IsIndex() || got[i].String() != p.String() {
				t.Fatalf("wrong segment %d: %v", i, got[i])
			}
		}
		return
	}
	t.Fatalf("missing path %s", text)
}
