package scenario_test

import (
	"testing"

	s "example.com/validation-eval/scenario"
	"github.com/muonsoft/validation"
)

func leaf(name string) s.Field { return s.Field{Name: name, Kind: "text", Value: "ok"} }
func TestTreeAccumulation(t *testing.T) {
	x := s.Form{Fields: []s.Field{{Name: "G", Kind: "group", Children: []s.Field{{Name: "A", Kind: "text"}, {Name: "A", Kind: "text"}, {Name: "", Kind: "invalid"}, leaf("B")}}, {Name: "", Kind: "text"}}}
	want(t, x.Validate(t.Context(), newValidator(t).AtProperty("form"), 3), expected{"form.fields[0].children", validation.ErrTooManyElements},
		expected{"form.fields[0].children[0].name", validation.ErrNotUnique},
		expected{"form.fields[0].children[1].name", validation.ErrNotUnique},
		expected{"form.fields[0].children[0].value", validation.ErrIsBlank},
		expected{"form.fields[0].children[1].value", validation.ErrIsBlank},
		expected{"form.fields[0].children[2].name", validation.ErrIsBlank},
		expected{"form.fields[0].children[2].kind", validation.ErrNoSuchChoice},
		expected{"form.fields[1].name", validation.ErrIsBlank},
		expected{"form.fields[1].value", validation.ErrIsBlank})
}
func TestTreePaths(t *testing.T) {
	for _, tc := range []struct{ key, path string }{{"a.b", "['a.b']"}, {"a/b~c", "['a/b~c']"}, {"0", "['0']"}, {"", "['']"}, {"a'b", "['a\\'b']"}, {`a\b`, `['a\\b']`}} {
		t.Run(tc.key, func(t *testing.T) {
			f := leaf("A")
			f.Attributes = map[string]string{tc.key: ""}
			x := s.Form{Fields: []s.Field{{Name: "G", Kind: "group", Children: []s.Field{leaf("B"), f}}}}
			prefix := "root.fields[0].children[1].attributes"
			err := x.Validate(t.Context(), newValidator(t).AtProperty("root"), 3)
			want(t, err, expected{prefix + tc.path, validation.ErrIsBlank})
			wantPath(t, err, prefix+tc.path, validation.PropertyName("root"), validation.PropertyName("fields"), validation.ArrayIndex(0), validation.PropertyName("children"), validation.ArrayIndex(1), validation.PropertyName("attributes"), validation.PropertyName(tc.key))
		})
	}
}
func TestTreeDepthAndScope(t *testing.T) {
	v := newValidator(t)
	x := s.Form{Fields: []s.Field{{Name: "A", Kind: "group", Children: []s.Field{leaf("A")}}, leaf("B")}}
	want(t, x.Validate(t.Context(), v, 2))
	x.Fields[1].Value = ""
	want(t, x.Validate(t.Context(), v, 1), expected{"fields[0].children[0]", s.ErrDepth},
		expected{"fields[1].value", validation.ErrIsBlank})
	want(t, (s.Form{}).Validate(t.Context(), v, 2), expected{"fields", validation.ErrTooFewElements})
	x = s.Form{Fields: []s.Field{{Name: "A", Kind: "group"}}}
	want(t, x.Validate(t.Context(), v, 2), expected{"fields[0].children", validation.ErrTooFewElements})
}
