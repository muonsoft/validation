package validate_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/muonsoft/validation/validate"
	"github.com/stretchr/testify/require"
)

func TestCSSColorFormats(t *testing.T) {
	tests := []struct {
		format validate.CSSColorFormat
		values []string
	}{
		{validate.CSSColorHexLong, []string{"#aAbBcC", "#000000", "#ffffff"}},
		{validate.CSSColorHexLongWithAlpha, []string{"#12345678", "#FFFFFFFF"}},
		{validate.CSSColorHexShort, []string{"#abc", "#000"}},
		{validate.CSSColorHexShortWithAlpha, []string{"#abcd", "#FFFF"}},
		{validate.CSSColorBasicNamedColors, []string{"red", "BLACK", "aqua"}},
		{validate.CSSColorExtendedNamedColors, []string{"aliceblue", "DarkSlateGrey", "yellowgreen", "red"}},
		{validate.CSSColorSystemColors, []string{"Canvas", "selecteditemtext", "BUTTONBORDER", "MarkText"}},
		{validate.CSSColorKeywords, []string{"transparent", "currentColor"}},
		{validate.CSSColorRGB, []string{"rgb(0,0,0)", "RGB(255, 254, 199)", "rgb(099,00,1)", "rgb(\t1,\n2,\r3\f\x0b )"}},
		{validate.CSSColorRGBA, []string{"rgba(0,255,1,0)", "rgba(1,2,3,.5)", "rgba(1,2,3,0.000)", "RGBA(1,2,3,1.0)"}},
		{validate.CSSColorHSL, []string{"hsl(0,0%,0%)", "HSL(360,100%,100%)", "hsl(099,09%,00%)"}},
		{validate.CSSColorHSLA, []string{"hsla(359,99%,50%,.9)", "hsla(0,0%,0%,1)"}},
	}
	for _, tt := range tests {
		for _, value := range tt.values {
			t.Run(string(tt.format)+"/"+value, func(t *testing.T) {
				require.NoError(t, validate.CSSColor(value, tt.format))
				require.NoError(t, validate.CSSColor(value))
				require.ErrorIs(t, validate.CSSColor(" "+value, tt.format), validate.ErrInvalidCSSColor)
				require.ErrorIs(t, validate.CSSColor(value+"\n", tt.format), validate.ErrInvalidCSSColor)
			})
		}
	}
	// Each format's representative must be rejected by unrelated formats.
	for _, tt := range tests {
		for _, other := range tests {
			if tt.format == other.format {
				continue
			}
			if tt.format == validate.CSSColorBasicNamedColors && other.format == validate.CSSColorExtendedNamedColors {
				continue
			}
			require.ErrorIs(t, validate.CSSColor(tt.values[0], other.format), validate.ErrInvalidCSSColor, "%s as %s", tt.format, other.format)
		}
	}
}

func TestCSSColorInvalid(t *testing.T) {
	values := []string{
		" ", "#", "#12", "#12345", "#1234567", "#123456789", "#ggg", "123456", "#１２３",
		"rebeccapurple", "inherit", "initial", "unset", "revert", "AccentColor", "blacK", "tranſparent", "unknown",
		"rgb(256,0,0)", "rgb(-1,0,0)", "rgb(+1,0,0)", "rgb(1.0,0,0)", "rgb(1e2,0,0)",
		"rgb(100%,0%,0%)", "rgb(1 2 3)", "rgb(1,2,3,.5)", "rgb (1,2,3)", "rgb(1 ,2,3)",
		"rgb(0000,0,0)", "rgb(١,2,3)", "rgb(\u00a01,2,3)", "rgb(1,2)", "rgb(1,2,3,)",
		"rgba(1,2,3)", "rgba(1,2,3,-.1)", "rgba(1,2,3,1.1)", "rgba(1,2,3,1.00)",
		"rgba(1,2,3,0.)", "rgba(1,2,3,50%)", "rgba(1 2 3 / .5)", "rgba(1,2,3,01)",
		"hsl(361,50%,50%)", "hsl(-1,50%,50%)", "hsl(10deg,50%,50%)", "hsl(1.5,50%,50%)",
		"hsl(0,101%,50%)", "hsl(0,50%,101%)", "hsl(0,0,0)", "hsl(0,1.5%,50%)", "hsl(0 50% 50%)",
		"hsla(0,50%,50%,2)", "hwb(0 0% 0%)", "lab(50% 0 0)", "oklch(50% .2 30)", "color(srgb 1 0 0)",
		"var(--color)", "rgb(calc(1),2,3)", "red;", "red\x00", "\xff",
	}
	for _, value := range values {
		t.Run(value, func(t *testing.T) { require.ErrorIs(t, validate.CSSColor(value), validate.ErrInvalidCSSColor) })
	}
}

func TestCSSColorBoundaries(t *testing.T) {
	for channel := 0; channel < 3; channel++ {
		for _, value := range []int{0, 1, 9, 10, 99, 100, 199, 200, 249, 250, 254, 255, 256} {
			components := []string{"0", "0", "0"}
			components[channel] = fmt.Sprint(value)
			want := validate.ErrInvalidCSSColor
			if value <= 255 {
				want = nil
			}
			require.ErrorIs(t, validate.CSSColor("rgb("+strings.Join(components, ",")+")"), want)
		}
	}
	for _, hue := range []int{0, 99, 100, 199, 200, 299, 300, 349, 350, 359, 360, 361} {
		want := validate.ErrInvalidCSSColor
		if hue <= 360 {
			want = nil
		}
		require.ErrorIs(t, validate.CSSColor(fmt.Sprintf("hsl(%d,100%%,0%%)", hue)), want)
	}
}

func TestCSSColorSelection(t *testing.T) {
	require.NoError(t, validate.CSSColor(""))
	require.NoError(t, validate.CSSColor("", validate.CSSColorRGB))
	require.NoError(t, validate.CSSColor("red", validate.CSSColorRGB, validate.CSSColorBasicNamedColors))
	require.NoError(t, validate.CSSColor("red", validate.CSSColorBasicNamedColors, validate.CSSColorBasicNamedColors))
	for _, value := range []string{"", "red"} {
		require.ErrorIs(t, validate.CSSColor(value, "unknown"), validate.ErrInvalidCSSColorFormats)
		require.ErrorIs(t, validate.CSSColor(value, validate.CSSColorBasicNamedColors, "unknown"), validate.ErrInvalidCSSColorFormats)
	}
}

func ExampleCSSColor() {
	fmt.Println(validate.CSSColor("#AbCd"))
	fmt.Println(validate.CSSColor("red", validate.CSSColorRGB))
	// Output:
	// <nil>
	// invalid CSS color
}
