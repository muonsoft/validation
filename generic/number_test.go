package generic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewNumber(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected Number
	}{
		{"uint", uint(123), Number{int: 123, float: 123, isInt: true}},
		{"uint8", uint8(123), Number{int: 123, float: 123, isInt: true}},
		{"uint16", uint16(123), Number{int: 123, float: 123, isInt: true}},
		{"uint32", uint32(123), Number{int: 123, float: 123, isInt: true}},
		{"uint64", uint64(123), Number{int: 123, float: 123, isInt: true}},
		{"int", 123, Number{int: 123, float: 123, isInt: true}},
		{"int8", int8(123), Number{int: 123, float: 123, isInt: true}},
		{"int16", int16(123), Number{int: 123, float: 123, isInt: true}},
		{"int32", int32(123), Number{int: 123, float: 123, isInt: true}},
		{"int64", int64(123), Number{int: 123, float: 123, isInt: true}},
		{"float32 as int", float32(123), Number{int: 123, float: 123, isInt: true}},
		{"float64 as int", float64(123), Number{int: 123, float: 123, isInt: true}},
		{"float32", float32(123.123), Number{int: 123, float: 123.12300109863281, isInt: false}},
		{"float64", float64(123.123), Number{int: 123, float: 123.123, isInt: false}},
		{"uint pointer", uintPointer(123), Number{int: 123, float: 123, isInt: true}},
		{"int pointer", intPointer(123), Number{int: 123, float: 123, isInt: true}},
		{"float64 pointer", floatPointer(123.123), Number{int: 123, float: 123.123, isInt: false}},
		{"nil uint pointer", nilUint, Number{isNil: true}},
		{"nil int pointer", nilInt, Number{isNil: true}},
		{"nil float64 pointer", nilFloat, Number{isNil: true}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			n, err := NewNumber(test.value)

			assert.NoError(t, err)
			if assert.NotNil(t, n) {
				assert.Equal(t, test.expected, *n)
			}
		})
	}
}

func TestNewNumberFromInt(t *testing.T) {
	n := NewNumberFromInt(1)

	assert.Equal(t, int64(1), n.int)
	assert.Equal(t, float64(1), n.float)
	assert.True(t, n.isInt)
}

func TestNewNumberFromFloat(t *testing.T) {
	n := NewNumberFromFloat(1.5)

	assert.Equal(t, int64(1), n.int)
	assert.Equal(t, 1.5, n.float)
	assert.False(t, n.isInt)
}

func TestNumber_String_WhenNumberFromInt_ExpectInt(t *testing.T) {
	n := NewNumberFromInt(123)

	s := n.String()

	assert.Equal(t, "123", s)
}

func TestNumber_String_WhenNumberFromFloat_ExpectFloat(t *testing.T) {
	n := NewNumberFromFloat(123.123)

	s := n.String()

	assert.Equal(t, "123.123", s)
}

func TestNewNumber_WhenNotNumeric_ExpectError(t *testing.T) {
	n, err := NewNumber(struct{}{})

	assert.Nil(t, n)
	assert.EqualError(t, err, "value of type struct is not numeric")
}

func TestMustNewNumber_WhenNotNumeric_ExpectPanic(t *testing.T) {
	assert.Panics(t, func() {
		MustNewNumber(struct{}{})
	})
}

func TestNumber_IsNil_WhenIsNil_ExpectTrue(t *testing.T) {
	n := MustNewNumber(nilInt)

	assert.True(t, n.IsNil())
}

func TestNumber_IsZero(t *testing.T) {
	tests := []struct {
		name   string
		n      Number
		isZero bool
	}{
		{"nil value is not zero", MustNewNumber(nilInt), false},
		{"zero int is zero", MustNewNumber(0), true},
		{"zero float is zero", MustNewNumber(0.0), true},
		{"float less than 1.0 is not zero", MustNewNumber(0.0123), false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			isZero := test.n.IsZero()

			assert.Equal(t, test.isZero, isZero)
		})
	}
}

func TestCompareNumbers(t *testing.T) {
	tests := []struct {
		name     string
		result   bool
		expected bool
	}{
		{"int equal to int", MustNewNumber(123).IsEqualTo(MustNewNumber(123)), true},
		{"int not equal to int", MustNewNumber(123).IsEqualTo(MustNewNumber(124)), false},
		{"nil not equal to any", MustNewNumber(nilInt).IsEqualTo(MustNewNumber(123)), false},
		{"any not equal to nil", MustNewNumber(123).IsEqualTo(MustNewNumber(nilInt)), false},
		{"float equal to float", MustNewNumber(123.123).IsEqualTo(MustNewNumber(123.123)), true},
		{"float not equal to float", MustNewNumber(123.123).IsEqualTo(MustNewNumber(123.124)), false},
		{"float equal to int", MustNewNumber(floatPointer(123)).IsEqualTo(MustNewNumber(123)), true},
		{"int greater than int", MustNewNumber(123).IsGreaterThan(MustNewNumber(122)), true},
		{"int not greater than int", MustNewNumber(123).IsGreaterThan(MustNewNumber(123)), false},
		{"nil not greater than any", MustNewNumber(nilInt).IsGreaterThan(MustNewNumber(123)), false},
		{"any not greater than nil", MustNewNumber(123).IsGreaterThan(MustNewNumber(nilInt)), false},
		{"float greater than float", MustNewNumber(123.123).IsGreaterThan(MustNewNumber(123.122)), true},
		{"float not greater than float", MustNewNumber(123.123).IsGreaterThan(MustNewNumber(123.123)), false},
		{"float greater than int", MustNewNumber(floatPointer(123.01)).IsGreaterThan(MustNewNumber(123)), true},
		{"int less than int", MustNewNumber(123).IsLessThan(MustNewNumber(124)), true},
		{"int not less than int", MustNewNumber(123).IsLessThan(MustNewNumber(123)), false},
		{"nil not less than any", MustNewNumber(nilInt).IsLessThan(MustNewNumber(123)), false},
		{"any not less than nil", MustNewNumber(123).IsLessThan(MustNewNumber(nilInt)), false},
		{"float less than float", MustNewNumber(123.123).IsLessThan(MustNewNumber(123.124)), true},
		{"float not less than float", MustNewNumber(123.123).IsLessThan(MustNewNumber(123.123)), false},
		{"float less than int", MustNewNumber(floatPointer(122.99)).IsLessThan(MustNewNumber(123)), true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, test.result)
		})
	}
}
