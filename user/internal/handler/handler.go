package handler

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/yrss1/my-shop/tree/main/user/docs"
	"github.com/yrss1/my-shop/tree/main/user/internal/config"
	"github.com/yrss1/my-shop/tree/main/user/internal/handler/http"
	"github.com/yrss1/my-shop/tree/main/user/internal/service/shop"
	"github.com/yrss1/my-shop/tree/main/user/pkg/server/router"
	"log"
	"net"
)

type Dependencies struct {
	Configs config.Configs

	ShopService *shop.Service
}
type Handler struct {
	dependencies Dependencies
	HTTP         *gin.Engine
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

		docs.SwaggerInfo.BasePath = h.dependencies.Configs.APP.Path
		h.HTTP.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		userHandler := http.NewUserHandler(h.dependencies.ShopService)

		api := h.HTTP.Group("/api/v1/")
		{
			userHandler.Routes(api)
		}
		return
	}
}

func WithGRPCHandler() Configuration {
	return func(h *Handler) (err error) {
		// Register gRPC service
		proto.RegisterUserServiceServer(h.GRPCServer, NewUserService(h.dependencies.ShopService))

		// Start gRPC server
		go func() {
			lis, err := net.Listen("tcp", ":50051") // Use your desired port
			if err != nil {
				log.Fatalf("Failed to listen: %v", err)
			}
			log.Println("gRPC server listening on port 50051")
			if err := h.GRPCServer.Serve(lis); err != nil {
				log.Fatalf("Failed to serve: %v", err)
			}
		}()

		return
	}
}
