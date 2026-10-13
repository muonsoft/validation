# Password strength

Use `it.HasPasswordStrength` for a minimum score, not a hand-written mix of length
and character-class rules unless those are the application's explicit contract.
The default estimator is a byte-based heuristic, with scores 0–4 and a default
minimum of Medium (2). It is not a breach check or dictionary/sequence detector.
Do not trim, normalize, or otherwise change a password during validation.

```go
package main

import (
    "context"
    "errors"
    "fmt"

    "github.com/muonsoft/validation"
    "github.com/muonsoft/validation/it"
    "github.com/muonsoft/validation/validate"
)

func main() {
    v, err := validation.NewValidator()
    if err != nil { panic(err) }
    rule := it.HasPasswordStrength().WithMinScore(validate.PasswordStrengthStrong)
    err = v.Validate(context.Background(), validation.StringProperty("password", "", rule))
    fmt.Println(errors.Is(err, validation.ErrPasswordTooWeak))
    err = v.Validate(context.Background(), validation.StringProperty("password", "abcdefghijklmnopqrstuvwxyz", rule))
    fmt.Println(err == nil)
}
// Output:
// true
// true
```

The alphabet example demonstrates the heuristic's limitation, not a recommended
password. Nil is skipped; an empty string is scored and rejected by the default
estimator. Adding `IsNotBlank` can produce an additional violation for empty
input, so choose the error contract deliberately. Optional password updates need
an explicit presence policy.

`WithEstimator(func(string) validate.PasswordStrengthScore)` supports an existing
application estimator. It receives the unchanged value, including empty strings;
its score must be 0–4, and the minimum must be 1–4. Invalid scores/options are
configuration errors. Shared estimators must be safe for concurrent calls.
Built-in violations contain strength/minimum parameters, not the password.
Keep passwords out of custom messages, parameters, and logs as well.
