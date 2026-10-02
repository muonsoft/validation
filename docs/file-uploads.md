# File upload validation

File names, extensions, declared MIME metadata, and detected content are separate
properties. These checks use only the standard library and existing dependencies.
They perform no I/O and do not change, decode, or sanitize input.

| Constraint | Input | Standalone function |
|------------|-------|---------------------|
| `it.IsFileName()` | `string` | `validate.FileName(name, options...)` |
| `it.HasFileExtension("jpg", "png")` | `string` | `validate.FileExtension(name, extensions...)` |
| `it.IsMIMEType(types...)` | `string` | `validate.MIMEType(value, types...)` |
| `it.HasContentType(types...)` | `[]byte` | `validate.ContentType(data, types, options...)` |

String constraints work with `validation.String`, `NilString`, `This`, and `Each`.
Content constraints work with `This` and `Each`. All support `WithError`,
`WithMessage`, `When`, and `WhenGroups`. No input values are automatically included
in violation messages or parameters, including file content.

Nil and empty inputs are skipped. Require names with `it.IsNotBlank()` and content
with a separate size check. A valid configuration is required even for empty
inputs; disabled constraints skip configuration checks too.

## Names and extensions

`IsFileName` rejects malformed UTF-8, Unicode control characters, `/`, `\`, `.` and
`..`, and ASCII drive prefixes such as `C:`. Unicode, spaces and dotfiles are
otherwise accepted. There is no universal filename length limit or filesystem
access check. `it.HasMaxLength` counts characters; use `validation.Number(len(name),
it.IsLessThanOrEqual(255))` when a particular storage policy requires a byte limit.

`WithWindowsRestrictions()` additionally rejects `< > : " | ? *`, trailing ASCII
spaces or periods, and device names such as `CON`, `NUL`, `COM1`, `LPT9`, `COM¹`,
`LPT²`, `CONIN$`, and `CONOUT$`, including names with extensions. Device stems are
case-insensitive and trailing spaces before an extension are ignored for this
check. The policy is conservative and identical on all operating systems; it does
not query the Windows version or registry. Standalone use:
`validate.FileName(name, validate.WithWindowsFileNameRestrictions())`.
See [Microsoft's naming rules](https://learn.microsoft.com/en-us/windows/win32/fileio/naming-a-file).

Extension allowlists accept an optional leading dot and compare ASCII letters
case-insensitively. Each dot-separated segment must be nonempty and contain only
ASCII letters, digits, `_`, `+`, or `-`. Compound suffixes are supported:
`tar.gz` matches `backup.TAR.GZ`, but not `backup.tar.gz.exe`. A suffix needs a
nonempty stem: `.jpg` does not match `jpg`. `report.final.pdf` is allowed by `pdf`.
Names must pass the basic filename check; paths are never reduced to their basename.
Windows restrictions remain a separate constraint. Empty or malformed allowlists
are configuration errors, not violations.

## MIME metadata and detected content

`IsMIMEType()` checks syntax; arguments add an allowlist. MIME values must have a
concrete `type/subtype`. Valid parameters such as `charset=utf-8` are parsed and
ignored for matching; malformed parameters are rejected. Allowlist entries must
be concrete types without parameters. Matching ignores case and surrounding
whitespace, but does not resolve aliases, wildcard patterns, or parent types.
There is no registry lookup: a syntactically valid private type is accepted.

`HasContentType` requires a nonempty allowlist and defaults to
[`http.DetectContentType`](https://pkg.go.dev/net/http#DetectContentType), which
examines at most the first 512 bytes. Unknown binary content returns
`application/octet-stream` and passes only if explicitly allowed. Detection is
not integrity validation: even a PNG signature alone can match `image/png`.

The `Content-Type` supplied by a client is untrusted metadata. Matching it cannot
replace content detection. Neither check establishes that a file is harmless or
complete. Decode images, inspect archives, or parse documents when structural
validity matters; apply resource limits appropriate to those parsers.
See [OWASP's file upload guidance](https://cheatsheetseries.owasp.org/cheatsheets/File_Upload_Cheat_Sheet.html).

## Multipart uploads and bounded reading

Limit the HTTP request body with `http.MaxBytesReader` before parsing multipart
input. `ParseMultipartForm`'s memory limit is not a total upload size limit.
Check `FileHeader.Size` against your application limit before reading content.
A streamed upload needs an enforced byte limit while receiving it; do not rely on
a client-declared size. Arrange `MultipartForm.RemoveAll()` after processing to
clean up temporary files.

Read a bounded prefix using a separate handle:

```go
func uploadPrefix(header *multipart.FileHeader, limit int64) ([]byte, error) {
    file, err := header.Open()
    if err != nil {
        return nil, err
    }
    defer file.Close()
    return io.ReadAll(io.LimitReader(file, limit))
}
```

Use a positive application-controlled limit, normally 512 for the default detector.
Validate the result only after checking the read error. Open the `FileHeader` again
when saving the upload; no seek or shared stream offset is involved. Go's multipart
parser may already reduce the original filename to a basename; validation applies
to the name it actually supplies.

```go
// header is a non-nil *multipart.FileHeader; v is a *validation.Validator.
err := v.Validate(ctx,
    validation.StringProperty("filename", header.Filename,
        it.IsNotBlank(),
        it.IsFileName().WithWindowsRestrictions(),
        it.HasFileExtension("jpg", "jpeg", "png"),
    ),
    validation.NumberProperty("size", header.Size,
        it.IsGreaterThan(int64(0)),
        it.IsLessThanOrEqual(int64(10 << 20)),
    ),
)
if err != nil {
    return err
}
prefix, err := uploadPrefix(header, 512)
if err != nil {
    return err // An I/O error is not a validation violation.
}
return v.Validate(ctx,
    validation.NumberProperty("contentSize", len(prefix), it.IsGreaterThan(0)),
    validation.This(prefix, it.HasContentType("image/jpeg", "image/png")).
        At(validation.PropertyName("content")),
)
```

For a one-pass reader, preserve the consumed prefix before passing it downstream:

```go
prefix, err := io.ReadAll(io.LimitReader(input, 512))
if err != nil {
    return err
}
// Validate prefix here, then use the reconstructed stream for saving.
replayed := io.MultiReader(bytes.NewReader(prefix), input)
```

The caller owns and closes the original reader. Do not re-read it directly after
sniffing: the prefix has already been consumed. Do not validate concurrently with
modifying the slice or consuming the same stream.

## Custom detectors and mimetype

`WithDetector(func([]byte) string)` replaces detection. It receives the entire
supplied slice, must not modify it, and must support concurrent calls if shared.
Nil restores the standard detector. A malformed/empty result is a constraint
configuration error; a valid but disallowed type is a violation.
Standalone configuration uses `validate.WithContentTypeDetector(detector)`.

Applications needing more formats can install
[`github.com/gabriel-vasile/mimetype`](https://pkg.go.dev/github.com/gabriel-vasile/mimetype)
in their own module. This library does not import it, including in tests:

```go
rule := it.HasContentType("image/png", "application/pdf").
    WithDetector(func(data []byte) string {
        return mimetype.Detect(data).String()
    })
```

`Detect` returns `*mimetype.MIME`, so a short adapter is needed. For containers such
as DOCX/XLSX, provide enough data for the external detector; 512 bytes are not a
universal detection budget. Its `SetLimit` is package-global: configure it once
at application startup if needed, not inside a constraint. Keep input bounded;
a truncated sample cannot establish what the rest of a file contains.

The external package recommends `MIME.Is` for aliases. Normalize selected aliases
explicitly when using this constraint:

```go
rule := it.HasContentType("application/zip").WithDetector(func(data []byte) string {
    detected := mimetype.Detect(data)
    if detected.Is("application/zip") {
        return "application/zip"
    }
    return detected.String()
})
```

This does not accept every descendant container as ZIP. If an application already
uses `mimetype.DetectReader`, handle its I/O error and then validate the returned
MIME string with `IsMIMEType`. Preserve or reopen the stream as above.

## Matching extensions to content

Two independent allowlists do not reject a PNG called `photo.jpg` when both PNG
and JPEG are allowed. Select the expected content types from an explicit policy
for the validated extension:

```go
allowed := map[string][]string{
    ".jpg": {"image/jpeg"},
    ".jpeg": {"image/jpeg"},
    ".png": {"image/png"},
}
// Run name/extension validation first so lookup cannot yield an empty policy.
types := allowed[strings.ToLower(path.Ext(header.Filename))]
return v.Validate(ctx, validation.This(prefix, it.HasContentType(types...)))
```

For compound extensions, match the longest configured suffix explicitly.
Do not derive this policy from `mime.TypeByExtension`: its results depend on
[system MIME databases or the Windows registry](https://pkg.go.dev/mime#TypeByExtension).
Generate storage names in the application; validation does not make the supplied
name a safe storage path or authorize serving active content.

## Errors

Violations use `validation.ErrInvalidFileName`, `ErrInvalidFileExtension`,
`ErrInvalidMIMEType`, and `ErrInvalidContentType`, with English/Russian messages.
Standalone checks return the corresponding `validate` errors. Invalid options
return `validate.ErrInvalidFileExtensions` or `validate.ErrInvalidMIMETypes`;
invalid detector output returns `validate.ErrInvalidDetectedContentType`.
Constraints convert these configuration/contract failures to constraint errors,
which stop validation instead of being collected as user-data violations.
