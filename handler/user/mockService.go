package user

import (
	"TaskManager2/models"
	"TaskManager2/utils"
)

type MockService struct {
	Check bool
}

func (MockService) Create(u *models.User) (int64, error) {
	var expID int64 = 999
	if u.ID == expID {
		return expID, utils.ErrTest
	}

	return 1, nil
}

func (MockService) GetByID(id int64) (*models.User, error) {
	var expID int64 = 999
	if id == expID {
		return nil, utils.ErrTest
	}

	return &models.User{ID: 1}, nil
}
