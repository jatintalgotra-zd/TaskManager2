package task

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"go.uber.org/mock/gomock"

	"TaskManager2/models"
	"TaskManager2/utils"
)

type errReader int

func (errReader) Read(_ []byte) (n int, err error) {
	return 0, utils.ErrTest
}

type errWriter struct {
	Code int
}

func (*errWriter) Header() http.Header {
	return nil
}
func (*errWriter) Write(_ []byte) (int, error) {
	return 0, utils.ErrTest
}
func (e *errWriter) WriteHeader(statusCode int) {
	e.Code = statusCode
}

func TestHandler_Post(t *testing.T) {
	controller := gomock.NewController(t)
	mockSvc := NewMockService(controller)
	taskHandler := New(mockSvc)

	testcases := []struct {
		description  string
		method       string
		reader       io.Reader
		response     http.ResponseWriter
		expectedCode int
		expectedBody string
		mockCall     bool
		mockTask     *models.Task
		mockErr      error
	}{
		{
			"success",
			http.MethodPost,
			func() io.Reader {
				user := models.User{}
				uBytes, _ := json.Marshal(user)
				return bytes.NewReader(uBytes)
			}(),
			httptest.NewRecorder(),
			http.StatusCreated,
			"1",
			true,
			&models.Task{},
			nil,
		},
		{
			"invalid method",
			http.MethodGet,
			http.NoBody,
			httptest.NewRecorder(),
			http.StatusMethodNotAllowed,
			"",
			false,
			&models.Task{},
			nil,
		},
		{
			"Read error",
			http.MethodPost,
			errReader(0),
			httptest.NewRecorder(),
			http.StatusBadRequest,
			"",
			false,
			&models.Task{},
			nil,
		},
		{
			"unmarshal fail",
			http.MethodPost,
			bytes.NewReader([]byte(`invalid json`)),
			httptest.NewRecorder(),
			http.StatusBadRequest,
			"",
			false,
			&models.Task{},
			nil,
		},
		{
			"service error",
			http.MethodPost,
			func() io.Reader {
				user2 := models.User{ID: 999}
				uBytes2, _ := json.Marshal(user2)
				return bytes.NewReader(uBytes2)
			}(),
			httptest.NewRecorder(),
			http.StatusBadRequest,
			"",
			true,
			&models.Task{ID: 999},
			utils.ErrTest,
		},
		{
			"write error",
			http.MethodPost,
			bytes.NewReader([]byte(`{"id":1}`)),
			&errWriter{},
			http.StatusBadRequest,
			"",
			true,
			&models.Task{ID: 1},
			nil,
		},
	}

	for _, tc := range testcases {
		if tc.mockCall {
			mockSvc.EXPECT().Create(tc.mockTask).Return(int64(1), tc.mockErr)
		}

		r := httptest.NewRequest(tc.method, "/", tc.reader)

		taskHandler.Post(tc.response, r)

		switch w := tc.response.(type) {
		case *errWriter:
			if w.Code != tc.expectedCode {
				t.Errorf("expected: %d, got: %d", tc.expectedCode, w.Code)
			}

		case *httptest.ResponseRecorder:
			if w.Code != tc.expectedCode {
				t.Errorf("expected: %d, got: %d", tc.expectedCode, w.Code)
			}

			if w.Body.String() != tc.expectedBody {
				t.Errorf("expected: %s, got: %s", tc.expectedBody, w.Body.String())
			}
		}
	}
}

