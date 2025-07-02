package user

import (
	"errors"
	"fmt"
	"testing"

	"go.uber.org/mock/gomock"
	"gofr.dev/pkg/gofr"

	"TaskManager2/models"
	"TaskManager2/utils"
)

func TestCreate(t *testing.T) {
	var ctx *gofr.Context

	controller := gomock.NewController(t)
	mockStore := NewMockStore(controller)
	userService := New(mockStore)

	testCases := []struct {
		description   string
		input         models.User
		expectedID    int64
		expectedError error
	}{
		{"success", models.User{ID: 1, Name: "test1"}, 1, nil},
		{"store create method error", models.User{ID: 1, Name: "test1"}, 0, utils.ErrTest},
	}

	for _, tc := range testCases {
		mockStore.EXPECT().Create(ctx, &tc.input).Return(tc.expectedID, tc.expectedError)

		id, err := userService.Create(ctx, &tc.input)
		if !errors.Is(err, tc.expectedError) {
			t.Errorf("Expected error %v, got %v", tc.expectedError, err)
		}

		if id != tc.expectedID {
			t.Errorf("Create(%+v) expected id = %d, actual = %d", tc.input, tc.expectedID, id)
		}
	}
}

func TestService_GetByID(t *testing.T) {
	var ctx *gofr.Context

	controller := gomock.NewController(t)
	mockStore := NewMockStore(controller)
	userService := New(mockStore)

	testCases := []struct {
		description   string
		input         int64
		expected      *models.User
		expectedError error
	}{
		{"success", 1, &models.User{ID: 1, Name: "test1"}, nil},
		{"store get method error", 10, &models.User{ID: 1, Name: "test2"}, utils.ErrTest},
	}

	for _, tc := range testCases {
		mockStore.EXPECT().GetByID(ctx, tc.input).Return(tc.expected, tc.expectedError)

		user, err := userService.GetByID(ctx, tc.input)
		if !errors.Is(err, tc.expectedError) {
			fmt.Println("check")
			t.Errorf("Expected error %v, got %v", tc.expectedError, err)
		}

		if user != nil && user.ID != tc.expected.ID {
			t.Errorf("Expected: %v, got %v", tc.expected.Name, user.Name)
		}
	}
}
