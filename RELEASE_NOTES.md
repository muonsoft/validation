# Release Notes

## v0.19.0

### Added

- **Slice validation**: `Slice`, `SliceProperty`, `Each`, and `EachProperty` for validating slices with per-element constraints.
- **HasUniqueValuesBy**: Constraint to ensure slice elements are unique by a key function.
- **Validate method on constraints**: Constraints can now be used directly with `Each` and `This`.

### Changed

- **Validator initialization**: Global validator replaced with atomic pointer for thread-safe usage. Use new validator setup methods in tests.
- **CheckNoViolations**: Now accepts variadic `errors` for improved error handling and termination conditions.
- **Documentation**: README restructured with installation, custom constraints, and property paths; expanded custom constraints guide with interface details and examples.
- **Skill docs**: SKILLS.md removed; new reference and skill documentation for adding validation constraints.

### Fixed

- Correct handling of single violations returned from validatable objects in `validateIt`.

### Breaking Changes

- Tests using `CheckNoViolations` or validator setup need to be updated to the new signatures and helpers.
