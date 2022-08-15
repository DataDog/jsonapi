package jsonapi

import (
	"fmt"
	"reflect"
	"runtime/debug"
)

func derefValue(v reflect.Value) reflect.Value {
	switch v.Kind() {
	case reflect.Pointer:
		if v.Elem().Kind() == reflect.Invalid {
			v.Set(reflect.New(v.Type().Elem()))
			return v.Elem()
		}
		return derefValue(v.Elem())
	}
	return v
}

func derefType(t reflect.Type) reflect.Type {
	switch t.Kind() {
	case reflect.Pointer:
		return derefType(t.Elem())
	}
	return t
}

func recoverError(rvr any) error {
	var err error
	switch e := rvr.(type) {
	case error:
		err = fmt.Errorf("unknown error: %w %s", e, debug.Stack())
	default:
		err = fmt.Errorf("%v %s", e, debug.Stack())
	}
	return err
}
