package user

import (
	"gofr.dev/pkg/gofr"

	"TaskManager2/models"
)

type Store interface {
	Create(*gofr.Context, *models.User) (int64, error)
	GetByID(*gofr.Context, int64) (*models.User, error)
}
