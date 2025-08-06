package handler

import (
	"net/http"

	"github.com/doddeeph/billing-engine/internal/billing/dto"
	"github.com/doddeeph/billing-engine/internal/billing/service"
	"github.com/doddeeph/billing-engine/utils"
	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	svc service.PaymentService
}

func NewPaymentHandler(svc service.PaymentService) *PaymentHandler {
	return &PaymentHandler{svc: svc}
}

func (h *PaymentHandler) RegisterRoutes(rg *gin.RouterGroup) {
	payment := rg.Group("/billings/:id/payments")
	// POST /billings/:id/payments
	payment.POST("", h.MakePayment)
}

func (h *PaymentHandler) MakePayment(c *gin.Context) {
	id := c.Param("id")
	billingID, err := utils.ConvertStringToUint(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var req dto.PaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	paymentResp, err := h.svc.MakePayment(c.Request.Context(), billingID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, paymentResp)
}
