package validation

import "fmt"

type ErrInapplicableConstraint struct {
	Code string
	Type string
}

func (err *ErrInapplicableConstraint) Error() string {
	return fmt.Sprintf("constraint with code '%s' cannot be applied to %s", err.Code, err.Type)
}
