# Password strength

`it.HasPasswordStrength()` requires a score of at least Medium (2):

```go
err := validator.Validate(context.Background(),
    validation.String(password,
        it.HasPasswordStrength().WithMinScore(validate.PasswordStrengthStrong),
    ),
)
```

`validate.EstimatePasswordStrength(password)` returns a `PasswordStrengthScore`:
`PasswordStrengthVeryWeak` (0), `PasswordStrengthWeak` (1),
`PasswordStrengthMedium` (2), `PasswordStrengthStrong` (3), or
`PasswordStrengthVeryStrong` (4). Minimum scores must be 1–4. Standalone
`validate.PasswordStrength` accepts
`validate.WithMinPasswordStrength` and `validate.WithPasswordStrengthEstimator`.

The estimator follows [Symfony 8.0](https://github.com/symfony/validator/blob/8.0/Constraints/PasswordStrengthValidator.php):
`E = U*log2(P) + (L-U)*log2(U)`, using byte length `L`, distinct bytes `U`,
and the combined sizes `P` of present byte categories. Categories are lowercase
ASCII (26), uppercase ASCII (26), digits (10), symbols including space (33),
control bytes (33), and non-ASCII bytes (128). Thresholds are 60, 80, 100, and 120.
UTF-8 bytes are counted separately without normalization; arbitrary byte strings
are accepted. Empty input has score 0.

This is a heuristic, not a guarantee of resistance to guessing. There are no
wordlists or sequence checks: the lowercase alphabet scores 4 and `ab` repeated
50 times scores 3. It is neither zxcvbn nor a Have I Been Pwned breach check.

`WithEstimator(func(string) validate.PasswordStrengthScore)` replaces the algorithm;
return values outside 0–4 are configuration errors. Nil restores the default.
Custom estimators must support concurrent calls when constraints are shared.
The estimator receives the original string unchanged and is called once per
validated non-nil value. Nil pointers are skipped; **empty strings are evaluated
and rejected by the default estimator**, unlike most string constraints. A custom
estimator also controls the score of empty strings. Invalid minimum scores are
configuration errors even for nil input. Disabled constraints and unmatched groups
skip validation entirely.

Violations use `validation.ErrPasswordTooWeak`, with English/Russian translations.
`WithError`, `WithMessage`, `When`, `WhenGroups`, `This`, and `Each` are supported.
Automatic message parameters are only `{{ strength }}` and `{{ minScore }}`:
the library does not add the password or a `{{ value }}` parameter to violations.
Avoid including passwords in custom messages, parameters, or estimator logs.
Standalone failures use `validate.ErrPasswordTooWeak`,
`validate.ErrInvalidPasswordStrengthMinimum`, or `validate.ErrInvalidPasswordStrengthScore`.
No dependencies or network access are required.

See [Usage](usage.md) for imports and validator setup, and the
[constraint catalog](constraints.md) for related checks.
