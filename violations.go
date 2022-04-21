package validation

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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

// ViolationList is a linked list of violations. It is the usual type of error that is returned from a validator.
type ViolationList struct {
	len   int
	first *ViolationListElement
	last  *ViolationListElement
}

// ViolationListElement points to violation build by validator. It also implements
// Violation and can be used as a proxy to underlying violation.
type ViolationListElement struct {
	next      *ViolationListElement
	violation Violation
}

// NewViolationList creates a new ViolationList, that can be immediately populated with
// variadic arguments of violations.
func NewViolationList(violations ...Violation) *ViolationList {
	list := &ViolationList{}
	list.Append(violations...)

	return list
}

// Len returns length of the linked list.
func (list *ViolationList) Len() int {
	if list == nil {
		return 0
	}

	return list.len
}

// Each can be used to iterate over ViolationList by a callback function. If callback returns
// any error, then it will be returned as a result of Each function.
func (list *ViolationList) Each(f func(i int, violation Violation) error) error {
	i := 0
	for violation := list.First(); violation != nil; violation = violation.Next() {
		err := f(i, violation)
		if err != nil {
			return err
		}
		i++
	}

	return nil
}

// First returns the first element of the linked list.
func (list *ViolationList) First() *ViolationListElement {
	return list.first
}

// Last returns the last element of the linked list.
func (list *ViolationList) Last() *ViolationListElement {
	return list.last
}

// Append appends violations to the end of the linked list.
func (list *ViolationList) Append(violations ...Violation) {
	for i := range violations {
		element := &ViolationListElement{violation: violations[i]}
		if list.first == nil {
			list.first = element
			list.last = element
		} else {
			list.last.next = element
			list.last = element
		}
	}

	list.len += len(violations)
}

// Join is used to append the given violation list to the end of the current list.
func (list *ViolationList) Join(violations *ViolationList) {
	if violations == nil || violations.len == 0 {
		return
	}

	if list.first == nil {
		list.first = violations.first
		list.last = violations.last
	} else {
		list.last.next = violations.first
		list.last = violations.last
	}

	list.len += violations.len
}

// Error returns a formatted list of errors as a string.
func (list *ViolationList) Error() string {
	if list == nil || list.len == 0 {
		return "the list of violations is empty, it looks like you forgot to use the AsError method somewhere"
	}

	var s strings.Builder
	s.Grow(32 * list.len)

	i := 0
	for e := list.first; e != nil; e = e.next {
		v := e.violation
		if i > 0 {
			s.WriteString("; ")
		}
		if iv, ok := v.(*internalViolation); ok {
			iv.writeToBuilder(&s)
		} else {
			s.WriteString(v.Error())
		}
		i++
	}

	return s.String()
}

// AppendFromError appends a single violation or a slice of violations into the end of a given slice.
// If an error does not implement the Violation or ViolationList interface, it will return an error itself.
// Otherwise nil will be returned.
func (list *ViolationList) AppendFromError(err error) error {
	if violation, ok := UnwrapViolation(err); ok {
		list.Append(violation)
	} else if violationList, ok := UnwrapViolationList(err); ok {
		list.Join(violationList)
	} else if err != nil {
		return err
	}

	return nil
}

// Has can be used to check that at least one of the violations contains one of the specific codes.
// For an empty list of codes, it should always returns false.
func (list *ViolationList) Has(codes ...string) bool {
	for e := list.First(); e != nil; e = e.next {
		if e.violation.Is(codes...) {
			return true
		}
	}

	return false
}

// Filter returns a new list of violations with violations of given codes.
func (list *ViolationList) Filter(codes ...string) *ViolationList {
	filtered := &ViolationList{}

	for e := list.First(); e != nil; e = e.next {
		if e.violation.Is(codes...) {
			filtered.Append(e.violation)
		}
	}

	return filtered
}

// AsError converts the list of violations to an error. This method correctly handles cases where
// the list of violations is empty. It returns nil on an empty list, indicating that the validation was successful.
func (list *ViolationList) AsError() error {
	if list == nil || list.len == 0 {
		return nil
	}

	return list
}

// AsSlice converts underlying linked list into slice of Violation.
func (list *ViolationList) AsSlice() []Violation {
	violations := make([]Violation, list.len)

	i := 0
	for e := list.First(); e != nil; e = e.next {
		violations[i] = e.violation
		i++
	}

	return violations
}

// MarshalJSON marshals the linked list into JSON. Usually, you should use
// json.Marshal function for marshaling purposes.
func (list *ViolationList) MarshalJSON() ([]byte, error) {
	b := bytes.Buffer{}
	b.WriteRune('[')
	i := 0
	for e := list.First(); e != nil; e = e.next {
		data, err := json.Marshal(e.violation)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal violation at %d: %w", i, err)
		}
		b.Write(data)
		if e.next != nil {
			b.WriteRune(',')
		}
		i++
	}
	b.WriteRune(']')

	return b.Bytes(), nil
}

// Next returns next element of the linked list.
func (element *ViolationListElement) Next() *ViolationListElement {
	return element.next
}

