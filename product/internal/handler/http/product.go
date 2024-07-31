package http

import (
	"errors"
	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	shopService *shop.Service
}

func NewProductHandler(s *shop.Service) *ProductHandler {
	return &ProductHandler{shopService: s}
}

func (h *ProductHandler) Routes(r *gin.RouterGroup) {
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

func (h *ProductHandler) list(c *gin.Context) {
	res, err := h.shopService.ListProducts(c)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.OK(c, res)
}

func (h *ProductHandler) add(c *gin.Context) {
	req := product.Request{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err, req)
		return
	}
	if err := req.Validate(); err != nil {
		response.BadRequest(c, err, req)
		return
	}

	res, err := h.shopService.CreateProduct(c, req)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.OK(c, res)
}

func (h *ProductHandler) get(c *gin.Context) {
	id := c.Param("id")

	res, err := h.shopService.GetProduct(c, id)
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

func (h *ProductHandler) update(c *gin.Context) {
	id := c.Param("id")
	req := product.Request{}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err, req)
		return
	}

	if err := req.IsEmpty(); err != nil {
		response.BadRequest(c, err, req)
		return
	}

	if err := h.shopService.UpdateProduct(c, id, req); err != nil {
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

func (h *ProductHandler) delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.shopService.DeleteProduct(c, id); err != nil {
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

func (h *ProductHandler) search(c *gin.Context) {
	req := product.Request{
		Name:     helpers.GetStringPtr(c.Query("name")),
		Category: helpers.GetStringPtr(c.Query("category")),
	}

	if err := req.IsEmpty(); err != nil {
		response.BadRequest(c, errors.New("incorrect query"), nil)
		return
	}

	res, err := h.shopService.SearchProduct(c, req)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.OK(c, res)
}
