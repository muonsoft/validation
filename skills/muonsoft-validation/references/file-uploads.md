# File uploads (v0.20.0+)

<!-- Minimum validation version: v0.20.0 -->

Validate filename, extension, declared MIME metadata, and detected content as
separate properties. `IsMIMEType` checks the declared string; `HasContentType`
checks bytes using `http.DetectContentType` by default (at most the first 512
bytes). Detection does not prove that a file is complete or structurally valid.
Two independent extension/content allowlists do not enforce matching pairs;
select content types from a policy for the validated extension when required.

```go
package main

import (
    "context"
    "fmt"
    "io"
    "strings"

    "github.com/muonsoft/validation"
    "github.com/muonsoft/validation/it"
)

func checkUpload(ctx context.Context, v *validation.Validator, name string, body io.Reader) error {
    if err := v.Validate(ctx, validation.StringProperty("name", name,
        it.IsNotBlank(), it.IsFileName().WithWindowsRestrictions(),
        it.HasFileExtension("txt"),
    )); err != nil { return err }
    prefix, err := io.ReadAll(io.LimitReader(body, 512))
    if err != nil { return err }
    return v.Validate(ctx,
        validation.NumberProperty("contentSize", len(prefix), it.IsGreaterThan(0)),
        validation.This(prefix, it.HasContentType("text/plain")).At(validation.PropertyName("content")),
    )
}

func main() {
    v, err := validation.NewValidator()
    if err != nil { panic(err) }
    fmt.Println(checkUpload(context.Background(), v, "notes.txt", strings.NewReader("Hello")))
}
// Output:
// <nil>
```

The caller owns the reader and any close operation. This example consumes a
prefix: reopen the upload or replay the prefix with `io.MultiReader` before
saving. Enforce a total upload byte limit at the receiving boundary; a sniffing
limit is not an upload size limit. Read failures are technical errors, not invalid
content. Reuse the supplied validator so prefixes/translations survive both stages.

Empty names/content are skipped by format constraints; require them separately.
`HasFileExtension` supports compound suffixes and checks basic filename validity;
Windows device names and trailing dots/spaces require `WithWindowsRestrictions`.
`WithDetector(func([]byte) string)` adapts an existing content detector without
adding a dependency to the library. A detector must not mutate the bytes and
must be safe for concurrent calls when shared. An empty/malformed detector result
or invalid allowlist is a configuration error, not a user-data violation.
