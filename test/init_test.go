package test

import (
	"testing"

	"github.com/muonsoft/validation"
)

func newValidator(tb testing.TB, options ...validation.ValidatorOption) *validation.Validator {
	tb.Helper()
	v, err := validation.NewValidator(options...)
	if err != nil {
		tb.Fatal("initialize validator:", err)
	}
	return v
}
