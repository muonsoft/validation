package validation_test

import (
	"testing"

	"github.com/muonsoft/validation"
	"github.com/stretchr/testify/assert"
)

func TestPropertyPath_String(t *testing.T) {
	var path *validation.PropertyPath
	path = path.WithProperty("array").WithIndex(1).WithProperty("property")

	formatted := path.String()

	assert.Equal(t, "array[1].property", formatted)
}

func TestPropertyPath_With(t *testing.T) {
	path := validation.NewPropertyPath(validation.PropertyName("top"), validation.ArrayIndex(0))

	path = path.With(
		validation.PropertyName("low"),
		validation.ArrayIndex(1),
		validation.PropertyName("property"),
	)

	assert.Equal(t, "top[0].low[1].property", path.String())
}
