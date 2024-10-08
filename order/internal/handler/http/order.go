package http

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/yrss1/my-shop/order/internal/domain/order"
	"github.com/yrss1/my-shop/order/internal/service/orderService"
	"github.com/yrss1/my-shop/order/pkg/helpers"
	"github.com/yrss1/my-shop/order/pkg/server/response"
	"github.com/yrss1/my-shop/order/pkg/store"
)

type OrderHandler struct {
	orderService *orderService.Service
}

func NewOrderHandler(s *orderService.Service) *OrderHandler {
	return &OrderHandler{orderService: s}
}

func (h *OrderHandler) Routes(r *gin.RouterGroup) {
	api := r.Group("/orders")
	{
		api.GET("/", h.list)
		api.POST("/", h.add)

		api.GET("/:id", h.get)
		api.PUT("/:id", h.update)
		api.DELETE("/:id", h.delete)

		api.GET("/search", h.search)

	}
}

// list godoc
// @Summary List orders
// @Description Get all orders
// @Tags orders
// @Accept  json
// @Produce  json
// @Success 200 {array} order.Response
// @Failure 500 {object} response.Object
// @Router /orders [get]
func (h *OrderHandler) list(c *gin.Context) {
	res, err := h.orderService.ListOrders(c)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.OK(c, res)
}

// add godoc
// @Summary Add an order
// @Description Add a new order
// @Tags orders
// @Accept  json
// @Produce  json
// @Param order body order.Request true "Order request"
// @Success 200 {object} order.Response
// @Failure 400 {object} response.Object
// @Failure 500 {object} response.Object
// @Router /orders [post]
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
	res, err := h.orderService.CreateOrder(c, req)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.OK(c, res)
}

// get godoc
// @Summary Get an order
// @Description Get order by ID
// @Tags orders
// @Accept  json
// @Produce  json
// @Param id path string true "Order ID"
// @Success 200 {object} order.Response
// @Failure 404 {object} response.Object
// @Failure 500 {object} response.Object
// @Router /orders/{id} [get]
func (h *OrderHandler) get(c *gin.Context) {
	id := c.Param("id")

	res, err := h.orderService.GetOrder(c, id)
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

// update godoc
// @Summary Update an order
// @Description Update order by ID
// @Tags orders
// @Accept  json
// @Produce  json
// @Param id path string true "Order ID"
// @Param order body order.Request true "Order request"
// @Success 200 {string} string "ok"
// @Failure 400 {object} response.Object
// @Failure 404 {object} response.Object
// @Failure 500 {object} response.Object
// @Router /orders/{id} [put]
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

	if err := h.orderService.UpdateOrder(c, id, req); err != nil {
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

// delete godoc
// @Summary Delete an order
// @Description Delete order by ID
// @Tags orders
// @Accept  json
// @Produce  json
// @Param id path string true "Order ID"
// @Success 200 {string} string "Order deleted"
// @Failure 404 {object} response.Object
// @Failure 500 {object} response.Object
// @Router /orders/{id} [delete]
func (h *OrderHandler) delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.orderService.DeleteOrder(c, id); err != nil {
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

// search godoc
// @Summary Search orders
// @Description Search orders by user ID or status
// @Tags orders
// @Accept  json
// @Produce  json
// @Param userId query string false "User ID"
// @Param status query string false "Status"
// @Success 200 {array} order.Response
// @Failure 400 {object} response.Object
// @Failure 500 {object} response.Object
// @Router /orders/search [get]
func (h *OrderHandler) search(c *gin.Context) {
	req := order.Request{
		UserID: helpers.GetStringPtr(c.Query("userId")),
		Status: helpers.GetStringPtr(c.Query("status")),
	}

	if err := req.IsEmpty("search"); err != nil {
		response.BadRequest(c, err, nil)
		return
	}

	res, err := h.orderService.SearchOrder(c, req)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.OK(c, res)
}
