# Draft and publish rules

Implement Document.Validate. Code is required in the default group. Title is required only in the publish group. The application selects default alone for drafts and default plus publish for publishing. Use ErrIsBlank. Support a caller-supplied property prefix and preserve its groups.

Modify only `solution.go`. Preserve the exported contract and existing behavior unless the task explicitly changes it. You may add private helpers there. Run your own checks if useful; do not add dependencies or change go.mod.
