package repository

import (
	"github.com/doddeeph/billing-engine/internal/billing/model"
	"gorm.io/gorm"
)

type BillingRepository interface {
	Create(billing *model.Billing) error
	FindByCustomerIdAndLoanId(customerID, loanID uint) (*model.Billing, error)
}

type billingRepository struct {
	db *gorm.DB
}

func NewBillingRepository(db *gorm.DB) BillingRepository {
	return &billingRepository{db}
}

func (r *billingRepository) Create(billing *model.Billing) error {
	return r.db.Create(billing).Error
}

func (r *billingRepository) FindByCustomerIdAndLoanId(customerID, loanID uint) (*model.Billing, error) {
	var billing model.Billing
	err := r.db.Where("customer_id = ? AND loan_id = ?", customerID, loanID).First(&billing).Error
	if err != nil {
		return nil, err
	}
	return &billing, nil
}
