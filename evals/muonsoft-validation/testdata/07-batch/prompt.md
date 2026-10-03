# Validate references in one batch

Implement CheckReferences for an already locally validated list. Make one Find call for a nonempty list, none for empty. A returned true means the ID exists. Report every missing original occurrence with ErrMissing at references[i], even if lookup IDs are deduplicated. Preserve the caller prefix and context. Return technical errors unchanged or wrapped with %w.

Modify only `solution.go`. Preserve the exported contract and existing behavior unless the task explicitly changes it. You may add private helpers there. Run your own checks if useful; do not add dependencies or change go.mod.
