// CSS color patterns adapted from Symfony Validator.
// Copyright (c) Fabien Potencier <fabien@symfony.com>.
// Distributed under the MIT license; see https://github.com/symfony/validator/blob/8.0/LICENSE.

package validate

import (
	"errors"
	"regexp"
)

// CSSColorFormat selects one of the CSS color syntaxes supported by Symfony 8.0.
type CSSColorFormat string

const (
	// CSSColorHexLong selects the hex_long format.
	CSSColorHexLong CSSColorFormat = "hex_long"
	// CSSColorHexLongWithAlpha selects the hex_long_with_alpha format.
	CSSColorHexLongWithAlpha CSSColorFormat = "hex_long_with_alpha"
	// CSSColorHexShort selects the hex_short format.
	CSSColorHexShort CSSColorFormat = "hex_short"
	// CSSColorHexShortWithAlpha selects the hex_short_with_alpha format.
	CSSColorHexShortWithAlpha CSSColorFormat = "hex_short_with_alpha"
	// CSSColorBasicNamedColors selects the basic_named_colors format.
	CSSColorBasicNamedColors CSSColorFormat = "basic_named_colors"
	// CSSColorExtendedNamedColors selects the extended_named_colors format.
	CSSColorExtendedNamedColors CSSColorFormat = "extended_named_colors"
	// CSSColorSystemColors selects the system_colors format.
	CSSColorSystemColors CSSColorFormat = "system_colors"
	// CSSColorKeywords selects the keywords format.
	CSSColorKeywords CSSColorFormat = "keywords"
	// CSSColorRGB selects the rgb format.
	CSSColorRGB CSSColorFormat = "rgb"
	// CSSColorRGBA selects the rgba format.
	CSSColorRGBA CSSColorFormat = "rgba"
	// CSSColorHSL selects the hsl format.
	CSSColorHSL CSSColorFormat = "hsl"
	// CSSColorHSLA selects the hsla format.
	CSSColorHSLA CSSColorFormat = "hsla"
)

var (
	// ErrInvalidCSSColor indicates an unsupported or malformed color.
	ErrInvalidCSSColor = errors.New("invalid CSS color")
	// ErrInvalidCSSColorFormats indicates an unknown format in the selection.
	ErrInvalidCSSColorFormats = errors.New("invalid CSS color formats")
)

