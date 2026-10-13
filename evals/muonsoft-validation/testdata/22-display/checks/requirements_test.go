package scenario_test

import (
	"reflect"
	"testing"

	s "example.com/validation-eval/scenario"
)

func TestDisplayGrouping(t *testing.T) {
	input := []s.Issue{{Path: "items[1].name", Message: "A"}, {Path: "", Message: "root"}, {Path: "items[1].name", Message: "A"}, {Path: "name", Message: "B"}}
	got := s.GroupErrors(input)
	want := []s.Group{{Path: "items[1].name", Messages: []string{"A", "A"}}, {Path: "", Messages: []string{"root"}}, {Path: "name", Messages: []string{"B"}}}
	if !reflect.DeepEqual(got, want) {
		t.Fatal(got)
	}
	got[0].Messages[0] = "changed"
	if !reflect.DeepEqual(s.GroupErrors(input), want) || input[0].Message != "A" {
		t.Fatal("shared state")
	}
}
func TestDisplayExistingBehavior(t *testing.T) {
	for _, input := range [][]s.Issue{nil, {}} {
		if s.GroupErrors(input) != nil {
			t.Fatal("empty must be nil")
		}
	}
	want := []s.Issue{{Path: "name", Message: "Required"}}
	if !reflect.DeepEqual(s.ExistingValidate(""), want) || s.ExistingValidate("valid") != nil {
		t.Fatal("validation changed")
	}
}
