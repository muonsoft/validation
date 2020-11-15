package generic

import (
	"fmt"
	"reflect"
)

type ErrNotNumeric struct {
	value reflect.Value
}

func (err ErrNotNumeric) Error() string {
	return fmt.Sprintf("value of type %v is not numeric", err.value.Kind())
}
