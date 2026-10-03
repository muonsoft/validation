package scenario

import (
	"context"
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
)

func (x Input) Validate(ctx context.Context, v *validation.Validator) error {
	return v.Validate(ctx, validation.StringProperty("name", x.Name, it.IsNotBlank()), validation.StringProperty("first", x.First, it.IsNotBlank()), validation.StringProperty("second", x.Second, it.IsNotBlank()))
}
func (s Service) Handle(ctx context.Context, v *validation.Validator, x Input) error {
	if err := v.ValidateIt(ctx, x); err != nil {
		return err
	}
	check := validation.Func[string](func(ctx context.Context, v *validation.Validator, id string) error {
		ok, err := s.Repository.Exists(ctx, id)
		if err != nil {
			return err
		}
		if !ok {
			return v.CreateViolation(ctx, ErrMissing, ErrMissing.Message())
		}
		return nil
	})
	if err := v.Validate(ctx, validation.This(x.First, check).At(validation.PropertyName("first")), validation.This(x.Second, check).At(validation.PropertyName("second"))); err != nil {
		return err
	}
	return s.Repository.Save(ctx, x)
}
