package service

import (
	"os"
	"strconv"

	"github.com/doddeeph/billing-engine/internal/billing/dto"
	"github.com/doddeeph/billing-engine/internal/billing/model"
	"github.com/doddeeph/billing-engine/internal/billing/repository"
	"gorm.io/gorm"
)

type BillingService interface {
	WithTransaction(tx *gorm.DB) BillingService
	CreateBilling(req dto.CreateBillingRequest) (*model.Billing, error)
	GetBilling(id string) (*model.Billing, error)
	GetOutstanding(id string) (int, error)
	IsDelinquent(id string) (bool, error)
	FindByCustomerIdAndLoanId(customerID, loanID uint, includePayments bool) (*model.Billing, error)
	UpdateOutstandingBalance(billingID uint, balance int) error
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

func (svc *billingServiceImpl) CreateBilling(req dto.CreateBillingRequest) (*model.Billing, error) {
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
		CustomerID:         req.CustomerID,
		LoanID:             req.LoanID,
		LoanAmount:         req.LoanAmount,
		LoanWeeks:          req.LoanWeeks,
		LoanInterest:       req.LoanInterest,
		OutstandingBalance: outstandingBalance,
		Payments:           payments,
	}
	if err := svc.repo.Create(billing); err != nil {
		return nil, err
	}
	return billing, nil
}

func (svc *billingServiceImpl) GetBilling(id string) (*model.Billing, error) {
	return svc.repo.FindByID(convertStringToUint(id))
}

func (svc *billingServiceImpl) GetOutstanding(id string) (int, error) {
	billing, err := svc.repo.FindByID(convertStringToUint(id))
	if err != nil {
		return 0, err
	}
	return billing.OutstandingBalance, nil
}

func (svc *billingServiceImpl) IsDelinquent(id string) (bool, error) {
	billing, err := svc.repo.FindByID(convertStringToUint(id))
	if err != nil {
		return false, err
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
	return missed >= svc.missedPaymentMax, nil
}

func (svc *billingServiceImpl) FindByCustomerIdAndLoanId(customerID uint, loanID uint, includePayments bool) (*model.Billing, error) {
	return svc.repo.FindByCustomerIdAndLoanId(customerID, loanID, includePayments)
}

func (svc *billingServiceImpl) UpdateOutstandingBalance(billingID uint, balance int) error {
	return svc.repo.UpdateOutstandingBalance(billingID, balance)
}

func convertStringToUint(s string) uint {
	billingID, _ := strconv.ParseUint(s, 10, 32)
	return uint(billingID)
}
