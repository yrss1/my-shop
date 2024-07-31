package http

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/yrss1/my-shop/tree/main/user/internal/domain/user"
	"github.com/yrss1/my-shop/tree/main/user/internal/service/shop"
	"github.com/yrss1/my-shop/tree/main/user/pkg/helpers"
	"github.com/yrss1/my-shop/tree/main/user/pkg/server/response"
	"github.com/yrss1/my-shop/tree/main/user/pkg/store"
)

type UserHandler struct {
	shopService *shop.Service
}

func NewUserHandler(s *shop.Service) *UserHandler {
	return &UserHandler{shopService: s}
}

func (h *UserHandler) Routes(r *gin.RouterGroup) {
	api := r.Group("")
	{
		api.GET("/", h.list)
		api.POST("/", h.add)

		api.GET("/:id", h.get)
		api.PUT("/:id", h.update)
		api.DELETE("/:id", h.delete)

		api.GET("/search", h.search)

	}
}

func (h *UserHandler) list(c *gin.Context) {
	res, err := h.shopService.ListUsers(c)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.OK(c, res)
}

func (h *UserHandler) add(c *gin.Context) {
	req := user.Request{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err, req)
		return
	}
	if err := req.Validate(); err != nil {
		response.BadRequest(c, err, req)
		return
	}

	res, err := h.shopService.CreateUser(c, req)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.OK(c, res)
}

func (h *UserHandler) get(c *gin.Context) {
	id := c.Param("id")

	res, err := h.shopService.GetUser(c, id)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrorNotFound):
			response.NotFound(c, err)
		default:
			response.InternalServerError(c, err)
		}
		return
	}

	response.OK(c, res)
}

func (h *UserHandler) update(c *gin.Context) {
	id := c.Param("id")
	req := user.Request{}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err, req)
		return
	}

	if err := req.IsEmpty("update"); err != nil {
		response.BadRequest(c, err, req)
		return
	}

	if err := h.shopService.UpdateUser(c, id, req); err != nil {
		switch {
		case errors.Is(err, store.ErrorNotFound):
			response.NotFound(c, err)
		default:
			response.InternalServerError(c, err)
		}
		return
	}

	response.OK(c, "ok")
}

func (h *UserHandler) delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.shopService.DeleteUser(c, id); err != nil {
		switch {
		case errors.Is(err, store.ErrorNotFound):
			response.NotFound(c, err)
		default:
			response.InternalServerError(c, err)
		}
		return
	}

	response.OK(c, id)
}

func (h *UserHandler) search(c *gin.Context) {
	req := user.Request{
		Name:  helpers.GetStringPtr(c.Query("name")),
		Email: helpers.GetStringPtr(c.Query("email")),
	}

	if err := req.IsEmpty("search"); err != nil {
		response.BadRequest(c, errors.New("incorrect query"), nil)
		return
	}

	res, err := h.shopService.SearchUser(c, req)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.OK(c, res)
}