// Patterns follow Symfony 8.0 CssColorValidator (MIT licensed):
// https://github.com/symfony/validator/blob/8.0/Constraints/CssColorValidator.php
// PCRE's ASCII whitespace includes vertical tab, unlike Go's regexp \s.
var cssColorPatterns = map[CSSColorFormat]*regexp.Regexp{
	CSSColorHexLong:             regexp.MustCompile(`(?i)^#[0-9a-f]{6}$`),
	CSSColorHexLongWithAlpha:    regexp.MustCompile(`(?i)^#[0-9a-f]{8}$`),
	CSSColorHexShort:            regexp.MustCompile(`(?i)^#[0-9a-f]{3}$`),
	CSSColorHexShortWithAlpha:   regexp.MustCompile(`(?i)^#[0-9a-f]{4}$`),
	CSSColorBasicNamedColors:    regexp.MustCompile(`(?i)^(black|silver|gray|white|maroon|red|purple|fuchsia|green|lime|olive|yellow|navy|blue|teal|aqua)$`),
	CSSColorExtendedNamedColors: regexp.MustCompile(`(?i)^(aliceblue|antiquewhite|aqua|aquamarine|azure|beige|bisque|black|blanchedalmond|blue|blueviolet|brown|burlywood|cadetblue|chartreuse|chocolate|coral|cornflowerblue|cornsilk|crimson|cyan|darkblue|darkcyan|darkgoldenrod|darkgray|darkgreen|darkgrey|darkkhaki|darkmagenta|darkolivegreen|darkorange|darkorchid|darkred|darksalmon|darkseagreen|darkslateblue|darkslategray|darkslategrey|darkturquoise|darkviolet|deeppink|deepskyblue|dimgray|dimgrey|dodgerblue|firebrick|floralwhite|forestgreen|fuchsia|gainsboro|ghostwhite|gold|goldenrod|gray|green|greenyellow|grey|honeydew|hotpink|indianred|indigo|ivory|khaki|lavender|lavenderblush|lawngreen|lemonchiffon|lightblue|lightcoral|lightcyan|lightgoldenrodyellow|lightgray|lightgreen|lightgrey|lightpink|lightsalmon|lightseagreen|lightskyblue|lightslategray|lightslategrey|lightsteelblue|lightyellow|lime|limegreen|linen|magenta|maroon|mediumaquamarine|mediumblue|mediumorchid|mediumpurple|mediumseagreen|mediumslateblue|mediumspringgreen|mediumturquoise|mediumvioletred|midnightblue|mintcream|mistyrose|moccasin|navajowhite|navy|oldlace|olive|olivedrab|orange|orangered|orchid|palegoldenrod|palegreen|paleturquoise|palevioletred|papayawhip|peachpuff|peru|pink|plum|powderblue|purple|red|rosybrown|royalblue|saddlebrown|salmon|sandybrown|seagreen|seashell|sienna|silver|skyblue|slateblue|slategray|slategrey|snow|springgreen|steelblue|tan|teal|thistle|tomato|turquoise|violet|wheat|white|whitesmoke|yellow|yellowgreen)$`),
	CSSColorSystemColors:        regexp.MustCompile(`(?i)^(Canvas|CanvasText|LinkText|VisitedText|ActiveText|ButtonFace|ButtonText|ButtonBorder|Field|FieldText|Highlight|HighlightText|SelectedItem|SelectedItemText|Mark|MarkText|GrayText)$`),
	CSSColorKeywords:            regexp.MustCompile(`(?i)^(transparent|currentColor)$`),
	CSSColorRGB:                 regexp.MustCompile(`(?i)^rgb\([\t\n\f\r\x0b ]*(0|255|25[0-4]|2[0-4]\d|1\d\d|0?\d?\d),[\t\n\f\r\x0b ]*(0|255|25[0-4]|2[0-4]\d|1\d\d|0?\d?\d),[\t\n\f\r\x0b ]*(0|255|25[0-4]|2[0-4]\d|1\d\d|0?\d?\d)[\t\n\f\r\x0b ]*\)$`),
	CSSColorRGBA:                regexp.MustCompile(`(?i)^rgba\([\t\n\f\r\x0b ]*(0|255|25[0-4]|2[0-4]\d|1\d\d|0?\d?\d),[\t\n\f\r\x0b ]*(0|255|25[0-4]|2[0-4]\d|1\d\d|0?\d?\d),[\t\n\f\r\x0b ]*(0|255|25[0-4]|2[0-4]\d|1\d\d|0?\d?\d),[\t\n\f\r\x0b ]*(0|0?\.\d+|1(\.0)?)[\t\n\f\r\x0b ]*\)$`),
	CSSColorHSL:                 regexp.MustCompile(`(?i)^hsl\([\t\n\f\r\x0b ]*(0|360|35\d|3[0-4]\d|[12]\d\d|0?\d?\d),[\t\n\f\r\x0b ]*(0|100|\d{1,2})%,[\t\n\f\r\x0b ]*(0|100|\d{1,2})%[\t\n\f\r\x0b ]*\)$`),
	CSSColorHSLA:                regexp.MustCompile(`(?i)^hsla\([\t\n\f\r\x0b ]*(0|360|35\d|3[0-4]\d|[12]\d\d|0?\d?\d),[\t\n\f\r\x0b ]*(0|100|\d{1,2})%,[\t\n\f\r\x0b ]*(0|100|\d{1,2})%,[\t\n\f\r\x0b ]*(0|0?\.\d+|1(\.0)?)[\t\n\f\r\x0b ]*\)$`),
}

// CSSColor validates a color using Symfony 8.0's supported syntaxes.
// No formats selects all formats; multiple formats use OR semantics.
// Names and function names are ASCII case-insensitive. Input is not trimmed.
// RGB/HSL functions use the legacy comma-separated syntax and strict ranges;
// modern CSS color functions, percentages in RGB, and CSS-wide keywords are
// unsupported. Empty values are valid; use [NotBlank] to require a value.
// Unknown formats return [ErrInvalidCSSColorFormats], even for empty input.
// Unsupported values return [ErrInvalidCSSColor].
func CSSColor(value string, formats ...CSSColorFormat) error {
	for _, format := range formats {
		if _, ok := cssColorPatterns[format]; !ok {
			return ErrInvalidCSSColorFormats
		}
	}
	if value == "" {
		return nil
	}
	for _, char := range value {
		if char > 127 {
			return ErrInvalidCSSColor
		}
	}
	if len(formats) == 0 {
		return matchCSSColor(value)
	}
	for _, format := range formats {
		if cssColorPatterns[format].MatchString(value) {
			return nil
		}
	}
	return ErrInvalidCSSColor
}

func matchCSSColor(value string) error {
	for _, pattern := range cssColorPatterns {
		if pattern.MatchString(value) {
			return nil
		}
	}
	return ErrInvalidCSSColor
}
