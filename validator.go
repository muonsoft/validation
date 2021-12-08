package validation

import (
	"context"
	"fmt"
	"time"

	"github.com/muonsoft/validation/message/translations"
	"golang.org/x/text/language"
	"golang.org/x/text/message/catalog"
)

// Validator is the main validation service. It can be created by NewValidator constructor.
// Also, you can use singleton version from the package "github.com/muonsoft/validation/validator".
type Validator struct {
	scope Scope
}

// Translator is used to translate violation messages. By default, validator uses an implementation from
// "github.com/muonsoft/validation/message/translations" package based on "golang.org/x/text" package.
// You can set up your own implementation by using SetTranslator option.
type Translator interface {
	Translate(tag language.Tag, message string, pluralCount int) string
}

// ValidatorOptions is a temporary structure for collecting functional options ValidatorOption.
type ValidatorOptions struct {
	translatorOptions []translations.TranslatorOption
	translator        Translator
	violationFactory  ViolationFactory
	constraints       map[string]Constraint
}

func newValidatorOptions() *ValidatorOptions {
	return &ValidatorOptions{
		constraints: map[string]Constraint{},
	}
}

// ValidatorOption is a base type for configuration options used to create a new instance of Validator.
type ValidatorOption func(options *ValidatorOptions) error

// NewValidator is a constructor for creating an instance of Validator.
// You can configure it by using validator options.
func NewValidator(options ...ValidatorOption) (*Validator, error) {
	var err error

	opts := newValidatorOptions()
	for _, setOption := range options {
		err = setOption(opts)
		if err != nil {
			return nil, err
		}
	}

	if opts.translator != nil && len(opts.translatorOptions) > 0 {
		return nil, errTranslatorOptionsDenied
	}
	if opts.translator == nil {
		opts.translator, err = translations.NewTranslator(opts.translatorOptions...)
		if err != nil {
			return nil, fmt.Errorf("failed to set up default translator: %w", err)
		}
	}
	if opts.violationFactory == nil {
		opts.violationFactory = newViolationFactory(opts.translator)
	}

	validator := &Validator{scope: newScope(
		opts.translator,
		opts.violationFactory,
		opts.constraints,
	)}

	return validator, nil
}

func newScopedValidator(scope Scope) *Validator {
	return &Validator{scope: scope}
}

// DefaultLanguage option is used to set up the default language for translation of violation messages.
func DefaultLanguage(tag language.Tag) ValidatorOption {
	return func(options *ValidatorOptions) error {
		options.translatorOptions = append(options.translatorOptions, translations.DefaultLanguage(tag))

		return nil
	}
}

// Translations option is used to load translation messages into the validator.
//
// By default, all violation messages are generated in the English language with pluralization capabilities.
// To use a custom language you have to load translations on validator initialization.
// Built-in translations are available in the sub-packages of the package "github.com/muonsoft/message/translations".
// The translation mechanism is provided by the "golang.org/x/text" package (be aware, it has no stable version yet).
func Translations(messages map[language.Tag]map[string]catalog.Message) ValidatorOption {
	return func(options *ValidatorOptions) error {
		options.translatorOptions = append(options.translatorOptions, translations.SetTranslations(messages))

		return nil
	}
}

// SetTranslator option is used to set up the custom implementation of message violation translator.
func SetTranslator(translator Translator) ValidatorOption {
	return func(options *ValidatorOptions) error {
		options.translator = translator

		return nil
	}
}

// SetViolationFactory option can be used to override the mechanism of violation creation.
func SetViolationFactory(factory ViolationFactory) ValidatorOption {
	return func(options *ValidatorOptions) error {
		options.violationFactory = factory

		return nil
	}
}

// StoredConstraint option can be used to store a constraint in an internal validator store.
// It can later be used by the validator.ValidateBy method. This can be useful for passing
// custom or prepared constraints to Validatable.
//
// If the constraint already exists, a ConstraintAlreadyStoredError will be returned.
func StoredConstraint(key string, constraint Constraint) ValidatorOption {
	return func(options *ValidatorOptions) error {
		if _, exists := options.constraints[key]; exists {
			return ConstraintAlreadyStoredError{Key: key}
		}

		options.constraints[key] = constraint

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
func (validator *Validator) ValidateBool(ctx context.Context, value bool, options ...Option) error {
	return validator.Validate(ctx, Bool(value, options...))
}

// ValidateNumber is an alias for validating a single numeric value (integer or float).
func (validator *Validator) ValidateNumber(ctx context.Context, value interface{}, options ...Option) error {
	return validator.Validate(ctx, Number(value, options...))
}

// ValidateString is an alias for validating a single string value.
func (validator *Validator) ValidateString(ctx context.Context, value string, options ...Option) error {
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
func (validator *Validator) ValidateTime(ctx context.Context, value time.Time, options ...Option) error {
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

// ValidateBy is used to get the constraint from the internal validator store.
// If the constraint does not exist, then the validator will
// return a ConstraintNotFoundError during the validation process.
// For storing a constraint you should use the StoredConstraint option.
func (validator *Validator) ValidateBy(constraintKey string) Constraint {
	if constraint, exists := validator.scope.constraints[constraintKey]; exists {
		return constraint
	}

	return notFoundConstraint{key: constraintKey}
}

// WithLanguage method creates a new scoped validator with a given language tag. All created violations
// will be translated into this language.
func (validator *Validator) WithLanguage(tag language.Tag) *Validator {
	return newScopedValidator(validator.scope.withLanguage(tag))
}

// AtProperty method creates a new scoped validator with injected property name element to scope property path.
func (validator *Validator) AtProperty(name string) *Validator {
	return newScopedValidator(validator.scope.AtProperty(name))
}

// AtIndex method creates a new scoped validator with injected array index element to scope property path.
func (validator *Validator) AtIndex(index int) *Validator {
	return newScopedValidator(validator.scope.AtIndex(index))
}

// BuildViolation can be used to build a custom violation on the client-side.
func (validator *Validator) BuildViolation(ctx context.Context, code, message string) *ViolationBuilder {
	return validator.scope.withContext(ctx).BuildViolation(code, message)
}
