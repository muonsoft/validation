package scenario_test

import (
	"reflect"
	"testing"

	s "example.com/validation-eval/scenario"
)

func TestIndependentCopy(t *testing.T) {
	original := map[string][]string{"items[2].name": {"required", "too short"}, "none": nil, "empty": {}}
	copied := s.CopyMessages(original)
	if !reflect.DeepEqual(copied, original) {
		t.Fatalf("copy changed values: %#v", copied)
	}
	copied["items[2].name"][0] = "changed"
	delete(copied, "none")
	copied["new"] = []string{"added"}
	if original["items[2].name"][0] != "required" {
		t.Fatal("shared slice")
	}
	if _, ok := original["none"]; !ok {
		t.Fatal("shared map")
	}
	if _, ok := original["new"]; ok {
		t.Fatal("shared map")
	}
	original["items[2].name"][1] = "source changed"
	if copied["items[2].name"][1] != "too short" {
		t.Fatal("shared slice")
	}
}
func TestNilAndExistingValidation(t *testing.T) {
	if s.CopyMessages(nil) != nil {
		t.Fatal("nil map changed")
	}
	if s.CopyMessages(map[string][]string{}) == nil {
		t.Fatal("empty map became nil")
	}
	if !s.ValidateCount(0) || !s.ValidateCount(2) || s.ValidateCount(-1) {
		t.Fatal("existing validation changed")
	}
}
