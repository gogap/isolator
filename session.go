package isolator

import (
	"github.com/pborman/uuid"
	"golang.org/x/net/context"
	"reflect"
	"sync"
	"time"
)

type SessionOption func(*Session)

type SessionOnErrorFunc func(*Session, error)
type SessionOnSuccessFunc func(*Session)

type Session struct {
	ID         string
	Context    context.Context
	CreateTime int64
	OnError    SessionOnErrorFunc
	OnSuccess  SessionOnSuccessFunc
}

func NewSession(opts ...SessionOption) *Session {
	s := &Session{
		Context:    context.TODO(),
		ID:         uuid.NewUUID().String(),
		CreateTime: time.Now().UnixNano(),
	}

	return s
}

func SessionOnError(fn SessionOnErrorFunc) SessionOption {
	return func(s *Session) {
		s.OnError = fn
	}
}

func SessionOnSuccess(fn SessionOnSuccessFunc) SessionOption {
	return func(s *Session) {
		s.OnSuccess = fn
	}
}

func (p *Session) Options(opts ...SessionOption) {
	if opts == nil {
		return
	}

	for _, o := range opts {
		o(p)
	}

	return
}

type Sessions struct {
	locker      sync.Mutex
	options     map[string][]SessionOption
	objectTypes map[string]reflect.Type
}

func NewSessions() *Sessions {
	return &Sessions{
		options:     make(map[string][]SessionOption),
		objectTypes: make(map[string]reflect.Type),
	}
}

func (p *Sessions) RegisterObjectOptions(obj Object, opts ...SessionOption) {

	objName, err := getStructName(obj)
	if err != nil {
		return
	}

	p.locker.Lock()
	defer p.locker.Unlock()
	p.options[objName] = opts
	p.objectTypes[objName] = lastSecondType(obj)

	return
}

func (p *Sessions) New(types ...reflect.Type) (session *Session) {
	session = NewSession()
	if types == nil {
		return
	}

	globalOptsMap := map[interface{}]bool{}
	globalOpts := []SessionOption{}

	for _, typ := range types {
		deepType := typ
		for deepType.Kind() == reflect.Ptr {
			deepType = deepType.Elem()
		}

		objTypeName := typ.String()

		if deepType.Kind() == reflect.Interface {
			for _, oType := range p.objectTypes {
				if oType.ConvertibleTo(typ) {
					for oType.Kind() == reflect.Ptr {
						oType = oType.Elem()
					}
					objTypeName = oType.String()
					break
				}
			}
		}

		if opts, exist := p.options[objTypeName]; exist {
			for i := 0; i < len(opts); i++ {
				if _, exist := globalOptsMap[&opts[i]]; !exist {
					globalOptsMap[&opts[i]] = true
					globalOpts = append(globalOpts, opts[i])
				}
			}
		}
	}

	session.Options(globalOpts...)
	return
}
