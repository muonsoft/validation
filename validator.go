package validation

import (
	"context"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message/catalog"
)

// Validator is the main validation service. It can be created by NewValidator constructor.
// Also, you can use singleton version from the package "github.com/muonsoft/validation/validator".
type Validator struct {
	scope Scope
}

// ValidatorOption is a base type for configuration options used to create a new instance of Validator.
type ValidatorOption func(validator *Validator) error

// NewValidator is a constructor for creating an instance of Validator.
// You can configure it by using validator options.
//
// Example
//  validator, err := validation.NewValidator(
//      validation.DefaultLanguage(language.Russian), // passing default language of translations
//      validation.Translations(russian.Messages), // setting up custom or built-in translations
//      validation.SetViolationFactory(userViolationFactory), // if you want to override creation of violations
//  )
//
//  // don't forget to check for errors
//  if err != nil {
//      fmt.Println(err)
//  }
func NewValidator(options ...ValidatorOption) (*Validator, error) {
	validator := &Validator{scope: newScope()}

	for _, setOption := range options {
		err := setOption(validator)
		if err != nil {
			return nil, err
		}
	}

	err := validator.scope.translator.init()
	if err != nil {
		return nil, err
	}

	return validator, nil
}

func newScopedValidator(scope Scope) *Validator {
	return &Validator{scope: scope}
}

// DefaultLanguage is used to set up the default language for translation of violation messages.
func DefaultLanguage(tag language.Tag) ValidatorOption {
	return func(validator *Validator) error {
		validator.scope.translator.defaultLanguage = tag
		return nil
	}
}

// Translations is used to load translation messages into the validator.
//
// By default, all violation messages are generated in the English language with pluralization capabilities.
// To use a custom language you have to load translations on validator initialization.
// Built-in translations are available in the sub-packages of the package "github.com/muonsoft/message/translations".
// The translation mechanism is provided by the "golang.org/x/text" package (be aware, it has no stable version yet).
//
// Example
//  // import "github.com/muonsoft/validation/message/translations/russian"
//
//  validator, err := validation.NewValidator(
//      validation.Translations(russian.Messages),
//  )
func Translations(messages map[language.Tag]map[string]catalog.Message) ValidatorOption {
	return func(validator *Validator) error {
		return validator.scope.translator.loadMessages(messages)
	}
}

// SetViolationFactory can be used to override the mechanism of violation creation.
func SetViolationFactory(factory ViolationFactory) ValidatorOption {
	return func(validator *Validator) error {
		validator.scope.violationFactory = factory

		return nil
	}
}

// StoredConstraint option can be used to store a constraint in an internal validator store.
// It can later be used by the validator.ValidateBy method. This can be useful for passing
// custom or prepared constraints to Validatable.
//
// If the constraint already exists, a ConstraintAlreadyStoredError will be returned.
//
// Example
//	validator, err := validation.NewValidator(
//		validation.StoredConstraint("isTagExists", isTagExistsConstraint)
//	)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	s := "
//	err = validator.ValidateString(&s, validator.ValidateBy("isTagExists"))
func StoredConstraint(key string, constraint Constraint) ValidatorOption {
	return func(validator *Validator) error {
		if _, exists := validator.scope.constraints[key]; exists {
			return ConstraintAlreadyStoredError{Key: key}
		}

		validator.scope.constraints[key] = constraint

		return nil
	}
}

// Validate is the main validation method. It accepts validation arguments. Arguments can be
// used to tune up the validation process or to pass values of a specific type.
func (validator *Validator) Validate(ctx context.Context, arguments ...Argument) error {
	args := &Arguments{scope: validator.scope.withContext(ctx)}
	for _, argument := range arguments {
		err := argument.set(args)
		if err != nil {
			return err
		}
	}

	violations := &ViolationList{}
	for _, validate := range args.validators {
		vs, err := validate(args.scope)
		if err != nil {
			return err
		}
		violations.Join(vs)
	}

	return violations.AsError()
}

