package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/yrss1/my-shop/tree/main/payment/internal/config"
	"github.com/yrss1/my-shop/tree/main/payment/internal/handler/http"
	"github.com/yrss1/my-shop/tree/main/payment/internal/service/epayment"
	"github.com/yrss1/my-shop/tree/main/payment/pkg/server/router"
)

type Dependencies struct {
	Configs config.Configs

	EpayService *epayment.Service
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
		paymentHandler := http.NewPaymentHandler(h.dependencies.EpayService)

		api := h.HTTP.Group("/api/v1/")
		{
			paymentHandler.Routes(api)
		}
		return
	}
}
