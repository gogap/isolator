package isolator

import (
	"errors"
	"reflect"
)

type Result []reflect.Value

func (p Result) End(fn interface{}) error {
	return mapTo(fn, p)
}

func mapTo(fn interface{}, values []reflect.Value) (err error) {

	if len(values) > 0 {
		lastV := values[len(values)-1]
		if lastV.IsValid() && !lastV.IsNil() && lastV.Type() == errType {
			err = values[len(values)-1].Interface().(error)
		}
	}

	if fn != nil {
		fnType := reflect.TypeOf(fn)
		if fnType.Kind() != reflect.Func {
			panic(errors.New("result mapper func should be func"))
		}

		if fnType.NumIn() != len(values) {

			if len(values) != 1 {
				panic(errors.New("mapper func args number error"))
			}

			newV := []reflect.Value{}
			for i := 0; i < fnType.NumIn()-1; i++ {
				newV = append(newV, reflect.Zero(fnType.In(i)))
			}

			newV = append(newV, values[0])
			values = newV
		}

		fnValue := reflect.ValueOf(fn)
		fnValue.Call(values)
	}

	return
}
