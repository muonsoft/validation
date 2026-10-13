package scenario

import (
	"context"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
)

func (x Request) Validate(ctx context.Context, v *validation.Validator, lookup Lookup) error {
	list := validation.NewViolationList()
	for i, item := range x.Items {
		child := v.AtProperty("items").AtIndex(i)
		err := child.Validate(ctx, validation.StringProperty("id", item.ID, it.IsNotBlank().WhenGroups(validation.DefaultGroup, "submit")), validation.NumberProperty("quantity", item.Quantity, it.IsPositive[int]().WhenGroups(validation.DefaultGroup, "submit")))
		if err != nil {
			if fatal := list.AppendFromError(err); fatal != nil {
				return fatal
			}
			continue
		}
		found, err := lookup.Exists(ctx, item.ID)
		if err != nil {
			return err
		}
		if !found {
			list.Append(child.AtProperty("id").CreateViolation(ctx, ErrMissing, ErrMissing.Message()))
		}
	}
	return list.AsError()
}
