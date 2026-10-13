# Implement a complete service request

Implement Request, Contact, Visit and Task validation. A request has a required
Title (max 40 characters), required contact Name and Email (email format), and
1–2 Visits. Each visit needs Address and 1–2 Tasks. Each task needs ServiceCode,
Variant in {basic, extended} (empty is invalid), and positive Duration.
Use the corresponding built-in errors. Task pairs (ServiceCode, Variant) must be
unique within a visit: ErrNotUnique on every duplicate task element, including
the first. The same pair in different visits is allowed. Collect all independent
errors, including collection bounds and child errors together. Child validation
must also work standalone. Do not mutate the request. JSON tags describe the
external names; the Go field names are not the error paths.
Preserve protected files and exported contracts. Edit the supplied implementation
files; you may add private helpers in new top-level Go files (not *_test.go).
Do not change dependencies. Use existing library/project rules when their semantics
match. Hidden tests assess behavior; later code review separately assesses reuse,
rule ownership, and unnecessary abstractions. You may run temporary tests, but
remove them before submission. Paths are PropertyPath.String() JavaScript-style
property paths, not JSON Pointer. Preserve any caller prefix.
