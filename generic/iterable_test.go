package generic

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewIterable_WhenNotIterable_ExpectNotIterableError(t *testing.T) {
	iterable, err := NewIterable(0)

	assert.Nil(t, iterable)
	assert.EqualError(t, err, "value of type int is not iterable")
}

func TestNewIterable_WhenNilSlice_ExpectIsNil(t *testing.T) {
	var slice []string

	iterable, err := NewIterable(slice)

	if assert.NoError(t, err) {
		assert.True(t, iterable.IsNil())
	}
}

func TestNewIterable_WhenNilMap_ExpectIsNil(t *testing.T) {
	var value map[string]string

	iterable, err := NewIterable(value)

	if assert.NoError(t, err) {
		assert.True(t, iterable.IsNil())
	}
}

func TestNewIterable_WhenSliceElementImplementsInterface_ExpectTrue(t *testing.T) {
	var value []testInterface
	interfaceType := reflect.TypeOf((*testInterface)(nil)).Elem()

	iterable, err := NewIterable(value)

	if assert.NoError(t, err) {
		implements := iterable.IsElementImplements(interfaceType)
		assert.True(t, implements)
	}
}

func TestNewIterable_WhenMapElementImplementsInterface_ExpectTrue(t *testing.T) {
	var value map[string]testInterface
	interfaceType := reflect.TypeOf((*testInterface)(nil)).Elem()

	iterable, err := NewIterable(value)

	if assert.NoError(t, err) {
		implements := iterable.IsElementImplements(interfaceType)
		assert.True(t, implements)
	}
}

func TestIterableArray_Iterate_WhenBreakAtFirstElement_ExpectCountIsOne(t *testing.T) {
	count := 0
	iterable, err := NewIterable([]string{"a", "b", "c"})
	if err != nil {
		t.Fatal("failed to initialize iterable from slice")
	}

	err = iterable.Iterate(func(key Key, value interface{}) error {
		count++
		return fmt.Errorf("error")
	})

	assert.EqualError(t, err, "error")
	assert.Equal(t, 1, count)
}

func TestIterableMap_Iterate_WhenBreakAtFirstElement_ExpectCountIsOne(t *testing.T) {
	count := 0
	iterable, err := NewIterable(map[string]string{"a": "a", "b": "b", "c": "c"})
	if err != nil {
		t.Fatal("failed to initialize iterable from map")
	}

	err = iterable.Iterate(func(key Key, value interface{}) error {
		count++
		return fmt.Errorf("error")
	})

	assert.EqualError(t, err, "error")
	assert.Equal(t, 1, count)
}

func TestNewIterable(t *testing.T) {
	tests := []struct {
		name            string
		value           interface{}
		expectedKey     string
		expectedIsIndex bool
		expectedIndex   int
		expectedValue   interface{}
	}{
		{
			name:            "array of strings",
			value:           [...]string{"value"},
			expectedKey:     "0",
			expectedIsIndex: true,
			expectedIndex:   0,
			expectedValue:   "value",
		},
		{
			name:            "array of string pointers",
			value:           [...]*string{stringPointer("value")},
			expectedKey:     "0",
			expectedIsIndex: true,
			expectedIndex:   0,
			expectedValue:   "value",
		},
		{
			name:            "array of string pointers with nil",
			value:           [...]*string{nil},
			expectedKey:     "0",
			expectedIsIndex: true,
			expectedIndex:   0,
			expectedValue:   nil,
		},
		{
			name:            "slice of strings",
			value:           []string{"value"},
			expectedKey:     "0",
			expectedIsIndex: true,
			expectedIndex:   0,
			expectedValue:   "value",
		},
		{
			name:            "slice of string pointers",
			value:           []*string{stringPointer("value")},
			expectedKey:     "0",
			expectedIsIndex: true,
			expectedIndex:   0,
			expectedValue:   "value",
		},
		{
			name:            "slice of string pointers with nil",
			value:           []*string{nil},
			expectedKey:     "0",
			expectedIsIndex: true,
			expectedIndex:   0,
			expectedValue:   nil,
		},
		{
			name:            "map of strings",
			value:           map[string]string{"key": "value"},
			expectedKey:     "key",
			expectedIsIndex: false,
			expectedIndex:   -1,
			expectedValue:   "value",
		},
		{
			name:            "map of string pointers",
			value:           map[string]*string{"key": stringPointer("value")},
			expectedKey:     "key",
			expectedIsIndex: false,
			expectedIndex:   -1,
			expectedValue:   "value",
		},
		{
			name:            "map of string pointers with nil",
			value:           map[string]*string{"key": nil},
			expectedKey:     "key",
			expectedIsIndex: false,
			expectedIndex:   -1,
			expectedValue:   nil,
		},
		{
			name:            "map of string pointers with pointer key",
			value:           map[*string]*string{stringPointer("key"): nil},
			expectedKey:     "key",
			expectedIsIndex: false,
			expectedIndex:   -1,
			expectedValue:   nil,
		},
		{
			name:            "map of string pointers with nil pointer key",
			value:           map[*string]*string{nil: nil},
			expectedKey:     "",
			expectedIsIndex: false,
			expectedIndex:   -1,
			expectedValue:   nil,
		},
		{
			name:            "map of string pointers with struct key",
			value:           map[mapKey]string{{Key: "key"}: "value"},
			expectedKey:     "<generic.mapKey Value>",
			expectedIsIndex: false,
			expectedIndex:   -1,
			expectedValue:   "value",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			iterable, err := NewIterable(test.value)

			if assert.NoError(t, err) {
				assert.Equal(t, 1, iterable.Count())
				assert.False(t, iterable.IsNil())
				i := 0
				err = iterable.Iterate(func(key Key, value interface{}) error {
					assert.Equal(t, test.expectedKey, key.String())
					assert.Equal(t, test.expectedIsIndex, key.IsIndex())
					assert.Equal(t, test.expectedIndex, key.Index())
					assert.Equal(t, test.expectedValue, value)
					i++

					return nil
				})
				assert.NoError(t, err)
				assert.Equal(t, 1, i)
			}
		})
	}
}
