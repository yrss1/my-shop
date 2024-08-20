package handler

import (
	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/yrss1/my-shop/order/docs"
	"github.com/yrss1/my-shop/order/internal/config"
	"github.com/yrss1/my-shop/order/internal/handler/http"
	"github.com/yrss1/my-shop/order/internal/service/orderService"
	"github.com/yrss1/my-shop/order/pkg/server/response"
	"github.com/yrss1/my-shop/order/pkg/server/router"
)

type Dependencies struct {
	Configs config.Configs

	OrderService *orderService.Service
}
type Handler struct {
	dependencies Dependencies
	HTTP         *gin.Engine
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

		orderHandler := http.NewOrderHandler(h.dependencies.OrderService)

		api := h.HTTP.Group(h.dependencies.Configs.APP.Path)
		{
			orderHandler.Routes(api)
		}
		return
	}
}
