package user

import (
	"errors"

	"TaskManager2/models"
)

type MockStore struct {}

func (MockStore) Create(u *models.User) (int64, error) {
	if u.Name == "test"{
		return 0, errors.New("User already exists")
	}

	return 1, nil
}

func (MockStore) GetByID(ID int64) (*models.User, error) {
	user := &models.User{
		ID: ID,
		Name: "test",
		Email: "test@user.com",
	}

	if ID == 10 {
		return nil, errors.New("User not found")
	}

	return user, nil
}
