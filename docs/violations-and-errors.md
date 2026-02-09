# Working with Violations and Errors

There are two types of errors returned from the validator. One is validation violations and another is internal errors (
for example, when attempting to apply a constraint on not applicable argument type). The best way to handle validation
errors is to check for implementing the `validation.ViolationList` struct. You can use the default way to unwrap errors.

```go
err := validator.Validate(/* validation arguments */)

if err != nil {
    var violations *validation.ViolationList
    if errors.As(err, &violations) {
        // handle violations
    } else {
        // handle internal error
    }
}
```

Also, you can use helper function `validation.UnwrapViolations()`.

```go
err := validator.Validate(/* validation arguments */)
if violations, ok := validation.UnwrapViolations(err); ok {
    // handle violations
} else if err != nil {
    // handle internal error
}
```

The validation error called violation consists of a few parameters.

* `error` - underlying static error. This error can be used as a unique, short, and semantic code of violation.
  You can use it to test `Violation` for specific static error by `errors.Is` from standard library.
  Built-in error values are defined in the `github.com/muonsoft/validation/errors.go`.
  Error code values are protected by backward compatibility rules, template values are not protected.
* `message` - translated message with injected values from constraint. It can be used to show a description of a
  violation to the end-user. Possible values for built-in constraints are defined in
  the `github.com/muonsoft/validation/message` package and can be changed at any time, even in patch versions.
* `messageTemplate` - template for rendering message. Alongside `parameters` it can be used to render the message on the
  client-side of the library.
* `parameters` is the map of the template variables and their values provided by the specific constraint.
* `propertyPath` points to the violated property as it described in the [Property paths](property-paths-and-structs.md#processing-property-paths) section.

Thanks to the static error codes provided, you can quickly test the resulting validation error for a specific violation
error using standard `errors.Is()` function.

```go
err := validator.Validate(context.Background(), validation.String("", it.IsNotBlank()))

fmt.Println("is validation.ErrIsBlank =", errors.Is(err, validation.ErrIsBlank))
// Output:
// is validation.ErrIsBlank = true
```

You can hook into process of violation generation by implementing `validation.ViolationFactory` interface and passing it
via `validation.SetViolationFactory()` option. Custom violation must implement `validation.Violation` interface.

## Storing violations in a database

If you have a need to store violations in persistent storage (database), then it is recommended to store only error code,
property path, and template parameters. It is not recommended to store message templates because they can contain
mistakes and can be changed more frequently than violation error codes. The better practice is to store messages in
separate storage with translations and to load them by violation error codes. So make sure that violation errors codes
are unique and have only one specific message template. To restore the violations from a storage load an error code,
property path, template parameters, and find a message template by the violation error code. To make a violation
error code unique it is recommended to use a namespaced value, for example `app: product: empty tags`.
