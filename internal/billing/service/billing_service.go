package service

import (
	"context"
	"os"
	"strconv"

	"github.com/doddeeph/billing-engine/internal/billing/dto"
	"github.com/doddeeph/billing-engine/internal/billing/model"
	"github.com/doddeeph/billing-engine/internal/billing/repository"
	"gorm.io/gorm"
)

type BillingService interface {
	WithTransaction(tx *gorm.DB) BillingService
	CreateBilling(ctx context.Context, req dto.CreateBillingRequest) (*model.Billing, error)
	GetBilling(ctx context.Context, id uint) (*model.Billing, error)
	GetOutstanding(ctx context.Context, id uint) (*dto.OutstandingResponse, error)
	IsDelinquent(ctx context.Context, id uint) (*dto.DelinquentResponse, error)
	FindByCustomerIdAndLoanId(ctx context.Context, customerID, loanID uint, includePayments bool) (*model.Billing, error)
	UpdateOutstanding(ctx context.Context, billingID uint, balance int) error
}

type billingServiceImpl struct {
	repo             repository.BillingRepository
	missedPaymentMax int
}

func getMissedPaymentMax() int {
	missedPaymentMaxStr := os.Getenv("MISSED_PAYMENT_MAX")
	missedPaymentMax, err := strconv.Atoi(missedPaymentMaxStr)
	if err != nil {
		missedPaymentMax = 2
	}
	return missedPaymentMax
}

func NewBillingService(repo repository.BillingRepository) BillingService {
	return &billingServiceImpl{repo: repo, missedPaymentMax: getMissedPaymentMax()}
}

func (svc *billingServiceImpl) WithTransaction(tx *gorm.DB) BillingService {
	return &billingServiceImpl{repo: svc.repo.WithTransaction(tx), missedPaymentMax: getMissedPaymentMax()}
}

func (svc *billingServiceImpl) CreateBilling(ctx context.Context, req dto.CreateBillingRequest) (*model.Billing, error) {
	outstandingBalance := req.LoanAmount + (req.LoanAmount * req.LoanInterest / 100)
	loanWeeklyAmount := outstandingBalance / req.LoanWeeks
	payments := make([]model.Payment, req.LoanWeeks)
	for i := range payments {
		payments[i] = model.Payment{
			Amount: loanWeeklyAmount,
			Week:   i + 1,
			Paid:   false,
		}
	}
	billing := &model.Billing{
		CustomerID:   req.CustomerID,
		LoanID:       req.LoanID,
		LoanAmount:   req.LoanAmount,
		LoanWeeks:    req.LoanWeeks,
		LoanInterest: req.LoanInterest,
		Outstanding:  outstandingBalance,
		Payments:     payments,
	}
	if err := svc.repo.Create(ctx, billing); err != nil {
		return nil, err
	}
	return billing, nil
}

func (svc *billingServiceImpl) GetBilling(ctx context.Context, id uint) (*model.Billing, error) {
	return svc.repo.FindByID(ctx, id)
}

func (svc *billingServiceImpl) GetOutstanding(ctx context.Context, id uint) (*dto.OutstandingResponse, error) {
	billing, err := svc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &dto.OutstandingResponse{
		BillingID:   billing.ID,
		CustomerID:  billing.CustomerID,
		LoanID:      billing.LoanID,
		Outstanding: billing.Outstanding,
	}, nil
}

func (svc *billingServiceImpl) IsDelinquent(ctx context.Context, id uint) (*dto.DelinquentResponse, error) {
	billing, err := svc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	missed := 0
	for _, p := range billing.Payments {
		if !p.Paid {
			missed++
		} else if p.Paid {
			missed = 0
		}
		if missed >= svc.missedPaymentMax {
			break
		}
	}
	return &dto.DelinquentResponse{
		BillingID:    billing.ID,
		CustomerID:   billing.CustomerID,
		LoanID:       billing.LoanID,
		IsDelinquent: missed >= svc.missedPaymentMax,
	}, nil
}

func (svc *billingServiceImpl) FindByCustomerIdAndLoanId(ctx context.Context, customerID uint, loanID uint, includePayments bool) (*model.Billing, error) {
	return svc.repo.FindByCustomerIdAndLoanId(ctx, customerID, loanID, includePayments)
}

func (svc *billingServiceImpl) UpdateOutstanding(ctx context.Context, billingID uint, balance int) error {
	return svc.repo.UpdateOutstanding(ctx, billingID, balance)
}
