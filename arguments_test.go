package validation_test

import (
	"context"
	"testing"

	"github.com/muonsoft/validation"
	"github.com/stretchr/testify/require"
)

func argumentDeriver[T validation.Argument](at func(...validation.PropertyPathElement) T) func(...validation.PropertyPathElement) validation.Argument {
	return func(path ...validation.PropertyPathElement) validation.Argument { return at(path...) }
}

func TestArgumentAt_CopiesAreIndependent(t *testing.T) {
	validator := newValidator(t)
	a, b, c := validation.PropertyName("a"), validation.PropertyName("b"), validation.PropertyName("c")
	invalid := validation.Check(false)
	cases := []struct {
		name   string
		derive func(...validation.PropertyPathElement) validation.Argument
	}{
		{"checker", argumentDeriver(invalid.At(a).At(b).At(c).At)},
		{"validator", argumentDeriver(validation.Valid(validation.ValidatableFunc(func(ctx context.Context, v *validation.Validator) error { return v.Validate(ctx, invalid) })).At(a).At(b).At(c).At)},
		{"when", argumentDeriver(validation.When(true).Then(invalid).At(a).At(b).At(c).At)},
		{"groups", argumentDeriver(validation.WhenGroups().Then(invalid).At(a).At(b).At(c).At)},
		{"sequential", argumentDeriver(validation.Sequentially(invalid).At(a).At(b).At(c).At)},
		{"at least one", argumentDeriver(validation.AtLeastOneOf(invalid).At(a).At(b).At(c).At)},
		{"all", argumentDeriver(validation.All(invalid).At(a).At(b).At(c).At)},
		{"async", argumentDeriver(validation.Async(invalid).At(a).At(b).At(c).At)},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			left := tc.derive(validation.PropertyName("left"))
			right := tc.derive(validation.PropertyName("right"))
			for _, item := range []struct {
				arg  validation.Argument
				path string
			}{{left, "a.b.c.left"}, {right, "a.b.c.right"}} {
				list, ok := validation.UnwrapViolations(validator.Validate(context.Background(), item.arg))
				require.True(t, ok)
				require.Equal(t, item.path, list.First().PropertyPath().String())
			}
		})
	}
}
