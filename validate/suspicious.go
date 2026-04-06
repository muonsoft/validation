package validate

import (
	"errors"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/language"
)

// Unicode script names from [unicode.Scripts] map keys (for goconst / clarity).
const (
	unicodeScriptCommon    = "Common"
	unicodeScriptInherited = "Inherited"
)

// Suspicious check flags mirror Symfony NoSuspiciousCharacters bit values for API familiarity.
// See https://symfony.com/doc/current/reference/constraints/NoSuspiciousCharacters.html
const (
	CheckSuspiciousInvisible     uint = 32  // CHECK_INVISIBLE
	CheckSuspiciousMixedNumbers  uint = 128 // CHECK_MIXED_NUMBERS
	CheckSuspiciousHiddenOverlay uint = 256 // CHECK_HIDDEN_OVERLAY
)

// DefaultSuspiciousChecks runs all standard spoofing checks (invisible, mixed digit scripts, hidden overlay).
const DefaultSuspiciousChecks = CheckSuspiciousInvisible | CheckSuspiciousMixedNumbers | CheckSuspiciousHiddenOverlay

// Restriction for script/locale-related validation (approximation of Symfony restriction levels).
type SuspiciousRestriction int

const (
	// SuspiciousRestrictionNone disables locale/script checks (like RESTRICTION_LEVEL_NONE).
	SuspiciousRestrictionNone SuspiciousRestriction = iota
	// SuspiciousRestrictionLocales restricts letters/marks to scripts associated with [NoSuspiciousCharactersOptions.Locales]
	// (like a moderate locale allow-list). When Locales is empty, [language.English] is used.
	SuspiciousRestrictionLocales
	// SuspiciousRestrictionSingleScript requires at most one non-Common, non-Inherited script among letters/marks/numbers.
	SuspiciousRestrictionSingleScript
)

var (
	// ErrSuspiciousInvisible is returned when invisible or format-control characters are present.
	ErrSuspiciousInvisible = errors.New("suspicious invisible characters")
	// ErrSuspiciousMixedNumbers is returned when decimal digits from more than one numbering system are mixed.
	ErrSuspiciousMixedNumbers = errors.New("suspicious mixed digit scripts")
	// ErrSuspiciousHiddenOverlay is returned for combining sequences that can hide in the preceding glyph (e.g. Latin i + U+0307).
	ErrSuspiciousHiddenOverlay = errors.New("suspicious hidden overlay")
	// ErrSuspiciousRestriction is returned when a rune’s script is not allowed for the configured locales or single-script rule.
	ErrSuspiciousRestriction = errors.New("suspicious script restriction")
)

// NoSuspiciousCharactersOptions configures [NoSuspiciousCharacters].
type NoSuspiciousCharactersOptions struct {
	Checks      uint
	checksSet   bool
	Restriction SuspiciousRestriction
	Locales     []string
}

// NoSuspiciousCharactersOption mutates [NoSuspiciousCharactersOptions].
type NoSuspiciousCharactersOption func(*NoSuspiciousCharactersOptions)

// WithSuspiciousChecks sets the bitmask of checks (Symfony-compatible values). Zero means “use default all checks”.
func WithSuspiciousChecks(checks uint) NoSuspiciousCharactersOption {
	return func(o *NoSuspiciousCharactersOptions) {
		o.Checks = checks
		o.checksSet = true
	}
}

// WithSuspiciousRestriction sets locale/script restriction mode.
func WithSuspiciousRestriction(r SuspiciousRestriction) NoSuspiciousCharactersOption {
	return func(o *NoSuspiciousCharactersOptions) {
		o.Restriction = r
	}
}

// WithSuspiciousLocales sets BCP 47 locale tags used for [SuspiciousRestrictionLocales] (e.g. "en", "en-US").
func WithSuspiciousLocales(locales ...string) NoSuspiciousCharactersOption {
	return func(o *NoSuspiciousCharactersOptions) {
		o.Locales = append([]string(nil), locales...)
	}
}

