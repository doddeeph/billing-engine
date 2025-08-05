package service

import (
	"sort"

	"github.com/doddeeph/billing-engine/internal/billing/dto"
	"github.com/doddeeph/billing-engine/internal/billing/model"
	"github.com/doddeeph/billing-engine/internal/billing/repository"
	"gorm.io/gorm"
)

type BillingService interface {
	WithTransaction(tx *gorm.DB) BillingService
	CreateBilling(req dto.CreateBillingRequest) (*model.Billing, error)
	GetOutstandingBalance(req dto.GetOutstandingRequest) (int, error)
	IsDelinquent(req dto.IsDelinquentRequest) (bool, error)
	FindByCustomerIdAndLoanId(customerID, loanID uint, includePayments bool) (*model.Billing, error)
	UpdateOutstandingBalance(billingID uint, balance int) error
}

type billingService struct {
	repo repository.BillingRepository
}

func NewBillingService(repo repository.BillingRepository) BillingService {
	return &billingService{repo}
}

func (svc *billingService) WithTransaction(tx *gorm.DB) BillingService {
	return &billingService{svc.repo.WithTransaction(tx)}
}

func (svc *billingService) CreateBilling(req dto.CreateBillingRequest) (*model.Billing, error) {
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

func (svc *billingService) GetOutstandingBalance(req dto.GetOutstandingRequest) (int, error) {
	billing, err := svc.repo.FindByCustomerIdAndLoanId(req.CustomerID, req.LoanID, false)
	if err != nil {
		return 0, err
	}
	return billing.OutstandingBalance, nil
}

func (svc *billingService) IsDelinquent(req dto.IsDelinquentRequest) (bool, error) {
	billing, err := svc.repo.FindByCustomerIdAndLoanId(req.CustomerID, req.LoanID, true)
	if err != nil {
		return false, err
	}
	payments := billing.Payments
	sort.Slice(payments, func(i, j int) bool {
		return payments[i].Week < payments[j].Week
	})
	missed := 0
	for i, p := range payments {
		if p.Week == i && !p.Paid {
			missed++
		}
		if missed > 2 {
			break
		}
	}
	return missed > 2, nil
}

func (svc *billingService) FindByCustomerIdAndLoanId(customerID uint, loanID uint, includePayments bool) (*model.Billing, error) {
	return svc.repo.FindByCustomerIdAndLoanId(customerID, loanID, includePayments)
}

func (svc *billingService) UpdateOutstandingBalance(billingID uint, balance int) error {
	return svc.repo.UpdateOutstandingBalance(billingID, balance)
}
