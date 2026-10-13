# Match image attachments to their declared format

Implement Attachment.Validate for a nonnil reader owned by the caller.
Accept PNG and JPEG attachments: names must be nonempty valid portable filenames
including Windows restrictions and end in .png, .jpg, or .jpeg, ignoring suffix
case. MIME is required, allows valid parameters/case differences, and must be
image/png or image/jpeg. Declared Size must be 1–1048576 inclusive; use ErrSize and
its message for either bound. The receiving boundary separately enforces actual
upload size; this method must not consume the whole upload to measure it.
Accumulate independent filename, extension, MIME, and size errors before I/O.
Use standard required/name/extension/MIME errors at name or mime and ErrSize at
size. Once metadata is valid, the MIME must also agree with the extension: PNG
requires image/png; both JPEG suffixes require image/jpeg. A mismatch is the
standard invalid-MIME error at mime and must not read the body.
Only then read at most 512 bytes. Propagate read errors unchanged, including an
error accompanied by bytes. Empty data yields the standard blank-value error at
content. Otherwise detect content using the library's default detector and require
the type implied by the extension; a mismatch uses the standard content-type error
at content. Do not close or rewind the reader. The caller reopens it for saving.
Preserve path prefixes. Matching signatures do not promise a complete image.

Preserve contract.go, existing protected files, and dependencies. Edit validation.go.
Private helpers in new top-level non-test Go files are allowed. Remove temporary
tests before submission. Use existing library/project rules when their semantics
match. Do not change input, normalize values, or create a replacement validator.