// NoSuspiciousCharacters reports whether value is free of common spoofing patterns (homoglyphs, invisible characters,
// mixed-script digits, risky combining sequences). Empty and whitespace-only strings are valid.
//
// Semantics are inspired by Symfony’s NoSuspiciousCharacters (ICU Spoofchecker) but implemented with the Go standard
// library and [golang.org/x/text/language] only; results may differ from ICU for edge cases.
//
// Possible errors:
//   - [ErrSuspiciousInvisible]
//   - [ErrSuspiciousMixedNumbers]
//   - [ErrSuspiciousHiddenOverlay]
//   - [ErrSuspiciousRestriction]
func NoSuspiciousCharacters(value string, options ...NoSuspiciousCharactersOption) error {
	opts := NoSuspiciousCharactersOptions{}
	for _, opt := range options {
		opt(&opts)
	}
	if !opts.checksSet {
		opts.Checks = DefaultSuspiciousChecks
	}
	if !utf8.ValidString(value) {
		return ErrSuspiciousInvisible
	}
	if strings.TrimSpace(value) == "" {
		return nil
	}
	rs := []rune(value)
	if err := applySuspiciousCheckBitmask(rs, opts.Checks); err != nil {
		return err
	}
	return applySuspiciousRestriction(rs, opts.Restriction, opts.Locales)
}

func applySuspiciousCheckBitmask(rs []rune, checks uint) error {
	if checks&CheckSuspiciousInvisible != 0 {
		if err := checkSuspiciousInvisible(rs); err != nil {
			return err
		}
	}
	if checks&CheckSuspiciousMixedNumbers != 0 {
		if err := checkSuspiciousMixedNumbers(rs); err != nil {
			return err
		}
	}
	if checks&CheckSuspiciousHiddenOverlay != 0 {
		if err := checkSuspiciousHiddenOverlay(rs); err != nil {
			return err
		}
	}
	return nil
}

func applySuspiciousRestriction(rs []rune, r SuspiciousRestriction, locales []string) error {
	switch r {
	case SuspiciousRestrictionLocales:
		loc := locales
		if len(loc) == 0 {
			loc = []string{"en"}
		}
		return checkSuspiciousLocales(rs, loc)
	case SuspiciousRestrictionSingleScript:
		return checkSuspiciousSingleScript(rs)
	default:
		return nil
	}
}

func checkSuspiciousInvisible(runes []rune) error {
	for _, r := range runes {
		if r == 0 {
			return ErrSuspiciousInvisible
		}
		// Category Cf contains many invisible / format controls used in spoofing (ZW*, BOM, bidi overrides, etc.).
		if unicode.Is(unicode.Cf, r) {
			return ErrSuspiciousInvisible
		}
		// Line and paragraph separators are often non-obvious in UI.
		if unicode.Is(unicode.Zl, r) || unicode.Is(unicode.Zp, r) {
			return ErrSuspiciousInvisible
		}
	}
	return nil
}

func checkSuspiciousMixedNumbers(runes []rune) error {
	// ASCII and many “European” digits have Unicode script Common; other decimal digits
	// belong to a specific script (e.g. Bengali, Arabic). Mixing those buckets matches
	// Symfony’s “mixed numbering systems” intent.
	var hasCommonDigit bool
	var scriptDigit string
	for _, r := range runes {
		if !unicode.Is(unicode.Nd, r) {
			continue
		}
		scr := scriptNameForRune(r)
		if scr == "" || scr == unicodeScriptCommon || scr == unicodeScriptInherited {
			hasCommonDigit = true
			if scriptDigit != "" {
				return ErrSuspiciousMixedNumbers
			}
			continue
		}
		if scriptDigit == "" {
			scriptDigit = scr
			if hasCommonDigit {
				return ErrSuspiciousMixedNumbers
			}
			continue
		}
		if scr != scriptDigit {
			return ErrSuspiciousMixedNumbers
		}
	}
	return nil
}

func checkSuspiciousHiddenOverlay(runes []rune) error {
	for i := 1; i < len(runes); i++ {
		prev, cur := runes[i-1], runes[i]
		if cur != '\u0307' {
			continue
		}
		// Latin small i/j/l + combining dot above (homograph for dotted Latin letters).
		if prev == 'i' || prev == 'j' || prev == 'l' {
			return ErrSuspiciousHiddenOverlay
		}
	}
	return nil
}

