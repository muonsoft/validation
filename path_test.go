package validation_test

import (
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/muonsoft/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPropertyPath_With(t *testing.T) {
	path := validation.NewPropertyPath(validation.PropertyName("top"), validation.ArrayIndex(0))

	path = path.With(
		validation.PropertyName("low"),
		validation.ArrayIndex(1),
		validation.PropertyName("property"),
	)

	assert.Equal(t, "top[0].low[1].property", path.String())
}

func TestPropertyPath_Elements(t *testing.T) {
	want := []validation.PropertyPathElement{
		validation.PropertyName("top"),
		validation.ArrayIndex(0),
		validation.PropertyName("low"),
		validation.ArrayIndex(1),
		validation.PropertyName("property"),
	}

	got := validation.NewPropertyPath(want...).Elements()

	assert.Equal(t, want, got)
}

func TestPropertyPath_String(t *testing.T) {
	tests := []struct {
		path *validation.PropertyPath
		want string
	}{
		{path: nil, want: ""},
		{path: validation.NewPropertyPath(), want: ""},
		{path: validation.NewPropertyPath(validation.PropertyName(" ")), want: "[' ']"},
		{path: validation.NewPropertyPath(validation.PropertyName("$")), want: "$"},
		{path: validation.NewPropertyPath(validation.PropertyName("_")), want: "_"},
		{path: validation.NewPropertyPath(validation.PropertyName("id$_")), want: "id$_"},
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

func TestPropertyPath_UnmarshalText(t *testing.T) {
	tests := []struct {
		pathString string
		want       []validation.PropertyPathElement
		wantError  string
	}{
		{pathString: "[", wantError: "parsing path element #0: incomplete array index"},
		{pathString: "]", wantError: "parsing path element #0 at char #0 ']': unexpected close bracket"},
		{pathString: "[]", wantError: "parsing path element #0 at char #1 ']': unexpected close bracket"},
		{pathString: "'", wantError: `parsing path element #0 at char #0 '\'': unexpected quote`},
		{pathString: ".", wantError: "parsing path element #0 at char #0 '.': unexpected point"},
		{pathString: "\\", wantError: `parsing path element #0 at char #0 '\\': unexpected backslash`},
		{pathString: "[[", wantError: "parsing path element #0 at char #1 '[': unexpected char"},
		{pathString: "0", wantError: "parsing path element #0 at char #0 '0': unexpected identifier character"},
		{pathString: "[0", wantError: "parsing path element #0: incomplete array index"},
		{pathString: "[0[", wantError: "parsing path element #0 at char #2 '[': unexpected char"},
		{pathString: "[0[]", wantError: "parsing path element #0 at char #2 '[': unexpected char"},
		{pathString: "[[0]", wantError: "parsing path element #0 at char #1 '[': unexpected char"},
		{pathString: "['property", wantError: "parsing path element #0: incomplete bracketed property name"},
		{pathString: "['property'", wantError: "parsing path element #0: incomplete bracketed property name"},
		{pathString: "['property]", wantError: "parsing path element #0: incomplete bracketed property name"},
		{pathString: "['property'][", wantError: "parsing path element #1: incomplete array index"},
		{pathString: "[0][1][invalid]", wantError: "parsing path element #2 at char #7 'i': unexpected array index character"},
		{pathString: "[0][1][0invalid]", wantError: "parsing path element #2 at char #8 'i': unexpected array index character"},
		{pathString: "[0][1][012345678901234567890123456789]", wantError: "parsing path element #2: value out of range: 012345678901234567890123456789"},
		{pathString: "[9227000000000000000]", wantError: "parsing path element #0: value out of range: 9227000000000000000"},
		{pathString: " ", wantError: "parsing path element #0 at char #0 ' ': unexpected identifier char"},
		{pathString: "A ", wantError: "parsing path element #0 at char #1 ' ': unexpected identifier char"},
		{pathString: "[0]A", wantError: "parsing path element #1 at char #3 'A': unexpected char"},
		{pathString: "A.", wantError: "parsing path element #1: incomplete property name"},
		{pathString: "A.[0]", wantError: "parsing path element #1 at char #2 '[': unexpected char"},
		{pathString: "[''A]", wantError: "parsing path element #0 at char #3 'A': unexpected char"},
		{pathString: "[''[]", wantError: "parsing path element #0 at char #3 '[': unexpected char"},
		{pathString: "", want: nil},
		{
			pathString: "[0]",
			want: []validation.PropertyPathElement{
				validation.ArrayIndex(0),
			},
		},
		{
			pathString: "[' ']",
			want: []validation.PropertyPathElement{
				validation.PropertyName(" "),
			},
		},
		{
			pathString: "[0][1][2][3][4][5][6][7][8][9][10]",
			want: []validation.PropertyPathElement{
				validation.ArrayIndex(0),
				validation.ArrayIndex(1),
				validation.ArrayIndex(2),
				validation.ArrayIndex(3),
				validation.ArrayIndex(4),
				validation.ArrayIndex(5),
				validation.ArrayIndex(6),
				validation.ArrayIndex(7),
				validation.ArrayIndex(8),
				validation.ArrayIndex(9),
				validation.ArrayIndex(10),
			},
		},
		{
			pathString: "array[1].property",
			want: []validation.PropertyPathElement{
				validation.PropertyName("array"),
				validation.ArrayIndex(1),
				validation.PropertyName("property"),
			},
		},
		{
			pathString: "foo1.bar",
			want: []validation.PropertyPathElement{
				validation.PropertyName("foo1"),
				validation.PropertyName("bar"),
			},
		},
		{
			pathString: "$foo.bar",
			want: []validation.PropertyPathElement{
				validation.PropertyName("$foo"),
				validation.PropertyName("bar"),
			},
		},
		{
			pathString: "_foo.bar",
			want: []validation.PropertyPathElement{
				validation.PropertyName("_foo"),
				validation.PropertyName("bar"),
			},
		},
		{
			pathString: "['@foo'].bar",
			want: []validation.PropertyPathElement{
				validation.PropertyName("@foo"),
				validation.PropertyName("bar"),
			},
		},
		{
			pathString: "['@foo'][0]",
			want: []validation.PropertyPathElement{
				validation.PropertyName("@foo"),
				validation.ArrayIndex(0),
			},
		},
		{
			pathString: "['foo.bar'].baz",
			want: []validation.PropertyPathElement{
				validation.PropertyName("foo.bar"),
				validation.PropertyName("baz"),
			},
		},
		{
			pathString: `['foo.\'bar\''].baz`,
			want: []validation.PropertyPathElement{
				validation.PropertyName("foo.'bar'"),
				validation.PropertyName("baz"),
			},
		},
		{
			pathString: `['0'].baz`,
			want: []validation.PropertyPathElement{
				validation.PropertyName(`0`),
				validation.PropertyName("baz"),
			},
		},
		{
			pathString: `['foo[0]'].baz`,
			want: []validation.PropertyPathElement{
				validation.PropertyName(`foo[0]`),
				validation.PropertyName("baz"),
			},
		},
		{
			pathString: `[''].baz`,
			want: []validation.PropertyPathElement{
				validation.PropertyName(``),
				validation.PropertyName("baz"),
			},
		},
		{
			pathString: `foo['']`,
			want: []validation.PropertyPathElement{
				validation.PropertyName("foo"),
				validation.PropertyName(""),
			},
		},
		{
			pathString: `['\''].baz`,
			want: []validation.PropertyPathElement{
				validation.PropertyName(`'`),
				validation.PropertyName("baz"),
			},
		},
		{
			pathString: `['\\'].baz`,
			want: []validation.PropertyPathElement{
				validation.PropertyName(`\`),
				validation.PropertyName("baz"),
			},
		},
		{
			pathString: `['\\\'foo'].baz`,
			want: []validation.PropertyPathElement{
				validation.PropertyName(`\'foo`),
				validation.PropertyName("baz"),
			},
		},
		{
			pathString: `фу.baz`,
			want: []validation.PropertyPathElement{
				validation.PropertyName(`фу`),
				validation.PropertyName("baz"),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.pathString, func(t *testing.T) {
			var got validation.PropertyPath
			err := got.UnmarshalText([]byte(test.pathString))

			if test.wantError == "" {
				assert.NoError(t, err)
				assert.Equal(t, test.want, got.Elements())
			} else {
				assert.Nil(t, got.Elements())
				assert.EqualError(t, err, test.wantError)
			}
		})
	}
}

func TestPropertyPath_UnmarshalJSON(t *testing.T) {
	jsonData := `"array[1].property"`

	var path validation.PropertyPath
	err := json.Unmarshal([]byte(jsonData), &path)

	require.NoError(t, err)
	assert.Equal(t,
		[]validation.PropertyPathElement{
			validation.PropertyName("array"),
			validation.ArrayIndex(1),
			validation.PropertyName("property"),
		},
		path.Elements(),
	)
}

func BenchmarkPropertyPath_String(b *testing.B) {
	// cpu: Intel(R) Core(TM) i9-9900K CPU @ 3.60GHz
	// BenchmarkPropertyPath_String
	// BenchmarkPropertyPath_String-16    	  225926	      7175 ns/op	    4352 B/op	       5 allocs/op
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

//nolint:stylecheck
func FuzzPropertyPath_UnmarshalText(f *testing.F) {
	f.Add("array[1].property")
	f.Add(" ")
	f.Add(fmt.Sprintf("[%d]", math.MaxInt))
	f.Add(fmt.Sprintf("[%d]", uint(math.MaxInt)+1))
	f.Add("A0000[0000].A000000") // valid path, but not symmetric
	f.Add("['A']")               // valid path, but not symmetric

	assertion := NewPropertyPathAssertion()

	f.Fuzz(func(t *testing.T, sourcePath string) {
		var path validation.PropertyPath
		err := path.UnmarshalText([]byte(sourcePath))
		if err != nil {
			return
		}
		assertion.Equal(t, sourcePath, path)
	})
}

type PropertyPathAssertion struct {
	zerosTrimmer *regexp.Regexp
}

func NewPropertyPathAssertion() *PropertyPathAssertion {
	return &PropertyPathAssertion{
		zerosTrimmer: regexp.MustCompile(`\[\d+]`),
	}
}

func (a *PropertyPathAssertion) Equal(t *testing.T, sourcePath string, parsedPath validation.PropertyPath) {
	t.Helper()

	encodedPathBytes, err := parsedPath.MarshalText()
	require.NoError(t, err)

	encodedPath := string(encodedPathBytes)
	bracketedPath := a.bracketPath(parsedPath)
	trimmedPath := a.trimIndexZeros(sourcePath)
	reencodedPath := a.reencodePath(trimmedPath)

	if trimmedPath != encodedPath && trimmedPath != bracketedPath && reencodedPath != encodedPath {
		assert.Fail(t, fmt.Sprintf(`paths not equal: source "%s", encoded "%s"`, sourcePath, encodedPath))
	}
}

func (a *PropertyPathAssertion) reencodePath(trimmedPath string) string {
	s := strings.Builder{}
	for _, c := range trimmedPath {
		s.WriteRune(c)
	}
	return s.String()
}

func (a *PropertyPathAssertion) trimIndexZeros(path string) string {
	return a.zerosTrimmer.ReplaceAllStringFunc(path, func(s string) string {
		i, err := strconv.Atoi(strings.Trim(s, "[]"))
		if err != nil {
			return s
		}
		return "[" + strconv.Itoa(i) + "]"
	})
}

func (a *PropertyPathAssertion) bracketPath(path validation.PropertyPath) string {
	s := strings.Builder{}

	for _, e := range path.Elements() {
		s.WriteString("[")
		if e.IsIndex() {
			s.WriteString(e.String())
		} else {
			s.WriteString("'")
			p := e.String()
			for _, c := range p {
				if c == '\'' || c == '\\' {
					s.WriteRune('\\')
				}
				s.WriteRune(c)
			}
			s.WriteString("'")
		}
		s.WriteString("]")
	}

	return s.String()
}
