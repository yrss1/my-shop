package http

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/yrss1/my-shop/tree/main/user/internal/domain/user"
	"github.com/yrss1/my-shop/tree/main/user/internal/service/userService"
	"github.com/yrss1/my-shop/tree/main/user/pkg/helpers"
	"github.com/yrss1/my-shop/tree/main/user/pkg/server/response"
	"github.com/yrss1/my-shop/tree/main/user/pkg/store"
)

type UserHandler struct {
	userService *userService.Service
}

func NewUserHandler(s *userService.Service) *UserHandler {
	return &UserHandler{userService: s}
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

// list godoc
// @Summary List users
// @Description Get all users
// @Tags users
// @Accept  json
// @Produce  json
// @Success 200 {array} user.Response
// @Failure 500 {object} response.Object
// @Router / [get]
func (h *UserHandler) list(c *gin.Context) {
	res, err := h.userService.ListUsers(c)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.OK(c, res)
}

// add godoc
// @Summary Add a user
// @Description Add a new user
// @Tags users
// @Accept  json
// @Produce  json
// @Param user body user.Request true "User request"
// @Success 200 {object} user.Response
// @Failure 400 {object} response.Object
// @Failure 500 {object} response.Object
// @Router / [post]
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

	res, err := h.userService.CreateUser(c, req)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.OK(c, res)
}

// get godoc
// @Summary Get a user
// @Description Get user by ID
// @Tags users
// @Accept  json
// @Produce  json
// @Param id path string true "User ID"
// @Success 200 {object} user.Response
// @Failure 404 {object} response.Object
// @Failure 500 {object} response.Object
// @Router /{id} [get]
func (h *UserHandler) get(c *gin.Context) {
	id := c.Param("id")

	res, err := h.userService.GetUser(c, id)
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

// update godoc
// @Summary Update a user
// @Description Update user by ID
// @Tags users
// @Accept  json
// @Produce  json
// @Param id path string true "User ID"
// @Param user body user.Request true "User request"
// @Success 200 {string} string "ok"
// @Failure 400 {object} response.Object
// @Failure 404 {object} response.Object
// @Failure 500 {object} response.Object
// @Router /{id} [put]
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

	if err := h.userService.UpdateUser(c, id, req); err != nil {
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

// delete godoc
// @Summary Delete a user
// @Description Delete user by ID
// @Tags users
// @Accept  json
// @Produce  json
// @Param id path string true "User ID"
// @Success 200 {string} string "User deleted"
// @Failure 404 {object} response.Object
// @Failure 500 {object} response.Object
// @Router /{id} [delete]
func (h *UserHandler) delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.userService.DeleteUser(c, id); err != nil {
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

// search godoc
// @Summary Search users
// @Description Search users by name or email
// @Tags users
// @Accept  json
// @Produce  json
// @Param name query string false "Name"
// @Param email query string false "Email"
// @Success 200 {array} user.Response
// @Failure 400 {object} response.Object
// @Failure 500 {object} response.Object
// @Router /search [get]
func (h *UserHandler) search(c *gin.Context) {
	req := user.Request{
		Name:  helpers.GetStringPtr(c.Query("name")),
		Email: helpers.GetStringPtr(c.Query("email")),
	}

	if err := req.IsEmpty("search"); err != nil {
		response.BadRequest(c, errors.New("incorrect query"), nil)
		return
	}

	res, err := h.userService.SearchUser(c, req)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.OK(c, res)
}
