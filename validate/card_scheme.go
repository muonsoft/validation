package validate

import (
	"errors"
	"regexp"
)

// CardSchemeName identifies a supported payment card numbering scheme.
type CardSchemeName string

const (
	// CardSchemeAMEX selects the AMEX numbering rules from Symfony 8.0.
	CardSchemeAMEX CardSchemeName = "AMEX"
	// CardSchemeChinaUnionPay selects the CHINA_UNIONPAY numbering rules from Symfony 8.0.
	CardSchemeChinaUnionPay CardSchemeName = "CHINA_UNIONPAY"
	// CardSchemeDiners selects the DINERS numbering rules from Symfony 8.0.
	CardSchemeDiners CardSchemeName = "DINERS"
	// CardSchemeDiscover selects the DISCOVER numbering rules from Symfony 8.0.
	CardSchemeDiscover CardSchemeName = "DISCOVER"
	// CardSchemeInstaPayment selects the INSTAPAYMENT numbering rules from Symfony 8.0.
	CardSchemeInstaPayment CardSchemeName = "INSTAPAYMENT"
	// CardSchemeJCB selects the JCB numbering rules from Symfony 8.0.
	CardSchemeJCB CardSchemeName = "JCB"
	// CardSchemeLaser selects the LASER numbering rules from Symfony 8.0.
	CardSchemeLaser CardSchemeName = "LASER"
	// CardSchemeMaestro selects the MAESTRO numbering rules from Symfony 8.0.
	CardSchemeMaestro CardSchemeName = "MAESTRO"
	// CardSchemeMastercard selects the MASTERCARD numbering rules from Symfony 8.0.
	CardSchemeMastercard CardSchemeName = "MASTERCARD"
	// CardSchemeMIR selects the MIR numbering rules from Symfony 8.0.
	CardSchemeMIR CardSchemeName = "MIR"
	// CardSchemeUATP selects the UATP numbering rules from Symfony 8.0.
	CardSchemeUATP CardSchemeName = "UATP"
	// CardSchemeVisa selects the VISA numbering rules from Symfony 8.0.
	CardSchemeVisa CardSchemeName = "VISA"
)

var (
	// ErrInvalidCardScheme indicates a number that does not match any selected scheme.
	ErrInvalidCardScheme = errors.New("invalid card scheme")
	// ErrInvalidCardSchemes indicates an empty selection or an unknown scheme name.
	ErrInvalidCardSchemes = errors.New("invalid card schemes")
)

// Numbering patterns follow Symfony 8.0 CardSchemeValidator:
// https://github.com/symfony/validator/blob/8.0/Constraints/CardSchemeValidator.php
var cardSchemePatterns = map[CardSchemeName]*regexp.Regexp{
	CardSchemeAMEX:          regexp.MustCompile(`^3[47][0-9]{13}$`),
	CardSchemeChinaUnionPay: regexp.MustCompile(`^62[0-9]{14,17}$`),
	CardSchemeDiners:        regexp.MustCompile(`^3(0[0-5]|[68][0-9])[0-9]{11}$`),
	CardSchemeDiscover:      regexp.MustCompile(`^(6011[0-9]{12}|64[4-9][0-9]{13}|65[0-9]{14}|622(12[6-9]|1[3-9][0-9]|[2-8][0-9]{2}|91[0-9]|92[0-5])[0-9]{10})$`),
	CardSchemeInstaPayment:  regexp.MustCompile(`^63[7-9][0-9]{13}$`),
	CardSchemeJCB:           regexp.MustCompile(`^(2131|1800|35[0-9]{3})[0-9]{11}$`),
	CardSchemeLaser:         regexp.MustCompile(`^(6304|670[69]|6771)[0-9]{12,15}$`),
	CardSchemeMaestro:       regexp.MustCompile(`^(50[0-9]{10,17}|5[6-9][0-9]{10,17}|6[0-9]{11,18})$`),
	CardSchemeMastercard:    regexp.MustCompile(`^(5[1-5][0-9]{14}|2(22[1-9][0-9]{12}|2[3-9][0-9]{13}|[3-6][0-9]{14}|7[0-1][0-9]{13}|720[0-9]{12}))$`),
	CardSchemeMIR:           regexp.MustCompile(`^220[0-4][0-9]{12,15}$`),
	CardSchemeUATP:          regexp.MustCompile(`^1[0-9]{14}$`),
	CardSchemeVisa:          regexp.MustCompile(`^4([0-9]{12}|[0-9]{15}|[0-9]{18})$`),
}

// CardScheme checks a card number against one or more explicitly selected schemes.
// Only ASCII digits are accepted; spaces, hyphens, and other formatting are not
// removed. This checks prefixes and lengths, not the LUHN checksum or whether
// a card exists or can be charged. See [LUHN] for a separate checksum check.
// Empty values are valid; use [NotBlank] to require a value.
// Missing or unknown schemes return [ErrInvalidCardSchemes], even for empty input.
// Nonmatching numbers return [ErrInvalidCardScheme].
func CardScheme(value string, schemes ...CardSchemeName) error {
	if len(schemes) == 0 {
		return ErrInvalidCardSchemes
	}
	for _, scheme := range schemes {
		if _, ok := cardSchemePatterns[scheme]; !ok {
			return ErrInvalidCardSchemes
		}
	}
	if value == "" {
		return nil
	}
	for _, scheme := range schemes {
		if cardSchemePatterns[scheme].MatchString(value) {
			return nil
		}
	}
	return ErrInvalidCardScheme
}
