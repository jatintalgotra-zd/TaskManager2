package task

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

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
	mockSvc := MockService{}
	taskHandler := New(&mockSvc)

	// testcase 1 - invalid method
	r := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	w := httptest.NewRecorder()
	taskHandler.Post(w, r)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("want %d, got %d", http.StatusMethodNotAllowed, w.Code)
		return
	}

	// testcase 2 - success
	task := models.Task{}

	tBytes, err := json.Marshal(task)
	if err != nil {
		t.Error(err)
		return
	}

	reader := bytes.NewReader(tBytes)
	r = httptest.NewRequest(http.MethodPost, "/", reader)
	w = httptest.NewRecorder()
	taskHandler.Post(w, r)

	if w.Code != http.StatusCreated {
		t.Errorf("want %d, got %d", http.StatusCreated, w.Code)
	}

	if w.Body.String() != "1" {
		t.Errorf("Expected 1, got %s", w.Body.String())
	}

	// testcase 3 - read error
	r = httptest.NewRequest(http.MethodPost, "/", errReader(0))
	w = httptest.NewRecorder()
	taskHandler.Post(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("want %d, got %d", http.StatusBadRequest, w.Code)
		return
	}

	// testcase 4 - unmarshal fail
	r = httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(`hello world`)))
	w = httptest.NewRecorder()
	taskHandler.Post(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("want %d, got %d", http.StatusBadRequest, w.Code)
		return
	}

	// testcase 5 - service error
	task = models.Task{ID: 999}

	tBytes, err = json.Marshal(task)
	if err != nil {
		t.Error(err)
		return
	}

	reader = bytes.NewReader(tBytes)
	r = httptest.NewRequest(http.MethodPost, "/", reader)
	w = httptest.NewRecorder()
	taskHandler.Post(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("want %d, got %d", http.StatusBadRequest, w.Code)
		return
	}

	// testcase 6 - write error
	r = httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(`{"id":1}`)))
	wErr := &errWriter{}
	taskHandler.Post(wErr, r)

	if wErr.Code != http.StatusBadRequest {
		t.Errorf("want %d, got %d", http.StatusBadRequest, wErr.Code)
	}
}

func TestHandler_Get(t *testing.T) {
	mockSvc := MockService{}
	taskHandler := New(&mockSvc)

	// testcase 1 - success
	r := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	w := httptest.NewRecorder()
	taskHandler.Get(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("want %d, got %d", http.StatusOK, w.Code)
		return
	}

	var tasks []models.Task

	err := json.Unmarshal(w.Body.Bytes(), &tasks)
	if err != nil {
		t.Error("failed to unmarshal response body")
		return
	}

	if len(tasks) != 1 {
		t.Errorf("Expected 1 mock task, got %+v", tasks)
		return
	}

	// testcase 2 - invalid method
	r = httptest.NewRequest(http.MethodPost, "/", http.NoBody)
	w = httptest.NewRecorder()
	taskHandler.Get(w, r)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("want %d, got %d", http.StatusMethodNotAllowed, w.Code)
		return
	}

	// testcase 3 - service error
	mockSvc.Check = true
	r = httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	w = httptest.NewRecorder()
	taskHandler.Get(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("want %d, got %d", http.StatusInternalServerError, w.Code)
		return
	}

	mockSvc.Check = false

	// testcase 4 - write error
	r = httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	wErr := &errWriter{}
	taskHandler.Get(wErr, r)

	if wErr.Code != http.StatusInternalServerError {
		t.Errorf("want %d, got %d", http.StatusInternalServerError, wErr.Code)
	}
}

