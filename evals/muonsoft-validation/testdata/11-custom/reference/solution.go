package scenario

import (
	"context"
	"github.com/muonsoft/validation"
	"golang.org/x/text/language"
	"golang.org/x/text/message/catalog"
	"strings"
)

func (c PrefixConstraint) ValidateString(ctx context.Context, v *validation.Validator, value *string) error {
	if value == nil || *value == "" || strings.HasPrefix(*value, c.Prefix) {
		return nil
	}
	return v.BuildViolation(ctx, ErrPrefix, ErrPrefix.Message()).WithParameter("{{ prefix }}", c.Prefix).Create()
}
func NewValidator() (*validation.Validator, error) {
	return validation.NewValidator(validation.Translations(map[language.Tag]map[string]catalog.Message{language.Russian: {ErrPrefix.Message(): catalog.String("Значение должно начинаться с {{ prefix }}.")}}))
}
