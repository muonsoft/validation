package scenario_test

import (
	"testing"

	s "example.com/validation-eval/scenario"
	"github.com/muonsoft/validation"
)

func TestEditorGroups(t *testing.T) {
	x := s.Article{Title: "Title"}
	want(t, x.Validate(t.Context(), newValidator(t)))
	for _, groups := range [][]string{{"publish"}, {validation.DefaultGroup, "publish"}} {
		want(t, x.Validate(t.Context(), newValidator(t).WithGroups(groups...).AtProperty("article")), expected{"article.contact", s.ErrContactRequired})
	}
	want(t, (s.Article{}).Validate(t.Context(), newValidator(t).WithGroups("archive")))
}
func TestEditorAccumulation(t *testing.T) {
	v := newValidator(t).WithGroups("publish").AtProperty("article")
	want(t, (s.Article{Title: "ab", Contact: "invalid"}).Validate(t.Context(), v), expected{"article.title", s.ErrTitleLength}, expected{"article.contact", s.ErrContactFormat})
	want(t, (s.Article{}).Validate(t.Context(), v), expected{"article.title", s.ErrTitleRequired}, expected{"article.contact", s.ErrContactRequired})
}
func TestEditorUnicodeAndMessages(t *testing.T) {
	for _, title := range []string{"猫猫猫", "abcdefghijkl", "   ", "e\u0301a"} {
		want(t, (s.Article{Title: title}).Validate(t.Context(), newValidator(t)))
	}
	for _, title := range []string{"猫猫", "abcdefghijklm"} {
		err := (s.Article{Title: title, Contact: "bad"}).Validate(t.Context(), newValidator(t))
		want(t, err, expected{"title", s.ErrTitleLength}, expected{"contact", s.ErrContactFormat})
		list, _ := validation.UnwrapViolations(err)
		for _, violation := range list.AsSlice() {
			message := s.ErrContactFormat.Message()
			if violation.PropertyPath().String() == "title" {
				message = s.ErrTitleLength.Message()
			}
			if violation.Message() != message {
				t.Fatalf("changed message: %s", violation.Message())
			}
		}
	}
	want(t, (s.Article{Title: "Valid", Contact: "person@example.com"}).Validate(t.Context(), newValidator(t).WithGroups("publish")))
}
