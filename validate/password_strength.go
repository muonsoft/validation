package validate

import (
	"errors"
	"math"
)

// PasswordStrengthScore is a heuristic score from 0 (very weak) to 4 (very strong).
type PasswordStrengthScore int

const (
	// PasswordStrengthVeryWeak is the lowest score, including empty passwords.
	PasswordStrengthVeryWeak PasswordStrengthScore = iota
	// PasswordStrengthWeak corresponds to an entropy estimate of at least 60.
	PasswordStrengthWeak
	// PasswordStrengthMedium corresponds to an entropy estimate of at least 80.
	PasswordStrengthMedium
	// PasswordStrengthStrong corresponds to an entropy estimate of at least 100.
	PasswordStrengthStrong
	// PasswordStrengthVeryStrong corresponds to an entropy estimate of at least 120.
	PasswordStrengthVeryStrong
)

var (
	// ErrPasswordTooWeak indicates a score below the configured minimum.
	ErrPasswordTooWeak = errors.New("password strength too low")
	// ErrInvalidPasswordStrengthMinimum indicates a minimum outside 1–4.
	ErrInvalidPasswordStrengthMinimum = errors.New("invalid password strength minimum")
	// ErrInvalidPasswordStrengthScore indicates an estimator result outside 0–4.
	ErrInvalidPasswordStrengthScore = errors.New("invalid password strength score")
)

// PasswordStrengthEstimator evaluates a password. It must return a score from 0 to 4.
// Implementations used concurrently must be safe for concurrent calls.
type PasswordStrengthEstimator func(string) PasswordStrengthScore

// PasswordStrengthOptions configures the minimum and optional custom estimator.
type PasswordStrengthOptions struct {
	minimum   PasswordStrengthScore
	estimator PasswordStrengthEstimator
}

// WithMinPasswordStrength sets the inclusive minimum score (1–4). Default is Medium.
func WithMinPasswordStrength(score PasswordStrengthScore) func(*PasswordStrengthOptions) {
	return func(o *PasswordStrengthOptions) { o.minimum = score }
}

// WithPasswordStrengthEstimator replaces the estimator. Nil restores the default.
func WithPasswordStrengthEstimator(estimator PasswordStrengthEstimator) func(*PasswordStrengthOptions) {
	return func(o *PasswordStrengthOptions) { o.estimator = estimator }
}

// PasswordStrength validates a password's score, defaulting to Medium.
// Empty strings are evaluated (score 0 by default), not skipped.
// It returns ErrPasswordTooWeak, ErrInvalidPasswordStrengthMinimum, or
// ErrInvalidPasswordStrengthScore on failure. Errors never contain the password.
func PasswordStrength(password string, options ...func(*PasswordStrengthOptions)) error {
	o := PasswordStrengthOptions{minimum: PasswordStrengthMedium}
	for _, option := range options {
		option(&o)
	}
	if o.minimum < PasswordStrengthWeak || o.minimum > PasswordStrengthVeryStrong {
		return ErrInvalidPasswordStrengthMinimum
	}
	if o.estimator == nil {
		o.estimator = EstimatePasswordStrength
	}
	score := o.estimator(password)
	if score < PasswordStrengthVeryWeak || score > PasswordStrengthVeryStrong {
		return ErrInvalidPasswordStrengthScore
	}
	if score < o.minimum {
		return ErrPasswordTooWeak
	}
	return nil
}

// EstimatePasswordStrength reproduces Symfony 8.0's byte-based heuristic.
// It uses byte length, unique byte count, and character category pool sizes.
// No Unicode normalization, dictionary, sequence detection, or breach lookup is
// performed. The score is not a guarantee of resistance to password guessing.
// Empty input returns VeryWeak. Runtime is O(len(password)), memory is O(1).
// See https://github.com/symfony/validator/blob/8.0/Constraints/PasswordStrengthValidator.php.
func EstimatePasswordStrength(password string) PasswordStrengthScore {
	if len(password) == 0 {
		return PasswordStrengthVeryWeak
	}
	var seen [256]bool
	var categories [6]bool
	unique := 0
	for i := 0; i < len(password); i++ {
		b := password[i]
		if !seen[b] {
			seen[b] = true
			unique++
			categories[passwordByteCategory(b)] = true
		}
	}
	pool := 0
	for i, size := range [...]int{33, 10, 26, 26, 128, 33} {
		if categories[i] {
			pool += size
		}
	}
	entropy := float64(unique)*math.Log2(float64(pool)) + float64(len(password)-unique)*math.Log2(float64(unique))
	return passwordEntropyScore(entropy)
}

func passwordByteCategory(b byte) int {
	switch {
	case b < 32 || b == 127:
		return 0
	case passwordByteInRange(b, '0', '9'):
		return 1
	case passwordByteInRange(b, 'A', 'Z'):
		return 2
	case passwordByteInRange(b, 'a', 'z'):
		return 3
	case b >= 128:
		return 4
	default:
		return 5
	}
}

func passwordEntropyScore(entropy float64) PasswordStrengthScore {
	switch {
	case entropy >= 120:
		return PasswordStrengthVeryStrong
	case entropy >= 100:
		return PasswordStrengthStrong
	case entropy >= 80:
		return PasswordStrengthMedium
	case entropy >= 60:
		return PasswordStrengthWeak
	default:
		return PasswordStrengthVeryWeak
	}
}

func passwordByteInRange(b, first, last byte) bool {
	return b >= first && b <= last
}
