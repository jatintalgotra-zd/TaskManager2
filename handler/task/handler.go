package task

import (
	"strconv"

	"gofr.dev/pkg/gofr"
	gofrhttp "gofr.dev/pkg/gofr/http"

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
		return nil, gofrhttp.ErrorInvalidParam{}
	}

	id, err := h.service.Create(ctx, &task)
	if err != nil {
		return nil, err
	}

	return id, nil
}

func (h *handler) GetAllHandler(ctx *gofr.Context) (any, error) {
	tasks, err := h.service.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (h *handler) GetByIDHandler(ctx *gofr.Context) (any, error) {
	id, err := strconv.Atoi(ctx.PathParam("id"))
	if err != nil {
		return nil, gofrhttp.ErrorInvalidParam{Params: []string{ctx.PathParam("id")}}
	}

	task, err := h.service.GetByID(ctx, int64(id))
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (h *handler) PutHandler(ctx *gofr.Context) (any, error) {
	var task models.Task

	err := ctx.Bind(&task)
	if err != nil {
		return nil, gofrhttp.ErrorInvalidParam{}
	}

	err = h.service.Update(ctx, &task)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (h *handler) DeleteHandler(ctx *gofr.Context) (any, error) {
	id, err := strconv.Atoi(ctx.PathParam("id"))
	if err != nil {
		return nil, gofrhttp.ErrorInvalidParam{Params: []string{ctx.PathParam("id")}}
	}

	err = h.service.Delete(ctx, int64(id))
	if err != nil {
		return nil, err
	}

	return nil, nil
}
