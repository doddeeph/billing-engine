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

func setupDB(t *testing.T) *gorm.DB {
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
			WithStartupTimeout(60 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})
	assert.NoError(t, err)

	t.Cleanup(func() {
		_ = container.Terminate(ctx)
	})

	host, err := container.Host(ctx)
	assert.NoError(t, err)

	port, err := container.MappedPort(ctx, "5432")
	assert.NoError(t, err)

	dsn := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		host, port.Port(), dbName, dbUser, dbPass,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	assert.NoError(t, err)

	return db
}

func TestIntegration_CreateBilling(t *testing.T) {
	db := setupDB(t)
	db.AutoMigrate(&model.Billing{}, &model.Payment{})

	repo := repository.NewBillingRepository(db)
	svc := service.NewBillingService(repo)

	req := dto.CreateBillingRequest{
		CustomerID:   1,
		LoanID:       1,
		LoanAmount:   5000000,
		LoanInterest: 10,
		LoanWeeks:    50,
	}

	billing, err := svc.CreateBilling(req)
	assert.NoError(t, err)

	assert.Equal(t, uint(1), billing.CustomerID)
	assert.Equal(t, uint(1), billing.LoanID)
	assert.Equal(t, 5000000, billing.LoanAmount)
	assert.Equal(t, 10, billing.LoanInterest)
	assert.Equal(t, 50, billing.LoanWeeks)
	assert.Equal(t, 5500000, billing.OutstandingBalance)
	assert.Equal(t, 110000, billing.LoanWeeklyAmount)
	assert.Equal(t, false, billing.IsDelinquent)
}
