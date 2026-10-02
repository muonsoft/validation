# Country, language, and locale validation

Use these checks for country fields, language metadata, and locale preferences.
They share the existing `golang.org/x/text/language` dependency and add no modules,
network requests, or separately maintained code lists. Each check has a different
contract:

| Constraint | Boolean helper | Standalone validator | Purpose |
|------------|----------------|----------------------|---------|
| `it.IsCountry()` | `is.Country` | `validate.Country` | Two-letter country or territory code according to x/text |
| `it.IsLanguage()` | `is.Language` | `validate.Language` | Standalone two- or three-letter language code |
| `it.IsLocale()` | `is.Locale` | `validate.Locale` | Locale identifier, optionally including script, region, variants, or extensions |

```go
err := validator.Validate(ctx,
    validation.StringProperty("country", "de", it.IsCountry()),
    validation.StringProperty("language", "eng", it.IsLanguage()),
    validation.StringProperty("locale", "zh-Hant-TW", it.IsLocale()),
)
```

All three checks ignore letter case, accept empty strings, reject surrounding
whitespace, and leave input unchanged. Constraints also skip nil values; add
`it.IsNotBlank()` for required fields. Standalone validators return `nil` or their
`validate.ErrInvalidCountry`, `validate.ErrInvalidLanguage`, or
`validate.ErrInvalidLocale` sentinel. Constraints produce the corresponding
`validation.ErrInvalid…` violation with English/Russian messages, customizable
messages/errors, conditions, and validation groups.

## Country: input codes, not service availability

Country fields in addresses and APIs often need more than a two-letter format
check. `Country` requires exactly two bytes, successful `language.ParseRegion`,
`Region.IsCountry()`, and no `Region.IsGroup()`. It accepts countries and
territories as classified by the installed x/text data. It rejects three-letter
codes, numeric regions, and geographic groups.

This is **not a strict list of currently assigned ISO 3166-1 codes**. x/text
includes CLDR additions, aliases, historical codes, and special codes in its
classification. For example, `DE`, `AX`, `XK`, `UK`, `SU`, `AN`, and even the
special code `UN` pass in the current dependency version; `EU`, `ZZ`, `AA`,
`DEU`, and `276` fail. Passing this check is not evidence that a country currently
exists or that your business delivers or provides service there.

## Language: standalone language metadata

`Language` uses `language.ParseBase`. It accepts two- and three-letter codes
recognized by x/text, such as `en`, `eng`, and `deu`, including legacy codes such
as `iw`. Composite identifiers such as `en-US` and `zh-Hant` fail: use `Locale`
when regional or script distinctions matter.

Special codes are intentionally accepted: `und` (undetermined), `mul` (multiple
languages), `zxx` (no linguistic content), and private-use codes `qaa` through
`qtz`. Unknown codes such as `zzz` fail. This check does not restrict input to
living languages or languages for which your application has translations.

## Locale: parser validity, not a catalogue of supported locales

`Locale` succeeds when `language.Parse` returns no error. It accepts BCP 47 and
Unicode locale forms supported by that parser, including `en`, `en-US`,
`en_US`, `zh-Hant-TW`, numeric regions such as `es-419`, variants such as
`de-CH-1901`, and extensions such as `en-u-ca-gregory`. A country is not required.

The parser also accepts `und`, `root`, legacy identifiers such as `i-klingon`,
private-use tags such as `x-private`, and private-use subtags. Numeric regions and
extension values follow the parser's permissive rules: `en-999` and
`en-u-ca-foobar` pass, without asserting that the region or calendar is meaningful.
This is not exhaustive registry validation of every component or combination.

Malformed input (`en--US`), unknown variants (`en-foobar`), POSIX locale strings
(`en_US.UTF-8`), and HTTP language preference lists (`en,fr;q=0.8`) fail. A parser
error always fails validation, even when the parser returns a usable partial tag.

## Application policy and data updates

These checks establish the contracts above; they do not prove translation or
formatting-resource availability, validate relationships between separate country
and language fields, or restrict input to your application's supported values.
Use an explicit allowlist, for example `it.IsOneOf("en-US", "de-DE")`, for that
policy. `IsOneOf` compares the actual string: normalize input yourself before
applying a canonical allowlist if you want case-insensitive or alias matching.
The validators do not perform that normalization.

Accepted values follow the version of x/text pinned in `go.mod`; updating it can
change recognition of codes and tags. Exact Symfony compatibility is not a goal.
No separate ICU installation or additional dependency is required.

See the upstream contracts for [ParseRegion](https://pkg.go.dev/golang.org/x/text/language#ParseRegion),
[IsCountry](https://pkg.go.dev/golang.org/x/text/language#Region.IsCountry),
[ParseBase](https://pkg.go.dev/golang.org/x/text/language#ParseBase), and
[Parse](https://pkg.go.dev/golang.org/x/text/language#Parse).
