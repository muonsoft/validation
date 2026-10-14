package scenario

import (
	"context"

	"github.com/muonsoft/validation"
)

var ErrTenantEmail = validation.NewError("tenant_email", "Email is not allowed for this tenant.")

type Signup struct{ Name, Email string }

// EmailPolicy is supplied by the application's tenant policy layer.
type EmailPolicy struct {
	Check func(context.Context, *validation.Validator, string) error
}
