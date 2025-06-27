package user

import "TaskManager2/models"

type Service interface {
	Create(*models.User) (int64, error)
	GetByID(int64) (*models.User, error)
}
