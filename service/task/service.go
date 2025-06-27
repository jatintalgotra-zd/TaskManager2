package task

import (
	"TaskManager2/handler/user"
	"TaskManager2/models"
)

type service struct {
	store       Store
	userService user.Service
}

func New(store Store, userSvc user.Service) *service {
	return &service{store: store, userService: userSvc}
}

func (s *service) Create(task *models.Task) (int64, error) {
	// validate if user exists
	_, err := s.userService.GetByID(task.UserID)
	if err != nil {
		return 0, err
	}

	id, err2 := s.store.Create(task)
	if err2 != nil {
		return 0, err2
	}

	return id, nil
}

func (s *service) GetAll() ([]models.Task, error) {
	tasks, err := s.store.GetAll()
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *service) GetByID(id int64) (*models.Task, error) {
	task, err := s.store.GetByID(id)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (s *service) Update(task *models.Task) error {
	err := s.store.Update(task)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) Delete(id int64) error {
	err := s.store.Delete(id)
	if err != nil {
		return err
	}

	return nil
}
