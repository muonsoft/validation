# CSS colors

`it.IsCSSColor()` accepts all color formats supported by
[Symfony 8.0 CssColor](https://github.com/symfony/validator/blob/8.0/Constraints/CssColorValidator.php).
Pass formats explicitly to restrict the allowed syntax:

```go
err := validator.Validate(context.Background(),
    validation.String(color, it.IsCSSColor(validate.CSSColorRGB, validate.CSSColorRGBA)),
)
```

The `validate.CSSColorFormat` constants are `CSSColorHexLong`,
`CSSColorHexLongWithAlpha`, `CSSColorHexShort`, `CSSColorHexShortWithAlpha`,
`CSSColorBasicNamedColors`, `CSSColorExtendedNamedColors`, `CSSColorSystemColors`,
`CSSColorKeywords`, `CSSColorRGB`, `CSSColorRGBA`, `CSSColorHSL`, and `CSSColorHSLA`.
Hex formats accept 6, 8, 3, and 4 digits respectively. Keywords are `transparent`
and `currentColor`. Names and function names are ASCII case-insensitive.
Standalone helpers `validate.CSSColor` and `is.CSSColor` accept the same optional
format list; multiple formats use OR semantics, and an empty list allows all.

This is the pinned Symfony syntax, not a complete modern CSS parser. RGB channels
are integers from 0 to 255; HSL hue is an integer from 0 to 360, with integer
percentages from 0 to 100. Functions require commas. Alpha accepts `0`, `1`, `1.0`,
or a decimal fraction such as `.5` or `0.25`; `1.00` is rejected. Whitespace is
allowed after `(`, after commas, and before `)`, but not before commas or around
the entire input. RGB percentages, space/slash syntax, `rebeccapurple`, CSS-wide
keywords, `var()`, `hwb()`, `lab()`, `oklch()`, and `color()` are unsupported.

Nil and empty values are skipped; use `it.IsNotBlank()` to require a value.
Unknown formats produce `validate.ErrInvalidCSSColorFormats` (a constraint
configuration error through `it`), even with empty input. Unlike Symfony's
mixed format lists, every selected format must be recognized. Disabled constraints
and unmatched groups skip validation entirely. Invalid colors produce
`validation.ErrInvalidCSSColor` or `validate.ErrInvalidCSSColor`.
`WithError`, `WithMessage` (with `{{ value }}` and custom parameters), `When`,
`WhenGroups`, `This`, and `Each` are supported, with English/Russian messages.
No dependencies are added.

See [Usage](usage.md) for imports and validator setup, and the
[constraint catalog](constraints.md) for related checks.
