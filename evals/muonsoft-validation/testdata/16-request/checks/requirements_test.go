package scenario_test

import (
	"reflect"
	"strings"
	"testing"

	s "example.com/validation-eval/scenario"
	"github.com/muonsoft/validation"
)

func goodRequest() s.Request {
	return s.Request{Title: "Visit", Contact: s.Contact{Name: "Ada", Email: "ada@example.com"}, Visits: []s.Visit{{Address: "A", Tasks: []s.Task{{ServiceCode: "S", Variant: "basic", Duration: 1}}}}}
}
func TestRequestAccumulation(t *testing.T) {
	x := goodRequest()
	x.Title = ""
	x.Contact = s.Contact{}
	x.Visits[0].Address = ""
	x.Visits[0].Tasks = []s.Task{{ServiceCode: "S", Variant: "basic"}, {ServiceCode: "S", Variant: "basic", Duration: -1}, {Variant: "other"}}
	err := newValidator(t).AtProperty("request").ValidateIt(t.Context(), x)
	want(t, err, expected{"request.title", validation.ErrIsBlank},
		expected{"request.contact.name", validation.ErrIsBlank},
		expected{"request.contact.email", validation.ErrIsBlank},
		expected{"request.visits[0].address", validation.ErrIsBlank},
		expected{"request.visits[0].tasks", validation.ErrTooManyElements},
		expected{"request.visits[0].tasks[0]", validation.ErrNotUnique},
		expected{"request.visits[0].tasks[1]", validation.ErrNotUnique},
		expected{"request.visits[0].tasks[0].duration", validation.ErrNotPositive},
		expected{"request.visits[0].tasks[1].duration", validation.ErrNotPositive},
		expected{"request.visits[0].tasks[2].serviceCode", validation.ErrIsBlank},
		expected{"request.visits[0].tasks[2].variant", validation.ErrNoSuchChoice},
		expected{"request.visits[0].tasks[2].duration", validation.ErrNotPositive})
}
func TestRequestPaths(t *testing.T) {
	x := goodRequest()
	x.Visits = append(x.Visits, x.Visits[0])
	x.Visits[1].Tasks = []s.Task{{ServiceCode: "A", Variant: "basic", Duration: 1}, {ServiceCode: "B", Variant: "basic"}}
	err := newValidator(t).AtProperty("request").ValidateIt(t.Context(), x)
	want(t, err, expected{"request.visits[1].tasks[1].duration", validation.ErrNotPositive})
	wantPath(t, err, "request.visits[1].tasks[1].duration", validation.PropertyName("request"), validation.PropertyName("visits"), validation.ArrayIndex(1), validation.PropertyName("tasks"), validation.ArrayIndex(1), validation.PropertyName("duration"))
	want(t, newValidator(t).ValidateIt(t.Context(), x.Visits[1].Tasks[1]), expected{"duration", validation.ErrNotPositive})
}
func TestRequestBoundsAndValid(t *testing.T) {
	v := newValidator(t)
	x := goodRequest()
	x.Visits = append(x.Visits, x.Visits[0])
	before := goodRequest()
	before.Visits = append(before.Visits, before.Visits[0])
	want(t, v.ValidateIt(t.Context(), x))
	if !reflect.DeepEqual(x, before) {
		t.Fatal("mutated input")
	}
	x.Visits = nil
	x.Title = strings.Repeat("界", 41)
	want(t, v.ValidateIt(t.Context(), x), expected{"title", validation.ErrTooLong},
		expected{"visits", validation.ErrTooFewElements})
	x = goodRequest()
	x.Visits[0].Tasks = nil
	want(t, v.ValidateIt(t.Context(), x), expected{"visits[0].tasks", validation.ErrTooFewElements})
	x = goodRequest()
	x.Contact.Email = "bad"
	x.Visits[0].Tasks[0].Variant = ""
	want(t, v.ValidateIt(t.Context(), x), expected{"contact.email", validation.ErrInvalidEmail},
		expected{"visits[0].tasks[0].variant", validation.ErrNoSuchChoice})
}
