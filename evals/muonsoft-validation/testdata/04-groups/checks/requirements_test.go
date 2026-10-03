package scenario_test

import (
	s "example.com/validation-eval/scenario"
	"github.com/muonsoft/validation"
	"testing"
)

func TestDraft(t *testing.T) {
	v := newValidator(t)
	wantViolations(t, v.ValidateIt(t.Context(), s.Document{Code: "A"}), nil)
	wantViolations(t, v.ValidateIt(t.Context(), s.Document{}), map[string]error{"code": validation.ErrIsBlank})
}
func TestPublish(t *testing.T) {
	v := newValidator(t).WithGroups(validation.DefaultGroup, "publish").AtProperty("document")
	wantViolations(t, v.ValidateIt(t.Context(), s.Document{}), map[string]error{"document.code": validation.ErrIsBlank, "document.title": validation.ErrIsBlank})
	wantViolations(t, v.ValidateIt(t.Context(), s.Document{Code: "A", Title: "B"}), nil)
}
