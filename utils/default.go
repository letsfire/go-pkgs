package utils

import (
	"reflect"
)

func SetDefault(src interface{}, def interface{}) {
	rv := reflect.ValueOf(src).Elem()
	if rv.CanAddr() == false {
		panic("the src must addressable")
	} else if rv.IsZero() {
		rv.Set(reflect.ValueOf(def))
	}
}

func FillSlice(ss interface{}, num int, value interface{}) {
	rv := reflect.ValueOf(ss).Elem()
	if rv.CanAddr() == false {
		panic("the slice must addressable")
	} else if n := num - rv.Len(); n > 0 {
		for i := 0; i < n; i++ {
			rv.Set(reflect.Append(rv, reflect.ValueOf(value)))
		}
	}
}
