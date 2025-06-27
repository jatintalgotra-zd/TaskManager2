package user

import "TaskManager2/models"

type Store interface {
	Create(*models.User) (int64, error)
	GetByID(int64) (*models.User, error)
}
