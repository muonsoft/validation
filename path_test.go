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
