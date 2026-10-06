# Validate independent item lookups
Implement Request.Validate(ctx, validator, lookup). Items may be empty. Each item
requires ID (ErrIsBlank) and positive Quantity (ErrNotPositive), at items[i].id
and items[i].quantity. Accumulate both local violations. Only an item whose local
rules pass calls Exists exactly once with its original ID and original context;
do not trim or deduplicate IDs. A false result yields ErrMissing at items[i].id.
Invalid or missing items must not prevent valid later siblings from being checked.
Local rules apply in default and submit groups; callers select either or both.
Visit items in original order. A technical lookup error stops all later work and
must remain discoverable with errors.Is; return it even if earlier items produced
violations. Preserve caller prefix, violation factory, and groups. Do not modify
input or perform lookups on invalid prerequisites.

Preserve contract.go and dependencies. Edit validation.go; private helpers in new
top-level non-test Go files are allowed. Remove temporary tests before submission.
Use existing library/project rules when their semantics match.
