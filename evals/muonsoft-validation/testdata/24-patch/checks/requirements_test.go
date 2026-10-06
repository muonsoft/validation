package scenario_test

import (
	"errors"
	s "example.com/validation-eval/scenario"

	"github.com/muonsoft/validation"
	"golang.org/x/text/language"
	"golang.org/x/text/message/catalog"
	"testing"
)

func TestPatchStates(t *testing.T) {
	for _, v := range []*validation.Validator{newValidator(t), newValidator(t).WithGroups("activate"), newValidator(t).WithGroups(validation.DefaultGroup, "activate")} {
		for _, n := range []int{0, 1, 2, 3, -1} {
			x := s.Patch{Mode: s.ModePatch{Present: false, Value: &n}}
			want(t, x.Validate(t.Context(), v))
			x.Mode.Present = true
			if n == 1 || n == 2 {
				want(t, x.Validate(t.Context(), v))
			} else {
				want(t, x.Validate(t.Context(), v), expected{"mode", s.ErrMode})
			}
		}
	}
	want(t, (s.Patch{Mode: s.ModePatch{Present: true}}).Validate(t.Context(), newValidator(t)))
	want(t, (s.Patch{Mode: s.ModePatch{Present: true}}).Validate(t.Context(), newValidator(t).WithGroups("activate")), expected{"mode", s.ErrClear})
}
func TestPatchContext(t *testing.T) {
	v, err := validation.NewValidator(validation.Translations(map[language.Tag]map[string]catalog.Message{language.Russian: {s.ErrMode.Message(): catalog.String("Неверный режим.")}}))
	if err != nil {
		t.Fatal(err)
	}
	n := 0
	note := ""
	x := s.Patch{Mode: s.ModePatch{Present: true, Value: &n}, Note: &note}
	err = x.Validate(t.Context(), v.WithLanguage(language.Russian).AtProperty("change"))
	want(t, err, expected{"change.mode", s.ErrMode}, expected{"change.note", validation.ErrIsBlank})
	list, _ := validation.UnwrapViolations(err)
	for _, item := range list.AsSlice() {
		if errors.Is(item, s.ErrMode) && item.Message() != "Неверный режим." {
			t.Fatal(item.Message())
		}
	}
	if n != 0 || note != "" || !x.Mode.Present {
		t.Fatal("mutated")
	}
}
func TestPatchBoundaries(t *testing.T) {
	for _, note := range []string{"界界界界界", " "} {
		want(t, (s.Patch{Note: &note}).Validate(t.Context(), newValidator(t)))
	}
	note := "界界界界界界"
	want(t, (s.Patch{Note: &note, Mode: s.ModePatch{Present: true}}).Validate(t.Context(), newValidator(t).WithGroups("activate")), expected{"note", validation.ErrTooLong}, expected{"mode", s.ErrClear})
}
