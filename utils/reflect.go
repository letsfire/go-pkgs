package utils

import (
	"reflect"
	"strings"
	"sync"
)

var fieldsCacheLocker = new(sync.RWMutex)
var fieldsCacheMapper = make(map[reflect.Type][]string)

func FieldsNameByTag(v interface{}, tag string, export bool, more ...string) []string {
	rt := reflect.TypeOf(v)
	for rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	fieldsCacheLocker.RLock()
	if v, ok := fieldsCacheMapper[rt]; ok {
		fieldsCacheLocker.RUnlock()
		return append(v, more...)
	}
	fieldsCacheLocker.RUnlock()
	fieldsCacheLocker.Lock()
	defer fieldsCacheLocker.Unlock()
	if v, ok := fieldsCacheMapper[rt]; ok {
		return v
	}
	var names = make([]string, 0)
	for i := 0; i < rt.NumField(); i++ {
		if export && !rt.Field(i).IsExported() {
			continue
		}
		if tagValue, exist := rt.Field(i).Tag.Lookup(tag); exist {
			if name := strings.SplitN(tagValue, ",", 2)[0]; name != "" {
				names = append(names, strings.TrimSpace(name))
			}
		}
	}
	fieldsCacheMapper[rt] = names
	return append(names, more...)
}

func IsZero(v interface{}) bool {
	return reflect.ValueOf(v).IsZero()
}
