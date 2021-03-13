package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPropertyPath_String(t *testing.T) {
	path := PropertyPath{PropertyNameElement("array"), ArrayIndexElement(1), PropertyNameElement("property")}

	formatted := path.String()

	assert.Equal(t, "array[1].property", formatted)
}
