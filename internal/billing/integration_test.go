package billing

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/doddeeph/billing-engine/internal/billing/dto"
	"github.com/doddeeph/billing-engine/internal/billing/handler"
	"github.com/doddeeph/billing-engine/internal/billing/model"
	"github.com/doddeeph/billing-engine/internal/billing/repository"
	"github.com/doddeeph/billing-engine/internal/billing/service"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db         *gorm.DB
	billingSvc service.BillingService
	paymentSvc service.PaymentService
	router     *gin.Engine
)

func setupTest(t *testing.T) func() {
	t.Helper()
	ctx := context.Background()

	err := godotenv.Load("../../.env")
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	dbImage := os.Getenv("DB_TEST_IMAGE")
	dbName := os.Getenv("DB_TEST_NAME")
	dbUser := os.Getenv("DB_TEST_USER")
	dbPass := os.Getenv("DB_TEST_PASSWORD")

	req := testcontainers.ContainerRequest{
		Image:        dbImage,
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_DB":       dbName,
			"POSTGRES_USER":     dbUser,
			"POSTGRES_PASSWORD": dbPass,
		},
		WaitingFor: wait.ForListeningPort("5432/tcp").
			WithStartupTimeout(30 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})
	assert.NoError(t, err)

	host, err := container.Host(ctx)
	assert.NoError(t, err)

	port, err := container.MappedPort(ctx, "5432")
	assert.NoError(t, err)

	dsn := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		host, port.Port(), dbName, dbUser, dbPass,
	)

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	assert.NoError(t, err)

	assert.NoError(t, db.AutoMigrate(&model.Billing{}, &model.Payment{}))

	billingRepo := repository.NewBillingRepository(db)
	billingSvc = service.NewBillingService(billingRepo)
	billingHandler := handler.NewBillingHandler(billingSvc)

	paymentRepo := repository.NewPaymentRepository(db)
	paymentSvc = service.NewPaymentService(paymentRepo, billingSvc)
	paymentHandler := handler.NewPaymentHandler(paymentSvc)

	gin.SetMode(gin.TestMode)
	router = gin.Default()
	router.POST("/billings", billingHandler.CreateBilling)
	router.GET("/billings/:id", billingHandler.GetBilling)
	router.GET("/billings/:id/outstanding", billingHandler.GetOutstanding)
	router.GET("/billings/:id/delinquent", billingHandler.IsDelinquent)
	router.POST("/billings/:id/payments", paymentHandler.MakePayment)

	return func() {
		_ = container.Terminate(ctx)
	}
}

func createTestBilling(t *testing.T) *model.Billing {
	req := dto.CreateBillingRequest{
		CreateBillingDTO: dto.CreateBillingDTO{
			CustomerID:   1,
			LoanID:       1,
			LoanAmount:   5000000,
			LoanInterest: 10,
			LoanWeeks:    50,
		},
	}
	billing, err := billingSvc.CreateBilling(t.Context(), req)
	assert.NoError(t, err)
	return billing
}

func TestIntegration_CreateBilling(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	payload := dto.CreateBillingRequest{
		CreateBillingDTO: dto.CreateBillingDTO{
			CustomerID:   1,
			LoanID:       1,
			LoanAmount:   5000000,
			LoanInterest: 10,
			LoanWeeks:    50,
		},
	}
	payloadBytes, _ := json.Marshal(payload)

	r, _ := http.NewRequest("POST", "/billings", bytes.NewBuffer(payloadBytes))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	assert.Equal(t, 201, w.Code)

	var resp dto.CreateBillingResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NotZero(t, resp.BillingID)
	assert.NotZero(t, resp.CustomerID)
	assert.NotZero(t, resp.LoanID)
	assert.Equal(t, 5000000, resp.LoanAmount)
	assert.Equal(t, 10, resp.LoanInterest)
	assert.Equal(t, 50, resp.LoanWeeks)
	assert.Equal(t, 5500000, resp.Outstanding)
}

func TestIntegration_GetBilling(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	billing := createTestBilling(t)
	assert.NotZero(t, billing.ID)
	assert.NotZero(t, billing.CustomerID)
	assert.NotZero(t, billing.LoanID)

	r, _ := http.NewRequest("GET", fmt.Sprintf("/billings/%d", billing.ID), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	assert.Equal(t, 200, w.Code)

	var resp model.Billing
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NotZero(t, resp.ID)
	assert.NotZero(t, resp.CustomerID)
	assert.NotZero(t, resp.LoanID)
	assert.Equal(t, 5000000, resp.LoanAmount)
	assert.Equal(t, 10, resp.LoanInterest)
	assert.Equal(t, 50, resp.LoanWeeks)
	assert.Equal(t, 5500000, resp.Outstanding)
	assert.Len(t, resp.Payments, 50)
	assert.Equal(t, 1, resp.Payments[0].Week)
	assert.Equal(t, 110000, resp.Payments[0].Amount)
	assert.False(t, resp.Payments[0].Paid)
	assert.Equal(t, 50, resp.Payments[49].Week)
	assert.Equal(t, 110000, resp.Payments[49].Amount)
	assert.False(t, resp.Payments[49].Paid)
}

func TestIntegration_GetOutstanding(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	billing := createTestBilling(t)
	assert.NotZero(t, billing.ID)
	assert.NotZero(t, billing.CustomerID)
	assert.NotZero(t, billing.LoanID)

	r, _ := http.NewRequest("GET", fmt.Sprintf("/billings/%d/outstanding", billing.ID), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	assert.Equal(t, 200, w.Code)

	var resp dto.OutstandingResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, 5500000, resp.Outstanding)
}

func TestIntegration_IsDelinquent(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	billing := createTestBilling(t)
	assert.NotZero(t, billing.ID)
	assert.NotZero(t, billing.CustomerID)
	assert.NotZero(t, billing.LoanID)

	paymentReq := dto.PaymentRequest{
		Week: 3,
	}
	_, err := paymentSvc.MakePayment(t.Context(), billing.ID, paymentReq)
	assert.NoError(t, err)

	r, _ := http.NewRequest("GET", fmt.Sprintf("/billings/%d/delinquent", billing.ID), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	assert.Equal(t, 200, w.Code)

	var resp dto.DelinquentResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.True(t, resp.IsDelinquent)
}

func TestIntregration_MakePayment(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	billing := createTestBilling(t)
	assert.NotZero(t, billing.ID)
	assert.NotZero(t, billing.CustomerID)
	assert.NotZero(t, billing.LoanID)

	payload := dto.PaymentRequest{
		Week: 1,
	}
	payloadBytes, _ := json.Marshal(payload)

	r, _ := http.NewRequest("POST", fmt.Sprintf("/billings/%d/payments", billing.ID), bytes.NewBuffer(payloadBytes))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	assert.Equal(t, 200, w.Code)

	var paymentResp dto.PaymentResponse
	json.Unmarshal(w.Body.Bytes(), &paymentResp)
	assert.Equal(t, 5390000, paymentResp.Outstanding)
	assert.True(t, paymentResp.Payment.Paid)
}
