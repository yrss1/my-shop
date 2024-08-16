package http

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/yrss1/my-shop/tree/main/product/internal/domain/product"
	"github.com/yrss1/my-shop/tree/main/product/internal/service/productService"
	"github.com/yrss1/my-shop/tree/main/product/pkg/helpers"
	"github.com/yrss1/my-shop/tree/main/product/pkg/server/response"
	"github.com/yrss1/my-shop/tree/main/product/pkg/store"
)

type ProductHandler struct {
	productService *productService.Service
}

func NewProductHandler(s *productService.Service) *ProductHandler {
	return &ProductHandler{
		productService: s,
	}
}

func (h *ProductHandler) Routes(r *gin.RouterGroup) {
	api := r.Group("/products")
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
// @Router /products [get]
func (h *ProductHandler) list(c *gin.Context) {
	res, err := h.productService.ListProducts(c)
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
// @Router /products [post]
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

	res, err := h.productService.CreateProduct(c, req)
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
// @Router /products/{id} [get]
func (h *ProductHandler) get(c *gin.Context) {
	id := c.Param("id")
	//message := pb.Message{Body: "hello"}
	//data, err := h.userGRPCService.SayHello(c, &message)
	//if err != nil {
	//	return
	//}
	//fmt.Println(data)
	res, err := h.productService.GetProduct(c, id)
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
// @Param product body product.Request true "Product request"
// @Success 200 {string} string "ok"
// @Failure 400 {object} response.Object
// @Failure 404 {object} response.Object
// @Failure 500 {object} response.Object
// @Router /products/{id} [put]
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

	if err := h.productService.UpdateProduct(c, id, req); err != nil {
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
// @Router /products/{id} [delete]
func (h *ProductHandler) delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.productService.DeleteProduct(c, id); err != nil {
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
// @Router /products/search [get]
func (h *ProductHandler) search(c *gin.Context) {
	req := product.Request{
		Name:     helpers.GetStringPtr(c.Query("name")),
		Category: helpers.GetStringPtr(c.Query("category")),
	}

	if err := req.IsEmpty(); err != nil {
		response.BadRequest(c, errors.New("incorrect query"), nil)
		return
	}

	res, err := h.productService.SearchProduct(c, req)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.OK(c, res)
}
