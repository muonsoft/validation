# Extend optional scheduling settings

Implement Document and Schedule validation while reusing Limits.Validate.
Document Title is always required. Schedule may be nil. When present: Note may
be nil, but provided empty is ErrIsBlank and length > 12 is ErrTooLong. Mode is
quiet or active; empty mode is allowed in draft but invalid in publish. Publish
is selected together with the default group. Nonempty invalid modes are always
ErrNoSuchChoice. Lower and Upper independently may be nil, but present numbers
must be nonnegative. If both exist and Lower >= Upper, also return ErrRange at
upper, even if a number is negative. All independent errors accumulate, including
multiple errors on upper. Limits is checked through its existing richer method
under schedule.limits, using the supplied maximum. Preserve input, context,
language, groups and caller prefix. An absent Schedule must not invoke Limits. A present Schedule always has a non-nil Limits dependency.
Preserve protected files and exported contracts. Edit the supplied implementation
files; you may add private helpers in new top-level Go files (not *_test.go).
Do not change dependencies. Use existing library/project rules when their semantics
match. Hidden tests assess behavior; later code review separately assesses reuse,
rule ownership, and unnecessary abstractions. You may run temporary tests, but
remove them before submission. Paths are PropertyPath.String() JavaScript-style
property paths, not JSON Pointer. Preserve any caller prefix.
