package scenario_test

import (
	"context"
	"errors"
	s "example.com/validation-eval/scenario"

	"github.com/muonsoft/validation"
	"reflect"
	"testing"
)

type lookup struct {
	ctx   context.Context
	calls []string
	fatal error
}

func (l *lookup) Exists(ctx context.Context, id string) (bool, error) {
	if ctx != l.ctx {
		panic("context lost")
	}
	l.calls = append(l.calls, id)
	if id == "fatal" {
		return false, l.fatal
	}
	return id == "ok", nil
}
func TestBranchAccumulation(t *testing.T) {
	l := &lookup{ctx: t.Context()}
	x := s.Request{Items: []s.Item{{}, {ID: "missing", Quantity: 1}, {ID: "ok", Quantity: 2}, {ID: "missing", Quantity: 3}}}
	want(t, x.Validate(t.Context(), newValidator(t).AtProperty("request"), l), expected{"request.items[0].id", validation.ErrIsBlank}, expected{"request.items[0].quantity", validation.ErrNotPositive}, expected{"request.items[1].id", s.ErrMissing}, expected{"request.items[3].id", s.ErrMissing})
	if !reflect.DeepEqual(l.calls, []string{"missing", "ok", "missing"}) {
		t.Fatal(l.calls)
	}

	l.calls = nil
	want(t, (s.Request{Items: []s.Item{{Quantity: 0}}}).Validate(t.Context(), newValidator(t).WithGroups("submit"), l), expected{"items[0].id", validation.ErrIsBlank}, expected{"items[0].quantity", validation.ErrNotPositive})
	if len(l.calls) != 0 {
		t.Fatal(l.calls)
	}
}

func TestBranchFatal(t *testing.T) {
	failure := errors.New("offline")
	l := &lookup{ctx: t.Context(), fatal: failure}
	x := s.Request{Items: []s.Item{{ID: "missing", Quantity: 1}, {ID: "fatal", Quantity: 1}, {ID: "ok", Quantity: 1}}}
	err := x.Validate(t.Context(), newValidator(t), l)
	if !errors.Is(err, failure) {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(l.calls, []string{"missing", "fatal"}) {
		t.Fatal(l.calls)
	}
	if _, ok := validation.UnwrapViolations(err); ok {
		t.Fatal("technical error converted")
	}
}
func TestBranchValid(t *testing.T) {
	l := &lookup{ctx: t.Context()}
	want(t, (s.Request{}).Validate(t.Context(), newValidator(t), l))
	if len(l.calls) != 0 {
		t.Fatal(l.calls)
	}
	x := s.Request{Items: []s.Item{{ID: "ok", Quantity: 1}, {ID: "ok", Quantity: 2}}}
	before := append([]s.Item(nil), x.Items...)
	want(t, x.Validate(t.Context(), newValidator(t).WithGroups("submit"), l))
	if !reflect.DeepEqual(l.calls, []string{"ok", "ok"}) || !reflect.DeepEqual(x.Items, before) {
		t.Fatal(l.calls, x)
	}
}
