package task

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"TaskManager2/models"
)

type handler struct {
	service Service
}

func New(service Service) *handler {
	return &handler{service: service}
}

// Post - Create method.
func (h *handler) Post(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var task models.Task

	tBytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(tBytes, &task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err2 := h.service.Create(&task)
	if err2 != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)

	_, err3 := w.Write([]byte(strconv.Itoa(int(id))))
	if err3 != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

// Get - get all.
func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	tasks, err := h.service.GetAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tBytes, err2 := json.Marshal(tasks)
	if err2 != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	_, err3 := w.Write(tBytes)
	if err3 != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// GetByID - get by id method.
func (h *handler) GetByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := r.PathValue("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	task, err2 := h.service.GetByID(int64(idInt))
	if err2 != nil {
		w.WriteHeader(http.StatusNotFound)
	}

	tBytes, err3 := json.Marshal(task)
	if err3 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err4 := w.Write(tBytes)
	if err4 != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// Put - update method.
func (h *handler) Put(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var task models.Task

	tBytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(tBytes, &task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.service.Update(&task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := r.PathValue("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.service.Delete(int64(idInt))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
