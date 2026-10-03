# Adapt an existing richer method

Implement ValidateDetail by composing Legacy.Validate under the relative property detail. Do not change the legacy signature or create a replacement validator. Preserve the outer prefix, selected groups, language, and caller context. Limit must be forwarded unchanged. The legacy rule runs only for group strict.

Modify only `solution.go`. Preserve the exported contract and existing behavior unless the task explicitly changes it. You may add private helpers there. Run your own checks if useful; do not add dependencies or change go.mod.
