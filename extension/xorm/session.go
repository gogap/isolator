package xorm

import (
	"github.com/go-xorm/xorm"
	"github.com/gogap/isolator"
	"golang.org/x/net/context"
)

type xormSessionKey struct{}
type xormSessionIsTransKey struct{}

type XORMEngines map[string]*xorm.Engine

func (p XORMEngines) NewXORMSession(engineName string, isTransaction bool) isolator.SessionOption {
	return func(s *isolator.Session) {
		if engine, exist := p[engineName]; exist {
			session := engine.NewSession()
			if isTransaction {
				session.Begin()
			}
			s.Context = context.WithValue(s.Context, xormSessionKey{}, session)
			s.Context = context.WithValue(s.Context, xormSessionIsTransKey{}, isTransaction)
		}
	}
}

func GetXORMSession(s *isolator.Session) *xorm.Session {
	if s == nil ||
		s.Context == nil {
		return nil
	}

	v := s.Context.Value(xormSessionKey{})

	if session, ok := v.(*xorm.Session); ok {
		return session
	}

	return nil
}

func IsXORMTransaction(s *isolator.Session) bool {
	if s == nil ||
		s.Context == nil {
		return false
	}

	v := s.Context.Value(xormSessionIsTransKey{})

	if isTrans, ok := v.(bool); ok {
		return isTrans
	}

	return false
}

func OnXORMSessionSuccess(session *isolator.Session) {
	if !IsXORMTransaction(session) {
		return
	}

	if session == nil ||
		session.Context == nil {
		return
	}

	xormSession := GetXORMSession(session)
	if xormSession == nil {
		return
	}

	xormSession.Commit()
}

func OnXORMSessionError(session *isolator.Session, err error) {
	if !IsXORMTransaction(session) {
		return
	}

	if session == nil ||
		session.Context == nil ||
		err == nil {
		return
	}

	xormSession := GetXORMSession(session)
	if xormSession == nil {
		return
	}

	xormSession.Rollback()
}
