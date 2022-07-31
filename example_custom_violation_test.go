package validation_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/message/translations"
	"golang.org/x/text/language"
)

// DomainError is the container that will pass the ID to DomainViolation.
type DomainError struct {
	ID      string // this ID will be passed into DomainViolation by DomainViolationFactory
	Message string
}

func (err *DomainError) Error() string {
	return err.Message
}

var ErrIsEmpty = &DomainError{ID: "IsEmpty", Message: "Value is empty."}

// DomainViolation is custom implementation of validation.Violation interface with domain
// data and custom marshaling to JSON.
type DomainViolation struct {
	id string // id passed from DomainError

	// required fields for implementing validation.Violation
	err             error
	message         string
	messageTemplate string
	parameters      []validation.TemplateParameter
	propertyPath    *validation.PropertyPath
}

func (v *DomainViolation) Unwrap() error                              { return v.err }
func (v *DomainViolation) Is(target error) bool                       { return errors.Is(v.err, target) }
func (v *DomainViolation) Error() string                              { return v.err.Error() }
func (v *DomainViolation) Message() string                            { return v.message }
func (v *DomainViolation) MessageTemplate() string                    { return v.messageTemplate }
func (v *DomainViolation) Parameters() []validation.TemplateParameter { return v.parameters }
func (v *DomainViolation) PropertyPath() *validation.PropertyPath     { return v.propertyPath }

// pathAsJSONPointer formats property path according to a JSON Pointer Syntax https://tools.ietf.org/html/rfc6901
func (v *DomainViolation) pathAsJSONPointer() string {
	var s strings.Builder
	for _, element := range v.propertyPath.Elements() {
		s.WriteRune('/')
		s.WriteString(element.String())
	}
	return s.String()
}

// MarshalJSON marshals violation data with id, message and path fields. Path is formatted
// according to JSON Pointer Syntax.
func (v *DomainViolation) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ID           string `json:"id"`
		Message      string `json:"message"`
		PropertyPath string `json:"path"`
	}{
		ID:           v.id,
		Message:      v.message,
		PropertyPath: v.pathAsJSONPointer(),
	})
}

// DomainViolationFactory is custom implementation for validation.ViolationFactory.
type DomainViolationFactory struct {
	// reuse translations and templating from BuiltinViolationFactory
	factory *validation.BuiltinViolationFactory
}

func NewDomainViolationFactory() (*DomainViolationFactory, error) {
	translator, err := translations.NewTranslator()
	if err != nil {
		return nil, err
	}

	return &DomainViolationFactory{factory: validation.NewViolationFactory(translator)}, nil
}

func (factory *DomainViolationFactory) CreateViolation(err error, messageTemplate string, pluralCount int, parameters []validation.TemplateParameter, propertyPath *validation.PropertyPath, lang language.Tag) validation.Violation {
	// extracting error ID from err if it implements DomainError
	id := ""
	var domainErr *DomainError
	if errors.As(err, &domainErr) {
		id = domainErr.ID
	}

	violation := factory.factory.CreateViolation(err, messageTemplate, pluralCount, parameters, propertyPath, lang)

	return &DomainViolation{
		id:              id,
		err:             err,
		message:         violation.Message(),
		messageTemplate: violation.MessageTemplate(),
		parameters:      violation.Parameters(),
		propertyPath:    violation.PropertyPath(),
	}
}

func ExampleSetViolationFactory() {
	violationFactory, err := NewDomainViolationFactory()
	if err != nil {
		log.Fatalln(err)
	}
	validator, err := validation.NewValidator(validation.SetViolationFactory(violationFactory))
	if err != nil {
		log.Fatalln(err)
	}

	err = validator.At(
		// property path will be formatted according to JSON Pointer Syntax
		validation.PropertyName("properties"),
		validation.ArrayIndex(1),
		validation.PropertyName("key"),
	).ValidateString(
		context.Background(),
		"",
		// passing DomainError implementation via a constraint method
		it.IsNotBlank().WithError(ErrIsEmpty).WithMessage(ErrIsEmpty.Message),
	)

	marshaled, err := json.MarshalIndent(err, "", "\t")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(marshaled))
	// Output:
	// [
	//	{
	//		"id": "IsEmpty",
	//		"message": "Value is empty.",
	//		"path": "/properties/1/key"
	//	}
	// ]
}
