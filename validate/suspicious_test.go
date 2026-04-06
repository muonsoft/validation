package validate

import (
	"strings"
	"testing"
	"unicode"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNoSuspiciousCharacters_invisible(t *testing.T) {
	t.Parallel()
	assert.NoError(t, NoSuspiciousCharacters("hello"))
	assert.ErrorIs(t, NoSuspiciousCharacters("a\u200b"), ErrSuspiciousInvisible)
	assert.ErrorIs(t, NoSuspiciousCharacters("\ufeffx"), ErrSuspiciousInvisible)
	assert.ErrorIs(t, NoSuspiciousCharacters(string([]byte{0xff, 0xfe})), ErrSuspiciousInvisible)
}

func TestNoSuspiciousCharacters_mixedNumbers(t *testing.T) {
	t.Parallel()
	assert.NoError(t, NoSuspiciousCharacters("123"))
	assert.ErrorIs(t, NoSuspiciousCharacters("8৪"), ErrSuspiciousMixedNumbers)
	assert.ErrorIs(t, NoSuspiciousCharacters("৪٤"), ErrSuspiciousMixedNumbers)
	assert.NoError(t, NoSuspiciousCharacters("৪৫")) // same Bengali script
	assert.NoError(t, NoSuspiciousCharacters("٤٥")) // same Arabic script digits
}

func TestNoSuspiciousCharacters_hiddenOverlay(t *testing.T) {
	t.Parallel()
	assert.NoError(t, NoSuspiciousCharacters("café"))
	assert.ErrorIs(t, NoSuspiciousCharacters("i\u0307"), ErrSuspiciousHiddenOverlay)
	assert.ErrorIs(t, NoSuspiciousCharacters("j\u0307"), ErrSuspiciousHiddenOverlay)
	assert.ErrorIs(t, NoSuspiciousCharacters("l\u0307"), ErrSuspiciousHiddenOverlay)
	assert.NoError(t, NoSuspiciousCharacters("x\u0307"))
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
	assert.NoError(t, NoSuspiciousCharacters("かな", WithSuspiciousRestriction(SuspiciousRestrictionLocales), WithSuspiciousLocales("ja-Hira")))
	assert.ErrorIs(t,
		NoSuspiciousCharacters("٤٥", WithSuspiciousRestriction(SuspiciousRestrictionLocales), WithSuspiciousLocales("en")),
		ErrSuspiciousRestriction,
	)
}

func TestNoSuspiciousCharacters_whitespaceOnlyValid(t *testing.T) {
	t.Parallel()
	assert.NoError(t, NoSuspiciousCharacters("   "))
	assert.NoError(t, NoSuspiciousCharacters("\t\n"))
}

func TestNoSuspiciousCharacters_restrictionNoneIgnoresLocales(t *testing.T) {
	t.Parallel()
	// Greek letters should not be restricted when mode is None even if locales are passed.
	assert.NoError(t, NoSuspiciousCharacters("πει", WithSuspiciousRestriction(SuspiciousRestrictionNone), WithSuspiciousLocales("en")))
}

func TestNoSuspiciousCharacters_localesEmptyDefaultsToEnglish(t *testing.T) {
	t.Parallel()
	assert.NoError(t, NoSuspiciousCharacters("hello", WithSuspiciousRestriction(SuspiciousRestrictionLocales)))
	assert.ErrorIs(t,
		NoSuspiciousCharacters("πει", WithSuspiciousRestriction(SuspiciousRestrictionLocales)),
		ErrSuspiciousRestriction,
	)
}

func TestNoSuspiciousCharacters_localesInvalidTagsStillAllowLatin(t *testing.T) {
	t.Parallel()
	// Unparseable locale tags add no script tables; [buildLocaleAllowedTables] still appends Latin.
	assert.NoError(t, NoSuspiciousCharacters("hello", WithSuspiciousRestriction(SuspiciousRestrictionLocales), WithSuspiciousLocales("not-a-locale-!!!")))
	assert.ErrorIs(t,
		NoSuspiciousCharacters("πει", WithSuspiciousRestriction(SuspiciousRestrictionLocales), WithSuspiciousLocales("not-a-locale-!!!")),
		ErrSuspiciousRestriction,
	)
}

func TestNoSuspiciousCharacters_checkMaskCombinations(t *testing.T) {
	t.Parallel()
	invis := "a\u200b"
	mixed := "8৪"
	overlay := "i\u0307"

	cases := []struct {
		name    string
		opts    []NoSuspiciousCharactersOption
		value   string
		wantErr error
	}{
		{"disable all", []NoSuspiciousCharactersOption{WithSuspiciousChecks(0)}, mixed, nil},
		{"only invisible on invis", []NoSuspiciousCharactersOption{WithSuspiciousChecks(CheckSuspiciousInvisible)}, invis, ErrSuspiciousInvisible},
		{"only invisible passes clean", []NoSuspiciousCharactersOption{WithSuspiciousChecks(CheckSuspiciousInvisible)}, "ok", nil},
		{"only invisible ignores mixed", []NoSuspiciousCharactersOption{WithSuspiciousChecks(CheckSuspiciousInvisible)}, mixed, nil},
		{"only mixed on mixed", []NoSuspiciousCharactersOption{WithSuspiciousChecks(CheckSuspiciousMixedNumbers)}, mixed, ErrSuspiciousMixedNumbers},
		{"only mixed ignores invis", []NoSuspiciousCharactersOption{WithSuspiciousChecks(CheckSuspiciousMixedNumbers)}, invis, nil},
		{"only overlay on overlay", []NoSuspiciousCharactersOption{WithSuspiciousChecks(CheckSuspiciousHiddenOverlay)}, overlay, ErrSuspiciousHiddenOverlay},
		{"invis+mixed", []NoSuspiciousCharactersOption{WithSuspiciousChecks(CheckSuspiciousInvisible | CheckSuspiciousMixedNumbers)}, invis, ErrSuspiciousInvisible},
		{"invis+mixed second", []NoSuspiciousCharactersOption{WithSuspiciousChecks(CheckSuspiciousInvisible | CheckSuspiciousMixedNumbers)}, mixed, ErrSuspiciousMixedNumbers},
		{"default all hits invis first", nil, invis, ErrSuspiciousInvisible},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := NoSuspiciousCharacters(tc.value, tc.opts...)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNoSuspiciousCharacters_restrictionWithPartialChecks(t *testing.T) {
	t.Parallel()
	// Cyrillic not allowed for en; mixed-numbers check off so we reach restriction.
	opts := []NoSuspiciousCharactersOption{
		WithSuspiciousChecks(CheckSuspiciousInvisible | CheckSuspiciousHiddenOverlay),
		WithSuspiciousRestriction(SuspiciousRestrictionLocales),
		WithSuspiciousLocales("en"),
	}
	assert.ErrorIs(t, NoSuspiciousCharacters("ы", opts...), ErrSuspiciousRestriction)
}

func TestUnicodeRangeTablesForISO15924(t *testing.T) {
	t.Parallel()
	assert.Nil(t, unicodeRangeTablesForISO15924(""))
	assert.Nil(t, unicodeRangeTablesForISO15924("Zzzz"))
	require.NotEmpty(t, unicodeRangeTablesForISO15924("Latn"))
	require.NotEmpty(t, unicodeRangeTablesForISO15924("Jpan"))
	require.NotEmpty(t, unicodeRangeTablesForISO15924("Hira"))
	require.NotEmpty(t, unicodeRangeTablesForISO15924("Kana"))
	// Unknown ISO code: if it matches a unicode.Scripts key, return one table.
	if rt := unicodeRangeTablesForISO15924("Thai"); len(rt) != 1 {
		t.Fatalf("Thai: got %d tables", len(rt))
	}
}

func TestScriptNameForRune(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "Latin", scriptNameForRune('a'))
	assert.Equal(t, "Greek", scriptNameForRune('π'))
	assert.Equal(t, "", scriptNameForRune(-1))
}

func TestBuildLocaleAllowedTables_includesLatin(t *testing.T) {
	t.Parallel()
	tab := buildLocaleAllowedTables([]string{"el"})
	var hasLatin bool
	for _, rt := range tab {
		if rt == unicode.Scripts["Latin"] {
			hasLatin = true
			break
		}
	}
	assert.True(t, hasLatin, "Latin should be appended for locale lists")
}

func TestApplySuspiciousCheckBitmask_emptyMaskNoop(t *testing.T) {
	t.Parallel()
	rs := []rune("a\u200b")
	assert.NoError(t, applySuspiciousCheckBitmask(rs, 0))
}

func TestNoSuspiciousCharacters_zlZpInvisible(t *testing.T) {
	t.Parallel()
	// U+2028 LINE SEPARATOR (Zl), U+2029 PARAGRAPH SEPARATOR (Zp)
	assert.ErrorIs(t, NoSuspiciousCharacters("\u2028x"), ErrSuspiciousInvisible)
	assert.ErrorIs(t, NoSuspiciousCharacters("\u2029x"), ErrSuspiciousInvisible)
}

func TestNoSuspiciousCharacters_nulInvisible(t *testing.T) {
	t.Parallel()
	assert.ErrorIs(t, NoSuspiciousCharacters("a\x00b"), ErrSuspiciousInvisible)
}

func TestCheckSuspiciousInvisible_bidiFormatCf(t *testing.T) {
	t.Parallel()
	// U+202E RIGHT-TO-LEFT OVERRIDE is Cf
	assert.ErrorIs(t, checkSuspiciousInvisible([]rune("a\u202eb")), ErrSuspiciousInvisible)
}

func TestNoSuspiciousCharacters_longStringStillScanned(t *testing.T) {
	t.Parallel()
	s := strings.Repeat("a", 500) + "\u200b"
	assert.ErrorIs(t, NoSuspiciousCharacters(s), ErrSuspiciousInvisible)
}

func TestNeedsScriptCheck_letterNumberOther(t *testing.T) {
	t.Parallel()
	// U+2160 ROMAN NUMERAL ONE is category Nl; exercise needsScriptCheck Nl branch with locale restriction.
	assert.True(t, needsScriptCheck('\u2160'))
	assert.NoError(t, NoSuspiciousCharacters("\u2160", WithSuspiciousRestriction(SuspiciousRestrictionLocales), WithSuspiciousLocales("en")))
}

func TestUnicodeRangeTablesForISO15924_unknownCode(t *testing.T) {
	t.Parallel()
	assert.Nil(t, unicodeRangeTablesForISO15924("Qabx")) // private-use script code, not in map / not a Scripts table key
}
