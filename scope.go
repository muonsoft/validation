package validation

import (
	"context"

	languagepkg "github.com/muonsoft/language"
	"golang.org/x/text/language"
)

// Scope holds the current state of validation. On the client-side of the package,
// it can be used to build violations.
type Scope struct {
	context          context.Context
	propertyPath     PropertyPath
	language         language.Tag
	translator       *Translator
	violationFactory ViolationFactory
}

// Context returns context value that was passed to the validator by Context argument or
// by creating scoped validator with the validator.WithContext method.
func (s Scope) Context() context.Context {
	return s.context
}

// BuildViolation is used to create violations in validation methods of constraints.
// This method automatically injects the property path and language of the current validation scope.
func (s Scope) BuildViolation(code, message string) *ViolationBuilder {
	b := s.violationFactory.BuildViolation(code, message)
	b.SetPropertyPath(s.propertyPath)

	if s.language != language.Und {
		b.SetLanguage(s.language)
	} else if s.context != nil {
		b.SetLanguage(languagepkg.FromContext(s.context))
	}

	return b
}

func (s *Scope) applyOptions(options ...Option) error {
	for _, option := range options {
		err := option.SetUp(s)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s Scope) withContext(ctx context.Context) Scope {
	s.context = ctx

	return s
}

func (s Scope) withLanguage(tag language.Tag) Scope {
	s.language = tag

	return s
}

func (s Scope) atProperty(name string) Scope {
	s.propertyPath = append(s.propertyPath, PropertyNameElement(name))

	return s
}

func (s Scope) atIndex(index int) Scope {
	s.propertyPath = append(s.propertyPath, ArrayIndexElement(index))

	return s
}

func newScope() Scope {
	translator := newTranslator()

	return Scope{
		context:          context.Background(),
		translator:       translator,
		violationFactory: newViolationFactory(translator),
	}
}
