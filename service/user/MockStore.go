package user

import (
	"errors"

	"TaskManager2/models"
)

var errTest = errors.New("test error")

type MockStore struct {
	Check int64
}

func (MockStore) Create(u *models.User) (int64, error) {
	testStr := "test"
	if u.Name == testStr {
		return 0, errTest
	}

	return 1, nil
}

func (MockStore) GetByID(id int64) (*models.User, error) {
	user := &models.User{ID: id, Name: "test", Email: "test@user.com"}

	var expID int64 = 10
	if id == expID {
		return nil, errTest
	}

	return user, nil
}
