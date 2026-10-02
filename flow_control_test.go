package validation_test

import (
	"bytes"
	"context"
	"errors"
	"runtime/pprof"
	"strings"
	"testing"
	"time"

	"github.com/muonsoft/validation"
	"github.com/stretchr/testify/require"
)

func TestAsync_FatalErrorDoesNotLeakWorkers(t *testing.T) {
	validator := newValidator(t)
	fatal := errors.New("backend failure")
	for range 10 {
		err := validator.Validate(context.Background(), validation.Async(
			validation.CheckNoViolations(fatal),
			validation.Valid(validation.ValidatableFunc(func(ctx context.Context, _ *validation.Validator) error {
				<-ctx.Done()
				return ctx.Err()
			})),
		))
		require.ErrorIs(t, err, fatal)
	}
	require.Eventually(t, func() bool {
		var stacks bytes.Buffer
		if err := pprof.Lookup("goroutine").WriteTo(&stacks, 2); err != nil {
			return false
		}
		return !strings.Contains(stacks.String(), "validation.AsyncArgument.validate.func")
	}, 5*time.Second, 10*time.Millisecond, "async workers must exit after cancellation")
}
