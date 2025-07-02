package user

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"go.uber.org/mock/gomock"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/container"
	gofrhttp "gofr.dev/pkg/gofr/http"

	"TaskManager2/models"
	"TaskManager2/utils"
)

func TestHandler_PostHandler(t *testing.T) {
	controller := gomock.NewController(t)
	mockSvc := NewMockService(controller)
	taskHandler := New(mockSvc)

	mockContainer, _ := container.NewMockContainer(t)

	ctx := &gofr.Context{
		Context:   t.Context(),
		Request:   nil,
		Container: mockContainer,
	}

	testcases := []struct {
		name             string
		requestBody      string
		mockExpect       func()
		expectedResponse any
		expectedError    error
	}{
		{
			"success",
			`{
							"name" : "test user",
							"email" : "test@email.com"
						}`,
			func() {
				mockSvc.EXPECT().Create(ctx, &models.User{Name: "test user", Email: "test@email.com"}).Return(int64(1), nil)
			},
			int64(1),
			nil,
		},
		{
			"bind error",
			`describe":"test task","status":false,"user_id":1}`,
			func() {},
			nil,
			gofrhttp.ErrorInvalidParam{},
		},
		{
			"service create error",
			`{
							"name" : "test user",
							"email" : "test@email.com"
						}`,
			func() {
				mockSvc.EXPECT().Create(ctx, &models.User{Name: "test user", Email: "test@email.com"}).Return(int64(0), utils.ErrTest)
			},
			nil,
			utils.ErrTest,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockExpect()

			body := bytes.NewReader([]byte(tc.requestBody))
			req := httptest.NewRequest(http.MethodPost, "/user", body)
			req.Header.Set("Content-Type", "application/json")

			ctx.Request = gofrhttp.NewRequest(req)

			id, err := taskHandler.PostHandler(ctx)
			if err != nil && err.Error() != tc.expectedError.Error() {
				t.Errorf("error, expected %v, got %v", tc.expectedError, err)
			}

			if id != tc.expectedResponse {
				t.Errorf("expected: %v, got: %v", tc.expectedResponse, id)
			}
		})
	}
}

func TestHandler_GetByIDHandler(t *testing.T) {
	controller := gomock.NewController(t)
	mockSvc := NewMockService(controller)
	taskHandler := New(mockSvc)

	mockContainer, _ := container.NewMockContainer(t)
	ctx := &gofr.Context{
		Context:   t.Context(),
		Request:   nil,
		Container: mockContainer,
	}

	testcases := []struct {
		name             string
		requestID        string
		mockExpect       func()
		expectedResponse any
		expectedError    error
	}{
		{
			"success",
			"1",
			func() {
				mockSvc.EXPECT().GetByID(ctx, int64(1)).Return(&models.User{ID: 1}, nil)
			},
			&models.User{ID: 1},
			nil,
		},
		{
			"Atoi error",
			"abc",
			func() {},
			nil,
			gofrhttp.ErrorInvalidParam{Params: []string{"abc"}},
		},
		{
			"service GetByID error",
			"1",
			func() {
				mockSvc.EXPECT().GetByID(ctx, int64(1)).Return(nil, utils.ErrTest)
			},
			nil,
			utils.ErrTest,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockExpect()

			req := httptest.NewRequest(http.MethodGet, "/user/{id}", http.NoBody)
			req.Header.Set("Content-Type", "application/json")

			req = mux.SetURLVars(req, map[string]string{"id": tc.requestID})
			ctx.Request = gofrhttp.NewRequest(req)

			task, err := taskHandler.GetByIDHandler(ctx)
			if err != nil && err.Error() != tc.expectedError.Error() {
				t.Errorf("error, expected %v, got %v", tc.expectedError, err)
			}

			switch res := task.(type) {
			case models.User:
				if res.ID != tc.expectedResponse.(models.User).ID {
					t.Errorf("expected: %v, got: %v", tc.expectedResponse, res)
				}
			case nil:
				if tc.expectedResponse != nil {
					t.Errorf("expected: %v, got: %v", tc.expectedResponse, res)
				}
			}
		})
	}
}
