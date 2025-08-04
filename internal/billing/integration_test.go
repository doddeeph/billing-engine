package billing

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/doddeeph/billing-engine/internal/billing/dto"
	"github.com/doddeeph/billing-engine/internal/billing/model"
	"github.com/doddeeph/billing-engine/internal/billing/repository"
	"github.com/doddeeph/billing-engine/internal/billing/service"
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
)

func setupTestDB(t *testing.T) func() {
	t.Helper()
	ctx := context.Background()

	err := godotenv.Load("../../.env")
	if err != nil {
		t.Log(".env Not Found.")
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

	return func() {
		_ = container.Terminate(ctx)
	}
}

func createTestBilling(t *testing.T) *model.Billing {
	req := dto.CreateBillingRequest{
		CustomerID:   1,
		LoanID:       1,
		LoanAmount:   5000000,
		LoanInterest: 10,
		LoanWeeks:    50,
	}
	billing, err := billingSvc.CreateBilling(req)
	assert.NoError(t, err)
	return billing
}

func TestIntegration_CreateBilling(t *testing.T) {
	teardown := setupTestDB(t)
	defer teardown()

	billing := createTestBilling(t)
	assert.NotZero(t, billing.CustomerID)
	assert.NotZero(t, billing.LoanID)
	assert.Equal(t, 5000000, billing.LoanAmount)
	assert.Equal(t, 10, billing.LoanInterest)
	assert.Equal(t, 50, billing.LoanWeeks)
	assert.Equal(t, 5500000, billing.OutstandingBalance)
	assert.Equal(t, 110000, billing.LoanWeeklyAmount)
	assert.Equal(t, false, billing.IsDelinquent)
}

func TestIntegration_GetOutstandingBalance(t *testing.T) {
	teardown := setupTestDB(t)
	defer teardown()

	billing := createTestBilling(t)
	assert.NotZero(t, billing.CustomerID)
	assert.NotZero(t, billing.LoanID)

	req := dto.GetOutstandingRequest{
		CustomerID: 1,
		LoanID:     1,
	}
	outstandingBalance, err := billingSvc.GetOutstandingBalance(req)
	assert.NoError(t, err)
	assert.Equal(t, 5500000, outstandingBalance)
}

func TestIntegration_IsDelinquent(t *testing.T) {
	teardown := setupTestDB(t)
	defer teardown()

	billing := createTestBilling(t)
	assert.NotZero(t, billing.CustomerID)
	assert.NotZero(t, billing.LoanID)

	req := dto.IsDelinquentRequest{
		CustomerID: 1,
		LoanID:     1,
	}
	isDelinquent, err := billingSvc.IsDelinquent(req)
	assert.NoError(t, err)
	assert.False(t, isDelinquent)
}
