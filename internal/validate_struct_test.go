package internal

import (
	"testing"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/code"
	"github.com/muonsoft/validation/validationtest"
	"github.com/stretchr/testify/assert"
)

func TestValidate_WhenStructWithComplexRules_ExpectViolations(t *testing.T) {
	p := Product{
		Name: "",
		Components: []Component{
			{
				ID:   1,
				Name: "",
			},
		},
	}

	err := validation.ValidateValue(p)

	validationtest.AssertIsViolationList(t, err, func(t *testing.T, violations validation.ViolationList) bool {
		t.Helper()
		if assert.Len(t, violations, 4) {
			assert.Equal(t, code.NotBlank, violations[0].GetCode())
			assert.Equal(t, "name", violations[0].GetPropertyPath().Format())
			assert.Equal(t, code.CountTooFew, violations[1].GetCode())
			assert.Equal(t, "tags", violations[1].GetPropertyPath().Format())
			assert.Equal(t, code.NotBlank, violations[2].GetCode())
			assert.Equal(t, "components[0].name", violations[2].GetPropertyPath().Format())
			assert.Equal(t, code.CountTooFew, violations[3].GetCode())
			assert.Equal(t, "components[0].tags", violations[3].GetPropertyPath().Format())
		}
		return true
	})
}
