package validation

import (
	"encoding/json"
	"errors"
	"strings"

	"golang.org/x/text/language"
)

// Violation is the abstraction for validator errors. You can use your own implementations on the application
// side to use it for your needs. In order for the validator to generate application violations,
// it is necessary to implement the ViolationFactory interface and inject it into the validator.
// You can do this by using the SetViolationFactory option in the NewValidator constructor.
type Violation interface {
	error

	// Code is unique, short, and semantic string that can be used to programmatically
	// test for specific violation. All "code" values are defined in the "github.com/muonsoft/validation/code" package
	// and are protected by backward compatibility rules.
	Code() string

	// Is can be used to check that the violation contains one of the specific codes.
	// For an empty list, it should always returns false.
	Is(codes ...string) bool

	// Message is a translated message with injected values from constraint. It can be used to show
	// a description of a violation to the end-user. Possible values for build-in constraints
	// are defined in the "github.com/muonsoft/validation/message" package and can be changed at any time,
	// even in patch versions.
	Message() string

	// MessageTemplate is a template for rendering message. Alongside parameters it can be used to
	// render the message on the client-side of the library.
	MessageTemplate() string

	// Parameters is the map of the template variables and their values provided by the specific constraint.
	Parameters() []TemplateParameter

	// PropertyPath is a path that points to the violated property.
	// See PropertyPath type description for more info.
	PropertyPath() *PropertyPath
}

// ViolationFactory is the abstraction that can be used to create custom violations on the application side.
// Use the SetViolationFactory option on the NewValidator constructor to inject your own factory into the validator.
type ViolationFactory interface {
	// CreateViolation creates a new instance of Violation.
	CreateViolation(
		code,
		messageTemplate string,
		pluralCount int,
		parameters []TemplateParameter,
		propertyPath *PropertyPath,
		lang language.Tag,
	) Violation
}

// ViolationList is a slice of violations. It is the usual type of error that is returned from a validator.
type ViolationList []Violation

// NewViolationFunc is an adapter that allows you to use ordinary functions as a ViolationFactory.
type NewViolationFunc func(
	code,
	messageTemplate string,
	pluralCount int,
	parameters []TemplateParameter,
	propertyPath *PropertyPath,
	lang language.Tag,
) Violation

// CreateViolation creates a new instance of a Violation.
func (f NewViolationFunc) CreateViolation(
	code,
	messageTemplate string,
	pluralCount int,
	parameters []TemplateParameter,
	propertyPath *PropertyPath,
	lang language.Tag,
) Violation {
	return f(code, messageTemplate, pluralCount, parameters, propertyPath, lang)
}

// Error returns a formatted list of errors as a string.
func (violations ViolationList) Error() string {
	if len(violations) == 0 {
		return "the list of violations is empty, it looks like you forgot to use the AsError method somewhere"
	}

	var s strings.Builder
	s.Grow(32 * len(violations))

	for i, v := range violations {
		if i > 0 {
			s.WriteString("; ")
		}
		if iv, ok := v.(*internalViolation); ok {
			iv.writeToBuilder(&s)
		} else {
			s.WriteString(v.Error())
		}
	}

	return s.String()
}

// AppendFromError appends a single violation or a slice of violations into the end of a given slice.
// If an error does not implement the Violation or ViolationList interface, it will return an error itself.
// Otherwise nil will be returned.
//
// Example
//  violations := make(ViolationList, 0)
//  err := violations.AppendFromError(previousError)
//  if err != nil {
//      // this error is not a violation, processing must fail
//      return err
//  }
//  // violations contain appended violations from the previousError and can be processed further
func (violations *ViolationList) AppendFromError(err error) error {
	if violation, ok := UnwrapViolation(err); ok {
		*violations = append(*violations, violation)
	} else if violationList, ok := UnwrapViolationList(err); ok {
		*violations = append(*violations, violationList...)
	} else if err != nil {
		return err
	}

	return nil
}

// Has can be used to check that at least one of the violations contains one of the specific codes.
// For an empty list of codes, it should always returns false.
func (violations ViolationList) Has(codes ...string) bool {
	for _, violation := range violations {
		if violation.Is(codes...) {
			return true
		}
	}

	return false
}

// Filter returns a new list of violations with violations of given codes.
func (violations ViolationList) Filter(codes ...string) ViolationList {
	filtered := make(ViolationList, 0, len(violations))

	for _, violation := range violations {
		if violation.Is(codes...) {
			filtered = append(filtered, violation)
		}
	}

	return filtered
}

// AsError converts the list of violations to an error. This method correctly handles cases where
// the list of violations is empty. It returns nil on an empty list, indicating that the validation was successful.
func (violations ViolationList) AsError() error {
	if len(violations) == 0 {
		return nil
	}

	return violations
}

// IsViolation can be used to verify that the error implements the Violation interface.
func IsViolation(err error) bool {
	var violation Violation

	return errors.As(err, &violation)
}

