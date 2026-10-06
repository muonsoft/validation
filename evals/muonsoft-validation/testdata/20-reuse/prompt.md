# Extend and simplify account validation

Extend Profile.Validate and Preferences.Validate. Profile requires DisplayName
(max 20 characters), Email with standard library email-constraint semantics, and
Age in [18,120]. Preferences may be nil; when present, Theme must be light or dark
(empty invalid), Quota may be nil but if supplied is in [0,100], and Tags have at
most 3 entries, each required and max 8 characters; duplicate tags produce one
ErrNotUnique at tags. Use ErrNotInRange for age/quota bounds, and the corresponding built-in codes for all other rules. Profile must validate its
existing Address under address, supplying the provided region constraint unchanged.
Preserve caller prefix and context. The existing Address and project rule remain
the owners of their rules: reuse them, do not copy/reimplement them. Collect all
independent violations. Prefer the existing built-in constraints over handwritten
length, email, range, enum and uniqueness validators. This reuse requirement is
reviewed separately from behavioral correctness. Do not mutate input.
Preserve protected files and exported contracts. Edit the supplied implementation
files; you may add private helpers in new top-level Go files (not *_test.go).
Do not change dependencies. Use existing library/project rules when their semantics
match. Hidden tests assess behavior; later code review separately assesses reuse,
rule ownership, and unnecessary abstractions. You may run temporary tests, but
remove them before submission. Paths are PropertyPath.String() JavaScript-style
property paths, not JSON Pointer. Preserve any caller prefix.