func TestHandler_Get(t *testing.T) {
	controller := gomock.NewController(t)
	mockSvc := NewMockService(controller)
	taskHandler := New(mockSvc)

	testcases := []struct {
		description  string
		method       string
		response     http.ResponseWriter
		expectedCode int
		expectedBody string
		mockCall     bool
		mockInput    int64
		mockOutput   []models.Task
		mockErr      error
	}{
		{
			"success",
			http.MethodGet,
			httptest.NewRecorder(),
			http.StatusOK,
			`[{"id":1,"desc":"","status":false,"user_id":0}]`,
			true,
			1,
			[]models.Task{{ID: 1}},
			nil,
		},
		{
			"invalid method",
			http.MethodPost,
			httptest.NewRecorder(),
			http.StatusMethodNotAllowed,
			"",
			false,
			0,
			nil,
			nil,
		},
		{
			"service error",
			http.MethodGet,
			httptest.NewRecorder(),
			http.StatusInternalServerError,
			"",
			true,
			999,
			nil,
			utils.ErrTest,
		},
		{
			"write error",
			http.MethodGet,
			&errWriter{},
			http.StatusInternalServerError,
			"",
			true,
			1,
			[]models.Task{{ID: 1}},
			nil,
		},
	}

	for _, tc := range testcases {
		if tc.mockCall {
			mockSvc.EXPECT().GetAll().Return(tc.mockOutput, tc.mockErr)
		}

		r := httptest.NewRequest(tc.method, "/", http.NoBody)

		taskHandler.Get(tc.response, r)

		switch w := tc.response.(type) {
		case *errWriter:
			if w.Code != tc.expectedCode {
				t.Errorf("expected: %d, got: %d", tc.expectedCode, w.Code)
			}

		case *httptest.ResponseRecorder:
			if w.Code != tc.expectedCode {
				t.Errorf("expected: %d, got: %d", tc.expectedCode, w.Code)
			}

			if w.Body.String() != tc.expectedBody {
				t.Errorf("expected: %s, got: %s", tc.expectedBody, w.Body.String())
			}
		}
	}
}

func TestHandler_GetByID(t *testing.T) {
	controller := gomock.NewController(t)
	mockSvc := NewMockService(controller)
	taskHandler := New(mockSvc)

	testcases := []struct {
		description  string
		method       string
		id           string
		response     http.ResponseWriter
		expectedCode int
		expectedBody string
		mockCall     bool
		mockInput    int64
		mockOutput   *models.Task
		mockErr      error
	}{
		{
			"success",
			http.MethodGet,
			"1",
			httptest.NewRecorder(),
			http.StatusOK,
			`{"id":1,"desc":"","status":false,"user_id":0}`,
			true,
			1,
			&models.Task{ID: 1},
			nil,
		},
		{
			"invalid method",
			http.MethodPost,
			"1",
			httptest.NewRecorder(),
			http.StatusMethodNotAllowed,
			"",
			false,
			0,
			nil,
			nil,
		},
		{
			"bad id (non-integer)",
			http.MethodGet,
			"abc",
			httptest.NewRecorder(),
			http.StatusBadRequest,
			"",
			false,
			0,
			nil,
			nil,
		},
		{
			"not found",
			http.MethodGet,
			"999",
			httptest.NewRecorder(),
			http.StatusNotFound,
			"null",
			true,
			999,
			nil,
			utils.ErrTest,
		},
		{
			"write error",
			http.MethodGet,
			"1",
			&errWriter{},
			http.StatusInternalServerError,
			"",
			true,
			1,
			&models.Task{ID: 1},
			nil,
		},
	}

	for _, tc := range testcases {
		if tc.mockCall {
			mockSvc.EXPECT().GetByID(tc.mockInput).Return(tc.mockOutput, tc.mockErr)
		}

		r := httptest.NewRequest(tc.method, "/", http.NoBody)
		r.SetPathValue("id", tc.id)

		taskHandler.GetByID(tc.response, r)

		switch w := tc.response.(type) {
		case *errWriter:
			if w.Code != tc.expectedCode {
				t.Errorf("expected: %d, got: %d", tc.expectedCode, w.Code)
			}

		case *httptest.ResponseRecorder:
			if w.Code != tc.expectedCode {
				t.Errorf("expected: %d, got: %d", tc.expectedCode, w.Code)
			}

			if w.Body.String() != tc.expectedBody {
				t.Errorf("expected: %s, got: %s", tc.expectedBody, w.Body.String())
			}
		}
	}
}

