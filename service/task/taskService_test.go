package task

import (
	"testing"

	"TaskManager2/models"
)

func TestService_Create(t *testing.T) {
	mockStore := &MockStore{}
	mockUserSvc := &MockUserService{}
	taskService := New(mockStore, mockUserSvc)

	// testcase 1 - success
	id, err := taskService.Create(&models.Task{})
	if err != nil {
		t.Error(err)
		return
	}

	if id != 1 {
		t.Errorf("Expected id 1, got %d", id)
		return
	}

	// testcase 2 - user not validated
	id, err = taskService.Create(&models.Task{UserID: 10})
	if err == nil {
		t.Error("Expected error")
		return
	}

	if id != 0 {
		t.Errorf("Expected id 0, got %d", id)
	}

	// testcase 3 - store.create error
	id, err = taskService.Create(&models.Task{ID: 11})
	if err == nil {
		t.Error("Expected error")
		return
	}

	if id != 0 {
		t.Errorf("Expected id 0, got %d", id)
	}
}

func TestService_GetAll(t *testing.T) {
	mockStore := &MockStore{}
	mockUserSvc := &MockUserService{}
	taskService := New(mockStore, mockUserSvc)

	// testcase 1 - success
	tasks, err := taskService.GetAll()
	if err != nil {
		t.Error(err)
		return
	}

	if len(tasks) != 1 {
		t.Error("Expected 1 task")
		return
	}

	// testcase 2 - store.GetAll error
	mockStore.check = true

	_, err = taskService.GetAll()
	if err == nil {
		t.Error("Expected error")
		return
	}
}

func TestService_GetByID(t *testing.T) {
	mockStore := &MockStore{}
	mockUserSvc := &MockUserService{}
	taskService := New(mockStore, mockUserSvc)

	// testcase 1 - success
	task, err := taskService.GetByID(1)
	if err != nil {
		t.Error(err)
		return
	}

	if task == nil || task.ID != 1 {
		t.Errorf("Expected task with ID 1, got %+v", task)
	}

	// testcase 2 - store.GetByID error
	task, err = taskService.GetByID(100)
	if err == nil {
		t.Error("Expected error")
		return
	}

	if task != nil {
		t.Errorf("Expected nil task on error, got %+v", task)
	}
}

func TestService_Update(t *testing.T) {
	mockStore := &MockStore{}
	mockUserSvc := &MockUserService{}
	taskService := New(mockStore, mockUserSvc)

	// testcase 1 - success
	err := taskService.Update(&models.Task{ID: 1})
	if err != nil {
		t.Error(err)
		return
	}

	// testcase 2 - store.Update error
	err = taskService.Update(&models.Task{ID: 101})
	if err == nil {
		t.Error("Expected error")
	}
}

func TestService_Delete(t *testing.T) {
	mockStore := &MockStore{}
	mockUserSvc := &MockUserService{}
	taskService := New(mockStore, mockUserSvc)

	// testcase 1 - success
	err := taskService.Delete(1)
	if err != nil {
		t.Error(err)
		return
	}

	// testcase 2 - store.Delete error
	err = taskService.Delete(102)
	if err == nil {
		t.Error("Expected error")
	}
}
