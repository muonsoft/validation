# Validate an explicit-presence patch
Implement Patch.Validate(ctx, validator). Mode.Present=false means no operation:
ignore Value even if nonnil and invalid. Present=true with Value=nil means clear;
it is accepted except when group activate is selected, when it yields ErrClear
at mode. Present=true with a value accepts only 1 or 2 in every group. All other
integers, including zero, yield exactly ErrMode at mode (not ErrNoSuchChoice or
ErrIsBlank). Note is optional; nil is accepted but a supplied empty string yields
ErrIsBlank, and more than 5 Unicode characters yields ErrTooLong, both at note.
Rules apply even when only activate is selected, without DefaultGroup. Accumulate
independent note and mode errors. Preserve caller paths, translations and error
factory. Do not mutate the patch or pointed-to values.

Preserve contract.go and dependencies. Edit validation.go; private helpers in new
top-level non-test Go files are allowed. Remove temporary tests before submission.
Use existing library/project rules when their semantics match.
