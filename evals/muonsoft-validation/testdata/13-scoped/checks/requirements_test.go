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
	calls int
	ctx   context.Context
	err   error
	ids   []string
}

func (l *lookup) Find(ctx context.Context, ids []string) (map[string]bool, error) {
	l.calls++
	l.ctx = ctx
	l.ids = append([]string(nil), ids...)
	found := make(map[string]bool)
	for _, id := range ids {
		found[id] = id == "ok"
	}
	return found, l.err
}
func TestBatchIndices(t *testing.T) {
	l := &lookup{}
	ctx := t.Context()
	err := s.CheckMembers(ctx, newValidator(t).AtProperty("team").AtProperty("members"), []string{"missing", "ok", "other", "missing"}, l)
	wantViolations(t, err, map[string]error{"team.members[0]": s.ErrMissing, "team.members[2]": s.ErrMissing, "team.members[3]": s.ErrMissing})
	gotIDs := make(map[string]bool)
	for _, id := range l.ids {
		gotIDs[id] = true
	}
	if !reflect.DeepEqual(gotIDs, map[string]bool{"missing": true, "ok": true, "other": true}) {
		t.Fatalf("wrong lookup IDs: %v", l.ids)
	}
	if l.calls != 1 || l.ctx != ctx {
		t.Fatal("not one batch with original context")
	}
}
func TestEmptyAndFound(t *testing.T) {
	l := &lookup{}
	v := newValidator(t)
	wantViolations(t, s.CheckMembers(t.Context(), v, nil, l), nil)
	if l.calls != 0 {
		t.Fatal("empty lookup")
	}
	wantViolations(t, s.CheckMembers(t.Context(), v, []string{"ok"}, l), nil)
}
func TestBatchFailure(t *testing.T) {
	l := &lookup{err: context.DeadlineExceeded}
	err := s.CheckMembers(t.Context(), newValidator(t), []string{"x"}, l)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatal(err)
	}
	if _, ok := validation.UnwrapViolations(err); ok {
		t.Fatal("technical error converted")
	}
}

func TestBatchPaths(t *testing.T) {
	for _, prefix := range []string{"", "teams[2].members"} {
		t.Run(prefix, func(t *testing.T) {
			v := newValidator(t)
			path := "[1]"
			if prefix != "" {
				v = v.AtProperty("teams").AtIndex(2).AtProperty("members")
				path = prefix + path
			}
			wantViolations(t, s.CheckMembers(t.Context(), v, []string{"ok", "missing"}, &lookup{}), map[string]error{path: s.ErrMissing})
		})
	}
}
