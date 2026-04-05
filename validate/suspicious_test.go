package validate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoSuspiciousCharacters_invisible(t *testing.T) {
	t.Parallel()
	assert.NoError(t, NoSuspiciousCharacters("hello"))
	assert.ErrorIs(t, NoSuspiciousCharacters("a\u200b"), ErrSuspiciousInvisible)
	assert.ErrorIs(t, NoSuspiciousCharacters("\ufeffx"), ErrSuspiciousInvisible)
}

func TestNoSuspiciousCharacters_mixedNumbers(t *testing.T) {
	t.Parallel()
	assert.NoError(t, NoSuspiciousCharacters("123"))
	assert.ErrorIs(t, NoSuspiciousCharacters("8৪"), ErrSuspiciousMixedNumbers)
	assert.ErrorIs(t, NoSuspiciousCharacters("৪٤"), ErrSuspiciousMixedNumbers)
}

func TestNoSuspiciousCharacters_hiddenOverlay(t *testing.T) {
	t.Parallel()
	assert.NoError(t, NoSuspiciousCharacters("café"))
	assert.ErrorIs(t, NoSuspiciousCharacters("i\u0307"), ErrSuspiciousHiddenOverlay)
}

func TestNoSuspiciousCharacters_singleScript(t *testing.T) {
	t.Parallel()
	assert.NoError(t, NoSuspiciousCharacters("hello", WithSuspiciousRestriction(SuspiciousRestrictionSingleScript)))
	assert.ErrorIs(t,
		NoSuspiciousCharacters("a\u0430", WithSuspiciousRestriction(SuspiciousRestrictionSingleScript)),
		ErrSuspiciousRestriction,
	)
}

func TestNoSuspiciousCharacters_locales(t *testing.T) {
	t.Parallel()
	assert.NoError(t, NoSuspiciousCharacters("hello", WithSuspiciousRestriction(SuspiciousRestrictionLocales), WithSuspiciousLocales("en")))
	assert.NoError(t, NoSuspiciousCharacters("πει", WithSuspiciousRestriction(SuspiciousRestrictionLocales), WithSuspiciousLocales("el")))
	assert.ErrorIs(t,
		NoSuspiciousCharacters("πει", WithSuspiciousRestriction(SuspiciousRestrictionLocales), WithSuspiciousLocales("en")),
		ErrSuspiciousRestriction,
	)
}

func TestNoSuspiciousCharacters_disableChecks(t *testing.T) {
	t.Parallel()
	assert.NoError(t, NoSuspiciousCharacters("8৪", WithSuspiciousChecks(0)))
}
