package validation

import (
	"context"
	"fmt"

	"github.com/muonsoft/language"
)

const DefaultGroup = "default"

// Scope holds the current state of validation. On the client-side of the package,
// it can be used to build violations.
type Scope struct {
	context          context.Context
	propertyPath     *PropertyPath
	language         language.Tag
	translator       Translator
	violationFactory ViolationFactory
	groups           []string
	constraints      map[string]Constraint // todo
}

// Context returns context value that was passed to the validator by Context argument or
// by creating scoped validator with the validator.WithContext method.
func (s Scope) Context() context.Context {
	return s.context
}

// NewConstraintError creates a new ConstraintError, which can be used to stop validation process
// if constraint is not properly configured.
func (s Scope) NewConstraintError(constraintName, description string) ConstraintError {
	return ConstraintError{
		ConstraintName: constraintName,
		Path:           s.propertyPath,
		Description:    description,
	}
}

// BuildViolation is used to create violations in validation methods of constraints.
// This method automatically injects the property path and language of the current validation scope.
func (s Scope) BuildViolation(code, message string) *ViolationBuilder {
	b := NewViolationBuilder(s.violationFactory).BuildViolation(code, message)
	b.SetPropertyPath(s.propertyPath)

	if s.language != language.Und {
		b.SetLanguage(s.language)
	} else if s.context != nil {
		b.SetLanguage(language.FromContext(s.context))
	}

	return b
}

// AtProperty returns a copy of the scope with property path appended by the given property name.
func (s Scope) AtProperty(name string) Scope {
	s.propertyPath = s.propertyPath.WithProperty(name)

	return s
}

// AtIndex returns a copy of the scope with property path appended by the given array index.
func (s Scope) AtIndex(index int) Scope {
	s.propertyPath = s.propertyPath.WithIndex(index)

	return s
}

// Validator creates a new validator for the given scope. This validator can be used to perform
// complex validation on a custom constraint using existing constraints.
func (s Scope) Validator() *Validator {
	return newScopedValidator(s)
}

// IsIgnored is the reverse condition for applying validation groups to the IsApplied method.
// It is recommended to use this method in every validation method of the constraint.
func (s Scope) IsIgnored(groups ...string) bool {
	return !s.IsApplied(groups...)
}

// IsApplied compares scope validation groups and constraint groups. If one of the scope groups intersects with
// the constraint groups, the validation scope should be applied (returns true).
// Empty groups are treated as DefaultGroup. To set validation groups use the validator.WithGroups() method.
func (s Scope) IsApplied(groups ...string) bool {
	if len(s.groups) == 0 {
		if len(groups) == 0 {
			return true
		}
		for _, g := range groups {
			if g == DefaultGroup {
				return true
			}
		}
	}

	for _, g1 := range s.groups {
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

func (s Scope) validate(arguments ...Argument) error {
	return s.Validator().Validate(s.context, arguments...)
}

func (s *Scope) applyOptions(options ...Option) error {
	for _, option := range options {
		err := option.SetUp(s)
		if err != nil {
			return fmt.Errorf(`failed to set up option: %w`, err)
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

func (s Scope) withGroups(groups ...string) Scope {
	s.groups = groups
	return s
}

func newScope(
	translator Translator,
	violationFactory ViolationFactory,
	constraints map[string]Constraint,
) Scope {
	return Scope{
		context:          context.Background(),
		translator:       translator,
		violationFactory: violationFactory,
		constraints:      constraints,
	}
}
