# Update a resource bundle with batch validation

Implement Bundle.Validate, Service.Update and CheckReferences. Update first loads
the bundle by command ID. Copy loaded state before changes: repository-owned data
must not be mutated, including on errors. Replace Name and the entire References
slice from the command; trim Name and every reference ID before validation. Ignore
PreviousName for validation: it represents persisted history. The repository method
is NameTaken(ctx, name, excludedID): pass the normalized new name first and the
loaded bundle ID second, i.e. NameTaken(ctx, normalizedName, loaded.ID).
Bundle requires Name (max 30 chars), 1–4 references; every
reference needs ID (ErrIsBlank at references[i].id) and positive Weight.
Accumulate all local errors before NameTaken or Find. After local success call
NameTaken then CheckReferences, accumulating ErrDuplicate at name and every
ErrMissing at references[i].id. A technical error stops later work, preserves
errors.Is identity, and is never converted to a violation. Save once only on full
success, with exactly the normalized validated bundle; loaded object remains
unchanged even after success. Forward the original context everywhere.
CheckReferences also works standalone with an already collection-scoped validator:
it adds only each original index and id. One Find for a nonempty list, none for
empty; it may deduplicate lookup IDs but must report every missing occurrence in
original order/positions. The map uses true for found; false or absent is missing.
The service owns the references prefix. Do not persist on any error.
Preserve protected files and exported contracts. Edit the supplied implementation
files; you may add private helpers in new top-level Go files (not *_test.go).
Do not change dependencies. Use existing library/project rules when their semantics
match. Hidden tests assess behavior; later code review separately assesses reuse,
rule ownership, and unnecessary abstractions. You may run temporary tests, but
remove them before submission. Paths are PropertyPath.String() JavaScript-style
property paths, not JSON Pointer. Preserve any caller prefix.
