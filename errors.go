package validation

import (
	"errors"
	"fmt"
	"strings"
)

// ConstraintError is used to return critical error from constraint that immediately
// stops the validation process. It is recommended to use Scope.NewConstraintError() method
// to initiate an error from current scope.
type ConstraintError struct {
	ConstraintName string
	Path           *PropertyPath
	Description    string
}

func (err *ConstraintError) Error() string {
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

func (err *ConstraintAlreadyStoredError) Error() string {
	return fmt.Sprintf(`constraint with key "%s" already stored`, err.Key)
}

// ConstraintNotFoundError is returned when trying to get a constraint
// from the validator store using a non-existent key.
type ConstraintNotFoundError struct {
	Key  string
	Type string
}

func (err *ConstraintNotFoundError) Error() string {
	return fmt.Sprintf(`constraint by key "%s" of type "%s" is not found`, err.Key, err.Type)
}

var errTranslatorOptionsDenied = errors.New("translation options denied when using custom translator")
