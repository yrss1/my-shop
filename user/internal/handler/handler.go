package handler

import (
	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/yrss1/my-shop/user/docs"
	"github.com/yrss1/my-shop/user/internal/config"
	"github.com/yrss1/my-shop/user/internal/handler/grpc_handler"
	"github.com/yrss1/my-shop/user/internal/handler/http"
	"github.com/yrss1/my-shop/user/internal/service/userService"
	"github.com/yrss1/my-shop/user/pkg/server/response"
	"github.com/yrss1/my-shop/user/pkg/server/router"
	pb "github.com/yrss1/proto-definitions/user"
	"google.golang.org/grpc"
)

type Dependencies struct {
	Configs config.Configs

	UserService *userService.Service
}
type Handler struct {
	dependencies Dependencies
	HTTP         *gin.Engine
	GRPCServer   *grpc.Server
}
type Configuration func(h *Handler) error

func New(d Dependencies, configs ...Configuration) (h *Handler, err error) {
	h = &Handler{
		dependencies: d,
	}

	for _, cfg := range configs {
		if err = cfg(h); err != nil {
			return
		}
	}

	return
}

func WithHTTPHandler() Configuration {
	return func(h *Handler) (err error) {
		h.HTTP = router.New()
		h.HTTP.Use(timeout.New(
			timeout.WithTimeout(h.dependencies.Configs.APP.Timeout),
			timeout.WithHandler(func(ctx *gin.Context) {
				ctx.Next()
			}),
			timeout.WithResponse(func(ctx *gin.Context) {
				response.StatusRequestTimeout(ctx)
			}),
		))

		docs.SwaggerInfo.BasePath = h.dependencies.Configs.APP.Path
		h.HTTP.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		userHandler := http.NewUserHandler(h.dependencies.UserService)

		api := h.HTTP.Group("/api/v1")
		{
			userHandler.Routes(api)
		}
		return
	}
}

func WithGRPCHandler() Configuration {
	return func(h *Handler) (err error) {
		h.GRPCServer = grpc.NewServer()
		pb.RegisterUserServiceServer(h.GRPCServer, grpc_handler.NewUserServiceServer(h.dependencies.UserService))
		return
	}
}
