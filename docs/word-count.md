# Word count

Use `it.HasMinWordCount(min)`, `it.HasMaxWordCount(max)`, or
`it.HasWordCountBetween(min, max)` to check inclusive word count bounds:

```go
err := validator.Validate(context.Background(),
    validation.String(content, it.IsNotBlank(), it.HasWordCountBetween(10, 200)),
)
```

A word starts with a Unicode letter or number and continues with letters, numbers,
or combining marks. Internal apostrophes (`'` and `’`) stay within a word;
hyphens, underscores, and other punctuation separate words. For example,
`don't` is one word, `well-known` is two, and `123` is one. Emoji and punctuation
alone count as zero words. Nil and empty strings are skipped; whitespace-only
strings have zero words. Use `it.IsNotBlank()` to require a value.

This dependency-free approximation is inspired by
[Symfony WordCount](https://symfony.com/doc/current/reference/constraints/WordCount.html),
but does not implement ICU's locale-aware segmentation. Unspaced Chinese, Japanese,
or Thai text is counted as one continuous word, and punctuation rules can differ
(for example, `3.14` counts as two words). Locale selection is not supported.

Bounds must be non-negative, with minimum no greater than maximum; invalid bounds
produce a constraint configuration error. Equal bounds require an exact count.
Violations use `validation.ErrTooFewWords` or `validation.ErrTooManyWords`.
Use `WithMinError`/`WithMaxError` and `WithMinMessage`/`WithMaxMessage` to customize
them; message parameters are `{{ count }}`, `{{ limit }}`, and `{{ value }}`
(the quoted input). Conditional validation, groups, and English/Russian plural
forms are supported.

See [Usage](usage.md) for imports and validator setup, and the
[constraint catalog](constraints.md) for related checks.
