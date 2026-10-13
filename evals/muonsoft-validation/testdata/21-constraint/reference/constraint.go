package scenario

import (
	"context"

	"github.com/muonsoft/validation"
	"golang.org/x/text/language"
	"golang.org/x/text/message/catalog"
)

type CodeConstraint struct {
	catalog  Catalog
	category string
}

func NewCodeConstraint(source Catalog, category string) CodeConstraint {
	return CodeConstraint{catalog: source, category: category}
}
func (c CodeConstraint) ValidateString(ctx context.Context, v *validation.Validator, value *string) error {
	if value == nil || *value == "" {
		return nil
	}
	entry, found, err := c.catalog.Lookup(ctx, *value)
	if err != nil {
		return err
	}
	if !found {
		return v.CreateViolation(ctx, ErrUnknown, ErrUnknown.Message())
	}
	if entry.Category != c.category {
		return v.BuildViolation(ctx, ErrCategory, ErrCategory.Message()).WithParameter("{{ category }}", c.category).Create()
	}
	return nil
}
func ValidatorOptions() []validation.ValidatorOption {
	return []validation.ValidatorOption{validation.Translations(map[language.Tag]map[string]catalog.Message{language.Russian: {ErrUnknown.Message(): catalog.String("Код не найден."), ErrCategory.Message(): catalog.String("Код должен относиться к категории {{ category }}.")}})}
}
