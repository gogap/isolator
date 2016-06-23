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

type Session struct {
	ID         string
	Context    context.Context
	CreateTime int64
	OnError    SessionOnErrorFunc
}

func NewSession(opts ...SessionOption) *Session {
	s := &Session{
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
	locker  sync.Mutex
	options map[string][]SessionOption
}

func NewSessions() *Sessions {
	return &Sessions{
		options: make(map[string][]SessionOption),
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
		for typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		}

		if opts, exist := p.options[typ.String()]; exist {
			for _, opt := range opts {
				if _, exist := globalOptsMap[opt]; !exist {
					globalOptsMap[&opt] = true
					globalOpts = append(globalOpts, opt)
				}
			}
		}
	}

	session.Options(globalOpts...)
	return
}
