package logic

import (
	"errors"
	"github.com/gogap/isolator/example/todo/models"
	"github.com/gogap/isolator/example/todo/repo"
	"github.com/pborman/uuid"
)

type TaskManager struct {
}

func (p *TaskManager) AddTask(userID, title, description string) (interface{}, error) {

	if len(userID) == 0 {
		return nil, errors.New("user id is empty")
	}

	if len(title) == 0 {
		return nil, errors.New("task title is empty")
	}

	return func(userRepo repo.UserRepository, taskRepo repo.TaskRepository) (taskID string, err error) {
		var exist bool
		if _, exist, err = userRepo.GetUser(userID); err != nil {
			return
		}

		if !exist {
			err = errors.New("user " + userID + " not exist.")
			return
		}

		newTask := models.Task{
			ID:          uuid.NewUUID().String(),
			OwnerID:     userID,
			Title:       title,
			Description: description,
		}

		taskID, err = taskRepo.AddTask(newTask)
		return
	}, nil
}

func (p *TaskManager) GetTask(taskID string) (interface{}, error) {

	if len(taskID) == 0 {
		return nil, errors.New("taskID is empty")
	}

	return func(taskRepo repo.TaskRepository) (task models.Task, err error) {
		var exist bool
		if task, exist, err = taskRepo.GetTaskByID(taskID); err != nil {
			return
		}

		if !exist {
			err = errors.New("task id of " + taskID + " not exist")
			return
		}

		return
	}, nil
}

func (p *TaskManager) GetAllTasks() (interface{}, error) {
	return func(taskRepo repo.TaskRepository) (tasks []models.Task, err error) {
		tasks, err = taskRepo.AllTasks()
		return
	}, nil
}
