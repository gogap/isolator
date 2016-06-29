package xorm

import (
	"github.com/go-xorm/xorm"
	"github.com/gogap/isolator"
	"golang.org/x/net/context"
)

type xormSessionKey struct{ EngineName string }
type xormSessionIsTransKey struct{ EngineName string }
type xormSessionsKey struct{}

type XORMEngines map[string]*xorm.Engine

func (p XORMEngines) NewXORMSession(engineName string, isTransaction bool) isolator.SessionOption {
	return func(s *isolator.Session) {
		if engine, exist := p[engineName]; exist {

			slist := getXORMSessionList(s)
			exist := false
			for _, name := range slist {
				if engineName == name {
					exist = true
					break
				}
			}

			if exist {
				return
			}

			session := engine.NewSession()
			if isTransaction {
				session.Begin()
			}

			slist = append(slist, engineName)

			s.Context = context.WithValue(s.Context, xormSessionKey{EngineName: engineName}, session)
			s.Context = context.WithValue(s.Context, xormSessionIsTransKey{EngineName: engineName}, isTransaction)
			s.Context = context.WithValue(s.Context, xormSessionsKey{}, slist)
		}
	}
}

func GetXORMSession(s *isolator.Session, engineName string) *xorm.Session {
	if s == nil ||
		s.Context == nil {
		return nil
	}

	v := s.Context.Value(xormSessionKey{EngineName: engineName})

	if session, ok := v.(*xorm.Session); ok {
		return session
	}

	return nil
}

func GetAllXORMSessions(s *isolator.Session, engineName string) map[string]*xorm.Session {
	if s == nil ||
		s.Context == nil {
		return nil
	}

	slist := getXORMSessionList(s)

	if len(slist) == 0 {
		return nil
	}

	engines := map[string]*xorm.Session{}

	for _, sName := range slist {
		session := GetXORMSession(s, sName)
		engines[sName] = session
	}

	return engines
}

func GetXORMSessions(s *isolator.Session, names ...string) map[string]*xorm.Session {
	if s == nil ||
		s.Context == nil ||
		names == nil {
		return nil
	}

	engines := map[string]*xorm.Session{}

	for _, sName := range names {
		session := GetXORMSession(s, sName)
		engines[sName] = session
	}

	return engines
}

func IsXORMTransaction(s *isolator.Session, engineName string) bool {
	if s == nil ||
		s.Context == nil {
		return false
	}

	v := s.Context.Value(xormSessionIsTransKey{EngineName: engineName})

	if isTrans, ok := v.(bool); ok {
		return isTrans
	}

	return false
}

func OnXORMSessionSuccess(session *isolator.Session) {

	if session == nil ||
		session.Context == nil {
		return
	}

	slist := getXORMSessionList(session)

	for _, name := range slist {
		if !IsXORMTransaction(session, name) {
			continue
		}

		xormSession := GetXORMSession(session, name)
		if xormSession == nil {
			continue
		}

		xormSession.Commit()
	}
}

func OnXORMSessionError(session *isolator.Session, err error) {
	if session == nil ||
		session.Context == nil ||
		err == nil {
		return
	}

	slist := getXORMSessionList(session)

	for _, name := range slist {
		if !IsXORMTransaction(session, name) {
			continue
		}

		xormSession := GetXORMSession(session, name)
		if xormSession == nil {
			continue
		}

		xormSession.Rollback()
	}
}

func getXORMSessionList(s *isolator.Session) []string {
	if s == nil {
		return nil
	}

	v := s.Context.Value(xormSessionsKey{})

	if list, ok := v.([]string); ok {
		return list
	}

	return nil
}
