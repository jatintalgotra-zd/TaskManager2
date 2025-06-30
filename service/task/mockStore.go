package task

import (
	"TaskManager2/models"
	"TaskManager2/utils"
)

type MockStore struct {
	check bool
}

func (MockStore) Create(t *models.Task) (int64, error) {
	var expID int64 = 11
	if t.ID == expID {
		return 0, utils.ErrTest
	}

	return 1, nil
}

func (m MockStore) GetAll() ([]models.Task, error) {
	if m.check {
		return nil, utils.ErrTest
	}

	tasks := []models.Task{
		{},
	}

	return tasks, nil
}

func (MockStore) GetByID(id int64) (*models.Task, error) {
	var expID int64 = 100
	if id == expID {
		return nil, utils.ErrTest
	}

	return &models.Task{ID: id, Desc: "Test task"}, nil
}

func (MockStore) Update(task *models.Task) error {
	var expID int64 = 101
	if task.ID == expID {
		return utils.ErrTest
	}

	return nil
}

func (MockStore) Delete(id int64) error {
	var expID int64 = 102

	if id == expID {
		return utils.ErrTest
	}

	return nil
}
