# Installation

Requires **Go 1.24.0 or later**, as declared in [go.mod](../go.mod). Use a
maintained Go toolchain with current patches when deploying applications.

From your application's Go module, run:

```bash
go get github.com/muonsoft/validation
```

For a reproducible upgrade, append `@vX.Y.Z` with the release you want. Review
[CHANGELOG.md](../CHANGELOG.md) before upgrading: v0 minor releases may include
breaking changes.

A typical application imports:

```go
import (
    "context"

    "github.com/muonsoft/validation"
    "github.com/muonsoft/validation/it"
    "github.com/muonsoft/validation/validator"
)
```

`validation` supplies typed arguments and validator configuration, `it` supplies
constraints, and `validator` supplies the default validator. For standalone
checks, use `is` (boolean result) or `validate` (error result).

Continue with the [runnable quick start](../README.md#basic-example) or
[Usage](usage.md) for configuration and empty-value behavior.
