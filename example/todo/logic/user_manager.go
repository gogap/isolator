package logic

import (
	"errors"
	"github.com/gogap/isolator/example/todo/models"
	"github.com/gogap/isolator/example/todo/repo"
)

type UserManager struct {
}

func (p *UserManager) AddUser(id, name string) (interface{}, error) {

	if len(id) == 0 {
		return nil, errors.New("user id is empty")
	}

	if len(name) == 0 {
		return nil, errors.New("user name is empty")
	}

	return func(userRepo repo.UserRepository) (err error) {

		var exist bool

		if _, exist, err = userRepo.GetUser(id); err != nil {
			return
		}

		if exist {
			err = errors.New("user of " + id + " already exist")
			return
		}

		newUser := models.User{
			ID:   id,
			Name: name,
		}

		err = userRepo.AddUser(newUser)
		return
	}, nil
}

func (p *UserManager) GetUser(userID string) (interface{}, error) {

	if len(userID) == 0 {
		return nil, errors.New("userID is empty")
	}

	return func(userRepo repo.UserRepository) (user models.User, err error) {
		var exist bool
		if user, exist, err = userRepo.GetUser(userID); err != nil {
			return
		}

		if !exist {
			err = errors.New("user " + userID + " not exist")
			return
		}
		return
	}, nil
}