// iso15924ToUnicodeScript maps BCP 47 / ISO 15924 script codes (from [language.Tag.Script])
// to keys in [unicode.Scripts]. Composite tags like Jpan/Kore expand to multiple scripts.
var iso15924ToUnicodeScripts = map[string][]string{
	"Latn": {"Latin"},
	"Cyrl": {"Cyrillic"},
	"Grek": {"Greek"},
	"Arab": {"Arabic"},
	"Hebr": {"Hebrew"},
	"Thai": {"Thai"},
	"Deva": {"Devanagari"},
	"Beng": {"Bengali"},
	"Guru": {"Gurmukhi"},
	"Gujr": {"Gujarati"},
	"Orya": {"Oriya"},
	"Taml": {"Tamil"},
	"Telu": {"Telugu"},
	"Knda": {"Kannada"},
	"Mlym": {"Malayalam"},
	"Sinh": {"Sinhala"},
	"Laoo": {"Lao"},
	"Tibt": {"Tibetan"},
	"Mymr": {"Myanmar"},
	"Geor": {"Georgian"},
	"Ethi": {"Ethiopic"},
	"Cher": {"Cherokee"},
	"Cans": {"Canadian_Aboriginal"},
	"Khmr": {"Khmer"},
	"Hans": {"Han"},
	"Hant": {"Han"},
	"Hani": {"Han"},
	"Jpan": {"Hiragana", "Katakana", "Han"},
	"Kore": {"Hangul", "Han"},
	// BCP 47 / ISO 15924 script subtag aliases (e.g. ja-Hira, ja-Kana).
	"Hira": {"Hiragana"},
	"Kana": {"Katakana"},
	"Hang": {"Hangul"},
	"Bopo": {"Bopomofo"},
}

func unicodeRangeTablesForISO15924(code string) []*unicode.RangeTable {
	if code == "" || code == "Zzzz" {
		return nil
	}
	if names, ok := iso15924ToUnicodeScripts[code]; ok {
		var out []*unicode.RangeTable
		for _, n := range names {
			if rt, ok := unicode.Scripts[n]; ok {
				out = append(out, rt)
			}
		}
		return out
	}
	if rt, ok := unicode.Scripts[code]; ok {
		return []*unicode.RangeTable{rt}
	}
	return nil
}

func buildLocaleAllowedTables(locales []string) []*unicode.RangeTable {
	var allowed []*unicode.RangeTable
	seen := make(map[*unicode.RangeTable]struct{})
	for _, loc := range locales {
		tag, err := language.Parse(loc)
		if err != nil {
			continue
		}
		scr, _ := tag.Script()
		for _, rt := range unicodeRangeTablesForISO15924(scr.String()) {
			if _, dup := seen[rt]; dup {
				continue
			}
			seen[rt] = struct{}{}
			allowed = append(allowed, rt)
		}
	}
	if latin, ok := unicode.Scripts["Latin"]; ok {
		if _, dup := seen[latin]; !dup {
			seen[latin] = struct{}{}
			allowed = append(allowed, latin)
		}
	}
	return allowed
}

func runeInAnyScriptTable(r rune, tables []*unicode.RangeTable) bool {
	for _, rt := range tables {
		if unicode.Is(rt, r) {
			return true
		}
	}
	return false
}

func checkSuspiciousLocales(runes []rune, locales []string) error {
	allowed := buildLocaleAllowedTables(locales)
	if len(allowed) == 0 {
		return nil
	}
	for _, r := range runes {
		if !needsScriptCheck(r) {
			continue
		}
		if unicodeIsOnlyCommonInherited(r) {
			continue
		}
		if !runeInAnyScriptTable(r, allowed) {
			return ErrSuspiciousRestriction
		}
	}
	return nil
}

func unicodeIsOnlyCommonInherited(r rune) bool {
	scr := scriptNameForRune(r)
	return scr == "" || scr == unicodeScriptCommon || scr == unicodeScriptInherited
}

func checkSuspiciousSingleScript(runes []rune) error {
	var primary string
	for _, r := range runes {
		if !needsScriptCheck(r) && !unicode.Is(unicode.Nd, r) {
			continue
		}
		scr := scriptNameForRune(r)
		if scr == "" || scr == unicodeScriptCommon || scr == unicodeScriptInherited {
			continue
		}
		if primary == "" {
			primary = scr
			continue
		}
		if scr != primary {
			return ErrSuspiciousRestriction
		}
	}
	return nil
}

func needsScriptCheck(r rune) bool {
	if unicode.IsLetter(r) {
		return true
	}
	if unicode.IsMark(r) {
		return true
	}
	if unicode.Is(unicode.Nd, r) {
		return true
	}
	if unicode.Is(unicode.Nl, r) || unicode.Is(unicode.No, r) {
		return true
	}
	return false
}

func scriptNameForRune(r rune) string {
	for name, table := range unicode.Scripts {
		if unicode.Is(table, r) {
			return name
		}
	}
	return ""
}