// IsViolationList can be used to verify that the error implements the ViolationList.
func IsViolationList(err error) bool {
	var violations ViolationList

	return errors.As(err, &violations)
}

// UnwrapViolation is a short function to unwrap Violation from the error.
func UnwrapViolation(err error) (Violation, bool) {
	var violation Violation

	as := errors.As(err, &violation)

	return violation, as
}

// UnwrapViolationList is a short function to unwrap ViolationList from the error.
func UnwrapViolationList(err error) (ViolationList, bool) {
	var violations ViolationList

	as := errors.As(err, &violations)

	return violations, as
}

type internalViolation struct {
	code            string
	message         string
	messageTemplate string
	parameters      []TemplateParameter
	propertyPath    *PropertyPath
}

func (v internalViolation) Is(codes ...string) bool {
	for _, code := range codes {
		if v.code == code {
			return true
		}
	}

	return false
}

func (v internalViolation) Error() string {
	var s strings.Builder
	s.Grow(32)
	v.writeToBuilder(&s)

	return s.String()
}

func (v internalViolation) writeToBuilder(s *strings.Builder) {
	s.WriteString("violation")
	if v.propertyPath != nil {
		s.WriteString(" at '" + v.propertyPath.String() + "'")
	}
	s.WriteString(": " + v.message)
}

func (v internalViolation) Code() string {
	return v.code
}

func (v internalViolation) Message() string {
	return v.message
}

func (v internalViolation) MessageTemplate() string {
	return v.messageTemplate
}

func (v internalViolation) Parameters() []TemplateParameter {
	return v.parameters
}

func (v internalViolation) PropertyPath() *PropertyPath {
	return v.propertyPath
}

func (v internalViolation) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Code         string        `json:"code"`
		Message      string        `json:"message"`
		PropertyPath *PropertyPath `json:"propertyPath,omitempty"`
	}{
		Code:         v.code,
		Message:      v.message,
		PropertyPath: v.propertyPath,
	})
}

type internalViolationFactory struct {
	translator *Translator
}

func newViolationFactory(translator *Translator) *internalViolationFactory {
	return &internalViolationFactory{translator: translator}
}

func (factory *internalViolationFactory) CreateViolation(
	code,
	messageTemplate string,
	pluralCount int,
	parameters []TemplateParameter,
	propertyPath *PropertyPath,
	lang language.Tag,
) Violation {
	message := factory.translator.translate(lang, messageTemplate, pluralCount)

	for i := range parameters {
		if parameters[i].NeedsTranslation {
			parameters[i].Value = factory.translator.translate(lang, parameters[i].Value, 0)
		}
	}

	return &internalViolation{
		code:            code,
		message:         renderMessage(message, parameters),
		messageTemplate: messageTemplate,
		parameters:      parameters,
		propertyPath:    propertyPath,
	}
}

// ViolationBuilder used to build an instance of a Violation.
type ViolationBuilder struct {
	code            string
	messageTemplate string
	pluralCount     int
	parameters      []TemplateParameter
	propertyPath    *PropertyPath
	language        language.Tag

	violationFactory ViolationFactory
}

// NewViolationBuilder creates a new ViolationBuilder.
func NewViolationBuilder(factory ViolationFactory) *ViolationBuilder {
	return &ViolationBuilder{violationFactory: factory}
}

// BuildViolation creates a new ViolationBuilder for composing Violation object fluently.
func (b *ViolationBuilder) BuildViolation(code, message string) *ViolationBuilder {
	return &ViolationBuilder{
		code:             code,
		messageTemplate:  message,
		violationFactory: b.violationFactory,
	}
}

// SetParameters sets template parameters that can be injected into the violation message.
func (b *ViolationBuilder) SetParameters(parameters ...TemplateParameter) *ViolationBuilder {
	b.parameters = parameters

	return b
}

// AddParameter adds one parameter into a slice of parameters.
func (b *ViolationBuilder) AddParameter(name, value string) *ViolationBuilder {
	b.parameters = append(b.parameters, TemplateParameter{Key: name, Value: value})

	return b
}

// SetPropertyPath sets a property path of violated attribute.
func (b *ViolationBuilder) SetPropertyPath(path *PropertyPath) *ViolationBuilder {
	b.propertyPath = path

	return b
}

// SetPluralCount sets a plural number that will be used for message pluralization during translations.
func (b *ViolationBuilder) SetPluralCount(pluralCount int) *ViolationBuilder {
	b.pluralCount = pluralCount

	return b
}

// SetLanguage sets language that will be used to translate the violation message.
func (b *ViolationBuilder) SetLanguage(tag language.Tag) *ViolationBuilder {
	b.language = tag

	return b
}

// CreateViolation creates a new violation with given parameters and returns it.
// Violation is created by calling the CreateViolation method of the ViolationFactory.
func (b *ViolationBuilder) CreateViolation() Violation {
	return b.violationFactory.CreateViolation(
		b.code,
		b.messageTemplate,
		b.pluralCount,
		b.parameters,
		b.propertyPath,
		b.language,
	)
}
