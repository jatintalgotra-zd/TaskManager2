package user

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

	var user models.User

	uBytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(uBytes, &user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err2 := h.service.Create(&user)
	if err2 != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)

	_, err3 := w.Write([]byte(strconv.Itoa(int(id))))
	if err3 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
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

	user, err2 := h.service.GetByID(int64(idInt))
	if err2 != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	uBytes, err3 := json.Marshal(user)
	if err3 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err4 := w.Write(uBytes)
	if err4 != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
