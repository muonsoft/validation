package test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidator_WithContext_WhenContextWithValue_ExpectContextShouldReturnValue(t *testing.T) {
	validator := newValidator(t)

	ctx := context.WithValue(context.Background(), defaultContextKey, "value")
	contextValidator := validator.WithContext(ctx)
	value := contextValidator.Context().Value(defaultContextKey)

	assert.Equal(t, "value", value)
}
