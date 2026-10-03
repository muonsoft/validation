package scenario_test

import (
	"context"
	"errors"
	s "example.com/validation-eval/scenario"
	"github.com/muonsoft/validation"
	"testing"
)

type lookup struct {
	calls int
	ctx   context.Context
	err   error
}

func (l *lookup) Find(ctx context.Context, ids []string) (map[string]bool, error) {
	l.calls++
	l.ctx = ctx
	return map[string]bool{"ok": true}, l.err
}
func TestBatchIndices(t *testing.T) {
	l := &lookup{}
	ctx := t.Context()
	err := s.CheckReferences(ctx, newValidator(t).AtProperty("input"), []string{"missing", "ok", "other", "missing"}, l)
	wantViolations(t, err, map[string]error{"input.references[0]": s.ErrMissing, "input.references[2]": s.ErrMissing, "input.references[3]": s.ErrMissing})
	if l.calls != 1 || l.ctx != ctx {
		t.Fatal("not one batch with original context")
	}
}
func TestEmptyAndFound(t *testing.T) {
	l := &lookup{}
	v := newValidator(t)
	wantViolations(t, s.CheckReferences(t.Context(), v, nil, l), nil)
	if l.calls != 0 {
		t.Fatal("empty lookup")
	}
	wantViolations(t, s.CheckReferences(t.Context(), v, []string{"ok"}, l), nil)
}
func TestBatchFailure(t *testing.T) {
	l := &lookup{err: context.DeadlineExceeded}
	err := s.CheckReferences(t.Context(), newValidator(t), []string{"x"}, l)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatal(err)
	}
	if _, ok := validation.UnwrapViolations(err); ok {
		t.Fatal("technical error converted")
	}
}
