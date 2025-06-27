package user

import (
	"TaskManager2/models"
)

type service struct {
	store Store
}

func New(store Store) *service {
	return &service{store: store}
}

func (s *service) Create(user *models.User) (int64, error) {
	id, err := s.store.Create(user)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *service) GetByID(id int64) (*models.User, error) {
	user, err := s.store.GetByID(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}
