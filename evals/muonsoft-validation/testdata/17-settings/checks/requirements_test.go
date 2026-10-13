package scenario_test

import (
	"context"
	"reflect"
	"strconv"
	"testing"

	s "example.com/validation-eval/scenario"
	"github.com/muonsoft/validation"
	"golang.org/x/text/language"
	"golang.org/x/text/message/catalog"
)

type limits struct {
	calls *int
	ctx   context.Context
	max   int
	fail  bool
	t     *testing.T
}

func (l limits) Validate(ctx context.Context, v *validation.Validator, max int) error {
	*l.calls++
	if ctx != l.ctx || max != l.max {
		l.t.Fatal("lost context/parameter")
	}
	if l.fail && !v.IsIgnoredForGroups("publish") {
		return v.BuildViolation(ctx, s.ErrLimit, s.ErrLimit.Message()).WithParameter("{{ maximum }}", strconv.Itoa(max)).AtProperty("value").Create()
	}
	return nil
}
func integer(n int) *int    { return &n }
func note(s string) *string { return &s }
func TestSettingsAccumulation(t *testing.T) {
	calls := 0
	x := s.Document{Schedule: &s.Schedule{Note: note(""), Mode: "wrong", Lower: integer(3), Upper: integer(-1), Limits: limits{&calls, t.Context(), 8, false, t}}}
	err := x.Validate(t.Context(), newValidator(t).AtProperty("request"), 8)
	want(t, err, expected{"request.title", validation.ErrIsBlank},
		expected{"request.schedule.note", validation.ErrIsBlank},
		expected{"request.schedule.mode", validation.ErrNoSuchChoice},
		expected{"request.schedule.upper", validation.ErrNotPositiveOrZero},
		expected{"request.schedule.upper", s.ErrRange})
	if calls != 1 {
		t.Fatal(calls)
	}
}
func TestSettingsOptional(t *testing.T) {
	v := newValidator(t)
	want(t, (s.Document{Title: "A"}).Validate(t.Context(), v, 8))
	for _, bounds := range [][2]*int{{nil, nil}, {integer(0), nil}, {nil, integer(0)}, {integer(1), integer(2)}} {
		calls := 0
		x := s.Document{Title: "A", Schedule: &s.Schedule{Lower: bounds[0], Upper: bounds[1], Limits: limits{&calls, t.Context(), 8, false, t}}}
		copy := *x.Schedule
		want(t, x.Validate(t.Context(), v, 8))
		if !reflect.DeepEqual(copy, *x.Schedule) {
			t.Fatal("mutated")
		}
	}
	calls := 0
	x := s.Document{Title: "A", Schedule: &s.Schedule{Note: note("abcdefghijklmn"), Lower: integer(-1), Limits: limits{&calls, t.Context(), 8, false, t}}}
	want(t, x.Validate(t.Context(), v, 8), expected{"schedule.note", validation.ErrTooLong},
		expected{"schedule.lower", validation.ErrNotPositiveOrZero})
}
func TestSettingsContext(t *testing.T) {
	calls := 0
	v, err := validation.NewValidator(validation.Translations(map[language.Tag]map[string]catalog.Message{language.Russian: {s.ErrLimit.Message(): catalog.String("Предел {{ maximum }}.")}}))
	if err != nil {
		t.Fatal(err)
	}
	x := s.Document{Title: "A", Schedule: &s.Schedule{Limits: limits{&calls, t.Context(), 7, true, t}}}
	err = x.Validate(t.Context(), v.WithGroups(validation.DefaultGroup, "publish").WithLanguage(language.Russian).AtProperty("root"), 7)
	want(t, err, expected{"root.schedule.mode", validation.ErrNoSuchChoice},
		expected{"root.schedule.limits.value", s.ErrLimit})
	list, _ := validation.UnwrapViolations(err)
	for _, item := range list.AsSlice() {
		if item.PropertyPath().String() == "root.schedule.limits.value" && item.Message() != "Предел 7." {
			t.Fatal(item.Message())
		}
	}
}
