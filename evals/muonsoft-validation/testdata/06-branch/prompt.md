# Repair prematurely executed validation

Repair Validate so a blank Reference produces ErrIsBlank at reference and does not invoke the dependency. Name is independently required and must still be checked. A valid Reference runs the dependency even if Name is blank, with the reference path prefix. Preserve all independent violations. The existing code accidentally invokes the dependency before the sequence can guard it.

Modify only `solution.go`. Preserve the exported contract and existing behavior unless the task explicitly changes it. You may add private helpers there. Run your own checks if useful; do not add dependencies or change go.mod.
