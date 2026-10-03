package scenario_test

import (
	s "example.com/validation-eval/scenario"
	"github.com/muonsoft/validation"
	"testing"
)

func TestCollectionAccumulation(t *testing.T) {
	err := newValidator(t).ValidateIt(t.Context(), s.Input{Entries: []s.Entry{{Code: "A"}, {Code: "A", Quantity: -1}}})
	wantViolations(t, err, map[string]error{"title": validation.ErrIsBlank, "entries": validation.ErrNotUnique, "entries[0].quantity": validation.ErrNotPositive, "entries[1].quantity": validation.ErrNotPositive})
	list, _ := validation.UnwrapViolations(err)
	for _, v := range list.AsSlice() {
		if v.PropertyPath().String() == "entries[1].quantity" {
			p := v.PropertyPath().Elements()
			if len(p) != 3 || !p[1].IsIndex() {
				t.Fatal("not an array path")
			}
		}
	}
}
func TestCollectionBounds(t *testing.T) {
	v := newValidator(t)
	wantViolations(t, v.ValidateIt(t.Context(), s.Input{Title: "A"}), map[string]error{"entries": validation.ErrTooFewElements})
	entries := []s.Entry{{Code: "A", Quantity: 1}, {Code: "B", Quantity: 1}, {Code: "C", Quantity: 1}, {Code: "D", Quantity: 1}}
	wantViolations(t, v.ValidateIt(t.Context(), s.Input{Title: "A", Entries: entries}), map[string]error{"entries": validation.ErrTooManyElements})
	wantViolations(t, v.ValidateIt(t.Context(), s.Input{Title: "A", Entries: entries[:1]}), nil)
}
