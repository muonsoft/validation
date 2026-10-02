# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased](https://github.com/muonsoft/validation/compare/v0.19.0...HEAD)

### Changed

- CI and local lint use golangci-lint `v2.13.2` (`gomodguard_v2`); `goconst` disabled in `.golangci.yml`.

### Added

- Add `it.IsCardScheme`, `validate.CardScheme`, and `is.CardScheme` with explicit selection of 12 payment card schemes, Symfony 8.0 prefix/length rules, and English/Russian messages. Checksum validation remains separate via `it.IsLUHN`; no new dependencies.

- Add ISO week validation with `it.IsWeek()`, `validate.Week`, and `is.Week`, inclusive min/max bounds, calendar-aware week 53 checks, and English/Russian messages. No new dependencies; document empty-value and literal early-year behavior.

- Add `it.HasMinWordCount`, `it.HasMaxWordCount`, and `it.HasWordCountBetween` with Unicode-aware counting, customizable errors/messages, and English/Russian plural forms. No new dependencies; document differences from ICU locale-aware segmentation.

- Add country, language, and locale validation with `it.IsCountry()`, `it.IsLanguage()`, `it.IsLocale()`, and matching `is`/`validate` helpers. Use existing `golang.org/x/text/language` data without new dependencies; include English/Russian messages and documentation of accepted formats, special codes, and application-support limitations.

- Add UTF-8 validation with `it.IsUTF8()`, `is.UTF8`, and `validate.UTF8`, English/Russian messages, and documentation of input-boundary use cases and encoding-only limitations. No new dependencies.

- ISO 4217 currency code validation: `it.IsCurrency()`, `validate.Currency`, `is.Currency`, with `validation.ErrInvalidCurrency` / `message.InvalidCurrency` and English and Russian translations (behavior aligned with Symfony `Currency`; recognized codes from `golang.org/x/text/currency.ParseISO`).
- ISBN validation: `it.IsISBN()` with `Only10` / `Only13`, `validate.ISBN` with `validate.ISBNOnly10` / `validate.ISBNOnly13`, `is.ISBN`; `validation.ErrInvalidISBN`, `ErrInvalidISBN10`, `ErrInvalidISBN13` / `message.InvalidISBN`, `InvalidISBN10`, `InvalidISBN13` and English and Russian translations (behavior aligned with Symfony `Isbn`).
- MAC address validation: `it.IsMacAddress()` with `WithType` (Symfony `MacAddress` type names: `validate.MacAddressTypeAll`, `MacAddressTypeBroadcast`, etc.), `validate.MacAddress` with `validate.WithMacAddressType`, `is.MACAddress`; `validation.ErrInvalidMAC` / `message.InvalidMAC` and English and Russian translations. Only 48-bit (6-octet) addresses accepted via [net.ParseMAC] (colon, hyphen, dot forms); EUI-64 and longer forms are rejected.
- ISSN (International Standard Serial Number) validation: `it.IsISSN()`, `validate.ISSN`, `is.ISSN`, with `validation.ErrInvalidISSN` / `message.InvalidISSN` and English and Russian translations (ISO 3297 mod 11 check digit; optional hyphen; behavior aligned with Symfony `Issn`).
- BIC / SWIFT validation: `it.IsBIC()` with `CaseInsensitive` and `WithIBAN` (and `WithIBANError` / `WithIBANMessage`), `validate.BIC` with `validate.BICCaseInsensitive` and `validate.BICWithIBAN`, `is.BIC`; `validation.ErrInvalidBIC`, `validation.ErrBICIBANCountryMismatch`, `message.InvalidBIC`, `message.BICNotAssociatedWithIBAN` and English and Russian translations (behavior aligned with Symfony `Bic` / `BicValidator`; country/territory check via `golang.org/x/text/language` regions plus Symfony’s BIC-to-IBAN territory map).
- IBAN validation: `it.IsIBAN()`, `validate.IBAN`, `is.IBAN`, with `validation.ErrInvalidIBAN` / `message.InvalidIBAN` and English and Russian translations (behavior aligned with Symfony `Iban`; country patterns from Symfony 7.2 `IbanValidator`).
- **NoSuspiciousCharacters** (spoofing / homoglyph checks): `it.HasNoSuspiciousCharacters()` with `CheckInvisible` / `CheckMixedNumbers` / `CheckHiddenOverlay`, `WithoutInvisible` / `WithoutMixedNumbers` / `WithoutHiddenOverlay` (bitmask-based), locale and single-script restrictions; `validate.NoSuspiciousCharacters`; errors `validation.ErrSuspiciousInvisible`, `ErrSuspiciousMixedNumbers`, `ErrSuspiciousHiddenOverlay`, `ErrSuspiciousCharactersRestriction` and English/Russian messages (behavior inspired by Symfony `NoSuspiciousCharacters`, implemented without CGO; may differ from ICU in edge cases).
- CIDR notation validation: `it.IsCIDR()` with `IPv4Only`, `IPv6Only`, `WithVersion`, `WithNetmaskRange`, and separate invalid vs netmask-range messages; `validate.CIDR`, `validate.CIDRViolationNetmaskBounds`, `is.CIDR`; `validation.ErrInvalidCIDR` / `validation.ErrCIDRNetmaskOutOfRange` and English and Russian translations (behavior aligned with Symfony `Cidr`). Exported names use the **CIDR** initialism per [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments).
- LUHN (mod 10 / Luhn) checksum validation: `it.IsLUHN()`, `validate.LUHN`, `is.LUHN`, with `validation.ErrInvalidLUHN` / `message.InvalidLUHN` and English and Russian translations (behavior aligned with Symfony `Luhn`).
- ISIN (International Securities Identification Number) validation: `it.IsISIN()`, `validate.ISIN`, `is.ISIN`, with `validation.ErrInvalidISIN` / `message.InvalidISIN` and English and Russian translations (behavior aligned with Symfony `Isin`).
- **HasUniqueValuesBy**: `SkipEmptyKeys()` on `it.UniqueByConstraint` skips elements whose key equals the zero value for `K`, so they are not counted toward uniqueness (e.g. optional IDs).
- IANA timezone validation: `it.IsTimezone()` with `WithZone` (`validate.TimezoneZoneAfrica`, `TimezoneZoneEurope`, etc.), `validate.Timezone` with `validate.WithTimezoneZone`, `is.Timezone`; `validation.ErrInvalidTimezone` / `message.InvalidTimezone` and English and Russian translations. Uses a bundled IANA tzdata 2026c identifier list; accepts `UTC` and names containing `/`, including legacy aliases.

### Fixed

- Reject malformed timezone identifiers and system-specific paths in `validate.Timezone`, `is.Timezone`, and `it.IsTimezone`, independently of system timezone files and `ZONEINFO`.
- Keep `it.TimezoneConstraint.WithZone` copies independent when deriving multiple constraints from a shared template.

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
