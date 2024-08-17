package handler

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	pb "github.com/yrss1/my-shop/proto/user"
	"github.com/yrss1/my-shop/user/docs"
	"github.com/yrss1/my-shop/user/internal/config"
	"github.com/yrss1/my-shop/user/internal/handler/grpc_handler"
	"github.com/yrss1/my-shop/user/internal/handler/http"
	"github.com/yrss1/my-shop/user/internal/service/userService"
	"github.com/yrss1/my-shop/user/pkg/server/router"
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
		HTTP:         router.New(),
		GRPCServer:   grpc.NewServer(),
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
		pb.RegisterUserServiceServer(h.GRPCServer, grpc_handler.NewUserServiceServer(h.dependencies.UserService))
		return
	}
}
