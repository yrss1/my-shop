package http

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/yrss1/my-shop/tree/main/api-gateway/intrenal/config"
	"github.com/yrss1/my-shop/tree/main/api-gateway/pkg/server/response"
	"io"
	"net/http"
)

func NewProxyHandler() *ProxyHandler {
	return &ProxyHandler{}
}

type ProxyHandler struct{}

func (h *ProxyHandler) Routes(routerGroup *gin.RouterGroup, config config.Configs) {
	routerGroup.Any("/order/*action", h.handleRequest(config.API.Order))
	routerGroup.Any("/payment/*action", h.handleRequest(config.API.Payment))
	routerGroup.Any("/product/*action", h.handleRequest(config.API.Product))
	routerGroup.Any("/user/*action", h.handleRequest(config.API.User))
}

func (h *ProxyHandler) handleRequest(targetURL string) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Param("action")
		method := c.Request.Method

		query := c.Request.URL.RawQuery

		target := targetURL + path
		if query != "" {
			target += "?" + query
		}

		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			response.InternalServerError(c, err)
			return
		}
		req, err := http.NewRequest(method, target, bytes.NewBuffer(body))
		if err != nil {
			response.InternalServerError(c, err)
			return
		}

		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		for key, values := range c.Request.Header {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			response.InternalServerError(c, err)
			return
		}
		defer resp.Body.Close()

		c.Writer.WriteHeader(resp.StatusCode)
		for key, values := range resp.Header {
			for _, value := range values {
				c.Writer.Header().Add(key, value)
			}
		}
		io.Copy(c.Writer, resp.Body)
	}
}
