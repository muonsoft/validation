package scenario

import (
	"context"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
)

type Profile struct {
	DisplayName, Email string
	Age                int
	Address            Address
	Preferences        *Preferences
}
type Preferences struct {
	Theme string
	Quota *int
	Tags  []string
}
type Address struct{ Street, Region string }

func (x Address) Validate(ctx context.Context, v *validation.Validator, region validation.StringConstraint) error {
	return v.Validate(ctx, validation.StringProperty("street", x.Street, it.IsNotBlank()), validation.StringProperty("region", x.Region, it.IsNotBlank(), region))
}