func TestHandler_Put(t *testing.T) {
	controller := gomock.NewController(t)
	mockSvc := NewMockService(controller)
	taskHandler := New(mockSvc)

	testcases := []struct {
		description  string
		method       string
		reader       io.Reader
		response     *httptest.ResponseRecorder
		expectedCode int
		mockCall     bool
		mockTask     *models.Task
		mockErr      error
	}{
		{
			"success",
			http.MethodPut,
			func() io.Reader {
				user := models.User{}
				uBytes, _ := json.Marshal(user)
				return bytes.NewReader(uBytes)
			}(),
			httptest.NewRecorder(),
			http.StatusNoContent,
			true,
			&models.Task{},
			nil,
		},
		{
			"invalid method",
			http.MethodGet,
			http.NoBody,
			httptest.NewRecorder(),
			http.StatusMethodNotAllowed,
			false,
			&models.Task{},
			nil,
		},
		{
			"Read error",
			http.MethodPut,
			errReader(0),
			httptest.NewRecorder(),
			http.StatusBadRequest,
			false,
			&models.Task{},
			nil,
		},
		{
			"unmarshal fail",
			http.MethodPut,
			bytes.NewReader([]byte(`invalid json`)),
			httptest.NewRecorder(),
			http.StatusBadRequest,
			false,
			&models.Task{},
			nil,
		},
		{
			"service error",
			http.MethodPut,
			func() io.Reader {
				user2 := models.User{ID: 999}
				uBytes2, _ := json.Marshal(user2)
				return bytes.NewReader(uBytes2)
			}(),
			httptest.NewRecorder(),
			http.StatusInternalServerError,
			true,
			&models.Task{ID: 999},
			utils.ErrTest,
		},
	}

	for _, tc := range testcases {
		if tc.mockCall {
			mockSvc.EXPECT().Update(tc.mockTask).Return(tc.mockErr)
		}

		r := httptest.NewRequest(tc.method, "/", tc.reader)

		taskHandler.Put(tc.response, r)

		w := tc.response
		if w.Code != tc.expectedCode {
			t.Errorf("expected: %d, got: %d", tc.expectedCode, w.Code)
		}
	}
}

func TestHandler_DeleteByID(t *testing.T) {
	controller := gomock.NewController(t)
	mockSvc := NewMockService(controller)
	taskHandler := New(mockSvc)

	testcases := []struct {
		description  string
		method       string
		id           string
		response     *httptest.ResponseRecorder
		expectedCode int
		mockCall     bool
		mockErr      error
	}{
		{
			"success",
			http.MethodDelete,
			"1",
			httptest.NewRecorder(),
			http.StatusNoContent,
			true,
			nil,
		},
		{
			"invalid method",
			http.MethodGet,
			"1",
			httptest.NewRecorder(),
			http.StatusMethodNotAllowed,
			false,
			nil,
		},
		{
			"invalid id (non-integer)",
			http.MethodDelete,
			"abc",
			httptest.NewRecorder(),
			http.StatusBadRequest,
			false,
			nil,
		},
		{
			"service error",
			http.MethodDelete,
			"999",
			httptest.NewRecorder(),
			http.StatusInternalServerError,
			true,
			utils.ErrTest,
		},
	}

	for _, tc := range testcases {
		if tc.mockCall {
			idInt, _ := strconv.Atoi(tc.id)
			mockSvc.EXPECT().Delete(int64(idInt)).Return(tc.mockErr)
		}

		r := httptest.NewRequest(tc.method, "/", http.NoBody)
		r.SetPathValue("id", tc.id)

		taskHandler.DeleteByID(tc.response, r)

		w := tc.response
		if w.Code != tc.expectedCode {
			t.Errorf("expected: %d, got: %d", tc.expectedCode, w.Code)
		}
	}
}
