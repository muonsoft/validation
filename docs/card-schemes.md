# Payment card schemes

Select one or more schemes explicitly with `it.IsCardScheme`:

```go
err := validator.Validate(context.Background(),
    validation.String(cardNumber,
        it.IsCardScheme(validate.CardSchemeVisa, validate.CardSchemeMastercard),
        it.IsLUHN(),
    ),
)
```

The constraint checks prefixes and lengths using the
[Symfony 8.0 CardScheme rules](https://github.com/symfony/validator/blob/8.0/Constraints/CardSchemeValidator.php).
It does not check the checksum, card existence, or ability to make payments.
`it.IsLUHN()` adds a separate checksum check when appropriate for the selected schemes.
Only ASCII digits are accepted; spaces and hyphens are not removed.

Supported constants in `validate` are `CardSchemeAMEX`, `CardSchemeChinaUnionPay`,
`CardSchemeDiners`, `CardSchemeDiscover`, `CardSchemeInstaPayment`, `CardSchemeJCB`,
`CardSchemeLaser`, `CardSchemeMaestro`, `CardSchemeMastercard`, `CardSchemeMIR`,
`CardSchemeUATP`, and `CardSchemeVisa`. Multiple schemes use OR semantics;
overlapping prefixes can match more than one system. These are the pinned Symfony
rules, including its legacy scheme formats, rather than a live issuer directory.

Standalone helpers `validate.CardScheme` and `is.CardScheme` accept the same
variadic `validate.CardSchemeName` constants. Nil/empty values are skipped;
combine with `it.IsNotBlank()` to require a number. An empty selection or any
unknown scheme returns `validate.ErrInvalidCardSchemes` (a configuration error
through `it`), including for empty input. This is stricter than Symfony's handling
of unknown names. Disabled constraints or unmatched groups skip validation entirely.

Invalid input returns `validation.ErrInvalidCardScheme` (or
`validate.ErrInvalidCardScheme` from the standalone validator). Nonnumeric input
and a nonmatching prefix/length share this error. `WithError` and `WithMessage`
customize violations; `{{ value }}` contains the original input. English/Russian
messages, `When`, `WhenGroups`, `This`, and `Each` are supported.

See [Usage](usage.md) for imports and validator setup, and the
[constraint catalog](constraints.md) for related checks.