func TestHandler_GetByID(t *testing.T) {
	mockSvc := MockService{}
	taskHandler := New(&mockSvc)

	// testcase 1 - success
	r := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	r.SetPathValue("id", "1")

	w := httptest.NewRecorder()
	taskHandler.GetByID(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("want %d, got %d", http.StatusOK, w.Code)
		return
	}

	var task models.Task

	err := json.Unmarshal(w.Body.Bytes(), &task)
	if err != nil {
		t.Error("failed to unmarshal response body")
		return
	}

	if task.ID != 1 {
		t.Errorf("unexpected task: %+v", task)
		return
	}

	// testcase 2 - invalid method
	r = httptest.NewRequest(http.MethodPost, "/", http.NoBody)
	r.SetPathValue("id", "1")

	w = httptest.NewRecorder()
	taskHandler.GetByID(w, r)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("want %d, got %d", http.StatusMethodNotAllowed, w.Code)
		return
	}

	// testcase 3 - bad id (non-int)
	r = httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	r.SetPathValue("id", "abc")

	w = httptest.NewRecorder()
	taskHandler.GetByID(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("want %d, got %d", http.StatusBadRequest, w.Code)
		return
	}

	// testcase 4 - not found
	r = httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	r.SetPathValue("id", "999")

	w = httptest.NewRecorder()
	taskHandler.GetByID(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("want %d, got %d", http.StatusNotFound, w.Code)
		return
	}

	// testcase 5 - write error
	r = httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	r.SetPathValue("id", "1")

	wErr := &errWriter{}
	taskHandler.GetByID(wErr, r)

	if wErr.Code != http.StatusInternalServerError {
		t.Errorf("want %d, got %d", http.StatusInternalServerError, wErr.Code)
	}
}

func TestHandler_Put(t *testing.T) {
	mockSvc := MockService{}
	taskHandler := New(&mockSvc)

	// testcase 1 - invalid method
	r := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	w := httptest.NewRecorder()
	taskHandler.Put(w, r)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("want %d, got %d", http.StatusMethodNotAllowed, w.Code)
		return
	}

	// testcase 2 - read error
	r = httptest.NewRequest(http.MethodPut, "/", errReader(0))
	w = httptest.NewRecorder()
	taskHandler.Put(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("want %d, got %d", http.StatusBadRequest, w.Code)
		return
	}

	// testcase 3 - unmarshal fail
	r = httptest.NewRequest(http.MethodPut, "/", bytes.NewReader([]byte(`not a json`)))
	w = httptest.NewRecorder()
	taskHandler.Put(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("want %d, got %d", http.StatusBadRequest, w.Code)
		return
	}

	// testcase 4 - service error
	task := models.Task{ID: 999}

	tBytes, err := json.Marshal(task)
	if err != nil {
		t.Error(err)
		return
	}

	reader := bytes.NewReader(tBytes)
	r = httptest.NewRequest(http.MethodPut, "/", reader)
	w = httptest.NewRecorder()
	taskHandler.Put(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("want %d, got %d", http.StatusInternalServerError, w.Code)
		return
	}

	// testcase 5 - success
	task = models.Task{ID: 1}

	tBytes, err = json.Marshal(task)
	if err != nil {
		t.Error(err)
		return
	}

	reader = bytes.NewReader(tBytes)
	r = httptest.NewRequest(http.MethodPut, "/", reader)
	w = httptest.NewRecorder()
	taskHandler.Put(w, r)

	if w.Code != http.StatusNoContent {
		t.Errorf("want %d, got %d", http.StatusNoContent, w.Code)
	}
}

func TestHandler_DeleteByID(t *testing.T) {
	mockSvc := MockService{}
	taskHandler := New(&mockSvc)

	// testcase 1 - invalid method
	r := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	r.SetPathValue("id", "1")

	w := httptest.NewRecorder()
	taskHandler.DeleteByID(w, r)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("want %d, got %d", http.StatusMethodNotAllowed, w.Code)
		return
	}

	// testcase 2 - bad id (non-integer)
	r = httptest.NewRequest(http.MethodDelete, "/", http.NoBody)
	r.SetPathValue("id", "abc")

	w = httptest.NewRecorder()
	taskHandler.DeleteByID(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("want %d, got %d", http.StatusBadRequest, w.Code)
		return
	}

	// testcase 3 - service error
	r = httptest.NewRequest(http.MethodDelete, "/", http.NoBody)
	r.SetPathValue("id", "999")

	w = httptest.NewRecorder()
	taskHandler.DeleteByID(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("want %d, got %d", http.StatusInternalServerError, w.Code)
		return
	}

	// testcase 4 - success
	r = httptest.NewRequest(http.MethodDelete, "/", http.NoBody)
	r.SetPathValue("id", "1")

	w = httptest.NewRecorder()
	taskHandler.DeleteByID(w, r)

	if w.Code != http.StatusNoContent {
		t.Errorf("want %d, got %d", http.StatusNoContent, w.Code)
	}
}
