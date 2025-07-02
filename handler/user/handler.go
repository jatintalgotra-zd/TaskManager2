package user

import (
	"strconv"

	"gofr.dev/pkg/gofr"
	gofrhttp "gofr.dev/pkg/gofr/http"
	"gofr.dev/pkg/gofr/http/response"

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
		return nil, gofrhttp.ErrorInvalidParam{}
	}

	id, err2 := h.service.Create(ctx, &task)
	if err2 != nil {
		return nil, err2
	}

	return id, nil
}

func (h *handler) GetByIDHandler(ctx *gofr.Context) (any, error) {
	id, err := strconv.Atoi(ctx.PathParam("id"))
	if err != nil {
		return nil, gofrhttp.ErrorInvalidParam{Params: []string{ctx.PathParam("id")}}
	}

	user, err2 := h.service.GetByID(ctx, int64(id))
	if err2 != nil {
		return nil, err2
	}

	return response.Raw{Data: user}, nil
}
