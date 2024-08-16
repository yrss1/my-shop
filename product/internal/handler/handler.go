package handler

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/yrss1/my-shop/tree/main/product/docs"
	"github.com/yrss1/my-shop/tree/main/product/internal/config"
	"github.com/yrss1/my-shop/tree/main/product/internal/handler/http"
	"github.com/yrss1/my-shop/tree/main/product/internal/service/productService"
	"github.com/yrss1/my-shop/tree/main/product/pkg/server/router"
)

type Dependencies struct {
	Configs config.Configs

	ProductService *productService.Service
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

		productHandler := http.NewProductHandler(h.dependencies.ProductService)

		api := h.HTTP.Group(h.dependencies.Configs.APP.Path)
		{
			productHandler.Routes(api)
		}
		return
	}
}
