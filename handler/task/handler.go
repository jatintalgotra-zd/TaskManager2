package task

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
	var task models.Task

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

func (h *handler) GetAllHandler(ctx *gofr.Context) (any, error) {
	tasks, err := h.service.GetAll()
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (h *handler) GetByIDHandler(ctx *gofr.Context) (any, error) {
	id, err := strconv.Atoi(ctx.PathParam("id"))
	if err != nil {
		return nil, err
	}

	task, err2 := h.service.GetByID(int64(id))
	if err2 != nil {
		return nil, err2
	}

	return task, nil
}

func (h *handler) PutHandler(ctx *gofr.Context) (any, error) {
	var task models.Task

	err := ctx.Bind(&task)
	if err != nil {
		return nil, err
	}

	err = h.service.Update(&task)
	if err != nil {
		return nil, err
	}

	// doubt
	return nil, nil
}

func (h *handler) DeleteHandler(ctx *gofr.Context) (any, error) {
	id, err := strconv.Atoi(ctx.PathParam("id"))
	if err != nil {
		return nil, err
	}

	err = h.service.Delete(int64(id))
	if err != nil {
		return nil, err
	}

	return nil, nil
}
