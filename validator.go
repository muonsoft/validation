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
	constraints       map[string]interface{}
}

func newValidatorOptions() *ValidatorOptions {
	return &ValidatorOptions{constraints: map[string]interface{}{}}
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
func StoredConstraint(key string, constraint interface{}) ValidatorOption {
	return func(options *ValidatorOptions) error {
		if _, exists := options.constraints[key]; exists {
			return &ConstraintAlreadyStoredError{Key: key}
		}

		options.constraints[key] = constraint

		return nil
	}
}

// Validate is the main validation method. It accepts validation arguments. executionContext can be
// used to tune up the validation process or to pass values of a specific type.
func (validator *Validator) Validate(ctx context.Context, arguments ...Argument) error {
	execContext := &executionContext{scope: validator.scope.withContext(ctx)}
	for _, argument := range arguments {
		argument.setUp(execContext)
	}

	violations := &ViolationList{}
	for _, validate := range execContext.validators {
		vs, err := validate(execContext.scope)
		if err != nil {
			return err
		}
		violations.Join(vs)
	}

	return violations.AsError()
}

// ValidateBool is an alias for validating a single boolean value.
func (validator *Validator) ValidateBool(ctx context.Context, value bool, constraints ...BoolConstraint) error {
	return validator.Validate(ctx, Bool(value, constraints...))
}

// ValidateInt is an alias for validating a single integer value.
func (validator *Validator) ValidateInt(ctx context.Context, value int, constraints ...NumberConstraint[int]) error {
	return validator.Validate(ctx, Number(value, constraints...))
}

// ValidateFloat is an alias for validating a single float value.
func (validator *Validator) ValidateFloat(ctx context.Context, value float64, constraints ...NumberConstraint[float64]) error {
	return validator.Validate(ctx, Number(value, constraints...))
}

// ValidateString is an alias for validating a single string value.
func (validator *Validator) ValidateString(ctx context.Context, value string, constraints ...StringConstraint) error {
	return validator.Validate(ctx, String(value, constraints...))
}

// ValidateStrings is an alias for validating slice of strings.
func (validator *Validator) ValidateStrings(ctx context.Context, values []string, constraints ...ComparablesConstraint[string]) error {
	return validator.Validate(ctx, Comparables(values, constraints...))
}

// ValidateCountable is an alias for validating a single countable value (an array, slice, or map).
func (validator *Validator) ValidateCountable(ctx context.Context, count int, constraints ...CountableConstraint) error {
	return validator.Validate(ctx, Countable(count, constraints...))
}

// ValidateTime is an alias for validating a single time value.
func (validator *Validator) ValidateTime(ctx context.Context, value time.Time, constraints ...TimeConstraint) error {
	return validator.Validate(ctx, Time(value, constraints...))
}

// ValidateEachString is an alias for validating each value of a strings slice.
func (validator *Validator) ValidateEachString(ctx context.Context, values []string, constraints ...StringConstraint) error {
	return validator.Validate(ctx, EachString(values, constraints...))
}

// ValidateIt is an alias for validating value that implements the Validatable interface.
func (validator *Validator) ValidateIt(ctx context.Context, validatable Validatable) error {
	return validator.Validate(ctx, Valid(validatable))
}

// GetConstraint is used to get the constraint from the internal validator store.
// If the constraint does not exist, then the validator will return nil.
// For storing a constraint you should use the StoredConstraint option.
//
// Experimental. This feature is experimental and may be changed in future versions.
func (validator *Validator) GetConstraint(key string) interface{} {
	return validator.scope.constraints[key]
}

// WithGroups is used to execute conditional validation based on validation groups. It creates
// a new scoped validation with a given set of groups.
//
// By default, when validating an object all constraints of it will be checked whether or not
// they pass. In some cases, however, you will need to validate an object against
// only some specific group of constraints. To do this, you can organize each constraint
// into one or more validation groups and then apply validation against one group of constraints.
//
// Validation groups are working together only with validation groups passed
// to a constraint by WhenGroups() method. This method is implemented in all built-in constraints.
// If you want to use validation groups for your own constraints do not forget to implement
// this method in your constraint.
//
// Be careful, empty groups are considered as the default group. Its value is equal to the DefaultGroup ("default").
func (validator *Validator) WithGroups(groups ...string) *Validator {
	return newScopedValidator(validator.scope.withGroups(groups...))
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

// CreateViolation can be used to quickly create a custom violation on the client-side.
func (validator *Validator) CreateViolation(ctx context.Context, code, message string, path ...PropertyPathElement) Violation {
	return validator.scope.withContext(ctx).CreateViolation(code, message, path...)
}

// BuildViolation can be used to build a custom violation on the client-side.
func (validator *Validator) BuildViolation(ctx context.Context, code, message string) *ViolationBuilder {
	return validator.scope.withContext(ctx).BuildViolation(code, message)
}

// BuildViolationList can be used to build a custom violation list on the client-side.
func (validator *Validator) BuildViolationList(ctx context.Context) *ViolationListBuilder {
	return validator.scope.withContext(ctx).BuildViolationList()
}
