package validation_test

import (
	"strings"
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
			path: validation.NewPropertyPath().WithProperty("foo").WithProperty(""),
			want: `foo['']`,
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

func BenchmarkPropertyPath_String(b *testing.B) {
	// cpu: Intel(R) Core(TM) i9-9900K CPU @ 3.60GHz
	// BenchmarkPropertyPath_String
	// BenchmarkPropertyPath_String-16    	  217718	      7428 ns/op	    4384 B/op	       5 allocs/op
	path := validation.NewPropertyPath(
		validation.PropertyName("array"),
		validation.ArrayIndex(1234567890),
		validation.PropertyName("@foo"),
		validation.ArrayIndex(1234567890),
		validation.PropertyName("foo.bar"),
		validation.PropertyName("foo.'bar'"),
		validation.PropertyName(`0123456789`),
		validation.PropertyName(`foo[0][1][2][3][4][5][6][7][8][9]`),
		validation.PropertyName(``),
		validation.PropertyName(`'`),
		validation.PropertyName(`\`),
		validation.PropertyName(`\'foo`),
		validation.PropertyName(`фу`),
		validation.PropertyName(strings.Repeat(`@foo.'bar'.[baz]`, 100)),
	)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = path.String()
	}
}
