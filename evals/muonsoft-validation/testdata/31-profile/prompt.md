# Validate a partial profile update

Implement ProfilePatch.Validate for a settings endpoint. Each field distinguishes
omission, clearing, and a replacement: Present=false ignores Value even if it is
non-nil; Present=true with nil clears and is allowed; a supplied string cannot be
empty. Empty strings yield the standard blank-value error at that field.
Country replacements must first be recognized two-letter countries/territories
according to the library's installed data, then belong to the business list
DE, US, GB. Locale replacements must first be valid locale identifiers, then be
one of en-US or de-DE. Lists compare the supplied spelling exactly, without
normalization. Valid but unsupported values yield ErrCountryPolicy or
ErrLocalePolicy using that error's message. Invalid formats use the library's
standard format error instead; never also add a policy error for the same input.
Both fields must be checked independently, so a failure in country cannot hide
a locale failure. Paths are country and locale under the caller's prefix.

Preserve contract.go, existing protected files, and dependencies. Edit validation.go.
Private helpers in new top-level non-test Go files are allowed. Remove temporary
tests before submission. Use existing library/project rules when their semantics
match. Do not change input, normalize values, or create a replacement validator.
