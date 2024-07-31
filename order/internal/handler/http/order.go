package http

import (
	"errors"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	shopService *shop.Service
}

func NewOrderHandler(s *shop.Service) *OrderHandler {
	return &OrderHandler{shopService: s}
}

func (h *OrderHandler) Routes(r *gin.RouterGroup) {
	api := r.Group("")
	{
		api.GET("/", h.list)
		api.POST("/", h.add)

		api.GET("/:id", h.get)
		api.PUT("/:id", h.update)
		api.DELETE("/:id", h.delete)

		api.GET("/search", h.search)

	}
}

func (h *OrderHandler) list(c *gin.Context) {
	res, err := h.shopService.ListOrders(c)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.OK(c, res)
}

func (h *OrderHandler) add(c *gin.Context) {
	req := order.Request{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err, req)
		return
	}
	if err := req.Validate(); err != nil {
		response.BadRequest(c, err, req)
		return
	}
	// нужно добавить проверку на продукты

	res, err := h.shopService.CreateOrder(c, req)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.OK(c, res)
}

func (h *OrderHandler) get(c *gin.Context) {
	id := c.Param("id")

	res, err := h.shopService.GetOrder(c, id)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrorNotFound):
			response.NotFound(c, err)
		default:
			response.InternalServerError(c, err)
		}
		return
	}

	response.OK(c, res)
}

func (h *OrderHandler) update(c *gin.Context) {
	id := c.Param("id")
	req := order.Request{}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err, req)
		return
	}

	if err := req.IsEmpty("update"); err != nil {
		response.BadRequest(c, err, req)
		return
	}

	if err := h.shopService.UpdateOrder(c, id, req); err != nil {
		switch {
		case errors.Is(err, store.ErrorNotFound):
			response.NotFound(c, err)
		default:
			response.InternalServerError(c, err)
		}
		return
	}

	response.OK(c, "ok")
}

func (h *OrderHandler) delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.shopService.DeleteOrder(c, id); err != nil {
		switch {
		case errors.Is(err, store.ErrorNotFound):
			response.NotFound(c, err)
		default:
			response.InternalServerError(c, err)
		}
		return
	}

	response.OK(c, id)
}

func (h *OrderHandler) search(c *gin.Context) {
	req := order.Request{
		UserID: helpers.GetStringPtr(c.Query("userId")),
		Status: helpers.GetStringPtr(c.Query("status")),
	}

	if err := req.IsEmpty("search"); err != nil {
		response.BadRequest(c, err, nil)
		return
	}

	res, err := h.shopService.SearchOrder(c, req)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.OK(c, res)
}
