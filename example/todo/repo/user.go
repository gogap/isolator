package repo

import (
	"errors"
	"github.com/gogap/isolator"
	"github.com/gogap/isolator/example/todo/models"
	"github.com/gogap/isolator/extension/xorm"
)

type UserRepository interface {
	isolator.Object

	AddUser(user models.User) (err error)
	GetUser(id string) (user models.User, exist bool, err error)
}

type UserRepo struct {
	session *isolator.Session
}

func NewUserRepo() UserRepository {
	return &UserRepo{}
}

func (p *UserRepo) Derive(session *isolator.Session) (obj isolator.Object, err error) {
	return &UserRepo{
		session: session,
	}, nil
}

func (p *UserRepo) AddUser(user models.User) (err error) {
	xormSession := xorm.GetXORMSession(p.session)
	if xormSession == nil {
		err = errors.New("xorm session is nil")
		return
	}

	if _, err = xormSession.InsertOne(&user); err != nil {
		return
	}

	return
}

func (p *UserRepo) GetUser(id string) (user models.User, exist bool, err error) {
	xormSession := xorm.GetXORMSession(p.session)
	if xormSession == nil {
		err = errors.New("xorm session is nil")
		return
	}

	u := models.User{ID: id}

	if exist, err = xormSession.Get(&u); err != nil {
		return
	}

	if exist {
		user = u
	}

	return
}
