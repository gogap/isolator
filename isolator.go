package isolator

import (
	"errors"
	"reflect"
)

var (
	errType = reflect.TypeOf((*error)(nil)).Elem()
)

type IsolatorOption func(*Isolator)

type Isolator struct {
	sessions      *Sessions
	ObjectBuilder ObjectBuilder
}

// Invoke is a func for reflect Call logic func
// fn:= func(arg1, arg2, arg3 string) (interface{}, error) {
//		if(arg1!="xxxx") {
//			return nil, error.New("arg1 error")
//		}
// 		return func(obj1 Repo1, obj2 Repo2) (result string, err error) {
// 			......
// 			return "good"
// 		}, nil
// }
//
// Repo1, Repo2 is build by object builder
func NewIsolator(opts ...IsolatorOption) *Isolator {
	ios := &Isolator{sessions: NewSessions(), ObjectBuilder: DefaultObjectBuilder}
	ios.Options(opts...)
	return ios
}

func IsolatorObjectBuilder(builder ObjectBuilder) IsolatorOption {
	return func(i *Isolator) {
		if builder == nil {
			i.ObjectBuilder = DefaultObjectBuilder
			return
		}
		i.ObjectBuilder = builder
	}
}

func (p *Isolator) Options(opts ...IsolatorOption) {
	if opts == nil {
		return
	}

	for _, opt := range opts {
		opt(p)
	}
}

func (p *Isolator) Invoke(fn interface{}, args ...interface{}) (ret Result) {
	var err error
	var session *Session

	defer func() {
		if err != nil {
			if session != nil && session.OnError != nil {
				session.OnError(session, err)
			}

			ret = []reflect.Value{reflect.ValueOf(err)}
		}
	}()

	if fn == nil {
		err = errors.New("invoke fn is nil")
		return
	}

	typ := reflect.TypeOf(fn)
	if typ.Kind() != reflect.Func {
		err = errors.New("invoke fn is not a func")
		return
	}

	fnVal := reflect.ValueOf(fn)
	argVals := objsToReflectValues(args...)

	fnRets := fnVal.Call(argVals)

	if len(fnRets) != 2 {
		err = errors.New("invoke fn return value length error, method:" + fnVal.String())
		return
	}

	fnLastV := fnRets[1]
	if fnLastV.IsValid() && !fnLastV.IsNil() && fnLastV.Type().ConvertibleTo(errType) {
		err = fnLastV.Interface().(error)
		return
	}

	logicFn := fnRets[0].Elem()

	if logicFn.Kind() != reflect.Func {
		err = errors.New("invoke fn return value is not a func, method:" + fnVal.String())
		return
	}

	typeLogicFn := logicFn.Type()
	lenLogicArgs := typeLogicFn.NumIn()
	logicFnArgTypes := make([]reflect.Type, lenLogicArgs)

	for i := 0; i < typeLogicFn.NumIn(); i++ {
		logicFnArgTypes[i] = typeLogicFn.In(i)
	}

	session = p.sessions.New(logicFnArgTypes...)

	var logicFnArgs []Object
	if logicFnArgs, err = p.ObjectBuilder.DeriveObjects(session, logicFnArgTypes...); err != nil {
		return
	}

	var logicFnArgsI = make([]interface{}, lenLogicArgs)
	for i := 0; i < len(logicFnArgs); i++ {
		logicFnArgsI[i] = logicFnArgs[i]
	}

	logicFnArgVals := objsToReflectValues(logicFnArgsI...)

	ret = logicFn.Call(logicFnArgVals)

	errorResult := false
	if len(ret) > 0 {
		lastV := ret[len(ret)-1]
		if lastV.IsValid() && !lastV.IsNil() && lastV.Type().ConvertibleTo(errType) {
			errorResult = true
			if session.OnError != nil {
				session.OnError(session, lastV.Interface().(error))
			}
		}
	}

	if !errorResult {
		if session.OnSuccess != nil {
			session.OnSuccess(session)
		}
	}

	return
}

func (p *Isolator) ObjectSessionOptions(obj Object, opts ...SessionOption) {
	p.sessions.RegisterObjectOptions(obj, opts...)
	return
}

func (p *Isolator) ObjectsSessionOptions(objs []Object, opts ...SessionOption) {
	if objs == nil || opts == nil {
		return
	}

	for i := 0; i < len(objs); i++ {
		p.sessions.RegisterObjectOptions(objs[i], opts...)
	}

	return
}
