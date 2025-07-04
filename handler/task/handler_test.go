package task

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

func TestHandler_Post(t *testing.T) {
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
							"desc" : "test task",
							"status" :  false,
							"user_id" : 2
						}`,
			func() {
				mockSvc.EXPECT().Create(ctx, &models.Task{Desc: "test task", Status: false, UserID: 2}).Return(int64(1), nil)
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
							"desc" : "test task",
							"status" :  false,
							"user_id" : 2
						}`,
			func() {
				mockSvc.EXPECT().Create(ctx, &models.Task{Desc: "test task", Status: false, UserID: 2}).Return(int64(0), utils.ErrTest)
			},
			nil,
			utils.ErrTest,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockExpect()

			body := bytes.NewReader([]byte(tc.requestBody))
			req := httptest.NewRequest(http.MethodPost, "/task", body)
			req.Header.Set("Content-Type", "application/json")

			ctx.Request = gofrhttp.NewRequest(req)

			id, err := taskHandler.Post(ctx)
			if err != nil && err.Error() != tc.expectedError.Error() {
				t.Errorf("error, expected %v, got %v", tc.expectedError, err)
			}

			if id != tc.expectedResponse {
				t.Errorf("expected: %v, got: %v", tc.expectedResponse, id)
			}
		})
	}
}

func TestHandler_GetAllHandler(t *testing.T) {
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
		mockExpect       func()
		expectedResponse any
		expectedError    error
	}{
		{
			"success",
			func() {
				mockSvc.EXPECT().GetAll(ctx).Return([]models.Task{{}}, nil)
			},
			[]models.Task{{}},
			nil,
		},
		{
			"service GetAll error",
			func() {
				mockSvc.EXPECT().GetAll(ctx).Return(nil, utils.ErrTest)
			},
			nil,
			utils.ErrTest,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockExpect()

			req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
			req.Header.Set("Content-Type", "application/json")
			ctx.Request = gofrhttp.NewRequest(req)

			tasks, err := taskHandler.GetAll(ctx)
			if err != nil && err.Error() != tc.expectedError.Error() {
				t.Errorf("error, expected %v, got %v", tc.expectedError, err)
			}

			switch res := tasks.(type) {
			case []models.Task:
				if len(res) != len(tc.expectedResponse.([]models.Task)) {
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

func TestHandler_GetByID(t *testing.T) {
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
				mockSvc.EXPECT().GetByID(ctx, int64(1)).Return(&models.Task{ID: 1}, nil)
			},
			&models.Task{ID: 1},
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

			req := httptest.NewRequest(http.MethodGet, "/task/{id}", http.NoBody)
			req.Header.Set("Content-Type", "application/json")

			req = mux.SetURLVars(req, map[string]string{"id": tc.requestID})
			ctx.Request = gofrhttp.NewRequest(req)

			task, err := taskHandler.GetByID(ctx)
			if err != nil && err.Error() != tc.expectedError.Error() {
				t.Errorf("error, expected %v, got %v", tc.expectedError, err)
			}

			switch res := task.(type) {
			case models.Task:
				if res.ID != tc.expectedResponse.(models.Task).ID {
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

func TestHandler_Put(t *testing.T) {
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
							"id" : 4,
							"desc": "test task",
							"status": false
						}`,
			func() {
				mockSvc.EXPECT().Update(ctx, &models.Task{ID: 4, Desc: "test task", Status: false}).Return(nil)
			},
			nil,
			nil,
		},
		{"bind error",
			`describe":"test task","status":false,"user_id":1}`,
			func() {},
			nil,
			gofrhttp.ErrorInvalidParam{},
		},
		{
			"service update error",
			`{
							"id" : 4,
							"desc": "test task",
							"status": false
						}`,
			func() {
				mockSvc.EXPECT().Update(ctx, &models.Task{ID: 4, Desc: "test task", Status: false}).Return(utils.ErrTest)
			},
			nil,
			utils.ErrTest,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockExpect()

			body := bytes.NewReader([]byte(tc.requestBody))
			req := httptest.NewRequest(http.MethodPut, "/1", body)
			req.Header.Set("Content-Type", "application/json")

			ctx.Request = gofrhttp.NewRequest(req)

			_, err := taskHandler.Put(ctx)
			if err != nil && err.Error() != tc.expectedError.Error() {
				t.Errorf("error, expected %v, got %v", tc.expectedError, err)
			}
		})
	}
}

func TestHandler_Delete(t *testing.T) {
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
				mockSvc.EXPECT().Delete(ctx, int64(1)).Return(nil)
			},
			nil,
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
				mockSvc.EXPECT().Delete(ctx, int64(1)).Return(utils.ErrTest)
			},
			nil,
			utils.ErrTest,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockExpect()

			req := httptest.NewRequest(http.MethodDelete, "/1", http.NoBody)
			req.Header.Set("Content-Type", "application/json")
			req = mux.SetURLVars(req, map[string]string{"id": tc.requestID})
			ctx.Request = gofrhttp.NewRequest(req)

			_, err := taskHandler.Delete(ctx)
			if err != nil && err.Error() != tc.expectedError.Error() {
				t.Errorf("error, expected %v, got %v", tc.expectedError, err)
			}
		})
	}
}
