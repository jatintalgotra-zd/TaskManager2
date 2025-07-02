package task

import (
	"gofr.dev/pkg/gofr"

	"TaskManager2/models"
)

type service struct {
	store       Store
	userService UserService
}

func New(store Store, userSvc UserService) *service {
	return &service{store: store, userService: userSvc}
}

func (s *service) Create(ctx *gofr.Context, task *models.Task) (int64, error) {
	// validate if user exists
	_, err := s.userService.GetByID(ctx, task.UserID)
	if err != nil {
		return 0, err
	}

	id, err2 := s.store.Create(ctx, task)
	if err2 != nil {
		return 0, err2
	}

	return id, nil
}

func (s *service) GetAll(ctx *gofr.Context) ([]models.Task, error) {
	tasks, err := s.store.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *service) GetByID(ctx *gofr.Context, id int64) (*models.Task, error) {
	task, err := s.store.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (s *service) Update(ctx *gofr.Context, task *models.Task) error {
	err := s.store.Update(ctx, task)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) Delete(ctx *gofr.Context, id int64) error {
	err := s.store.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
