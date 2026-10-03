# Validate the normalized state on update

Implement Service.Update. Load the record by ID first, set Name to strings.TrimSpace(name), and use the existing Record.Validate. On local violations do not run NameTaken or Save. Otherwise check NameTaken(name, excludeID) using the normalized name and loaded identity; report ErrDuplicate at name if taken. Save the exact validated record only on success. Preserve context and technical errors. A preliminary Load is required, not forbidden I/O.

Modify only `solution.go`. Preserve the exported contract and existing behavior unless the task explicitly changes it. You may add private helpers there. Run your own checks if useful; do not add dependencies or change go.mod.
