package handler

import (
	"net/http"

	"github.com/doddeeph/billing-engine/internal/billing/dto"
	"github.com/doddeeph/billing-engine/internal/billing/service"
	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	svc service.PaymentService
}

func NewPaymentHandler(svc service.PaymentService) *PaymentHandler {
	return &PaymentHandler{svc: svc}
}

func (h *PaymentHandler) RegisterRoutes(rg *gin.RouterGroup) {
	payment := rg.Group("/payments")
	// POST /payments
	payment.POST("", h.MakePayment)
}

func (h *PaymentHandler) MakePayment(c *gin.Context) {
	var req dto.PaymetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	err := h.svc.MakePayment(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, nil)
}
