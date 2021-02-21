package generic

import (
	"fmt"
	"reflect"
	"strconv"
)

type Key interface {
	IsIndex() bool
	fmt.Stringer
}

type Iterable interface {
	Iterate(next func(key Key, value interface{}))
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

func (k intKey) String() string {
	return strconv.Itoa(int(k))
}

type stringKey string

func (s stringKey) IsIndex() bool {
	return false
}

func (s stringKey) String() string {
	return string(s)
}

type iterableArray struct {
	value     reflect.Value
	valueType reflect.Type
}

func (iterable *iterableArray) Iterate(next func(key Key, value interface{})) {
	for i := 0; i < iterable.value.Len(); i++ {
		v := iterable.value.Index(i)
		next(intKey(i), getInterface(v))
	}
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

func (iterable *iterableMap) Iterate(next func(key Key, value interface{})) {
	for _, k := range iterable.value.MapKeys() {
		key := getString(k)
		value := getInterface(iterable.value.MapIndex(k))
		next(stringKey(key), value)
	}
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
