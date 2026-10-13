# Choose built-in constraints

Read this before implementing format checks, counting, uniqueness, or parsing in
an application rule. Check the installed module's version and symbol-specific
`go doc`; the table is a selection guide, not a reason to upgrade dependencies.
Prefer a built-in when its input, accepted values, nil/empty behavior, error
identity, and options match the application's contract. Compose requiredness,
format, and business policy separately. Preserve project rules that add semantics.

## Available in v0.19.0

| Requirement | Start with | Decision to preserve |
| --- | --- | --- |
| Required string or pointer | `it.IsNotBlank`, `it.IsNotNil` | Blank means zero value, not trimmed whitespace; `WithAllowedNil` supports optional nonblank strings. |
| Enum or allowlist | `it.IsOneOf` | Add `WithoutBlank` to reject zero values; allowlists express application policy. |
| String length or numeric bounds | `it.HasMinLength`, `it.HasMaxLength`, `it.IsGreaterThanOrEqual`, `it.IsLessThanOrEqual` | String length counts characters; use `len` and a numeric constraint for bytes. |
| Email, URL, hostname, IP, UUID | `it.IsEmail`, `it.IsURL`, `it.IsHostname`, `it.IsIP`, `it.IsUUID` | Inspect options before writing a regex; requiredness remains separate. |
| Collection size and duplicates | `it.HasCountBetween`, `it.HasUniqueValues`, `it.HasUniqueValuesBy` | Validate each element as well as collection rules; derive uniqueness keys from the actual contract. |

## Added in v0.20.0

| Requirement | Start with | Decision to preserve |
| --- | --- | --- |
| File name, suffix, MIME metadata, detected content | `it.IsFileName`, `it.HasFileExtension`, `it.IsMIMEType`, `it.HasContentType` | These check different properties. Content takes `[]byte` through `This`/`Each`. See [file uploads](file-uploads.md). |
| Password score | `it.HasPasswordStrength`; `validate.EstimatePasswordStrength` for a score | Default minimum is Medium (2); empty strings are evaluated, nil skipped. See [passwords](passwords.md). |
| CSS color | `it.IsCSSColor` | Select allowed `validate.CSSColorFormat` values when the contract restricts syntax. |
| Payment card scheme and checksum | `it.IsCardScheme`, `it.IsLUHN` | Select schemes explicitly using `validate.CardScheme…`; scheme matching does not check the checksum. Neither proves a card exists. |
| IBAN, BIC/SWIFT, ISIN | `it.IsIBAN`, `it.IsBIC`, `it.IsISIN` | Inspect BIC case/country options. Format/checksum validity is distinct from account or instrument existence. |
| ISBN, ISSN, currency code | `it.IsISBN`, `it.IsISSN`, `it.IsCurrency` | ISBN can restrict versions; currency recognition is distinct from currencies supported by the application. |
| Country, language, locale | `it.IsCountry`, `it.IsLanguage`, `it.IsLocale` | Country follows x/text classification including aliases/special codes. Language accepts standalone codes; locale accepts composite tags. Add a business allowlist separately. |
| ISO week with inclusive bounds | `it.IsWeek().WithMin(...).WithMax(...)` | Validates calendar week 53; distinguishes format, nonexistent week, and out-of-range errors. |
| Word count | `it.HasMinWordCount`, `it.HasMaxWordCount`, `it.HasWordCountBetween` | Unicode letters/numbers form words; internal apostrophes stay within a word, hyphens split. This is not `strings.Fields` or locale-aware segmentation. |
| CIDR, MAC address, IANA timezone | `it.IsCIDR`, `it.IsMacAddress`, `it.IsTimezone` | Inspect IP/netmask, address-type, or zone options. Exact casing matters: the boolean MAC helper is `is.MACAddress`. |
| UTF-8 and suspicious characters | `it.IsUTF8`, `it.HasNoSuspiciousCharacters` | Encoding validity and suspicious-character policy are separate; inspect the latter's configurable checks. |
| Uniqueness excluding unset keys | `it.HasUniqueValuesBy(...).SkipEmptyKeys()` | Skips keys equal to their type's zero value, not arbitrary “empty” application values. Validate required keys separately if needed. |

Use `validation` arguments plus `it` for object validation with paths, groups,
translations, and accumulated violations. Use `is` for a boolean predicate and
`validate` for a standalone error check; not every constraint has all three APIs.
For an uncertain helper, inspect its actual name/signature rather than deriving it
from the constructor name.

Before replacing a manual rule, compare invalid and boundary cases as well as a
valid value. A built-in's broader accepted set must not silently broaden business
policy. Conversely, do not replace a format check with a narrow list of examples.
Nil/empty rules and configuration errors vary: most format constraints skip empty
values, but password strength evaluates them, and invalid options can fail even
when input is empty. Configuration failures remain ordinary errors that stop
validation, rather than user-data violations.
