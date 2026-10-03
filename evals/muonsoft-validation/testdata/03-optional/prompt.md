# Optional fields and dependent bounds

Implement Input.Validate. Note may be nil, but an explicitly supplied empty string is invalid (ErrIsBlank at note). Status must be draft or ready; reject empty/unknown status (ErrNoSuchChoice at status). Lower and Upper may independently be nil. When both exist, require Lower < Upper (ErrRange at upper). Accumulate independent errors, do not panic on nil, and do not modify input.

Modify only `solution.go`. Preserve the exported contract and existing behavior unless the task explicitly changes it. You may add private helpers there. Run your own checks if useful; do not add dependencies or change go.mod.
