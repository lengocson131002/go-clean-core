package util

import (
	"fmt"
	"reflect"
)

const (
	POINTER_PREFIX = ""
)

func GetType(obj interface{}) string {
	if t := reflect.TypeOf(obj); t.Kind() == reflect.Ptr {
		return fmt.Sprintf("%s%s", POINTER_PREFIX, t.Elem().Name())
	} else {
		return t.Name()
	}
}
