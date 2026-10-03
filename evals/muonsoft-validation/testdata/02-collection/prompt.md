# Extend nested collection validation

Implement validation for Input and Entry. Title and every entry Code are required; Quantity must be positive. Entries must contain 1–3 elements and Codes must be unique. Preserve all independent violations, including title, entries, and entries[i].quantity. Use built-in codes. Attribute duplicate codes to entries. Paths must contain actual array indices, not bracketed property names.

Modify only `solution.go`. Preserve the exported contract and existing behavior unless the task explicitly changes it. You may add private helpers there. Run your own checks if useful; do not add dependencies or change go.mod.
