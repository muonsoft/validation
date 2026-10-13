package scenario

import (
	"context"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
)

func (x Field) Validate(ctx context.Context, v *validation.Validator, depth, maxDepth int) error {
	if depth > maxDepth {
		return v.CreateViolation(ctx, ErrDepth, ErrDepth.Message())
	}
	args := []validation.Argument{
		validation.StringProperty("name", x.Name, it.IsNotBlank()),
		validation.ComparableProperty("kind", x.Kind, it.IsOneOf("text", "group").WithoutBlank()),
		validation.StringProperty("value", x.Value, it.IsNotBlank()).When(x.Kind == "text"),
		validation.ValidProperty("children", fieldList(x.Children, depth+1, maxDepth)).When(x.Kind == "group"),
	}
	for key, value := range x.Attributes {
		args = append(args, validation.AtProperty("attributes", validation.StringProperty(key, value, it.IsNotBlank())))
	}
	return v.Validate(ctx, args...)
}
