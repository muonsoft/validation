# Validate a text upload before saving

Implement Upload.Validate for a nonnil body reader owned by the caller.
Name is required, must satisfy the library's filename policy including Windows
restrictions, and must end in .txt (case-insensitive). MIME is required and must
be syntactically valid text/plain, allowing valid MIME parameters and ignoring
case. Accumulate metadata failures at name and mime before reading any bytes.
Use standard library violation identities for requiredness/name/extension/MIME;
independent name and extension rules may both fail at name.
Only after metadata passes, read at most 512 bytes without closing the reader.
Preserve any read error as a technical error. Empty content produces ErrIsBlank at
content. Nonempty bytes must be detected as text/plain using the library's default
content detector, otherwise ErrInvalidContentType at content. The caller enforces
the total request size and will reopen the body for saving; do not read the entire
upload, trust MIME as proof of content, or claim structural file validation.

Preserve contract.go and dependencies. Edit validation.go; private helpers in new
top-level non-test Go files are allowed. Remove temporary tests before submission.
Use existing library/project rules when their semantics match. Do not mutate input.
Preserve the supplied validator and its path prefix.
