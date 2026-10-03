package scenario_test

import (
	"context"
	"errors"
	s "example.com/validation-eval/scenario"
	"github.com/muonsoft/validation"
	"testing"
)

type updateRepo struct {
	loads, checks, writes int
	name, exclude         string
	saved                 s.Record
	taken                 bool
	err                   error
}

func (r *updateRepo) Load(context.Context, string) (s.Record, error) {
	r.loads++
	return s.Record{ID: "stored-id", Name: "Old"}, r.err
}
func (r *updateRepo) NameTaken(_ context.Context, name, id string) (bool, error) {
	r.checks++
	r.name = name
	r.exclude = id
	return r.taken, nil
}
func (r *updateRepo) Save(_ context.Context, x s.Record) error { r.writes++; r.saved = x; return nil }
func TestNormalizedUpdate(t *testing.T) {
	r := &updateRepo{}
	err := (s.Service{Repository: r}).Update(t.Context(), newValidator(t), "request-id", " New ")
	wantViolations(t, err, nil)
	if r.loads != 1 || r.checks != 1 || r.writes != 1 || r.name != "New" || r.exclude != "stored-id" || r.saved.Name != "New" {
		t.Fatal("wrong state/identity/calls")
	}
}
func TestUpdateBarrier(t *testing.T) {
	r := &updateRepo{}
	err := (s.Service{Repository: r}).Update(t.Context(), newValidator(t), "id", "  ")
	wantViolations(t, err, map[string]error{"name": validation.ErrIsBlank})
	if r.loads != 1 || r.checks != 0 || r.writes != 0 {
		t.Fatal("wrong boundary")
	}
}
func TestDuplicateAndLoadFailure(t *testing.T) {
	r := &updateRepo{taken: true}
	err := (s.Service{Repository: r}).Update(t.Context(), newValidator(t), "id", "New")
	wantViolations(t, err, map[string]error{"name": s.ErrDuplicate})
	if r.writes != 0 {
		t.Fatal("saved duplicate")
	}
	r = &updateRepo{err: context.Canceled}
	err = (s.Service{Repository: r}).Update(t.Context(), newValidator(t), "id", "New")
	if !errors.Is(err, context.Canceled) || r.checks != 0 || r.writes != 0 {
		t.Fatal("masked load error")
	}
}
