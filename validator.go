package validation

import (
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message/catalog"
)

type Validator struct {
	scope      Scope
	translator *Translator
}

type ValidatorOption func(validator *Validator) error

func NewValidator(options ...ValidatorOption) (*Validator, error) {
	validator := &Validator{scope: newScope()}
	validator.translator = &Translator{}
	validator.scope.violationFactory = newViolationFactory(validator.translator)

	for _, setOption := range options {
		err := setOption(validator)
		if err != nil {
			return nil, err
		}
	}

	err := validator.translator.init()
	if err != nil {
		return nil, err
	}

	return validator, nil
}

func DefaultLanguage(tag language.Tag) ValidatorOption {
	return func(validator *Validator) error {
		validator.translator.defaultLanguage = tag
		return nil
	}
}

func Translations(messages map[language.Tag]map[string]catalog.Message) ValidatorOption {
	return func(validator *Validator) error {
		return validator.translator.loadMessages(messages)
	}
}

func OverrideViolationFactory(factory ViolationFactory) ValidatorOption {
	return func(validator *Validator) error {
		validator.scope.violationFactory = factory

		return nil
	}
}

func OverrideDefaults(options ...ValidatorOption) error {
	for _, setOption := range options {
		err := setOption(validator)
		if err != nil {
			return err
		}
	}

	return nil
}

func ResetDefaults() {
	validator.scope = newScope()
	validator.translator.defaultLanguage = language.English
	validator.scope.violationFactory = newViolationFactory(validator.translator)
}

var validator, _ = NewValidator()

type validateByConstraintFunc func(constraint Constraint, scope Scope) error

func (validator *Validator) Validate(arguments ...Argument) error {
	args := &Arguments{scope: validator.scope}
	for _, argument := range arguments {
		err := argument.set(args)
		if err != nil {
			return err
		}
	}

	violations := make(ViolationList, 0)
	for _, validate := range args.validators {
		vs, err := validate(args.scope)
		if err != nil {
			return err
		}
		violations = append(violations, vs...)
	}

	return violations.AsError()
}

func (validator *Validator) InScope(scope Scope) *Validator {
	return &Validator{
		scope:      scope,
		translator: validator.translator,
	}
}

func (validator *Validator) ValidateValue(value interface{}, options ...Option) error {
	return validator.Validate(Value(value, options...))
}

func (validator *Validator) ValidateBool(value *bool, options ...Option) error {
	return validator.Validate(Bool(value, options...))
}

func (validator *Validator) ValidateNumber(value interface{}, options ...Option) error {
	return validator.Validate(Number(value, options...))
}

func (validator *Validator) ValidateString(value *string, options ...Option) error {
	return validator.Validate(String(value, options...))
}

func (validator *Validator) ValidateIterable(value interface{}, options ...Option) error {
	return validator.Validate(Iterable(value, options...))
}

func (validator *Validator) ValidateCountable(count int, options ...Option) error {
	return validator.Validate(Countable(count, options...))
}

func (validator *Validator) ValidateTime(value *time.Time, options ...Option) error {
	return validator.Validate(Time(value, options...))
}

func (validator *Validator) ValidateValidatable(validatable Validatable, options ...Option) error {
	return validator.Validate(Valid(validatable, options...))
}

func (validator *Validator) ValidateEach(value interface{}, options ...Option) error {
	return validator.Validate(Each(value, options...))
}

func (validator *Validator) ValidateEachString(values []string, options ...Option) error {
	return validator.Validate(EachString(values, options...))
}

func (validator *Validator) GetScope() Scope {
	return validator.scope
}
