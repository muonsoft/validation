package generic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewIterable_WhenNotIterable_ExpectNotIterableError(t *testing.T) {
	iterable, err := NewIterable(0)

	assert.Nil(t, iterable)
	assert.EqualError(t, err, "value of type int is not iterable")
}

func TestNewIterable(t *testing.T) {
	tests := []struct {
		name            string
		value           interface{}
		expectedKey     string
		expectedIsIndex bool
		expectedValue   interface{}
	}{
		{
			name:            "array of strings",
			value:           [...]string{"value"},
			expectedKey:     "0",
			expectedIsIndex: true,
			expectedValue:   "value",
		},
		{
			name:            "array of string pointers",
			value:           [...]*string{stringPointer("value")},
			expectedKey:     "0",
			expectedIsIndex: true,
			expectedValue:   "value",
		},
		{
			name:            "array of string pointers with nil",
			value:           [...]*string{nil},
			expectedKey:     "0",
			expectedIsIndex: true,
			expectedValue:   nil,
		},
		{
			name:            "slice of strings",
			value:           []string{"value"},
			expectedKey:     "0",
			expectedIsIndex: true,
			expectedValue:   "value",
		},
		{
			name:            "slice of string pointers",
			value:           []*string{stringPointer("value")},
			expectedKey:     "0",
			expectedIsIndex: true,
			expectedValue:   "value",
		},
		{
			name:            "slice of string pointers with nil",
			value:           []*string{nil},
			expectedKey:     "0",
			expectedIsIndex: true,
			expectedValue:   nil,
		},
		{
			name:            "map of strings",
			value:           map[string]string{"key": "value"},
			expectedKey:     "key",
			expectedIsIndex: false,
			expectedValue:   "value",
		},
		{
			name:            "map of string pointers",
			value:           map[string]*string{"key": stringPointer("value")},
			expectedKey:     "key",
			expectedIsIndex: false,
			expectedValue:   "value",
		},
		{
			name:            "map of string pointers with nil",
			value:           map[string]*string{"key": nil},
			expectedKey:     "key",
			expectedIsIndex: false,
			expectedValue:   nil,
		},
		{
			name:            "map of string pointers with pointer key",
			value:           map[*string]*string{stringPointer("key"): nil},
			expectedKey:     "key",
			expectedIsIndex: false,
			expectedValue:   nil,
		},
		{
			name:            "map of string pointers with nil pointer key",
			value:           map[*string]*string{nil: nil},
			expectedKey:     "",
			expectedIsIndex: false,
			expectedValue:   nil,
		},
		{
			name:            "map of string pointers with struct key",
			value:           map[mapKey]string{{Key: "key"}: "value"},
			expectedKey:     "<generic.mapKey Value>",
			expectedIsIndex: false,
			expectedValue:   "value",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			iterable, err := NewIterable(test.value)

			if assert.NoError(t, err) {
				assert.Equal(t, 1, iterable.Count())
				i := 0
				iterable.Iterate(func(key Key, value interface{}) {
					assert.Equal(t, test.expectedKey, key.String())
					assert.Equal(t, test.expectedIsIndex, key.IsIndex())
					assert.Equal(t, test.expectedValue, value)
					i++
				})
				assert.Equal(t, 1, i)
			}
		})
	}
}
