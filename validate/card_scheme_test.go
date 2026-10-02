package validate_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/muonsoft/validation/validate"
	"github.com/stretchr/testify/require"
)

func TestCardSchemePrefixesAndLengths(t *testing.T) {
	tests := []struct {
		scheme   validate.CardSchemeName
		prefixes []string
		lengths  []int
	}{
		{validate.CardSchemeAMEX, []string{"34", "37"}, []int{15}},
		{validate.CardSchemeChinaUnionPay, []string{"62"}, []int{16, 17, 18, 19}},
		{validate.CardSchemeDiners, []string{"300", "301", "302", "303", "304", "305", "36", "38"}, []int{14}},
		{validate.CardSchemeDiscover, []string{"6011", "644", "645", "649", "65", "622126", "622129", "622130", "622199", "622200", "622899", "622910", "622919", "622920", "622925"}, []int{16}},
		{validate.CardSchemeInstaPayment, []string{"637", "638", "639"}, []int{16}},
		{validate.CardSchemeJCB, []string{"2131", "1800"}, []int{15}},
		{validate.CardSchemeJCB, []string{"35"}, []int{16}},
		{validate.CardSchemeLaser, []string{"6304", "6706", "6709", "6771"}, []int{16, 17, 18, 19}},
		{validate.CardSchemeMaestro, []string{"50", "56", "57", "58", "59", "60", "675900", "675999", "69"}, []int{12, 13, 14, 15, 16, 17, 18, 19}},
		{validate.CardSchemeMastercard, []string{"51", "52", "53", "54", "55", "222100", "222999", "223000", "229999", "230000", "269999", "270000", "271999", "272000", "272099"}, []int{16}},
		{validate.CardSchemeMIR, []string{"2200", "2201", "2202", "2203", "2204"}, []int{16, 17, 18, 19}},
		{validate.CardSchemeUATP, []string{"1"}, []int{15}},
		{validate.CardSchemeVisa, []string{"4"}, []int{13, 16, 19}},
	}
	for _, tt := range tests {
		for _, prefix := range tt.prefixes {
			for length := 11; length <= 20; length++ {
				t.Run(fmt.Sprintf("%s/%s/%d", tt.scheme, prefix, length), func(t *testing.T) {
					value := prefix + strings.Repeat("0", length-len(prefix))
					want := validate.ErrInvalidCardScheme
					for _, allowed := range tt.lengths {
						if length == allowed {
							want = nil
						}
					}
					require.ErrorIs(t, validate.CardScheme(value, tt.scheme), want)
				})
			}
		}
	}
}

func TestCardSchemeInvalidPrefixes(t *testing.T) {
	tests := []struct {
		scheme   validate.CardSchemeName
		prefixes []string
		length   int
	}{
		{validate.CardSchemeAMEX, []string{"33", "35", "36", "38"}, 15},
		{validate.CardSchemeChinaUnionPay, []string{"61", "63"}, 16},
		{validate.CardSchemeDiners, []string{"299", "306", "35", "37", "39"}, 14},
		{validate.CardSchemeDiscover, []string{"6010", "6012", "643", "622125", "622900", "622909", "622926"}, 16},
		{validate.CardSchemeInstaPayment, []string{"636", "640"}, 16},
		{validate.CardSchemeJCB, []string{"2130", "2132", "1799", "1801"}, 15},
		{validate.CardSchemeJCB, []string{"34", "36"}, 16},
		{validate.CardSchemeLaser, []string{"6303", "6305", "6705", "6707", "6708", "6710", "6770", "6772"}, 16},
		{validate.CardSchemeMaestro, []string{"49", "51", "55", "70"}, 12},
		{validate.CardSchemeMastercard, []string{"50", "56", "222099", "272100"}, 16},
		{validate.CardSchemeMIR, []string{"2199", "2205"}, 16},
		{validate.CardSchemeUATP, []string{"0", "2"}, 15},
		{validate.CardSchemeVisa, []string{"3", "5"}, 16},
	}
	for _, tt := range tests {
		for _, prefix := range tt.prefixes {
			t.Run(string(tt.scheme)+"/"+prefix, func(t *testing.T) {
				require.ErrorIs(t, validate.CardScheme(prefix+strings.Repeat("0", tt.length-len(prefix)), tt.scheme), validate.ErrInvalidCardScheme)
			})
		}
	}
}

func TestCardSchemeInputAndSelection(t *testing.T) {
	for _, value := range []string{" ", "4111 1111 1111 1111", "4111-1111-1111-1111", "4111111111111111\n", " 4111111111111111", "4111111111111111 ", "+4111111111111111", "4.11111111111111", "4e15", "４１１１１１１１１１１１１１１１", "411111111111111a", "411111111111111\xff"} {
		require.ErrorIs(t, validate.CardScheme(value, validate.CardSchemeVisa), validate.ErrInvalidCardScheme, value)
	}
	require.NoError(t, validate.CardScheme("", validate.CardSchemeVisa))
	require.NoError(t, validate.CardScheme("5100000000000000", validate.CardSchemeVisa, validate.CardSchemeMastercard))
	require.NoError(t, validate.CardScheme("4111111111111111", validate.CardSchemeVisa, validate.CardSchemeVisa))
	require.ErrorIs(t, validate.CardScheme("5100000000000000", validate.CardSchemeVisa), validate.ErrInvalidCardScheme)
	for _, value := range []string{"", "4111111111111111"} {
		require.ErrorIs(t, validate.CardScheme(value), validate.ErrInvalidCardSchemes)
		require.ErrorIs(t, validate.CardScheme(value, "visa"), validate.ErrInvalidCardSchemes)
		require.ErrorIs(t, validate.CardScheme(value, validate.CardSchemeVisa, "unknown"), validate.ErrInvalidCardSchemes)
	}
	// Scheme matching deliberately does not validate the checksum.
	require.NoError(t, validate.CardScheme("4111111111111112", validate.CardSchemeVisa))
	require.Error(t, validate.LUHN("4111111111111112"))
}

func ExampleCardScheme() {
	fmt.Println(validate.CardScheme("4111111111111111", validate.CardSchemeVisa))
	fmt.Println(validate.CardScheme("4111111111111111", validate.CardSchemeMastercard))
	// Output:
	// <nil>
	// invalid card scheme
}
