package handler

import (
	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
	"github.com/yrss1/my-shop/auth/internal/config"
	"github.com/yrss1/my-shop/auth/internal/handler/http"
	"github.com/yrss1/my-shop/auth/internal/service/authService"
	"github.com/yrss1/my-shop/auth/pkg/server/response"
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

		//docs.SwaggerInfo.BasePath = h.dependencies.Configs.APP.Path
		//h.HTTP.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		userHandler := http.NewUserHandler(h.dependencies.AuthService)

		api := h.HTTP.Group(h.dependencies.Configs.APP.Path)
		{
			userHandler.Routes(api)
		}
		return
	}
}

//func WithGRPCHandler() Configuration {
//	return func(h *Handler) (err error) {
//		h.GRPCServer = grpc.NewServer()
//		pb.RegisterUserServiceServer(h.GRPCServer, grpc_handler.NewUserServiceServer(h.dependencies.AuthService))
//		return
//	}
//}