// Violation returns underlying violation value.
func (element *ViolationListElement) Violation() Violation {
	return element.violation
}

func (element *ViolationListElement) Error() string {
	return element.violation.Error()
}

func (element *ViolationListElement) Code() string {
	return element.violation.Code()
}

func (element *ViolationListElement) Is(codes ...string) bool {
	return element.violation.Is(codes...)
}

func (element *ViolationListElement) Message() string {
	return element.violation.Message()
}

func (element *ViolationListElement) MessageTemplate() string {
	return element.violation.MessageTemplate()
}

func (element *ViolationListElement) Parameters() []TemplateParameter {
	return element.violation.Parameters()
}

func (element *ViolationListElement) PropertyPath() *PropertyPath {
	return element.violation.PropertyPath()
}

// IsViolation can be used to verify that the error implements the Violation interface.
func IsViolation(err error) bool {
	var violation Violation

	return errors.As(err, &violation)
}

// IsViolationList can be used to verify that the error implements the ViolationList.
func IsViolationList(err error) bool {
	var violations *ViolationList

	return errors.As(err, &violations)
}

// UnwrapViolation is a short function to unwrap Violation from the error.
func UnwrapViolation(err error) (Violation, bool) {
	var violation Violation

	as := errors.As(err, &violation)

	return violation, as
}

// UnwrapViolationList is a short function to unwrap ViolationList from the error.
func UnwrapViolationList(err error) (*ViolationList, bool) {
	var violations *ViolationList

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
	translator Translator
}

