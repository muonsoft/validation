# Validate a password update

Implement PasswordUpdate.Validate. A nil Password means no update and succeeds
without calling Estimate. A supplied password, including an empty string, must
score at least Strong (3) using the library's strength estimator. If Estimate is
provided, it replaces that estimator, receives the unchanged password exactly
once, and controls the score even for empty input. Scores 0–2 yield exactly one
ErrPasswordTooWeak at password. Scores 3–4 pass; scores outside 0–4 are technical
configuration errors, never user-data violations. Preserve validator context.
Do not trim or normalize input or put the password in errors, message parameters,
or logs. Nil Estimate selects the library's default heuristic, not an error.

Preserve contract.go and dependencies. Edit validation.go; private helpers in new
top-level non-test Go files are allowed. Remove temporary tests before submission.
Use existing library/project rules when their semantics match. Do not mutate input.
Preserve the supplied validator and its path prefix.
