# Go validation library

[![Go Reference](https://pkg.go.dev/badge/github.com/muonsoft/validation.svg)](https://pkg.go.dev/github.com/muonsoft/validation)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/muonsoft/validation)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/muonsoft/validation)
![GitHub](https://img.shields.io/github/license/muonsoft/validation)
[![tests](https://github.com/muonsoft/validation/actions/workflows/tests.yml/badge.svg)](https://github.com/muonsoft/validation/actions/workflows/tests.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/muonsoft/validation)](https://goreportcard.com/report/github.com/muonsoft/validation)
[![Scrutinizer Code Quality](https://scrutinizer-ci.com/g/muonsoft/validation/badges/quality-score.png?b=main)](https://scrutinizer-ci.com/g/muonsoft/validation/?branch=main)
[![Scrutinizer Code Coverage](https://scrutinizer-ci.com/g/muonsoft/validation/badges/coverage.png?b=main)](https://scrutinizer-ci.com/g/muonsoft/validation/?branch=main)
[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-2.0-4baaaa.svg)](CODE_OF_CONDUCT.md)

Validate Go values with typed, composable constraints. Define rules in code,
validate nested structures, and return translated violations with property paths.

This project is inspired by [Symfony Validator component](https://symfony.com/index.php/doc/current/validation.html).

## Key features

* Typed API using Go generics
* Declarative style of describing a validation process in code
* Validation of different types: booleans, numbers, strings, slices, maps, and time
* Validation of custom data types that implement the `Validatable` interface
* Customizable validation errors with translations and pluralization supported out of the box
* Custom validation rules with context propagation and message translations

## Version policy

The library is in the v0 series. Minor releases (`0.x.0`) may introduce breaking
changes; patch releases contain bug fixes. Review the [changelog](CHANGELOG.md)
before upgrading.

## Quick start

### Installation

Requires Go 1.24.0 or later. See [Installation](docs/installation.md).

```bash
go get github.com/muonsoft/validation
```

### Basic example

Validate a string and a number, then inspect each violation:

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/muonsoft/validation"
    "github.com/muonsoft/validation/it"
    "github.com/muonsoft/validation/validator"
)

func main() {
    err := validator.Validate(context.Background(),
        validation.String("not-an-email", it.IsEmail()),
        validation.Number(15, it.IsGreaterThanOrEqual(18)),
    )
    if violations, ok := validation.UnwrapViolations(err); ok {
        for _, violation := range violations.All() {
            fmt.Println(violation.Message())
        }
    } else if err != nil {
        log.Fatal(err)
    }
}
// Output:
// This value is not a valid email address.
// This value should be greater than or equal to 18.
```

### Conceptual example

Declarative validation with nested structs, property paths, and reusable `Validatable` types:

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/muonsoft/validation"
    "github.com/muonsoft/validation/it"
    "github.com/muonsoft/validation/validator"
)

type Product struct {
    Name       string
    Tags       []string
    Components []Component
}

func (p Product) Validate(ctx context.Context, v *validation.Validator) error {
    return v.Validate(ctx,
        validation.StringProperty("name", p.Name, it.IsNotBlank()),
        validation.AtProperty("tags",
            validation.Countable(len(p.Tags), it.HasMinCount(2)),
            validation.Comparables[string](p.Tags, it.HasUniqueValues[string]()),
            validation.EachString(p.Tags, it.IsNotBlank()),
        ),
        validation.AtProperty("components",
            validation.Countable(len(p.Components), it.HasMinCount(1)),
            validation.ValidSlice(p.Components),
        ),
    )
}

type Component struct {
    Name string
}

func (c Component) Validate(ctx context.Context, v *validation.Validator) error {
    return v.Validate(ctx,
        validation.StringProperty("name", c.Name, it.IsNotBlank()),
    )
}

func main() {
    err := validator.ValidateIt(context.Background(), Product{
        Name: "",
        Tags: []string{"a", "", "a"},
        Components: []Component{{Name: ""}},
    })
    if violations, ok := validation.UnwrapViolations(err); ok {
        for _, violation := range violations.All() {
            fmt.Println(violation)
        }
    } else if err != nil {
        log.Fatal(err)
    }
}
// Output:
// violation at "name": "This value should not be blank."
// violation at "tags": "This collection should contain only unique elements."
// violation at "tags[1]": "This value should not be blank."
// violation at "components[0].name": "This value should not be blank."
```

See [Usage](docs/usage.md) for validator setup and validation arguments.

## Documentation

| Topic | Description |
|-------|-------------|
| [Installation](docs/installation.md) | How to install the package |
| [Constraint catalog](docs/constraints.md) | Find constraints by input type and standalone helpers |
| [Usage](docs/usage.md) | Basic concepts, validator, validation arguments |
| [Property paths & structs](docs/property-paths-and-structs.md) | Property paths, struct validation, conditional validation, groups |
| [Violations and errors](docs/violations-and-errors.md) | Handling violations, error structure, storing in database |
| [Country, language & locale](docs/international.md) | Formats, special codes, and application-support limitations |
| [File uploads](docs/file-uploads.md) | Names, extensions, MIME metadata, content detection, and custom detectors |
| [UTF-8 validation](docs/utf8.md) | Reject malformed text at input boundaries; scope and limitations |
| [Translations](docs/translations.md) | Multi-language support and custom messages |
| [Custom constraints](docs/custom-constraints.md) | Creating your own constraints |

Full index: [docs/README.md](docs/README.md).

API reference: [pkg.go.dev/github.com/muonsoft/validation](https://pkg.go.dev/github.com/muonsoft/validation)

## Contributing

You may help this project by

* reporting an [issue](https://github.com/muonsoft/validation/issues);
* making translations for error messages;
* suggesting an improvement or [discussing](https://github.com/muonsoft/validation/discussions) the usability of the package.

If you'd like to contribute, see [the contribution guide](CONTRIBUTING.md). Pull requests are welcome.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
