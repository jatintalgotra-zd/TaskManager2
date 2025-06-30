package user

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
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
	userHandler := New(mockSvc)

	testcases := []struct {
		description  string
		method       string
		reader       io.Reader
		response     http.ResponseWriter
		expectedCode int
		expectedBody string
		mockCall     bool
		mockUser     *models.User
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
			&models.User{},
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
			&models.User{},
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
			&models.User{},
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
			&models.User{},
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
			&models.User{ID: 999},
			utils.ErrTest,
		},
		{
			"write error",
			http.MethodPost,
			bytes.NewReader([]byte(`{"id":1}`)),
			&errWriter{},
			http.StatusInternalServerError,
			"",
			true,
			&models.User{ID: 1},
			nil,
		},
	}

	for _, tc := range testcases {
		if tc.mockCall {
			mockSvc.EXPECT().Create(tc.mockUser).Return(int64(1), tc.mockErr)
		}

		r := httptest.NewRequest(tc.method, "/", tc.reader)

		userHandler.Post(tc.response, r)

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
	userHandler := New(mockSvc)

	testcases := []struct {
		description  string
		method       string
		id           string
		response     http.ResponseWriter
		expectedCode int
		expectedBody string
		mockCall     bool
		mockInput    int64
		mockOutput   *models.User
		mockErr      error
	}{
		{
			"success",
			http.MethodGet,
			"1",
			httptest.NewRecorder(),
			http.StatusOK,
			`{"id":1,"name":"","email":""}`,
			true,
			1,
			&models.User{ID: 1},
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
			"",
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
			&models.User{ID: 1},
			nil,
		},
	}

	for _, tc := range testcases {
		if tc.mockCall {
			mockSvc.EXPECT().GetByID(tc.mockInput).Return(tc.mockOutput, tc.mockErr)
		}

		r := httptest.NewRequest(tc.method, "/", http.NoBody)
		r.SetPathValue("id", tc.id)

		userHandler.GetByID(tc.response, r)

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
