package scenario_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	s "example.com/validation-eval/scenario"
	"github.com/muonsoft/validation"
)

type repo struct {
	t       *testing.T
	ctx     context.Context
	loaded  *s.Bundle
	calls   []string
	found   map[string]bool
	taken   bool
	failAt  string
	failure error
	saved   *s.Bundle
	ids     []string
}

func (r *repo) call(ctx context.Context, name string) error {
	if ctx != r.ctx {
		r.t.Fatal("context")
	}
	r.calls = append(r.calls, name)
	if r.failAt == name {
		return r.failure
	}
	return nil
}
func (r *repo) Load(ctx context.Context, id string) (*s.Bundle, error) {
	if id != "lookup" {
		r.t.Fatal(id)
	}
	return r.loaded, r.call(ctx, "load")
}
func (r *repo) NameTaken(ctx context.Context, name, id string) (bool, error) {
	if name != "New" || id != "loaded" {
		r.t.Fatalf("unnormalized/wrong identity %q %q", name, id)
	}
	return r.taken, r.call(ctx, "name")
}
func (r *repo) Find(ctx context.Context, ids []string) (map[string]bool, error) {
	r.ids = append([]string(nil), ids...)
	return r.found, r.call(ctx, "find")
}
func (r *repo) Save(ctx context.Context, b *s.Bundle) error { r.saved = b; return r.call(ctx, "save") }
func repository(t *testing.T) *repo {
	return &repo{t: t, ctx: t.Context(), loaded: &s.Bundle{ID: "loaded", Name: "Old", PreviousName: "history", References: []s.Reference{{ID: "old", Weight: 1}}}, found: map[string]bool{"A": true}, failure: errors.New("offline")}
}
func command() s.Command {
	return s.Command{ID: "lookup", Name: " New ", References: []s.Reference{{ID: " A ", Weight: 2}}}
}
func TestUpdateLocalBarrier(t *testing.T) {
	r := repository(t)
	c := command()
	c.Name = " "
	c.References = []s.Reference{{ID: " ", Weight: 0}, {ID: "", Weight: -1}}
	err := (s.Service{Repository: r, Validator: newValidator(t).AtProperty("input")}).Update(t.Context(), c)
	want(t, err, expected{"input.name", validation.ErrIsBlank},
		expected{"input.references[0].id", validation.ErrIsBlank},
		expected{"input.references[0].weight", validation.ErrNotPositive},
		expected{"input.references[1].id", validation.ErrIsBlank},
		expected{"input.references[1].weight", validation.ErrNotPositive})
	if !reflect.DeepEqual(r.calls, []string{"load"}) || r.loaded.Name != "Old" || r.loaded.References[0].ID != "old" {
		t.Fatal(r)
	}
}
func TestUpdateRemoteAccumulation(t *testing.T) {
	r := repository(t)
	r.taken = true
	c := command()
	c.References = []s.Reference{{ID: "A", Weight: 1}, {ID: "X", Weight: 1}, {ID: "X", Weight: 1}}
	err := (s.Service{Repository: r, Validator: newValidator(t).AtProperty("input")}).Update(t.Context(), c)
	want(t, err, expected{"input.name", s.ErrDuplicate},
		expected{"input.references[1].id", s.ErrMissing},
		expected{"input.references[2].id", s.ErrMissing})
	if !reflect.DeepEqual(r.calls, []string{"load", "name", "find"}) {
		t.Fatal(r.calls)
	}
}
func TestUpdateBatchPaths(t *testing.T) {
	for _, prefix := range []string{"references", "input.references"} {
		r := repository(t)
		v := newValidator(t)
		if prefix == "input.references" {
			v = v.AtProperty("input")
		}
		v = v.AtProperty("references")
		refs := []s.Reference{{ID: "X"}, {ID: "A"}, {ID: "X"}}
		err := s.CheckReferences(t.Context(), v, r, refs)
		want(t, err, expected{prefix + "[0].id", s.ErrMissing},
			expected{prefix + "[2].id", s.ErrMissing})
		if !reflect.DeepEqual(r.calls, []string{"find"}) {
			t.Fatal(r.calls)
		}
		set := map[string]bool{}
		for _, id := range r.ids {
			set[id] = true
		}
		if !reflect.DeepEqual(set, map[string]bool{"X": true, "A": true}) {
			t.Fatal(r.ids)
		}
	}
	r := repository(t)
	want(t, s.CheckReferences(t.Context(), newValidator(t), r, nil))
	if len(r.calls) != 0 {
		t.Fatal(r.calls)
	}
}
func TestUpdatePersistenceAndFatal(t *testing.T) {
	for _, stage := range []string{"", "load", "name", "find", "save"} {
		t.Run(stage, func(t *testing.T) {
			r := repository(t)
			r.failAt = stage
			c := command()
			err := (s.Service{Repository: r, Validator: newValidator(t)}).Update(t.Context(), c)
			if stage == "" {
				want(t, err)
				if r.saved == r.loaded || r.saved.Name != "New" || r.saved.References[0].ID != "A" || r.saved.PreviousName != "history" {
					t.Fatal(r.saved)
				}
			} else if !errors.Is(err, r.failure) {
				t.Fatal(err)
			}
			expectedCalls := []string{"load", "name", "find", "save"}
			if stage != "" {
				for i, name := range expectedCalls {
					if name == stage {
						expectedCalls = expectedCalls[:i+1]
						break
					}
				}
			}
			if !reflect.DeepEqual(r.calls, expectedCalls) {
				t.Fatal(r.calls)
			}
			if r.loaded.Name != "Old" || r.loaded.References[0].ID != "old" || c.References[0].ID != " A " {
				t.Fatal("mutated source")
			}
		})
	}
	r := repository(t)
	r.taken = true
	r.failAt = "find"
	err := (s.Service{Repository: r, Validator: newValidator(t)}).Update(t.Context(), command())
	if !errors.Is(err, r.failure) || len(r.calls) != 3 {
		t.Fatal(err, r.calls)
	}
}
