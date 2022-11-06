package is_test

import (
	"fmt"
	"testing"

	"github.com/muonsoft/validation/is"
	"github.com/stretchr/testify/assert"
)

func TestDivisibleBy(t *testing.T) {
	tests := []struct {
		divisible float64
		divisor   float64
		want      bool
	}{
		{-7, 1, true},
		{0, 3.1415, true},
		{42, 42, true},
		{42, 21, true},
		{10.12, 0.01, true},
		{10.12, 0.001, true},
		{1.133, 0.001, true},
		{1.1331, 0.0001, true},
		{1.13331, 0.00001, true},
		{1.13331, 0.000001, true},
		{1, 0.1, true},
		{1, 0.01, true},
		{1, 0.001, true},
		{1, 0.0001, true},
		{1, 0.00001, true},
		{1, 0.000001, true},
		{3.25, 0.25, true},
		{100, 10, true},
		{4.1, 0.1, true},
		{-4.1, 0.1, true},
		{1, 2, false},
		{10, 3, false},
		{10, 0, false},
		{22, 10, false},
		{4.15, 0.1, false},
		{10.123, 0.01, false},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("divisible %v, divisior %v", test.divisible, test.divisor), func(t *testing.T) {
			got := is.DivisibleBy(test.divisible, test.divisor)

			assert.Equal(t, test.want, got)
		})
	}
}
