# Golang validation framework

[![Go Reference](https://pkg.go.dev/badge/github.com/muonsoft/validation.svg)](https://pkg.go.dev/github.com/muonsoft/validation)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/muonsoft/validation)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/muonsoft/validation)
![GitHub](https://img.shields.io/github/license/muonsoft/validation)
[![tests](https://github.com/muonsoft/validation/actions/workflows/tests.yml/badge.svg)](https://github.com/muonsoft/validation/actions/workflows/tests.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/muonsoft/validation)](https://goreportcard.com/report/github.com/muonsoft/validation)
[![Scrutinizer Code Quality](https://scrutinizer-ci.com/g/muonsoft/validation/badges/quality-score.png?b=main)](https://scrutinizer-ci.com/g/muonsoft/validation/?branch=main)
[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-2.0-4baaaa.svg)](CODE_OF_CONDUCT.md)

Golang validation framework based on static typing and generics. Designed to create complex validation rules with
abilities to hook into the validation process.

This project is inspired by [Symfony Validator component](https://symfony.com/index.php/doc/current/validation.html).

## Key features

* Flexible and customizable API built in mind to use benefits of static typing and generics
* Declarative style of describing a validation process in code
* Validation of different types: booleans, numbers, strings, slices, maps, and time
* Validation of custom data types that implements `Validatable` interface
* Customizable validation errors with translations and pluralization supported out of the box
* Easy way to create own validation rules with context propagation and message translations

## Work-in-progress notice

This package is under active development and API may be changed until the first major version will be released. Minor
versions `n` 0.n.m may contain breaking changes. Patch versions `m` 0.n.m may contain only bug fixes.

Goals before stable release:

* [x] implementation of static type arguments by generics;
* [x] mechanism for asynchronous validation (lazy violations by async/await pattern);
* [ ] implement all common constraints.

## Quick start

### Installation

```bash
go get -u github.com/muonsoft/validation
```

### Basic example

Validation of string and number with typed constraints:

```go
err := validator.Validate(context.Background(),
    validation.String("not-an-email", it.IsEmail()),
    validation.Number(15, it.IsGreaterThanOrEqual(18)),
)
// violations:
//   This value is not a valid email address.
//   This value should be 18 or more.
```

### Conceptual example

Declarative validation with nested structs, property paths, and reusable `Validatable` types:

```go
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

// Usage:
err := validator.ValidateIt(context.Background(), Product{
    Name: "",
    Tags: []string{"a", "", "a"},
    Components: []Component{{Name: ""}},
})
// violations:
//   'name': This value should not be blank.
//   'tags': This collection should contain 2 elements or more.
//   'tags': This collection should contain only unique elements.
//   'tags[1]': This value should not be blank.
//   'components[0].name': This value should not be blank.
```

See [Usage](docs/usage.md) for validator setup and validation arguments.

## Documentation

| Topic | Description |
|-------|-------------|
| [Installation](docs/installation.md) | How to install the package |
| [Usage](docs/usage.md) | Basic concepts, validator, validation arguments |
| [Property paths & structs](docs/property-paths-and-structs.md) | Property paths, struct validation, conditional validation, groups |
| [Violations and errors](docs/violations-and-errors.md) | Handling violations, error structure, storing in database |
| [Translations](docs/translations.md) | Multi-language support and custom messages |
| [Custom constraints](docs/custom-constraints.md) | Creating your own constraints |

Full index: [docs/README.md](docs/README.md).

API reference: [pkg.go.dev/github.com/muonsoft/validation](https://pkg.go.dev/github.com/muonsoft/validation)

## Contributing

You may help this project by

* reporting an [issue](https://github.com/muonsoft/validation/issues);
* making translations for error messages;
* suggest an improvement or [discuss](https://github.com/muonsoft/validation/discussions) the usability of the package.

If you'd like to contribute, see [the contribution guide](CONTRIBUTING.md). Pull requests are welcome.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
