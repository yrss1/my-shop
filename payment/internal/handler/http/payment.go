package http

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/yrss1/my-shop/tree/main/payment/internal/domain/payment"
	"github.com/yrss1/my-shop/tree/main/payment/internal/provider/epay"
	"github.com/yrss1/my-shop/tree/main/payment/internal/service/epayment"
	"github.com/yrss1/my-shop/tree/main/payment/pkg/helpers"
	"github.com/yrss1/my-shop/tree/main/payment/pkg/server/response"
	"github.com/yrss1/my-shop/tree/main/payment/pkg/store"
)

type PaymentHandler struct {
	epayService *epayment.Service
}

func NewPaymentHandler(s *epayment.Service) *PaymentHandler {
	return &PaymentHandler{epayService: s}
}

func (h *PaymentHandler) Routes(r *gin.RouterGroup) {
	api := r.Group("/payments")
	{
		api.GET("/", h.list)
		api.POST("/", h.add)
		api.GET("/:id", h.get)
		api.PUT("/:id", h.update)
		api.DELETE("/:id", h.delete)
		api.GET("/search", h.search)

		//api.POST("/token", h.getToken)
		//api.GET("/status/:id", h.getStatus)
		//api.GET("/pay", h.pay)

	}
}

// list godoc
// @Summary List payments
// @Description Get all payments
// @Tags payments
// @Accept  json
// @Produce  json
// @Success 200 {array} payment.Response
// @Failure 500 {object} response.Object
// @Router /payments [get]
func (h *PaymentHandler) list(c *gin.Context) {
	res, err := h.epayService.ListPayments(c)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.OK(c, res)
}

// add godoc
// @Summary Add a payment
// @Description Add a new payment
// @Tags payments
// @Accept  json
// @Produce  json
// @Param payment body payment.Request true "Payment request"
// @Success 200 {object} payment.Response
// @Failure 400 {object} response.Object
// @Failure 500 {object} response.Object
// @Router /payments [post]
func (h *PaymentHandler) add(c *gin.Context) {
	req := payment.Request{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err, req)
		return
	}

	if err := req.Validate(); err != nil {
		response.BadRequest(c, err, req)
		return
	}

	epayReq := epay.PaymentRequest{
		Amount:     *req.Amount,
		Currency:   "KZT",
		InvoiceID:  *req.OrderID,
		TerminalID: "67e34d63-102f-4bd1-898e-370781d0074d",
	}

	token, err := h.epayService.GetToken(c, &epayReq)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	payRes, err := h.epayService.Pay(c, token, &epayReq)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}
	switch payRes.Status {
	case "NEW", "AUTH", "EXPIRED":
		req.Status = helpers.GetStringPtr("pending")
	case "CHARGE":
		req.Status = helpers.GetStringPtr("successful")
	default:
		req.Status = helpers.GetStringPtr("unsuccessful")
	}

	res, err := h.epayService.CreatePayment(c, req)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}
	response.OK(c, res)
}

// get godoc
// @Summary Get a payment
// @Description Get payment by ID
// @Tags payments
// @Accept  json
// @Produce  json
// @Param id path string true "Payment ID"
// @Success 200 {object} payment.Response
// @Failure 404 {object} response.Object
// @Failure 500 {object} response.Object
// @Router /payments/{id} [get]
func (h *PaymentHandler) get(c *gin.Context) {
	id := c.Param("id")

	res, err := h.epayService.GetPayment(c, id)
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
// @Summary Update a payment
// @Description Update payment by ID
// @Tags payments
// @Accept  json
// @Produce  json
// @Param id path string true "Payment ID"
// @Param payment body payment.Request true "Payment request"
// @Success 200 {string} string "ok"
// @Failure 400 {object} response.Object
// @Failure 404 {object} response.Object
// @Failure 500 {object} response.Object
// @Router /payments/{id} [put]
func (h *PaymentHandler) update(c *gin.Context) {
	id := c.Param("id")
	req := payment.Request{}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err, req)
		return
	}

	if err := req.IsEmpty("update"); err != nil {
		response.BadRequest(c, err, req)
		return
	}

	if err := h.epayService.UpdatePayment(c, id, req); err != nil {
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
// @Summary Delete a payment
// @Description Delete payment by ID
// @Tags payments
// @Accept  json
// @Produce  json
// @Param id path string true "Payment ID"
// @Success 200 {string} string "Payment deleted"
// @Failure 404 {object} response.Object
// @Failure 500 {object} response.Object
// @Router /payments/{id} [delete]
func (h *PaymentHandler) delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.epayService.DeletePayment(c, id); err != nil {
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
// @Summary Search payments
// @Description Search payments by user ID or order ID
// @Tags payments
// @Accept  json
// @Produce  json
// @Param userId query string false "User ID"
// @Param orderId query string false "Order ID"
// @Param status query string false "Status"
// @Success 200 {array} payment.Response
// @Failure 400 {object} response.Object
// @Failure 500 {object} response.Object
// @Router /payments/search [get]
func (h *PaymentHandler) search(c *gin.Context) {
	req := payment.Request{
		UserID:  helpers.GetStringPtr(c.Query("userId")),
		OrderID: helpers.GetStringPtr(c.Query("orderId")),
		Status:  helpers.GetStringPtr(c.Query("status")),
	}

	if err := req.IsEmpty("search"); err != nil {
		response.BadRequest(c, errors.New("incorrect query"), nil)
		return
	}

	res, err := h.epayService.SearchPayment(c, req)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.OK(c, res)
}

//func (h *PaymentHandler) getToken(c *gin.Context) {
//	req := epay.PaymentRequest{}
//	if err := c.ShouldBindJSON(&req); err != nil {
//		response.BadRequest(c, err, req)
//		return
//	}
//
//	tok, err := h.epayService.GetToken(c, &req)
//	if err != nil {
//		response.BadRequest(c, err, "err")
//	}
//
//	response.OK(c, tok)
//}
//
//func (h *PaymentHandler) getStatus(c *gin.Context) {
//	invoiceID := c.Param("id")
//	authHeader := c.GetHeader("Authorization")
//	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
//		token := authHeader[7:]
//		status, err := h.epayService.GetStatus(c, token, invoiceID)
//		if err != nil {
//			response.BadRequest(c, err, "err")
//			return
//		}
//		response.OK(c, status)
//	} else {
//		response.BadRequest(c, nil, "Invalid token")
//	}
//}

//func (h *PaymentHandler) pay(c *gin.Context) {
//	invoiceID := "19090919000919"
//	req := epay.PaymentRequest{
//		Amount:     "100",
//		Currency:   "KZT",
//		InvoiceID:  invoiceID,
//		TerminalID: "67e34d63-102f-4bd1-898e-370781d0074d",
//	}
//	token, err := h.epayService.GetToken(c, &req)
//	if err != nil {
//		response.BadRequest(c, err, req)
//	}
//	res, err := h.epayService.Pay(c, token, invoiceID)
//	if err != nil {
//		response.BadRequest(c, err, "err")
//		return
//	}
//	response.OK(c, res)
//	//authHeader := c.GetHeader("Authorization")
//	//if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
//	//	token := authHeader[7:]
//	//	status, err := h.epayService.Pay(c, token, invoiceID)
//	//	if err != nil {
//	//		response.BadRequest(c, err, "err")
//	//		return
//	//	}
//	//	response.OK(c, status)
//	//} else {
//	//	response.BadRequest(c, nil, "Invalid token")
//	//}
//}
