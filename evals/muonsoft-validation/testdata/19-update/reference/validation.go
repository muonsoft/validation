package scenario

import (
	"context"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
)

func (x Bundle) Validate(ctx context.Context, v *validation.Validator) error {
	return v.Validate(ctx,
		validation.StringProperty("name", x.Name, it.IsNotBlank(), it.HasMaxLength(30)),
		validation.AtProperty("references", validation.Countable(len(x.References), it.HasCountBetween(1, 4)),
			validation.Each(x.References, validation.Func[Reference](func(ctx context.Context, child *validation.Validator, r Reference) error {
				return child.Validate(ctx, validation.StringProperty("id", r.ID, it.IsNotBlank()), validation.NumberProperty("weight", r.Weight, it.IsPositive[int]()))
			}))))
}
func CheckReferences(ctx context.Context, v *validation.Validator, r Repository, refs []Reference) error {
	if len(refs) == 0 {
		return nil
	}
	ids := make([]string, len(refs))
	for i, ref := range refs {
		ids[i] = ref.ID
	}
	found, err := r.Find(ctx, ids)
	if err != nil {
		return err
	}
	violations := v.BuildViolationList(ctx)
	for i, ref := range refs {
		if !found[ref.ID] {
			violations.BuildViolation(ErrMissing, ErrMissing.Message()).AtIndex(i).AtProperty("id").Add()
		}
	}
	return violations.Create().AsError()
}
