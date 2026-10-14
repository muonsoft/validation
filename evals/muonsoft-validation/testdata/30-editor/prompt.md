# Refactor article validation without changing its error contract

Replace the hand-written field checks in Article.Validate with the library's
existing rules where they fit, and correct the incomplete implementation.
The editor runs with the default group; publishing may select only "publish"
or both default and "publish". Both modes require Title and enforce 3–12 Unicode
code points inclusive. An empty title produces only ErrTitleRequired; nonempty
out-of-range titles produce ErrTitleLength. Do not trim: three spaces are valid.
Contact is optional in the editor, but required when "publish" is selected.
A supplied contact must satisfy the library's ordinary email-address policy.
Empty contact on publish produces only ErrContactRequired; malformed nonempty
contact produces ErrContactFormat. Each violation uses its sentinel's Message().
When neither default nor publish is selected, none of these rules run.
Keep independent title/contact failures together and preserve the caller's prefix.
External field names are title and contact. Do not alter the exported DTO/errors.

Preserve contract.go, existing protected files, and dependencies. Edit validation.go.
Private helpers in new top-level non-test Go files are allowed. Remove temporary
tests before submission. Use existing library/project rules when their semantics
match. Do not change input, normalize values, or create a replacement validator.
