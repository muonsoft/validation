package validation

import (
	"context"

	languagepkg "github.com/muonsoft/language"
	"golang.org/x/text/language"
)

type Scope struct {
	context          context.Context
	propertyPath     PropertyPath
	language         language.Tag
	violationFactory ViolationFactory
}

func (s *Scope) Context() context.Context {
	return s.context
}

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
		err := option.Set(s)
		if err != nil {
			return err
		}
	}

	return nil
}

func newScope() Scope {
	return Scope{context: context.Background()}
}
