package scenario_test

import (
	s "example.com/validation-eval/scenario"

	"github.com/muonsoft/validation"
	"reflect"
	"testing"
)

func TestCatalogAccumulation(t *testing.T) {
	x := s.Catalog{Roots: map[string]*s.Node{"a": {Children: map[string]*s.Node{"b": {}}}, "z": nil}}
	want(t, x.Validate(t.Context(), newValidator(t).AtProperty("payload"), 2), expected{"payload.roots.a.label", validation.ErrIsBlank}, expected{"payload.roots.a.children.b.label", validation.ErrIsBlank}, expected{"payload.roots.z", s.ErrMissing})
}
func TestCatalogPaths(t *testing.T) {
	for _, key := range []string{"", "0", "a.b", "a/b~c", "a'b", "a\\b"} {
		x := s.Catalog{Roots: map[string]*s.Node{key: {}}}
		err := x.Validate(t.Context(), newValidator(t).AtProperty("payload"), 0)
		parts := []validation.PropertyPathElement{validation.PropertyName("payload"), validation.PropertyName("roots"), validation.PropertyName(key), validation.PropertyName("label")}
		path := validation.NewPropertyPath(parts...).String()
		want(t, err, expected{path, validation.ErrIsBlank})
		wantPath(t, err, path, parts...)
	}
}
func TestCatalogDepth(t *testing.T) {
	leaf := &s.Node{Label: "same"}
	x := s.Catalog{Roots: map[string]*s.Node{"a": {Label: "same", Children: map[string]*s.Node{"b": leaf}}, "z": {}}}
	want(t, x.Validate(t.Context(), newValidator(t), 0), expected{"roots.a.children.b", s.ErrDepth}, expected{"roots.z.label", validation.ErrIsBlank})
	want(t, x.Validate(t.Context(), newValidator(t), 1), expected{"roots.z.label", validation.ErrIsBlank})
	if !reflect.DeepEqual(*leaf, s.Node{Label: "same"}) {
		t.Fatal("mutated")
	}
	blocked := s.Catalog{Roots: map[string]*s.Node{"root": {Label: "ok", Children: map[string]*s.Node{"deep": {Children: map[string]*s.Node{"bad": nil}}}}}}
	want(t, blocked.Validate(t.Context(), newValidator(t), 0), expected{"roots.root.children.deep", s.ErrDepth})

	want(t, (s.Catalog{}).Validate(t.Context(), newValidator(t), 0))
	want(t, (s.Catalog{Roots: map[string]*s.Node{"a": {Label: "same"}, "b": {Label: "same"}}}).Validate(t.Context(), newValidator(t), 0))
}
