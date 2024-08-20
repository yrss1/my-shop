package handler

import (
	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
	"github.com/yrss1/my-shop/api-gateway/intrenal/config"
	"github.com/yrss1/my-shop/api-gateway/intrenal/handler/http"
	"github.com/yrss1/my-shop/api-gateway/pkg/server/response"
	"github.com/yrss1/my-shop/api-gateway/pkg/server/router"
)

type Dependencies struct {
	Configs config.Configs
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
		proxyHandler := http.NewProxyHandler()

		api := h.HTTP.Group(h.dependencies.Configs.APP.Path)
		{
			proxyHandler.Routes(api, h.dependencies.Configs)
		}
		return
	}
}