func newViolationFactory(translator Translator) *internalViolationFactory {
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
	message := factory.translator.Translate(lang, messageTemplate, pluralCount)

	for i := range parameters {
		if parameters[i].NeedsTranslation {
			parameters[i].Value = factory.translator.Translate(lang, parameters[i].Value, 0)
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
// Deprecated: use WithParameters instead.
func (b *ViolationBuilder) SetParameters(parameters ...TemplateParameter) *ViolationBuilder {
	b.parameters = parameters

	return b
}

// AddParameter adds one parameter into a slice of parameters.
// Deprecated: use WithParameter instead.
func (b *ViolationBuilder) AddParameter(name, value string) *ViolationBuilder {
	b.parameters = append(b.parameters, TemplateParameter{Key: name, Value: value})

	return b
}

// SetPropertyPath resets a base property path of violated attributes.
func (b *ViolationBuilder) SetPropertyPath(path *PropertyPath) *ViolationBuilder {
	b.propertyPath = path

	return b
}

// SetPluralCount sets a plural number that will be used for message pluralization during translations.
// Deprecated: use WithPluralCount instead.
func (b *ViolationBuilder) SetPluralCount(pluralCount int) *ViolationBuilder {
	b.pluralCount = pluralCount

	return b
}

// SetLanguage sets language that will be used to translate the violation message.
// Deprecated: use WithLanguage instead.
func (b *ViolationBuilder) SetLanguage(tag language.Tag) *ViolationBuilder {
	b.language = tag

	return b
}

// WithParameters sets template parameters that can be injected into the violation message.
func (b *ViolationBuilder) WithParameters(parameters ...TemplateParameter) *ViolationBuilder {
	b.parameters = parameters

	return b
}

// WithParameter adds one parameter into a slice of parameters.
func (b *ViolationBuilder) WithParameter(name, value string) *ViolationBuilder {
	b.parameters = append(b.parameters, TemplateParameter{Key: name, Value: value})

	return b
}

// At appends a property path of violated attribute.
func (b *ViolationBuilder) At(path ...PropertyPathElement) *ViolationBuilder {
	b.propertyPath = b.propertyPath.With(path...)

	return b
}

// AtProperty adds a property name to property path of violated attribute.
func (b *ViolationBuilder) AtProperty(propertyName string) *ViolationBuilder {
	b.propertyPath = b.propertyPath.WithProperty(propertyName)

	return b
}

// AtIndex adds an array index to property path of violated attribute.
func (b *ViolationBuilder) AtIndex(index int) *ViolationBuilder {
	b.propertyPath = b.propertyPath.WithIndex(index)

	return b
}

// WithPluralCount sets a plural number that will be used for message pluralization during translations.
func (b *ViolationBuilder) WithPluralCount(pluralCount int) *ViolationBuilder {
	b.pluralCount = pluralCount

	return b
}

// WithLanguage sets language that will be used to translate the violation message.
func (b *ViolationBuilder) WithLanguage(tag language.Tag) *ViolationBuilder {
	b.language = tag

	return b
}

// CreateViolation creates a new violation with given parameters and returns it.
// Violation is created by calling the CreateViolation method of the ViolationFactory.
// Deprecated: use Create instead.
func (b *ViolationBuilder) CreateViolation() Violation {
	return b.Create()
}

// Create creates a new violation with given parameters and returns it.
// Violation is created by calling the CreateViolation method of the ViolationFactory.
func (b *ViolationBuilder) Create() Violation {
	return b.violationFactory.CreateViolation(
		b.code,
		b.messageTemplate,
		b.pluralCount,
		b.parameters,
		b.propertyPath,
		b.language,
	)
}

// ViolationListBuilder is used to build a ViolationList by fluent interface.
type ViolationListBuilder struct {
	violations       *ViolationList
	violationFactory ViolationFactory

	propertyPath *PropertyPath
	language     language.Tag
}

// ViolationListElementBuilder is used to build Violation that will be added into ViolationList
// of the ViolationListBuilder.
type ViolationListElementBuilder struct {
	listBuilder *ViolationListBuilder

	code            string
	messageTemplate string
	pluralCount     int
	parameters      []TemplateParameter
	propertyPath    *PropertyPath
}

// NewViolationListBuilder creates a new ViolationListBuilder.
func NewViolationListBuilder(factory ViolationFactory) *ViolationListBuilder {
	return &ViolationListBuilder{violationFactory: factory, violations: NewViolationList()}
}

// BuildViolation initiates a builder for violation that will be added into ViolationList.
func (b *ViolationListBuilder) BuildViolation(code, message string) *ViolationListElementBuilder {
	return &ViolationListElementBuilder{
		listBuilder:     b,
		code:            code,
		messageTemplate: message,
		propertyPath:    b.propertyPath,
	}
}

// AddViolation can be used to quickly add a new violation using only code, message
// and optional property path elements.
func (b *ViolationListBuilder) AddViolation(code, message string, path ...PropertyPathElement) *ViolationListBuilder {
	return b.add(code, message, 0, nil, b.propertyPath.With(path...))
}

// SetPropertyPath resets a base property path of violated attributes.
func (b *ViolationListBuilder) SetPropertyPath(path *PropertyPath) *ViolationListBuilder {
	b.propertyPath = path

	return b
}

// At appends a property path of violated attribute.
func (b *ViolationListBuilder) At(path ...PropertyPathElement) *ViolationListBuilder {
	b.propertyPath = b.propertyPath.With(path...)

	return b
}

// AtProperty adds a property name to the base property path of violated attributes.
func (b *ViolationListBuilder) AtProperty(propertyName string) *ViolationListBuilder {
	b.propertyPath = b.propertyPath.WithProperty(propertyName)

	return b
}

// AtIndex adds an array index to the base property path of violated attributes.
func (b *ViolationListBuilder) AtIndex(index int) *ViolationListBuilder {
	b.propertyPath = b.propertyPath.WithIndex(index)

	return b
}

// Create returns a ViolationList with built violations.
func (b *ViolationListBuilder) Create() *ViolationList {
	return b.violations
}

func (b *ViolationListBuilder) add(
	code, template string,
	count int,
	parameters []TemplateParameter,
	path *PropertyPath,
) *ViolationListBuilder {
	b.violations.Append(b.violationFactory.CreateViolation(
		code,
		template,
		count,
		parameters,
		path,
		b.language,
	))

	return b
}

// WithLanguage sets language that will be used to translate the violation message.
func (b *ViolationListBuilder) WithLanguage(tag language.Tag) *ViolationListBuilder {
	b.language = tag

	return b
}

// WithParameters sets template parameters that can be injected into the violation message.
func (b *ViolationListElementBuilder) WithParameters(parameters ...TemplateParameter) *ViolationListElementBuilder {
	b.parameters = parameters

	return b
}

// WithParameter adds one parameter into a slice of parameters.
func (b *ViolationListElementBuilder) WithParameter(name, value string) *ViolationListElementBuilder {
	b.parameters = append(b.parameters, TemplateParameter{Key: name, Value: value})

	return b
}

// At appends a property path of violated attribute.
func (b *ViolationListElementBuilder) At(path ...PropertyPathElement) *ViolationListElementBuilder {
	b.propertyPath = b.propertyPath.With(path...)

	return b
}

// AtProperty adds a property name to property path of violated attribute.
func (b *ViolationListElementBuilder) AtProperty(propertyName string) *ViolationListElementBuilder {
	b.propertyPath = b.propertyPath.WithProperty(propertyName)

	return b
}

// AtIndex adds an array index to property path of violated attribute.
func (b *ViolationListElementBuilder) AtIndex(index int) *ViolationListElementBuilder {
	b.propertyPath = b.propertyPath.WithIndex(index)

	return b
}

// WithPluralCount sets a plural number that will be used for message pluralization during translations.
func (b *ViolationListElementBuilder) WithPluralCount(pluralCount int) *ViolationListElementBuilder {
	b.pluralCount = pluralCount

	return b
}

// Add creates a Violation and appends it into the end of the ViolationList.
// It returns a ViolationListBuilder to continue process of creating a ViolationList.
func (b *ViolationListElementBuilder) Add() *ViolationListBuilder {
	return b.listBuilder.add(b.code, b.messageTemplate, b.pluralCount, b.parameters, b.propertyPath)
}

func unwrapViolationList(err error) (*ViolationList, error) {
	violations := NewViolationList()
	fatal := violations.AppendFromError(err)
	if fatal != nil {
		return nil, fatal
	}

	return violations, nil
}
