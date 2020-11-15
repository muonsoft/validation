package generic

import (
	"reflect"
)

type Number struct {
	int   int64
	float float64
	isInt bool
	isNil bool
}

func NewNumber(value interface{}) (*Number, error) {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return &Number{isNil: true}, nil
		}

		return numberFromValue(v.Elem())
	}

	return numberFromValue(v)
}

func MustNewNumber(value interface{}) Number {
	n, err := NewNumber(value)
	if err != nil {
		panic(err)
	}
	return *n
}

func numberFromValue(v reflect.Value) (*Number, error) {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i := v.Int()
		return &Number{int: i, float: float64(i), isInt: true}, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i := int64(v.Uint())
		return &Number{int: i, float: float64(i), isInt: true}, nil
	case reflect.Float32, reflect.Float64:
		f := v.Float()
		i := int64(f)
		return &Number{int: i, float: f, isInt: f == float64(i)}, nil
	}

	return nil, ErrNotNumeric{value: v}
}

func (n Number) IsNil() bool {
	return n.isNil
}

func (n Number) IsZero() bool {
	if n.isNil {
		return false
	}
	if n.isInt {
		return n.int == 0
	}
	return n.float == 0
}

func (n Number) IsEqualTo(v Number) bool {
	if n.isNil || v.isNil {
		return false
	}
	if n.isInt && v.isInt {
		return n.int == v.int
	}

	return n.float == v.float
}

func (n Number) IsGreaterThan(v Number) bool {
	if n.isNil || v.isNil {
		return false
	}
	if n.isInt && v.isInt {
		return n.int > v.int
	}

	return n.float > v.float
}

func (n Number) IsLessThan(v Number) bool {
	if n.isNil || v.isNil {
		return false
	}
	if n.isInt && v.isInt {
		return n.int < v.int
	}

	return n.float < v.float
}
