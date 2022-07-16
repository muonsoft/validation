package validation_test

import (
	"testing"

	"github.com/muonsoft/validation"
	"github.com/stretchr/testify/assert"
)

func TestPropertyPath_String(t *testing.T) {
	tests := []struct {
		path *validation.PropertyPath
		want string
	}{
		{path: nil, want: ""},
		{path: validation.NewPropertyPath(), want: ""},
		{
			path: validation.NewPropertyPath().WithProperty("array").WithIndex(1).WithProperty("property"),
			want: "array[1].property",
		},
		{
			path: validation.NewPropertyPath().WithProperty("@foo").WithProperty("bar"),
			want: "['@foo'].bar",
		},
		{
			path: validation.NewPropertyPath().WithProperty("@foo").WithIndex(0),
			want: "['@foo'][0]",
		},
		{
			path: validation.NewPropertyPath().WithProperty("foo.bar").WithProperty("baz"),
			want: "['foo.bar'].baz",
		},
		{
			path: validation.NewPropertyPath().WithProperty("foo.'bar'").WithProperty("baz"),
			want: `['foo.\'bar\''].baz`,
		},
		{
			path: validation.NewPropertyPath().WithProperty(`0`).WithProperty("baz"),
			want: `['0'].baz`,
		},
		{
			path: validation.NewPropertyPath().WithProperty(`foo[0]`).WithProperty("baz"),
			want: `['foo[0]'].baz`,
		},
		{
			path: validation.NewPropertyPath().WithProperty(``).WithProperty("baz"),
			want: `[''].baz`,
		},
		{
			path: validation.NewPropertyPath().WithProperty(`'`).WithProperty("baz"),
			want: `['\''].baz`,
		},
		{
			path: validation.NewPropertyPath().WithProperty(`\`).WithProperty("baz"),
			want: `['\\'].baz`,
		},
		{
			path: validation.NewPropertyPath().WithProperty(`\'foo`).WithProperty("baz"),
			want: `['\\\'foo'].baz`,
		},
		{
			path: validation.NewPropertyPath().WithProperty(`фу`).WithProperty("baz"),
			want: `фу.baz`,
		},
	}
	for _, test := range tests {
		t.Run(test.want, func(t *testing.T) {
			got := test.path.String()

			assert.Equal(t, test.want, got)
		})
	}
}

func TestPropertyPath_With(t *testing.T) {
	path := validation.NewPropertyPath(validation.PropertyName("top"), validation.ArrayIndex(0))

	path = path.With(
		validation.PropertyName("low"),
		validation.ArrayIndex(1),
		validation.PropertyName("property"),
	)

	assert.Equal(t, "top[0].low[1].property", path.String())
}
