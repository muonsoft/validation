# ISO weeks

Use `it.IsWeek()` for strings such as `2020-W53`. Optional `WithMin` and `WithMax`
set inclusive bounds, including ranges across ISO years:

```go
err := validator.Validate(context.Background(),
    validation.String(week, it.IsWeek().WithMin("2020-W01").WithMax("2021-W01")),
)
```

The format is exactly `YYYY-Www`: four ASCII year digits, uppercase `W`, and a
zero-padded week from `01` to `53`. Week 53 is accepted only in ISO years that
have it. Years `0000`–`9999` are interpreted literally in the proleptic Gregorian
calendar, with no two-digit-year remapping. No whitespace is trimmed.

The behavior is based on [Symfony Week](https://symfony.com/doc/current/reference/constraints/Week.html).
Unlike Symfony, empty strings are skipped, following this library's convention;
nil pointers are also skipped. Use `it.IsNotBlank()` when a value is required.
Literal early years follow Go's calendar, without PHP `mktime`'s remapping of
some years below 100. No local timezone or new dependencies are involved.

Standalone helpers are `validate.Week` and `is.Week`, with optional
`validate.WithMinWeek` and `validate.WithMaxWeek`. Empty bounds remove restrictions;
malformed, nonexistent, or reversed bounds produce `validate.ErrInvalidWeekBounds`
(or a constraint configuration error through `it`), even for empty values.

The four violation errors are `validation.ErrInvalidWeekFormat`,
`validation.ErrInvalidWeekNumber`, `validation.ErrWeekTooEarly`, and
`validation.ErrWeekTooLate`. Customize them with `WithFormatError`,
`WithWeekNumberError`, `WithMinError`, and `WithMaxError`; corresponding `Message`
methods accept `{{ value }}`, `{{ min }}`, `{{ max }}`, and custom parameters.
The constraint supports `When`, `WhenGroups`, `This`, `Each`, and English/Russian
translations. Disabled constraints and unmatched groups skip validation entirely.

See [Usage](usage.md) for imports and validator setup, and the
[constraint catalog](constraints.md) for related checks.
