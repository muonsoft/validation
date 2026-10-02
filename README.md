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

### Word count

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

### ISO weeks

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

## Documentation

| Topic | Description |
|-------|-------------|
| [Installation](docs/installation.md) | How to install the package |
| [Usage](docs/usage.md) | Basic concepts, validator, validation arguments |
| [Property paths & structs](docs/property-paths-and-structs.md) | Property paths, struct validation, conditional validation, groups |
| [Violations and errors](docs/violations-and-errors.md) | Handling violations, error structure, storing in database |
| [Country, language & locale](docs/international.md) | Formats, special codes, and application-support limitations |
| [UTF-8 validation](docs/utf8.md) | Reject malformed text at input boundaries; scope and limitations |
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
