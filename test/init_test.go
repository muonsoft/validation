package test

import (
	"testing"

	"github.com/muonsoft/validation"
)

func newValidator(t *testing.T, options ...validation.ValidatorOption) *validation.Validator {
	t.Helper()
	v, err := validation.NewValidator(options...)
	if err != nil {
		t.Fatal("initialize validator:", err)
	}
	return v
}
