# Fix a display formatter without changing validation

Complete DisplayName: trim each name part, join nonempty parts with one space, return "Anonymous" when both are empty. This project uses its existing ValidateName function; preserve it and do not introduce any validation framework. Only this formatting behavior is requested.

Modify only `solution.go`. Preserve the exported contract and existing behavior unless the task explicitly changes it. You may add private helpers there. Run your own checks if useful; do not add dependencies or change go.mod.
