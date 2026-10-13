# Validate a catalog tree
Implement Catalog.Validate(ctx, validator, maxDepth) for an acyclic tree stored
in Roots. Roots may be empty. Each root has depth 0; maxDepth is nonnegative and
inclusive. Each map entry is a node: nil yields ErrMissing at that entry path.
A nonnil node beyond maxDepth yields only ErrDepth there; do not inspect its
label or descendants. Other branches must still be validated. Within the depth
limit Label is required (ErrIsBlank at label), and Children are recursively
validated even when Label is invalid. Empty Children is valid. Labels need not
be unique, either among siblings or across branches. Never mutate the input.
External paths use roots, children and label. All map keys are literal property
segments, including empty, numeric and punctuated keys; preserve the caller's
prefix and use library path rendering. No sorting or error-list order required.

Preserve contract.go and dependencies. Edit validation.go; private helpers in new
top-level non-test Go files are allowed. Remove temporary tests before submission.
Use existing library/project rules when their semantics match.
