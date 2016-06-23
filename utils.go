package isolator

import (
	"errors"
	"reflect"
)

func getStructName(v interface{}) (name string, err error) {
	typ := reflect.TypeOf(v)
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	if typ.Kind() != reflect.Struct {
		err = errors.New("object is not a struct")
		return
	}

	name = typ.String()

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
