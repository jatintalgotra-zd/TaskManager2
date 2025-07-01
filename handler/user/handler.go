package user

import (
	"strconv"

	"gofr.dev/pkg/gofr"

	"TaskManager2/models"
)

type handler struct {
	service Service
}

func New(service Service) *handler {
	return &handler{service: service}
}

func (h *handler) PostHandler(ctx *gofr.Context) (any, error) {
	var task models.User

	err := ctx.Bind(&task)
	if err != nil {
		return nil, err
	}

	id, err2 := h.service.Create(&task)
	if err2 != nil {
		return nil, err2
	}

	return id, nil
}

func (h *handler) GetByIDHandler(ctx *gofr.Context) (any, error) {
	id, err := strconv.Atoi(ctx.PathParam("id"))
	if err != nil {
		return nil, err
	}

	user, err2 := h.service.GetByID(int64(id))
	if err2 != nil {
		return nil, err2
	}

	return user, nil
}
