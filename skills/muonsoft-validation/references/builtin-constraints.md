# Built-in constraint inventory

This inventory covers every exported constructor in package `it` in the current
library, grouped by task, with corresponding standalone helpers. Names in each
package column belong to that package; a dash means no matching helper exists.
Generic type arguments and option methods are omitted from the main tables.
Use symbol-specific `go doc` in the application to confirm signatures and options.

Prefer a built-in when its input, accepted values, nil/empty behavior, error
identity, and options match the application's contract. Compose requiredness,
format, and business policy separately; retain project rules that add semantics.

## Presence, booleans, and choices

| Check | Input | `it` constructors | `validate` helpers | `is` helpers |
| --- | --- | --- | --- | --- |
| Required and blank values | `string, bool, number, comparable, count, time` | `IsNotBlank`, `IsBlank`, `IsNotNil`, `IsNil` | — | — |
| Boolean values | `bool` | `IsTrue`, `IsFalse` | — | — |
| Equality and choice | `comparable` | `IsEqualTo`, `IsNotEqualTo`, `IsOneOf` | — | `InList` |

`IsNotBlank` rejects zero/empty values, including false; it does not trim strings.
`IsNotNil` allows a non-nil pointer to an empty value. `WithAllowedNil` makes
`IsNotBlank` useful for an optional field that must be nonblank when supplied.
`IsOneOf` accepts zero values by default; use `WithoutBlank` for strict membership.
Business-specific choices can use `IsOneOf(...).WithError(err).WithMessage(err.Message())`;
a custom error identity does not require a custom membership predicate. See the
[patch example](optional-and-conditional.md#patch-fields-with-format-and-business-choices)
for sequencing format and policy.

Generic presence constructors are also available explicitly:

| Rule | Numeric values | Comparable values |
| --- | --- | --- |
| Not blank | `it.IsNotBlankNumber[T]` | `it.IsNotBlankComparable[T]` |
| Blank | `it.IsBlankNumber[T]` | `it.IsBlankComparable[T]` |
| Not nil | `it.IsNotNilNumber[T]` | `it.IsNotNilComparable[T]` |
| Nil | `it.IsNilNumber[T]` | `it.IsNilComparable[T]` |

## Numbers and time values

| Check | Input | `it` constructors | `validate` helpers | `is` helpers |
| --- | --- | --- | --- | --- |
| Numeric comparisons | `number` | `IsLessThan`, `IsLessThanOrEqual`, `IsGreaterThan`, `IsGreaterThanOrEqual`, `IsBetween` | — | — |
| Signs and divisibility | `number` | `IsPositive`, `IsPositiveOrZero`, `IsNegative`, `IsNegativeOrZero`, `IsDivisibleBy`, `IsDivisibleByFloat` | — | `DivisibleBy` |
| Time comparisons | `time.Time` | `IsEarlierThan`, `IsEarlierThanOrEqual`, `IsLaterThan`, `IsLaterThanOrEqual`, `IsBetweenTime` | — | — |

Range bounds are inclusive. Choose the integer or floating-point divisibility
constructor for the input type; a zero divisor is a configuration error.
Time comparisons operate on `time.Time`; date/time string parsing is listed below.
For a business interval with one sentinel for either bound, use
`it.IsBetween[int64](min, max).WithError(err).WithMessage(err.Message())`.
This retains the built-in range rule instead of duplicating it in a predicate or
custom numeric constraint. Use separate comparisons when the contract requires
different errors for the two bounds.

## Collections

| Check | Input | `it` constructors | `validate` helpers | `is` helpers |
| --- | --- | --- | --- | --- |
| Collection size | `int count` | `HasMinCount`, `HasMaxCount`, `HasCountBetween`, `HasExactCount`, `HasCountDivisibleBy` | — | — |
| Uniqueness | `slice` | `HasUniqueValues`, `HasUniqueValuesBy` | — | `Unique` |

Use `validation.Countable(len(values), ...)` for size, `Comparables` for comparable
value uniqueness, and `Slice` with `HasUniqueValuesBy` for uniqueness by a key.
`HasUniqueValuesBy(...).SkipEmptyKeys()` skips keys equal to their type's zero
value. Required keys and per-element rules remain separate checks.

## Text, patterns, and encoding

| Check | Input | `it` constructors | `validate` helpers | `is` helpers |
| --- | --- | --- | --- | --- |
| String length | `string` | `HasMinLength`, `HasMaxLength`, `HasLengthBetween`, `HasExactLength` | — | — |
| Patterns and JSON | `string` | `Matches`, `DoesNotMatch`, `IsJSON` | — | `JSON` |
| Numeric text | `string` | `IsInteger`, `IsNumeric` | — | `Integer`, `Number` |
| Word count | `string` | `HasMinWordCount`, `HasMaxWordCount`, `HasWordCountBetween` | — | — |
| UTF-8 | `string` | `IsUTF8` | `UTF8` | `UTF8` |
| Suspicious characters | `string` | `HasNoSuspiciousCharacters` | `NoSuspiciousCharacters` | — |

Length counts Unicode code points (runes), not bytes or grapheme clusters; use
`len` with a numeric constraint for bytes. Length constraints skip empty strings.
Word counting uses Unicode letters/numbers and combining marks; internal
apostrophes stay within a word, hyphens split. It is not `strings.Fields` or
locale-aware segmentation. `Matches`/`DoesNotMatch` take a compiled
`*regexp.Regexp`. JSON syntax does not establish a payload schema.
`IsUTF8` checks encoding; `HasNoSuspiciousCharacters` has separate configurable
invisible-character, mixed-number, overlay, locale, and script checks.

## Email, networks, and date/time text

| Check | Input | `it` constructors | `validate` helpers | `is` helpers |
| --- | --- | --- | --- | --- |
| Email | `string` | `IsEmail`, `IsHTML5Email` | — | `Email`, `HTML5Email` |
| Hostnames | `string` | `IsHostname`, `IsLooseHostname` | — | `StrictHostname`, `Hostname` |
| URLs and IP addresses | `string` | `IsURL`, `IsIP`, `IsIPv4`, `IsIPv6`, `IsCIDR` | `URL`, `IP`, `IPv4`, `IPv6`, `CIDR` | `URL`, `IP`, `IPv4`, `IPv6`, `CIDR` |
| MAC addresses | `string` | `IsMacAddress` | `MacAddress` | `MACAddress` |
| Date, time, and timezone text | `string` | `IsDate`, `IsTime`, `IsDateTime`, `IsTimezone` | `Timezone` | `Timezone` |
| ISO weeks | `string` | `IsWeek` | `Week` | `Week` |

Choose email and hostname strictness deliberately; note the hostname helper
mapping (`IsHostname` → `is.StrictHostname`, `IsLooseHostname` → `is.Hostname`).
Inspect URL schemes/relative-URL options and IP/CIDR version/netmask options before
writing a parser. MAC API casing differs across packages, as shown in the table.
`IsDateTime` defaults to RFC3339, `IsDate` to `2006-01-02`, and `IsTime` to
`15:04:05`; `WithLayout` changes the Go time layout. `IsTimezone().WithZone(...)`
restricts recognized IANA names by geographical region. `IsWeek().WithMin(...).
WithMax(...)` uses inclusive bounds and calendar-aware week 53 checks, with
distinct format, week-number, and range errors.

## Codes, identifiers, and barcodes

| Check | Input | `it` constructors | `validate` helpers | `is` helpers |
| --- | --- | --- | --- | --- |
| Country, language, locale | `string` | `IsCountry`, `IsLanguage`, `IsLocale` | `Country`, `Language`, `Locale` | `Country`, `Language`, `Locale` |
| Currency codes | `string` | `IsCurrency` | `Currency` | `Currency` |
| Identifiers | `string` | `IsUUID`, `IsULID` | `UUID`, `ULID` | `UUID`, `ULID` |
| Bank and security identifiers | `string` | `IsIBAN`, `IsBIC`, `IsISIN` | `IBAN`, `BIC`, `ISIN` | `IBAN`, `BIC`, `ISIN` |
| Publication identifiers | `string` | `IsISBN`, `IsISSN` | `ISBN`, `ISSN` | `ISBN`, `ISSN` |
| Barcodes | `string` | `IsEAN8`, `IsEAN13`, `IsUPCA`, `IsUPCE` | `EAN8`, `EAN13`, `UPCA`, `UPCE` | `EAN8`, `EAN13`, `UPCA`, `UPCE` |
| Payment card schemes and checksums | `string` | `IsCardScheme`, `IsLUHN` | `CardScheme`, `LUHN` | `CardScheme`, `LUHN` |

Country/language/locale recognition uses the installed x/text data, including
aliases and special codes. Country is not a shipping/service allowlist; language
is a standalone code, whereas locale may include script, region, and extensions.
Use a separate application allowlist when required; these checks do not normalize
input. Inspect UUID version/format, ISBN version, and BIC case/country options.
`IsCardScheme` requires explicit `validate.CardScheme…` choices and checks
prefix/length; compose `IsLUHN` for a checksum. Format/checksum validity does not
establish the existence of a card, account, or financial instrument.

## Colors, passwords, and files

| Check | Input | `it` constructors | `validate` helpers | `is` helpers |
| --- | --- | --- | --- | --- |
| CSS colors | `string` | `IsCSSColor` | `CSSColor` | `CSSColor` |
| Password strength | `string` | `HasPasswordStrength` | `PasswordStrength`, `EstimatePasswordStrength` | — |
| File names and metadata | `string` | `IsFileName`, `HasFileExtension`, `IsMIMEType` | `FileName`, `FileExtension`, `MIMEType` | — |
| File content | `[]byte` | `HasContentType` | `ContentType` | — |

`IsCSSColor` can restrict accepted syntax with `validate.CSSColorFormat` values.
Password strength defaults to Medium (2); it scores empty strings and skips nil.
See [passwords](passwords.md) for estimator options and secret handling.
Filename, extension, declared MIME, and detected content are separate properties;
`HasContentType` takes bytes through `This`/`Each`. See [file uploads](file-uploads.md)
for bounded reading, Windows restrictions, detectors, and matching suffix/content.

## Apply the selected rule

Use `validation.String`, `Number`, `Bool`, or `Time` with typed constraints, and
`NilString`, `NilNumber`, `NilBool`, or `NilTime` for pointers. `This` and `Each`
accept constraints implementing `Constraint[T]` for their input type. Use `it`
with these arguments to retain paths, groups, translations, and accumulated
violations. Use `is` for a boolean predicate and `validate` for a standalone error.
Helper option signatures may differ from their `it` counterparts.

Use `When`/`WhenGroups` for conditional rules and customize violations with the
constraint's error/message methods. Multi-error constraints may have separate
min/max or format-specific setters. See [validation flow](validation-flow.md)
for sequencing and alternatives, [nested collections](nested-collections.md) for
elements, and [custom constraints](custom-constraints.md) for rules not covered here.

Before replacing a manual rule, compare invalid and boundary cases as well as a
valid value. Do not broaden application policy just because a format accepts
more values, or replace format validation with a narrow list of examples.
Most format constraints skip nil/empty values, but requiredness, counting, and
password rules differ. Invalid options can fail even when input is empty;
configuration failures stop validation as ordinary errors rather than being
collected as user-data violations.
