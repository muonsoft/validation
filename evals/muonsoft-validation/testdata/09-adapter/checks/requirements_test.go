package scenario_test

import (
	s "example.com/validation-eval/scenario"
	"github.com/muonsoft/validation"
	"golang.org/x/text/language"
	"golang.org/x/text/message/catalog"
	"testing"
)

func TestAdapterContext(t *testing.T) {
	v, err := validation.NewValidator(validation.Translations(map[language.Tag]map[string]catalog.Message{language.Russian: {s.ErrLimit.Message(): catalog.String("Превышен предел.")}}))
	if err != nil {
		t.Fatal(err)
	}
	v = v.AtProperty("payload").WithGroups("strict").WithLanguage(language.Russian)
	err = s.ValidateDetail(t.Context(), v, s.Legacy{Value: 5}, 3)
	wantViolations(t, err, map[string]error{"payload.detail.value": s.ErrLimit})
	list, _ := validation.UnwrapViolations(err)
	if list.First().Violation().Message() != "Превышен предел." {
		t.Fatal("lost language")
	}
}
func TestAdapterGroupsAndLimit(t *testing.T) {
	v := newValidator(t)
	wantViolations(t, s.ValidateDetail(t.Context(), v, s.Legacy{Value: 5}, 3), nil)
	wantViolations(t, s.ValidateDetail(t.Context(), v.WithGroups("strict"), s.Legacy{Value: 5}, 5), nil)
}
