package billing

import (
	"fmt"
	"log"

	"github.com/doddeeph/billing-engine/internal/billing/handler"
	"github.com/doddeeph/billing-engine/internal/billing/repository"
	"github.com/doddeeph/billing-engine/internal/billing/service"
	"github.com/doddeeph/billing-engine/pkg/billing/config"
	"github.com/gin-gonic/gin"
)

type BillingApp struct {
	AppPort        string
	BillingHandler *handler.BillingHandler
	PaymentHandler *handler.PaymentHandler
}

func NewBillingApp() *BillingApp {
	config.Load()
	appPort := config.GetAppPort()
	db := config.InitBillingDB()

	billingRepo := repository.NewBillingRepository(db)
	billingSvc := service.NewBillingService(billingRepo)
	billingHandler := handler.NewBillingHandler(billingSvc)

	paymentRepo := repository.NewPaymentRepository(db)
	paymentSvc := service.NewPaymentService(paymentRepo, billingSvc)
	paymentHandler := handler.NewPaymentHandler(paymentSvc)

	return &BillingApp{
		AppPort:        fmt.Sprintf(":%s", appPort),
		BillingHandler: billingHandler,
		PaymentHandler: paymentHandler,
	}
}

func (app *BillingApp) Start() {
	r := gin.Default()

	apiV1 := r.Group("/api/v1")
	app.BillingHandler.RegisterRoutes(apiV1)
	app.PaymentHandler.RegisterRoutes(apiV1)

	if err := r.Run(app.AppPort); err != nil {
		log.Fatalf("Failed to run Billing Engine: %v", err)
	}
	log.Printf("Billing Engine started at %s", app.AppPort)
}
