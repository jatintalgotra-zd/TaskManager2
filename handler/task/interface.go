package task

import "TaskManager2/models"

type Service interface {
	Create(*models.Task) (int64, error)
	GetAll() ([]models.Task, error)
	GetByID(int64) (*models.Task, error)
	Update(*models.Task) error
	Delete(int64) error
}
