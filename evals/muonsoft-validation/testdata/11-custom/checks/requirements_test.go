package scenario_test

import (
	s "example.com/validation-eval/scenario"
	"github.com/muonsoft/validation"
	"golang.org/x/text/language"
	"testing"
)

func TestCustomRule(t *testing.T) {
	var _ validation.StringConstraint = s.PrefixConstraint{}
	v, err := s.NewValidator()
	if err != nil {
		t.Fatal(err)
	}
	c := s.PrefixConstraint{Prefix: "REF-"}
	for _, value := range []*string{nil, ptr(""), ptr("REF-1")} {
		wantViolations(t, v.Validate(t.Context(), validation.NilStringProperty("code", value, c)), nil)
	}
	err = v.Validate(t.Context(), validation.StringProperty("code", "bad", c))
	wantViolations(t, err, map[string]error{"code": s.ErrPrefix})
	list, _ := validation.UnwrapViolations(err)
	if list.First().Violation().Message() != "Value must start with REF-." {
		t.Fatal("wrong fallback")
	}
}
func TestTranslation(t *testing.T) {
	v, err := s.NewValidator()
	if err != nil {
		t.Fatal(err)
	}
	err = v.WithLanguage(language.Russian).AtProperty("input").Validate(t.Context(), validation.StringProperty("code", "bad", s.PrefixConstraint{Prefix: "X-"}))
	wantViolations(t, err, map[string]error{"input.code": s.ErrPrefix})
	list, _ := validation.UnwrapViolations(err)
	if list.First().Violation().Message() != "Значение должно начинаться с X-." {
		t.Fatal("wrong translation")
	}
}
func ptr(s string) *string { return &s }
