# Custom reusable constraint and translation

Implement PrefixConstraint.ValidateString and NewValidator. Skip nil and empty strings; requiredness is composed separately. Accept strings starting with Prefix, otherwise create ErrPrefix using its template and parameter {{ prefix }}. Register a Russian translation: Значение должно начинаться с {{ prefix }}. Default fallback stays English. Use the supplied validator so paths and language survive; return constructor errors.

Modify only `solution.go`. Preserve the exported contract and existing behavior unless the task explicitly changes it. You may add private helpers there. Run your own checks if useful; do not add dependencies or change go.mod.
