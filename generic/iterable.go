package generic

import (
	"fmt"
	"reflect"
	"strconv"
)

type Key interface {
	IsIndex() bool
	Index() int
	fmt.Stringer
}

// NextElementFunc used to iterate over iterable via closure. To break cycle closure function needs
// to return error.
type NextElementFunc func(key Key, value interface{}) error

type Iterable interface {
	Iterate(next NextElementFunc) error
	Count() int
	IsNil() bool
	IsElementImplements(reflectType reflect.Type) bool
}

func NewIterable(value interface{}) (Iterable, error) {
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Map:
		return &iterableMap{value: v, valueType: reflect.TypeOf(value)}, nil
	case reflect.Slice, reflect.Array:
		return &iterableArray{value: v, valueType: reflect.TypeOf(value)}, nil
	}

	return nil, NotIterableError{value: v}
}

type intKey int

func (k intKey) IsIndex() bool {
	return true
}

func (k intKey) Index() int {
	return int(k)
}

func (k intKey) String() string {
	return strconv.Itoa(int(k))
}

type stringKey string

func (s stringKey) IsIndex() bool {
	return false
}

func (s stringKey) Index() int {
	return -1
}

func (s stringKey) String() string {
	return string(s)
}

type iterableArray struct {
	value     reflect.Value
	valueType reflect.Type
}

func (iterable *iterableArray) Iterate(next NextElementFunc) error {
	for i := 0; i < iterable.value.Len(); i++ {
		v := iterable.value.Index(i)
		err := next(intKey(i), getInterface(v))
		if err != nil {
			return err
		}
	}

	return nil
}

func (iterable *iterableArray) Count() int {
	return iterable.value.Len()
}

func (iterable *iterableArray) IsNil() bool {
	return iterable.value.Kind() == reflect.Slice && iterable.value.IsNil()
}

func (iterable *iterableArray) IsElementImplements(reflectType reflect.Type) bool {
	return iterable.valueType.Elem().Implements(reflectType)
}

type iterableMap struct {
	value     reflect.Value
	valueType reflect.Type
}

func (iterable *iterableMap) Iterate(next NextElementFunc) error {
	for _, k := range iterable.value.MapKeys() {
		key := getString(k)
		value := getInterface(iterable.value.MapIndex(k))
		err := next(stringKey(key), value)
		if err != nil {
			return err
		}
	}

	return nil
}

func (iterable *iterableMap) Count() int {
	return iterable.value.Len()
}

func (iterable *iterableMap) IsNil() bool {
	return iterable.value.IsNil()
}

func (iterable *iterableMap) IsElementImplements(reflectType reflect.Type) bool {
	return iterable.valueType.Elem().Implements(reflectType)
}

func getInterface(value reflect.Value) interface{} {
	switch value.Kind() {
	case reflect.Ptr, reflect.Interface:
		if value.IsNil() {
			return nil
		}
		return value.Elem().Interface()
	}

	return value.Interface()
}

func getString(value reflect.Value) string {
	switch value.Kind() {
	case reflect.Ptr, reflect.Interface:
		if value.IsNil() {
			return ""
		}
		return value.Elem().String()
	}

	return value.String()
}
