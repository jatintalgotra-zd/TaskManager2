package user

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"TaskManager2/models"
)

var errTest = errors.New("test error")

type errReader int

func (errReader) Read(_ []byte) (n int, err error) {
	return 0, errTest
}

type errWriter struct {
	Code int
}

func (*errWriter) Header() http.Header {
	return nil
}
func (*errWriter) Write(_ []byte) (int, error) {
	return 0, errTest
}
func (e *errWriter) WriteHeader(statusCode int) {
	e.Code = statusCode
}

func TestHandler_Post(t *testing.T) {
	mockSvc := MockService{}
	userHandler := New(&mockSvc)

	// testcase 1 - invalid method
	r := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	w := httptest.NewRecorder()
	userHandler.Post(w, r)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("want %d, got %d", http.StatusMethodNotAllowed, w.Code)
		return
	}

	// testcase 2 - success
	user := models.User{}

	uBytes, err := json.Marshal(user)
	if err != nil {
		t.Error(err)
	}

	reader := bytes.NewReader(uBytes)
	r = httptest.NewRequest(http.MethodPost, "/", reader)
	w = httptest.NewRecorder()
	userHandler.Post(w, r)

	if w.Code != http.StatusCreated {
		t.Errorf("want %d, got %d", http.StatusCreated, w.Code)
	}

	if w.Body.String() != "1" {
		t.Errorf("Expected 1, got %s", w.Body.String())
	}

	// testcase 3 - read error
	r = httptest.NewRequest(http.MethodPost, "/", errReader(0))
	w = httptest.NewRecorder()
	userHandler.Post(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("want %d, got %d", http.StatusBadRequest, w.Code)
		return
	}

	// testcase 4 - unmarshal fail
	r = httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(`invalid json`)))
	w = httptest.NewRecorder()
	userHandler.Post(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("want %d, got %d", http.StatusBadRequest, w.Code)
		return
	}

	// testcase 5 - service error
	user = models.User{ID: 999}

	uBytes, err = json.Marshal(user)
	if err != nil {
		t.Error(err)
		return
	}

	reader = bytes.NewReader(uBytes)
	r = httptest.NewRequest(http.MethodPost, "/", reader)
	w = httptest.NewRecorder()
	userHandler.Post(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("want %d, got %d", http.StatusBadRequest, w.Code)
		return
	}

	// testcase 6 - write error
	r = httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(`{"id":1}`)))
	wErr := &errWriter{}
	userHandler.Post(wErr, r)

	if wErr.Code != http.StatusInternalServerError {
		t.Errorf("want %d, got %d", http.StatusInternalServerError, wErr.Code)
	}
}

func TestHandler_GetByID(t *testing.T) {
	mockSvc := MockService{}
	userHandler := New(&mockSvc)

	// testcase 1 - success
	r := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	r.SetPathValue("id", "1")

	w := httptest.NewRecorder()
	userHandler.GetByID(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("want %d, got %d", http.StatusOK, w.Code)
		return
	}

	var user models.User

	err := json.Unmarshal(w.Body.Bytes(), &user)
	if err != nil {
		t.Error("failed to unmarshal response body")
		return
	}

	if user.ID != 1 {
		t.Errorf("unexpected user: %+v", user)
		return
	}

	// testcase 2 - invalid method
	r = httptest.NewRequest(http.MethodPost, "/", http.NoBody)
	r.SetPathValue("id", "1")

	w = httptest.NewRecorder()
	userHandler.GetByID(w, r)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("want %d, got %d", http.StatusMethodNotAllowed, w.Code)
		return
	}

	// testcase 3 - bad id (non-integer)
	r = httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	r.SetPathValue("id", "abc")

	w = httptest.NewRecorder()
	userHandler.GetByID(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("want %d, got %d", http.StatusBadRequest, w.Code)
		return
	}

	// testcase 4 - not found
	r = httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	r.SetPathValue("id", "999")

	w = httptest.NewRecorder()
	userHandler.GetByID(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("want %d, got %d", http.StatusNotFound, w.Code)
		return
	}

	// testcase 5 - write error
	r = httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	r.SetPathValue("id", "1")

	wErr := &errWriter{}
	userHandler.GetByID(wErr, r)

	if wErr.Code != http.StatusInternalServerError {
		t.Errorf("want %d, got %d", http.StatusInternalServerError, wErr.Code)
	}
}
