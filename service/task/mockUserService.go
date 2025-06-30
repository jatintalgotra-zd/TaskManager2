package task

import (
	"TaskManager2/models"
	"TaskManager2/utils"
)

type MockUserService struct {
	Check int64
}

func (MockUserService) GetByID(id int64) (*models.User, error) {
	var exp int64 = 10
	if id == exp {
		return nil, utils.ErrTest
	}

	return nil, nil
}
