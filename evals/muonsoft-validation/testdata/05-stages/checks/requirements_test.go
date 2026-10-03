package scenario_test

import (
	"context"
	s "example.com/validation-eval/scenario"
	"github.com/muonsoft/validation"
	"testing"
)

type repo struct {
	calls, writes int
	found         bool
	ctx           context.Context
}

func (r *repo) Exists(ctx context.Context, _ string) (bool, error) {
	r.calls++
	r.ctx = ctx
	return r.found, nil
}
func (r *repo) Save(ctx context.Context, _ s.Input) error { r.writes++; r.ctx = ctx; return nil }
func TestLocalBarrier(t *testing.T) {
	r := &repo{}
	err := (s.Service{Repository: r}).Handle(t.Context(), newValidator(t), s.Input{})
	wantViolations(t, err, map[string]error{"name": validation.ErrIsBlank, "first": validation.ErrIsBlank, "second": validation.ErrIsBlank})
	if r.calls != 0 || r.writes != 0 {
		t.Fatal("premature I/O")
	}
}
func TestAllRemoteViolations(t *testing.T) {
	r := &repo{}
	err := (s.Service{Repository: r}).Handle(t.Context(), newValidator(t), s.Input{Name: "A", First: "x", Second: "y"})
	wantViolations(t, err, map[string]error{"first": s.ErrMissing, "second": s.ErrMissing})
	if r.calls != 2 || r.writes != 0 {
		t.Fatal("wrong calls")
	}
}
func TestSaveAfterSuccess(t *testing.T) {
	r := &repo{found: true}
	ctx := t.Context()
	err := (s.Service{Repository: r}).Handle(ctx, newValidator(t), s.Input{Name: "A", First: "x", Second: "y"})
	wantViolations(t, err, nil)
	if r.calls != 2 || r.writes != 1 || r.ctx != ctx {
		t.Fatal("wrong calls/context")
	}
}
