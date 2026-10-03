# Check members using an already scoped validator

Implement CheckMembers for a locally validated list. The caller supplies a validator already scoped to the members collection; do not append another collection property. Query Find once for a nonempty list, never for empty. Its map reports membership by ID. Report ErrMissing for each missing original occurrence, preserving original indices, caller context and path. Deduplication for the lookup is allowed. Preserve technical error identity and return technical failures without turning them into violations.
