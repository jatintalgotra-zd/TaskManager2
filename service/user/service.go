package user

import (
	"gofr.dev/pkg/gofr"

	"TaskManager2/models"
)

type service struct {
	store Store
}

func New(store Store) *service {
	return &service{store: store}
}

func (s *service) Create(ctx *gofr.Context, user *models.User) (int64, error) {
	id, err := s.store.Create(ctx, user)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *service) GetByID(ctx *gofr.Context, id int64) (*models.User, error) {
	user, err := s.store.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}
