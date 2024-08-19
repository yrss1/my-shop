package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/yrss1/my-shop/auth/internal/config"
	"github.com/yrss1/my-shop/auth/internal/handler/http"
	"github.com/yrss1/my-shop/auth/internal/service/authService"
	"github.com/yrss1/my-shop/auth/pkg/server/router"
	"google.golang.org/grpc"
)

type Dependencies struct {
	Configs config.Configs

	AuthService *authService.Service
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

		//docs.SwaggerInfo.BasePath = h.dependencies.Configs.APP.Path
		//h.HTTP.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		userHandler := http.NewUserHandler(h.dependencies.AuthService)

		api := h.HTTP.Group("/api/v1")
		{
			userHandler.Routes(api)
		}
		return
	}
}

//func WithGRPCHandler() Configuration {
//	return func(h *Handler) (err error) {
//		pb.RegisterUserServiceServer(h.GRPCServer, grpc_handler.NewUserServiceServer(h.dependencies.AuthService))
//		return
//	}
//}
