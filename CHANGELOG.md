# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Luhn (mod 10) checksum validation: `it.IsLuhn()`, `validate.Luhn`, `is.Luhn`, with `validation.ErrInvalidLuhn` / `message.InvalidLuhn` and English and Russian translations (behavior aligned with Symfony `Luhn`).
- ISIN (International Securities Identification Number) validation: `it.IsISIN()`, `validate.ISIN`, `is.ISIN`, with `validation.ErrInvalidISIN` / `message.InvalidISIN` and English and Russian translations (behavior aligned with Symfony `Isin`).

## [0.19.0] - 2026-02-09

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
[0.19.0]: https://github.com/muonsoft/validation/releases/tag/v0.19.0
