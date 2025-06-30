package user

import (
	"testing"

	"TaskManager2/models"
)

func TestCreate(t *testing.T) {
	mockStore := MockStore{}
	userService := New(&mockStore)

	// case 1 - pass
	user := &models.User{}

	id, err := userService.Create(user)
	if err != nil {
		t.Errorf("Error creating user: %v", err)
		return
	}

	if id != 1 {
		t.Errorf("Expected id 1, got %v", id)
		return
	}

	// case 2 - fail
	user2 := &models.User{Name: "test", Email: "test@user.com"}

	_, err2 := userService.Create(user2)
	if err2 == nil {
		t.Errorf("Test should have failed but did not")
		return
	}
}

func TestService_GetByID(t *testing.T) {
	mockStore := MockStore{}
	userService := New(&mockStore)

	// case 1 - pass
	user, err := userService.GetByID(1)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if user == nil {
		t.Error("user is nil")
		return
	}

	if user.ID != 1 {
		t.Errorf("Expected: 1, Got: %d", user.ID)
		return
	}

	if user.Name != "test" {
		t.Errorf("Expected: %s, Got: %s", "test", user.Name)
		return
	}

	// case 2 - fail
	_, err2 := userService.GetByID(10)
	if err2 == nil {
		t.Errorf("Test should have failed but did not")
	}
}
