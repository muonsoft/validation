package scenario

import (
	"context"
	"github.com/muonsoft/validation"
)

func CheckReferences(ctx context.Context, v *validation.Validator, ids []string, lookup Lookup) error {
	if len(ids) == 0 {
		return nil
	}
	found, err := lookup.Find(ctx, ids)
	if err != nil {
		return err
	}
	list := v.AtProperty("references").BuildViolationList(ctx)
	for i, id := range ids {
		if !found[id] {
			list.AddViolation(ErrMissing, ErrMissing.Message(), validation.ArrayIndex(i))
		}
	}
	return list.Create().AsError()
}
