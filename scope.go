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
	constraints      []Constraint
	violationFactory ViolationFactory
}

func (s *Scope) Context() context.Context {
	return s.context
}

func (s *Scope) AddConstraint(constraint Constraint) {
	s.constraints = append(s.constraints, constraint)
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

func GetScope() Scope {
	return validator.GetScope()
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

func (s *Scope) applyNonConstraints(options ...Option) error {
	for _, option := range options {
		if _, isConstraint := option.(Constraint); isConstraint {
			continue
		}

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

func extendAndPassOptions(extendedScope *Scope, passedOptions ...Option) Option {
	return OptionFunc(func(scope *Scope) error {
		scope.context = extendedScope.context
		scope.propertyPath = append(scope.propertyPath, extendedScope.propertyPath...)
		scope.violationFactory = extendedScope.violationFactory

		return scope.applyNonConstraints(passedOptions...)
	})
}
