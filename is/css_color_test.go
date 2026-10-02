package is_test

import (
	"fmt"
	"testing"

	"github.com/muonsoft/validation/is"
	"github.com/muonsoft/validation/validate"
	"github.com/stretchr/testify/require"
)

func TestCSSColor(t *testing.T) {
	require.True(t, is.CSSColor(""))
	require.True(t, is.CSSColor("red"))
	require.True(t, is.CSSColor("#abc", validate.CSSColorHexShort, validate.CSSColorRGB))
	require.False(t, is.CSSColor("red", validate.CSSColorRGB))
	require.False(t, is.CSSColor("bad"))
	require.False(t, is.CSSColor("", "unknown"))
}

func ExampleCSSColor() {
	fmt.Println(is.CSSColor("rgba(0,0,0,.5)"))
	// Output:
	// true
}
