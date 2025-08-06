package handler

import (
	"net/http"

	"github.com/doddeeph/billing-engine/internal/billing/dto"
	"github.com/doddeeph/billing-engine/internal/billing/service"
	"github.com/doddeeph/billing-engine/utils"
	"github.com/gin-gonic/gin"
)

type BillingHandler struct {
	svc service.BillingService
}

func NewBillingHandler(svc service.BillingService) *BillingHandler {
	return &BillingHandler{svc: svc}
}

func (h *BillingHandler) RegisterRoutes(rg *gin.RouterGroup) {
	billing := rg.Group("/billings")
	// POST /billings
	billing.POST("", h.CreateBilling)
	// GET /billings/1
	billing.GET("/:id", h.GetBilling)
	// GET /billings/1/outstanding
	billing.GET("/:id/outstanding", h.GetOutstanding)
	// GET /billings/1/delinquent
	billing.GET("/:id/delinquent", h.IsDelinquent)
}

func (h *BillingHandler) CreateBilling(c *gin.Context) {
	var req dto.CreateBillingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	billing, err := h.svc.CreateBilling(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, billing)
}

func (h *BillingHandler) GetBilling(c *gin.Context) {
	id := c.Param("id")
	billingID, err := utils.ConvertStringToUint(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	billing, err := h.svc.GetBilling(billingID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, billing)
}

func (h *BillingHandler) GetOutstanding(c *gin.Context) {
	id := c.Param("id")
	billingID, err := utils.ConvertStringToUint(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	outstandingResp, err := h.svc.GetOutstanding(billingID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, outstandingResp)
}

func (h *BillingHandler) IsDelinquent(c *gin.Context) {
	id := c.Param("id")
	billingID, err := utils.ConvertStringToUint(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	delinquentResp, err := h.svc.IsDelinquent(billingID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, delinquentResp)
}
