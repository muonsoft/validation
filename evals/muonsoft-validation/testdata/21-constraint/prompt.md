# Create a reusable category-aware code constraint

Implement CodeConstraint, NewCodeConstraint, ValidatorOptions, Signup.Validate
and Transfer.Validate. A nonempty code must be returned by Catalog.Lookup and its
Category must match the category configured on the constraint. Missing entry ->
ErrUnknown; wrong category -> ErrCategory, with {{ category }} set to the required
category. Catalog errors are technical errors with preserved errors.Is identity.
The constraint skips nil and empty values without querying the catalog. Signup
requires Code under code; Transfer requires Destination under destinationCode.
Each DTO also requires Label. Requiredness is separate from the reusable rule.
Use the supplied constraint for the code field, accumulate independent violations,
and preserve the caller validator (path, locale, factory) and context.
ValidatorOptions registers Russian messages: 'Код не найден.' and
'Код должен относиться к категории {{ category }}.' English defaults are in the
contract. Caller constructs its validator with these options. Support repeated
and concurrent reuse of the same constraint with different caller paths/locales,
without mutation or retained per-call state. One lookup per nonempty validation.
No fluent builder, caching, or extra dependencies are required.
Preserve protected files and exported contracts. Edit the supplied implementation
files; you may add private helpers in new top-level Go files (not *_test.go).
Do not change dependencies. Use existing library/project rules when their semantics
match. Hidden tests assess behavior; later code review separately assesses reuse,
rule ownership, and unnecessary abstractions. You may run temporary tests, but
remove them before submission. Paths are PropertyPath.String() JavaScript-style
property paths, not JSON Pointer. Preserve any caller prefix.
