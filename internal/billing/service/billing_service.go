package service

import (
	"github.com/doddeeph/billing-engine/internal/billing/dto"
	"github.com/doddeeph/billing-engine/internal/billing/model"
	"github.com/doddeeph/billing-engine/internal/billing/repository"
)

type BillingService interface {
	CreateBilling(req dto.CreateBillingRequest) (*model.Billing, error)
	GetOutstandingBalance(customerID, loanID uint) (int, error)
	IsDelinquent(customerID, loanID uint) (bool, error)
}

type billingService struct {
	repo repository.BillingRepository
}

func NewBillingService(repo repository.BillingRepository) BillingService {
	return &billingService{repo}
}

func (svc *billingService) CreateBilling(req dto.CreateBillingRequest) (*model.Billing, error) {
	outstandingBalance := req.LoanAmount + (req.LoanAmount * req.LoanInterest / 100)
	loanWeeklyAmount := outstandingBalance / req.LoanWeeks
	billing := &model.Billing{
		CustomerID:         req.CustomerID,
		LoanID:             req.LoanID,
		LoanAmount:         req.LoanAmount,
		LoanWeeks:          req.LoanWeeks,
		LoanInterest:       req.LoanInterest,
		OutstandingBalance: outstandingBalance,
		LoanWeeklyAmount:   loanWeeklyAmount,
		IsDelinquent:       false,
	}
	if err := svc.repo.Create(billing); err != nil {
		return nil, err
	}
	return billing, nil
}

func (svc *billingService) GetOutstandingBalance(customerID uint, loanID uint) (int, error) {
	billing, err := svc.repo.FindByCustomerIdAndLoanId(customerID, loanID)
	if err != nil {
		return 0, err
	}
	return billing.OutstandingBalance, nil
}

func (svc *billingService) IsDelinquent(customerID uint, loanID uint) (bool, error) {
	billing, err := svc.repo.FindByCustomerIdAndLoanId(customerID, loanID)
	if err != nil {
		return false, err
	}
	return billing.IsDelinquent, nil
}
