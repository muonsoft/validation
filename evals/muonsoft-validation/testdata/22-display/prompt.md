# Group existing errors for display

Implement GroupErrors for the existing application. Group already produced
Issue values by exact Path. Preserve first-seen path order, and message order
within each group, including duplicate messages and the empty root path. Return
nil for nil/empty input. Result slices must be independent of the input and of
other invocations. Preserve ExistingValidate and all unrelated behavior. This is
a presentation change, with no new validation rules or dependencies.
Edit group.go only; private helpers in new top-level Go files are allowed. Preserve all protected contracts and dependencies.
