package validation

import (
	"context"
	"fmt"
	"time"

	"github.com/muonsoft/language"
	"github.com/muonsoft/validation/message/translations"
	"golang.org/x/text/message/catalog"
)

// Validator is the root validation service. It can be created by [NewValidator] constructor.
// Also, you can use singleton version from the package [github.com/muonsoft/validation/validator].
type Validator struct {
	propertyPath     *PropertyPath
	language         language.Tag
	translator       Translator
	violationFactory ViolationFactory
	groups           []string
}

// Translator is used to translate violation messages. By default, validator uses an implementation from
// [github.com/muonsoft/validation/message/translations] package based on [golang.org/x/text] package.
// You can set up your own implementation by using [SetTranslator] option.
type Translator interface {
	Translate(tag language.Tag, message string, pluralCount int) string
}

// ValidatorOptions is a temporary structure for collecting functional options [ValidatorOption].
type ValidatorOptions struct {
	translatorOptions []translations.TranslatorOption
	translator        Translator
	violationFactory  ViolationFactory
}

func newValidatorOptions() *ValidatorOptions {
	return &ValidatorOptions{}
}

// ValidatorOption is a base type for configuration options used to create a new instance of [Validator].
type ValidatorOption func(options *ValidatorOptions) error

// NewValidator is a constructor for creating an instance of [Validator].
// You can configure it by using the [ValidatorOption].
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
			return nil, fmt.Errorf("set up default translator: %w", err)
		}
	}
	if opts.violationFactory == nil {
		opts.violationFactory = NewViolationFactory(opts.translator)
	}

	validator := &Validator{
		translator:       opts.translator,
		violationFactory: opts.violationFactory,
	}

	return validator, nil
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
// Built-in translations are available in the sub-packages of the package [github.com/muonsoft/message/translations].
// The translation mechanism is provided by the [golang.org/x/text] package (be aware, it has no stable version yet).
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

// Validate is the main validation method. It accepts validation arguments that can be
// used to tune up the validation process or to pass values of a specific type.
func (validator *Validator) Validate(ctx context.Context, arguments ...Argument) error {
	execContext := &executionContext{}
	for _, argument := range arguments {
		argument.setUp(execContext)
	}

	violations := &ViolationList{}
	for _, validate := range execContext.validations {
		vs, err := validate(ctx, validator)
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

// ValidateIt is an alias for validating value that implements the [Validatable] interface.
func (validator *Validator) ValidateIt(ctx context.Context, validatable Validatable) error {
	return validator.Validate(ctx, Valid(validatable))
}

// WithGroups is used to execute conditional validation based on validation groups. It creates
// a new context validator with a given set of groups.
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
// Be careful, empty groups are considered as the default group. Its value is equal to the [DefaultGroup] ("default").
func (validator *Validator) WithGroups(groups ...string) *Validator {
	v := validator.copy()
	v.groups = groups

	return v
}

// IsAppliedForGroups compares current validation groups and constraint groups. If one of the validator groups
// intersects with the constraint groups, the validation process should be applied (returns true).
// Empty groups are treated as [DefaultGroup]. To create a new validator with the validation groups
// use the [Validator.WithGroups] method.
func (validator *Validator) IsAppliedForGroups(groups ...string) bool {
	if len(validator.groups) == 0 {
		if len(groups) == 0 {
			return true
		}
		for _, g := range groups {
			if g == DefaultGroup {
				return true
			}
		}
	}

	for _, g1 := range validator.groups {
		if len(groups) == 0 {
			if g1 == DefaultGroup {
				return true
			}
		}
		for _, g2 := range groups {
			if g1 == g2 {
				return true
			}
		}
	}

	return false
}

// IsIgnoredForGroups is the reverse condition for applying validation groups
// to the [Validator.IsAppliedForGroups] method. It is recommended to use this method in
// every validation method of the constraint.
func (validator *Validator) IsIgnoredForGroups(groups ...string) bool {
	return !validator.IsAppliedForGroups(groups...)
}

// CreateConstraintError creates a new [ConstraintError], which can be used to stop validation process
// if constraint is not properly configured.
func (validator *Validator) CreateConstraintError(constraintName, description string) *ConstraintError {
	return &ConstraintError{
		ConstraintName: constraintName,
		Path:           validator.propertyPath,
		Description:    description,
	}
}

// WithLanguage method creates a new context validator with a given language tag. All created violations
// will be translated into this language.
//
// The priority of language selection methods:
//
//   - [Validator.WithLanguage] has the highest priority and will override any other options;
//   - if the validator language is not specified, the validator will try to get the language from the context;
//   - in all other cases, the default language specified in the translator will be used.
func (validator *Validator) WithLanguage(tag language.Tag) *Validator {
	v := validator.copy()
	v.language = tag

	return v
}

// At method creates a new context validator with appended property path.
func (validator *Validator) At(path ...PropertyPathElement) *Validator {
	v := validator.copy()
	v.propertyPath = v.propertyPath.With(path...)

	return v
}

// AtProperty method creates a new context validator with appended property name to the property path.
func (validator *Validator) AtProperty(name string) *Validator {
	v := validator.copy()
	v.propertyPath = v.propertyPath.WithProperty(name)

	return v
}

// AtIndex method creates a new context validator with appended array index to the property path.
func (validator *Validator) AtIndex(index int) *Validator {
	v := validator.copy()
	v.propertyPath = v.propertyPath.WithIndex(index)

	return v
}

// CreateViolation can be used to quickly create a custom violation on the client-side.
func (validator *Validator) CreateViolation(ctx context.Context, err error, message string, path ...PropertyPathElement) Violation {
	return validator.BuildViolation(ctx, err, message).At(path...).Create()
}

// BuildViolation can be used to build a custom violation on the client-side.
func (validator *Validator) BuildViolation(ctx context.Context, err error, message string) *ViolationBuilder {
	b := NewViolationBuilder(validator.violationFactory).BuildViolation(err, message)
	b = b.SetPropertyPath(validator.propertyPath)

	if validator.language != language.Und {
		b = b.WithLanguage(validator.language)
	} else if ctx != nil {
		b = b.WithLanguage(language.FromContext(ctx))
	}

	return b
}

// BuildViolationList can be used to build a custom violation list on the client-side.
func (validator *Validator) BuildViolationList(ctx context.Context) *ViolationListBuilder {
	b := NewViolationListBuilder(validator.violationFactory)
	b = b.SetPropertyPath(validator.propertyPath)

	if validator.language != language.Und {
		b = b.WithLanguage(validator.language)
	} else if ctx != nil {
		b = b.WithLanguage(language.FromContext(ctx))
	}

	return b
}

func (validator *Validator) copy() *Validator {
	return &Validator{
		propertyPath:     validator.propertyPath,
		language:         validator.language,
		translator:       validator.translator,
		violationFactory: validator.violationFactory,
		groups:           validator.groups,
	}
}
