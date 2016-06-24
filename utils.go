package isolator

import (
	"errors"
	"reflect"
)

func getStructName(v interface{}) (name string, err error) {
	typ := finnalType(v)

	if typ.Kind() != reflect.Struct {
		err = errors.New("object is not a struct")
		return
	}

	name = typ.String()

	return
}

func finnalType(v interface{}) (typ reflect.Type) {
	t := reflect.TypeOf(v)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	typ = t
	return
}

func lastSecondType(v interface{}) (typ reflect.Type) {
	t := reflect.TypeOf(v)
	for t.Kind() == reflect.Ptr {
		if t.Elem().Kind() == reflect.Ptr {
			t = t.Elem()
		} else {
			typ = t
			return
		}
	}
	typ = t
	return
}

func objsToReflectValues(objs ...interface{}) []reflect.Value {
	if objs == nil {
		return []reflect.Value{}
	}

	values := []reflect.Value{}

	for _, arg := range objs {
		v := reflect.ValueOf(arg)
		values = append(values, v)
	}

	return values
}

func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Func, reflect.Map, reflect.Slice:
		return v.IsNil()
	case reflect.Array:
		z := true
		for i := 0; i < v.Len(); i++ {
			z = z && isZero(v.Index(i))
		}
		return z
	case reflect.Struct:
		z := true
		for i := 0; i < v.NumField(); i++ {
			z = z && isZero(v.Field(i))
		}
		return z
	}

	z := reflect.Zero(v.Type())
	return v.Interface() == z.Interface()
}