// ValidateValue is an alias for validating a single value of any supported type.
func (validator *Validator) ValidateValue(ctx context.Context, value interface{}, options ...Option) error {
	return validator.Validate(ctx, Value(value, options...))
}

// ValidateBool is an alias for validating a single boolean value.
func (validator *Validator) ValidateBool(ctx context.Context, value *bool, options ...Option) error {
	return validator.Validate(ctx, Bool(value, options...))
}

// ValidateNumber is an alias for validating a single numeric value (integer or float).
func (validator *Validator) ValidateNumber(ctx context.Context, value interface{}, options ...Option) error {
	return validator.Validate(ctx, Number(value, options...))
}

// ValidateString is an alias for validating a single string value.
func (validator *Validator) ValidateString(ctx context.Context, value *string, options ...Option) error {
	return validator.Validate(ctx, String(value, options...))
}

// ValidateStrings is an alias for validating slice of strings.
func (validator *Validator) ValidateStrings(ctx context.Context, values []string, options ...Option) error {
	return validator.Validate(ctx, Strings(values, options...))
}

// ValidateIterable is an alias for validating a single iterable value (an array, slice, or map).
func (validator *Validator) ValidateIterable(ctx context.Context, value interface{}, options ...Option) error {
	return validator.Validate(ctx, Iterable(value, options...))
}

// ValidateCountable is an alias for validating a single countable value (an array, slice, or map).
func (validator *Validator) ValidateCountable(ctx context.Context, count int, options ...Option) error {
	return validator.Validate(ctx, Countable(count, options...))
}

// ValidateTime is an alias for validating a single time value.
func (validator *Validator) ValidateTime(ctx context.Context, value *time.Time, options ...Option) error {
	return validator.Validate(ctx, Time(value, options...))
}

// ValidateEach is an alias for validating each value of an iterable (an array, slice, or map).
func (validator *Validator) ValidateEach(ctx context.Context, value interface{}, options ...Option) error {
	return validator.Validate(ctx, Each(value, options...))
}

// ValidateEachString is an alias for validating each value of a strings slice.
func (validator *Validator) ValidateEachString(ctx context.Context, values []string, options ...Option) error {
	return validator.Validate(ctx, EachString(values, options...))
}

// ValidateValidatable is an alias for validating value that implements the Validatable interface.
func (validator *Validator) ValidateValidatable(ctx context.Context, validatable Validatable, options ...Option) error {
	return validator.Validate(ctx, Valid(validatable, options...))
}

// WithLanguage method creates a new scoped validator with a given language tag. All created violations
// will be translated into this language.
//
// Example
//  err := validator.WithLanguage(language.Russian).Validate(
//      validation.ValidateString(&s, it.IsNotBlank()), // violation from this constraint will be translated
//  )
func (validator *Validator) WithLanguage(tag language.Tag) *Validator {
	return newScopedValidator(validator.scope.withLanguage(tag))
}

// AtProperty method creates a new scoped validator with injected property name element to scope property path.
func (validator *Validator) AtProperty(name string) *Validator {
	return newScopedValidator(validator.scope.atProperty(name))
}

// AtIndex method creates a new scoped validator with injected array index element to scope property path.
func (validator *Validator) AtIndex(index int) *Validator {
	return newScopedValidator(validator.scope.atIndex(index))
}

// BuildViolation can be used to build a custom violation on the client-side.
func (validator *Validator) BuildViolation(code, message string) *ViolationBuilder {
	return validator.scope.BuildViolation(code, message)
}

// ValidateBy is used to get the constraint from the internal validator store.
// If the constraint does not exist, then the validator will
// return a ConstraintNotFoundError during the validation process.
// For storing a constraint you should use the StoreConstraint method.
func (validator *Validator) ValidateBy(constraintKey string) Constraint {
	if constraint, exists := validator.scope.constraints[constraintKey]; exists {
		return constraint
	}

	return notFoundConstraint{key: constraintKey}
}
