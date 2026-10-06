package scenario

import (
	"context"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
)

func (x Catalog) Validate(ctx context.Context, v *validation.Validator, maxDepth int) error {
	return validateNodes(ctx, v.AtProperty("roots"), x.Roots, 0, maxDepth)
}
func validateNodes(ctx context.Context, v *validation.Validator, nodes map[string]*Node, depth, maxDepth int) error {
	list := validation.NewViolationList()
	for key, node := range nodes {
		child := v.AtProperty(key)
		var err error
		switch {
		case node == nil:
			err = child.CreateViolation(ctx, ErrMissing, ErrMissing.Message())
		case depth > maxDepth:
			err = child.CreateViolation(ctx, ErrDepth, ErrDepth.Message())
		default:
			err = child.Validate(ctx, validation.StringProperty("label", node.Label, it.IsNotBlank()), validation.ValidProperty("children", validation.ValidatableFunc(func(ctx context.Context, v *validation.Validator) error {
				return validateNodes(ctx, v, node.Children, depth+1, maxDepth)
			})))
		}
		if fatal := list.AppendFromError(err); fatal != nil {
			return fatal
		}
	}
	return list.AsError()
}
