# Stop on technical errors without losing identity

Repair ValidateAndSave. Run first, second, third in order, accumulating validation violations. A technical error must stop later checks and preserve errors.Is identity, even after an earlier violation. Call save once only when all checks succeed, with the original context. The existing Filter call computes every check too early.

Modify only `solution.go`. Preserve the exported contract and existing behavior unless the task explicitly changes it. You may add private helpers there. Run your own checks if useful; do not add dependencies or change go.mod.
