# Preserve a tenant-specific email policy

Implement Signup.Validate while retaining the injected EmailPolicy from policy.go.
Name and Email are required. A supplied Email must satisfy the library's ordinary
email syntax policy before the tenant rule runs. Use standard required/email-format
errors at name and email. Do not replace tenant policy with a hard-coded domain
list or infer its behavior from its name: it can change between requests.
For a syntactically valid nonempty email, call the supplied policy exactly once,
even when Name is invalid. For a missing or malformed email, do not call it.
The policy receives the original context and a validator already positioned at
email under the caller's prefix. Preserve the caller's groups, translations and
violation factory. Policy violations accumulate with independent Name failures;
ordinary policy errors propagate as technical failures, preserving errors.Is.
The caller always selects the default group, possibly alongside tenant-specific
groups. No input normalization or network work belongs in this method.

Preserve contract.go, existing protected files, and dependencies. Edit validation.go.
Private helpers in new top-level non-test Go files are allowed. Remove temporary
tests before submission. Use existing library/project rules when their semantics
match. Do not change input, normalize values, or create a replacement validator.
