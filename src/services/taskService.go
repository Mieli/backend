package services

import (
	pkgmodels "github.com/mieli/backend/src/models"
	pkgrepositories "github.com/mieli/backend/src/repositories"
)

type TaskService struct {
	Repository *pkgrepositories.TaskRepository
}

func NewTaskService(repository *pkgrepositories.TaskRepository) *TaskService {
	return &TaskService{
		Repository: repository,
	}
}

func (s *TaskService) FindAll() ([]pkgmodels.Task, error) {
	tasks, err := s.Repository.FindAll()
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (s *TaskService) FindById(id string) (*pkgmodels.Task, error) {
	task, err := s.Repository.FindById(id)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TaskService) Add(task pkgmodels.Task) error {
	err := s.Repository.Add(task)
	if err != nil {
		return err
	}
	return nil
}

func (s *TaskService) Update(id string, updatedTask pkgmodels.Task) error {
	err := s.Repository.Update(id, updatedTask)
	if err != nil {
		return err
	}
	return nil
}

func (s *TaskService) Remove(id string) error {
	err := s.Repository.Remove(id)
	if err != nil {
		return err
	}
	return nil
}
