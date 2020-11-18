package generic

import (
	"fmt"
	"reflect"
)

type NotNumericError struct {
	value reflect.Value
}

func (err NotNumericError) Error() string {
	return fmt.Sprintf("value of type %v is not numeric", err.value.Kind())
}

type NotIterableError struct {
	value reflect.Value
}

func (err NotIterableError) Error() string {
	return fmt.Sprintf("value of type %v is not iterable", err.value.Kind())
}
