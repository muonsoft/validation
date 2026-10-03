# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.20.0] - 2026-10-03

### Added

- Add file upload validation with `it.IsFileName`, optional Windows restrictions, `it.HasFileExtension`, `it.IsMIMEType`, and `it.HasContentType`, matching `validate` helpers, and replaceable content detection. See [file uploads](docs/file-uploads.md).
- Add `it.HasPasswordStrength`, `validate.PasswordStrength`, and `validate.EstimatePasswordStrength` with configurable minimum scores and estimators. The default uses Symfony 8.0’s byte-based heuristic; empty passwords are evaluated, and passwords are not automatically included in violations. See [password strength](docs/password-strength.md).
- Add `it.IsCSSColor`, `validate.CSSColor`, and `is.CSSColor` with 12 selectable color formats. See [CSS colors](docs/css-colors.md) for supported syntax and limitations.
- Add `it.IsCardScheme`, `validate.CardScheme`, and `is.CardScheme` with explicit selection of 12 payment card schemes using Symfony 8.0 prefix/length rules. Checksum validation remains separate via `it.IsLUHN`. See [card schemes](docs/card-schemes.md).
- Add `it.IsWeek`, `validate.Week`, and `is.Week` with inclusive bounds and calendar-aware week 53 checks. See [ISO weeks](docs/iso-weeks.md).
- Add `it.HasMinWordCount`, `it.HasMaxWordCount`, and `it.HasWordCountBetween` with Unicode-aware counting and English/Russian plural forms. See [word count](docs/word-count.md) for segmentation limitations.
- Add country, language, and locale validation with `it.IsCountry`, `it.IsLanguage`, `it.IsLocale`, and matching `is`/`validate` helpers using existing `golang.org/x/text/language` data. See [international validation](docs/international.md) for accepted formats and special codes.
- Add UTF-8 validation with `it.IsUTF8`, `is.UTF8`, and `validate.UTF8`. See [UTF-8 validation](docs/utf8.md) for input-boundary use cases and encoding-only limitations.
- Add ISO 4217 currency code validation with `it.IsCurrency`, `validate.Currency`, and `is.Currency`, using codes recognized by `golang.org/x/text/currency.ParseISO`.
- Add ISBN validation with `it.IsISBN`, `validate.ISBN`, and `is.ISBN`, including options to accept only ISBN-10 or ISBN-13 and independent options when deriving constraints.
- Add MAC address validation with `it.IsMacAddress`, `validate.MacAddress`, and `is.MACAddress`, including address-type filtering. Keep options independent when deriving constraints. Accept only 48-bit addresses in colon, hyphen, or dot notation; reject EUI-64 and longer forms.
- Add ISSN validation with `it.IsISSN`, `validate.ISSN`, and `is.ISSN`, including the ISO 3297 mod 11 check digit and optional hyphen.
- Add BIC/SWIFT validation with `it.IsBIC`, `validate.BIC`, and `is.BIC`, including optional case-insensitive matching and IBAN country matching.
- Add IBAN validation with `it.IsIBAN`, `validate.IBAN`, and `is.IBAN`, using country patterns from Symfony 7.2.
- Add `it.HasNoSuspiciousCharacters` and `validate.NoSuspiciousCharacters` with configurable invisible-character, mixed-number, hidden-overlay, locale, and single-script checks. The implementation uses no CGO and may differ from ICU in edge cases.
- Add CIDR validation with `it.IsCIDR`, `validate.CIDR`, and `is.CIDR`, including IP-version selection, netmask bounds, and separate invalid-input and netmask-range messages. Keep options independent when deriving constraints; cap IPv6 prefixes and reported bounds at 128 bits.
- Add Luhn checksum validation with `it.IsLUHN`, `validate.LUHN`, and `is.LUHN`.
- Add ISIN validation with `it.IsISIN`, `validate.ISIN`, and `is.ISIN`.
- Add `SkipEmptyKeys` to `it.UniqueByConstraint` to exclude elements whose key is the zero value for its type from uniqueness checks.
- Add IANA timezone validation with `it.IsTimezone`, `validate.Timezone`, and `is.Timezone`, including geographical region filtering. Use a bundled tzdata 2026c identifier list, independent of system files and `ZONEINFO`; accept `UTC` and known names containing `/`, including legacy aliases, and reject malformed identifiers and system-specific paths. Constraints derived with `WithZone` keep their options independent.
- Include English and Russian messages for all new constraints without adding dependencies.

### Changed

- Reorganize the documentation around a shorter quick start and constraint catalog; correct API examples, validation groups, translations, and release instructions.
- Show the Scrutinizer test coverage badge in the README next to the code quality badge.

### Fixed

- Prevent goroutine leaks in `validation.Async` after fatal errors.
- Keep `ViolationList.Join` inputs independent and handle repeated or self joins without cycles; preserve every branch of joined errors without swallowing fatal errors.
- Keep derived argument paths and URL, IP, and UUID constraint options independent.
- Copy translated message parameters before rendering to preserve reusable rules across languages and concurrent validations.
- Correct floating-point divisibility for negative operands and report zero divisors as constraint configuration errors.
- Enforce hostname length including separators and reject reserved top-level domains regardless of case.
- Round-trip property names containing Unicode digits and clear existing property paths when unmarshaling empty text.

## [0.19.0](https://github.com/muonsoft/validation/releases/tag/v0.19.0) - 2026-02-09

### Added

- **Slice validation**: `Slice`, `SliceProperty`, `Each`, and `EachProperty` for validating slices with per-element constraints.
- **HasUniqueValuesBy**: Constraint to ensure slice elements are unique by a key function.
- **Validate method on constraints**: Constraints can now be used directly with `Each` and `This`.

### Changed

- **Validator initialization**: Global validator replaced with atomic pointer for thread-safe usage. Use new validator setup methods in tests.
- **CheckNoViolations**: Now accepts variadic `errors` for improved error handling and termination conditions.
- **Documentation**: README restructured with installation, custom constraints, and property paths; expanded custom constraints guide with interface details and examples.
- **Skill docs**: SKILLS.md removed; new reference and skill documentation for adding validation constraints.

### Breaking

- Tests using `CheckNoViolations` or validator setup need to be updated to the new signatures and helpers.

### Fixed

- Correct handling of single violations returned from validatable objects in `validateIt`.

[Unreleased]: https://github.com/muonsoft/validation/compare/v0.19.0...HEAD
