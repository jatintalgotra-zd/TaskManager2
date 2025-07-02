package task

import (
	"errors"
	"testing"

	"go.uber.org/mock/gomock"
	"gofr.dev/pkg/gofr"

	"TaskManager2/models"
	"TaskManager2/utils"
)

func TestService_Create(t *testing.T) {
	var ctx *gofr.Context

	controller := gomock.NewController(t)
	mockStore := NewMockStore(controller)
	mockUserSvc := NewMockUserService(controller)
	taskService := New(mockStore, mockUserSvc)

	tests := []struct {
		description string
		input       *models.Task
		expectedID  int64
		expectedErr error
	}{
		{"success", &models.Task{}, 0, nil},
		{"user not validated", &models.Task{UserID: 10}, 0, utils.ErrTest},
		{"create error", &models.Task{UserID: 11}, 0, utils.ErrTest},
	}

	for _, tc := range tests {
		if tc.input.UserID != 10 {
			mockStore.EXPECT().Create(ctx, tc.input).Return(tc.expectedID, tc.expectedErr)
			mockUserSvc.EXPECT().GetByID(ctx, tc.input.UserID).Return(&models.User{}, nil)
		}

		if tc.input.UserID == 10 {
			mockUserSvc.EXPECT().GetByID(ctx, tc.input.UserID).Return(nil, tc.expectedErr)
		}

		id, err := taskService.Create(ctx, tc.input)
		if !errors.Is(err, tc.expectedErr) {
			t.Errorf("expected error %s, got %s", tc.expectedErr, err)
		}

		if id != tc.expectedID {
			t.Errorf("Expected id %d, got %d", tc.expectedID, id)
		}
	}
}

func TestService_GetAll(t *testing.T) {
	var ctx *gofr.Context

	controller := gomock.NewController(t)
	mockStore := NewMockStore(controller)
	mockUserSvc := NewMockUserService(controller)
	taskService := New(mockStore, mockUserSvc)

	testcases := []struct {
		description   string
		expected      []models.Task
		expectedError error
	}{
		{"success", []models.Task{{}}, nil},
		{"user not validated", nil, utils.ErrTest},
	}

	for _, test := range testcases {
		mockStore.EXPECT().GetAll(ctx).Return(test.expected, test.expectedError)

		_, err := taskService.GetAll(ctx)
		if !errors.Is(err, test.expectedError) {
			t.Errorf("Test Failed: (%s) Expected: (%s) Actual: (%s)", test.description, test.expectedError, err)
		}
	}
}

func TestService_GetByID(t *testing.T) {
	var ctx *gofr.Context

	controller := gomock.NewController(t)
	mockStore := NewMockStore(controller)
	mockUserSvc := NewMockUserService(controller)
	taskService := New(mockStore, mockUserSvc)

	testcases := []struct {
		description   string
		input         int64
		expected      *models.Task
		expectedError error
	}{
		{"success", 1, &models.Task{}, nil},
		{"store GetByID method error", 1, nil, utils.ErrTest},
	}

	for _, tc := range testcases {
		mockStore.EXPECT().GetByID(ctx, tc.input).Return(tc.expected, tc.expectedError)

		_, err := taskService.GetByID(ctx, tc.input)
		if !errors.Is(err, tc.expectedError) {
			t.Errorf("Expected error: %s, got %s", tc.expectedError, err)
		}
	}
}

func TestService_Update(t *testing.T) {
	var ctx *gofr.Context

	controller := gomock.NewController(t)
	mockStore := NewMockStore(controller)
	mockUserSvc := NewMockUserService(controller)
	taskService := New(mockStore, mockUserSvc)

	testcases := []struct {
		description   string
		input         *models.Task
		expectedError error
	}{
		{"success", &models.Task{}, nil},
		{"store Update method error", nil, utils.ErrTest},
	}

	for _, tc := range testcases {
		mockStore.EXPECT().Update(ctx, tc.input).Return(tc.expectedError)

		err := taskService.Update(ctx, tc.input)
		if !errors.Is(err, tc.expectedError) {
			t.Errorf("Expected error: %s, got %s", tc.expectedError, err)
		}
	}
}

func TestService_Delete(t *testing.T) {
	var ctx *gofr.Context

	controller := gomock.NewController(t)
	mockStore := NewMockStore(controller)
	mockUserSvc := NewMockUserService(controller)
	taskService := New(mockStore, mockUserSvc)

	testcases := []struct {
		description   string
		input         int64
		expectedError error
	}{
		{"success", 1, nil},
		{"store Delete method error", 0, utils.ErrTest},
	}

	for _, tc := range testcases {
		mockStore.EXPECT().Delete(ctx, tc.input).Return(tc.expectedError)

		err := taskService.Delete(ctx, tc.input)
		if !errors.Is(err, tc.expectedError) {
			t.Errorf("Expected error: %s, got %s", tc.expectedError, err)
		}
	}
}
