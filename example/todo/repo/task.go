package repo

import (
	"errors"
	"github.com/gogap/isolator"
	"github.com/gogap/isolator/example/todo/models"
	"github.com/gogap/isolator/extension/xorm"
)

type TaskRepository interface {
	isolator.Object

	AddTask(task models.Task) (id string, err error)
	GetTaskByID(id string) (task models.Task, exist bool, err error)
	GetUserTasks(userID string) (task []models.Task, err error)
	AllTasks() (task []models.Task, err error)
}

type TaskRepo struct {
	session *isolator.Session
}

func NewTaskRepo() TaskRepository {
	return &TaskRepo{}
}

func (p *TaskRepo) Derive(session *isolator.Session) (obj isolator.Object, err error) {
	return &TaskRepo{
		session: session,
	}, nil
}

func (p *TaskRepo) AddTask(task models.Task) (id string, err error) {
	xormSession := xorm.GetXORMSession(p.session, "todo")
	if xormSession == nil {
		err = errors.New("xorm session is nil")
		return
	}

	if _, err = xormSession.InsertOne(&task); err != nil {
		return
	}

	return
}

func (p *TaskRepo) GetTaskByID(id string) (task models.Task, exist bool, err error) {
	xormSession := xorm.GetXORMSession(p.session, "todo")
	if xormSession == nil {
		err = errors.New("xorm session is nil")
		return
	}

	t := models.Task{ID: id}

	if exist, err = xormSession.Get(&t); err != nil {
		return
	}

	if exist {
		task = t
	}

	return
}

func (p *TaskRepo) GetUserTasks(userID string) (task []models.Task, err error) {
	xormSession := xorm.GetXORMSession(p.session, "todo")
	if xormSession == nil {
		err = errors.New("xorm session is nil")
		return
	}

	if err = xormSession.Where("`owner_id` = ?", userID).Find(&task); err != nil {
		return
	}

	return
}

func (p *TaskRepo) AllTasks() (task []models.Task, err error) {
	xormSession := xorm.GetXORMSession(p.session, "todo")
	if xormSession == nil {
		err = errors.New("xorm session is nil")
		return
	}

	if err = xormSession.Find(&task); err != nil {
		return
	}

	return
}
