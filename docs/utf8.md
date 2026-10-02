# UTF-8 validation

Go strings can contain arbitrary bytes. Use `it.IsUTF8()` when an input contract
requires UTF-8 and the source does not already guarantee it, for example file
imports, legacy integrations, or raw text payloads. The constraint rejects
malformed or truncated byte sequences before storage or further processing might
reject them or silently replace them with the replacement character (`�`). It
reports a validation violation that can be associated with the input field.

```go
err := validator.Validate(ctx,
    validation.StringProperty("text", input, it.IsUTF8()),
)
```

For standalone checks, `is.UTF8(input)` returns a boolean and
`validate.UTF8(input)` returns `nil` or `validate.ErrInvalidUTF8`.
The constraint produces `validation.ErrInvalidUTF8`, with English and Russian
messages and the usual custom-message, condition, and group support.

## What passes and fails

| Input | Result |
|-------|--------|
| ASCII, Cyrillic, Japanese, emoji | Pass |
| Empty string; nil value for the constraint | Pass; add `it.IsNotBlank()` to require a value |
| Truncated sequence such as `"\xe2\x82"` | Fail |
| Invalid byte such as `"\xff"` | Fail |
| Overlong encodings, encoded surrogates, values above U+10FFFF | Fail |
| HTML, NUL, invisible characters, combining marks | Pass if correctly encoded |
| Literal replacement character `�` or validly encoded mojibake such as `Ã©` | Pass |

## Scope and limitations

This implements the **UTF-8-only** scope of the Charset feature. Unlike Symfony's
configurable Charset constraint, it does not accept a choice of encodings. It uses
Go's standard `unicode/utf8.ValidString` and requires no new dependencies.

The check validates encoding, not text meaning, language, readability, or safety.
It does not sanitize HTML, detect suspicious characters, normalize Unicode,
repair corrupted text, or identify or convert another encoding. Input from
another encoding fails only when its bytes are not also valid UTF-8; ASCII bytes,
for example, are shared by many encodings. Binary input that happens to form
valid UTF-8 also passes.

Validate the original bytes, converted directly to a Go string, **before** any
lossy decoding or replacement. Once another layer replaces invalid bytes with
`�`, this check cannot recover or identify the original corruption. If an earlier
layer already guarantees valid UTF-8, repeating this check adds little value.
