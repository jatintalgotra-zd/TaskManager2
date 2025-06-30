package task

import (
	"TaskManager2/models"
	"TaskManager2/utils"
)

type MockService struct {
	Check bool
}

func (MockService) Create(t *models.Task) (int64, error) {
	var expID int64 = 999
	if t.ID == expID {
		return expID, utils.ErrTest
	}

	return 1, nil
}

func (m MockService) GetAll() ([]models.Task, error) {
	if m.Check {
		return nil, utils.ErrTest
	}

	tasks := []models.Task{
		{},
	}

	return tasks, nil
}

func (MockService) GetByID(id int64) (*models.Task, error) {
	var expID int64 = 999
	if id == expID {
		return nil, utils.ErrTest
	}

	return &models.Task{ID: 1}, nil
}

func (MockService) Update(t *models.Task) error {
	var expID int64 = 999
	if t.ID == expID {
		return utils.ErrTest
	}

	return nil
}

func (MockService) Delete(id int64) error {
	var expID int64 = 999

	if id == expID {
		return utils.ErrTest
	}

	return nil
}
