package validation

import (
	"errors"
	"fmt"
	"strings"
)

type ConstraintError struct {
	ConstraintName string
	Path           *PropertyPath
	Description    string
}

func (err ConstraintError) Error() string {
	var s strings.Builder
	s.WriteString("failed to validate by " + err.ConstraintName)
	if err.Path != nil {
		s.WriteString(` at path "` + err.Path.String() + `"`)
	}
	s.WriteString(": " + err.Description)

	return s.String()
}

// ConstraintAlreadyStoredError is returned when trying to put a constraint
// in the validator store using an existing key.
type ConstraintAlreadyStoredError struct {
	Key string
}

func (err ConstraintAlreadyStoredError) Error() string {
	return fmt.Sprintf(`constraint with key "%s" already stored`, err.Key)
}

// ConstraintNotFoundError is returned when trying to get a constraint
// from the validator store using a non-existent key.
type ConstraintNotFoundError struct {
	Key string
}

func (err ConstraintNotFoundError) Error() string {
	return fmt.Sprintf(`constraint with key "%s" is not stored in the validator`, err.Key)
}

var (
	errTranslatorOptionsDenied       = errors.New("translation options denied when using custom translator")
	errThenBranchNotSet              = errors.New("then branch of conditional constraint not set")
	errSequentiallyConstraintsNotSet = errors.New("constraints for sequentially validation not set")
	errAtLeastOneOfConstraintsNotSet = errors.New("constraints for at least one of validation not set")
	errCompoundConstraintsNotSet     = errors.New("constraints for compound validation not set")
)
