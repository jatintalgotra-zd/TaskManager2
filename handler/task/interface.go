package task

import (
	"gofr.dev/pkg/gofr"

	"TaskManager2/models"
)

type Service interface {
	Create(*gofr.Context, *models.Task) (int64, error)
	GetAll(*gofr.Context) ([]models.Task, error)
	GetByID(*gofr.Context, int64) (*models.Task, error)
	Update(*gofr.Context, *models.Task) error
	Delete(*gofr.Context, int64) error
}
