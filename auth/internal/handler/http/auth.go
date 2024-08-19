package http

import (
	"github.com/gin-gonic/gin"
	"github.com/yrss1/my-shop/auth/internal/service/authService"
	"github.com/yrss1/my-shop/auth/pkg/server/response"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserHandler struct {
	authService *authService.Service
}

func NewUserHandler(s *authService.Service) *UserHandler {
	return &UserHandler{authService: s}
}

func (h *UserHandler) Routes(r *gin.RouterGroup) {
	api := r.Group("/auth")
	{
		api.GET("/", h.hello)
		api.GET("/user", h.getUserByEmail)
	}
}

func (h *UserHandler) hello(c *gin.Context) {
	response.OK(c, "ok")
}

func (h *UserHandler) getUserByEmail(c *gin.Context) {
	email := c.Query("email")

	res, err := h.authService.GetUserByEmail(c, email)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			response.NotFound(c, err)
		} else {
			response.InternalServerError(c, err)
		}
		return
	}

	response.OK(c, res)
}
