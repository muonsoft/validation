# Repair recursive form validation

Implement Form and Field validation for a finite tree (no cycles). Root Fields
and every group Children contain 1–3 elements. Field Name is required. Kind must
be text or group (empty/unknown -> ErrNoSuchChoice). Text fields require Value;
groups validate Children, ignoring Value. Names are unique among siblings:
ErrNotUnique on each duplicate sibling's name, including the first. Repetition
in another branch is allowed. Text fields ignore Children. A root field has depth
1. A field deeper than maxDepth gets ErrDepth on that field; do not descend into
that field, but continue siblings. maxDepth is always positive. Every field's
Attributes map has required string values; keys are arbitrary literal property
names, including dots, slashes, tildes, quotes, backslashes, empty and numeric
strings. Accumulate independent errors, preserve original indices and prefixes.
Paths follow the library's JavaScript notation; do not implement JSON Pointer.
Preserve protected files and exported contracts. Edit the supplied implementation
files; you may add private helpers in new top-level Go files (not *_test.go).
Do not change dependencies. Use existing library/project rules when their semantics
match. Hidden tests assess behavior; later code review separately assesses reuse,
rule ownership, and unnecessary abstractions. You may run temporary tests, but
remove them before submission. Paths are PropertyPath.String() JavaScript-style
property paths, not JSON Pointer. Preserve any caller prefix.
