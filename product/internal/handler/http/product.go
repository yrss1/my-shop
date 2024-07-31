package http

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/yrss1/my-shop/tree/main/product/internal/domain/product"
	"github.com/yrss1/my-shop/tree/main/product/internal/service/shop"
	"github.com/yrss1/my-shop/tree/main/product/pkg/helpers"
	"github.com/yrss1/my-shop/tree/main/product/pkg/server/response"
	"github.com/yrss1/my-shop/tree/main/product/pkg/store"
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

// list godoc
// @Summary List products
// @Description Get all products
// @Tags products
// @Accept  json
// @Produce  json
// @Success 200 {array} product.Response
// @Failure 500 {object} response.Object
// @Router / [get]
func (h *ProductHandler) list(c *gin.Context) {
	res, err := h.shopService.ListProducts(c)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.OK(c, res)
}

// add godoc
// @Summary Add a product
// @Description Add a new product
// @Tags products
// @Accept  json
// @Produce  json
// @Param product body product.Request true "Product request"
// @Success 200 {object} product.Response
// @Failure 400 {object} response.Object
// @Failure 500 {object} response.Object
// @Router / [post]
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

// get godoc
// @Summary Get a product
// @Description Get product by ID
// @Tags products
// @Accept  json
// @Produce  json
// @Param id path string true "Product ID"
// @Success 200 {object} product.Response
// @Failure 404 {object} response.Object
// @Failure 500 {object} response.Object
// @Router /{id} [get]
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

// update godoc
// @Summary Update a product
// @Description Update product by ID
// @Tags products
// @Accept  json
// @Produce  json
// @Param id path string true "Product ID"
// @Param user body product.Request true "Product request"
// @Success 200 {string} string "ok"
// @Failure 400 {object} response.Object
// @Failure 404 {object} response.Object
// @Failure 500 {object} response.Object
// @Router /{id} [put]
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

// delete godoc
// @Summary Delete a product
// @Description Delete product by ID
// @Tags products
// @Accept  json
// @Produce  json
// @Param id path string true "Product ID"
// @Success 200 {string} string "Product deleted"
// @Failure 404 {object} response.Object
// @Failure 500 {object} response.Object
// @Router /{id} [delete]
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

// search godoc
// @Summary Search products
// @Description Search products by name or email
// @Tags products
// @Accept  json
// @Produce  json
// @Param name query string false "Name"
// @Param category query string false "Category"
// @Success 200 {array} product.Response
// @Failure 400 {object} response.Object
// @Failure 500 {object} response.Object
// @Router /search [get]
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
